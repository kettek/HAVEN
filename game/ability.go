package game

type AbilityType string

const (
	AbilityNone         AbilityType = ""
	AbilityBlock                    = "BLOCK"
	AbilityPerfectHit               = "PERFECT HIT"
	AbilityPerfectBlock             = "PERFECT BLOCK"
	AbilityHardy                    = "HARDY"
	AbilityCleave                   = "CLEAVE"
	AbilityRandomDamage             = "RANDOM DAMAGE"
)

type AbilityDescription string

const (
	AbilityDescriptionNone         AbilityDescription = ""
	AbilityDescriptionBlock                           = "Block the next 2*Ability damage up to Turns."
	AbilityDescriptionPerfectHit                      = "Deal 2*Ability Tier damage on the next attacks up to Turns."
	AbilityDescriptionPerfectBlock                    = "Block the next 2*Ability Tier attacks up to Turns."
	AbilityDescriptionHardy                           = "Avoid becoming infected up to Turns."
	AbilityDescriptionCleave                          = "Halves enemy INTEGRITY, FIREWALL, or PENETRATION."
	AbilityDescriptionRandomDamage                    = "Deals random bonus damage up to Tier for Turns."
)

type Ability struct {
	Name  string
	Tier  int
	Turns int

	cooldown    int
	turnsActive int
}

func (b *Ability) Reset() {
	b.cooldown = 0
	b.turnsActive = 0
}

func (b *Ability) ReduceCooldown() {
	if b.cooldown > 0 {
		b.cooldown--
	}
}

func (b *Ability) Activate() {
	b.turnsActive = b.Turns
}

func (b *Ability) OnCooldown() bool {
	return b.cooldown > 0
}

func (b *Ability) IsActive() bool {
	return b.turnsActive > 0
}

func (b *Ability) Turn() {
	if b.turnsActive > 0 {
		b.turnsActive--
		if b.turnsActive == 0 {
			b.Reset()
		}
	}
}
