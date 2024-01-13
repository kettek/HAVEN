package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/inputs"
)

type Actor interface {
	Update(r *Room) (cmd commands.Command)
	Ready() bool
	SetReady(bool)
	TakeTurn() (cmd commands.Command)
	Draw(screen *ebiten.Image, r *Room, geom ebiten.GeoM, drawMode DrawMode)
	DrawPost(screen, post *ebiten.Image, r *Room, geom ebiten.GeoM, drawMode DrawMode)
	SetPosition(int, int, int)
	Position() (int, int, int)
	Command(cmd commands.Command)
	Input(in inputs.Input) bool
	Hover(bool)
	Hovered() bool
	SetTag(string)
	Tag() string
	SetName(string)
	Name() string
	SpriteStack() *SpriteStack
	Interact(w *World, r *Room, other Actor) commands.Command
	Blocks() bool
	SetBlocks(bool)
}

type ActorCommand struct {
	Actor Actor
	Cmd   commands.Command
}

type CombatActor interface {
	Name() string
	CurrentStats() (int, int, int)
	MaxStats() (int, int, int)
	ApplyDamage(int, int, int) (int, int, int)
	ApplyBoost(int, int, int) (int, int, int)
	Level() int
	SetLevel(int)
	Exp() int
	AddExp(int)
	ExpValue() int
	RollBoost() (pen, fire, inte int)
	RollAttack() (pen int)
}
