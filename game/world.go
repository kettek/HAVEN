package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/ebihack23/res"
)

type World struct {
	Rooms        []*Room
	Room         *Room
	Camera       *Camera
	RoutineChan  chan func() bool
	RoutineChans []func() bool
	Messages     []Message
}

func NewWorld() *World {
	return &World{
		Camera:      NewCamera(),
		RoutineChan: make(chan func() bool),
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

	// Process camera.
	w.Camera.Update()

	// Process room.
	if w.Room != nil {
		for _, cmd := range w.Room.Update(w) {
			fmt.Println("handle room command", cmd)
		}
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}

	geom.Translate(-w.Camera.W/2, -w.Camera.H/2)
	geom.Rotate(w.Camera.Rotation)
	geom.Translate(w.Camera.W/2, w.Camera.H/2)
	geom.Translate(-w.Camera.X+w.Camera.W/2, -w.Camera.Y+w.Camera.H/2)
	geom.Scale(w.Camera.Zoom, w.Camera.Zoom)
	if w.Room != nil {
		w.Room.Draw(screen, geom)
	}

	if len(w.Messages) > 0 {
		m := w.Messages[0]
		text.Draw(screen, m.Text, m.Font, 16, 32, m.Color)
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
			w.Room = room
		})
		w.Room.OnEnter(w, w.Room)
		<-w.FuncR(func() {
			w.Room.Activate()
		})
	}()
}

func (w *World) MessageR(msg Message) chan bool {
	done := make(chan bool)

	msg.start = time.Now()
	msg.id = messageID
	if msg.Color.A == 0 {
		msg.Color = color.NRGBA{0, 0, 0, 255}
	}
	if msg.Font == nil {
		msg.Font = res.Font
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
