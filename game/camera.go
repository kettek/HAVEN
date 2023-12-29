package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	rotation float64
	x        float64
	y        float64
	w        float64
	h        float64
	ticker   int
}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) Update() error {
	c.ticker++
	w, h := ebiten.WindowSize()
	/*if c.w != float64(w) || c.h != float64(h) {
		c.x = float64(w) / 2
		c.y = float64(h) / 2
	}*/
	c.w = float64(w)
	c.h = float64(h)

	/*if ebiten.IsKeyPressed(ebiten.KeyQ) {
		c.rotation -= 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		c.rotation += 0.05
	}*/
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.x -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.x += 1
	}

	// Wiggle rotation.
	c.rotation = math.Sin(float64(c.ticker)/100) / 100

	return nil
}
