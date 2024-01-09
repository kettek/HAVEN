package actors

type Combat struct {
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

func (c *Combat) CurrentStats() (int, int, int) {
	return c.penetration, c.firewall, c.integrity
}

func (c *Combat) MaxStats() (int, int, int) {
	return c.maxPenetration, c.maxFirewall, c.maxIntegrity
}
