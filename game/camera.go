package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	Rotation           float64
	X                  float64
	Y                  float64
	W                  float64
	H                  float64
	Zoom               float64
	pendingX, pendingY float64
	targetX, targetY   float64
	targetTicker       int
	ticker             int
}

func NewCamera() *Camera {
	//w, h := ebiten.WindowSize()
	w, h := 640, 480
	return &Camera{
		W:    float64(w) / 3,
		H:    float64(h) / 3,
		Zoom: 3.0,
	}
}

func (c *Camera) MoveTo(x, y float64) {
	c.pendingX = c.X
	c.pendingY = c.Y
	c.targetX = x
	c.targetY = y
	c.targetTicker = 0
}

func (c *Camera) CenterTo(x, y float64) {
	c.pendingX = c.X
	c.pendingY = c.Y
	c.targetX = x - c.W
	c.targetY = y - c.H
	c.targetTicker = 0
}

func (c *Camera) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
	c.pendingX = x
	c.pendingY = y
	c.targetX = x
	c.targetY = y
	c.targetTicker = 0
}

func (c *Camera) CenterOn(x, y float64) {
	c.X = x - c.W
	c.Y = y - c.H
	c.CenterTo(c.X, c.Y)
}

func (c *Camera) Update() error {
	c.ticker++
	//w, h := ebiten.WindowSize()
	w, h := 640, 480
	c.W = float64(w) / c.Zoom
	c.H = float64(h) / c.Zoom

	if c.targetX != c.pendingX || c.targetY != c.pendingY {
		c.targetTicker++
		c.X = c.pendingX + (c.targetX-c.pendingX)*float64(c.targetTicker)/60
		c.Y = c.pendingY + (c.targetY-c.pendingY)*float64(c.targetTicker)/60
		if c.X == c.targetX && c.Y == c.targetY {
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
		c.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.X += 1
	}

	// Wiggle rotation.
	c.Rotation = math.Sin(float64(c.ticker)/100) / 100

	return nil
}
