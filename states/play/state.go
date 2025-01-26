package play

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/kettek/mobifire/states"
)

type State struct {
	container *fyne.Container
	mb        *multiBoard
}

func (s *State) Enter(next func(states.State)) (leave func()) {

	s.mb = newMultiBoard(11, 11, 8)

	s.container = container.New(layout.NewCenterLayout(), s.mb.container)

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}
