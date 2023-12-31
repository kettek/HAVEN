package game

type Rooms interface {
	BuildRoom(name string) *Room
}
