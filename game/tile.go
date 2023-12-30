package game

import (
	"math"

	"github.com/kettek/ebihack23/res"
)

type Tile struct {
	spriteStack *SpriteStack
	ticker      int
	glitchion   float64
	rotation    float64
}

func (t *Tile) Update() {
	t.ticker++
	if t.glitchion == 0 {
		t.rotation = 0
	} else {
		t.rotation = math.Cos(float64(t.ticker)/t.glitchion) / t.glitchion
	}
	t.spriteStack.rotation = t.rotation
}

func GetTileIsoPosition(x, y int) (float64, float64) {
	return float64(x)*res.TileStep - float64(y)*(res.TileStep), float64(y)*res.TileStep + float64(x)*(res.TileXStep) - (float64(y) * res.TileXStep)
}

func GetTilePosition(x, y int) (float64, float64) {
	return float64(x) * res.TileWidth, float64(y) * res.TileHeight
}
