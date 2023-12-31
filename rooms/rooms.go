package rooms

import "github.com/kettek/ebihack23/game"

type rooms struct {
	room map[string]Room
}

func (r *rooms) BuildRoom(name string) *game.Room {
	room, ok := r.room[name]
	if !ok {
		return nil
	}
	gRoom := room.ToGameRoom()
	return gRoom
}

var Rooms = rooms{
	room: make(map[string]Room),
}
