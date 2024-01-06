package rooms

import (
	"fmt"
	"image/color"
	"time"

	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/res"
)

func init() {
	doorLocked := true
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
		},
		entities: `
		  DT 
		     
		  @  
		     
		     
		`,
		entityMap: EntityDefs{
			"@": {
				Actor: "player",
			},
			"D": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to ![unknown]")
					d.SpriteStack().SetSprite("haven-door")
					d.SetTag("door")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					fmt.Println("it be a door interacted with")
					if doorLocked {
						fmt.Println("it is locked")
					}
					return nil
				},
			},
			"T": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					s.(*actors.Interactable).SetName("terminal")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					prompts := []string{"Query Mainframe", "Manage Safeguard", "Leave"}
					//res.PlaySound("button")
					poweron := res.PlaySound("poweron")
					powered := res.GetSound("powered")
					poweroff := res.GetSound("poweroff")
					poweron.Next = powered

					powered.Looping = true
					powered.Next = poweroff
					return commands.Prompt{
						Items: prompts,
						Handler: func(index int, result string) bool {
							if index == 0 {
								w.AddPrompt([]string{"Return"}, "Mainframe status... corrupted.\nSolution: purge system", func(index int, result string) bool {
									return true
								})
								return false
							} else if index == 1 {
								status := "Safeguard: "
								if doorLocked {
									status += "locked"
								} else {
									status += "unlocked"
								}
								w.AddPrompt([]string{"Lock", "Unlock", "Return"}, status, func(index int, result string) bool {
									if index == 0 {
										if !doorLocked {
											res.PlaySound("lock")
										}
										doorLocked = true
										r.GetActorByTag("door").SpriteStack().SetSprite("haven-door")
										w.Prompts[len(w.Prompts)-1].Message = "Safeguard: locked"
										return false
									} else if index == 1 {
										if doorLocked {
											res.PlaySound("unlock")
										}
										doorLocked = false
										r.GetActorByTag("door").SpriteStack().SetSprite("haven-door-unlocked")
										w.Prompts[len(w.Prompts)-1].Message = "Safeguard: unlocked"
										return false
									}
									return true
								})
								return false
							}
							poweron.Next = poweroff // Set poweron's next to poweroff just in case the player exits the menu quickly.
							powered.Looping = false
							powered.Pause()
							return true
						},
					}
				},
			},
		},
		metadata: make(map[string]interface{}),
		enter: func(w *game.World, r *game.Room) {
			// Get our player.
			for _, a := range r.Actors {
				if a.Name() == "player" {
					w.PlayerActor = a
					break
				}
			}
			fmt.Println("enter called")
			makeBigMsg := func(s string, d time.Duration, c color.NRGBA) game.Message {
				return game.Message{Text: s, Duration: d, Color: c, Font: &res.BigFont}
			}
			<-w.FuncR(func() {
				r.Color = color.NRGBA{0, 0, 0, 255}
			})
			delayTimeR(2 * time.Second)
			clr := color.NRGBA{0, 255, 0, 255}
			/*s := "/activate SHOU"
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
			<-w.MessageR(makeBigMsg("defense system <SHOU> online", 500*time.Millisecond, color.NRGBA{205, 205, 180, 255}))*/
			<-w.FuncR(func() {
				r.SetColor(color.NRGBA{205, 205, 180, 255})
			})
			go func() {
				<-w.MessageR(makeBigMsg("<LMB> to investigate, <RMB> to act", 8000*time.Millisecond, clr))
			}()
		},
		leave: func(w *game.World, r *game.Room) {
			fmt.Println("left spawn")
		},
		update: func(w *game.World, r *game.Room) {
		},
	}
}
