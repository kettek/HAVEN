package actors

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/inputs"
)

type Glitch struct {
	name             string
	tag              string
	X, Y             int
	Z                int
	movingTicker     int
	targetX, targetY int
	shadow           *game.SpriteStack
	spriteStack      *game.SpriteStack
	onInteract       InteractFunc
	pendingCommands  []commands.Command
	warble           int
	Skews            bool
	Floats           bool
}

func (g *Glitch) Command(cmd commands.Command) {
	// ???
	switch cmd := cmd.(type) {
	case commands.Face:
		if cmd.X < g.X {
			g.spriteStack.Rotation = math.Pi * 3 / 2
		} else if cmd.X > g.X {
			g.spriteStack.Rotation = math.Pi / 2
		} else if cmd.Y < g.Y {
			g.spriteStack.Rotation = 0
		} else if cmd.Y > g.Y {
			g.spriteStack.Rotation = math.Pi
		}
	case commands.Move:
		//res.PlaySound("step") // splort?
		if cmd.X < g.targetX {
			g.spriteStack.Rotation = math.Pi * 3 / 2
		} else if cmd.X > g.targetX {
			g.spriteStack.Rotation = math.Pi / 2
		} else if cmd.Y < g.targetY {
			g.spriteStack.Rotation = 0
		} else if cmd.Y > g.targetY {
			g.spriteStack.Rotation = math.Pi
		}
		g.movingTicker = 10
		g.targetX = cmd.X
		g.targetY = cmd.Y
	}
}

func (g *Glitch) Update(room *game.Room) (cmd commands.Command) {
	g.warble++
	//g.spriteStack.Rotation += math.Sin(float64(g.warble)/100)/100 + math.Cos(float64(g.warble)/50)/50
	if g.Skews {
		g.spriteStack.SkewX = math.Sin(float64(g.warble)/600) / 600
		g.spriteStack.SkewY = math.Cos(float64(g.warble)/300) / 300
	}
	if g.movingTicker > 0 {
		g.movingTicker--
		if g.movingTicker == 0 {
			g.X = g.targetX
			g.Y = g.targetY
		}
		return nil
	}

	if len(g.pendingCommands) > 0 {
		cmd = g.pendingCommands[0]
		g.pendingCommands = g.pendingCommands[1:]
	}

	return cmd
}

func (g *Glitch) Input(in inputs.Input) bool {
	return false
}

func (g *Glitch) Draw(screen *ebiten.Image, r *game.Room, geom ebiten.GeoM, drawMode game.DrawMode) {
	var gg ebiten.GeoM
	var ratio float64
	var offsetY float64
	var offsetYRatio float64
	if g.Z > 0 && r.DrawMode != game.DrawModeFlat {
		gg, ratio := r.GetTilePositionGeoM(g.X, g.Y-2*g.Z)
		offsetY = gg.Element(1, 2)
		offsetYRatio = 1.0 - ratio
	}
	if g.movingTicker > 0 {
		moveRatio := float64(g.movingTicker) / 10
		g2, _ := r.GetTilePositionGeoM(g.X, g.Y)
		g1, _ := r.GetTilePositionGeoM(g.targetX, g.targetY)
		gg.SetElement(0, 2, g1.Element(0, 2)*(1-moveRatio)+g2.Element(0, 2)*(moveRatio))
		gg.SetElement(1, 2, g1.Element(1, 2)*(1-moveRatio)+g2.Element(1, 2)*(moveRatio))
		offsetY -= g2.Element(1, 2) * moveRatio
	} else {
		gg, ratio = r.GetTilePositionGeoM(g.X, g.Y)
		offsetY -= gg.Element(1, 2)
	}

	if g.Z > 0 {
		gg2 := ebiten.GeoM{}
		gg2.Concat(gg)
		gg2.Concat(geom)
		g.shadow.Draw(screen, gg2, drawMode, ratio)
	}

	if g.Z > 0 && r.DrawMode != game.DrawModeFlat {
		gg.Translate(0, offsetY*offsetYRatio)
	}

	if g.Floats {
		gg.Translate(0, math.Sin(float64(g.warble)/40)*3)
	}

	gg.Concat(geom)

	g.spriteStack.Draw(screen, gg, drawMode, ratio)
}

func (g *Glitch) SetPosition(x, y, z int) {
	g.X = x
	g.Y = y
	g.Z = z
}

func (g *Glitch) Position() (int, int, int) {
	return g.X, g.Y, g.Z
}

func (g *Glitch) Hover(h bool) {
	if h {
		g.spriteStack.Highlight = true
	} else {
		g.spriteStack.Highlight = false
	}
}

func (g *Glitch) Hovered() bool {
	return g.spriteStack.Highlight
}

func (g *Glitch) Name() string {
	return g.name
}
func (g *Glitch) SetName(s string) {
	g.name = s
}

func (g *Glitch) SetTag(s string) {
}

func (g *Glitch) Tag() string {
	return g.tag
}

func (g *Glitch) SpriteStack() *game.SpriteStack {
	return g.spriteStack
}

func (g *Glitch) Interact(w *game.World, r *game.Room, o game.Actor) commands.Command {
	if g.onInteract != nil {
		return g.onInteract(w, r, g, o)
	}
	return nil
}

func init() {
	actors["glitch"] = func(x, y int, ctor CreateFunc, interact InteractFunc) game.Actor {
		ss := game.NewSpriteStack("glitch")
		ss.Shaded = true
		ss.YScale = 0.5
		ss.LayerDistance = -1
		shadow := game.NewSpriteStack("shadow")
		p := &Glitch{
			X:           x,
			Y:           y,
			spriteStack: ss,
			shadow:      shadow,
			onInteract:  interact,
		}
		if ctor != nil {
			ctor(p)
		}
		return p
	}
}
