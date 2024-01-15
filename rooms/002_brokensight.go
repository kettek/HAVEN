package rooms

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
)

func init() {
	first := true
	rooms["002_brokensight"] = Room{
		name:     "brokensight",
		song:     "infrequent-lament",
		darkness: 2,
		color:    color.NRGBA{26, 7, 7, 255},
		tiles: `
		            ###########         ########
		     _______..1........_________.......E
		     _      #....1....#         #......#
		     _      #####.#####         #......#
		#####.####       _              ###.####
		#..2....2#       _                 _    
		#........#       _                 _    
		#####.####  #####.#####          ##.####
		     _      #.........#          #.....#
		     _______....2...2..__________......#
		      _     #..........__________...1..#
		      _     ###########    __    #.....#
		      _                    __    #.....#
		  ####.###                 __    #######
		  #......#            #####..#          
		  #......#            #......#          
		  #.......____________....1..#          
		  #......#            #......#          
		  ###.####            ########          
		    #D#                                 
		`,
		tileDefs: TileDefs{
			"#": {
				Name:       "wall of brokensight",
				Sprite:     "brokensight-wall",
				BlocksMove: true,
			},
			".": {
				Name:   "floor of brokensight",
				Sprite: "brokensight-floor",
			},
			"_": {
				Name:   "path of brokensight",
				Sprite: "brokensight-path",
			},
		},
		entities: `
                                        
              1                        E
                 1                      
                                        
                                        
   2    2                               
                                        
                                        
                                        
                2   2                   
                                    1   
                                        
                                        
                                        
                                        
                                        
                          1             
                                        
                                        
     D                                  
		`,
		entityMap: EntityDefs{
			"D": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to harbinger")
					d.SetTag("harbinger-to-brokensight-door")
					d.SpriteStack().SetSprite("harbinger-door-unlocked")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "001_harbinger",
						Tag:     "harbinger-to-brokensight-door",
						OffsetY: 1,
						Target:  other,
					}
				},
			},
			"E": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to the end")
					d.SetTag("brokensight-to-end-door")
					d.SpriteStack().SetSprite("harbinger-door-unlocked")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "003_source",
						Tag:     "brokensight-to-end-door",
						OffsetX: 1,
						Target:  other,
					}
				},
			},
			"1": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("glitch-warp")
					s.(*actors.Glitch).SetName("warp")
					s.(*actors.Glitch).SetLevel(2 + rand.Intn(4))
					s.(*actors.Glitch).Skews = true
					s.(*actors.Glitch).SetStats(10, 12, 10)
					s.(*actors.Glitch).SetAbility(&game.Ability{
						Name:     game.AbilityRandomDamage,
						Tier:     2 + rand.Intn(2),
						Turns:    2 + rand.Intn(3),
						Cooldown: 1 + rand.Intn(3),
					})
				},
			},
			"2": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("glitch-tripo")
					s.(*actors.Glitch).SetName("tripo")
					s.(*actors.Glitch).SetLevel(2 + rand.Intn(4))
					s.(*actors.Glitch).Z = 1
					s.(*actors.Glitch).Floats = true
					s.(*actors.Glitch).SetStats(16, 8, 10)
					s.(*actors.Glitch).SetAbility(&game.Ability{
						Name:     game.AbilityCleave,
						Tier:     1,
						Turns:    1,
						Cooldown: 2 + rand.Intn(3),
					})
				},
			},
		},
		enter: func(w *game.World, r *game.Room) {
			if !first {
				return
			}
			first = false
			<-w.MessageR(game.Message{
				Duration:   5 * time.Second,
				Color:      color.NRGBA{0, 0, 0, 255},
				Background: color.NRGBA{255, 255, 255, 255},
				Text:       "<THINK>\never closer to the SOURCE",
			})
		},
	}
}
