package automaton_s

import "backend/internal/domain/model"

// 删除经CompleteDFA添加的陷阱状态以及相关转换，如果deaed为""，则默认采用model.SinkState，即"_sink_"
func UnCompleteDFA(a *model.Automaton, sink model.State) *model.Automaton {
	if sink == "" {
		sink = model.SinkState
	}
	// 判断是否有经过CompleteDFA添加的陷阱状态
	hasSinkState := false
	newStates := make([]model.State, 0, len(a.States)-1)
	for _, s := range a.States {
		if s == sink {
			hasSinkState = true
			break
		}
		newStates = append(newStates, s)
	}
	if !hasSinkState {
		return a
	}

	// 删除相关转换函数
	newTransitions := make([]model.Transition, 0, len(a.Transitions))
	for _, t := range a.Transitions {
		if t.ToStates[0] != sink {
			newTransitions = append(newTransitions, t)
		}
	}
	a.Transitions = newTransitions
	a.States = newStates

	return a
}
