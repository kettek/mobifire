package join

import (
	"fmt"
	"time"

	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/handshake"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// State provides an initial joining to a server.
type State struct {
	messages.MessageHandler
	container *fyne.Container
	Hostname  string
	Port      int
	conn      *net.Connection
}

// Enter attempts a connection to the server and either continues to handshake state or shows an error and returns to the metaserver.
func (s *State) Enter(next func(states.State)) (leave func()) {
	label := widget.NewLabel("Joining " + s.Hostname + ":" + fmt.Sprint(s.Port) + "...")

	serverName := s.Hostname
	if s.Port != 0 {
		serverName += ":" + fmt.Sprint(s.Port)
	}
	s.container = container.New(layout.NewCenterLayout(), label)

	s.conn = &net.Connection{}

	go func() {
		if err := s.conn.Join(serverName); err != nil {
			label.SetText("Failed to join " + serverName + ": " + err.Error())
			time.AfterFunc(3*time.Second, func() {
				next(nil)
			})
		} else {
			s.conn.SetMessageHandler(nil) // Set to nil to ensure any messages are queued.
			next(handshake.NewState(s.conn))
		}
	}()

	return nil
}

// Container returns the container.
func (s *State) Container() *fyne.Container {
	return s.container
}
