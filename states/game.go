package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebihack23/actors"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/rooms"
	"github.com/kettek/ebihack23/settings"
)

type Game struct {
	camera           *game.Camera
	room             *game.Room
	cursorX, cursorY int
	hoveredTile      *game.Tile
	hoveredActor     game.Actor
	enterChan        chan struct{}
	processChan      chan func() bool
	entering         bool
}

func NewGame() *Game {
	return &Game{
		enterChan:   make(chan struct{}),
		processChan: make(chan func() bool),
	}
}

func (g *Game) Update() error {
	g.camera.Update()

	g.room.Update()

	if g.entering {
		select {
		case <-g.enterChan:
			g.entering = false
		default:
			return nil
		}
		fmt.Println("done")
	}

	cx, cy := ebiten.CursorPosition()
	x := float64(cx) / g.camera.Zoom
	y := float64(cy) / g.camera.Zoom
	if x >= 0 && x <= g.camera.W && y >= 0 && y <= g.camera.H {
		// Convert to world coordinates
		x += g.camera.X - g.camera.W/2
		y += g.camera.Y - g.camera.H/2
		px, py := g.room.GetTilePositionFromCoordinate(float64(x), float64(y))
		rw, rh := g.room.Size()
		if g.hoveredTile != nil {
			if g.hoveredTile.SpriteStack != nil {
				g.hoveredTile.SpriteStack.Highlight = false
			}
			g.hoveredTile = nil
		}
		if g.hoveredActor != nil {
			g.hoveredActor.Hover(false)
			g.hoveredActor = nil
		}
		if px >= 0 && px < rw && py >= 0 && py < rh {
			g.cursorX, g.cursorY = px, py
			if tile := g.room.GetTile(px, py); tile != nil {
				if tile.SpriteStack != nil {
					tile.SpriteStack.Highlight = true
				}
				g.hoveredTile = tile
			}
			if actor := g.room.GetActor(px, py); actor != nil {
				actor.Hover(true)
				g.hoveredActor = actor
			}
		} else {
			g.cursorX, g.cursorY = -1, -1
		}
	}

	if g.cursorX != -1 && g.cursorY != -1 {
		lmb := false
		rmb := false
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			lmb = true
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			rmb = true
		}
		if lmb || rmb {
			for _, actor := range g.room.Actors {
				switch actor := actor.(type) {
				case *actors.Player:
					if rmb {
						g.room.PendingCommands = append(g.room.PendingCommands, game.ActorCommand{
							Actor: actor,
							Cmd:   commands.Move{X: g.cursorX, Y: g.cursorY},
						})
					} else if lmb {
						g.room.PendingCommands = append(g.room.PendingCommands, game.ActorCommand{
							Actor: actor,
							Cmd:   commands.Investigate{X: g.cursorX, Y: g.cursorY},
						})
					}
				}
			}
		}
	}

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.room.ToIso()
		x, y := g.room.CenterIso()
		g.camera.MoveTo(x, y)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		g.room.ToFlat()
		x, y := g.room.Center()
		g.camera.MoveTo(x, y)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyM) {
		if settings.FilterMode == settings.MayoMode {
			settings.FilterMode = settings.ClarityMode
		} else {
			settings.FilterMode = settings.MayoMode
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyN) {
		settings.StackShading = !settings.StackShading
	}

	return nil
}

func (g *Game) Enter() {
	if g.camera == nil {
		g.camera = game.NewCamera()
	}
	if g.room == nil {
		g.room = rooms.BuildRoom("000_spawn")
		go func() {
			g.entering = true
			g.room.OnEnter(g.room)
			g.enterChan <- struct{}{}
		}()
	}
}
func (g *Game) Leave() {
}

func (g *Game) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}

	geom.Translate(-g.camera.W/2, -g.camera.H/2)
	geom.Rotate(g.camera.Rotation)
	geom.Translate(g.camera.W/2, g.camera.H/2)
	geom.Translate(-g.camera.X+g.camera.W/2, -g.camera.Y+g.camera.H/2)
	geom.Scale(g.camera.Zoom, g.camera.Zoom)
	g.room.Draw(screen, geom)
}
