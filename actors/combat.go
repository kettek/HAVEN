package actors

import (
	"math/rand"

	"github.com/kettek/ebihack23/game"
)

type Combat struct {
	level              int
	exp                int
	penetration        int
	maxPenetration     int
	firewall           int
	maxFirewall        int
	integrity          int
	maxIntegrity       int
	penaltyPenetration int
	penaltyFirewall    int
	penaltyIntegrity   int
	killed             bool
	captured           bool
	glitches           []game.GlitchActor
	currentGlitch      game.GlitchActor
	ability            game.Ability
}

func roll(count int) (result int) {
	for i := 0; i < count; i++ {
		result += rand.Intn(2) + 1
	}
	return
}

func (c *Combat) ReduceDamage(pen, fire, inte int) (rpen, rfire, rinte int) {
	_, f, _ := c.CurrentStats()

	if f <= 1 {
		f = 2
	}
	f = roll(f)

	if pen > 0 {
		pen -= f
		rpen = pen
		if rpen < 0 {
			rpen = 1
		}
	}
	if fire > 0 {
		fire -= f
		rfire = fire
		if rfire < 0 {
			rfire = 1
		}
	}
	if inte > 0 {
		inte -= f
		rinte = inte
		if inte < 0 {
			rinte = 1
		}
	}

	return
}

func (c *Combat) ApplyDamage(pen, fire, inte int) (rpen, rfire, rinte int) {
	if pen >= 0 {
		if c.penetration-pen < 0 {
			pen = c.penetration
		}
		c.penetration -= pen
	}
	if fire >= 0 {
		if c.firewall-fire < 0 {
			fire = c.firewall
		}
		c.firewall -= fire
	}
	if inte >= 0 {
		if c.integrity-inte < 0 {
			inte = c.integrity
		}
		c.integrity -= inte
	}

	return pen, fire, inte
}

func (c *Combat) ApplyBoost(pen, fire, inte int) (int, int, int) {
	if pen >= 0 {
		if c.penetration > c.maxPenetration {
			pen /= 2
		}
		if pen == 0 {
			pen = 1
		}
		c.penetration += pen
	}
	if fire > 0 {
		if c.firewall > c.maxFirewall {
			fire /= 2
		}
		if fire == 0 {
			fire = 1
		}
		c.firewall += fire
	}
	if inte > 0 {
		if c.integrity > c.maxIntegrity {
			inte /= 2
		}
		if inte == 0 {
			inte = 1
		}
		c.integrity += inte
	}

	return pen, fire, inte
}

func (c *Combat) SetStats(pen, fire, inte int) {
	c.maxPenetration = pen
	c.maxFirewall = fire
	c.maxIntegrity = inte
	c.RestoreStats()
}

func (c *Combat) RestoreStats() {
	p, f, i := c.MaxStats()
	c.penetration = p
	c.firewall = f
	c.integrity = i
}

func (c *Combat) CurrentStats() (int, int, int) {
	p, f, i := c.penetration, c.firewall, c.integrity

	return p - c.penaltyPenetration, f - c.penaltyFirewall, i - c.penaltyIntegrity
}

func (c *Combat) MaxStats() (int, int, int) {
	p, f, i := c.maxPenetration, c.maxFirewall, c.maxIntegrity

	p += c.level * c.maxPenetration / 10
	f += c.level * c.maxFirewall / 10
	i += c.level * c.maxIntegrity / 10

	return p, f, i
}

func (c *Combat) Level() int {
	return c.level
}

func (c *Combat) SetLevel(l int) {
	c.level = l
}

func (c *Combat) Exp() int {
	return c.exp
}

func (c *Combat) AddExp(e int) int {
	c.exp += e
	lvl := 0
	for c.exp >= 100 {
		c.exp -= 100
		c.level++
		lvl++
		c.RestoreStats()
	}
	return lvl
}

func (c *Combat) ExpValue() int {
	p, f, i := c.MaxStats()
	return p + f + i + c.level*5
}

func (c *Combat) RollBoost() (int, int, int) {
	mp, mf, mi := c.MaxStats()
	mp2, mf2, mi2 := mp, mf, mi
	p, f, i := c.CurrentStats()
	if p > 0 {
		mp /= p
	}
	if f > 0 {
		mf /= f
	}
	if i > 0 {
		mi /= i
	}
	if mp <= 1 {
		mp = 2
	}
	if mf <= 1 {
		mf = 2
	}
	if mi <= 1 {
		mi = 2
	}

	mp = roll(mp)
	mf = roll(mf)
	mi = roll(mi)

	if mp > mp2/3 {
		mp = mp2 / 3
	}
	if mf > mf2/3 {
		mf = mf2 / 3
	}
	if mi > mi2/3 {
		mi = mi2 / 3
	}

	return mp, mf, mi
}

func (c *Combat) RollAttack() int {
	p, _, _ := c.CurrentStats()

	if p <= 0 {
		p = 1
	}

	p = roll(p)

	return p
}

func (c *Combat) HasGlitch() bool {
	return len(c.glitches) > 0
}

func (c *Combat) Glitches() []game.GlitchActor {
	return c.glitches
}

func (c *Combat) SetGlitch(g game.GlitchActor) {
	for _, g2 := range c.glitches {
		if g == g2 {
			c.currentGlitch = g
			return
		}
	}
}

func (c *Combat) AddGlitch(g game.GlitchActor) {
	c.glitches = append(c.glitches, g)
	if c.currentGlitch == nil {
		c.currentGlitch = g
	}
	if gl, ok := g.(*Glitch); ok {
		gl.spriteStack.Highlight = false
		gl.spriteStack.SkewX = 0
		gl.spriteStack.SkewY = 0
	}
}

func (c *Combat) RemoveGlitch(g game.GlitchActor) {
	for i, g2 := range c.glitches {
		if g == g2 {
			c.glitches = append(c.glitches[:i], c.glitches[i+1:]...)
			if c.currentGlitch == g {
				if len(c.glitches) > 0 {
					c.currentGlitch = c.glitches[0]
				} else {
					c.currentGlitch = nil
				}
			}
			break
		}
	}
}

func (c *Combat) CurrentGlitch() game.GlitchActor {
	return c.currentGlitch
}

func (c *Combat) Penalize(pen, fire, inte int) {
	c.penaltyPenetration += pen
	c.penaltyFirewall += fire
	c.penaltyIntegrity += inte
}

func (c *Combat) ClearPenalties() {
	c.penaltyPenetration = 0
	c.penaltyFirewall = 0
	c.penaltyIntegrity = 0
}

func (c *Combat) Killed() bool {
	return c.killed
}

func (c *Combat) Kill() {
	c.killed = true
}

func (c *Combat) Captured() bool {
	return c.captured
}

func (c *Combat) Capture() {
	c.captured = true
}

func (c *Combat) Ability() game.Ability {
	return c.ability
}

func (c *Combat) SetAbility(a game.Ability) {
	c.ability = a
}
