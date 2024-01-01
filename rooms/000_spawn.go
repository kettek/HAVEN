package rooms

import (
	"fmt"
	"image/color"
	"time"

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
			r.Color = color.NRGBA{0, 0, 0, 255}
			delayTimeR(2 * time.Second)
			clr := color.NRGBA{0, 255, 0, 255}
			s := "/activate SHOU"
			for i := range s {
				u := ""
				if i%2 == 0 {
					u = "_"
				}
				r.MessageR(game.Message{Text: string(s[:i]) + u, Duration: 200 * time.Millisecond, Color: clr})
			}
			r.MessageR(game.Message{Text: string(s), Duration: 1000 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: ".", Duration: 500 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: "..", Duration: 500 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: "...", Duration: 500 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: ".", Duration: 500 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: "..", Duration: 500 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: "...", Duration: 500 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: "defense system <SHOU> online", Duration: 500 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: "defense system <SHOU> online", Duration: 500 * time.Millisecond, Color: color.NRGBA{205, 205, 180, 255}})
			r.MessageR(game.Message{Text: "defense system <SHOU> online", Duration: 500 * time.Millisecond, Color: clr})
			r.MessageR(game.Message{Text: "defense system <SHOU> online", Duration: 500 * time.Millisecond, Color: color.NRGBA{205, 205, 180, 255}})
			r.SetColor(color.NRGBA{205, 205, 180, 255})
		},
		leave: func(r *game.Room) {
			fmt.Println("left spawn")
		},
		update: func(r *game.Room) {
			defaultUpdate(r)
		},
	}
}
