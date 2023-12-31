package rooms

import (
	"strings"

	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/game"
)

type Room struct {
	tiles     string
	tileMap   map[string]string
	entities  string
	entityMap map[string]string
	metadata  map[string]interface{}
	enter     func(r *game.Room)
	leave     func(r *game.Room)
	update    func(r *game.Room)
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
	height := len(lines)
	g := game.NewRoom(width, height)
	g.OnUpdate = r.update
	g.OnEnter = r.enter
	g.OnLeave = r.leave

	for y, line := range lines {
		for x, char := range line {
			if char == ' ' || char == '\t' {
				continue
			}
			tile, ok := r.tileMap[string(char)]
			if !ok {
				continue
			}
			g.Tiles[y][x].SpriteStack = game.NewSpriteStack(tile)
		}
	}

	// Make them entities.
	lines = strings.Split(r.entities, "\n")
	lines = lines[1:]
	for y, line := range lines {
		rx := 1
		for _, char := range line {
			if char == '\t' {
				continue
			}
			if char == ' ' {
				continue
			}
			rx++
			entity, ok := r.entityMap[string(char)]
			if !ok {
				continue
			}
			actor := actors.New(entity, rx, y)
			if actor == nil {
				continue
			}
			g.Actors = append(g.Actors, actor)
		}
	}

	return g
}
