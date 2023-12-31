package rooms

import (
	"fmt"
	"image/color"

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
			"D": "haven-door",
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
			r.Color = color.NRGBA{205, 205, 180, 255}
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
