package automaton_s

import "github.com/chenzanhong/formallanglab-master/internal/domain/model"

// 删除经 CompleteDFA 添加的陷阱状态以及相关转换，如果 sink 为 ""，则默认采用 model.SinkState，即"_sink_"
func UnCompleteDFA(a *model.Automaton, sink model.State) *model.Automaton {
	if a == nil {
		return nil
	}
	if sink == "" {
		sink = model.SinkState
	}

	hasSinkState := false
	newStates := make([]model.State, 0, len(a.States))
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
