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
	glitchDead := false
	rooms["003_source"] = Room{
		name:     "source",
		song:     "damaged-haven",
		darkness: 2,
		color:    color.NRGBA{55, 55, 30, 255},
		tiles: `
  #.............#   
  #.............#   
  ##...........##   
  #######_#######   
         ___        
           __       
            _       
           __       
         ___        
         _          
         _          
         __         
          __        
           __       
            _       
           __       
#        ___        
 _________          
#                   
		`,
		tileDefs: TileDefs{
			"#": {
				Name:       "wall of source",
				Sprite:     "brokensight-wall",
				BlocksMove: true,
			},
			".": {
				Name:   "?",
				Sprite: "missing",
			},
			"_": {
				Name:   "path of source",
				Sprite: "harbinger-path",
			},
		},
		entities: `
                    
         1          
                    
                    
                    
                    
                    
                    
                    
                    
                    
                    
                    
                    
                    
                    
                    
E                   
                    
		`,
		entityMap: EntityDefs{
			"E": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to brokensight")
					d.SetTag("brokensight-to-end-door")
					d.SpriteStack().SetSprite("harbinger-door-unlocked")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "002_brokensight",
						Tag:     "brokensight-to-end-door",
						OffsetX: -1,
						Target:  other,
					}
				},
			},
			"1": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("badplayer")
					s.(*actors.Glitch).SetTag("evil")
					s.(*actors.Glitch).SetName("SHOU-09")
					s.(*actors.Glitch).SetLevel(10)
					s.(*actors.Glitch).Skews = true
					s.(*actors.Glitch).SetStats(10, 10, 10)
					s.(*actors.Glitch).SetAbility(&game.Ability{
						Name:     game.AbilityRandomDamage,
						Tier:     10 + rand.Intn(10),
						Turns:    2 + rand.Intn(3),
						Cooldown: 1 + rand.Intn(3),
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
				Duration:   4 * time.Second,
				Color:      color.NRGBA{0, 0, 0, 255},
				Background: color.NRGBA{255, 255, 255, 255},
				Text:       "<SENSE>\ncausation is near",
			})
			<-w.MessageR(game.Message{
				Duration:   3 * time.Second,
				Color:      color.NRGBA{0, 0, 0, 255},
				Background: color.NRGBA{255, 255, 255, 255},
				Text:       "<ACT>\ndestroy",
			})
		},
		turn: func(w *game.World, r *game.Room) {
			if !glitchDead {
				if g := r.GetActorByTag("evil"); g == nil {
					glitchDead = true
					go func() {
						<-w.MessageR(game.Message{
							Duration:   4 * time.Second,
							Color:      color.NRGBA{0, 0, 0, 255},
							Background: color.NRGBA{255, 255, 255, 255},
							Text:       "<KNOW>\nsource gone, haven safe",
						})
						<-time.After(1 * time.Second)
						<-w.MessageR(game.Message{
							Duration:   8 * time.Second,
							Color:      color.NRGBA{205, 205, 180, 255},
							Background: color.NRGBA{0, 0, 0, 200},
							Text:       "THE END\nthanx 4 playin",
						})
					}()
				}
			}
		},
	}
}
