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
	Pending:   []State{Scheduled},
	Scheduled: []State{Scheduled, Running, Failed},
	Running:   []State{Running, Completed, Failed},
	Completed: []State{},
	Failed:    []State{},
}

func ValidStateTransition(src State, dest State) bool {
	return Contains(stateTransitionMap[src], dest)

}
