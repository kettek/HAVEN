package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	rotation           float64
	x                  float64
	y                  float64
	w                  float64
	h                  float64
	pendingX, pendingY float64
	targetX, targetY   float64
	targetTicker       int
	ticker             int
}

func NewCamera() *Camera {
	w, h := ebiten.WindowSize()
	return &Camera{
		w: float64(w),
		h: float64(h),
	}
}

func (c *Camera) MoveTo(x, y float64) {
	c.pendingX = c.x
	c.pendingY = c.y
	c.targetX = x
	c.targetY = y
	c.targetTicker = 0
}

func (c *Camera) CenterTo(x, y float64) {
	c.pendingX = c.x
	c.pendingY = c.y
	c.targetX = x - c.w/2
	c.targetY = y - c.h/2
	c.targetTicker = 0
}

func (c *Camera) CenterOn(x, y float64) {
	c.x = x - c.w/2
	c.y = y - c.h/2
	c.CenterTo(c.x, c.y)
}

func (c *Camera) Update() error {
	c.ticker++
	w, h := ebiten.WindowSize()
	c.w = float64(w)
	c.h = float64(h)

	if c.targetX != c.pendingX || c.targetY != c.pendingY {
		c.targetTicker++
		c.x = c.pendingX + (c.targetX-c.pendingX)*float64(c.targetTicker)/60
		c.y = c.pendingY + (c.targetY-c.pendingY)*float64(c.targetTicker)/60
		if c.x == c.targetX && c.y == c.targetY {
			c.targetX = c.pendingX
			c.targetY = c.pendingY
		}
	}

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
