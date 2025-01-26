package login

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

	s.Once(&messages.MessageVersion{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
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
		s.conn.Send(&messages.MessageVersion{CLVersion: "1030", SVName: "mobilefire"})
		fmt.Println("looks good!")
		// TODO: Send setup!!!
	})

	// TODO: Handle Face2 to store, since it gets sent here... need access to a files cache.

	// TODO: Handle Setup receive, as that's a confirm to our request in sending Version -- this is where we swap to actual login credentials (TODO: Move this state to handshake???)

	label := widget.NewLabel("TODO: Login")

	s.container = container.New(layout.NewCenterLayout(), label)

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}
