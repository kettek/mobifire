package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/mobifire/states/metaserver"
)

type Game struct {
	app        fyne.App
	window     fyne.Window
	firstState states.State // Used to ensure Server state is returned to.
	priorState states.State // Absolute bogus handle to just bounce back to last state.
	state      states.State
	leaveCb    func()
}

func (g *Game) SetNext(state states.State) {
	if g.firstState == nil {
		g.firstState = state
	}
	if g.leaveCb != nil {
		g.leaveCb()
	}
	var priorState states.State
	priorState = g.state
	g.state = state
	if state != nil {
		// Prior state a lil hacky, but oh well~~~
		if state == states.Prior {
			if g.priorState != nil {
				g.state = g.priorState
				state = g.priorState
				priorState = g.priorState
			} else {
				g.state = g.firstState
				state = g.firstState
				priorState = g.firstState
			}
		}

		// Set window if interface conforms.
		if s, ok := state.(states.StateWithWindow); ok {
			s.SetWindow(g.window)
		}

		g.priorState = priorState
		g.leaveCb = state.Enter(g.SetNext)
		g.window.SetContent(state.Container())
	} else if g.firstState != nil { // Bump back to first state if we can! This should be guaranteed to be the metaserver.
		g.leaveCb = g.firstState.Enter(g.SetNext)
		g.window.SetContent(g.firstState.Container())
	}
}

func NewGame() *Game {
	g := &Game{
		app: app.New(),
	}
	g.window = g.app.NewWindow("Crossfire Mobile")
	g.window.Resize(fyne.NewSize(800, 360))
	g.window.SetFixedSize(true)

	// Set our initial state...
	g.SetNext(&metaserver.State{})

	g.window.ShowAndRun()

	return g
}
