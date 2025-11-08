package automaton_s

import "backend/internal/domain/model"

func getDFANextStateFromMap(automaton *model.Automaton, from model.State, input model.Symbol) model.State {
	// 使用 TransMap 提高查找效率
	if stateMap, exists := automaton.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists && len(targets) == 1 {
			return targets[0] //
		}
	}
	return ""
}

func getNFANextStatesFromMap(automaton *model.Automaton, from model.State, input model.Symbol) []model.State {
	// 使用 TransMap 提高查找效率
	if stateMap, exists := automaton.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists {
			return targets
		}
	}
	return nil
}

// containsState 检查状态是否在状态列表中
func containsState(states []model.State, s model.State) bool {
	for _, st := range states {
		if st == s {
			return true
		}
	}
	return false
}

// epsilonClosure 计算状态集合的 ε-闭包
func epsilonClosure(nfa *model.Automaton, states []model.State) []model.State {
	closure := make([]model.State, 0, len(states))
	visited := make(map[model.State]bool)

	// 初始化：加入输入状态
	for _, s := range states {
		if !visited[s] {
			closure = append(closure, s)
			visited[s] = true
		}
	}

	// 使用栈进行 DFS 遍历 ε 转移
	stack := append([]model.State(nil), states...)

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// 查找所有 ε 转移
		for _, t := range nfa.Transitions {
			if t.FromState == current && t.Input == model.Epsilon {
				for _, target := range t.ToStates {
					if !visited[target] {
						visited[target] = true
						closure = append(closure, target)
						stack = append(stack, target)
					}
				}
			}
		}
	}

	return closure
}
