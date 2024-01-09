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
)

type Combat struct {
	x, y          float64
	Attacker      CombatActor
	Defender      CombatActor
	attackerFloat float64
	defenderFloat float64
	turn          int
	done          bool
	showGlitches  bool
	image         *ebiten.Image
}

func NewCombat(w, h int, attacker, defender CombatActor) *Combat {
	c := &Combat{
		Attacker: attacker,
		Defender: defender,
		image:    ebiten.NewImage(w, h),
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
	if c.done {
		return commands.CombatResult{Winner: c.Attacker, Loser: c.Defender}
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
		c.done = true
	case inputs.Confirm:
		c.showGlitches = true
	case inputs.Direction:
		// TODO
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
	c.image.Fill(color.NRGBA{66, 20, 66, 200})
	pt := c.image.Bounds().Size()

	vector.StrokeRect(c.image, 0, 0, float32(pt.X), float32(pt.Y), 4, color.NRGBA{245, 120, 245, 255}, true)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(geom)

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

		mp, mf, mi := c.Attacker.MaxStats()
		cp, cf, ci := c.Attacker.CurrentStats()
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		x := 100
		y := c.image.Bounds().Dy() - 72
		res.Text.Draw(c.image, fmt.Sprintf("INTEGRITY   %d/%d", ci, mi), x, y)
		y += res.DefFont.Size
		res.Text.Draw(c.image, fmt.Sprintf("FIREWALL    %d/%d", cf, mf), x, y)
		y += res.DefFont.Size
		res.Text.Draw(c.image, fmt.Sprintf("PENETRATION %d/%d", cp, mp), x, y)
		res.Text.Utils().RestoreState()
	}
	{
		defender := c.Defender.(Actor)
		geom := ebiten.GeoM{}
		geom.Scale(8, 8)
		geom.Translate(float64(c.image.Bounds().Size().X)-64, 32)
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

		mp, mf, mi := c.Defender.MaxStats()
		cp, cf, ci := c.Defender.CurrentStats()
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		x := 16
		y := 32
		res.Text.Draw(c.image, defender.Name(), x, y)
		y += res.DefFont.Size * 2
		res.Text.Draw(c.image, fmt.Sprintf("INTEGRITY   %d/%d", ci, mi), x, y)
		y += res.DefFont.Size
		res.Text.Draw(c.image, fmt.Sprintf("FIREWALL    %d/%d", cf, mf), x, y)
		y += res.DefFont.Size
		res.Text.Draw(c.image, fmt.Sprintf("PENETRATION %d/%d", cp, mp), x, y)
		res.Text.Utils().RestoreState()
	}

	screen.DrawImage(c.image, op)
}
