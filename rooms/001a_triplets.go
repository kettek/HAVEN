package rooms

import (
	"image/color"
	"time"

	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/res"
)

func init() {
	first := true
	rooms["001a_triplets"] = Room{
		name:     "the triplets",
		song:     "infrequent-lament",
		darkness: 2,
		color:    color.NRGBA{15, 7, 26, 255},
		tiles: `// First line is ignored because lazy.
    ###.###   
   ##  _  ##  
  ##_______## 
 ## _  _  _ ##
 #  .  .  .  #
 #  _  _  _  #
 #  _  _  _  #
 ##.........##
  ###.....### 
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
       E   
             
             
             
    1  2  3  
             
             
             
             
       D     
		`,
		entityMap: EntityDefs{
			"D": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to hall")
					d.SpriteStack().SetSprite("haven-door-unlocked")
					d.SetTag("hall-to-triplets-door")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "000a_hall",
						Tag:     "hall-to-triplets-door",
						OffsetY: 1,
						Target:  other,
					}
				},
			},
			"E": {
				Actor: "interactable",
				OnCreate: func(s game.Actor) {
					d := s.(*actors.Interactable)
					d.SetName("door to harbinger")
					d.SpriteStack().SetSprite("harbinger-door")
					d.SetTag("triplets-to-harbinger-door")
				},
				OnInteract: func(w *game.World, r *game.Room, s game.Actor, other game.Actor) commands.Command {
					return commands.Travel{
						Room:    "001_harbinger",
						Tag:     "triplets-to-harbinger-door",
						OffsetY: -1,
						Target:  other,
					}
				},
			},
			"1": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					g := s.(*actors.Glitch)
					g.SpriteStack().SetSprite("minion-pen")
					g.SetName("minpen")
					g.Z = 1
					g.Floats = true
					g.Wanders = false
					g.SpriteStack().YScale = 1.0
					g.SpriteStack().Shaded = true
					g.SetLevel(2)
					g.SetStats(10, 5, 5)
					g.SetAbility(&game.Ability{
						Name:     game.AbilityPerfectHit,
						Tier:     2,
						Turns:    2,
						Cooldown: 2,
					})
				},
			},
			"2": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					g := s.(*actors.Glitch)
					g.SpriteStack().SetSprite("minion-wall")
					g.SetName("minwall")
					g.Z = 1
					g.Floats = true
					g.Wanders = false
					g.SpriteStack().YScale = 1.0
					g.SpriteStack().Shaded = true
					g.SetLevel(2)
					g.SetStats(5, 5, 10)
					g.SetAbility(&game.Ability{
						Name:     game.AbilityBlock,
						Tier:     2,
						Turns:    4,
						Cooldown: 3,
					})
				},
			},
			"3": {
				Actor: "glitch",
				OnCreate: func(s game.Actor) {
					g := s.(*actors.Glitch)
					g.SpriteStack().SetSprite("minion-shel")
					g.SetName("minshel")
					g.Z = 1
					g.Floats = true
					g.Wanders = false
					g.SpriteStack().YScale = 1.0
					g.SpriteStack().Shaded = true
					g.SetLevel(2)
					g.SetStats(5, 10, 5)
					g.SetAbility(&game.Ability{
						Name:     game.AbilityHardy,
						Tier:     1,
						Turns:    3,
						Cooldown: 4,
					})
				},
			},
		},
		metadata: make(map[string]interface{}),
		enter: func(w *game.World, r *game.Room) {
			if !first {
				return
			}
			first = false
			makeBigMsg := func(s string, d time.Duration, c color.NRGBA) game.Message {
				return game.Message{Text: s, Duration: d, Color: c, Font: &res.BigFont}
			}
			delayTimeR(1 * time.Second)
			clr := color.NRGBA{200, 64, 200, 255}
			<-w.MessageR(makeBigMsg("we are three with systems corrupted", 4000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("one piercing,", 2000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("one defending,", 2000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("one hardy", 2000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("defeat any", 2000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("and fix this place", 2000*time.Millisecond, clr))
			<-w.MessageR(makeBigMsg("...please", 1000*time.Millisecond, clr))
		},
		leave: func(w *game.World, r *game.Room) {
		},
		update: func(w *game.World, r *game.Room) {
		},
	}
}
