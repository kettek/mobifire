package login

import (
	"errors"
	"fmt"

	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/chars"
	"github.com/kettek/mobifire/states/play"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

type State struct {
	messages.MessageHandler
	window    fyne.Window
	container *fyne.Container
	conn      *net.Connection
	faces     []messages.MessageFace2
}

func NewState(conn *net.Connection) *State {
	return &State{
		conn: conn,
	}
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

	s.On(&messages.MessageAccountLogin{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		if mf != nil {
			fmt.Println("Failed to login: ", mf.Reason)
			dialog.ShowError(errors.New(mf.Reason), s.window)
			return
		}
		next(&play.State{})
	})

	s.On(&messages.MessageAccountPlayers{}, &messages.MessageAccountLogin{}, func(msg messages.Message, failure *messages.MessageFailure) {
		if failure != nil {
			dialog.ShowError(errors.New(failure.Reason), s.window)
			return
		}
		m := msg.(*messages.MessageAccountPlayers)
		next(chars.NewState(s.conn, m.Characters, s.faces))
	})

	s.On(&messages.MessageFace2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m, ok := msg.(*messages.MessageFace2)
		if !ok {
			return
		}
		s.faces = append(s.faces, *m)
	})

	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		},
		OnSubmit: func() {
			s.conn.Send(&messages.MessageAccountLogin{Account: usernameEntry.Text, Password: passwordEntry.Text})
		},
	}

	s.container = container.New(layout.NewVBoxLayout(), form)

	return nil
}

func (s *State) SetWindow(window fyne.Window) {
	s.window = window
}

func (s *State) Container() *fyne.Container {
	return s.container
}
