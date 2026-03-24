package automaton_s

import (
	"backend/internal/domain/model"
)

func getDFANextStateFromMap(automaton *model.Automaton, from model.State, input model.Symbol) model.State {
	// 使用 TransMap 提高查找效率
	if stateMap, exists := automaton.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists && len(targets) == 1 {
			return targets[0] //
		}
	}

	return ""
}

// 获取NFA的下一个状态，不含空字符的NFA状态转移，非 ε-闭包
func getNFANextStatesFromMap(automaton *model.Automaton, from model.State, input model.Symbol) []model.State {
	// 使用 TransMap 提高查找效率
	if stateMap, exists := automaton.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists {
			return targets
		}
	}

	return nil
}

// getEpsilonNFANextStatesFromMap 从一组状态出发，对给定输入符号 sym，返回所有可达的下一状态（不含 ε-闭包）
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

// containsState 检查状态是否在状态列表中
func ContainsState(states []model.State, s model.State) bool {
	for _, st := range states {
		if st == s {
			return true
		}
	}

	return false
}

// computeEpsilonClosure 计算状态集合的 ε-闭包
func computeEpsilonClosure(nfa *model.Automaton, states []model.State) []model.State {
	closure := make([]model.State, 0, len(states))
	visited := make(map[model.State]bool)

	// 初始化：加入输入状态
	for _, s := range states {
		if !visited[s] {
			closure = append(closure, s)
			visited[s] = true
		}
	}

	// 使用队列进行 DFS 遍历 ε 转移
	queue := append([]model.State(nil), states...)

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		// 查找所有 ε 转移
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

// computeEpsilonClosureWithMap 使用 automaton.TransMap 加速计算状态集合的 ε-闭包
func computeEpsilonClosureWithMap(nfa *model.Automaton, states []model.State) []model.State {
	closure := make([]model.State, 0, len(states))
	visited := make(map[model.State]bool)

	// 初始化：加入输入状态
	for _, s := range states {
		if !visited[s] {
			closure = append(closure, s)
			visited[s] = true
		}
	}

	// 使用队列进行 DFS 遍历 ε 转移
	queue := append([]model.State(nil), states...)

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		// 查找所有 ε 转移
		for _, target := range nfa.TransMap[current][model.Epsilon] {
			if !visited[target] {
				visited[target] = true
				closure = append(closure, target)
				queue = append(queue, target)
			}
		}
	}

	return closure
}

// GetNonAcceptingStates 返回所有非接受状态
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
