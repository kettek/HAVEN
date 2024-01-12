package game

import (
	"fmt"
	"image/color"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/inputs"
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
	Darkness        float64
	colorTicker     int
	iso             bool
	transition      int
	active          bool
	DrawMode        DrawMode
	OnUpdate        func(*World, *Room)
	OnEnter         func(*World, *Room)
	OnLeave         func(*World, *Room)
	OnTurn          func(*World, *Room)
	RoutineChan     chan func() bool
	RoutineChans    []func() bool
	Song            string
	turn            int
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

	if r.DrawMode == DrawModeFlatToIso {
		if r.transition > 0 {
			r.transition--
		} else {
			r.DrawMode = DrawModeIso
		}
	} else if r.DrawMode == DrawModeIsoToFlat {
		if r.transition > 0 {
			r.transition--
		} else {
			r.DrawMode = DrawModeFlat
		}
	}

	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			r.Tiles[i][j].Update()
			if w.PlayerActor != nil {
				x, y, _ := w.PlayerActor.Position()
				if r.Tiles[i][j].SpriteStack != nil {
					r.Tiles[i][j].SpriteStack.Alpha = 1.0 - float32((x-j)*(x-j)+(y-i)*(y-i))/100*float32(r.Darkness)
				}
			}
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
			/*r.PendingCommands = append(r.PendingCommands,
			ActorCommand{
				Actor: a,
				Cmd:   cmd,
			})*/
		}
		if w.PlayerActor != nil {
			j, i, _ := a.Position()
			x, y, _ := w.PlayerActor.Position()
			if a.SpriteStack() != nil {
				a.SpriteStack().Alpha = 1.0 - float32((x-j)*(x-j)+(y-i)*(y-i))/100*float32(r.Darkness)
			}
		}
	}

	if w.PlayerActor != nil && w.PlayerActor.Ready() {
		for _, a := range r.Actors {
			a.SetReady(false)
			if cmd := a.TakeTurn(); cmd != nil {
				r.PendingCommands = append(r.PendingCommands,
					ActorCommand{
						Actor: a,
						Cmd:   cmd,
					})
			}
		}
		if r.OnTurn != nil {
			r.OnTurn(w, r)
		}
	}

	// Resort actors by their Z + X - Y position.
	if len(r.PendingCommands) > 0 {
		sort.Slice(r.Actors, func(i, j int) bool {
			x1, y1, z1 := r.Actors[i].Position()
			x2, y2, z2 := r.Actors[j].Position()
			return z1 < z2 || (z1 == z2 && x2-x1 > y1-y2)
		})
	}

	if r.OnUpdate != nil {
		r.OnUpdate(w, r)
	}

	var results []commands.Command
	if len(r.PendingCommands) > 0 {
		for _, cmd := range r.HandlePendingCommands(w) {
			switch c := cmd.(type) {
			default:
				results = append(results, c)
			}
		}
	}

	return results
}

func (r *Room) Input(w *World, in inputs.Input) bool {
	if w.PlayerActor != nil {
		if w.PlayerActor.Input(in) {
			return true
		}
	}
	return false
}

func (r *Room) HandlePendingCommands(w *World) (results []commands.Command) {

	// Resolve potential collisions first.
	var collidedActors []Actor
	var collisionResults []commands.Command
	for _, cmd := range r.PendingCommands {
		a := cmd.Actor
		bail := false
		for _, actor := range collidedActors {
			if actor == a {
				bail = true
				break
			}
		}
		if bail {
			continue
		}
		switch c := cmd.Cmd.(type) {
		case commands.Step:
			ax, ay, _ := cmd.Actor.Position()
			x := ax + c.X
			y := ay + c.Y
			// First check if an actor is there.
			if actor := r.GetActor(x, y); actor != nil && actor != cmd.Actor {
				if cmd := actor.Interact(w, r, cmd.Actor); cmd != nil {
					collisionResults = append(collisionResults, cmd)
					collidedActors = append(collidedActors, actor)
					collidedActors = append(collidedActors, a)
					x1, y1, _ := actor.Position()
					x2, y2, _ := a.Position()
					a.Command(commands.Face{X: x1, Y: y1})
					actor.Command(commands.Face{X: x2, Y: y2})
				} else if actor != nil && w.PlayerActor == a {
					var s string
					if actor.Name() == "" {
						s = "something"
					} else {
						s = fmt.Sprintf("<%s>", actor.Name())
					}
					r.TileMessage(Message{Text: fmt.Sprintf("%s is there...\n", s), Duration: 3 * time.Second, Font: &res.SmallFont, X: ax, Y: ay})
				}
			}
		}
	}
	results = append(results, collisionResults...)

	for _, cmd := range r.PendingCommands {
		a := cmd.Actor
		bail := false
		for _, actor := range collidedActors {
			if actor == a {
				bail = true
				break
			}
		}
		if bail {
			continue
		}
		switch c := cmd.Cmd.(type) {
		case commands.Step:
			cmd.Actor.Command(commands.Face{X: c.X, Y: c.Y})
			ax, ay, _ := cmd.Actor.Position()
			x := ax + c.X
			y := ay + c.Y
			// Check if our destination is blocked.
			if actor := r.GetActor(x, y); actor != nil && actor != cmd.Actor {
				if cmd := actor.Interact(w, r, cmd.Actor); cmd != nil {
					results = append(results, cmd)
					collidedActors = append(collidedActors, actor)
					collidedActors = append(collidedActors, a)
				} else if actor != nil && w.PlayerActor == a {
					var s string
					if actor.Name() == "" {
						s = "something"
					} else {
						s = fmt.Sprintf("<%s>", actor.Name())
					}
					r.TileMessage(Message{Text: fmt.Sprintf("%s is there...\n", s), Duration: 3 * time.Second, Font: &res.SmallFont, X: ax, Y: ay})
				}
			} else if tile := r.GetTile(x, y); tile != nil {
				if tile.SpriteStack == nil {
					if w.PlayerActor == a {
						r.TileMessage(Message{Text: "the void gazes at you", Duration: 1 * time.Second, Font: &res.SmallFont, X: ax, Y: ay})
					}
				} else if !tile.BlocksMove {
					cmd.Actor.Command(c)
				} else {
					if w.PlayerActor == a {
						r.TileMessage(Message{Text: "the way is blocked", Duration: 1 * time.Second, Font: &res.SmallFont, X: ax, Y: ay})
					}
					res.PlaySound("bump")
				}
			} else {
				if w.PlayerActor == a {
					r.TileMessage(Message{Text: "impossible", Duration: 1 * time.Second, Font: &res.SmallFont, X: ax, Y: ay})
				}
			}
		case commands.Investigate:
			cmd.Actor.Command(commands.Face{X: c.X, Y: c.Y})
			ax, ay, _ := cmd.Actor.Position()
			s := "feel"
			if ax-c.X < -1 || ax-c.X > 1 || ay-c.Y < -1 || ay-c.Y > 1 {
				s = "see"
			}
			if actor := r.GetActor(c.X, c.Y); actor != nil {
				r.TileMessage(Message{Text: fmt.Sprintf("i %s thing <%s>", s, actor.Name()), Duration: 1 * time.Second, Font: &res.SmallFont, X: ax, Y: ay})
				continue
			}
			if tile := r.GetTile(c.X, c.Y); tile != nil && tile.Name != "" {
				r.TileMessage(Message{Text: fmt.Sprintf("i %s <%s>", s, tile.Name), Duration: 1 * time.Second, Font: &res.SmallFont, X: ax, Y: ay})
				continue
			}
			r.TileMessage(Message{Text: fmt.Sprintf("i %s nil", s), Duration: 1 * time.Second, Font: &res.SmallFont, X: ax, Y: ay})
		default:
			fmt.Println("handle", cmd.Actor, "wants to", cmd.Cmd)
		}
	}
	r.PendingCommands = nil
	return
}

func (r *Room) ToIso() {
	if r.DrawMode == DrawModeIso {
		return
	}
	r.transition = 60
	r.DrawMode = DrawModeFlatToIso
}

func (r *Room) ToFlat() {
	if r.DrawMode == DrawModeFlat {
		return
	}
	r.transition = 60
	r.DrawMode = DrawModeIsoToFlat
}

func (r *Room) GetTilePositionGeoM(x, y int) (g ebiten.GeoM, ratio float64) {
	if r.DrawMode == DrawModeFlatToIso {
		ratio = float64(r.transition) / 60
		x1, y1 := GetTilePosition(x, y)
		x2, y2 := GetTileIsoPosition(x, y)
		x, y := x1*ratio+x2*(1-ratio), y1*ratio+y2*(1-ratio)
		g.Translate(x, y)
	} else if r.DrawMode == DrawModeIsoToFlat {
		ratio = float64(r.transition) / 60
		x1, y1 := GetTilePosition(x, y)
		x2, y2 := GetTileIsoPosition(x, y)
		x, y := x1*(1-ratio)+x2*ratio, y1*(1-ratio)+y2*ratio
		g.Translate(x, y)
		ratio = 1.0 - ratio
	} else if r.DrawMode == DrawModeIso {
		g.Translate(GetTileIsoPosition(x, y))
	} else {
		g.Translate(GetTilePosition(x, y))
	}
	return g, ratio
}

func (r *Room) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	for i := range r.Tiles {
		for j := range r.Tiles[i] {
			if r.Tiles[i][j].SpriteStack == nil {
				continue
			}
			g, ratio := r.GetTilePositionGeoM(j, i)
			g.Concat(geom)
			r.Tiles[i][j].SpriteStack.Draw(screen, g, r.DrawMode, ratio)
		}
	}
	for _, a := range r.Actors {
		a.Draw(screen, r, geom, r.DrawMode)
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
		res.Text.SetSize(float64(m.Font.Size))
		res.Text.SetFont(m.Font.Font)
		res.Text.SetColor(m.Color)
		res.Text.Draw(screen, m.Text, int(gx), int(gy))
	}
	res.Text.Utils().RestoreState()
}

func (r *Room) DrawPost(screen, post *ebiten.Image, geom ebiten.GeoM) {
	for _, a := range r.Actors {
		a.DrawPost(screen, post, r, geom, r.DrawMode)
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
	if r.DrawMode == DrawModeIso {
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
		ax, ay, _ := a.Position()
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
		m.Color = color.NRGBA{255 - r.Color.R, 255 - r.Color.G, 255 - r.Color.B, 255}
	}
	if m.Font == nil {
		m.Font = &res.DefFont
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
		m.Font = &res.DefFont
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

func (r *Room) GetActorByTag(tag string) Actor {
	for _, a := range r.Actors {
		if a.Tag() == tag {
			return a
		}
	}
	return nil
}

func (r *Room) RemoveActor(a Actor) {
	for i, actor := range r.Actors {
		if actor == a {
			r.Actors = append(r.Actors[:i], r.Actors[i+1:]...)
			return
		}
	}
}

func (r *Room) AddActor(a Actor) {
	r.Actors = append(r.Actors, a)
}
