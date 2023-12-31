package rooms

import "github.com/kettek/ebihack23/game"

func BuildRoom(name string) *game.Room {
	room, ok := rooms[name]
	if !ok {
		return nil
	}
	gRoom := room.ToGameRoom()
	return gRoom
}

var rooms = make(map[string]Room)
