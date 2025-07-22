package task

func Contains(states []State, s State) bool {
	for _, v := range states {
		if v == s {
			return true
		}
	}
	return false
}

var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

func ValidStateTransition(src State, dest State) bool {
	return Contains(stateTransitionMap[src], dest)

}
