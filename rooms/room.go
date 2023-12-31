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
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	height := len(lines)
	g := game.NewRoom(width, height)
	g.OnUpdate = r.update
	g.OnEnter = r.enter
	g.OnLeave = r.leave

	for y, line := range lines {
		for x, char := range line {
			if char == ' ' {
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
