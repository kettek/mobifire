package join

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/handshake"

	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// State provides an initial joining to a server.
type State struct {
	messages.MessageHandler
	Hostname string
	Port     int
	conn     *net.Connection
}

// Enter attempts a connection to the server and either continues to handshake state or shows an error and returns to the metaserver.
func (s *State) Enter(next func(states.State)) (leave func()) {
	// Show joining info/label.

	serverName := s.Hostname
	if s.Port != 0 {
		serverName += ":" + fmt.Sprint(s.Port)
	}

	s.conn = &net.Connection{}
	// Set an OnLoss handler to boot back to the top state on failure... probably should show an error...
	s.conn.OnLoss = func(err error) {
		next(nil)
	}

	go func() {
		if err := s.conn.Join(serverName); err != nil {
			// TODO: Show error label.
			time.AfterFunc(3*time.Second, func() {
				next(nil)
			})
		} else {
			s.conn.SetMessageHandler(nil) // Set to nil to ensure any messages are queued.
			// TODO: Bump to handshaking.
			next(handshake.NewState(s.conn))
		}
	}()

	return nil
}

func (s *State) Update() error {
	return nil
}

func (s *State) Draw(screen *ebiten.Image) {
	// TODO
}
