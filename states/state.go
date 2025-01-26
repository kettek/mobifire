package states

import "fyne.io/fyne/v2"

type State interface {
	Enter(next func(State)) (leave func())
	Container() *fyne.Container // Adopted after Enter()
}

// StatePrior is some terrible B.S....
type StatePrior struct {
}

// Enter means nothing
func (s *StatePrior) Enter(next func(State)) (leave func()) {
	return nil
}

// Container means nothing
func (s *StatePrior) Container() *fyne.Container {
	return nil
}

// StateWithWindow is an extension of State that allows setting the fyne window (needed for dialog.Show* funcs)
type StateWithWindow interface {
	State
	SetWindow(window fyne.Window)
}
