package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/ebihack23/commands"
	"github.com/kettek/ebihack23/inputs"
)

type Combat struct {
	x, y         float64
	Attacker     CombatActor
	Defender     CombatActor
	turn         int
	done         bool
	showGlitches bool
	image        *ebiten.Image
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
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(geom)
	screen.DrawImage(c.image, op)
}
