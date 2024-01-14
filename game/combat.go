package game

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

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
	menu          *CombatMenu
	menus         CombatMenus
	turn          int
	doneCommand   commands.Command
	showGlitches  bool
	image         *ebiten.Image
	isEnemyTurn   bool
	menuMode      CombatMenuMode
	action        CombatAction
	report        []CombatLine
}

type CombatLine struct {
	icon  *ebiten.Image
	text  string
	color color.NRGBA
}

var attackColor = color.NRGBA{255, 50, 50, 200}
var defenseColor = color.NRGBA{50, 255, 50, 200}
var neutralColor = color.NRGBA{200, 200, 200, 200}
var infoColor = color.NRGBA{255, 255, 50, 200}
var importantColor = color.NRGBA{255, 50, 255, 200}

type CombatAction interface {
	Done(c *Combat) bool
	Update(c *Combat)
	IsAttacker() bool
}

type CombatActionAttack struct {
	stat       string
	isAttacker bool
	timer      int
}

func (c CombatActionAttack) Done(cmb *Combat) bool {
	if c.timer >= 120 {
		var defender CombatActor
		if c.isAttacker {
			defender = cmb.Defender
		} else {
			defender = cmb.Attacker
		}
		p, f, i := defender.CurrentStats()
		if p <= 0 && c.stat == "PENETRATION" {
			cmb.AddReport(fmt.Sprintf("%s's penetration is down!", defender.Name()), nil, neutralColor)
		} else if f <= 0 && c.stat == "FIREWALL" {
			cmb.AddReport(fmt.Sprintf("%s's firewall is down!", defender.Name()), nil, neutralColor)
		} else if i <= 0 && c.stat == "INTEGRITY" {
			cmb.AddReport(fmt.Sprintf("%s's integrity is down!", defender.Name()), nil, neutralColor)
			cmb.AddReport(fmt.Sprintf("next attack will destroy %s", defender.Name()), res.LoadImage("icon-exclamation"), importantColor)
		}
		return true
	}
	return false
}

func (c *CombatActionAttack) Update(cmb *Combat) {
	c.timer++
	if c.timer == 10 {
		icon := res.LoadImage("icon-attack")
		/*if c.stat == "INTEGRITY" {
			icon = res.LoadImage("icon-integrity")
		} else if c.stat == "FIREWALL" {
			icon = res.LoadImage("icon-firewall")
		} else if c.stat == "PENETRATION" {
			icon = res.LoadImage("icon-penetration")
		}*/
		var attacker, defender CombatActor
		if c.isAttacker {
			attacker = cmb.Attacker
			defender = cmb.Defender
		} else {
			attacker = cmb.Defender
			defender = cmb.Attacker
		}
		v := attacker.RollAttack()
		if v <= 0 {
			cmb.AddReport(fmt.Sprintf("%s attacks %s, but misses!", attacker.Name(), c.stat), icon, infoColor)
			return
		}
		if c.stat == "INTEGRITY" {
			_, _, v = defender.ApplyDamage(0, 0, v)
		} else if c.stat == "FIREWALL" {
			_, v, _ = defender.ApplyDamage(0, v, 0)
		} else if c.stat == "PENETRATION" {
			v, _, _ = defender.ApplyDamage(v, 0, 0)
		}
		if v <= 0 {
			cmb.AddReport(fmt.Sprintf("%s attacks %s, but is denied!", attacker.Name(), c.stat), icon, infoColor)
			return
		}
		cmb.AddReport(fmt.Sprintf("%s attacks %s for %d!", attacker.Name(), c.stat, v), icon, attackColor)
	}
}

func (c CombatActionAttack) IsAttacker() bool {
	return c.isAttacker
}

type CombatActionBoost struct {
	stat       string
	isAttacker bool
	timer      int
}

func (c CombatActionBoost) Done(cmb *Combat) bool {
	if c.timer >= 120 {
		return true
	}
	return false
}

func (c *CombatActionBoost) Update(cmb *Combat) {
	c.timer++
	if c.timer == 10 {
		icon := res.LoadImage("icon-boost")
		/*if c.stat == "INTEGRITY" {
			icon = res.LoadImage("icon-integrity")
		} else if c.stat == "FIREWALL" {
			icon = res.LoadImage("icon-firewall")
		} else if c.stat == "PENETRATION" {
			icon = res.LoadImage("icon-penetration")
		}*/

		var attacker CombatActor
		if c.isAttacker {
			attacker = cmb.Attacker
		} else {
			attacker = cmb.Defender
		}
		p, f, i := attacker.RollBoost()
		var v int
		if c.stat == "INTEGRITY" {
			_, _, v = attacker.ApplyBoost(0, 0, i)
		} else if c.stat == "FIREWALL" {
			_, v, _ = attacker.ApplyBoost(0, f, 0)
		} else if c.stat == "PENETRATION" {
			v, _, _ = attacker.ApplyBoost(p, 0, 0)
		}
		cmb.AddReport(fmt.Sprintf("%s boosts %s for %d!", attacker.Name(), c.stat, v), icon, defenseColor)
	}
}

func (c CombatActionBoost) IsAttacker() bool {
	return c.isAttacker
}

type CombatActionFlee struct {
	isAttacker bool
	canFlee    bool
	timer      int
}

func (c CombatActionFlee) Done(cmb *Combat) bool {
	return c.timer >= 120
}

func (c *CombatActionFlee) Update(cmb *Combat) {
	c.timer++
	if c.timer == 60 {
		if c.canFlee {
			cmb.AddReport(fmt.Sprintf("%s flees successfully!", cmb.Attacker.Name()), nil, neutralColor)
		} else {
			cmb.AddReport("escape is denied!", nil, neutralColor)
		}
	}
}

func (c CombatActionFlee) IsAttacker() bool {
	return c.isAttacker
}

type CombatMenu struct {
	items         []CombatMenuItem
	selectedIndex int
}

type CombatMenus struct {
	main, attack, boost, swap, use CombatMenu
}

type CombatMenuMode int

const (
	CombatMenuModeMain CombatMenuMode = iota
	CombatMenuModeAttackStat
	CombatMenuModeBoostStat
	CombatMenuModeUseGlitch
	CombatMenuModeSwapGlitch
)

type CombatMenuItem struct {
	Icon     *ebiten.Image
	SubIcon  *ebiten.Image
	Text     string
	Bounds   image.Rectangle
	Disabled bool
	Trigger  func()
}

func (c *Combat) SwapMenu(mode CombatMenuMode) {
	c.menuMode = mode
	switch mode {
	case CombatMenuModeMain:
		c.menu = &c.menus.main
	case CombatMenuModeAttackStat:
		c.menu = &c.menus.attack
	case CombatMenuModeBoostStat:
		c.menu = &c.menus.boost
	case CombatMenuModeSwapGlitch:
		c.menu = &c.menus.swap
	case CombatMenuModeUseGlitch:
		c.menu = &c.menus.use
	}
}

func (c *Combat) SetAction(action CombatAction) {
	c.action = action
}

func (c *Combat) AddReport(text string, icon *ebiten.Image, color color.NRGBA) {
	c.report = append(c.report, CombatLine{
		icon:  icon,
		text:  text,
		color: color,
	})

	res.Text.Utils().StoreState()
	res.Text.SetSize(float64(res.DefFont.Size))
	res.Text.SetFont(res.DefFont.Font)
	res.Text.SetAlign(etxt.Top | etxt.Left)

	// This is incorrect.
	y := 4
	cull := 0
	for i := len(c.report) - 1; i >= 0; i-- {
		line := c.report[i]
		t := line.text
		if line.icon != nil {
			t = "  " + t
		}
		y += res.Text.MeasureWithWrap(t, 190-26).IntHeight()
		if y > 240 {
			cull++
		}
	}
	c.report = c.report[cull:]
	res.Text.Utils().RestoreState()
}

func NewCombat(w, h int, attacker, defender CombatActor) *Combat {
	var c *Combat
	c = &Combat{
		Attacker: attacker,
		Defender: defender,
		image:    ebiten.NewImage(w, h),
		menus: CombatMenus{
			main: CombatMenu{
				items: []CombatMenuItem{
					{
						Icon: res.LoadImage("icon-attack"),
						Text: "ATTACK",
						Trigger: func() {
							c.SwapMenu(CombatMenuModeAttackStat)
						},
					},
					{
						Icon: res.LoadImage("icon-boost"),
						Text: "BOOST STAT",
						Trigger: func() {
							c.SwapMenu(CombatMenuModeBoostStat)
						},
					},
					{
						Icon: res.LoadImage("icon-useGlitch"),
						Text: "USE GLITCH",
						Trigger: func() {
							c.SwapMenu(CombatMenuModeUseGlitch)
						},
						Disabled: !attacker.HasGlitch(),
					},
					{
						Icon: res.LoadImage("icon-swapGlitch"),
						Text: "SWAP GLITCH",
						Trigger: func() {
							c.SwapMenu(CombatMenuModeSwapGlitch)
						},
						Disabled: !attacker.HasGlitch(),
					},
					{
						Icon: res.LoadImage("icon-captureGlitch"),
						Text: "CAPTURE GLITCH",
						Trigger: func() {
							// TODO: sort of an attack.
						},
						Disabled: !attacker.HasGlitch(),
					},
					{
						Icon: res.LoadImage("icon-escape"),
						Text: "FLEE",
						Trigger: func() {
							// We pre-roll the flee results.
							_, f, i := c.Defender.CurrentStats()
							a := &CombatActionFlee{isAttacker: true}
							r := rand.Intn(100 + f + i)
							if r < 50 {
								a.canFlee = true
							}
							c.SetAction(a)
							c.AddReport(
								fmt.Sprintf("%s attempts to flee!", c.Attacker.Name()),
								res.LoadImage("icon-escape"),
								neutralColor,
							)
						},
					},
				},
			},
			attack: CombatMenu{
				items: []CombatMenuItem{
					{
						Icon:    res.LoadImage("icon-integrity"),
						SubIcon: res.LoadImage("subicon-attack"),
						Text:    "INTEGRITY",
						Trigger: func() {
							c.SetAction(&CombatActionAttack{stat: "INTEGRITY", isAttacker: true})
						},
					},
					{
						Icon:    res.LoadImage("icon-firewall"),
						SubIcon: res.LoadImage("subicon-attack"),
						Text:    "FIREWALL",
						Trigger: func() {
							c.SetAction(&CombatActionAttack{stat: "FIREWALL", isAttacker: true})
						},
					},
					{
						Icon:    res.LoadImage("icon-penetration"),
						SubIcon: res.LoadImage("subicon-attack"),
						Text:    "PENETRATION",
						Trigger: func() {
							c.SetAction(&CombatActionAttack{stat: "PENETRATION", isAttacker: true})
						},
					},
					{
						Text: "CANCEL",
						Trigger: func() {
							c.SwapMenu(CombatMenuModeMain)
						},
					},
				},
			},
			boost: CombatMenu{
				items: []CombatMenuItem{
					{
						Icon:    res.LoadImage("icon-integrity"),
						SubIcon: res.LoadImage("subicon-boost"),
						Text:    "INTEGRITY",
						Trigger: func() {
							c.SetAction(&CombatActionBoost{stat: "INTEGRITY", isAttacker: true})
						},
					},
					{
						Icon:    res.LoadImage("icon-firewall"),
						SubIcon: res.LoadImage("subicon-boost"),
						Text:    "FIREWALL",
						Trigger: func() {
							c.SetAction(&CombatActionBoost{stat: "FIREWALL", isAttacker: true})
						},
					},
					{
						Icon:    res.LoadImage("icon-penetration"),
						SubIcon: res.LoadImage("subicon-boost"),
						Text:    "PENETRATION",
						Trigger: func() {
							c.SetAction(&CombatActionBoost{stat: "PENETRATION", isAttacker: true})
						},
					},
					{
						Text: "CANCEL",
						Trigger: func() {
							c.SwapMenu(CombatMenuModeMain)
						},
					},
				},
			},
		},
	}
	c.SwapMenu(CombatMenuModeMain)
	c.Refresh()

	return c
}

func (c *Combat) Refresh() {
	c.image.Fill(color.NRGBA{66, 20, 66, 200})
	pt := c.image.Bounds().Size()

	vector.StrokeRect(c.image, 0, 0, float32(pt.X), float32(pt.Y), 4, color.NRGBA{245, 120, 245, 255}, true)
}

func (c *Combat) GenerateEnemyAction() CombatAction {
	// FIXME: This is a placeholder. Enemy actions should be based on the enemy's stats as well as tendencies.
	targets := []string{"INTEGRITY", "FIREWALL", "PENETRATION"}
	t := targets[rand.Intn(len(targets))]
	if rand.Intn(2) == 0 {
		return &CombatActionAttack{stat: t, isAttacker: false}
	}
	return &CombatActionBoost{stat: t, isAttacker: false}
}

func (c *Combat) Update(w *World, r *Room) (cmd commands.Command) {
	c.attackerFloat += 0.025
	c.defenderFloat += 0.05
	if c.doneCommand != nil {
		return c.doneCommand
	}
	if c.action != nil {
		c.action.Update(c)
		if c.action.Done(c) {
			// If the action is not the attacker (which is always the player), that means the turn is over and we can swap back to main menu.
			if !c.action.IsAttacker() {
				c.SwapMenu(CombatMenuModeMain)
				c.action = nil
			} else {
				// Otherwise, it means the enemy should do an action (if not fleeing).
				if a, ok := c.action.(*CombatActionFlee); ok {
					if a.canFlee {
						c.doneCommand = commands.CombatResult{Fled: true}
					} else {
						c.SetAction(c.GenerateEnemyAction())
					}
				} else {
					c.SetAction(c.GenerateEnemyAction())
				}
			}
		}
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
	if c.action != nil {
		return
	}
	switch in := in.(type) {
	case inputs.Cancel:
		c.SwapMenu(CombatMenuModeMain)
	case inputs.Confirm:
		if c.menu.selectedIndex >= 0 && c.menu.selectedIndex < len(c.menu.items) {
			c.menu.items[c.menu.selectedIndex].Trigger()
		}
	case inputs.Direction:
		if in.Y < 0 {
			for {
				c.menu.selectedIndex--
				if c.menu.selectedIndex < 0 {
					c.menu.selectedIndex = 0
					break
				}
				if !c.menu.items[c.menu.selectedIndex].Disabled {
					break
				}
			}
		} else if in.Y > 0 {
			for {
				c.menu.selectedIndex++
				if c.menu.selectedIndex > len(c.menu.items)-1 {
					c.menu.selectedIndex = len(c.menu.items) - 1
					break
				}
				if !c.menu.items[c.menu.selectedIndex].Disabled {
					break
				}
			}
		}
	case inputs.Click:
		x := in.X - c.x
		y := in.Y - c.y
		for i, item := range c.menu.items {
			if x >= float64(item.Bounds.Min.X) && x <= float64(item.Bounds.Max.X) && y >= float64(item.Bounds.Min.Y) && y <= float64(item.Bounds.Max.Y) {
				if !item.Disabled {
					c.menu.selectedIndex = i
					item.Trigger()
				}
			}
		}
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

	// Draw combat report.
	{
		mx := int(cw) + 10
		my := 0
		mw := 190
		mh := 240
		vector.DrawFilledRect(c.image, float32(mx), float32(my), float32(mw), float32(mh), color.NRGBA{66, 66, 60, 220}, true)
		vector.StrokeRect(c.image, float32(mx), float32(my), float32(mw), float32(mh), 4, color.NRGBA{245, 245, 220, 255}, true)
		mx += 6
		my += 4
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		res.Text.SetAlign(etxt.Top | etxt.Left)
		y := my
		for _, line := range c.report {
			if line.color.A == 0 {
				res.Text.SetColor(color.NRGBA{200, 200, 200, 200})
			} else {
				res.Text.SetColor(line.color)
			}
			x := 0
			t := line.text
			if line.icon != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(mx), float64(y))
				c.image.DrawImage(line.icon, op)
				t = "  " + t
			}
			res.Text.DrawWithWrap(c.image, t, mx+x, y, mw-26)
			y += res.Text.MeasureWithWrap(t, mw-26).IntHeight()
		}
		res.Text.Utils().RestoreState()
	}

	// Draw right menu.
	mx := cw + 10
	my := ch - 150
	mw := 190
	mh := 150
	if c.action == nil {
		vector.DrawFilledRect(c.image, float32(mx), float32(my), float32(mw), float32(mh), color.NRGBA{66, 66, 60, 220}, true)
		vector.StrokeRect(c.image, float32(mx), float32(my), float32(mw), float32(mh), 4, color.NRGBA{245, 245, 220, 255}, true)
		mx += 6
		my += 4
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		res.Text.SetAlign(etxt.Top | etxt.Left)

		res.Text.SetColor(color.NRGBA{219, 86, 32, 200})
		for i, item := range c.menu.items {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(mx), float64(my))
			if item.Icon != nil {
				c.image.DrawImage(item.Icon, op)
				if item.SubIcon != nil {
					op.GeoM.Translate(16, 0)
					c.image.DrawImage(item.SubIcon, op)
				}
			}
			s := item.Text
			if i == c.menu.selectedIndex {
				s = "   > " + s
			} else {
				s = "     " + s
			}
			if item.Disabled {
				res.Text.SetColor(color.NRGBA{219, 86, 32, 100})
			} else {
				res.Text.SetColor(color.NRGBA{219, 86, 32, 255})
			}

			res.Text.Draw(c.image, s, int(mx), int(my))
			item.Bounds = image.Rect(int(mx), int(my), int(mx)+res.Text.Measure(s).IntWidth(), int(my)+res.Text.Measure(s).IntHeight())
			my += float64(res.DefFont.Size)
			c.menu.items[i] = item
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
