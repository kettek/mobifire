package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/states"
)

type Game struct {
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

		g.priorState = priorState
		g.leaveCb = state.Enter(g.SetNext)
	} else if g.firstState != nil { // Bump back to first state if we can! This should be guaranteed to be the metaserver.
		g.leaveCb = g.firstState.Enter(g.SetNext)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.state.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	return g.state.Update()
}
