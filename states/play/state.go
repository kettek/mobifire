package play

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/kettek/mobifire/states"
)

type State struct {
	container *fyne.Container
}

func (s *State) Enter(next func(states.State)) {

	mb := NewMultiBoard(11, 11, 8)

	s.container = container.New(layout.NewCenterLayout(), mb.container)
}

func (s *State) Leave() {
}

func (s *State) Container() *fyne.Container {
	return s.container
}
