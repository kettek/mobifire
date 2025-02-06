package play

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/mobifire/states/play/layouts"
	"github.com/kettek/termfire/messages"
)

// State provides the actual play state of the game.
type State struct {
	messages.MessageHandler
	window          fyne.Window
	container       *fyne.Container
	mb              *multiBoard
	commandsManager commandsManager
	character       string
	conn            *net.Connection
	messages        []messages.MessageDrawExtInfo
	pendingImages   []boardPendingImage
	// To be moved to a character-specific location.
	sayOptions []string
}

// NewState creates a new State from a connection and a desired character to play as.
func NewState(conn *net.Connection, character string) *State {
	return &State{
		conn:      conn,
		character: character,
		commandsManager: commandsManager{
			conn: conn,
		},
		sayOptions: []string{"hi", "yes", "no"},
	}
}

// Enter sets up all the necessary UI and network handling.
func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)
	s.conn.Send(&messages.MessageAccountPlay{Character: s.character})
	// It's a little silly, but we have to handle character select failure here, as Crossfire's protocol is all over the place with state confirmations.
	s.On(&messages.MessageAccountPlay{}, &messages.MessageAccountPlay{}, func(m messages.Message, failure *messages.MessageFailure) {
		err := dialog.NewError(errors.New(failure.Reason), s.window)
		err.SetOnClosed(func() {
			next(states.Prior)
		})
		err.Show()
	})

	// Setup commands to show in the commands list.
	s.commandsManager.commands = []command{
		{
			Name: "command",
			OnActivate: func() {
				s.ShowInput("Command", "Submit", func(cmd string) {
					s.conn.SendCommand(cmd, 0)
				})
			},
		},
		{
			Name: "say",
			OnActivate: func() {
				s.ShowInputWithOptions("Say", "Say", &s.sayOptions, func(cmd string) {
					s.conn.SendCommand("say "+cmd, 0)
				})
			},
		},
		{
			Name: "who",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("who", messages.MessageTypeCommand, messages.SubMessageTypeCommandWho)
			},
		},
		{
			Name: "statistics",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("statistics", messages.MessageTypeCommand, messages.SubMessageTypeCommandStatistics)
			},
		},
		{
			Name: "body",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("body", messages.MessageTypeCommand, messages.SubMessageTypeCommandBody)
			},
		},
		{
			Name: "inventory",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("inventory", messages.MessageTypeCommand, messages.SubMessageTypeCommandInventory)
			},
		},
		{
			Name: "skills",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("skills", messages.MessageTypeSkill, messages.SubMessageTypeSkillList)
			},
		},
		{
			Name: "maps",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("maps", messages.MessageTypeCommand, messages.SubMessageTypeCommandMaps)
			},
		},
		{
			Name: "hiscore",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("hiscore", messages.MessageTypeAdmin, messages.SubMessageTypeAdminHiscore)
			},
		},
		{
			Name: "news",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("news", messages.MessageTypeAdmin, messages.SubMessageTypeAdminNews)
			},
		},
		{
			Name: "rules",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("rules", messages.MessageTypeAdmin, messages.SubMessageTypeAdminRules)
			},
		},
		{
			Name: "motd",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommand("motd", messages.MessageTypeMOTD, 0)
			},
		},
		{
			Name: "help",
			OnActivate: func() {
				s.commandsManager.QuerySimpleCommandWithInput("help", messages.MessageTypeCommand, messages.SubMessageTypeCommandInfo).Repeat = true
			},
		},
		{
			Name: "title",
			OnActivate: func() {
				q := s.commandsManager.QuerySimpleCommandWithInput("title", messages.MessageTypeCommand, messages.SubMessageTypeCommandConfig)
				q.SubmitText = "Set Title"
			},
		},
	}
	s.commandsManager.OnCommandComplete = func(c *queryCommand) {
		if c.Text == "" {
			// Skip if no text was every received for this query.
			return
		}
		if c.HasInput {
			s.ShowTextDialogWithInput(c.Command, c.Text, c.SubmitText, func(cmd string) {
				if c.Repeat {
					query := *c
					query.Command = c.OriginalCommand + " " + cmd
					s.commandsManager.QueryCommand(query)
				} else {
					s.commandsManager.QueryComplexCommand(c.OriginalCommand+" "+cmd, c.OriginalCommand, c.MT, c.ST)
				}
			})
		} else {
			s.ShowTextDialog(c.OriginalCommand, c.Text)
		}
	}

	// Command response packet handling.
	s.On(&messages.MessageCommandCompleted{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageCommandCompleted)
		if s.commandsManager.checkCommandCompleted(msg) {
			return
		}
	})

	// Setup message handling.
	s.On(&messages.MessageSetup{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageSetup)
		if msg.MapSize.Use {
			parts := strings.Split(msg.MapSize.Value, "x")
			if len(parts) != 2 {
				fmt.Println("Invalid map size:", msg.MapSize.Value)
				return
			}
			rows, err := strconv.Atoi(parts[0])
			if err != nil {
				fmt.Println("Invalid map size:", msg.MapSize.Value)
				return
			}
			cols, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid map size:", msg.MapSize.Value)
			}
			s.mb.SetBoardSize(rows+1, cols+1)
		}
	})

	// Image and animation message processing.
	s.On(&messages.MessageFace2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageFace2)
		if _, ok := data.GetFace(int(msg.Num)); !ok {
			s.conn.Send(&messages.MessageAskFace{Face: int32(msg.Num)})
		}
	})

	s.On(&messages.MessageImage2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageImage2)
		data.AddFaceImage(*msg)
		for i := len(s.pendingImages) - 1; i >= 0; i-- {
			if s.pendingImages[i].Num == int16(msg.Face) {
				faceImage, _ := data.GetFace(int(msg.Face))
				s.mb.SetCell(s.pendingImages[i].X, s.pendingImages[i].Y, s.pendingImages[i].Z, &faceImage)
				s.pendingImages = append(s.pendingImages[:i], s.pendingImages[i+1:]...)
			}
		}
	})

	s.On(&messages.MessageMap2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageMap2)

		for _, m := range msg.Coords {
			if len(m.Data) == 0 {
				// TODO ???
				continue
			}
			for _, c := range m.Data {
				switch d := c.(type) {
				case messages.MessageMap2CoordDataDarkness:
					// TODO
				case messages.MessageMap2CoordDataAnim:
					// TODO
				case messages.MessageMap2CoordDataClear:
					s.mb.SetCells(m.X, m.Y, nil)
				case messages.MessageMap2CoordDataClearLayer:
					s.mb.SetCell(m.X, m.Y, int(d.Layer), nil)
				case messages.MessageMap2CoordDataImage:
					if d.FaceNum == 0 {
						s.mb.SetCell(m.X, m.Y, int(d.Layer), nil)
						continue
					}
					faceImage, ok := data.GetFace(int(d.FaceNum))
					if !ok {
						s.pendingImages = append(s.pendingImages, boardPendingImage{X: m.X, Y: m.Y, Z: int(d.Layer), Num: int16(d.FaceNum)})
						continue
					}
					s.mb.SetCell(m.X, m.Y, int(d.Layer), &faceImage)
				}
			}
		}
	})

	s.On(&messages.MessageNewMap{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		s.mb.Clear()
	})

	messagesList := widget.NewList(
		func() int {
			return len(s.messages)
		},
		func() fyne.CanvasObject {
			txt := canvas.NewText("", color.Black)
			return txt
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = s.messages[i].Message
			o.(*canvas.Text).Color = data.Color(s.messages[i].Color)
			o.Refresh() // ???
		},
	)
	messagesList.HideSeparators = true

	// Messages.
	lastVOffset := float32(0)
	s.On(&messages.MessageDrawExtInfo{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageDrawExtInfo)

		if s.commandsManager.checkDrawExtInfo(msg) {
			return
		}

		// This isn't right to do, but for now, split all messages if their columns exceed 50.
		var msgs []messages.MessageDrawExtInfo
		for i := 0; i < len(msg.Message); i += 50 {
			end := i + 50
			if end > len(msg.Message) {
				end = len(msg.Message)
				if msg.Message[end-1] == '\n' {
					end--
				}
			}
			msgs = append(msgs, messages.MessageDrawExtInfo{
				Color:   msg.Color,
				Type:    msg.Type,
				Subtype: msg.Subtype,
				Message: msg.Message[i:end],
			})
		}

		if lastVOffset == 0 {
			lastVOffset = messagesList.GetScrollOffset()
		}
		// Automatically scroll to end if user has not scrolled up.
		if messagesList.GetScrollOffset() == lastVOffset {
			s.messages = append(s.messages, msgs...)
			messagesList.Refresh()
			messagesList.ScrollToBottom()
			lastVOffset = messagesList.GetScrollOffset()
		} else {
			s.messages = append(s.messages, msgs...)
			messagesList.Refresh()
		}
	})

	// Use our current face set for the board... could we make setting the faceset dynamic...??
	faceset := data.CurrentFaceSet()
	s.mb = newMultiBoard(11, 11, 10, faceset.Width, faceset.Height)

	// Setup hooks for requesting map resizes.
	s.mb.onSizeChanged = func(rows, cols int) {
		s.conn.Send(&messages.MessageSetup{
			MapSize: struct {
				Use   bool
				Value string
			}{Use: true, Value: fmt.Sprintf("%dx%d", rows, cols)},
		})
	}

	// Right-hand toolbar stuff
	var toolbar *Toolbar
	{
		commandsPopup := widget.NewPopUpMenu(fyne.NewMenu("Commands", s.commandsManager.toMenuItems()...), s.window.Canvas())
		// TODO: Make our own custom hotkey sort of thing.
		var toolbarCmdAction *widget.ToolbarAction
		var toolbarApplyAction *widget.ToolbarAction
		var toolbarGetAction *widget.ToolbarAction
		toolbarCmdAction = widget.NewToolbarAction(resourceCommandsPng, func() {
			commandsPopup.ShowAtRelativePosition(fyne.NewPos(-toolbarCmdAction.ToolbarObject().Size().Width, 0), toolbarCmdAction.ToolbarObject())
		})
		toolbarApplyAction = widget.NewToolbarAction(resourceApplyPng, func() {
			s.conn.SendCommand("apply", 0)
		})
		toolbarGetAction = widget.NewToolbarAction(resourceGetPng, func() {
			s.conn.SendCommand("get", 0)
		})
		toolbar = NewToolbar(
			toolbarCmdAction,
			toolbarApplyAction,
			toolbarGetAction,
			widget.NewToolbarAction(resourceInventoryPng, func() {
				fmt.Println("Toolbar action 4")
			}),
			widget.NewToolbarAction(resourceInventoryPng, func() {
				fmt.Println("Toolbar action 5")
			}),
		)
	}

	sizedTheme := myTheme{}

	toolbarSized := container.NewThemeOverride(toolbar, sizedTheme)
	toolbars := container.NewHBox(layout.NewSpacer(), toolbarSized)

	thumbPad := &thumbpadWidget{}
	thumbPad.onCommand = func(cmd string) {
		s.conn.SendCommand(cmd, 0)
	}
	thumbPadContainer := container.New(layout.NewStackLayout(), thumbPad)

	leftAreaToolbarTop := container.NewThemeOverride(container.New(layout.NewGridLayout(3),
		widget.NewButtonWithIcon("", resourceInventoryPng, func() {
			fmt.Println("Toolbar action 1")
		}),
		widget.NewButtonWithIcon("", resourceInventoryPng, func() {
			fmt.Println("Toolbar action 1")
		}),
		widget.NewButtonWithIcon("", resourceInventoryPng, func() {
			fmt.Println("Toolbar action 1")
		}),
	), sizedTheme)
	leftAreaToolbarBot := container.NewThemeOverride(container.New(layout.NewGridLayout(3),
		widget.NewButtonWithIcon("", resourceInventoryPng, func() {
			fmt.Println("Toolbar action 1")
		}),
		widget.NewButtonWithIcon("", resourceInventoryPng, func() {
			fmt.Println("Toolbar action 2")
		}),
		widget.NewButtonWithIcon("", resourceInventoryPng, func() {
			fmt.Println("Toolbar action 3")
		}),
	), sizedTheme)

	leftArea := container.New(&layouts.Left{}, leftAreaToolbarTop, thumbPadContainer, leftAreaToolbarBot)

	s.container = container.New(&layouts.Game{
		Board:    s.mb.container,
		Messages: messagesList,
		Left:     leftArea,
		Right:    toolbars,
	}, s.mb.container, container.NewThemeOverride(messagesList, sizedTheme), leftArea, toolbars)

	//s.container = container.New(layout.NewCenterLayout(), vcontainer)

	return nil
}

// Container returns the container.
func (s *State) Container() *fyne.Container {
	return s.container
}

// SetWindow sets the window for dialog usage.
func (s *State) SetWindow(window fyne.Window) {
	s.window = window
}

type dialogContainer struct {
	window fyne.Window
}

func (d *dialogContainer) MinSize(objects []fyne.CanvasObject) fyne.Size {
	size := d.window.Canvas().Size()
	// Not sure if we have a flag somewhere for landscape vs. portrait, but...
	padding := float32(0)
	if size.Width > size.Height {
		padding = size.Height / 2
	} else {
		padding = size.Width / 2
	}
	return fyne.NewSize(size.Width-padding, size.Height-padding)
}

func (d *dialogContainer) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(size)
	}
}

// ShowTextDialog shows a near fullscreen dialog, wow.
func (s *State) ShowTextDialog(title string, content string) {
	segments := data.TextToRichTextSegments(content)

	text := widget.NewRichText(segments...)
	text.Wrapping = fyne.TextWrapWord
	cnt := &dialogContainer{
		window: s.window,
	}
	dialog.ShowCustom(title, "Close", container.New(cnt, container.NewVScroll(text)), s.window)
}

// ShowTextDialogWithInput is like ShowTextDialog, but with an input entry.
func (s *State) ShowTextDialogWithInput(title string, content string, submit string, cb func(string)) {
	segments := data.TextToRichTextSegments(content)

	text := widget.NewRichText(segments...)
	text.Wrapping = fyne.TextWrapWord
	cnt := &dialogContainer{
		window: s.window,
	}
	entry := widget.NewEntry()
	if submit == "" {
		submit = "Submit"
	}
	dialog.ShowCustomConfirm(title, submit, "Cancel", container.New(cnt, container.NewBorder(nil, entry, nil, nil, container.NewVScroll(text))), func(b bool) {
		if b {
			cb(entry.Text)
		}
	}, s.window)
}

func (s *State) ShowInput(title string, submit string, cb func(string)) {
	entry := widget.NewEntry()
	if submit == "" {
		submit = "Submit"
	}
	dialog.ShowForm(title, submit, "Cancel", []*widget.FormItem{
		{Text: "", Widget: entry},
	}, func(b bool) {
		if b {
			cb(entry.Text)
		}
	}, s.window)
}

func (s *State) ShowInputWithOptions(title string, submit string, opts *[]string, cb func(string)) {
	var entry *widget.SelectEntry
	var addEntry *widget.Button
	var removeEntry *widget.Button
	entry = widget.NewSelectEntry(*opts)
	entry.Resize(fyne.NewSize(200, 30))
	if submit == "" {
		submit = "Submit"
	}
	adjustButtons := func() {
		if entry.Text == "" {
			addEntry.Disable()
			removeEntry.Disable()
			return
		}
		found := false
		for _, o := range *opts {
			if o == entry.Text {
				found = true
				break
			}
		}
		if found {
			addEntry.Disable()
			removeEntry.Enable()
		} else {
			addEntry.Enable()
			removeEntry.Disable()
		}
	}
	entry.OnChanged = func(s string) {
		adjustButtons()
	}
	addEntry = widget.NewButton("+", func() {
		*opts = append(*opts, entry.Text)
		entry.SetOptions(*opts)
		adjustButtons()
		entry.Refresh()
	})
	addEntry.Disable()
	addEntry.Importance = widget.SuccessImportance
	removeEntry = widget.NewButton("-", func() {
		for i, o := range *opts {
			if o == entry.Text {
				*opts = append((*opts)[:i], (*opts)[i+1:]...)
				entry.SetOptions(*opts)
				adjustButtons()
				entry.Refresh()
				return
			}
		}
	})
	removeEntry.Disable()
	removeEntry.Importance = widget.DangerImportance
	entries := container.NewHBox(addEntry, removeEntry)
	dialog.ShowForm(title, submit, "Cancel", []*widget.FormItem{
		{Text: "", Widget: entry},
		{Text: "", Widget: entries},
	}, func(b bool) {
		if b {
			cb(entry.Text)
		}
	}, s.window)
}
