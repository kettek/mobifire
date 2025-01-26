package states

import "fyne.io/fyne/v2"

type State interface {
	Enter(next func(State)) (leave func())
	Container() *fyne.Container // Adopted after Enter()
}
