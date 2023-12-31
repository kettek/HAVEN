package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebihack23/game"
	"github.com/kettek/ebihack23/rooms"
)

type Game struct {
	camera            *game.Camera
	room              *game.Room
	hoveredTileSprite *game.SpriteStack
}

func (g *Game) Update() error {
	g.camera.Update()
	g.room.Update()

	//if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
	cx, cy := ebiten.CursorPosition()
	x := float64(cx) / g.camera.Zoom
	y := float64(cy) / g.camera.Zoom
	if x >= 0 && x <= g.camera.W && y >= 0 && y <= g.camera.H {
		// Convert to world coordinates
		x += g.camera.X - g.camera.W/2
		y += g.camera.Y - g.camera.H/2
		px, py := g.room.GetTilePositionFromCoordinate(float64(x), float64(y))
		rw, rh := g.room.Size()
		if g.hoveredTileSprite != nil {
			g.hoveredTileSprite.Highlight = false
		}
		if px >= 0 && px < rw && py >= 0 && py < rh {
			if tile := g.room.GetTile(px, py); tile != nil && tile.SpriteStack != nil {
				tile.SpriteStack.Highlight = true
				g.hoveredTileSprite = tile.SpriteStack
			}
		} else {
			if g.hoveredTileSprite != nil {
				g.hoveredTileSprite.Highlight = false
				g.hoveredTileSprite = nil
			}
		}
	}
	//}

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.room.ToIso()
		x, y := g.room.CenterIso()
		g.camera.MoveTo(x, y)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		g.room.ToFlat()
		x, y := g.room.Center()
		g.camera.MoveTo(x, y)
	}

	return nil
}

func (g *Game) Enter() {
	if g.camera == nil {
		g.camera = game.NewCamera()
	}
	if g.room == nil {
		g.room = rooms.BuildRoom("000_spawn")
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
