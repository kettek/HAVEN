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
		g.camera.CenterTo(x*3.0, y*3.0)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		g.room.Transition = -60
		x, y := g.room.Center()
		g.camera.CenterTo(x*3.0, y*3.0)
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
	geom.Scale(3, 3)

	geom.Translate(-g.camera.X, -g.camera.Y)
	geom.Translate(-g.camera.W/2, -g.camera.H/2)
	geom.Rotate(g.camera.Rotation)
	geom.Translate(g.camera.W/2, g.camera.H/2)
	g.room.Draw(screen, geom)
}
