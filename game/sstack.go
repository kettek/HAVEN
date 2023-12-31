package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/res"
)

type SpriteStack struct {
	layers        []*ebiten.Image
	LayerDistance float64
	Alpha         float32
	Rotation      float64
}

func NewSpriteStack(sprite string) *SpriteStack {
	ss := &SpriteStack{
		LayerDistance: -1,
		Alpha:         1.0,
	}

	layers, err := res.LoadSpriteStack(sprite)
	if err != nil {
		panic(err)
	}
	ss.layers = layers

	return ss
}

func (ss *SpriteStack) IsoGeoM(geom ebiten.GeoM) ebiten.GeoM {
	geom.Rotate(math.Pi / 4)
	geom.Scale(1, 0.5)
	geom.Translate(-res.TileHalfWidth, -res.TileHalfHeight)
	geom.Rotate(ss.Rotation)
	geom.Translate(res.TileHalfWidth, res.TileHalfHeight)
	return geom
}

func (ss *SpriteStack) GeoM(geom ebiten.GeoM) ebiten.GeoM {
	geom.Translate(-res.TileHalfWidth, -res.TileHalfHeight)
	geom.Rotate(ss.Rotation)
	geom.Translate(res.TileHalfWidth, res.TileHalfHeight)
	return geom
}

func (ss *SpriteStack) DrawMixed(screen *ebiten.Image, geom ebiten.GeoM, ratio float64) {
	op := &ebiten.DrawImageOptions{}

	geom1 := ss.GeoM(ebiten.GeoM{})
	geom2 := ss.IsoGeoM(ebiten.GeoM{})
	//geom1.Scale(ratio, ratio)
	//geom2.Scale(1-ratio, 1-ratio)

	// Get entire matrix from geom1 and geom2 and multiply the elements by ratio.
	a := geom1.Element(0, 0)*ratio + geom2.Element(0, 0)*(1-ratio)
	b := geom1.Element(0, 1)*ratio + geom2.Element(0, 1)*(1-ratio)
	c := geom1.Element(0, 2)*ratio + geom2.Element(0, 2)*(1-ratio)
	d := geom1.Element(1, 0)*ratio + geom2.Element(1, 0)*(1-ratio)
	e := geom1.Element(1, 1)*ratio + geom2.Element(1, 1)*(1-ratio)
	f := geom1.Element(1, 2)*ratio + geom2.Element(1, 2)*(1-ratio)

	// And set our real op to that.
	op.GeoM.SetElement(0, 0, a)
	op.GeoM.SetElement(0, 1, b)
	op.GeoM.SetElement(0, 2, c)
	op.GeoM.SetElement(1, 0, d)
	op.GeoM.SetElement(1, 1, e)
	op.GeoM.SetElement(1, 2, f)

	op.GeoM.Concat(geom)

	for i := 0; i < len(ss.layers); i++ {
		op.GeoM.Translate(0, ss.LayerDistance*(1-ratio))
		op.ColorScale.ScaleAlpha(ss.Alpha)
		screen.DrawImage(ss.layers[i], op)
	}
}

func (ss *SpriteStack) Draw(screen *ebiten.Image, geom ebiten.GeoM, mode DrawMode, ratio float64) {
	if mode == DrawModeIsoToFlat {
		ss.DrawMixed(screen, geom, ratio)
	} else if mode == DrawModeFlatToIso {
		ss.DrawMixed(screen, geom, ratio)
	} else if mode == DrawModeIso {
		ss.DrawIso(screen, geom)
	} else {
		ss.DrawFlat(screen, geom)
	}
}

func (ss *SpriteStack) DrawFlat(screen *ebiten.Image, geom ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = ss.GeoM(op.GeoM)
	op.GeoM.Concat(geom)
	for i := 0; i < len(ss.layers); i++ {
		//op.GeoM.Translate(0, ss.LayerDistance)
		op.ColorScale.ScaleAlpha(ss.Alpha)
		screen.DrawImage(ss.layers[i], op)
	}
}

func (ss *SpriteStack) DrawIso(screen *ebiten.Image, geom ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = ss.IsoGeoM(op.GeoM)
	op.GeoM.Concat(geom)
	for i := 0; i < len(ss.layers); i++ {
		op.GeoM.Translate(0, ss.LayerDistance)
		op.ColorScale.ScaleAlpha(ss.Alpha)
		screen.DrawImage(ss.layers[i], op)
	}
}

type DrawMode int

const (
	DrawModeFlat DrawMode = iota
	DrawModeIso
	DrawModeIsoToFlat
	DrawModeFlatToIso
)
