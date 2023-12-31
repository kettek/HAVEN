package game

import (
	"math"

	"github.com/kettek/ebihack23/res"
)

type Tile struct {
	SpriteStack *SpriteStack
	Ticker      int
	Glitchion   float64
	Rotation    float64
}

func (t *Tile) Update() {
	t.Ticker++
	if t.Glitchion == 0 {
		t.Rotation = 0
	} else {
		t.Rotation = math.Cos(float64(t.Ticker)/t.Glitchion) / t.Glitchion
	}
	if t.SpriteStack != nil {
		t.SpriteStack.Rotation = t.Rotation
	}
}

func GetTileIsoPosition(x, y int) (float64, float64) {
	return float64(x)*res.TileStep - float64(y)*(res.TileStep), float64(y)*res.TileStep + float64(x)*(res.TileXStep) - (float64(y) * res.TileXStep)
}

func GetTilePosition(x, y int) (float64, float64) {
	return float64(x) * res.TileWidth, float64(y) * res.TileHeight
}
