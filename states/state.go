package states

import "fyne.io/fyne/v2"

// State is our interface for any distinct part of the program that should be considered the main process.
type State interface {
	Enter(next func(State)) (leave func())
	Container() *fyne.Container // Adopted after Enter()
}

// See Prior
type statePrior struct {
}

// Enter means nothing
func (s *statePrior) Enter(next func(State)) (leave func()) {
	return nil
}

// Container means nothing
func (s *statePrior) Container() *fyne.Container {
	return nil
}

// Prior is used to return back to the previous state.
var Prior = &statePrior{}

// StateWithWindow is an extension of State that allows setting the fyne window (needed for dialog.Show* funcs)
type StateWithWindow interface {
	State
	SetWindow(window fyne.Window)
}

// StateWithApp is an extension of State that allows setting the fyne app (needed for Preferences)
type StateWithApp interface {
	State
	SetApp(app fyne.App)
}
