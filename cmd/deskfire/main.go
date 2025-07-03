package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/game"
	"github.com/kettek/mobifire/states/metaserver"
)

func main() {
	game := &game.Game{}
	game.SetNext(&metaserver.State{})

	ebiten.SetWindowSize(800, 360)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
