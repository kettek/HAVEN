package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/ebihack23/inputs"
	"github.com/kettek/ebihack23/res"
	"github.com/tinne26/etxt"
)

// Prompt system. It's kinda jank, but it works well enough for this project.
type Prompt struct {
	image    *ebiten.Image
	Message  string
	Items    []string
	Selected int
	cb       func(int, string) bool
}

func NewPrompt(w, h int, items []string, msg string, cb func(int, string) bool) *Prompt {
	p := &Prompt{
		image:    ebiten.NewImage(w, h),
		Message:  msg,
		Items:    items,
		Selected: 0,
		cb:       cb,
	}
	p.Refresh()
	return p
}

func (p *Prompt) Refresh() {
	p.image.Fill(color.NRGBA{66, 66, 60, 200})

	pt := p.image.Bounds().Size()

	vector.StrokeRect(p.image, 0, 0, float32(pt.X), float32(pt.Y), 4, color.NRGBA{245, 245, 220, 255}, true)

	x := 4
	y := 2
	res.Text.Utils().StoreState()
	res.Text.SetAlign(etxt.Left | etxt.Top)
	res.Text.SetSize(float64(res.DefFont.Size))
	res.Text.SetFont(res.DefFont.Font)

	msg := fmt.Sprintf("ebiOS %s\n", res.EbiOS)
	res.Text.SetColor(color.NRGBA{219, 86, 32, 200})
	res.Text.DrawWithWrap(p.image, msg, x, y, pt.X-8)
	y += res.Text.MeasureWithWrap(msg, pt.X-8).IntHeight()

	msg = p.Message + "\n"
	res.Text.SetColor(color.NRGBA{255, 255, 255, 200})
	res.Text.DrawWithWrap(p.image, msg, x, y, pt.X-8)
	y += res.Text.MeasureWithWrap(msg, pt.X-8).IntHeight()

	// Magic numbers... for now.
	if y < 50 {
		y = 50
	}

	res.Text.SetColor(color.NRGBA{0, 255, 44, 200})
	for i, item := range p.Items {
		s := item
		if p.Selected == i {
			s = "> " + s
		} else {
			s = "  " + s
		}
		res.Text.Draw(p.image, s, x, y)
		// Ugh, screw it.
		y += 16
	}
	res.Text.Utils().RestoreState()
}

func (p *Prompt) Update() {
}

func (p *Prompt) Input(in inputs.Input) {
	switch in := in.(type) {
	case inputs.Direction:
		if in.Y < 0 {
			p.Selected--
		}
		if in.Y > 0 {
			p.Selected++
		}
		if p.Selected < 0 {
			p.Selected = 0
		}
		if p.Selected >= len(p.Items) {
			p.Selected = len(p.Items) - 1
		}
	case inputs.Confirm:
		p.cb(p.Selected, p.Items[p.Selected])
	case inputs.Cancel:
		p.cb(-1, "")
	}
	p.Refresh()
}

func (p *Prompt) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(geom)
	screen.DrawImage(p.image, op)
}
