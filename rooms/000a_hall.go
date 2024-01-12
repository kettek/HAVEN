package rooms

import (
	"fmt"
	"math"

	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/res"
)

func init() {
	doorLocked := true
	rooms["000a_hall"] = Room{
		tiles: `// First line is ignored because lazy.
		#   ##   ##   ##   ##   ##    ###v###             ### ###       
		##tv###tv###tv###tv###tv##      #_#                 #d#         
		# __   __   __   __   __         _                   _          
		#  _    _    _    _    _       .......               _        # 
		##........................   ....   ....  ......     ...      ##
		#.............................        ..................   ,,,,,
		##........................   ....   ....  ......     ...      ##
		#  _    _    _    _    _       .......                        # 
		#  __   __   __   __   __        _                   ,          
		###^T###^T###^D###^T###^T#      #_#                 #,#         
		 #   ##   ##   ##   ##   #    ###^###             ###,###       
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
			"_": {
				Name:   "path of haven",
				Sprite: "haven-path",
			},
		},
		entities: `
		     
		                                                     E
		     
		     
		     
		     
		     
		     
		     
		             DT
		`,
		entityMap: EntityDefs{
			"E": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to ![harbinger]")
					d.SpriteStack().SetSprite("harbinger-door")
					d.SetTag("hall-door")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "001_harbinger",
						Tag:     "hall-door",
						OffsetY: -1,
						Target:  other,
					}
				},
			},
			"D": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to ![haven]")
					d.SpriteStack().SetSprite("haven-door")
					d.SetTag("haven-door")
					d.SpriteStack().Rotation = math.Pi
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					if doorLocked {
						return nil
					}

					return commands.Travel{
						Room:    "000_spawn",
						Tag:     "haven-door",
						OffsetY: 1,
						Target:  other,
					}
				},
			},
			"T": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					s.(*actors.Interactable).SetName("terminal")
					s.SpriteStack().SetSprite("terminal-off")
					s.SetTag("terminal")
					s.SpriteStack().Rotation = math.Pi
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					r.GetActorByTag("terminal").SpriteStack().SetSprite("terminal")
					prompts := []string{"Query z-level SHOU", "Manage Safeguard", "Leave"}
					//res.PlaySound("button")
					poweron := res.PlaySound("poweron")
					powered := res.GetSound("powered")
					poweroff := res.GetSound("poweroff")
					poweron.Next = powered

					powered.Looping = true
					powered.Next = poweroff
					return commands.Prompt{
						Items:        prompts,
						ShowVersions: true,
						Handler: func(index int, result string) bool {
							if index == 0 {
								w.AddPrompt([]string{"Return"}, "01-05: lost\n06   : released\n07-08: missing\n09   : ???", func(index int, result string) bool {
									return true
								}, true)
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
										r.GetActorByTag("haven-door").SpriteStack().SetSprite("haven-door")
										w.Prompts[len(w.Prompts)-1].Message = "Safeguard: locked"
										return false
									} else if index == 1 {
										if doorLocked {
											res.PlaySound("unlock")
										}
										doorLocked = false
										r.GetActorByTag("haven-door").SpriteStack().SetSprite("haven-door-unlocked")
										w.Prompts[len(w.Prompts)-1].Message = "Safeguard: unlocked"
										return false
									}
									return true
								}, true)
								return false
							}
							poweron.Next = poweroff // Set poweron's next to poweroff just in case the player exits the menu quickly.
							powered.Looping = false
							powered.Pause()
							r.GetActorByTag("terminal").SpriteStack().SetSprite("terminal-off")
							return true
						},
					}
				},
			},
		},
		metadata: make(map[string]interface{}),
		enter: func(w *game.World, r *game.Room) {
		},
		leave: func(w *game.World, r *game.Room) {
			fmt.Println("left spawn")
		},
		update: func(w *game.World, r *game.Room) {
		},
	}
}
