package game

type Ability interface {
	Use()
	Cooldown() int
	Active() bool
	Turn()
}

type Block struct {
	level       int
	cooldown    int
	turnsActive int
}

func (b *Block) Reset() {
	b.cooldown = 0
}

func (b *Block) Use() {
	b.cooldown = b.level
}

func (b *Block) Cooldown() int {
	return b.cooldown
}

func (b *Block) Active() bool {
	return b.turnsActive > 0
}

func (b *Block) Turn(o CombatActor) {
	if b.turnsActive > 0 {
		b.turnsActive--
	}
}
