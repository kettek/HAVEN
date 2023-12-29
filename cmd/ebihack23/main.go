package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/game"
)

func main() {
	// Start dat gizame.
	g := game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
