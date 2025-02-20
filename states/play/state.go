package play

import (
	"errors"
	"fmt"
	"image/color"

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
	"github.com/kettek/mobifire/states/play/managers"
	"github.com/kettek/mobifire/states/play/managers/action"
	"github.com/kettek/mobifire/states/play/managers/board"
	"github.com/kettek/mobifire/states/play/managers/face"
	"github.com/kettek/mobifire/states/play/managers/items"
	"github.com/kettek/mobifire/states/play/managers/skills"
	"github.com/kettek/mobifire/states/play/managers/spells"
	"github.com/kettek/termfire/messages"
)

// State provides the actual play state of the game.
type State struct {
	messages.MessageHandler
	app             fyne.App
	window          fyne.Window
	container       *fyne.Container
	commandsManager commandsManager
	character       string
	conn            *net.Connection
	messages        []messages.MessageDrawExtInfo
	// To be moved to a character-specific location.
	sayOptions []string
	//
	playerTag int32
	//
	managers managers.Managers
}

// NewState creates a new State from a connection and a desired character to play as.
func NewState(conn *net.Connection, character string) *State {
	state := &State{
		conn:      conn,
		character: character,
		commandsManager: commandsManager{
			conn: conn,
		},
		sayOptions: []string{"hi", "yes", "no"},
	}

	state.managers.Add(face.NewManager())
	state.managers.Add(board.NewManager())
	state.managers.Add(skills.NewManager())
	state.managers.Add(spells.NewManager())
	state.managers.Add(items.NewManager())
	state.managers.Add(action.NewManager())
	return state
}

// Enter sets up all the necessary UI and network handling.
func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)
	s.managers.SetupAccessors(s.app, s.window, s.conn, &s.MessageHandler)

	s.managers.PreInit()

	// Now actually try to join the world.
	s.conn.Send(&messages.MessageAccountPlay{Character: s.character})
	// It's a little silly, but we have to handle character select failure here, as Crossfire's protocol is all over the place with state confirmations.
	s.On(&messages.MessageAccountPlay{}, &messages.MessageAccountPlay{}, func(m messages.Message, failure *messages.MessageFailure) {
		err := dialog.NewError(errors.New(failure.Reason), s.window)
		err.SetOnClosed(func() {
			s.conn.SetMessageHandler(nil)
			next(states.Prior)
		})
		err.Show()
	})

	s.managers.Init()

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

	// Leave game handling.
	s.On(&messages.MessagePlayer{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessagePlayer)
		// I believe all blank means the player used a bed to reality.
		if msg.Name == "" {
			s.conn.SetMessageHandler(nil)
			next(states.Prior)
		} else {
			s.playerTag = msg.Tag
		}
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
	//messagesListBackground := canvas.NewVerticalGradient(color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 0})
	messagesListBackground := canvas.NewRectangle(color.NRGBA{255, 0, 0, 255})

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

	// Right-hand toolbar stuff
	var toolbar *Toolbar
	{
		commandsPopup := widget.NewPopUpMenu(fyne.NewMenu("Commands", s.commandsManager.toMenuItems()...), s.window.Canvas())
		// TODO: Make our own custom hotkey sort of thing.
		var toolbarCmdAction *widget.ToolbarAction
		var toolbarApplyAction *widget.ToolbarAction
		var toolbarGetAction *widget.ToolbarAction
		toolbarCmdAction = widget.NewToolbarAction(data.GetResource("icon_commands.png"), func() {
			commandsPopup.ShowAtRelativePosition(fyne.NewPos(-toolbarCmdAction.ToolbarObject().Size().Width, 0), toolbarCmdAction.ToolbarObject())
		})
		toolbarApplyAction = widget.NewToolbarAction(data.GetResource("icon_apply.png"), func() {
			s.conn.SendCommand("apply", 0)
		})
		toolbarGetAction = widget.NewToolbarAction(data.GetResource("icon_pickup.png"), func() {
			s.conn.SendCommand("get", 0)
		})
		toolbar = NewToolbar(
			toolbarCmdAction,
			toolbarApplyAction,
			toolbarGetAction,
			widget.NewToolbarAction(data.GetResource("icon_inventory.png"), func() {
				im := s.managers.GetByType(&items.Manager{}).(*items.Manager)
				im.ShowInventory(s.playerTag, func(item *items.Item) bool {
					return false
				})
			}),
			widget.NewToolbarAction(data.GetResource("icon_inventory.png"), func() {
				sm := s.managers.GetByType(&skills.Manager{}).(*skills.Manager)
				sm.ShowSkillsList()
				fmt.Println("Toolbar action 5")
			}),
			widget.NewToolbarAction(data.GetResource("icon_inventory.png"), func() {
				sm := s.managers.GetByType(&spells.Manager{}).(*spells.Manager)
				sm.ShowSpellsList(nil)
				fmt.Println("Toolbar action 6")
			}),
		)
	}

	sizedTheme := layouts.Theme{}

	toolbarSized := container.NewThemeOverride(toolbar, sizedTheme)
	toolbars := container.NewHBox(layout.NewSpacer(), toolbarSized)

	thumbPad := &thumbpadWidget{}
	thumbPad.onCommand = func(cmd string) {
		s.conn.SendCommand(cmd, 0)
		s.managers.GetByType(&action.Manager{}).(*action.Manager).SetDirectionFromString(cmd)
	}
	thumbPadContainer := container.New(layout.NewStackLayout(), thumbPad)

	leftAreaToolbarTop := container.NewThemeOverride(container.New(layout.NewGridLayout(4),
		s.managers.GetByType(&action.Manager{}).(*action.Manager).AcquireButton(0),
		s.managers.GetByType(&action.Manager{}).(*action.Manager).AcquireButton(1),
		s.managers.GetByType(&action.Manager{}).(*action.Manager).AcquireButton(2),
		s.managers.GetByType(&action.Manager{}).(*action.Manager).AcquireButton(3),
	), sizedTheme)
	leftAreaToolbarBot := container.NewThemeOverride(container.New(layout.NewGridLayout(4),
		s.managers.GetByType(&action.Manager{}).(*action.Manager).AcquireButton(4),
		s.managers.GetByType(&action.Manager{}).(*action.Manager).AcquireButton(5),
		s.managers.GetByType(&action.Manager{}).(*action.Manager).AcquireButton(6),
		s.managers.GetByType(&action.Manager{}).(*action.Manager).AcquireButton(7),
	), sizedTheme)

	leftArea := container.New(&layouts.Left{}, leftAreaToolbarTop, thumbPadContainer, leftAreaToolbarBot)

	board := s.managers.GetByType(&board.Manager{}).(*board.Manager).CanvasObject()

	s.container = container.New(&layouts.Game{
		Board:    board,
		Messages: messagesList,
		Left:     leftArea,
		Right:    toolbars,
	}, board, container.NewStack(messagesListBackground, container.NewThemeOverride(messagesList, sizedTheme)), leftArea, toolbars)

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

// SetApp sets the app for preferences usage.
func (s *State) SetApp(app fyne.App) {
	s.app = app
}

// ShowTextDialog shows a near fullscreen dialog, wow.
func (s *State) ShowTextDialog(title string, content string) {
	segments := data.TextToRichTextSegments(content)

	text := widget.NewRichText(segments...)
	text.Wrapping = fyne.TextWrapWord
	cnt := layouts.NewDialog(s.window)
	dialog.ShowCustom(title, "Close", container.New(cnt, container.NewVScroll(text)), s.window)
}

// ShowTextDialogWithInput is like ShowTextDialog, but with an input entry.
func (s *State) ShowTextDialogWithInput(title string, content string, submit string, cb func(string)) {
	segments := data.TextToRichTextSegments(content)

	text := widget.NewRichText(segments...)
	text.Wrapping = fyne.TextWrapWord
	cnt := layouts.NewDialog(s.window)
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
