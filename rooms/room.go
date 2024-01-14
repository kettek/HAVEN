package rooms

import (
	"image/color"
	"strings"

	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/game"
)

type Room struct {
	tiles     string
	tileDefs  TileDefs
	entities  string
	entityMap EntityDefs
	metadata  map[string]interface{}
	enter     func(w *game.World, r *game.Room)
	leave     func(w *game.World, r *game.Room)
	update    func(w *game.World, r *game.Room)
	turn      func(w *game.World, r *game.Room)
	name      string
	song      string
}

func (r *Room) ToGameRoom() *game.Room {
	lines := strings.Split(r.tiles, "\n")
	lines = lines[1:]
	width := 0
	for i, line := range lines {
		c := strings.TrimLeft(line, "\t")
		if len(c) > width {
			width = len(c)
		}
		lines[i] = c
	}
	height := len(lines) - 1
	g := game.NewRoom(width, height)
	g.OnUpdate = r.update
	g.OnEnter = r.enter
	g.OnLeave = r.leave
	g.OnTurn = r.turn
	g.Song = r.song
	g.Color = color.NRGBA{0, 0, 0, 255}
	g.Name = strings.ToUpper(r.name)

	for y, line := range lines {
		for x, char := range line {
			if char == ' ' || char == '\t' {
				continue
			}
			tileDef, ok := r.tileDefs[string(char)]
			if !ok {
				continue
			}
			g.Tiles[y][x].SpriteStack = game.NewSpriteStack(tileDef.Sprite)
			if g.Tiles[y][x].SpriteStack != nil {
				g.Tiles[y][x].SpriteStack.ExtraRotation = tileDef.Rotation
			}
			g.Tiles[y][x].BlocksMove = tileDef.BlocksMove
			g.Tiles[y][x].Name = tileDef.Name
		}
	}

	// Make them entities.
	lines = strings.Split(r.entities, "\n")
	lines = lines[1 : len(lines)-1]
	for y, line := range lines {
		l := strings.TrimLeft(line, "\t")
		for x, char := range l {
			if char == ' ' {
				continue
			}
			entity, ok := r.entityMap[string(char)]
			if !ok {
				continue
			}
			actor := actors.New(entity.Actor, x, y, entity.OnCreate, entity.OnInteract)
			if actor == nil {
				continue
			}
			if _, ok := actor.(*actors.Glitch); ok {
				g.MaxGlitches++
			}
			g.Actors = append(g.Actors, actor)
		}
	}

	return g
}
