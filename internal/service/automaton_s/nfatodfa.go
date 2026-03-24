package automaton_s

import (
	"fmt"
	"sort"
	"strings"

	"backend/internal/domain/model"
)

// NFAToDFA 将 NFA 转换为等价的 DFA， 子集构造法
// 假设输入 NFA 已经有效，ε 转移用 model.Epsilon 表示
func NFAToDFA(nfa *model.Automaton) *model.Automaton {
	if nfa == nil {
		return nil
	}
	if nfa.Type != model.NFA && nfa.Type != model.EpsilonNFA {
		return nil
	}

	// 构建 DFA 字母表
	var alphabet []model.Symbol
	for _, sym := range nfa.Alphabet {
		if sym != model.Epsilon { // 确保没有 ε
			alphabet = append(alphabet, sym)
		}
	}

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
	initialSet := computeEpsilonClosure(nfa, []model.State{nfa.InitialState})
	initialState := stateName(initialSet)
	seen[string(initialState)] = true
	dfaStates = append(dfaStates, initialState)

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
							if !ContainsState(nextSet, target) {
								nextSet = append(nextSet, target)
							}
						}
					}
				}
			}

			// 4. 对 nextSet 取 ε-闭包
			closure := computeEpsilonClosure(nfa, nextSet)
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
			if ContainsState(nfa.AcceptingStates, s) {
				dfaAccepting = append(dfaAccepting, currentName)
				break
			}
		}
	}

	return &model.Automaton{
		States:          dfaStates,
		Alphabet:        alphabet,
		Transitions:     dfaTransitions,
		InitialState:    initialState,
		AcceptingStates: dfaAccepting,
		Type:            model.DFA,
	}
}

func NFAToDFAWithProcess(nfa *model.Automaton) *model.NFADeterminizationProcess {
	if nfa == nil {
		return nil
	}
	if nfa.Type != model.NFA && nfa.Type != model.EpsilonNFA {
		return nil
	}

	// 构建 DFA 字母表
	var alphabet []model.Symbol
	for _, sym := range nfa.Alphabet {
		if sym != model.Epsilon { // 确保没有 ε
			alphabet = append(alphabet, sym)
		}
	}

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
	initialSet := computeEpsilonClosure(nfa, []model.State{nfa.InitialState})
	initialState := stateName(initialSet)
	seen[string(initialState)] = true
	dfaStates = append(dfaStates, initialState)

	// BFS 队列：存储 NFA 状态集合
	queue := [][]model.State{initialSet}

	// 记录转换过程
	steps := []model.NFADeterminizationStep{} // 使用空切片初始化，确保JSON序列化为[]而不是null
	step := 0
	// 初始快照
	initialDFA := &model.Automaton{
		States:       []model.State{initialState},
		Alphabet:     alphabet,
		Transitions:  make([]model.Transition, 0),
		InitialState: initialState,
		AcceptingStates: func() []model.State {
			for _, s := range initialSet {
				if ContainsState(nfa.AcceptingStates, s) {
					return []model.State{initialState}
				}
			}

			return nil
		}(),
		Type: model.DFA,
	}

	// 生成初始步骤描述
	initialDesc := fmt.Sprintf("步骤 %d: 初始状态为 NFA 初始状态 %s 的 ε-闭包 → %s",
		step, nfa.InitialState, initialState)

	var finalAutomaton *model.Automaton
	finalAutomaton = initialDFA
	steps = append(steps, model.NFADeterminizationStep{
		Step:          step,
		Description:   initialDesc,
		AutomatonFlow: initialDFA.ToReactFlow(),
	})
	step++

	// 2. 子集构造主循环
	for len(queue) > 0 {
		currentSet := queue[0]
		queue = queue[1:]
		currentName := stateName(currentSet)

		// 用于收集本轮新增的转移和新状态（用于描述）
		var processedTransitions []string

		// 3. 对每个非 ε 输入符号，计算转移
		for _, sym := range alphabet {
			var nextSet []model.State

			// 遍历 currentSet 中每个状态
			for _, state := range currentSet {
				for _, t := range nfa.Transitions {
					if t.FromState == state && t.Input == sym {
						// 添加所有目标状态（去重）
						for _, target := range t.ToStates {
							if !ContainsState(nextSet, target) {
								nextSet = append(nextSet, target)
							}
						}
					}
				}
			}

			// 4. 对 nextSet 取 ε-闭包
			closure := computeEpsilonClosure(nfa, nextSet)
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

			// 记录这条转移用于描述
			processedTransitions = append(processedTransitions,
				fmt.Sprintf("δ(%s, %s) = %s", currentName, sym, nextName))
		}

		// 7. 判断当前状态是否为接受状态
		for _, s := range currentSet {
			if ContainsState(nfa.AcceptingStates, s) {
				if !ContainsState(dfaAccepting, currentName) {
					dfaAccepting = append(dfaAccepting, currentName)
				}

				break
			}
		}

		// 生成描述
		var desc string
		if len(processedTransitions) == 0 {
			desc = fmt.Sprintf("步骤 %d: 状态 %s 无有效转移", step, currentName)
		} else {
			desc = fmt.Sprintf("步骤 %d: 处理状态 %s 的转移:\n%s",
				step, currentName, strings.Join(processedTransitions, "\n"))
		}

		// 生成当前 DFA 的完整快照
		currentDFA := &model.Automaton{
			States:          copySlice(dfaStates), // 深拷贝
			Alphabet:        alphabet,
			Transitions:     copyTransitions(dfaTransitions),
			InitialState:    initialState,
			AcceptingStates: copySlice(dfaAccepting),
			Type:            model.DFA,
		}
		finalAutomaton = currentDFA
		steps = append(steps, model.NFADeterminizationStep{
			Step:          step,
			Description:   desc,
			AutomatonFlow: currentDFA.ToReactFlow(),
		})
		step++
	}

	return &model.NFADeterminizationProcess{
		FinalAutomaton: finalAutomaton,
		Steps:          steps,
	}
}

func copySlice(slice []model.State) []model.State {
	copy := append([]model.State{}, slice...)
	return copy
}

func copyTransitions(transitions []model.Transition) []model.Transition {
	copy := append([]model.Transition{}, transitions...)
	return copy
}
