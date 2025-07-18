package join

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/handshake"
	"github.com/kettek/rebui"
	"github.com/kettek/rebui/widgets"

	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// State provides an initial joining to a server.
type State struct {
	messages.MessageHandler
	Hostname string
	Port     int
	conn     *net.Connection
	layout   rebui.Layout
}

// Enter attempts a connection to the server and either continues to handshake state or shows an error and returns to the metaserver.
func (s *State) Enter(next func(states.State)) (leave func()) {
	// Show joining info/label.
	node := s.layout.AddNode(rebui.Node{
		Type:            "Text",
		Text:            "connecting...",
		TextWrap:        rebui.WrapWord,
		VerticalAlign:   rebui.AlignMiddle,
		HorizontalAlign: rebui.AlignCenter,
		Width:           "100%",
		Height:          "100%",
		X:               "50%",
		Y:               "50%",
		OriginX:         "-50%",
		OriginY:         "-50%",
	})
	widget := node.Widget.(*widgets.Text)

	serverName := s.Hostname
	if s.Port != 0 {
		serverName += ":" + fmt.Sprint(s.Port)
	}

	s.conn = &net.Connection{}
	// Set an OnLoss handler to boot back to the top state on failure... probably should show an error...
	s.conn.OnLoss = func(err error) {
		widget.AssignText("Loss:\n" + err.Error())
		<-time.After(2 * time.Second)
		next(nil)
	}

	go func() {
		if err := s.conn.Join(serverName); err != nil {
			widget.AssignText("Failed to join:\n" + err.Error())
			// TODO: Show error label.
			time.AfterFunc(3*time.Second, func() {
				next(nil)
			})
		} else {
			s.conn.SetMessageHandler(nil) // Set to nil to ensure any messages are queued.
			// TODO: Bump to handshaking.
			next(handshake.NewState(s.conn, s.Hostname))
		}
	}()

	return nil
}

func (s *State) Update() error {
	s.layout.Update()
	return nil
}

func (s *State) Draw(screen *ebiten.Image) {
	s.layout.Draw(screen)
}
