package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/res"
	"github.com/tinne26/etxt"
)

type Room struct {
	Tiles           [][]Tile
	Actors          []Actor
	PendingCommands []ActorCommand
	TileMessages    []Message
	Color           color.NRGBA
	targetColor     color.NRGBA
	colorTicker     int
	iso             bool
	transition      int
	active          bool
	drawMode        DrawMode
	OnUpdate        func(*World, *Room)
	OnEnter         func(*World, *Room)
	OnLeave         func(*World, *Room)
	RoutineChan     chan func() bool
	RoutineChans    []func() bool
}

func NewRoom(w, h int) *Room {
	r := &Room{
		iso:         true,
		RoutineChan: make(chan func() bool),
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

func (r *Room) Activate() {
	r.active = true
}

func (r *Room) Update(w *World) []commands.Command {
	// Routine routines.
	for done := false; !done; {
		select {
		case fnc := <-r.RoutineChan:
			r.RoutineChans = append(r.RoutineChans, fnc)
		default:
			done = true
		}
	}
	routines := r.RoutineChans[:0]
	for _, ch := range r.RoutineChans {
		if !ch() {
			routines = append(routines, ch)
		}
	}
	r.RoutineChans = routines

	if r.colorTicker > 0 {
		ratio := float64(r.colorTicker) / 60

		r.Color.R = uint8(float64(r.Color.R)*ratio + float64(r.targetColor.R)*(1-ratio))
		r.Color.G = uint8(float64(r.Color.G)*ratio + float64(r.targetColor.G)*(1-ratio))
		r.Color.B = uint8(float64(r.Color.B)*ratio + float64(r.targetColor.B)*(1-ratio))
		r.Color.A = uint8(float64(r.Color.A)*ratio + float64(r.targetColor.A)*(1-ratio))
		r.colorTicker--
		if r.colorTicker == 0 {
			r.Color = r.targetColor
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

	// Routine small messages.
	messages := r.TileMessages[:0]
	for _, m := range r.TileMessages {
		if time.Since(m.start) < m.Duration {
			messages = append(messages, m)
		}
	}
	r.TileMessages = messages

	// Do not process actors or commands if we're not active.
	if !r.active {
		return nil
	}

	// Process actors.
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
		r.OnUpdate(w, r)
	}

	var results []commands.Command
	for _, cmd := range r.HandlePendingCommands(w) {
		switch c := cmd.(type) {
		default:
			results = append(results, c)
		}
	}

	return results
}

func (r *Room) HandlePendingCommands(w *World) (results []commands.Command) {
	for _, cmd := range r.PendingCommands {
		switch c := cmd.Cmd.(type) {
		case commands.Move:
			cmd.Actor.Command(commands.Face{X: c.X, Y: c.Y})
			ax, ay := cmd.Actor.Position()
			if ax-c.X >= -1 && ax-c.X <= 1 && ay-c.Y >= -1 && ay-c.Y <= 1 {
				// First check if an actor is there.
				if actor := r.GetActor(c.X, c.Y); actor != nil {
					r.TileMessage(Message{Text: "something is there", Duration: 1 * time.Second, Font: res.SmallFont, X: ax, Y: ay})
					if cmd := actor.Interact(w, r, cmd.Actor); cmd != nil {
						results = append(results, cmd)
					}
					continue
				}
				if tile := r.GetTile(c.X, c.Y); tile != nil {
					if tile.SpriteStack == nil {
						r.TileMessage(Message{Text: "the void gazes at you", Duration: 1 * time.Second, Font: res.SmallFont, X: ax, Y: ay})
					} else if !tile.BlocksMove {
						cmd.Actor.Command(c)
					} else {
						r.TileMessage(Message{Text: "the way is blocked", Duration: 1 * time.Second, Font: res.SmallFont, X: ax, Y: ay})
					}
				} else {
					r.TileMessage(Message{Text: "impossible", Duration: 1 * time.Second, Font: res.SmallFont, X: ax, Y: ay})
				}
			}
		case commands.Investigate:
			cmd.Actor.Command(commands.Face{X: c.X, Y: c.Y})
			ax, ay := cmd.Actor.Position()
			s := "feel"
			if ax-c.X < -1 || ax-c.X > 1 || ay-c.Y < -1 || ay-c.Y > 1 {
				s = "see"
			}
			if actor := r.GetActor(c.X, c.Y); actor != nil {
				r.TileMessage(Message{Text: fmt.Sprintf("i %s thing <%s>", s, actor.Name()), Duration: 1 * time.Second, Font: res.SmallFont, X: ax, Y: ay})
				continue
			}
			if tile := r.GetTile(c.X, c.Y); tile != nil && tile.Name != "" {
				r.TileMessage(Message{Text: fmt.Sprintf("i %s <%s>", s, tile.Name), Duration: 1 * time.Second, Font: res.SmallFont, X: ax, Y: ay})
				continue
			}
			r.TileMessage(Message{Text: fmt.Sprintf("i %s nil", s), Duration: 1 * time.Second, Font: res.SmallFont, X: ax, Y: ay})
		default:
			fmt.Println("handle", cmd.Actor, "wants to", cmd.Cmd)
		}
	}
	r.PendingCommands = nil
	return
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
		a.Draw(screen, r, geom, r.drawMode)
	}

	lastX, lastY := -1, -1
	res.Text.Utils().StoreState()
	res.Text.SetAlign(etxt.Center)
	for _, m := range r.TileMessages {
		g, _ := r.GetTilePositionGeoM(m.X, m.Y)
		if m.X == lastX && m.Y == lastY {
			g.Translate(0, -4)
		}
		lastX = m.X
		lastY = m.Y
		g.Concat(geom)
		gx := g.Element(0, 2)
		gy := g.Element(1, 2) - res.TileHeight*2 // Offset so the text doesn't appear right over the target
		res.Text.SetSize(12)
		res.Text.SetFont(m.Font)
		res.Text.SetColor(m.Color)
		res.Text.Draw(screen, m.Text, int(gx), int(gy))
	}
	res.Text.Utils().RestoreState()
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

func (r *Room) GetActor(x, y int) Actor {
	for _, a := range r.Actors {
		ax, ay := a.Position()
		if ax == x && ay == y {
			return a
		}
	}
	return nil
}

func (r *Room) SetColor(c color.NRGBA) {
	r.targetColor = c
	r.colorTicker = 60
}

func (r *Room) TileMessage(m Message) {
	m.start = time.Now()
	if m.Color.A == 0 {
		m.Color = color.NRGBA{255, 255, 255, 255}
	}
	if m.Font == nil {
		m.Font = res.Font
	}
	r.TileMessages = append(r.TileMessages, m)
}

func (r *Room) TileMessageR(m Message) chan bool {
	done := make(chan bool)
	m.start = time.Now()
	if m.Color.A == 0 {
		m.Color = color.NRGBA{0, 0, 0, 255}
	}
	if m.Font == nil {
		m.Font = res.Font
	}
	fnc := func() bool {
		r.TileMessages = append(r.TileMessages, m)
		done <- true
		return true
	}
	r.RoutineChan <- fnc
	return done
}

func (r *Room) FuncR(fnc func()) chan bool {
	done := make(chan bool)
	f := func() bool {
		fnc()
		done <- true
		return true
	}
	r.RoutineChan <- f
	return done
}

func (r *Room) DropInR() chan bool {
	done := make(chan bool)

	count := 0

	fnc := func() bool {
		count++
		if count >= 60 {
			for i := range r.Tiles {
				for j := range r.Tiles[i] {
					if r.Tiles[i][j].SpriteStack == nil {
						continue
					}
					r.Tiles[i][j].SpriteStack.LayerDistance = -1
					r.Tiles[i][j].SpriteStack.Alpha = 1
				}
			}
			for _, a := range r.Actors {
				if sp := a.SpriteStack(); sp != nil {
					sp.LayerDistance = -1
					sp.Alpha = 1
				}
			}
			done <- true
			return true
		} else {
			for i := range r.Tiles {
				for j := range r.Tiles[i] {
					if r.Tiles[i][j].SpriteStack == nil {
						continue
					}
					r.Tiles[i][j].SpriteStack.LayerDistance = -1 + (1.0-float64(count)/60)*-10
					r.Tiles[i][j].SpriteStack.Alpha = float32(count) / 60
				}
			}
			for _, a := range r.Actors {
				if sp := a.SpriteStack(); sp != nil {
					sp.LayerDistance = -1 + (1.0-float64(count)/60)*-10
					sp.Alpha = float32(count) / 60
				}
			}
			return false
		}
	}
	r.RoutineChan <- fnc
	return done
}
