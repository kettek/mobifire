package play

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
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
	window        fyne.Window
	container     *fyne.Container
	mb            *multiBoard
	character     string
	conn          *net.Connection
	messages      []messages.MessageDrawExtInfo
	pendingImages []boardPendingImage
}

// NewState creates a new State from a connection and a desired character to play as.
func NewState(conn *net.Connection, character string) *State {
	return &State{
		conn:      conn,
		character: character,
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

	// Image and animation message processing.
	s.On(&messages.MessageFace2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageFace2)
		if _, ok := data.GetFace(int(msg.Num)); !ok {
			s.conn.Send(&messages.MessageAskFace{Face: uint32(msg.Num)})
		}
	})

	s.On(&messages.MessageImage2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageImage2)
		data.AddFaceImage(*msg)
		for i := len(s.pendingImages) - 1; i >= 0; i-- {
			if s.pendingImages[i].Num == uint16(msg.Face) {
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
				case *messages.MessageMap2CoordDataClear:
					s.mb.SetCells(m.X, m.Y, nil)
				case *messages.MessageMap2CoordDataClearLayer:
					s.mb.SetCell(m.X, m.Y, int(d.Layer), nil)
				case *messages.MessageMap2CoordDataImage:
					if d.FaceNum == 0 {
						s.mb.SetCell(m.X, m.Y, int(d.Layer), nil)
						continue
					}
					faceImage, ok := data.GetFace(int(d.FaceNum))
					if !ok {
						s.pendingImages = append(s.pendingImages, boardPendingImage{X: m.X, Y: m.Y, Z: int(d.Layer), Num: d.FaceNum})
						continue
					}
					s.mb.SetCell(m.X, m.Y, int(d.Layer), &faceImage)
				}
			}
		}
	})

	// Text processing.
	s.On(&messages.MessageDrawExtInfo{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageDrawExtInfo)
		s.messages = append(s.messages, *msg)
		fmt.Println("got draw ext info", msg)
	})

	messagesList := widget.NewList(
		func() int {
			return len(s.messages)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(s.messages[i].Message)
		},
	)
	messagesList.HideSeparators = true

	// Use our current face set for the board... could we make setting the faceset dynamic...??
	faceset := data.CurrentFaceSet()
	s.mb = newMultiBoard(11, 11, 8, faceset.Width, faceset.Height)

	// TODO: Make our own custom hotkey sort of thing.
	toolbar := NewToolbar(
		widget.NewToolbarAction(resourceInventoryPng, func() {
			fmt.Println("Toolbar action 1")
		}),
		widget.NewToolbarAction(resourceInventoryPng, func() {
			fmt.Println("Toolbar action 2")
		}),
		widget.NewToolbarAction(resourceInventoryPng, func() {
			fmt.Println("Toolbar action 3")
		}),
		widget.NewToolbarAction(resourceInventoryPng, func() {
			fmt.Println("Toolbar action 4")
		}),
		widget.NewToolbarAction(resourceInventoryPng, func() {
			fmt.Println("Toolbar action 5")
		}),
	)

	sizedTheme := myTheme{}

	toolbarSized := container.NewThemeOverride(toolbar, sizedTheme)
	toolbars := container.NewHBox(layout.NewSpacer(), toolbarSized)

	thumbPad := &thumbpadWidget{}
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

	s.container = container.New(&layouts.Game{}, leftArea, s.mb.container, toolbars)

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
