package metaserver

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/mobifire/states/play"
)

type State struct {
	container *fyne.Container
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	button := widget.NewButton("nextie", func() {
		next(&play.State{})
	})

	s.container = container.New(layout.NewCenterLayout(), button)

	return nil
}

func (s *State) Container() *fyne.Container {
	return s.container
}
