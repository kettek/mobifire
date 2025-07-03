package main

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/kettek/mobifire/game"
	"github.com/kettek/mobifire/states/metaserver"
)

func init() {
	// yourgame.Game must implement ebiten.Game interface.
	// For more details, see
	// * https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Game
	game := &game.Game{}
	game.SetNext(&metaserver.State{})
	mobile.SetGame(game)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
