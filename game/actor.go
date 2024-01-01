package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
)

type Actor interface {
	Update(r *Room) (cmd commands.Command)
	Draw(screen *ebiten.Image, r *Room, geom ebiten.GeoM, drawMode DrawMode, ratio float64)
	Position() (int, int)
	SetPosition(x, y int)
	Hover(bool)
	Hovered() bool
	Name() string
}

type ActorCommand struct {
	Actor Actor
	Cmd   commands.Command
}
