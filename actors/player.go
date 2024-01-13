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
	Combat
	X, Y             int
	movingTicker     int
	targetX, targetY int
	spriteStack      *game.SpriteStack
	onInteract       InteractFunc
	pendingCommand   commands.Command
	ready            bool
	ghosting         bool
}

func (p *Player) Ready() bool {
	return p.ready && p.movingTicker == 0 && p.pendingCommand != nil
}

func (p *Player) SetReady(r bool) {
	p.ready = r
}

func (p *Player) TakeTurn() (cmd commands.Command) {
	if p.pendingCommand != nil {
		cmd = p.pendingCommand
		p.pendingCommand = nil
	}
	return cmd
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
	case commands.Step:
		res.PlaySound("step")
		x := p.X + cmd.X
		y := p.Y + cmd.Y
		if x < p.targetX {
			p.spriteStack.Rotation = math.Pi * 3 / 2
		} else if x > p.targetX {
			p.spriteStack.Rotation = math.Pi / 2
		} else if y < p.targetY {
			p.spriteStack.Rotation = 0
		} else if y > p.targetY {
			p.spriteStack.Rotation = math.Pi
		}
		p.movingTicker = 10
		p.targetX = x
		p.targetY = y
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

	return nil
}

func (p *Player) Input(in inputs.Input) bool {
	var cmd commands.Command
	switch in := in.(type) {
	case inputs.Direction:
		if in.Mod {
			cmd = commands.Investigate{X: p.X + in.X, Y: p.Y + in.Y}
		} else {
			cmd = commands.Step{X: in.X, Y: in.Y}
		}
	case inputs.MapClick:
		if in.Which == ebiten.MouseButtonLeft {
			cmd = commands.Investigate{X: in.X, Y: in.Y}
		} else if in.Which == ebiten.MouseButtonRight {
			var x, y int
			if in.X < p.X {
				x = -1
			} else if in.X > p.X {
				x = 1
			}
			if in.Y < p.Y {
				y = -1
			} else if in.Y > p.Y {
				y = 1
			}
			cmd = commands.Step{X: x, Y: y}
		}
	}
	if cmd != nil {
		p.pendingCommand = cmd
		p.SetReady(true)
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
func (p *Player) DrawPost(screen, post *ebiten.Image, r *game.Room, geom ebiten.GeoM, drawMode game.DrawMode) {
}

func (p *Player) SetPosition(x, y, z int) {
	p.X = x
	p.Y = y
}

func (p *Player) Position() (int, int, int) {
	return p.X, p.Y, 0
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

func (p *Player) SetName(s string) {
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
		if cmd := p.onInteract(w, r, p, o); cmd != nil {
			return cmd
		}
	}
	if _, ok := o.(*Glitch); ok {
		return commands.Combat{
			Defender: o,
			Attacker: p,
		}
	}
	return nil
}

func (p *Player) Blocks() bool {
	return true
}

func (p *Player) SetBlocks(b bool) {
}

func (p *Player) Ghosting() bool {
	return p.ghosting
}

func (p *Player) SetGhosting(b bool) {
	p.ghosting = b
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
