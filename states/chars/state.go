package chars

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

	characterList := container.New(layout.NewVBoxLayout())

	for _, character := range s.characters {
		if character.Name == "" {
			// Skip the weird bogus empty char.
			continue
		}
		content := container.New(layout.NewHBoxLayout(), widget.NewLabel(character.Map), widget.NewButton("Play", func() {
			next(play.NewState(s.conn, character.Name))
		}))
		card := widget.NewCard(character.Name, fmt.Sprintf("%d %s %s", character.Level, character.Race, character.Class), content)
		characterList.Add(card)
	}

	s.container = container.New(layout.NewVBoxLayout(), characterList)

	return nil
}

func (s *State) SetWindow(window fyne.Window) {
	s.window = window
}

func (s *State) Container() *fyne.Container {
	return s.container
}
