package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
)

type Player struct {
	X, Y               int
	moving             bool
	targetX, targetY   int
	pendingX, pendingY float64
	spriteStack        *game.SpriteStack
}

func (p *Player) Command(cmd commands.Command) {
	switch cmd := cmd.(type) {
	case commands.Move:
		p.targetX = cmd.X
		p.targetY = cmd.Y
	}
}

func (p *Player) Update(room *game.Room) (cmd commands.Command) {
	// FIXME: This isn't supposed to be here.
	var x, y int
	if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		x = -1
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		x = 1
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		y = -1
	} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		y = 1
	}
	if x != 0 || y != 0 {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			return commands.Investigate{X: p.X + x, Y: p.Y + y}
		} else {
			return commands.Move{X: p.X + x, Y: p.Y + y}
		}
	}

	if p.moving {
		p.X = p.targetX
		p.Y = p.targetY
		// TODO: increase pendingX/pendingY until targetX/targetY is reached.
		p.moving = false
	}
	return nil
}

func (p *Player) Draw(screen *ebiten.Image, r *game.Room, geom ebiten.GeoM, drawMode game.DrawMode, ratio float64) {
	p.spriteStack.Draw(screen, geom, drawMode, ratio)
}

func (p *Player) Position() (int, int) {
	return p.X, p.Y
}

func (p *Player) SetPosition(x, y int) {
	p.X = x
	p.Y = y
	p.targetX = x
	p.targetY = y
}

func (p *Player) Hover(h bool) {
	if h {
		p.spriteStack.Highlight = true
	} else {
		p.spriteStack.Highlight = false
	}
}

func (p *Player) Hovered() bool {
	return p.spriteStack.Highlight
}

func (p *Player) Name() string {
	return "player"
}

func init() {
	actors["player"] = func(x, y int) game.Actor {
		ss := game.NewSpriteStack("player")
		ss.LayerDistance = -1
		return &Player{
			X:           x,
			Y:           y,
			spriteStack: ss,
		}
	}
}
