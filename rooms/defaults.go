package rooms

import (
	"fmt"

	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
)

func defaultUpdate(r *game.Room) {
	for _, cmd := range r.PendingCommands {
		switch c := cmd.Cmd.(type) {
		case commands.Move:
			ax, ay := cmd.Actor.Position()
			if ax-c.X >= -1 && ax-c.X <= 1 && ay-c.Y >= -1 && ay-c.Y <= 1 {
				cmd.Actor.SetPosition(c.X, c.Y)
			}
		default:
			fmt.Println("handle", cmd.Actor, "wants to", cmd.Cmd)
		}
	}
	r.PendingCommands = nil
}
