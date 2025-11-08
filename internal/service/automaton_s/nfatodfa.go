package automaton_s

import (
	"backend/internal/domain/model"
	"sort"
	"strings"
)

// NFAToDFA 将 NFA 转换为等价的 DFA， 子集构造法
// 假设输入 NFA 已经有效，ε 转移用 model.Epsilon 表示
func NFAToDFA(nfa *model.Automaton) *model.Automaton {
	// 构建 DFA 字母表
	var alphabet []model.Symbol = nfa.Alphabet

	// 状态集合 → 名称映射函数，如 {q0,q1}
	stateName := func(states []model.State) model.State {
		sorted := make([]string, len(states))
		for i, s := range states {
			sorted[i] = string(s)
		}
		sort.Strings(sorted)
		return model.State("{" + strings.Join(sorted, ",") + "}")
	}

	// 记录已生成的状态名称，避免重复
	seen := make(map[string]bool)
	var dfaStates []model.State
	var dfaTransitions []model.Transition
	var dfaAccepting []model.State

	// 1. 初始状态：NFA 初始状态的 ε-闭包
	initialSet := epsilonClosure(nfa, []model.State{nfa.InitialState})
	initialName := stateName(initialSet)
	seen[string(initialName)] = true
	dfaStates = append(dfaStates, initialName)

	// BFS 队列：存储 NFA 状态集合
	queue := [][]model.State{initialSet}

	// 2. 子集构造主循环
	for len(queue) > 0 {
		currentSet := queue[0]
		queue = queue[1:]
		currentName := stateName(currentSet)

		// 3. 对每个非 ε 输入符号，计算转移
		for _, sym := range alphabet {
			var nextSet []model.State

			// 遍历 currentSet 中每个状态
			for _, state := range currentSet {
				for _, t := range nfa.Transitions {
					if t.FromState == state && t.Input == sym {
						// 添加所有目标状态（去重）
						for _, target := range t.ToStates {
							if !containsState(nextSet, target) {
								nextSet = append(nextSet, target)
							}
						}
					}
				}
			}

			// 4. 对 nextSet 取 ε-闭包
			closure := epsilonClosure(nfa, nextSet)
			if len(closure) == 0 {
				continue // 无有效状态，跳过
			}

			nextName := stateName(closure)

			// 5. 添加 DFA 转移
			dfaTransitions = append(dfaTransitions, model.Transition{
				FromState: currentName,
				Input:     sym,
				ToStates:  []model.State{nextName}, // DFA 单目标
			})

			// 6. 如果是新状态，加入队列
			if !seen[string(nextName)] {
				seen[string(nextName)] = true
				dfaStates = append(dfaStates, nextName)
				queue = append(queue, closure)
			}
		}

		// 7. 判断当前状态是否为接受状态
		for _, s := range currentSet {
			if containsState(nfa.AcceptingStates, s) {
				dfaAccepting = append(dfaAccepting, currentName)
				break
			}
		}
	}

	return &model.Automaton{
		States:          dfaStates,
		Alphabet:        alphabet,
		Transitions:     dfaTransitions,
		InitialState:    initialName,
		AcceptingStates: dfaAccepting,
		IsDFA:           true,
	}
}


