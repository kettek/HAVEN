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
	AutoCapture   bool
	ability       Ability
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
var abilityColor = color.NRGBA{50, 50, 255, 200}

type CombatAction interface {
	Done(c *Combat) (CombatAction, bool)
	Update(c *Combat)
	IsAttacker() bool
}

type CombatActionDone struct {
	Result     commands.CombatResult
	isAttacker bool
	timer      int
}

func (c CombatActionDone) Done(combat *Combat) (CombatAction, bool) {
	if c.timer >= 120 {
		combat.doneCommand = c.Result
		return nil, true
	}
	return nil, false
}

func (c *CombatActionDone) Update(combat *Combat) {
	c.timer++
	if c.timer == 10 {
		if c.Result.Fled {
			// ???
		} else if c.isAttacker {
			if !c.Result.Fled && !c.Result.Destroyed {
				// Yeah, yeah, capturing isn't handled here because I'm lazy.
			} else {
				combat.AddReport(fmt.Sprintf("%s destroys %s!", combat.Attacker.Name(), combat.Defender.Name()), nil, attackColor)
				res.PlaySound("death")
			}
		} else {
			combat.AddReport(fmt.Sprintf("%s infects %s!", combat.Defender.Name(), combat.Attacker.Name()), nil, attackColor)
			res.PlaySound("death")
		}
	}
}

func (c CombatActionDone) IsAttacker() bool {
	return c.isAttacker
}

type CombatActionAttack struct {
	stat       string
	isAttacker bool
	timer      int
	next       CombatAction
}

func (c CombatActionAttack) Done(cmb *Combat) (CombatAction, bool) {
	if c.timer >= 120 {
		if c.next != nil {
			return c.next, true
		}
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
			res.PlaySound("candie")
		}
		return nil, true
	}
	return nil, false
}

func (c *CombatActionAttack) Update(cmb *Combat) {
	c.timer++

	// This isn't efficient, but it's easier if I throw this here.
	var attackerAbility *Ability
	var defenderAbility *Ability
	if c.isAttacker {
		if glitch := cmb.Attacker.CurrentGlitch(); glitch != nil {
			attackerAbility = glitch.Ability()
		}
	} else {
		if glitch := cmb.Attacker.CurrentGlitch(); glitch != nil {
			defenderAbility = glitch.Ability()
		}
	}
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
		var bonus int
		if attackerAbility != nil && attackerAbility.Name == AbilityPerfectHit {
			v = attackerAbility.Tier * 2
			bonus = 0
			cmb.AddReport(fmt.Sprintf("%s for %d!", attackerAbility.Name, v), nil, abilityColor)
		}
		if attackerAbility != nil && attackerAbility.Name == AbilityRandomDamage {
			b := rand.Intn(1 + attackerAbility.Tier)
			bonus += b
			cmb.AddReport(fmt.Sprintf("%s for +%d!", attackerAbility.Name, b), nil, abilityColor)
		}
		if v+bonus <= 0 {
			cmb.AddReport(fmt.Sprintf("%s attacks %s, but misses!", attacker.Name(), c.stat), icon, infoColor)
			res.PlaySound("miss")
			return
		}
		_, _, inte := defender.CurrentStats()
		if inte <= 0 {
			if defenderAbility != nil && defenderAbility.Name == AbilityHardy {
				cmb.AddReport(fmt.Sprintf("%s!", defenderAbility.Name), nil, infoColor)
				res.PlaySound("miss")
				return
			} else {
				defender.Kill()
				c.next = &CombatActionDone{
					isAttacker: c.isAttacker,
					Result: commands.CombatResult{
						Winner:    attacker,
						Loser:     defender,
						Destroyed: true,
						ExpGained: defender.ExpValue(),
					},
				}
				c.timer = 120
				return
			}
		}
		if defenderAbility != nil {
			if defenderAbility.Name == AbilityBlock {
				red := defenderAbility.Tier * 2
				// reduce v by red, and if red still has some left, apply it to bonus.
				v -= red
				if v < 0 {
					bonus += v
					if bonus < 0 {
						bonus = 0
					}
					v = 0
				}
				cmb.AddReport(fmt.Sprintf("%s!", defenderAbility.Name), nil, abilityColor)
			} else if defenderAbility.Name == AbilityPerfectBlock {
				cmb.AddReport(fmt.Sprintf("%s!", defenderAbility.Name), nil, abilityColor)
				v = 0
				bonus = 0
			}
		}
		if c.stat == "INTEGRITY" {
			_, _, v = defender.ReduceDamage(-1, -1, v+bonus)
			_, _, v = defender.ApplyDamage(-1, -1, v+bonus)
		} else if c.stat == "FIREWALL" {
			_, v, _ = defender.ReduceDamage(-1, v+bonus, -1)
			_, v, _ = defender.ApplyDamage(-1, v+bonus, -1)
		} else if c.stat == "PENETRATION" {
			v, _, _ = defender.ReduceDamage(v+bonus, -1, -1)
			v, _, _ = defender.ApplyDamage(v+bonus, -1, -1)
		}
		if v <= 0 {
			cmb.AddReport(fmt.Sprintf("%s attacks %s, but is denied!", attacker.Name(), c.stat), icon, infoColor)
			res.PlaySound("miss")
			return
		}
		if bonus > 0 {
			cmb.AddReport(fmt.Sprintf("%s attacks %s for %d(%d+%d)!", attacker.Name(), c.stat, v, v-bonus, bonus), icon, attackColor)
		} else {
			cmb.AddReport(fmt.Sprintf("%s attacks %s for %d!", attacker.Name(), c.stat, v), icon, attackColor)
		}
		res.PlaySound("hit")
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

func (c CombatActionBoost) Done(cmb *Combat) (CombatAction, bool) {
	return nil, c.timer >= 120
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
			_, _, v = attacker.ApplyBoost(-1, -1, i)
		} else if c.stat == "FIREWALL" {
			_, v, _ = attacker.ApplyBoost(-1, f, -1)
		} else if c.stat == "PENETRATION" {
			v, _, _ = attacker.ApplyBoost(p, -1, -1)
		}
		if v <= 0 {
			cmb.AddReport(fmt.Sprintf("%s fails to boost %s", attacker.Name(), c.stat), icon, infoColor)
			res.PlaySound("miss")
		} else {
			cmb.AddReport(fmt.Sprintf("%s boosts %s for %d!", attacker.Name(), c.stat, v), icon, defenseColor)
			res.PlaySound("boost")
		}
	}
}

func (c CombatActionBoost) IsAttacker() bool {
	return c.isAttacker
}

type CombatActionAbility struct {
	isAttacker bool
	timer      int
}

func (c CombatActionAbility) Done(cmb *Combat) (CombatAction, bool) {
	return nil, c.timer >= 120
}

func (c *CombatActionAbility) Update(cmb *Combat) {
	c.timer++
	//var attacker CombatActor
	var defender CombatActor
	var attackerAbility *Ability
	var defenderAbility *Ability
	if c.isAttacker {
		if glitch := cmb.Attacker.CurrentGlitch(); glitch != nil {
			attackerAbility = glitch.Ability()
		}
		if glitch := cmb.Defender.CurrentGlitch(); glitch != nil {
			defenderAbility = glitch.Ability()
		}
		//attacker = cmb.Attacker
		defender = cmb.Defender
	} else {
		if glitch := cmb.Attacker.CurrentGlitch(); glitch != nil {
			defenderAbility = glitch.Ability()
		}
		if glitch := cmb.Defender.CurrentGlitch(); glitch != nil {
			attackerAbility = glitch.Ability()
		}
		//attacker = cmb.Defender
		defender = cmb.Attacker
	}
	if c.timer == 10 {
		if attackerAbility != nil {
			if attackerAbility.Name == AbilityCleave {
				p, f, i := defender.CurrentStats()
				r := rand.Intn(3)
				if r == 0 {
					p /= 2
					f = -1
					i = -i
				} else if r == 1 {
					p = -1
					f /= 2
					i = -1
				} else {
					p = -1
					f = -1
					i /= 2
				}
				if defenderAbility != nil {
					if defenderAbility.Name == AbilityBlock {
						if r == 0 {
							p -= defenderAbility.Tier * 2
							if p < 0 {
								p = 0
							}
						} else if r == 1 {
							f -= defenderAbility.Tier * 2
							if f < 0 {
								f = 0
							}
						} else {
							i -= defenderAbility.Tier * 2
							if i < 0 {
								i = 0
							}
						}
					} else if defenderAbility.Name == AbilityPerfectBlock {
						p = -1
						f = -1
						i = -1
					}
				}
				p, f, i = defender.ReduceDamage(p, f, i)
				defender.ApplyDamage(p, f, i)
			}
		}
	}
}

func (c CombatActionAbility) IsAttacker() bool {
	return c.isAttacker
}

type CombatActionFlee struct {
	isAttacker bool
	canFlee    bool
	timer      int
}

func (c CombatActionFlee) Done(cmb *Combat) (CombatAction, bool) {
	return nil, c.timer >= 120
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

type CombatActionCapture struct {
	isAttacker bool
	timer      int
	caught     bool
}

func (c CombatActionCapture) IsAttacker() bool {
	return c.isAttacker
}

func (c CombatActionCapture) Done(cmb *Combat) (CombatAction, bool) {
	if cmb.AutoCapture {
		return &CombatActionDone{
			isAttacker: c.isAttacker,
			Result: commands.CombatResult{
				Winner:    cmb.Attacker,
				Loser:     cmb.Defender,
				ExpGained: cmb.Defender.ExpValue(),
				Destroyed: false,
				Fled:      false,
			},
		}, true
	}

	if c.timer < 360 {
		return nil, false

	}
	if c.caught {
		return &CombatActionDone{
			isAttacker: c.isAttacker,
			Result: commands.CombatResult{
				Winner:    cmb.Attacker,
				Loser:     cmb.Defender,
				ExpGained: cmb.Defender.ExpValue(),
				Destroyed: false,
				Fled:      false,
			},
		}, true
	}
	return nil, true
}

func (c *CombatActionCapture) Try(cmb *Combat) bool {
	return rand.Float64() < cmb.CaptureChance()
}

func (cmb *Combat) CaptureChance() float64 {
	ap, _, _ := cmb.Attacker.CurrentStats()
	apm, _, _ := cmb.Attacker.MaxStats()
	dp, df, di := cmb.Defender.CurrentStats()
	dpm, dfm, dim := cmb.Defender.MaxStats()

	defi := float64(di) / float64(dim)
	defp := float64(dp) / float64(dpm)
	defm := float64(df) / float64(dfm)
	defv := (defp + defm + defi) / 3

	atkv := float64(ap) / float64(apm)

	resv := math.Min(100, math.Max(0, atkv-defv))

	return resv
}

func (c *CombatActionCapture) Update(cmb *Combat) {
	c.timer++
	if c.timer == 1 {
		if len(cmb.Attacker.Glitches()) >= 9 {
			cmb.AddReport(fmt.Sprintf("%s attempts to capture %s, but the quarantine is full!", cmb.Attacker.Name(), cmb.Defender.Name()), nil, neutralColor)
			c.timer = 301
			res.PlaySound("miss")
		} else {
			cmb.AddReport(fmt.Sprintf("%s attempts to capture %s!", cmb.Attacker.Name(), cmb.Defender.Name()), nil, neutralColor)
		}
	} else if c.timer == 60 {
		cmb.AddReport("maybe...", nil, neutralColor)
	} else if c.timer == 120 {
		if c.Try(cmb) {
			cmb.AddReport(fmt.Sprintf("%s captures %s!", cmb.Attacker.Name(), cmb.Defender.Name()), nil, neutralColor)
			c.timer = 301
			c.caught = true
			res.PlaySound("caught")
			return
		}
	} else if c.timer == 180 {
		cmb.AddReport("maybe...!", nil, neutralColor)
	} else if c.timer == 240 {
		if c.Try(cmb) {
			cmb.AddReport(fmt.Sprintf("%s captures %s!", cmb.Attacker.Name(), cmb.Defender.Name()), nil, neutralColor)
			res.PlaySound("caught")
			c.timer = 301
			c.caught = true
			return
		}
	} else if c.timer == 300 {
		if !c.Try(cmb) {
			cmb.AddReport(fmt.Sprintf("%s failed to capture %s!", cmb.Attacker.Name(), cmb.Defender.Name()), nil, neutralColor)
			res.PlaySound("miss")
		} else {
			cmb.AddReport(fmt.Sprintf("%s captures %s!", cmb.Attacker.Name(), cmb.Defender.Name()), nil, neutralColor)
			res.PlaySound("caught")
			c.caught = true
		}
	}
}

type CombatActionSwapGlitch struct {
	isAttacker bool
	timer      int
	glitch     GlitchActor
}

func (c CombatActionSwapGlitch) Done(cmb *Combat) (CombatAction, bool) {
	return nil, c.timer >= 120
}

func (c *CombatActionSwapGlitch) Update(cmb *Combat) {
	c.timer++
	if c.timer == 10 {
		cmb.AddReport(fmt.Sprintf("%s swaps to %s!", cmb.Attacker.Name(), c.glitch.Name()), nil, neutralColor)
	} else if c.timer == 60 {
		cmb.Attacker.SetGlitch(c.glitch)
	}
}

func (c CombatActionSwapGlitch) IsAttacker() bool {
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
	Glitch   GlitchActor
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
func (c *Combat) RefreshGlitchUse() {
	var glitchMenuItems []CombatMenuItem
	if gl := c.Attacker.CurrentGlitch(); gl != nil {
		if abil := gl.Ability(); abil != nil {
			text := abil.Name
			if abil.OnCooldown() || abil.IsActive() {
				text += fmt.Sprintf(" (%d)", abil.cooldown+(abil.Turns-abil.turnsActive))
			}
			glitchMenuItems = append(glitchMenuItems, CombatMenuItem{
				Text:     text,
				Disabled: abil.OnCooldown() || abil.IsActive(),
				Trigger: func() {
					abil.Activate()
					c.SetAction(&CombatActionAbility{
						isAttacker: true,
					})
				},
			})
		}
	}
	glitchMenuItems = append(glitchMenuItems, CombatMenuItem{
		Text: "CANCEL",
		Trigger: func() {
			c.SwapMenu(CombatMenuModeMain)
		},
	})
	c.menus.use.items = glitchMenuItems
	c.menus.use.selectedIndex = 0
}

func (c *Combat) RefreshGlitchSwap() {
	var glitchMenuItems []CombatMenuItem
	for _, g := range c.Attacker.Glitches() {
		func(g GlitchActor) {
			glitchMenuItems = append(glitchMenuItems, CombatMenuItem{
				Glitch:   g,
				Disabled: c.Attacker.CurrentGlitch() == g,
				Text:     fmt.Sprintf("%s (%d)", g.Name(), g.Level()),
				Trigger: func() {
					if c.Attacker.CurrentGlitch() == g {
						return
					}
					c.SetAction(&CombatActionSwapGlitch{
						isAttacker: true,
						glitch:     g,
					})
				},
			})
		}(g)
	}
	glitchMenuItems = append(glitchMenuItems, CombatMenuItem{
		Text: "CANCEL",
		Trigger: func() {
			c.SwapMenu(CombatMenuModeMain)
		},
	})

	c.menus.swap.items = glitchMenuItems
	c.menus.swap.selectedIndex = 0
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
							c.SetAction(&CombatActionCapture{isAttacker: true})
						},
						//Disabled: !attacker.HasGlitch(),
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
	c.RefreshGlitchUse()
	c.RefreshGlitchSwap()
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

func (c *Combat) RefreshAbilities() {
	actors := []CombatActor{c.Attacker, c.Defender}
	for _, actor := range actors {
		for _, glitch := range actor.Glitches() {
			if abil := glitch.Ability(); abil != nil {
				if abil.IsActive() {
					abil.Turn()
				} else {
					abil.ReduceCooldown()
				}
			}
		}
	}
}

func (c *Combat) Update(w *World, r *Room) (cmd commands.Command) {
	c.attackerFloat += 0.025
	c.defenderFloat += 0.05
	if c.doneCommand != nil {
		return c.doneCommand
	}
	if c.action != nil {
		c.action.Update(c)
		next, done := c.action.Done(c)
		if done {
			c.RefreshAbilities()
			if next != nil {
				c.SetAction(next)
			} else {
				if c.action.IsAttacker() {
					// Otherwise, it means the enemy should do an action (if not fleeing).
					if a, ok := c.action.(*CombatActionFlee); ok {
						if a.canFlee {
							c.doneCommand = commands.CombatResult{Fled: true, Winner: c.Attacker, Loser: c.Defender}
						} else {
							c.SetAction(c.GenerateEnemyAction())
						}
					} else {
						c.SetAction(c.GenerateEnemyAction())
					}
				} else {
					// If the action is not the attacker (which is always the player), that means the turn is over and we can swap back to main menu.
					c.SwapMenu(CombatMenuModeMain)
					c.action = nil
					// ugh, this is stupid, but having some dumb issues
					c.RefreshGlitchSwap()
					c.RefreshGlitchUse()
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
			if item.Glitch != nil {
				op.GeoM.Translate(float64(mx)+6, float64(my))
				item.Glitch.SpriteStack().DrawFlat(c.image, op.GeoM)
			} else if item.Icon != nil {
				op.GeoM.Translate(float64(mx), float64(my))
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

		mp, mf, mi := c.Attacker.MaxStats()
		cp, cf, ci := c.Attacker.CurrentStats()
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		x := 100
		y := c.image.Bounds().Dy() - 72
		res.Text.SetColor(color.NRGBA{50, 255, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("INTEGRITY   %d/%d", ci, mi), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(color.NRGBA{255, 50, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("FIREWALL    %d/%d", cf, mf), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(color.NRGBA{255, 255, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("PENETRATION %d/%d", cp, mp), x, y)
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

		mp, mf, mi := c.Defender.MaxStats()
		cp, cf, ci := c.Defender.CurrentStats()
		res.Text.Utils().StoreState()
		res.Text.SetSize(float64(res.DefFont.Size))
		res.Text.SetFont(res.DefFont.Font)
		x := 16
		y := 32
		res.Text.Draw(c.image, fmt.Sprintf("LVL %d %s", c.Defender.Level(), defender.Name()), x, y)
		y += res.DefFont.Size * 2
		res.Text.SetColor(color.NRGBA{50, 255, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("INTEGRITY   %d/%d", ci, mi), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(color.NRGBA{255, 50, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("FIREWALL    %d/%d", cf, mf), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(color.NRGBA{255, 255, 50, 200})
		res.Text.Draw(c.image, fmt.Sprintf("PENETRATION %d/%d", cp, mp), x, y)
		y += res.DefFont.Size
		res.Text.SetColor(neutralColor)
		res.Text.Draw(c.image, fmt.Sprintf("CAPTURE %d%%", int(c.CaptureChance()*100)), x, y)
		res.Text.Utils().RestoreState()
	}

	screen.DrawImage(c.image, op)
}
