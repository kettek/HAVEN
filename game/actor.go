package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
)

type Actor interface {
	Update(r *Room) (cmd commands.Command)
	Draw(screen *ebiten.Image, r *Room, geom ebiten.GeoM, drawMode DrawMode)
	Position() (int, int)
	Command(cmd commands.Command)
	Hover(bool)
	Hovered() bool
	Name() string
	SpriteStack() *SpriteStack
	Interact(w *World, r *Room, other Actor) commands.Command
}

type ActorCommand struct {
	Actor Actor
	Cmd   commands.Command
}
