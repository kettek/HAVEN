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

func GetRoom(name string) *game.Room {
	room, ok := cachedRooms[name]
	if !ok {
		room = BuildRoom(name)
		cachedRooms[name] = room
	}
	return room
}

func GetRoomNames() (names []string) {
	for name := range rooms {
		names = append(names, name)
	}
	return
}

var rooms = make(map[string]Room)
var cachedRooms = make(map[string]*game.Room)
