package login

import (
	"github.com/kettek/mobifire/net"

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
	Hostname  string
	Port      int
	conn      net.Connection
}

func NewState(conn net.Connection) *State {
	return &State{
		conn: conn,
	}
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	label := widget.NewLabel("TODO: Login")

	s.container = container.New(layout.NewCenterLayout(), label)

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}
