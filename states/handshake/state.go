package handshake

import (
	"fmt"

	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/login"

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
	conn      *net.Connection
}

func NewState(conn *net.Connection) *State {
	return &State{
		conn: conn,
	}
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

	// Setup receive just sends to actual login.
	s.Once(&messages.MessageSetup{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		fmt.Println("got setup message!", m.(*messages.MessageSetup), failure)
		// FIXME: Uh... do we have to handle for MessageSetup fail?
		next(login.NewState(s.conn))
	})

	s.Once(&messages.MessageVersion{}, &messages.MessageVersion{}, func(m messages.Message, failure *messages.MessageFailure) {
		msg, ok := m.(*messages.MessageVersion)
		if !ok {
			fmt.Println("not a version message...")
			next(nil)
			return
		}
		if msg.SVVersion != "1030" {
			fmt.Println("Server version is not 1030")
			next(nil)
			return
		}
		fmt.Println("got version!", msg)
		if err := s.conn.Send(&messages.MessageVersion{CLVersion: "1030", SVName: "mobilefire"}); err != nil {
			fmt.Println("Failed to send version message:", err)
			next(nil)
			return
		}
		// FIXME: This isn't optimized, as I'm working relative to termfire.
		if err := s.conn.Send(&messages.MessageSetup{
			FaceCache: struct {
				Use   bool
				Value bool
			}{Use: true, Value: false}, // Changed to false so I can get _all_ the delicious PNGs.
			LoginMethod: struct {
				Use   bool
				Value string
			}{Use: true, Value: "2"},
			ExtendedStats: struct {
				Use   bool
				Value bool
			}{Use: true, Value: true},
			Sound2: struct {
				Use   bool
				Value uint8
			}{Use: true, Value: 1},
		}); err != nil {
			fmt.Println("Failed to send setup message:", err)
			next(nil)
			return
		}
		fmt.Println("...ok?")
	})

	// TODO: timeout? maybe from s.conn?

	label := widget.NewLabel("handshaking...")

	s.container = container.New(layout.NewCenterLayout(), label)

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}
