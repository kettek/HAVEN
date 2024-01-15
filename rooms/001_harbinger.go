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
	rooms["001_harbinger"] = Room{
		name:     "harbinger",
		song:     "infrequent-lament",
		darkness: 2,
		color:    color.NRGBA{15, 7, 26, 255},
		tiles: `
		##### ########################
		#  .....                     #
		#   ...                     .#
		#    _                     ..#
		#    ______________________.. 
		#         _                ..#
		#         _                 .#
		#         _                  #
		#     .....................  #
		#    ... ....... ........... #
		#  __....................... #
		#  _   ....................  #
		#  _   .     _     .      _  #
		# ...  .     _     .      _  #
		#  .   _     _     _      _  #
		#      _     _     _      _  #
		#      _____________      _  #
		#         _        _      _  #
		#     .   _        ________  #
		#    ...___               _  #
		#     .                   _  #
		# .......................... #
		#............................#
		#............................#
		#............................#
		# .......................... #
		#  .   _   .   _    .        #
		#  _   _   _   _    _        #
		#  _   _   _   _    _        #
		###_###_###_###_####_#########
		`,
		tileDefs: TileDefs{
			"#": {
				Name:       "wall of harbinger",
				Sprite:     "harbinger-wall",
				BlocksMove: true,
			},
			".": {
				Name:   "floor of harbinger",
				Sprite: "harbinger-floor",
			},
			"_": {
				Name:   "path of harbinger",
				Sprite: "harbinger-path",
			},
		},
		entities: `
     T                        
                              
                              
                              
                              
                              
                              
                              
                              
             V                
          V   V     w w       
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
      v        v              
           w w    v    v      
                              
                              
                              
   D   B   B   B   B          
		`,
		entityMap: EntityDefs{
			"D": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to triplets")
					d.SetTag("triplets-to-harbinger-door")
					d.SpriteStack().SetSprite("harbinger-door-unlocked")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "001a_triplets",
						Tag:     "triplets-to-harbinger-door",
						OffsetY: 1,
						Target:  other,
					}
				},
			},
			"B": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("harbinger-door")
				},
			},
			"T": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("harbinger-door-unlocked")
					s.SetTag("harbinger-to-brokensight-door")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "002_brokensight",
						Tag:     "harbinger-to-brokensight-door",
						OffsetY: -1,
						Target:  other,
					}
				},
			},
			"e": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("harbinger-door")
				},
			},
			"V": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("glitch-slime")
					s.(*actors.Glitch).SetName("slime")
					s.(*actors.Glitch).SetLevel(rand.Intn(2))
					s.(*actors.Glitch).Skews = true
					s.(*actors.Glitch).SetStats(4, 4, 8)
				},
			},
			"v": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("glitch-eye")
					s.(*actors.Glitch).SetName("eye")
					s.(*actors.Glitch).SetLevel(rand.Intn(2))
					s.(*actors.Glitch).Z = 1
					s.(*actors.Glitch).Floats = true
					s.(*actors.Glitch).SetStats(8, 4, 6)
				},
			},
			"w": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					g := s.(*actors.Glitch)
					g.SpriteStack().SetSprite("glitch-wanderer")
					g.SetName("wanderer")
					g.Wanders = true
					g.SetStats(2, 8, 4)
				},
			},
		},
		enter: func(w *game.World, r *game.Room) {
			if !first {
				return
			}
			first = true
			<-w.MessageR(game.Message{
				Duration:   3 * time.Second,
				Color:      color.NRGBA{0, 0, 0, 255},
				Background: color.NRGBA{255, 255, 255, 255},
				Text:       "<SENSE>\nhigh corruption variety",
			})
		},
		leave: func(w *game.World, r *game.Room) {
		},
		update: func(w *game.World, r *game.Room) {
		},
	}
}
