package inputs

import "github.com/hajimehoshi/ebiten/v2"

type Input interface{}

type Direction struct {
	X, Y int
	Mod  bool
}

type Confirm struct{}

type Cancel struct{}

type Click struct {
	X, Y  float64
	Which ebiten.MouseButton
	Mod   bool
}
