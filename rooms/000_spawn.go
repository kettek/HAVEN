package rooms

import (
	"fmt"

	"github.com/kettek/ebihack23/game"
)

func init() {
	rooms["000_spawn"] = Room{
		tiles: `// First line is ignored because lazy.
		##DT#
		# ..#
		#...#
		# . #
		#####
		`,
		tileMap: map[string]string{
			"#": "haven-wall",
			".": "haven-floor",
			"D": "door",
			"T": "terminal",
		},
		entities: `
		     
		     
		  @  
		     
		     
		`,
		entityMap: map[string]string{
			"@": "player",
		},
		metadata: make(map[string]interface{}),
		enter: func(r *game.Room) {
			fmt.Println("entered spawn")
		},
		leave: func(r *game.Room) {
			fmt.Println("left spawn")
		},
		update: func(r *game.Room) {
			defaultUpdate(r)
		},
	}
}
