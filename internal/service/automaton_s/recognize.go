package automaton_s

import (
	"backend/internal/domain/model"
	"fmt"
)

// Recognize 判断自动机是否接受输入字符串 str
// 根据自动机类型（DFA或NFA或EPSILON-NFA）选择相应的识别算法
// 默认已经对自动机有效性和类型进行了判断
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

	switch automaton.Type {
	case model.DFA:
		return recognizeDFA(automaton, str)
	case model.NFA:
		return recognizeNFA(automaton, str)
	case model.EpsilonNFA:
		return recognizeEpsilonNFA(automaton, str)
	default:
		return &model.RecognitionResult{
			IsAccepted: false,
			Steps:      []model.RecognitionStep{},
		}, fmt.Errorf("未知的自动机类型")
	}
}

// recognizeDFA 识别DFA是否接受输入字符串 str
func recognizeDFA(automaton *model.Automaton, str string) (*model.RecognitionResult, error) {
	var mustFailed bool
	symbols, err := automaton.SplitString(str)
	if err != nil {
		if len(symbols) == 0 { // 第一个字符就非法了
			return &model.RecognitionResult{
				IsAccepted: false,
				Steps:      []model.RecognitionStep{},
			}, fmt.Errorf("无法识别该字符串，第一个字符就非法了")
		}
		// 不直接返回，而是记录结果，继续处理已分词部分
		mustFailed = true
	}

	currentState := automaton.InitialState
	steps := []model.RecognitionStep{}
	for _, sym := range symbols {
		if automaton.TransMap[currentState][sym] == nil || automaton.TransMap[currentState][sym][0] == "" {
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
		nextState := automaton.TransMap[currentState][sym][0]
		// fmt.Printf("currentState: %s, nextState: %s", currentState, nextState)
		steps = append(steps, model.RecognitionStep{
			State:     currentState,
			Input:     sym,
			NextState: nextState,
		})
		currentState = nextState
	}

	// 检查最终状态是否为接收状态
	for _, st := range automaton.AcceptingStates {
		// fmt.Printf("DFA: %+v		%s\n", automaton.AcceptingStates, currentState)
		if st == currentState {
			return &model.RecognitionResult{
				IsAccepted: !mustFailed,
				Steps:      steps,
			}, nil
		}
	}

	return &model.RecognitionResult{
		IsAccepted: false,
		Steps:      steps,
	}, fmt.Errorf("该自动机不接受输入字符串 %s", str)
}

// ===================== NFA =====================

// recognizeNFA 识别NFA是否接受输入字符串 str，DFS
func recognizeNFA(automaton *model.Automaton, str string) (*model.RecognitionResult, error) {
	// 处理空字符串
	if str == "" {
		if ContainsState(automaton.AcceptingStates, automaton.InitialState) {
			return &model.RecognitionResult{IsAccepted: true}, nil
		}
		return &model.RecognitionResult{IsAccepted: false}, fmt.Errorf("该自动机无法识别空字符串")
	}

	var mustFailed bool
	// 第一步：将字符串分词为符号序列
	symbols, err := automaton.SplitString(str)
	if err != nil {
		if len(symbols) == 0 { // 第一个字符就非法了
			return &model.RecognitionResult{
				IsAccepted: false,
				Steps:      []model.RecognitionStep{},
			}, fmt.Errorf("无法识别该字符串，第一个字符就非法了")
		}
		// 不直接返回，而是记录结果，继续处理已分词部分
		mustFailed = true
	}

	deepest := []model.RecognitionStep{} // 使用空切片初始化，确保JSON序列化为[]而不是null
	pathFound, finalSteps := dfsForNFA(
		automaton,
		automaton.InitialState,
		symbols,
		0,
		[]model.RecognitionStep{},
		make(map[string]bool),
		&deepest,
	)

	if pathFound {
		return &model.RecognitionResult{
			IsAccepted: !mustFailed,
			Steps:      finalSteps,
		}, nil
	}

	return &model.RecognitionResult{
		IsAccepted: false,
		Steps:      deepest,
	}, fmt.Errorf("该自动机不接受输入字符串 %s", str)
}

func dfsForNFA(
	automaton *model.Automaton,
	currentState model.State,
	symbols []model.Symbol,
	pos int,
	steps []model.RecognitionStep,
	visited map[string]bool,
	deepest *[]model.RecognitionStep,
) (bool, []model.RecognitionStep) {

	// 更新最深路径
	if len(steps) > len(*deepest) {
		*deepest = append([]model.RecognitionStep(nil), steps...)
	}

	key := fmt.Sprintf("%s@%d", currentState, pos)
	if visited[key] {
		return false, steps
	}
	visited[key] = true
	defer delete(visited, key)

	// 输入处理完毕：检查是否为接受状态
	if pos == len(symbols) {
		if ContainsState(automaton.AcceptingStates, currentState) {
			return true, steps
		}
		return false, steps
	}

	input := symbols[pos]

	// 尝试所有可能的转移（NFA 允许多个）
	for _, next := range automaton.TransMap[currentState][input] {
		newSteps := append(append([]model.RecognitionStep(nil), steps...), model.RecognitionStep{
			State:     currentState,
			Input:     input,
			NextState: next,
		})
		if found, finalSteps := dfsForNFA(automaton, next, symbols, pos+1, newSteps, visited, deepest); found {
			return true, finalSteps
		}
	}

	return false, steps
}

// ========================== EpsilonNFA ============================

// recognizeEpsilonNFA 识别带ε的NFA是否接受输入字符串 str，DFS
func recognizeEpsilonNFA(automaton *model.Automaton, str string) (*model.RecognitionResult, error) {
	// 处理空字符串
	if str == "" {
		// 如果初始状态本身就是接受状态？
		if ContainsState(automaton.AcceptingStates, automaton.InitialState) {
			return &model.RecognitionResult{IsAccepted: true, Steps: []model.RecognitionStep{}}, nil
		}
		// 否则，查找是否存在 ε-路径到达接受状态
		visited := make(map[model.State]bool)
		if steps, found := findEpsilonPathToAccept(automaton, automaton.InitialState, []model.RecognitionStep{}, visited); found {
			return &model.RecognitionResult{IsAccepted: true, Steps: steps}, nil // 此时的 steps 包含了 ε-转移
		}
		return &model.RecognitionResult{IsAccepted: false, Steps: []model.RecognitionStep{}}, fmt.Errorf("该自动机无法识别空字符串")
	}

	// 记录是否确定一定会识别失败
	var mustFailed bool = false

	// 分词（可能部分成功，即部分识别）
	symbols, err := automaton.SplitString(str)
	if err != nil { // 不可完全分词
		if len(symbols) == 0 { // 第一个字符就非法了
			// 直接返回
			return &model.RecognitionResult{IsAccepted: false, Steps: []model.RecognitionStep{}}, fmt.Errorf("无法识别该字符串，第一个字符就非法了")
		}
		// 部分可分词，记录结果，继续处理已分词部分
		mustFailed = true
	}

	// 第一阶段：快速判断是否接受
	if !isAcceptedForEpsilonNFA(automaton, symbols) {
		// 不直接返回，而是设置结果，确保即时识别失败也有识别路径返回
		mustFailed = true
	}

	deepest := []model.RecognitionStep{} // 使用空切片初始化，确保JSON序列化为[]而不是null
	pathFound, finalSteps := dfsWithEpsilon(
		automaton,
		automaton.InitialState,
		symbols,
		0,
		[]model.RecognitionStep{},
		make(map[string]bool), // 防止无限循环（状态+输入位置作为 key）
		&deepest,
	)
	if pathFound {
		return &model.RecognitionResult{
			IsAccepted: !mustFailed, // 分词部分识别成功，且分词部分为完整的输入字符串
			Steps:      finalSteps,
		}, nil
	}
	// 失败也返回已识别的路径
	return &model.RecognitionResult{
		IsAccepted: false,
		Steps:      deepest,
	}, fmt.Errorf("ε-NFA 无法接受输入字符串 %s", str)
}

// findEpsilonPathToAccept 从currentState出发，通过 ε 转移找到一条到接受状态的路径
// 但对于一开始currentState就是接受态的情况，不会记录自身的转移，返回的[]model.RecognitionStep为currentSteps
func findEpsilonPathToAccept(automaton *model.Automaton, currentState model.State, currentSteps []model.RecognitionStep, visited map[model.State]bool) ([]model.RecognitionStep, bool) {
	if ContainsState(automaton.AcceptingStates, currentState) {
		return currentSteps, true
	}

	if visited[currentState] {
		return currentSteps, false // 避免循环
	}
	visited[currentState] = true
	defer delete(visited, currentState) // 回溯

	// 递归查找 ε 转移到的状态
	for _, next := range automaton.TransMap[currentState][model.Epsilon] {
		if finalSteps, found := findEpsilonPathToAccept(automaton, next, append(append([]model.RecognitionStep(nil), currentSteps...), model.RecognitionStep{
			State:     currentState,
			Input:     model.Epsilon,
			NextState: next,
		}), visited); found {
			return finalSteps, true
		}
	}

	return currentSteps, false // 未找到接受路径
}

// dfsWithEpsilon 执行带路径记录的 DFS
func dfsWithEpsilon(
	automaton *model.Automaton,
	currentState model.State,
	symbols []model.Symbol,
	pos int, // 当前输入符号位置
	steps []model.RecognitionStep,
	visited map[string]bool,
	deepest *[]model.RecognitionStep, // 记录最深路径，用于返回识别失败时的识别路径
) (bool, []model.RecognitionStep) {

	// 更新最新路径
	if len(steps) > len(*deepest) {
		*deepest = append([]model.RecognitionStep(nil), steps...)
	}

	key := fmt.Sprintf("%s@%d", currentState, pos)
	if visited[key] {
		return false, steps // 避免循环
	}
	visited[key] = true
	defer delete(visited, key) // 回溯

	// 输入已处理完输入串：检查是否能通过 ε 转移到达接受状态，并记录路径
	if pos == len(symbols) {
		visitedEps := make(map[model.State]bool)
		if acceptSteps, found := findEpsilonPathToAccept(automaton, currentState, steps, visitedEps); found {
			return true, acceptSteps
		}
		return false, steps
	}

	input := symbols[pos]

	// 先尝试 ε 转移
	for _, next := range automaton.TransMap[currentState][model.Epsilon] {
		// 使用 append 复制 steps，避免 slices 共享底层数组导致的路径污染
		newSteps := append(append([]model.RecognitionStep(nil), steps...), model.RecognitionStep{
			State:     currentState,
			Input:     model.Epsilon,
			NextState: next,
		})
		if found, finalSteps := dfsWithEpsilon(automaton, next, symbols, pos, newSteps, visited, deepest); found {
			return true, finalSteps
		}
	}

	// 尝试消耗当前输入符号
	for _, next := range automaton.TransMap[currentState][input] {
		newSteps := append(append([]model.RecognitionStep(nil), steps...), model.RecognitionStep{
			State:     currentState,
			Input:     input,
			NextState: next,
		})
		if found, finalSteps := dfsWithEpsilon(automaton, next, symbols, pos+1, newSteps, visited, deepest); found {
			return true, finalSteps
		}
	}

	return false, steps
}

// computeEpsilonClosureWithPath 返回闭包和对应的 ε 转移步骤（从给定 states 出发）
func computeEpsilonClosureWithPath(states []model.State, automaton model.Automaton) ([]model.State, []model.RecognitionStep) {
	closure := make([]model.State, 0)
	visited := make(map[model.State]bool)
	steps := []model.RecognitionStep{} // 使用空切片初始化，确保JSON序列化为[]而不是null
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

		for _, target := range automaton.TransMap[current][model.Epsilon] {
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

	return closure, steps
}

// isAccepted 先判断是否接受
func isAcceptedForEpsilonNFA(automaton *model.Automaton, symbols []model.Symbol) bool {
	// 先计算初始状态的 ε-闭包
	current := computeEpsilonClosureWithMap(automaton, []model.State{automaton.InitialState})

	for _, sym := range symbols {
		next := getEpsilonNFANextStatesFromMap(automaton, current, sym)
		if len(next) == 0 {
			return false
		}
		current = computeEpsilonClosureWithMap(automaton, next)
	}

	for _, s := range current {
		if ContainsState(automaton.AcceptingStates, s) {
			return true
		}
	}
	return false
}

// func findAcceptPathForEpsilonNFA(automaton *model.Automaton, symbols []model.Symbol) ([]model.RecognitionStep, bool) {
// 	initialClosure, initialSteps := computeEpsilonClosureWithPath([]model.State{automaton.InitialState}, *automaton)

// 	// 空字符串情况
// 	if len(symbols) == 0 {
// 		for _, s := range initialClosure {
// 			if ContainsState(automaton.AcceptingStates, s) {
// 				return initialSteps, true
// 			}
// 		}
// 		return initialSteps, false
// 	}

// 	// 从初始闭包的每个状态出发尝试
// 	for _, s := range initialClosure {
// 		if found, steps := dfsWithEpsilon(automaton, s, symbols, 0, initialSteps, make(map[string]bool), ); found {
// 			return steps, true
// 		}
// 	}
// 	return initialSteps, false
// }
