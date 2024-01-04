package actors

import (
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
)

func New(actor string, y, x int, ctor CreateFunc, interact InteractFunc) game.Actor {
	if f, ok := actors[actor]; ok {
		return f(y, x, ctor, interact)
	}
	return nil
}

type CreateFunc func(s game.Actor)
type InteractFunc func(w *game.World, r *game.Room, s, o game.Actor) commands.Command

var actors = make(map[string]func(y, x int, ctor CreateFunc, interact InteractFunc) game.Actor)
