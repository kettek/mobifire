package states

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// State is our interface for any distinct part of the program that should be considered the main process.
type State interface {
	Enter(next func(State)) (leave func())
	Draw(*ebiten.Image)
	Update() error
}

// See Prior
type statePrior struct {
}

// Enter means nothing
func (s *statePrior) Enter(next func(State)) (leave func()) {
	return nil
}

// Prior is used to return back to the previous state.
var Prior = &statePrior{}
