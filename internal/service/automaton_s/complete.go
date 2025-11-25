package automaton_s

import "backend/internal/domain/model"

// 补充死状态以及相关转换函数
func CompleteDFA(automaton *model.Automaton, dead model.State) *model.Automaton {
	if dead == "" {
		dead = model.DeadState
	}
	// 添加 dead state
	newStates := append(automaton.States, dead)
	newAccepting := automaton.AcceptingStates // dead 不是接受状态

	// 构建新转移
	newTransitions := make([]model.Transition, 0, len(automaton.Transitions))
	newTransitions = append(newTransitions, automaton.Transitions...)

	// 补全缺失转移
	for _, s := range automaton.States {
		for _, a := range automaton.Alphabet {
			found := false
			for _, t := range automaton.Transitions {
				if t.FromState == s && t.Input == a {
					found = true
					break
				}
			}
			if !found {
				newTransitions = append(newTransitions, model.Transition{
					FromState: s,
					Input:     a,
					ToStates:  []model.State{dead},
				})
			}
		}
	}

	// dead state 的转移：所有输入都指向自己
	for _, a := range automaton.Alphabet {
		newTransitions = append(newTransitions, model.Transition{
			FromState: dead,
			Input:     a,
			ToStates:  []model.State{dead},
		})
	}

	return &model.Automaton{
		States:          newStates,
		Alphabet:        automaton.Alphabet,
		Transitions:     newTransitions,
		InitialState:    automaton.InitialState,
		AcceptingStates: newAccepting,
		Type:            model.DFA,
	}
}

// 删除经CompleteDFA添加的死状态以及相关转换，如果deaed为“”，则默认采用model.DeadState，即"dead_state"
func UnCompleteDFA(a *model.Automaton, dead model.State) *model.Automaton {
	if dead == "" {
		dead = model.DeadState
	}
	// 判断是否有经过CompleteDFA添加的死状态
	var hasDeadState bool = false
	newStates := make([]model.State, 0, len(a.States)-1)
	for _, s := range a.States {
		if s == dead {
			hasDeadState = true
			break
		}
		newStates = append(newStates, s)
	}
	if !hasDeadState {
		return a
	}

	// 删除相关转换函数
	newTransitions := make([]model.Transition, 0, len(a.Transitions))
	for _, t := range a.Transitions {
		if t.ToStates[0] != dead {
			newTransitions = append(newTransitions, t)
		}
	}
	a.Transitions = newTransitions[:len(newTransitions)]
	a.States = newStates
	return a
}
