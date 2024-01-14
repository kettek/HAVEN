package rooms

import (
	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
)

func init() {
	rooms["001a_triplets"] = Room{
		name: "the triplets",
		song: "infrequent-lament",
		tiles: `// First line is ignored because lazy.
    #######   
   ##     ##  
  ##       ## 
 ##         ##
 #  .  .  .  #
 #  _  _  _  #
 #  _  _  _  #
 #  _  _  _  #
 ##.........##
  ##.......## 
   ##.....##  
    ###.###   
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
           
             
             
             
    1  2  3  
             
             
             
             
             
             
       D     
		`,
		entityMap: EntityDefs{
			"D": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to harbinger")
					d.SpriteStack().SetSprite("harbinger-door")
					d.SetTag("harbinger-door")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "001_harbinger",
						Tag:     "triplets-door",
						OffsetY: 1,
						Target:  other,
					}
				},
			},
			"1": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					g := s.(*actors.Glitch)
					g.SpriteStack().SetSprite("minion-pen")
					g.SetName("minion-pen")
					g.Z = 1
					g.Floats = true
					g.Wanders = false
					g.SpriteStack().YScale = 1.0
					g.SpriteStack().Shaded = true
					g.SetLevel(2)
					g.SetStats(10, 5, 5)
				},
			},
			"2": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					g := s.(*actors.Glitch)
					g.SpriteStack().SetSprite("minion-wall")
					g.SetName("minion-wall")
					g.Z = 1
					g.Floats = true
					g.Wanders = false
					g.SpriteStack().YScale = 1.0
					g.SpriteStack().Shaded = true
					g.SetLevel(2)
					g.SetStats(5, 5, 10)
				},
			},
			"3": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					g := s.(*actors.Glitch)
					g.SpriteStack().SetSprite("minion-shel")
					g.SetName("minion-shel")
					g.Z = 1
					g.Floats = true
					g.Wanders = false
					g.SpriteStack().YScale = 1.0
					g.SpriteStack().Shaded = true
					g.SetLevel(2)
					g.SetStats(5, 10, 5)
				},
			},
		},
		metadata: make(map[string]interface{}),
		enter: func(w *game.World, r *game.Room) {
			/*makeBigMsg := func(s string, d time.Duration, c color.NRGBA) game.Message {
				return game.Message{Text: s, Duration: d, Color: c, Font: &res.BigFont}
			}
			delayTimeR(1 * time.Second)
			clr := color.NRGBA{0, 255, 0, 255}
			<-w.MessageR(makeBigMsg("We are the triplets, three.", 3000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("It hurts to be...", 3000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("but not due to thee", 3000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("but rather, GRE!", 3000*time.Millisecond, clr))*/
		},
		leave: func(w *game.World, r *game.Room) {
		},
		update: func(w *game.World, r *game.Room) {
		},
	}
}
