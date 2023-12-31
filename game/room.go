package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Room struct {
	Tiles           [][]Tile
	color           color.NRGBA
	droppingIn      bool
	droppingInCount int
	iso             bool
	Transition      int
	OnUpdate        func()
	OnEnter         func()
	OnLeave         func()
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
			//r.Tiles[i][j].spriteStack = NewSpriteStack("floor")
			//r.Tiles[i][j].glitchion = 10
			//r.Tiles[i][j].ticker = rand.Intn(100)
			r.Tiles[i][j].Ticker = j*h + j
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
					if r.Tiles[i][j].SpriteStack == nil {
						continue
					}
					r.Tiles[i][j].SpriteStack.layerDistance = -1
					r.Tiles[i][j].SpriteStack.Alpha = 1
				}
			}
		} else {
			for i := range r.Tiles {
				for j := range r.Tiles[i] {
					if r.Tiles[i][j].SpriteStack == nil {
						continue
					}
					r.Tiles[i][j].SpriteStack.layerDistance = -1 + (1.0-float64(r.droppingInCount)/60)*-10
					r.Tiles[i][j].SpriteStack.Alpha = float32(r.droppingInCount) / 60
				}
			}
		}
	}

	if r.Transition > 0 {
		r.Transition--
		if r.Transition == 0 {
			r.iso = true
		}
	} else if r.Transition < 0 {
		r.Transition++
		if r.Transition == 0 {
			r.iso = false
		}
	}

	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			r.Tiles[i][j].Update()
		}
	}

	if r.OnUpdate != nil {
		r.OnUpdate()
	}

	return nil
}

func (r *Room) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	screen.Fill(r.color)
	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			if r.Tiles[i][j].SpriteStack == nil {
				continue
			}
			g := ebiten.GeoM{}
			if r.Transition != 0 {
				if r.Transition > 0 {
					ratio := float64(r.Transition) / 60
					x1, y1 := GetTilePosition(i, j)
					x2, y2 := GetTileIsoPosition(i, j)
					x, y := x1*ratio+x2*(1-ratio), y1*ratio+y2*(1-ratio)
					g.Translate(x, y)
					g.Concat(geom)
					r.Tiles[i][j].SpriteStack.DrawMixed(screen, g, ratio)
				} else if r.Transition < 0 {
					ratio := math.Abs(float64(r.Transition)) / 60
					x1, y1 := GetTilePosition(i, j)
					x2, y2 := GetTileIsoPosition(i, j)
					x, y := x1*(1-ratio)+x2*ratio, y1*(1-ratio)+y2*ratio
					g.Translate(x, y)
					g.Concat(geom)
					r.Tiles[i][j].SpriteStack.DrawMixed(screen, g, 1.0-ratio)
				}
			} else {
				if r.iso {
					g.Translate(GetTileIsoPosition(i, j))
					g.Concat(geom)
					r.Tiles[i][j].SpriteStack.DrawIso(screen, g)
				} else {
					g.Translate(GetTilePosition(i, j))
					g.Concat(geom)
					r.Tiles[i][j].SpriteStack.Draw(screen, g)
				}
			}
		}
	}
}

func (r *Room) Center() (float64, float64) {
	x, y := GetTilePosition(len(r.Tiles[0]), len(r.Tiles))
	return x / 2, y / 2
}

func (r *Room) CenterIso() (float64, float64) {
	x, y := GetTileIsoPosition(len(r.Tiles[0]), len(r.Tiles))
	return x / 2, y / 2
}
