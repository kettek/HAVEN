package states

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/inputs"
	"github.com/kettek/ebihack23/res"
	"github.com/kettek/ebihack23/rooms"
	"github.com/kettek/ebihack23/settings"
)

type Game struct {
	world            *game.World
	cursorX, cursorY int
	hoveredTile      *game.Tile
	hoveredActor     game.Actor
	processChan      chan func() bool
	Cheats           bool
	cheatEngine      *CheatEngine
}

func NewGame() *Game {
	g := &Game{
		processChan: make(chan func() bool),
		Cheats:      true,
		cheatEngine: NewCheatEngine(),
	}
	g.cheatEngine.AddCheat("NOCLIP", func(g *Game) {
		if g.world.PlayerActor != nil {
			g.world.PlayerActor.SetGhosting(!g.world.PlayerActor.Ghosting())
			fmt.Println("hey, you're not supposed to be a ghost!")
		}
	})
	g.cheatEngine.AddCheat("WEAK", func(g *Game) {
		if g.world.PlayerActor != nil {
			g.world.PlayerActor.(game.CombatActor).Penalize(100, 100, 100)
			fmt.Println("weak as puny gnat, you are")
		}
	})
	warpTo := func(r string) {
		fmt.Printf("warping to \"%s\", you dirty cheater\n", r)
		g.world.Room.RemoveActor(g.world.PlayerActor)
		room := rooms.GetRoom(r)
		room.AddActor(g.world.PlayerActor)
		g.world.PlayerActor.SetPosition(1, 1, 0)
		g.world.EnterRoom(room)
	}

	for _, m := range rooms.GetRoomNames() {
		code := "WARP"
		for _, c := range m {
			if c >= '0' && c <= '9' {
				code += "Digit" + string(c)
			} else if c == '_' {
				code += "Minus"
				break
			} else {
				code += strings.ToUpper(string(c))
				break
			}
		}
		func(m string) {
			g.cheatEngine.AddCheat(code, func(g *Game) {
				warpTo(m)
			})
		}(m)
	}
	return g
}

func (g *Game) Camera() *game.Camera {
	return g.world.Camera
}

func (g *Game) Room() *game.Room {
	return g.world.Room
}

func (g *Game) Update() error {
	if g.Cheats {
		g.cheatEngine.Update(g)
	}
	res.UpdateSounds()
	res.Jukebox.Update()
	// Get inputs.
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		g.world.Input(inputs.Cancel{})
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) || inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.world.Input(inputs.Confirm{})
	}
	{
		var x int
		var y int
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			y--
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			y++
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			x--
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			x++
		}
		if x != 0 || y != 0 {
			g.world.Input(inputs.Direction{X: x, Y: y, Mod: ebiten.IsKeyPressed(ebiten.KeyShift)})
		}
	}
	{
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			g.world.Input(inputs.Click{X: float64(x), Y: float64(y), Which: ebiten.MouseButtonLeft, Mod: ebiten.IsKeyPressed(ebiten.KeyShift)})
		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			g.world.Input(inputs.Click{X: float64(x), Y: float64(y), Which: ebiten.MouseButtonRight, Mod: ebiten.IsKeyPressed(ebiten.KeyShift)})
		}
	}

	g.world.Update()

	if g.Room() == nil {
		return nil
	}

	cx, cy := ebiten.CursorPosition()
	x := float64(cx) / g.Camera().Zoom
	y := float64(cy) / g.Camera().Zoom
	if x >= 0 && x <= g.Camera().W && y >= 0 && y <= g.Camera().H {
		// Convert to world coordinates
		x += g.Camera().X - g.Camera().W/2
		y += g.Camera().Y - g.Camera().H/2
		px, py := g.Room().GetTilePositionFromCoordinate(float64(x), float64(y))
		rw, rh := g.Room().Size()
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
			if tile := g.Room().GetTile(px, py); tile != nil {
				if tile.SpriteStack != nil {
					tile.SpriteStack.Highlight = true
				}
				g.hoveredTile = tile
			}
			if actor := g.Room().GetActor(px, py); actor != nil {
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
		if lmb {
			g.world.Input(inputs.MapClick{X: g.cursorX, Y: g.cursorY, Which: ebiten.MouseButtonLeft, Mod: ebiten.IsKeyPressed(ebiten.KeyShift)})
		} else if rmb {
			g.world.Input(inputs.MapClick{X: g.cursorX, Y: g.cursorY, Which: ebiten.MouseButtonRight, Mod: ebiten.IsKeyPressed(ebiten.KeyShift)})
		}
	}

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.Room().ToIso()
		x, y := g.Room().CenterIso()
		g.Camera().MoveTo(x, y)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		g.Room().ToFlat()
		x, y := g.Room().Center()
		g.Camera().MoveTo(x, y)
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
	if g.world == nil {
		g.world = game.NewWorld(rooms.GetRoom)
		g.world.EnterRoom(rooms.GetRoom("000_spawn"))
	}
}
func (g *Game) Leave() {
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Draw(screen)
}

type CheatEngine struct {
	cheats        map[string]func(g *Game)
	keys          []ebiten.Key
	maxCheatLen   int
	currentString string
	lastType      int
}

func NewCheatEngine() *CheatEngine {
	c := &CheatEngine{
		cheats: make(map[string]func(g *Game)),
		keys:   make([]ebiten.Key, 60),
	}
	return c
}

func (c *CheatEngine) AddCheat(str string, cb func(g *Game)) {
	c.cheats[str] = cb
	if len(str) > c.maxCheatLen {
		c.maxCheatLen = len(str)
	}
}

func (c *CheatEngine) Update(g *Game) {
	c.lastType++
	if c.lastType > 60 {
		c.lastType = 0
		c.currentString = ""
	}
	c.keys = c.keys[:0]
	c.keys = inpututil.AppendJustReleasedKeys(c.keys)
	for _, k := range c.keys {
		c.currentString += k.String()
		c.lastType = 0
	}
	var fnc func(g *Game)
	for str, cb := range c.cheats {
		if c.currentString == str {
			fnc = cb
			break
		}
	}
	if fnc != nil {
		fnc(g)
		c.currentString = ""
	} else if len(c.currentString) > c.maxCheatLen {
		c.currentString = ""
	}
}
