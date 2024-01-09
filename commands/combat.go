package commands

type Combat struct {
	Attacker interface{}
	Defender interface{}
}

type CombatResult struct {
	Winner    interface{}
	Loser     interface{}
	ExpGained int
	Destroyed bool
	Fled      bool
}
