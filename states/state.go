package states

import "github.com/hajimehoshi/ebiten/v2"

type State interface {
	Draw(screen *ebiten.Image)
	Update() error
	Enter()
	Leave()
}
