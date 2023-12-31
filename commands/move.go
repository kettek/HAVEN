package commands

type Command interface {
}

type Move struct {
	X, Y int
}
