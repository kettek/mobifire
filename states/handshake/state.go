package handshake

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/login"
	"github.com/kettek/rebui"
	"github.com/kettek/rebui/widgets"

	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// State provides a handshake step to connecting to a server.
type State struct {
	messages.MessageHandler
	Hostname string
	Port     int
	conn     *net.Connection
	layout   rebui.Layout
}

// NewState creates a new state from a given connection.
func NewState(conn *net.Connection) *State {
	return &State{
		conn: conn,
	}
}

// Enter handles the version & setup messages from/to the server. Failure boots back to metaserver.
func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

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

	// Setup receive just sends to actual login.
	s.Once(&messages.MessageSetup{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		fmt.Println("got setup message!", m.(*messages.MessageSetup), failure)
		// FIXME: Uh... do we have to handle for MessageSetup fail?
		next(login.NewState(s.conn))
	})

	s.Once(&messages.MessageVersion{}, &messages.MessageVersion{}, func(m messages.Message, failure *messages.MessageFailure) {
		msg, ok := m.(*messages.MessageVersion)
		if !ok {
			s.Bail(widget, "expected a version message", next)
			return
		}
		if msg.SVVersion != "1030" {
			s.Bail(widget, "expected a server version 1030", next)
			return
		}
		if err := s.conn.Send(&messages.MessageVersion{CLVersion: "1030", SVName: "mobilefire"}); err != nil {
			s.Bail(widget, "we failed to send version message:\n"+err.Error(), next)
			return
		}
		// FIXME: This isn't optimized, as I'm working relative to termfire.
		if err := s.conn.Send(&messages.MessageSetup{
			FaceCache: struct {
				Use   bool
				Value bool
			}{Use: true, Value: true}, // Changed to false so I can get _all_ the delicious PNGs.
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
			SpellMon: struct {
				Use   bool
				Value uint8
			}{Use: true, Value: 2},
			Tick: struct {
				Use   bool
				Value uint8
			}{Use: true, Value: 1},
		}); err != nil {
			s.Bail(widget, "we failed to send setup message:\n"+err.Error(), next)
			return
		}
	})

	// TODO: Show handshaking label.

	return nil
}

func (s *State) Bail(widget *widgets.Text, msg string, next func(states.State)) {
	widget.AssignText(msg)
	fmt.Println("wut")

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("dang")
		next(nil)
		fmt.Println("ree")
	})
}

func (s *State) Update() error {
	s.layout.Update()
	return nil
}

func (s *State) Draw(screen *ebiten.Image) {
	s.layout.Draw(screen)
	// TODO
}
