package automaton_s

import (
	"backend/internal/domain/model"
	"fmt"
)

// Recognize 判断自动机是否接受输入字符串 str
// 根据自动机类型（DFA或NFA）选择相应的识别算法
//
// 算法说明：
// 对于DFA（确定性有限自动机）：
//  1. 从初始状态开始，依次读取输入字符串中的每个符号
//  2. 根据当前状态和输入符号，按照转移函数转移到下一个状态
//  3. 读取完所有符号后，检查当前状态是否为接受状态
//
// 对于NFA（非确定性有限自动机）：
//  1. 从初始状态开始，维护一个可能的当前状态集合
//  2. 依次读取输入字符串中的每个符号
//  3. 对于当前状态集合中的每个状态，计算在该符号下的所有可能转移
//  4. 将所有可能的转移目标合并为新的当前状态集合
//  5. 读取完所有符号后，检查当前状态集合中是否包含任何接受状态
//
// Symbol 是 string 类型，使用 map + 固定长度集合进行高效匹配
func Recognize(automaton *model.Automaton, str string) (*model.RecognitionResult, error) {
	if automaton == nil || len(automaton.States) == 0 {
		return &model.RecognitionResult{
			IsAccepted: false,
			Steps:      []model.RecognitionStep{},
		}, fmt.Errorf("自动机状态为空")
	}

	// 初始化 TransMap 以提高识别效率
	automaton.InitTransMap()

	if automaton.IsDFA {
		return recognizeDFA(automaton, str)
	} else {
		return recognizeNFA(automaton, str)
	}
}

// recognizeDFA 识别DFA是否接受输入字符串 str
func recognizeDFA(automaton *model.Automaton, str string) (*model.RecognitionResult, error) {
	symbols, err := automaton.SplitString(str)
	if err != nil {
		return &model.RecognitionResult{
			IsAccepted: false,
			Steps:      []model.RecognitionStep{},
		}, fmt.Errorf("DFA 识别失败：无法将输入字符串分词 -> %w", err)
	}
	fmt.Printf("symbols: %+v", symbols)

	currentState := automaton.InitialState
	steps := []model.RecognitionStep{}
	for _, sym := range symbols {
		nextState := automaton.TransMap[currentState][sym][0]
		if nextState == "" {
			if len(steps) > 0 {
				steps = append(steps, model.RecognitionStep{
					State:     steps[len(steps)-1].NextState,
					Input:     sym,
					NextState: "",
				})
			}
			return &model.RecognitionResult{
				IsAccepted: false,
				Steps:      steps,
			}, fmt.Errorf("当前状态 %s 下输入 %s 没有相关的有效转移", currentState, sym)
		}
		fmt.Printf("currentState: %s, nextState: %s", currentState, nextState)
		steps = append(steps, model.RecognitionStep{
			State:     currentState,
			Input:     sym,
			NextState: nextState,
		})
		currentState = nextState
	}

	// 检查最终状态是否为接收状态
	for _, st := range automaton.AcceptingStates {
		fmt.Printf("DFA: %+v		%s\n", automaton.AcceptingStates, currentState)
		if st == currentState {
			return &model.RecognitionResult{
				IsAccepted: true,
				Steps:      steps,
			}, nil
		}
	}

	return &model.RecognitionResult{
		IsAccepted: false,
		Steps:      steps,
	}, fmt.Errorf("输入字符串 %s 被完整识别，但是未到达接收状态", str)
}

// recognizeNFA 识别NFA是否接受输入字符串 str，DFS
func recognizeNFA(automaton *model.Automaton, str string) (*model.RecognitionResult, error) {
	// 第一步：将字符串分词为符号序列
	symbols, err := automaton.SplitString(str)
	if err != nil {
		return &model.RecognitionResult{
			IsAccepted: false,
			Steps:      []model.RecognitionStep{},
		}, fmt.Errorf("NFA 识别失败：无法将输入字符串分词 -> %w", err)
	}

	// 当前可能处于的状态集合（NFA 的核心：状态集合）
	currentStates := []model.State{automaton.InitialState}
	steps := []model.RecognitionStep{}

	// 逐个处理每个符号
	for _, sym := range symbols {
		var nextStates []model.State
		found := false

		// 对当前每个可能的状态，查找该符号的转移
		for _, state := range currentStates {
			targets := automaton.TransMap[state][sym]
			if len(targets) > 0 {
				found = true
				// 将目标状态加入 nextStates（去重）
				for _, t := range targets {
					if !containsState(nextStates, t) {
						nextStates = append(nextStates, t)
					}
					steps = append(steps, model.RecognitionStep{
						State:     state,
						Input:     sym,
						NextState: t,
					})
				}
			}
		}

		// 如果没有任何状态能处理当前符号
		if !found {
			return &model.RecognitionResult{
				IsAccepted: false,
				Steps: append(steps, model.RecognitionStep{
					State:     steps[len(steps)-1].NextState,
					Input:     sym,
					NextState: "",
				}),
			}, fmt.Errorf("NFA 识别失败：在输入符号 '%s' 时，当前状态集合中没有状态可以转移", sym)
		}

		// 更新当前状态集合
		currentStates = nextStates
	}

	// 最终：检查当前状态集合中是否有任意一个接受状态
	for _, s := range currentStates {
		fmt.Printf("NFA: %+v		%+v\n", automaton.AcceptingStates, currentStates)
		if containsState(automaton.AcceptingStates, s) {
			return &model.RecognitionResult{
				IsAccepted: true,
				Steps:      steps,
			}, nil
		}
	}

	return &model.RecognitionResult{
		IsAccepted: false,
		Steps:      steps,
	}, fmt.Errorf("输入字符串 %s 被完整识别，但未到达任何接受状态", str)
}

// // recognizeNFA_e 识别带ε的NFA是否接受输入字符串 str，DFS
func recognizeNFA_e(automaton model.Automaton, str string) (model.RecognitionResult, error) {
	symbols, err := automaton.SplitString(str)
	if err != nil {
		return model.RecognitionResult{
			IsAccepted: false,
			Steps:      []model.RecognitionStep{},
		}, fmt.Errorf("ε-NFA 识别失败：无法将输入字符串分词 -> %w", err)
	}

	// 先计算初始 ε-闭包，并记录初始 ε 路径
	initialClosure, initialSteps := computeEpsilonClosureWithPath([]model.State{automaton.InitialState}, automaton)

	// 尝试从每个初始闭包中的状态开始 DFS（因为 NFA 可以“同时”处于多个状态，但我们找一条路径）
	for _, startState := range initialClosure {
		var currentSteps []model.RecognitionStep
		currentSteps = append(currentSteps, initialSteps...) // 包含初始 ε 转移

		pathFound, finalSteps := dfsWithEpsilon(
			automaton,
			startState,
			symbols,
			0,
			currentSteps,
			make(map[string]bool), // 防止无限循环（状态+输入位置作为 key）
		)
		if pathFound {
			return model.RecognitionResult{
				IsAccepted: true,
				Steps:      finalSteps,
			}, nil
		}
	}

	// 如果所有路径都失败，尝试构造一条失败路径（可选）
	// 这里简化：返回初始闭包 + 第一个无法转移的步骤
	lastState := automaton.InitialState
	if len(initialClosure) > 0 {
		lastState = initialClosure[0]
	}
	if len(symbols) == 0 {
		// 空串情况
		for _, s := range initialClosure {
			if containsState(automaton.AcceptingStates, s) {
				return model.RecognitionResult{IsAccepted: true, Steps: initialSteps}, nil
			}
		}
		return model.RecognitionResult{IsAccepted: false, Steps: initialSteps}, nil
	}

	steps := initialSteps
	steps = append(steps, model.RecognitionStep{
		State:     lastState,
		Input:     symbols[0],
		NextState: "",
	})
	return model.RecognitionResult{
		IsAccepted: false,
		Steps:      steps,
	}, fmt.Errorf("ε-NFA 无法接受输入字符串 %s", str)
}

// dfsWithEpsilon 执行带路径记录的 DFS
func dfsWithEpsilon(
	automaton model.Automaton,
	currentState model.State,
	symbols []model.Symbol,
	pos int,
	steps []model.RecognitionStep,
	visited map[string]bool,
) (bool, []model.RecognitionStep) {

	key := fmt.Sprintf("%s@%d", currentState, pos)
	if visited[key] {
		return false, steps // 避免循环
	}
	visited[key] = true
	defer delete(visited, key) // 回溯

	// 如果已处理完所有输入
	if pos == len(symbols) {
		// 检查当前状态是否为接受状态（或其 ε-闭包中包含接受状态）
		closure, _ := computeEpsilonClosureWithPath([]model.State{currentState}, automaton)
		for _, s := range closure {
			if containsState(automaton.AcceptingStates, s) {
				// 把 closure 中的 ε 转移也加进来（如果有的话）
				_, epsSteps := computeEpsilonClosureWithPath([]model.State{currentState}, automaton)
				finalSteps := append(steps, epsSteps...)
				return true, finalSteps
			}
		}
		return false, steps
	}

	input := symbols[pos]

	// Option 1: 先尝试 ε 转移（不消耗输入）
	for _, t := range automaton.Transitions {
		if t.FromState == currentState && t.Input == model.Epsilon {
			for _, next := range t.ToStates {
				newSteps := append(steps, model.RecognitionStep{
					State:     currentState,
					Input:     model.Epsilon,
					NextState: next,
				})
				if found, finalSteps := dfsWithEpsilon(automaton, next, symbols, pos, newSteps, visited); found {
					return true, finalSteps
				}
			}
		}
	}

	// Option 2: 尝试消耗当前输入符号
	for _, t := range automaton.Transitions {
		if t.FromState == currentState && t.Input == input {
			for _, next := range t.ToStates {
				newSteps := append(steps, model.RecognitionStep{
					State:     currentState,
					Input:     input,
					NextState: next,
				})
				if found, finalSteps := dfsWithEpsilon(automaton, next, symbols, pos+1, newSteps, visited); found {
					return true, finalSteps
				}
			}
		}
	}

	return false, steps
}

// computeEpsilonClosureWithPath 返回闭包和对应的 ε 转移步骤（从给定 states 出发）
func computeEpsilonClosureWithPath(states []model.State, automaton model.Automaton) ([]model.State, []model.RecognitionStep) {
	closure := make([]model.State, 0)
	visited := make(map[model.State]bool)
	var steps []model.RecognitionStep
	queue := append([]model.State(nil), states...)

	for _, s := range states {
		if !visited[s] {
			visited[s] = true
			closure = append(closure, s)
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, t := range automaton.Transitions {
			if t.FromState == current && t.Input == model.Epsilon {
				for _, target := range t.ToStates {
					if !visited[target] {
						visited[target] = true
						closure = append(closure, target)
						queue = append(queue, target)
						steps = append(steps, model.RecognitionStep{
							State:     current,
							Input:     model.Epsilon,
							NextState: target,
						})
					}
				}
			}
		}
	}

	return closure, steps
}
