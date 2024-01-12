package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/inputs"
	"github.com/kettek/ebihack23/res"
)

type World struct {
	PlayerActor      Actor
	Rooms            []*Room
	LastRoom         *Room
	Room             *Room
	Camera           *Camera
	RoutineChan      chan func() bool
	RoutineChans     []func() bool
	Messages         []Message
	Prompts          []*Prompt
	Combat           *Combat
	roomBuilder      func(string) *Room
	Color            color.NRGBA
	colorTicker      int
	postProcessImage *ebiten.Image
}

func NewWorld(roomBuilder func(string) *Room) *World {
	return &World{
		Camera:      NewCamera(),
		RoutineChan: make(chan func() bool),
		roomBuilder: roomBuilder,
	}
}

func (w *World) Update() {
	// Process routines.
	for {
		select {
		case fnc := <-w.RoutineChan:
			w.RoutineChans = append(w.RoutineChans, fnc)
		default:
			goto done
		}
	}
done:
	routines := w.RoutineChans[:0]
	for _, r := range w.RoutineChans {
		if !r() {
			routines = append(routines, r)
		}
	}
	w.RoutineChans = routines

	if len(w.Prompts) != 0 {
		w.Prompts[len(w.Prompts)-1].Update()
	}

	if w.Combat != nil {
		if c := w.Combat.Update(w, w.Room); c != nil {
			if cmd, ok := c.(commands.CombatResult); ok {
				fmt.Println("combat result", cmd.Winner, "won", cmd.Destroyed, cmd.ExpGained)
				w.Combat = nil
				res.Jukebox.Play(w.Room.Song)
			} else if cmd, ok := c.(commands.Prompt); ok {
				w.AddPrompt(cmd.Items, cmd.Message, cmd.Handler, cmd.ShowVersions)
			}
		}
	}

	if w.PlayerActor != nil {
		x, y, _ := w.PlayerActor.Position()
		geom, _ := w.Room.GetTilePositionGeoM(x, y)
		w.Camera.MoveTo(geom.Element(0, 2), geom.Element(1, 2))
	}

	// Process camera.
	w.Camera.Update()

	w.colorTicker++
	if w.LastRoom != nil && w.Room != nil && w.colorTicker < 30 {
		w.Color = color.NRGBA{
			uint8(float64(w.LastRoom.Color.R)*(1-float64(w.colorTicker)/30) + float64(w.Room.Color.R)*(float64(w.colorTicker)/30)),
			uint8(float64(w.LastRoom.Color.G)*(1-float64(w.colorTicker)/30) + float64(w.Room.Color.G)*(float64(w.colorTicker)/30)),
			uint8(float64(w.LastRoom.Color.B)*(1-float64(w.colorTicker)/30) + float64(w.Room.Color.B)*(float64(w.colorTicker)/30)),
			255,
		}
	} else if w.Room != nil {
		w.Color = w.Room.Color
	} else {
		w.Color = color.NRGBA{0, 0, 0, 255}
	}

	// Process room.
	if w.Room != nil {
		cmds := w.Room.Update(w)
		for _, cmd := range cmds {
			switch cmd := cmd.(type) {
			case commands.Prompt:
				w.AddPrompt(cmd.Items, "", cmd.Handler, cmd.ShowVersions)
			case commands.Travel:
				room := w.roomBuilder(cmd.Room)
				var targetActor Actor
				if cmd.Target != nil {
					targetActor = cmd.Target.(Actor)
				} else if w.PlayerActor != nil {
					targetActor = w.PlayerActor
				}
				var x, y int
				if actor := room.GetActorByTag(cmd.Tag); actor != nil {
					x, y, _ = actor.Position()
					x += cmd.OffsetX
					y += cmd.OffsetY
				}
				if targetActor != nil {
					w.Room.RemoveActor(targetActor)
					room.AddActor(targetActor)
					targetActor.SetPosition(x, y, 0)
				}
				w.EnterRoom(room)
			case commands.Combat:
				attacker := cmd.Attacker.(CombatActor)
				defender := cmd.Defender.(CombatActor)
				//w.Combat = NewCombat(384, 288, attacker, defender)
				w.Combat = NewCombat(500, 388, attacker, defender)
				res.Jukebox.Play("bad-health")
			default:
				fmt.Println("unhandled room->world command", cmd)
			}
		}
	}
}

func (w *World) Input(in inputs.Input) {
	if len(w.Prompts) > 0 {
		w.Prompts[len(w.Prompts)-1].Input(in)
	} else if w.Combat != nil {
		w.Combat.Input(in)
	} else {
		if !w.Room.Input(w, in) {
			// TODO: Send to camera?
		}
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	if w.postProcessImage == nil {
		w.postProcessImage = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	}
	screen.Fill(w.Color)

	geom := ebiten.GeoM{}

	geom.Translate(-w.Camera.W/2, -w.Camera.H/2)
	geom.Rotate(w.Camera.Rotation)
	geom.Translate(w.Camera.W/2, w.Camera.H/2)
	geom.Translate(-w.Camera.X+w.Camera.W/2, -w.Camera.Y+w.Camera.H/2)
	geom.Scale(w.Camera.Zoom, w.Camera.Zoom)
	if w.Room != nil {
		w.Room.Draw(screen, geom)
	}

	w.postProcessImage.Clear()
	if w.Room != nil {
		w.Room.DrawPost(screen, w.postProcessImage, geom)
	}
	screen.DrawImage(w.postProcessImage, nil)

	res.Text.Utils().StoreState()
	if len(w.Messages) > 0 {
		m := w.Messages[0]
		res.Text.SetColor(m.Color)
		res.Text.Draw(screen, m.Text, 16, 32)
	}
	res.Text.Utils().RestoreState()

	if w.Combat != nil {
		geom := ebiten.GeoM{}
		w.Combat.x = float64(screen.Bounds().Dx()/2) - float64(w.Combat.image.Bounds().Dx()/2)
		w.Combat.y = float64(screen.Bounds().Dy()/2) - float64(w.Combat.image.Bounds().Dy()/2)
		geom.Translate(w.Combat.x, w.Combat.y)
		w.Combat.Draw(screen, geom)
	}

	if len(w.Prompts) != 0 {
		geom := ebiten.GeoM{}
		prompt := w.Prompts[len(w.Prompts)-1]
		prompt.x = float64(screen.Bounds().Dx()/2) - float64(prompt.image.Bounds().Dx()/2)
		prompt.y = float64(screen.Bounds().Dy()/2) - float64(prompt.image.Bounds().Dy()/2)
		geom.Translate(prompt.x, prompt.y)
		prompt.Draw(screen, geom)
	}
}

func (w *World) EnterRoom(room *Room) {
	go func() {
		<-w.FuncR(func() {
			if w.Room != nil {
				w.Room.OnLeave(w, w.Room)
				w.Room.active = false
			}
		})
		<-w.FuncR(func() {
			w.LastRoom = w.Room
			w.Room = room
			if w.LastRoom != nil {
				w.Room.DrawMode = w.LastRoom.DrawMode
			}
			if w.Room.Song != "" {
				res.Jukebox.Play(w.Room.Song)
			}
			w.colorTicker = 0
			if w.PlayerActor != nil {
				x, y, _ := w.PlayerActor.Position()
				geom, _ := w.Room.GetTilePositionGeoM(x, y)
				w.Camera.SetPosition(geom.Element(0, 2), geom.Element(1, 2))
			}

		})
		w.Room.OnEnter(w, w.Room)
		<-w.FuncR(func() {
			w.Room.Activate()
		})
	}()
}

func (w *World) AddPrompt(items []string, msg string, cb func(int, string) bool, showExtra bool) {
	w.Prompts = append(w.Prompts, NewPrompt(320, 200, items, msg, func(i int, s string) bool {
		done := false
		if i == -1 {
			done = true
			cb(i, s)
		} else {
			done = cb(i, s)
		}
		if done {
			w.Prompts = w.Prompts[:len(w.Prompts)-1]
		}
		return done
	}, showExtra))
}

func (w *World) MessageR(msg Message) chan bool {
	done := make(chan bool)

	msg.start = time.Now()
	msg.id = messageID
	if msg.Color.A == 0 {
		msg.Color = color.NRGBA{0, 0, 0, 255}
	}
	if msg.Font == nil {
		msg.Font = &res.DefFont
	}
	messageID++

	first := true
	fnc := func() bool {
		if first {
			w.Messages = append(w.Messages, msg)
			first = false
		}
		if time.Since(msg.start) >= msg.Duration {
			w.Messages = w.Messages[1:]
			done <- true
			return true
		}
		return false
	}

	w.RoutineChan <- fnc
	return done
}

func (w *World) FuncR(fnc func()) chan bool {
	done := make(chan bool)
	f := func() bool {
		fnc()
		done <- true
		return true
	}
	w.RoutineChan <- f
	return done
}
