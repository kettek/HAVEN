package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/res"
)

type SpriteStack struct {
	layers        []*ebiten.Image
	layerDistance float64
	alpha         float32
	rotation      float64
}

func NewSpriteStack(sprite string) *SpriteStack {
	ss := &SpriteStack{
		rotation:      math.Pi / 4,
		layerDistance: -1,
		alpha:         1.0,
	}

	layers, err := res.LoadSpriteStack(sprite)
	if err != nil {
		panic(err)
	}
	ss.layers = layers

	return ss
}

func (ss *SpriteStack) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(ss.rotation)
	op.GeoM.Translate(-res.TileHalfWidth, -res.TileHalfHeight)
	op.GeoM.Scale(1, 0.5)
	op.GeoM.Concat(geom)
	for i := 0; i < len(ss.layers); i++ {
		op.GeoM.Translate(0, ss.layerDistance)
		op.ColorScale.ScaleAlpha(ss.alpha)
		screen.DrawImage(ss.layers[i], op)
	}
}
