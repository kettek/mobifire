package chars

import (
	"fmt"

	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// State provides the character selection and creation screen.
type State struct {
	window fyne.Window
	messages.MessageHandler
	container  *fyne.Container
	conn       *net.Connection
	characters []messages.Character
	faces      []messages.MessageFace2
}

// NewState provides a new State from a connection, Character, and Face messages.
func NewState(conn *net.Connection, characters []messages.Character, faces []messages.MessageFace2) *State {
	return &State{
		conn:       conn,
		characters: characters,
		faces:      faces,
	}
}

// Enter sets up the necessary state.
func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

	// Request faces sent during login.
	for _, face := range s.faces {
		s.conn.Send(&messages.MessageAskFace{Face: uint32(face.Num)})
	}

	s.On(&messages.MessageImage2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageImage2)
		data.AddFaceImage(*msg)
	})

	// Selection
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

	// Creation
	creationContainer := container.New(layout.NewVBoxLayout(), widget.NewLabel("Create a new character!"))

	// Tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Create", creationContainer),
		container.NewTabItem("Select", characterList),
	)
	if len(s.characters) > 1 {
		tabs.SelectIndex(1)
	}

	s.container = container.New(layout.NewVBoxLayout(), tabs)

	//s.container = container.New(layout.NewVBoxLayout(), characterList)

	return nil
}

// SetWindow sets the window for dialog functions.
func (s *State) SetWindow(window fyne.Window) {
	s.window = window
}

// Container returns the container.
func (s *State) Container() *fyne.Container {
	return s.container
}
