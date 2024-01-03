package rooms

import (
	"fmt"
	"image/color"
	"time"

	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/res"
)

func init() {
	rooms["000_spawn"] = Room{
		tiles: `// First line is ignored because lazy.
		##D #
		# ..#
		#...#
		# . #
		#####
		`,
		tileDefs: TileDefs{
			"#": {
				Name:       "wall of haven",
				Sprite:     "haven-wall",
				BlocksMove: true,
			},
			".": {
				Name:   "floor of haven",
				Sprite: "haven-floor",
			},
			"D": {
				Name:       "door to <unknown>",
				Sprite:     "haven-door",
				BlocksMove: true,
			},
		},
		entities: `
		   T 
		     
		  @  
		     
		     
		`,
		entityMap: map[string]string{
			"@": "player",
			"T": "terminal",
		},
		metadata: make(map[string]interface{}),
		enter: func(w *game.World, r *game.Room) {
			fmt.Println("enter called")
			makeBigMsg := func(s string, d time.Duration, c color.NRGBA) game.Message {
				return game.Message{Text: s, Duration: d, Color: c, Font: res.BigFont}
			}
			<-w.FuncR(func() {
				r.Color = color.NRGBA{0, 0, 0, 255}
			})
			delayTimeR(2 * time.Second)
			clr := color.NRGBA{0, 255, 0, 255}
			s := "/activate SHOU"
			for i := range s {
				u := ""
				if i%2 == 0 {
					u = "_"
				}
				<-w.MessageR(makeBigMsg(string(s[:i])+u, 200*time.Millisecond, clr))
			}
			<-w.MessageR(makeBigMsg(s, 1000*time.Millisecond, clr))
			<-r.DropInR()
			<-w.MessageR(makeBigMsg(".", 500*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("..", 500*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("..", 500*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg(".", 500*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("..", 500*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("..", 500*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("defense system <SHOU> online", 500*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("defense system <SHOU> online", 500*time.Millisecond, color.NRGBA{205, 205, 180, 255}))
			<-w.MessageR(makeBigMsg("defense system <SHOU> online", 500*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("defense system <SHOU> online", 500*time.Millisecond, color.NRGBA{205, 205, 180, 255}))
			<-w.FuncR(func() {
				r.SetColor(color.NRGBA{205, 205, 180, 255})
			})
		},
		leave: func(w *game.World, r *game.Room) {
			fmt.Println("left spawn")
		},
		update: func(w *game.World, r *game.Room) {
		},
	}
}
