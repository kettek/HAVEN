package game

type World interface {
	Room() *Room
	Camera() *Camera
}
