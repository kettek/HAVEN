package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	// FIXME: Just a test.
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.currentRoom.transition = 60
		x, y := g.currentRoom.CenterIso()
		g.camera.CenterTo(x*3.0, y*3.0)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		g.currentRoom.transition = -60
		x, y := g.currentRoom.Center()
		g.camera.CenterTo(x*3.0, y*3.0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(3, 3)

	geom.Translate(-g.camera.x, -g.camera.y)
	geom.Translate(-g.camera.w/2, -g.camera.h/2)
	geom.Rotate(g.camera.rotation)
	geom.Translate(g.camera.w/2, g.camera.h/2)

	g.currentRoom.Draw(screen, geom)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
