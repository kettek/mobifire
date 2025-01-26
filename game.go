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
	g.state = nil
	if state != nil {
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
	g.window.Resize(fyne.NewSize(360, 800))

	// Set our initial state...
	g.SetNext(&metaserver.State{})

	g.window.ShowAndRun()

	return g
}
