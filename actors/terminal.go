package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
)

type Terminal struct {
	X, Y               int
	moving             bool
	targetX, targetY   int
	pendingX, pendingY float64
	spriteStack        *game.SpriteStack
}

func (p *Terminal) Command(cmd commands.Command) {
	switch cmd := cmd.(type) {
	case commands.Move:
		p.targetX = cmd.X
		p.targetY = cmd.Y
	}
}

func (p *Terminal) Update(room *game.Room) (cmd commands.Command) {
	// If active, check for... something.
	return nil
}

func (p *Terminal) Draw(screen *ebiten.Image, r *game.Room, geom ebiten.GeoM, drawMode game.DrawMode) {
	g, ratio := r.GetTilePositionGeoM(p.X, p.Y)
	g.Concat(geom)
	p.spriteStack.Draw(screen, g, drawMode, ratio)
}

func (p *Terminal) Position() (int, int) {
	return p.X, p.Y
}

func (p *Terminal) SetPosition(x, y int) {
	p.X = x
	p.Y = y
	p.targetX = x
	p.targetY = y
}

func (p *Terminal) Hover(h bool) {
	if h {
		p.spriteStack.Highlight = true
	} else {
		p.spriteStack.Highlight = false
	}
}

func (p *Terminal) Hovered() bool {
	return p.spriteStack.Highlight
}

func (p *Terminal) Name() string {
	return "terminal"
}

func init() {
	actors["terminal"] = func(x, y int) game.Actor {
		ss := game.NewSpriteStack("terminal")
		ss.LayerDistance = -1
		return &Terminal{
			X:           x,
			Y:           y,
			spriteStack: ss,
		}
	}
}
