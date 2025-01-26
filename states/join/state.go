package join

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
)

type State struct {
	container *fyne.Container
	Hostname  string
	Port      int
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	fmt.Println("joining", s.Hostname, s.Port)

	label := widget.NewLabel("Joining " + s.Hostname + ":" + fmt.Sprint(s.Port))

	s.container = container.New(layout.NewCenterLayout(), label)

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}
