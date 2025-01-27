package play

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

type State struct {
	messages.MessageHandler
	window    fyne.Window
	container *fyne.Container
	mb        *multiBoard
	character string
	conn      *net.Connection
	messages  []messages.MessageDrawExtInfo
}

func NewState(conn *net.Connection, character string) *State {
	return &State{
		conn:      conn,
		character: character,
	}
}

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
		fmt.Println("got face", msg)
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

	s.mb = newMultiBoard(11, 11, 8)

	// TODO: Make our own custom hotkey sort of thing.
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(resourceMarkPng, func() {
			fmt.Println("Toolbar action 1")
		}),
		widget.NewToolbarAction(resourceMarkPng, func() {
			fmt.Println("Toolbar action 2")
		}),
	)

	borderContainer := container.NewBorder(nil, nil, messagesList, toolbar, s.mb.container)
	s.container = borderContainer

	//s.container = container.New(layout.NewCenterLayout(), vcontainer)

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}

func (s *State) SetWindow(window fyne.Window) {
	s.window = window
}
