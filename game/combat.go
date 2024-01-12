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
}

type CombatAction interface {
	Done() bool
	Update(c *Combat)
	IsAttacker() bool
}

type CombatActionAttack struct {
	stat       string
	isAttacker bool
	timer      int
}

func (c CombatActionAttack) Done() bool {
	return c.timer >= 60
}

func (c *CombatActionAttack) Update(cmb *Combat) {
	c.timer++
}

func (c CombatActionAttack) IsAttacker() bool {
	return c.isAttacker
}

type CombatActionBoost struct {
	stat       string
	isAttacker bool
	timer      int
}

func (c CombatActionBoost) Done() bool {
	return c.timer >= 60
}

func (c *CombatActionBoost) Update(cmd *Combat) {
	c.timer++
}

func (c CombatActionBoost) IsAttacker() bool {
	return c.isAttacker
}

type CombatActionFlee struct {
	isAttacker bool
	canFlee    bool
	timer      int
}

func (c CombatActionFlee) Done() bool {
	return c.timer >= 30
}

func (c *CombatActionFlee) Update(cmb *Combat) {
	c.timer++
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
	Icon    *ebiten.Image
	SubIcon *ebiten.Image
	Text    string
	Bounds  image.Rectangle
	Trigger func()
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
						Text: "INCREASE STAT",
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
					},
					{
						Icon: res.LoadImage("icon-swapGlitch"),
						Text: "SWAP GLITCH",
						Trigger: func() {
							c.SwapMenu(CombatMenuModeSwapGlitch)
						},
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
		if c.action.Done() {
			// If the action is not the attacker (which is always the player), that means the turn is over and we can swap back to main menu.
			if !c.action.IsAttacker() {
				c.SwapMenu(CombatMenuModeMain)
				c.action = nil
			} else {
				// Otherwise, it means the enemy should do an action (if not fleeing).
				if a, ok := c.action.(*CombatActionFlee); ok {
					if a.canFlee {
						fmt.Println("fled successfully")
						c.doneCommand = commands.CombatResult{Fled: true}
					} else {
						fmt.Println("escape is denied!")
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
			c.menu.selectedIndex--
			if c.menu.selectedIndex < 0 {
				c.menu.selectedIndex = 0
			}
		} else if in.Y > 0 {
			c.menu.selectedIndex++
			if c.menu.selectedIndex > len(c.menu.items)-1 {
				c.menu.selectedIndex = len(c.menu.items) - 1
			}
		}
	case inputs.Click:
		x := in.X - c.x
		y := in.Y - c.y
		for i, item := range c.menu.items {
			if x >= float64(item.Bounds.Min.X) && x <= float64(item.Bounds.Max.X) && y >= float64(item.Bounds.Min.Y) && y <= float64(item.Bounds.Max.Y) {
				c.menu.selectedIndex = i
				item.Trigger()
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
