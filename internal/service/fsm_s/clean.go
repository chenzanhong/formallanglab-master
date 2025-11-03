package fsm_s

import (
	"backend/internal/domain/model" // 请根据你的项目路径调整
)

// Cleanup 清理自动机：移除不可达状态和未使用的符号
func Cleanup(fsm *model.Automaton) {
	if fsm == nil {
		return
	}

	// 用于存储有效（可达）状态
	reachable := make(map[model.State]bool)
	// 用于存储被使用的符号
	usedSymbols := make(map[model.Symbol]bool)

	// ========== 1. 使用 BFS 找出所有从初始状态可达的状态 ==========
	queue := []model.State{fsm.InitialState}
	reachable[fsm.InitialState] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 遍历所有从 current 出发的转移
		for _, t := range fsm.Transitions {
			if t.FromState != current {
				continue
			}

			// 标记使用的符号
			usedSymbols[t.Input] = true

			// 遍历所有目标状态
			for _, to := range t.ToStates {
				if !reachable[to] {
					reachable[to] = true
					queue = append(queue, to)
				}
			}
		}
	}

	// ========== 2. 构建新的状态列表（只保留可达状态） ==========
	var newStates []model.State
	stateSet := make(map[model.State]bool) // 用于快速查找
	for _, s := range fsm.States {
		if reachable[s] {
			newStates = append(newStates, s)
			stateSet[s] = true
		}
	}
	fsm.States = newStates

	// ========== 3. 构建新的字母表（只保留被使用的符号） ==========
	var newAlphabet []model.Symbol
	for _, sym := range fsm.Alphabet {
		if usedSymbols[sym] {
			newAlphabet = append(newAlphabet, sym)
		}
	}
	fsm.Alphabet = newAlphabet

	// ========== 4. 清理接受状态：只保留可达的 ==========
	var newAcceptingStates []model.State
	for _, acc := range fsm.AcceptingStates {
		if reachable[acc] {
			newAcceptingStates = append(newAcceptingStates, acc)
		}
	}
	fsm.AcceptingStates = newAcceptingStates

	// ========== 5. 清理转移规则：只保留起始状态在可达集合中的 ==========
	var newTransitions []model.Transition
	for _, t := range fsm.Transitions {
		if reachable[t.FromState] {
			// 同时清理 ToStates：只保留可达的目标（虽然通常都可达，但更安全）
			var validToStates []model.State
			for _, to := range t.ToStates {
				if reachable[to] {
					validToStates = append(validToStates, to)
				}
			}
			if len(validToStates) > 0 {
				t.ToStates = validToStates
				newTransitions = append(newTransitions, t)
			}
			// 注意：即使 ToStates 被清空，说明该转移无效，应丢弃
		}
	}
	fsm.Transitions = newTransitions
}
