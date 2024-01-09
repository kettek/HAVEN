package actors

type Combat struct {
	level          int
	exp            int
	penetration    int
	maxPenetration int
	firewall       int
	maxFirewall    int
	integrity      int
	maxIntegrity   int
}

func (c *Combat) ApplyDamage(pen, fire, inte int) {
	c.penetration -= pen
	c.firewall -= fire
	c.integrity -= inte
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

	return p, f, i
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
