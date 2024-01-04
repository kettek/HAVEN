package actors

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
)

type Player struct {
	X, Y             int
	movingTicker     int
	targetX, targetY int
	spriteStack      *game.SpriteStack
	onInteract       InteractFunc
}

func (p *Player) Command(cmd commands.Command) {
	switch cmd := cmd.(type) {
	case commands.Face:
		if cmd.X < p.X {
			p.spriteStack.Rotation = math.Pi * 3 / 2
		} else if cmd.X > p.X {
			p.spriteStack.Rotation = math.Pi / 2
		} else if cmd.Y < p.Y {
			p.spriteStack.Rotation = 0
		} else if cmd.Y > p.Y {
			p.spriteStack.Rotation = math.Pi
		}
	case commands.Move:
		if cmd.X < p.targetX {
			p.spriteStack.Rotation = math.Pi * 3 / 2
		} else if cmd.X > p.targetX {
			p.spriteStack.Rotation = math.Pi / 2
		} else if cmd.Y < p.targetY {
			p.spriteStack.Rotation = 0
		} else if cmd.Y > p.targetY {
			p.spriteStack.Rotation = math.Pi
		}
		p.movingTicker = 10
		p.targetX = cmd.X
		p.targetY = cmd.Y
	case commands.Investigate:
		if cmd.X < p.targetX {
			p.spriteStack.Rotation = math.Pi * 3 / 2
		} else if cmd.X > p.targetX {
			p.spriteStack.Rotation = math.Pi / 2
		} else if cmd.Y < p.targetY {
			p.spriteStack.Rotation = 0
		} else if cmd.Y > p.targetY {
			p.spriteStack.Rotation = math.Pi
		}
	}
}

func (p *Player) Update(room *game.Room) (cmd commands.Command) {
	if p.movingTicker > 0 {
		p.movingTicker--
		if p.movingTicker == 0 {
			p.X = p.targetX
			p.Y = p.targetY
		}
		return nil
	}

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
	return nil
}

func (p *Player) Draw(screen *ebiten.Image, r *game.Room, geom ebiten.GeoM, drawMode game.DrawMode) {
	var g ebiten.GeoM
	var ratio float64
	if p.movingTicker > 0 {
		moveRatio := float64(p.movingTicker) / 10
		g2, _ := r.GetTilePositionGeoM(p.X, p.Y)
		g1, _ := r.GetTilePositionGeoM(p.targetX, p.targetY)
		g.SetElement(0, 2, g1.Element(0, 2)*(1-moveRatio)+g2.Element(0, 2)*(moveRatio))
		g.SetElement(1, 2, g1.Element(1, 2)*(1-moveRatio)+g2.Element(1, 2)*(moveRatio))
	} else {
		g, ratio = r.GetTilePositionGeoM(p.X, p.Y)
	}

	g.Concat(geom)

	p.spriteStack.Draw(screen, g, drawMode, ratio)
}

func (p *Player) Position() (int, int) {
	return p.X, p.Y
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

func (p *Player) SpriteStack() *game.SpriteStack {
	return p.spriteStack
}

func (p *Player) Interact(w *game.World, r *game.Room, o game.Actor) commands.Command {
	if p.onInteract != nil {
		return p.onInteract(w, r, p, o)
	}
	return nil
}

func init() {
	actors["player"] = func(x, y int, ctor CreateFunc, interact InteractFunc) game.Actor {
		ss := game.NewSpriteStack("player")
		ss.Shaded = true
		ss.YScale = 1
		ss.LayerDistance = -1
		p := &Player{
			X:           x,
			Y:           y,
			spriteStack: ss,
			onInteract:  interact,
		}
		if ctor != nil {
			ctor(p)
		}
		return p
	}
}
