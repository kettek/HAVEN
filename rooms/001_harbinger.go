package rooms

import (
	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
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
     E                        
                              
                              
                              
                             e
                              
                              
                              
                              
                              
          V                   
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
                              
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
			"E": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					s.SpriteStack().SetSprite("harbinger-door")
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
				},
			},
		},
		enter: func(w *game.World, r *game.Room) {
		},
		leave: func(w *game.World, r *game.Room) {
		},
		update: func(w *game.World, r *game.Room) {
		},
	}
}
