package states

var CurrentState State

func NextState(state State) {
	if CurrentState != nil {
		CurrentState.Leave()
	}
	state.Enter()
	CurrentState = state
}
