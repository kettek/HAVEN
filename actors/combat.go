package actors

import (
	"fmt"
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

func (c *Combat) ApplyDamage(pen, fire, inte int) (rpen, rfire, rinte int) {
	_, f, _ := c.CurrentStats()

	fmt.Println("apply damage", pen, fire, inte)

	if pen > 0 {
		pen -= f
		rpen = pen
	}
	if pen < 0 {
		pen = 0
		rpen = 0
	}
	if fire > 0 {
		fire -= f
		rfire = fire
	}
	if fire < 0 {
		fire = 0
		rfire = 0
	}
	if inte > 0 {
		inte -= f
		rinte = inte
	}
	if inte < 0 {
		inte = 0
		rinte = 0
	}

	fmt.Println("apply damage2", pen, fire, inte)

	c.penetration -= pen
	c.firewall -= fire
	c.integrity -= inte

	if c.penetration < 0 {
		c.penetration = 0
	}
	if c.firewall < 0 {
		c.firewall = 0
	}
	if c.integrity < 0 {
		c.integrity = 0
	}
	return
}

func (c *Combat) ApplyBoost(pen, fire, inte int) (int, int, int) {
	if pen != 0 && c.penetration > c.maxPenetration {
		pen /= 2
	}
	if fire != 0 && c.firewall > c.maxFirewall {
		fire /= 2
	}
	if inte != 0 && c.integrity > c.maxIntegrity {
		inte /= 2
	}

	c.penetration += pen
	c.firewall += fire
	c.integrity += inte
	return pen, fire, inte
}

func (c *Combat) SetStats(pen, fire, inte int) {
	c.penetration = pen
	c.maxPenetration = pen
	c.firewall = fire
	c.maxFirewall = fire
	c.integrity = inte
	c.maxIntegrity = inte
}

func (c *Combat) CurrentStats() (int, int, int) {
	p, f, i := c.penetration, c.firewall, c.integrity

	p += c.level * c.maxPenetration / 10
	f += c.level * c.maxFirewall / 10
	i += c.level * c.maxIntegrity / 10

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
	p, f, i := c.MaxStats()
	p /= 5
	f /= 5
	i /= 5
	if p <= 0 {
		p = 2
	}
	if f <= 0 {
		f = 2
	}
	if i <= 0 {
		i = 2
	}
	return rand.Intn(p), rand.Intn(f), rand.Intn(i)
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
