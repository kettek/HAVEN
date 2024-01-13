package rooms

type TileDef struct {
	Name       string
	Sprite     string
	BlocksMove bool
	Rotation   float64
}

type TileDefs map[string]TileDef
