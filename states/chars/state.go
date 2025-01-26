package chars

import (
	"fmt"

	"github.com/kettek/mobifire/net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

type State struct {
	window fyne.Window
	messages.MessageHandler
	container  *fyne.Container
	conn       *net.Connection
	characters []messages.Character
	faces      []messages.MessageFace2
}

func NewState(conn *net.Connection, characters []messages.Character, faces []messages.MessageFace2) *State {
	return &State{
		conn:       conn,
		characters: characters,
		faces:      faces,
	}
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

	fmt.Println("we got", s.characters)

	label := widget.NewLabel("Select a character:")

	s.container = container.New(layout.NewVBoxLayout(), label)

	return nil
}

func (s *State) SetWindow(window fyne.Window) {
	s.window = window
}

func (s *State) Container() *fyne.Container {
	return s.container
}
