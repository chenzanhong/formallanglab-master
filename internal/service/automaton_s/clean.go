package automaton_s

import (
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// Cleanup 清理自动机：移除不可达状态和未使用的符号（正确支持 ε-NFA）
func Cleanup(automaton *model.Automaton) {
	if automaton == nil {
		return
	}

	// Step 1: 构建 ε-closure 缓存（可选优化）
	epsilonClosure := func(state model.State) map[model.State]bool {
		closure := make(map[model.State]bool)
		stack := []model.State{state}
		closure[state] = true

		for len(stack) > 0 {
			s := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			for _, t := range automaton.Transitions {
				if t.FromState == s && t.Input == model.Epsilon {
					for _, to := range t.ToStates {
						if !closure[to] {
							closure[to] = true
							stack = append(stack, to)
						}
					}
				}
			}
		}

		return closure
	}

	// Step 2: 计算所有从初始状态可达的状态（含 ε 路径）
	reachable := make(map[model.State]bool)
	usedSymbols := make(map[model.Symbol]bool)

	// 初始状态的 ε-闭包
	initClosure := epsilonClosure(automaton.InitialState)
	queue := make([]model.State, 0, len(initClosure))
	for s := range initClosure {
		reachable[s] = true
		queue = append(queue, s)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, t := range automaton.Transitions {
			if t.FromState != current {
				continue
			}

			if t.Input == model.Epsilon {
				// ε 转移已在 closure 中处理，此处可跳过（避免重复）
				continue
			}

			// 记录使用的输入符号（排除 ε）
			usedSymbols[t.Input] = true

			for _, to := range t.ToStates {
				toClosure := epsilonClosure(to)
				for s := range toClosure {
					if !reachable[s] {
						reachable[s] = true
						queue = append(queue, s)
					}
				}
			}
		}
	}

	// Step 3: 重建 States
	var newStates []model.State
	for _, s := range automaton.States {
		if reachable[s] {
			newStates = append(newStates, s)
		}
	}
	automaton.States = newStates

	// Step 4: 重建 Alphabet（排除 ε！）
	var newAlphabet []model.Symbol
	for _, sym := range automaton.Alphabet {
		if sym != model.Epsilon && usedSymbols[sym] {
			newAlphabet = append(newAlphabet, sym)
		}
	}
	automaton.Alphabet = newAlphabet

	// Step 5: 重建 AcceptingStates
	var newAccepting []model.State
	for _, s := range automaton.AcceptingStates {
		if reachable[s] {
			newAccepting = append(newAccepting, s)
		}
	}
	automaton.AcceptingStates = newAccepting

	// Step 6: 重建 Transitions（只保留起点和终点都可达的）
	var newTransitions []model.Transition
	for _, t := range automaton.Transitions {
		if !reachable[t.FromState] {
			continue
		}
		var validTo []model.State
		for _, to := range t.ToStates {
			if reachable[to] {
				validTo = append(validTo, to)
			}
		}
		if len(validTo) > 0 {
			t.ToStates = validTo
			newTransitions = append(newTransitions, t)
		}
	}
	automaton.Transitions = newTransitions
}

// Cleanup 清理自动机：移除不可达状态和未使用的符号，未考虑带空转移的
func Cleanup1(automaton *model.Automaton) {
	if automaton == nil {
		return
	}

	// 用于存储有效（可达）状态
	reachable := make(map[model.State]bool)
	// 用于存储被使用的符号
	usedSymbols := make(map[model.Symbol]bool)

	// ========== 1. 使用 BFS 找出所有从初始状态可达的状态 ==========
	queue := []model.State{automaton.InitialState}
	reachable[automaton.InitialState] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 遍历所有从 current 出发的转移
		for _, t := range automaton.Transitions {
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
	for _, s := range automaton.States {
		if reachable[s] {
			newStates = append(newStates, s)
			stateSet[s] = true
		}
	}
	automaton.States = newStates

	// ========== 3. 构建新的字母表（只保留被使用的符号） ==========
	var newAlphabet []model.Symbol
	for _, sym := range automaton.Alphabet {
		if usedSymbols[sym] {
			newAlphabet = append(newAlphabet, sym)
		}
	}
	automaton.Alphabet = newAlphabet

	// ========== 4. 清理接受状态：只保留可达的 ==========
	var newAcceptingStates []model.State
	for _, acc := range automaton.AcceptingStates {
		if reachable[acc] {
			newAcceptingStates = append(newAcceptingStates, acc)
		}
	}
	automaton.AcceptingStates = newAcceptingStates

	// ========== 5. 清理转移规则：只保留起始状态在可达集合中的 ==========
	var newTransitions []model.Transition
	for _, t := range automaton.Transitions {
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
	automaton.Transitions = newTransitions
}
