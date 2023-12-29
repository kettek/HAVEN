package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	camera *Camera
	// ...
	currentRoom *Room
}

func NewGame() *Game {
	g := &Game{}

	g.camera = NewCamera()
	g.currentRoom = NewRoom(10, 10)

	return g
}

func (g *Game) Update() error {
	g.camera.Update()
	g.currentRoom.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{128, 128, 128, 255})

	geom := ebiten.GeoM{}
	geom.Translate(50, 50)
	geom.Scale(3, 3)

	geom.Translate(g.camera.x, g.camera.y)
	geom.Translate(-g.camera.w/2, -g.camera.h/2)
	geom.Rotate(g.camera.rotation)
	geom.Translate(g.camera.w/2, g.camera.h/2)

	g.currentRoom.Draw(screen, geom)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
