package actors

import (
	"math"
	"math/rand"
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
	Glitches           []*Glitch
}

func (c *Combat) ReduceDamage(pen, fire, inte int) (rpen, rfire, rinte int) {
	_, f, _ := c.CurrentStats()

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

func (c *Combat) AddExp(e int) {
	c.exp += e
	if c.exp >= c.level*100 {
		c.exp -= c.level * 100
		c.level++
	}
}

func (c *Combat) ExpValue() int {
	p, f, i := c.MaxStats()
	return p + f + i + c.level*5
}

func (c *Combat) RollBoost() (int, int, int) {
	mp, mf, mi := c.MaxStats()
	p, f, i := c.CurrentStats()
	mp /= p
	mf /= f
	mi /= i
	if mp <= 1 {
		mp = 2
	}
	if mf <= 1 {
		mf = 2
	}
	if mi <= 1 {
		mi = 2
	}

	mp = rand.Intn(mp)
	mf = rand.Intn(mf)
	mi = rand.Intn(mi)

	return mp, mf, mi
}

func (c *Combat) RollAttack() int {
	p, _, _ := c.CurrentStats()

	if p <= 0 {
		p = 1
	}

	p = int(math.Ceil(math.Max(float64(p)/4, float64(rand.Intn(p)))))

	return p
}

func (c *Combat) HasGlitch() bool {
	return len(c.Glitches) > 0
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
