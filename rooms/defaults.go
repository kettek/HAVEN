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
			cmd.Actor.SetPosition(c.X, c.Y)
		default:
			fmt.Println("handle", cmd.Actor, "wants to", cmd.Cmd)
		}
	}
	r.PendingCommands = nil
}
