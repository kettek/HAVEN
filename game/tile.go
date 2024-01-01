package game

import (
	"math"

	"github.com/kettek/ebihack23/res"
)

type Tile struct {
	Name        string
	SpriteStack *SpriteStack
	BlocksMove  bool
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
	return float64(x-y) * res.TileYStep, float64(x+y) * res.TileXStep
}

func GetTilePosition(x, y int) (float64, float64) {
	return float64(x) * res.TileWidth, float64(y) * res.TileHeight
}

func GetTileIsoPositionFromCoordinate(x, y float64) (int, int) {
	// This needs more slop, yo.
	return int(math.Round(x/res.TileYStep+y/res.TileXStep) / 2), int(math.Round((y/res.TileXStep - x/res.TileYStep) / 2))
}

func GetTilePositionFromCoordinate(x, y float64) (int, int) {
	return int(math.Round((x - res.TileHalfWidth) / res.TileWidth)), int(math.Round((y - res.TileHalfHeight) / res.TileHeight))
}
