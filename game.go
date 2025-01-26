package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/mobifire/states/metaserver"
)

type Game struct {
	app    fyne.App
	window fyne.Window
	state  states.State
}

func (g *Game) SetNext(state states.State) {
	if g.state != nil {
		g.state.Leave()
	}
	g.state = nil
	if state != nil {
		state.Enter(g.SetNext)
		g.window.SetContent(state.Container())
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
