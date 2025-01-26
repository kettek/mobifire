package join

import (
	"fmt"
	"time"

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
	conn      net.Connection
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	label := widget.NewLabel("Joining " + s.Hostname + ":" + fmt.Sprint(s.Port) + "...")

	serverName := s.Hostname
	if s.Port != 0 {
		serverName += ":" + fmt.Sprint(s.Port)
	}
	s.container = container.New(layout.NewCenterLayout(), label)

	go func() {
		if err := s.conn.Join(serverName); err != nil {
			label.SetText("Failed to join " + serverName + ": " + err.Error())
			time.AfterFunc(3*time.Second, func() {
				next(nil)
			})
		} else {
			s.conn.SetMessageHandler(nil)
			next(login.NewState(&s.conn))
		}
	}()

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}
