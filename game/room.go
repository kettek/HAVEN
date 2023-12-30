package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Room struct {
	Tiles           [][]Tile
	color           color.NRGBA
	droppingIn      bool
	droppingInCount int
	iso             bool
	transition      int
}

func NewRoom(w, h int) *Room {
	r := &Room{
		droppingIn: true,
		iso:        true,
	}

	r.Tiles = make([][]Tile, h)
	for i := range r.Tiles {
		r.Tiles[i] = make([]Tile, w)
		for j := range r.Tiles[i] {
			r.Tiles[i][j].spriteStack = NewSpriteStack("floor")
			r.Tiles[i][j].glitchion = 10
			//r.Tiles[i][j].ticker = rand.Intn(100)
			r.Tiles[i][j].ticker = j*h + j
		}
	}

	return r
}

func (r *Room) Update() error {
	if r.droppingIn {
		r.droppingInCount++
		if r.droppingInCount >= 60 {
			r.droppingIn = false
			r.droppingInCount = 0
			for i := range r.Tiles {
				for j := range r.Tiles[i] {
					r.Tiles[i][j].spriteStack.layerDistance = -1
					r.Tiles[i][j].spriteStack.alpha = 1
				}
			}
		} else {
			for i := range r.Tiles {
				for j := range r.Tiles[i] {
					r.Tiles[i][j].spriteStack.layerDistance = -1 + (1.0-float64(r.droppingInCount)/60)*-10
					r.Tiles[i][j].spriteStack.alpha = float32(r.droppingInCount) / 60
				}
			}
		}
	}

	if r.transition > 0 {
		r.transition--
		if r.transition == 0 {
			r.iso = true
		}
	} else if r.transition < 0 {
		r.transition++
		if r.transition == 0 {
			r.iso = false
		}
	}

	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			r.Tiles[i][j].Update()
		}
	}

	// FIXME: Just a test.
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		r.transition = 60
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		r.transition = -60
	}

	return nil
}

func (r *Room) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	screen.Fill(r.color)
	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			g := ebiten.GeoM{}
			if r.transition != 0 {
				if r.transition > 0 {
					ratio := float64(r.transition) / 60
					x1, y1 := GetTilePosition(i, j)
					x2, y2 := GetTileIsoPosition(i, j)
					x, y := x1*ratio+x2*(1-ratio), y1*ratio+y2*(1-ratio)
					g.Translate(x, y)
					g.Concat(geom)
					r.Tiles[i][j].spriteStack.DrawMixed(screen, g, ratio)
				} else if r.transition < 0 {
					ratio := math.Abs(float64(r.transition)) / 60
					x1, y1 := GetTilePosition(i, j)
					x2, y2 := GetTileIsoPosition(i, j)
					x, y := x1*(1-ratio)+x2*ratio, y1*(1-ratio)+y2*ratio
					g.Translate(x, y)
					g.Concat(geom)
					r.Tiles[i][j].spriteStack.DrawMixed(screen, g, 1.0-ratio)
				}
			} else {
				if r.iso {
					g.Translate(GetTileIsoPosition(i, j))
					g.Concat(geom)
					r.Tiles[i][j].spriteStack.DrawIso(screen, g)
				} else {
					g.Translate(GetTilePosition(i, j))
					g.Concat(geom)
					r.Tiles[i][j].spriteStack.Draw(screen, g)
				}
			}
		}
	}
}
