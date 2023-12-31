package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Room struct {
	Tiles           [][]Tile
	Actors          []Actor
	PendingCommands []ActorCommand
	Color           color.NRGBA
	droppingIn      bool
	droppingInCount int
	iso             bool
	transition      int
	drawMode        DrawMode
	OnUpdate        func(*Room)
	OnEnter         func(*Room)
	OnLeave         func(*Room)
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
					r.Tiles[i][j].SpriteStack.LayerDistance = -1
					r.Tiles[i][j].SpriteStack.Alpha = 1
				}
			}
		} else {
			for i := range r.Tiles {
				for j := range r.Tiles[i] {
					if r.Tiles[i][j].SpriteStack == nil {
						continue
					}
					r.Tiles[i][j].SpriteStack.LayerDistance = -1 + (1.0-float64(r.droppingInCount)/60)*-10
					r.Tiles[i][j].SpriteStack.Alpha = float32(r.droppingInCount) / 60
				}
			}
		}
	}

	if r.drawMode == DrawModeFlatToIso {
		if r.transition > 0 {
			r.transition--
		} else {
			r.drawMode = DrawModeIso
		}
	} else if r.drawMode == DrawModeIsoToFlat {
		if r.transition > 0 {
			r.transition--
		} else {
			r.drawMode = DrawModeFlat
		}
	}

	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			r.Tiles[i][j].Update()
		}
	}

	for _, a := range r.Actors {
		if cmd := a.Update(r); cmd != nil {
			r.PendingCommands = append(r.PendingCommands,
				ActorCommand{
					Actor: a,
					Cmd:   cmd,
				})
		}
	}

	if r.OnUpdate != nil {
		r.OnUpdate(r)
	}

	return nil
}

func (r *Room) ToIso() {
	if r.drawMode == DrawModeIso {
		return
	}
	r.transition = 60
	r.drawMode = DrawModeFlatToIso
}

func (r *Room) ToFlat() {
	if r.drawMode == DrawModeFlat {
		return
	}
	r.transition = 60
	r.drawMode = DrawModeIsoToFlat
}

func (r *Room) GetTilePositionGeoM(x, y int) (g ebiten.GeoM, ratio float64) {
	if r.drawMode == DrawModeFlatToIso {
		ratio = float64(r.transition) / 60
		x1, y1 := GetTilePosition(x, y)
		x2, y2 := GetTileIsoPosition(x, y)
		x, y := x1*ratio+x2*(1-ratio), y1*ratio+y2*(1-ratio)
		g.Translate(x, y)
	} else if r.drawMode == DrawModeIsoToFlat {
		ratio = float64(r.transition) / 60
		x1, y1 := GetTilePosition(x, y)
		x2, y2 := GetTileIsoPosition(x, y)
		x, y := x1*(1-ratio)+x2*ratio, y1*(1-ratio)+y2*ratio
		g.Translate(x, y)
		ratio = 1.0 - ratio
	} else if r.drawMode == DrawModeIso {
		g.Translate(GetTileIsoPosition(x, y))
	} else {
		g.Translate(GetTilePosition(x, y))
	}
	return g, ratio
}

func (r *Room) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	screen.Fill(r.Color)
	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			if r.Tiles[i][j].SpriteStack == nil {
				continue
			}
			g, ratio := r.GetTilePositionGeoM(j, i)
			g.Concat(geom)
			r.Tiles[i][j].SpriteStack.Draw(screen, g, r.drawMode, ratio)
		}
	}
	for _, a := range r.Actors {
		x, y := a.Position()
		g, ratio := r.GetTilePositionGeoM(x, y)
		g.Concat(geom)
		a.Draw(screen, r, g, r.drawMode, ratio)
	}
}

func (r *Room) Size() (int, int) {
	return len(r.Tiles[0]), len(r.Tiles)
}

func (r *Room) Center() (float64, float64) {
	x, y := GetTilePosition(len(r.Tiles[0]), len(r.Tiles))
	return x / 2, y / 2
}

func (r *Room) CenterIso() (float64, float64) {
	x, y := GetTileIsoPosition(len(r.Tiles[0]), len(r.Tiles))
	return x / 2, y / 2
}

func (r *Room) GetTilePositionFromCoordinate(x, y float64) (int, int) {
	if r.drawMode == DrawModeIso {
		return GetTileIsoPositionFromCoordinate(x, y)
	} else {
		return GetTilePositionFromCoordinate(x, y)
	}
}

func (r *Room) GetTile(x, y int) *Tile {
	if x < 0 || y < 0 || y >= len(r.Tiles) || x >= len(r.Tiles[y]) {
		return nil
	}
	return &r.Tiles[y][x]
}
