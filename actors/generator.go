package actors

import "github.com/kettek/ebihack23/game"

func New(actor string, y, x int) game.Actor {
	if f, ok := actors[actor]; ok {
		return f(y, x)
	}
	return nil
}

var actors = make(map[string]func(y, x int) game.Actor)
