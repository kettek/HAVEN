package rooms

import (
	"strings"

	"github.com/kettek/ebihack23/game"
)

type Room struct {
	tiles     string
	tileMap   map[string]string
	entities  string
	entityMap map[string]string
	metadata  map[string]interface{}
	enter     func()
	leave     func()
	update    func()
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

	return g
}

var data = map[string]Room{}
