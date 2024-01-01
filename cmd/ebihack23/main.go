package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/states"
)

type ebihack struct {
}

func (e *ebihack) Update() error {
	return states.CurrentState.Update()
}
func (e *ebihack) Draw(screen *ebiten.Image) {
	states.CurrentState.Draw(screen)
}
func (e *ebihack) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	// Start dat gizame.
	g := &ebihack{}

	states.NextState(states.NewGame())

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
