package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/rooms"
)

type Game struct {
	camera *game.Camera
	room   *game.Room
}

func (g *Game) Update() error {
	g.camera.Update()
	g.room.Update()

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.room.Transition = 60
		x, y := g.room.CenterIso()
		g.camera.MoveTo(x, y)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		g.room.Transition = -60
		x, y := g.room.Center()
		g.camera.MoveTo(x, y)
	}

	return nil
}

func (g *Game) Enter() {
	if g.camera == nil {
		g.camera = game.NewCamera()
	}
	if g.room == nil {
		g.room = rooms.BuildRoom("000_spawn")
	}
}
func (g *Game) Leave() {
}

func (g *Game) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}

	geom.Translate(-g.camera.W/2, -g.camera.H/2)
	geom.Rotate(g.camera.Rotation)
	geom.Translate(g.camera.W/2, g.camera.H/2)
	geom.Translate(-g.camera.X+g.camera.W/2, -g.camera.Y+g.camera.H/2)
	geom.Scale(g.camera.Zoom, g.camera.Zoom)
	g.room.Draw(screen, geom)
}
