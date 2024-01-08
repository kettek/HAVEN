package actors

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/inputs"
	"github.com/kettek/ebihack23/res"
)

type Player struct {
	X, Y             int
	movingTicker     int
	targetX, targetY int
	spriteStack      *game.SpriteStack
	onInteract       InteractFunc
	pendingCommands  []commands.Command
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
		res.PlaySound("step")
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

	if len(p.pendingCommands) > 0 {
		cmd = p.pendingCommands[0]
		p.pendingCommands = p.pendingCommands[1:]
	}

	return cmd
}

func (p *Player) Input(in inputs.Input) bool {
	var cmd commands.Command
	switch in := in.(type) {
	case inputs.Direction:
		if in.Mod {
			cmd = commands.Investigate{X: p.X + in.X, Y: p.Y + in.Y}
		} else {
			cmd = commands.Move{X: p.X + in.X, Y: p.Y + in.Y}
		}
	}
	if cmd != nil {
		if len(p.pendingCommands) < 10 {
			p.pendingCommands = append(p.pendingCommands, cmd)
		}
		return true
	}
	return false
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

func (p *Player) SetPosition(x, y int) {
	p.X = x
	p.Y = y
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

func (p *Player) SetTag(s string) {
}

func (p *Player) Tag() string {
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
