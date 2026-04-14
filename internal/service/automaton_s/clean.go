package automaton_s

import (
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// 清理自动机：移除不可达状态和未使用的符号
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
