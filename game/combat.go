package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/inputs"
	"github.com/kettek/ebihack23/res"
	"github.com/tinne26/etxt"
)

type Combat struct {
	x, y          float64
	Attacker      CombatActor
	Defender      CombatActor
	attackerFloat float64
	defenderFloat float64
	selectedItem  int
	items         []CombatMenuItem
	turn          int
	doneCommand   commands.Command
	showGlitches  bool
	image         *ebiten.Image
}

type CombatMenuItem struct {
	Icon *ebiten.Image
	Text string
}

func NewCombat(w, h int, attacker, defender CombatActor) *Combat {
	c := &Combat{
		Attacker: attacker,
		Defender: defender,
		image:    ebiten.NewImage(w, h),
		items: []CombatMenuItem{
			{
				Icon: res.LoadImage("icon-attack"),
				Text: "ATTACK",
			},
			{
				Icon: res.LoadImage("icon-boost"),
				Text: "INCREASE STAT",
			},
			{
				Icon: res.LoadImage("icon-useGlitch"),
				Text: "USE GLITCH",
			},
			{
				Icon: res.LoadImage("icon-swapGlitch"),
				Text: "SWAP GLITCH",
			},
			{
				Icon: res.LoadImage("icon-escape"),
				Text: "FLEE",
			},
		},
	}
	c.Refresh()

	return c
}

func (c *Combat) Refresh() {
	c.image.Fill(color.NRGBA{66, 20, 66, 200})
	pt := c.image.Bounds().Size()

	vector.StrokeRect(c.image, 0, 0, float32(pt.X), float32(pt.Y), 4, color.NRGBA{245, 120, 245, 255}, true)
}

func (c *Combat) Update(w *World, r *Room) (cmd commands.Command) {
	c.attackerFloat += 0.025
	c.defenderFloat += 0.05
	if c.doneCommand != nil {
		return c.doneCommand
	}
	if c.showGlitches {
		c.showGlitches = false
		return commands.Prompt{
			Items:   []string{"Glitch 1", "Glitch 2", "Glitch 3"},
			Message: "Select a glitch to use.",
			Handler: func(i int, s string) bool {
				fmt.Println("got", i, s)
				return true
			},
		}
	}
	return nil
}

func (c *Combat) Input(in inputs.Input) {
	switch in := in.(type) {
	case inputs.Cancel:
	case inputs.Confirm:
		if c.items[c.selectedItem].Text == "FLEE" {
			c.doneCommand = commands.CombatResult{Fled: true}
		} else if c.items[c.selectedItem].Text == "SWAP GLITCH" {
			c.showGlitches = true
		}
	case inputs.Direction:
		if in.Y < 0 {
			c.selectedItem--
			if c.selectedItem < 0 {
				c.selectedItem = 0
			}
		} else if in.Y > 0 {
			c.selectedItem++
			if c.selectedItem > len(c.items)-1 {
				c.selectedItem = len(c.items) - 1
			}
		}
	case inputs.Click:
		x := in.X - c.x
		y := in.Y - c.y
		fmt.Println("click", x, y)
		// TODO
	}
	c.Refresh()
}

func (c *Combat) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	c.image.Clear()

	// combat area width
	cw := float64(c.image.Bounds().Size().X) - 200
	ch := float64(c.image.Bounds().Size().Y)

	vector.DrawFilledRect(c.image, 0, 0, float32(cw), float32(ch), color.NRGBA{66, 20, 66, 220}, true)
	//pt := c.image.Bounds().Size()

	vector.StrokeRect(c.image, 0, 0, float32(cw), float32(ch), 4, color.NRGBA{245, 120, 245, 255}, true)

	// Draw right menu.
	mx := cw + 10
	my := ch - 100
	mw := 190
	mh := 100
	{
		vector.DrawFilledRect(c.image, float32(mx), float32(my), float32(mw), float32(mh), color.NRGBA{66, 66, 60, 220}, true)
		vector.StrokeRect(c.image, float32(mx), float32(my), float32(mw), float32(mh), 4, color.NRGBA{245, 245, 220, 255}, true)
		mx += 6
		my += 4
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		res.Text.SetAlign(etxt.Top | etxt.Left)

		res.Text.SetColor(color.NRGBA{219, 86, 32, 200})
		for i, item := range c.items {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(mx), float64(my))
			c.image.DrawImage(item.Icon, op)
			s := item.Text
			if i == c.selectedItem {
				s = "   > " + s
			} else {
				s = "     " + s
			}
			res.Text.Draw(c.image, s, int(mx), int(my))
			my += float64(res.DefFont.Size)
		}

		res.Text.Utils().RestoreState()
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(geom)

	// Draw combat scene.
	{
		attacker := c.Attacker.(Actor)
		geom := ebiten.GeoM{}
		geom.Scale(12, 12)
		geom.Translate(32, float64(c.image.Bounds().Dy())-72)
		geom.Translate(math.Cos(c.attackerFloat)*4, math.Sin(c.attackerFloat)*4)
		r := attacker.SpriteStack().Rotation
		attacker.SpriteStack().Rotation = -math.Pi * 2
		attacker.SpriteStack().DrawIso(c.image, geom)
		attacker.SpriteStack().Rotation = r

		//mp, mf, mi := c.Attacker.MaxStats()
		cp, cf, ci := c.Attacker.CurrentStats()
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		x := 100
		y := c.image.Bounds().Dy() - 72
		res.Text.SetColor(color.NRGBA{50, 255, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("INTEGRITY   %d", ci), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(color.NRGBA{255, 50, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("FIREWALL    %d", cf), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(color.NRGBA{255, 255, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("PENETRATION %d", cp), x, y)
		res.Text.Utils().RestoreState()
	}
	{
		defender := c.Defender.(Actor)
		geom := ebiten.GeoM{}
		geom.Scale(8, 8)
		geom.Translate(float64(cw)-64, 32)
		if defender.SpriteStack().SkewY == 0 {
			geom.Translate(math.Sin(c.defenderFloat)*4, math.Cos(c.defenderFloat)*4)
		}
		r := defender.SpriteStack().Rotation
		d := defender.SpriteStack().LayerDistance
		defender.SpriteStack().Rotation = -math.Pi * 2
		defender.SpriteStack().LayerDistance = -2
		defender.SpriteStack().DrawIso(c.image, geom)
		defender.SpriteStack().Rotation = r
		defender.SpriteStack().LayerDistance = d

		//mp, mf, mi := c.Defender.MaxStats()
		cp, cf, ci := c.Defender.CurrentStats()
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		x := 16
		y := 32
		res.Text.Draw(c.image, fmt.Sprintf("LVL %d %s", c.Defender.Level(), defender.Name()), x, y)
		y += res.DefFont.Size * 2
		res.Text.SetColor(color.NRGBA{50, 255, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("INTEGRITY   %d", ci), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(color.NRGBA{255, 50, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("FIREWALL    %d", cf), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(color.NRGBA{255, 255, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("PENETRATION %d", cp), x, y)
		res.Text.Utils().RestoreState()
	}

	screen.DrawImage(c.image, op)
}
