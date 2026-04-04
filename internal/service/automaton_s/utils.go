package automaton_s

import (
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// 使用 TransMap 查找下一个状态
func getDFANextStateFromMap(automaton *model.Automaton, from model.State, input model.Symbol) model.State {
	if stateMap, exists := automaton.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists && len(targets) == 1 {
			return targets[0]
		}
	}

	return ""
}

// 使用 TransMap 查找下一个可能的状态，非 ε-闭包
func getNFANextStatesFromMap(automaton *model.Automaton, from model.State, input model.Symbol) []model.State {
	if stateMap, exists := automaton.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists {
			return targets
		}
	}

	return nil
}

func getEpsilonNFANextStatesFromMap(automaton *model.Automaton, fromStates []model.State, sym model.Symbol) []model.State {
	var nextStates []model.State
	seen := make(map[model.State]bool)

	for _, s := range fromStates {
		if targets, exists := automaton.TransMap[s][sym]; exists {
			for _, t := range targets {
				if !seen[t] {
					seen[t] = true
					nextStates = append(nextStates, t)
				}
			}
		}
	}

	return nextStates
}

func ContainsState(states []model.State, s model.State) bool {
	for _, st := range states {
		if st == s {
			return true
		}
	}

	return false
}

// 计算状态集合的 ε-闭包
func computeEpsilonClosure(nfa *model.Automaton, states []model.State) []model.State {
	closure := make([]model.State, 0, len(states))
	visited := make(map[model.State]bool)

	for _, s := range states {
		if !visited[s] {
			closure = append(closure, s)
			visited[s] = true
		}
	}

	queue := append([]model.State(nil), states...)

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		for _, t := range nfa.Transitions {
			if t.FromState == current && t.Input == model.Epsilon {
				for _, target := range t.ToStates {
					if !visited[target] {
						visited[target] = true
						closure = append(closure, target)
						queue = append(queue, target)
					}
				}
			}
		}
	}

	return closure
}

// 使用 automaton.TransMap 加速计算状态集合的 ε-闭包
func computeEpsilonClosureWithMap(nfa *model.Automaton, states []model.State) []model.State {
	closure := make([]model.State, 0, len(states))
	visited := make(map[model.State]bool)

	for _, s := range states {
		if !visited[s] {
			closure = append(closure, s)
			visited[s] = true
		}
	}

	queue := append([]model.State(nil), states...)

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if targets, exists := nfa.TransMap[current][model.Epsilon]; exists {
			for _, target := range targets {
				if !visited[target] {
					visited[target] = true
					closure = append(closure, target)
					queue = append(queue, target)
				}
			}
		}
	}

	return closure
}

// 返回所有非接受状态
func GetNonAcceptingStates(a *model.Automaton) []model.State {
	acceptSet := make(map[model.State]bool)
	for _, s := range a.AcceptingStates {
		acceptSet[s] = true
	}
	var nonAccept []model.State
	for _, s := range a.States {
		if !acceptSet[s] {
			nonAccept = append(nonAccept, s)
		}
	}

	return nonAccept
}
