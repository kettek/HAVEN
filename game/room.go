package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/res"
)

type Room struct {
	Tiles           [][]Tile
	droppingIn      bool
	droppingInCount int
}

func NewRoom(w, h int) *Room {
	r := &Room{
		droppingIn: true,
	}

	r.Tiles = make([][]Tile, h)
	for i := range r.Tiles {
		r.Tiles[i] = make([]Tile, w)
		for j := range r.Tiles[i] {
			r.Tiles[i][j].spriteStack = NewSpriteStack("floor")
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

	return nil
}

func (r *Room) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	//rw := float64(len(r.Tiles[0])) * 9
	//rh := float64(len(r.Tiles))*9 + float64(len(r.Tiles[0]))*4.5
	//geom.Translate(-rw*2, -rh/2*2)
	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			g := ebiten.GeoM{}
			g.Translate(float64(j)*res.TileStep, float64(i)*res.TileStep+float64(j)*res.TileXStep)
			g.Concat(geom)
			r.Tiles[i][j].spriteStack.Draw(screen, g)
		}
	}
}
