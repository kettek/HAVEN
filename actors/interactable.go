package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/inputs"
)

type Interactable struct {
	X, Y               int
	name               string
	tag                string
	moving             bool
	targetX, targetY   int
	pendingX, pendingY float64
	spriteStack        *game.SpriteStack
	onInteract         InteractFunc
}

func (p *Interactable) Command(cmd commands.Command) {
	switch cmd := cmd.(type) {
	case commands.Move:
		p.targetX = cmd.X
		p.targetY = cmd.Y
	}
}

func (p *Interactable) Update(room *game.Room) (cmd commands.Command) {
	// If active, check for... something.
	return nil
}

func (p *Interactable) Input(in inputs.Input) bool {
	return false
}

func (p *Interactable) Draw(screen *ebiten.Image, r *game.Room, geom ebiten.GeoM, drawMode game.DrawMode) {
	g, ratio := r.GetTilePositionGeoM(p.X, p.Y)
	g.Concat(geom)
	p.spriteStack.Draw(screen, g, drawMode, ratio)
}

func (p *Interactable) Position() (int, int) {
	return p.X, p.Y
}

func (p *Interactable) Hover(h bool) {
	if h {
		p.spriteStack.Highlight = true
	} else {
		p.spriteStack.Highlight = false
	}
}

func (p *Interactable) Hovered() bool {
	return p.spriteStack.Highlight
}

func (p *Interactable) SetName(s string) {
	p.name = s
}

func (p *Interactable) Name() string {
	return p.name
}

func (p *Interactable) SetTag(s string) {
	p.tag = s
}

func (p *Interactable) Tag() string {
	return p.tag
}

func (p *Interactable) SpriteStack() *game.SpriteStack {
	return p.spriteStack
}

func (p *Interactable) Interact(w *game.World, r *game.Room, o game.Actor) commands.Command {
	if p.onInteract != nil {
		return p.onInteract(w, r, p, o)
	}
	return nil
}

func init() {
	actors["interactable"] = func(x, y int, ctor CreateFunc, interact InteractFunc) game.Actor {
		ss := game.NewSpriteStack("terminal")
		//ss.Shaded = true
		ss.LayerDistance = -1
		t := &Interactable{
			X:           x,
			Y:           y,
			name:        "terminal",
			spriteStack: ss,
			onInteract:  interact,
		}
		if ctor != nil {
			ctor(t)
		}
		return t
	}
}
