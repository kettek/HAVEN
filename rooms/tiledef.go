package rooms

type TileDef struct {
	Name       string
	Sprite     string
	BlocksMove bool
}

type TileDefs map[string]TileDef
