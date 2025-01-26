package login

import (
	"fmt"

	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

type State struct {
	messages.MessageHandler
	container *fyne.Container
	conn      *net.Connection
}

func NewState(conn *net.Connection) *State {
	return &State{
		conn: conn,
	}
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

	// TODO: Handle Face2 to store, since it gets sent here... need access to a files cache.

	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		},
		OnSubmit: func() {
			fmt.Println("Username:", usernameEntry.Text)
			next(&play.State{})
		},
	}

	s.container = container.New(layout.NewVBoxLayout(), form)

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}
