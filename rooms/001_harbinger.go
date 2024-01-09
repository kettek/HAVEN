package rooms

import (
	"math/rand"

	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/res"
)

func init() {
	rooms["001_harbinger"] = Room{
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
                              
                              
                              
                             e
                              
                              
                              
                              
                              
          V                   
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
      v                       
                              
                              
                              
                              
   D   B   B   B   B          
		`,
		entityMap: EntityDefs{
			"D": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to haven")
					d.SetTag("haven-door")
					d.SpriteStack().SetSprite("haven-door-unlocked")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "000_spawn",
						Tag:     "haven-door",
						OffsetY: 1,
						Target:  other,
					}
				},
			},
			"B": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("harbinger-door-sealed")
				},
			},
			"T": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("harbinger-door")
					s.SetTag("triplets-door")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "001a_triplets",
						Tag:     "harbinger-door",
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
					s.(*actors.Glitch).SetStats(5, 5, 15)
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
					s.(*actors.Glitch).SetStats(15, 5, 5)
				},
			},
		},
		enter: func(w *game.World, r *game.Room) {
			res.Jukebox.Play("infrequent-lament")
		},
		leave: func(w *game.World, r *game.Room) {
		},
		update: func(w *game.World, r *game.Room) {
		},
	}
}
