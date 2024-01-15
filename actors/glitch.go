package actors

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/inputs"
)

type WarpEffect struct {
	lifetime         int
	x, y             int
	w, h             int
	offsetX, offsetY float64
	image            *ebiten.Image
}

type Glitch struct {
	Combat
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
	ready            bool
	Target           game.Actor
	thinkTicker      int
	Skews            bool
	Floats           bool
	Wanders          bool
	ghosting         bool
	ability          *game.Ability
	warpEffects      []WarpEffect
}

func (g *Glitch) Ready() bool {
	return g.ready
}

func (g *Glitch) SetReady(r bool) {
	g.ready = r
}

func (g *Glitch) TakeTurn() (cmd commands.Command) {
	g.thinkTicker--
	if len(g.pendingCommands) > 0 {
		cmd = g.pendingCommands[0]
		g.pendingCommands = g.pendingCommands[1:]
	}
	return cmd
}

func (g *Glitch) Command(cmd commands.Command) {
	// ???
	switch cmd := cmd.(type) {
	case commands.Face:
		if cmd.X < g.X {
			g.spriteStack.Rotation = math.Pi / 2
		} else if cmd.X > g.X {
			g.spriteStack.Rotation = math.Pi * 3 / 2
		} else if cmd.Y < g.Y {
			g.spriteStack.Rotation = math.Pi
		} else if cmd.Y > g.Y {
			g.spriteStack.Rotation = 0
		}
	case commands.Step:
		x := g.X + cmd.X
		y := g.Y + cmd.Y
		if x < g.targetX {
			g.spriteStack.Rotation = math.Pi / 2
		} else if x > g.targetX {
			g.spriteStack.Rotation = math.Pi * 3 / 2
		} else if y < g.targetY {
			g.spriteStack.Rotation = math.Pi
		} else if y > g.targetY {
			g.spriteStack.Rotation = 0
		}
		g.movingTicker = 10
		g.targetX = x
		g.targetY = y
	}
}

func (g *Glitch) Update(room *game.Room) (cmd commands.Command) {
	if len(g.warpEffects) == 0 {
		for i := 0; i < 3; i++ {
			x1 := -1 - rand.Intn(2)
			y1 := -1 - rand.Intn(2)
			x2 := 1 + rand.Intn(2)
			y2 := 1 + rand.Intn(2)
			g.warpEffects = append(g.warpEffects, WarpEffect{
				lifetime: 10 + rand.Intn(30),
				x:        x1,
				y:        y1,
				w:        x2 - x1,
				h:        y2 - y1,
				offsetX:  float64(-12 + rand.Intn(24)),
				offsetY:  float64(-12 + rand.Intn(24)),
			})
		}
	}

	for i := range g.warpEffects {
		g.warpEffects[i].lifetime--
		if g.warpEffects[i].lifetime <= 0 {
			x1 := -1 - rand.Intn(2)
			y1 := -1 - rand.Intn(2)
			x2 := 1 + rand.Intn(2)
			y2 := 1 + rand.Intn(2)
			g.warpEffects[i].lifetime = rand.Intn(30) + 10
			g.warpEffects[i].x = x1
			g.warpEffects[i].y = y1
			g.warpEffects[i].w = x2 - x1
			g.warpEffects[i].h = y2 - y1
			g.warpEffects[i].offsetX = float64(-12 + rand.Intn(24))
			g.warpEffects[i].offsetY = float64(-12 + rand.Intn(24))
		}
	}

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

	// Do a little wandering if our brain is empty.
	if g.thinkTicker <= 0 {
		g.thinkTicker = 2
		if g.Target != nil {
			tx, ty, _ := g.Target.Position()
			xdir := 0
			if tx < g.X {
				xdir = -1
			} else if tx > g.X {
				xdir = 1
			}
			ydir := 0
			if ty < g.Y {
				ydir = -1
			} else if ty > g.Y {
				ydir = 1
			}
			if xdir != 0 && ydir != 0 {
				if rand.Intn(2) == 0 {
					xdir = 0
				} else {
					ydir = 0
				}
			}
			g.pendingCommands = append(g.pendingCommands, commands.Step{
				X: xdir,
				Y: ydir,
			})

		} else if g.Wanders {
			if len(g.pendingCommands) == 0 {
				x := rand.Intn(3) - 1
				y := rand.Intn(3) - 1
				if x != 0 && y != 0 {
					if rand.Intn(2) == 0 {
						x = 0
					} else {
						y = 0
					}
				}
				g.pendingCommands = append(g.pendingCommands, commands.Step{
					X: x,
					Y: y,
				})
			}
		}
	}

	return nil
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
		offsetY -= g2.Element(1, 2)
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
func (g *Glitch) DrawPost(screen, post *ebiten.Image, r *game.Room, geom ebiten.GeoM, drawMode game.DrawMode) {
	// Get our position in the world.
	tg, _ := r.GetTilePositionGeoM(g.X, g.Y)
	for _, warp := range g.warpEffects {
		var gg ebiten.GeoM
		var gg2 ebiten.GeoM
		gg.Concat(tg)
		gg2.Concat(tg)
		gg.Translate(float64(warp.x), float64(warp.y))
		gg2.Translate(float64(warp.x), float64(warp.y))
		gg.Concat(geom)
		x := gg.Element(0, 2)
		y := gg.Element(1, 2)
		warp.image = screen.SubImage(image.Rect(int(x), int(y), int(x)+warp.w, int(y)+warp.h)).(*ebiten.Image)
		gg2.Translate(warp.offsetX, warp.offsetY)
		gg2.Concat(geom)
		post.DrawImage(warp.image, &ebiten.DrawImageOptions{
			GeoM: gg2,
		})
	}
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
	g.tag = s
}

func (g *Glitch) Tag() string {
	return g.tag
}

func (g *Glitch) SpriteStack() *game.SpriteStack {
	return g.spriteStack
}

func (g *Glitch) Interact(w *game.World, r *game.Room, o game.Actor) commands.Command {
	if g.onInteract != nil {
		cmd := g.onInteract(w, r, g, o)
		if cmd != nil {
			return cmd
		}
	}
	if _, ok := o.(*Player); ok {
		return commands.Combat{
			Defender: g,
			Attacker: o,
		}
	}
	return nil
}

func (g *Glitch) Blocks() bool {
	return true
}

func (g *Glitch) SetBlocks(b bool) {
}

func (g *Glitch) Ghosting() bool {
	return g.ghosting
}

func (g *Glitch) SetGhosting(b bool) {
	g.ghosting = b
}

func (g *Glitch) Glitch() bool {
	return true
}

func (g *Glitch) SetAbility(a *game.Ability) {
	g.ability = a
}

func (g *Glitch) Ability() *game.Ability {
	return g.ability
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
			Wanders:     true,
		}
		if ctor != nil {
			ctor(p)
		}
		return p
	}
}
