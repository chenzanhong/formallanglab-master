package automaton_s

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// DFAMinimize 是 DFA 最小化的统一入口（默认使用 Hopcroft 算法）
// 备选算法为表格填充算法
func DFAMinimize(automaton *model.Automaton) *model.Automaton {
	if automaton == nil {
		return nil
	}

	minimized, _ := minimizeByHopcroft(automaton)

	return minimized
	// return minimizeByTableFilling(automaton)
}

func DFAMinimizeWithProcess(automaton *model.Automaton) (*model.Automaton, *model.MinimizationProcess) {
	if automaton == nil {
		return nil, nil
	}

	return minimizeByHopcroft(automaton)
}

// 使用表格填充法（Table-Filling Method）对 DFA 进行最小化
func minimizeByTableFilling(automaton *model.Automaton) *model.Automaton {
	if automaton == nil {
		return nil
	}

	// 去除不可达状态
	reachable := getReachableStates(automaton)
	var states []model.State
	for s := range reachable {
		states = append(states, s)
	}

	newAccepting := []model.State{}
	acceptingSet := make(map[model.State]bool)
	for _, s := range automaton.AcceptingStates {
		if reachable[s] {
			newAccepting = append(newAccepting, s)
			acceptingSet[s] = true
		}
	}

	if len(states) <= 1 {
		return &model.Automaton{
			States:          states,
			Alphabet:        automaton.Alphabet,
			Transitions:     filterTransitions(automaton.Transitions, reachable),
			InitialState:    automaton.InitialState,
			AcceptingStates: newAccepting,
			Type:            model.DFA,
		}
	}

	// 构建仅含可达状态的子自动机（用于后续转移查询）
	reduced := &model.Automaton{
		States:          states,
		Alphabet:        automaton.Alphabet,
		Transitions:     filterTransitions(automaton.Transitions, reachable),
		InitialState:    automaton.InitialState,
		AcceptingStates: newAccepting,
		Type:            model.DFA,
	}

	// 表格填充法标记可区分状态对
	distinguishable := markDistinguishablePairs(reduced, acceptingSet)

	// 构建等价类
	classes := findEquivalenceClasses(reduced.States, distinguishable)

	// 构建最小 DFA
	return buildMinimizedDFA(reduced, classes, acceptingSet)
}

// minimizeByHopcroft 使用 Hopcroft 算法对 DFA 进行最小化
func minimizeByHopcroft(automaton *model.Automaton) (*model.Automaton, *model.MinimizationProcess) {
	if automaton == nil {
		return nil, nil
	}
	// 先确保 automaton 为完备的，但是可能添加陷阱状态"dead_state"，用户体验不好，这里先不采用完备化
	// if err := automaton.CompleteDFA(); err != nil { // 一般不会 err，handler 层确保了是有效 DFA。
	// 	return automaton, nil
	// }
	// 创建过程记录器
	process := &model.MinimizationProcess{
		Steps: []model.MinimizationStep{},
	}
	stepCount := 0

	// 移除不可达状态
	reachable := getReachableStates(automaton)
	actions := []string{"移除不可达状态"}
	reduced := &model.Automaton{
		States:          filterStates(automaton.States, reachable),
		Alphabet:        automaton.Alphabet,
		Transitions:     filterTransitions(automaton.Transitions, reachable),
		InitialState:    automaton.InitialState,
		AcceptingStates: filterStates(automaton.AcceptingStates, reachable),
		Type:            model.DFA,
	}

	stepCount++
	partition := [][]model.State{}
	stateList := make([]model.State, 0, len(reduced.States))
	stateList = append(stateList, reduced.States...)
	partition = append(partition, stateList)
	process.Steps = append(process.Steps, model.MinimizationStep{
		Step:          stepCount,
		Partition:     partition,
		Actions:       actions,
		AutomatonFlow: reduced.ToReactFlow(),
	})

	if len(reduced.States) <= 1 {
		stepCount++
		finalPartition := [][]model.State{}
		finalStateList := make([]model.State, 0, len(reduced.States))
		finalStateList = append(finalStateList, reduced.States...)
		finalPartition = append(finalPartition, finalStateList)
		process.Steps = append(process.Steps, model.MinimizationStep{
			Step:          stepCount,
			Partition:     finalPartition,
			Actions:       []string{"状态数量 <= 1，无需进一步最小化"},
			AutomatonFlow: automaton.ToReactFlow(),
		})

		return reduced, process
	}

	// 构建反向转移图（用于 Hopcroft）
	revTrans := buildReverseTransitions(reduced)
	actions = []string{"构建反向转移图"}

	// 初始划分：接受状态 vs 非接受状态
	acceptingSet := make(map[model.State]bool)
	for _, s := range reduced.AcceptingStates {
		acceptingSet[s] = true
	}

	var P [][]model.State
	if len(reduced.AcceptingStates) > 0 {
		P = append(P, reduced.AcceptingStates)
	}
	nonAccepting := []model.State{}
	for _, s := range reduced.States {
		if !acceptingSet[s] {
			nonAccepting = append(nonAccepting, s)
		}
	}
	if len(nonAccepting) > 0 {
		P = append(P, nonAccepting)
	}

	stepCount++
	actions = append(actions, "根据接受/非接受状态进行初始划分")
	initialPartition := [][]model.State{}
	for _, block := range P {
		stateList := make([]model.State, 0, len(block))
		stateList = append(stateList, block...)
		initialPartition = append(initialPartition, stateList)
	}

	tempStateToClass := make(map[model.State][]model.State)
	for _, cls := range P {
		for _, s := range cls {
			tempStateToClass[s] = cls
		}
	}
	currentDFA := buildMinimizedDFAFromClasses(reduced, P, acceptingSet, tempStateToClass)
	process.Steps = append(process.Steps, model.MinimizationStep{
		Step:          stepCount,
		Partition:     initialPartition,
		Actions:       actions,
		AutomatonFlow: currentDFA.ToReactFlow(),
	})

	// W 是待处理的划分块（初始为接受状态块）
	W := [][]model.State{reduced.AcceptingStates}

	// Hopcroft 主循环
	for len(W) > 0 {
		A := W[0]
		W = W[1:]

		for _, c := range reduced.Alphabet {
			actions = []string{fmt.Sprintf("选择划分块 %v 进行处理", A)}
			// 找到所有能通过 c 转移到 A 中状态的前驱状态
			X := make(map[model.State]bool)
			preStates := make([]model.State, 0)
			for _, state := range A {
				if preds, ok := revTrans[state][c]; ok {
					for _, p := range preds {
						X[p] = true
						preStates = append(preStates, p)
					}
				}
			}
			actions = append(actions, fmt.Sprintf("找到所有能通过符号 %s 转移到 %v 中状态的前驱状态：%v", c, A, preStates))

			// 对 P 中每个块 Y，检查是否需要分裂
			var newP [][]model.State
			splitted := false
			for _, Y := range P {
				YInX := []model.State{}
				YNotInX := []model.State{}

				for _, s := range Y {
					if X[s] {
						YInX = append(YInX, s)
					} else {
						YNotInX = append(YNotInX, s)
					}
				}

				if len(YInX) == 0 || len(YNotInX) == 0 {
					// 无需分裂
					newP = append(newP, Y)
				} else {
					// 分裂 Y 为 YInX 和 YNotInX
					newP = append(newP, YInX, YNotInX)
					splitted = true

					// 更新 W
					found := false
					for i, w := range W {
						if equalSet(w, Y) {
							// 替换 W 中的 Y 为两个新块
							W = append(W[:i], W[i+1:]...)
							// // 较小者入 W（前提是需 automaton 是完备的，但是可能会引入 dead_state，用户体验可能不好）
							// if len(YInX) <= len(YNotInX) {
							// 	W = append(W, YInX)
							// } else {
							// 	W = append(W, YNotInX)
							// }
							// 这里直接都添加，不考虑是否完备
							W = append(W, YInX)
							W = append(W, YNotInX)
							found = true

							break
						}
					}
					if !found {
						// // 较小者入 W（前提是需 automaton 是完备的，但是可能会引入 dead_state，用户体验可能不好）
						// if len(YInX) <= len(YNotInX) {
						// 	W = append(W, YInX)
						// } else {
						// 	W = append(W, YNotInX)
						// }
						// 这里直接都添加，不考虑是否完备
						W = append(W, YInX)
						W = append(W, YNotInX)
					}
					actions = append(actions, fmt.Sprintf("将划分块 %v 分裂为新划分块 %v 和 %v，并将新划分块加入待处理列表", Y, YInX, YNotInX))
				}
			}

			// 如果有分裂发生，记录新的划分
			newPartition := [][]model.State{}
			if splitted {
				for _, block := range newP {
					stateList := make([]model.State, 0, len(block))
					stateList = append(stateList, block...)
					newPartition = append(newPartition, stateList)
				}
			} else {
				// 无需分裂，保持原划分，保持为上一步的划分
				if stepCount > 0 {
					newPartition = process.Steps[stepCount-1].Partition
				}
				actions = append(actions, "没有划分块需要分裂")
			}

			stepCount++
			tempStateToClass := make(map[model.State][]model.State)
			for _, cls := range P {
				for _, s := range cls {
					tempStateToClass[s] = cls
				}
			}
			currentDFA = buildMinimizedDFAFromClasses(reduced, P, acceptingSet, tempStateToClass)
			process.Steps = append(process.Steps, model.MinimizationStep{
				Step:          stepCount,
				Partition:     newPartition,
				Actions:       actions,
				AutomatonFlow: currentDFA.ToReactFlow(),
			})
			P = newP
		}
	}

	// 记录最终步骤：最小化完成
	stepCount++
	finalPartition := P
	process.Steps = append(process.Steps, model.MinimizationStep{
		Step:          stepCount,
		Partition:     finalPartition,
		Actions:       []string{"划分不再变化，最小化完成"},
		AutomatonFlow: currentDFA.ToReactFlow(),
	})

	// 构建等价类映射
	classes := P
	stateToClass := make(map[model.State][]model.State)
	for _, cls := range classes {
		for _, s := range cls {
			stateToClass[s] = cls
		}
	}

	// 构建新 DFA
	minimizedDFA := buildMinimizedDFAFromClasses(reduced, classes, acceptingSet, stateToClass)

	// 把可能添加的陷阱状态去除，经过 buildMinimizedDFAFromClasses 后，model.DeadState 被[]包裹
	// UnCompleteDFA(minimizedDFA, "["+model.SinkState+"]")

	return minimizedDFA, process
}

func filterStates(states []model.State, reachable map[model.State]bool) []model.State {
	var res []model.State
	for _, s := range states {
		if reachable[s] {
			res = append(res, s)
		}
	}

	return res
}

func filterTransitions(trans []model.Transition, reachable map[model.State]bool) []model.Transition {
	var res []model.Transition
	for _, t := range trans {
		if reachable[t.FromState] && reachable[t.ToStates[0]] {
			res = append(res, t)
		}
	}

	return res
}

func getDFANextState(automaton *model.Automaton, from model.State, input model.Symbol) model.State {
	for _, t := range automaton.Transitions {
		if t.FromState == from && t.Input == input {
			return t.ToStates[0]
		}
	}

	return ""
}

func markDistinguishablePairs(automaton *model.Automaton, acceptingSet map[model.State]bool) map[model.State]map[model.State]bool {
	if automaton == nil {
		return nil
	}
	states := automaton.States
	distinguishable := make(map[model.State]map[model.State]bool)
	for _, s := range states {
		distinguishable[s] = make(map[model.State]bool)
	}

	for i, p := range states {
		for _, q := range states[i+1:] {
			if acceptingSet[p] != acceptingSet[q] {
				distinguishable[p][q] = true
				distinguishable[q][p] = true
			}
		}
	}

	changed := true
	for changed {
		changed = false
		for i, p := range states {
			for _, q := range states[i+1:] {
				if distinguishable[p][q] {
					continue
				}
				for _, a := range automaton.Alphabet {
					nextP := getDFANextState(automaton, p, a)
					nextQ := getDFANextState(automaton, q, a)
					// 处理缺失转移：一个有转移、一个没有 → 可区分
					if (nextP == "") != (nextQ == "") {
						distinguishable[p][q] = true
						distinguishable[q][p] = true
						changed = true

						break
					}
					if nextP == "" && nextQ == "" {
						continue // 都无转移，不影响区分
					}
					// 若后继状态可区分，则当前也可区分
					if distinguishable[nextP][nextQ] {
						distinguishable[p][q] = true
						distinguishable[q][p] = true
						changed = true

						break
					}
				}
			}
		}
	}

	return distinguishable
}

func buildReverseTransitions(automaton *model.Automaton) map[model.State]map[model.Symbol][]model.State {
	if automaton == nil {
		return nil
	}
	rev := make(map[model.State]map[model.Symbol][]model.State)
	for _, t := range automaton.Transitions {
		to := t.ToStates[0]
		if rev[to] == nil {
			rev[to] = make(map[model.Symbol][]model.State)
		}
		rev[to][t.Input] = append(rev[to][t.Input], t.FromState)
	}

	return rev
}

func equalSet(a, b []model.State) bool {
	if len(a) != len(b) {
		return false
	}
	setA := make(map[model.State]bool)
	for _, s := range a {
		setA[s] = true
	}
	for _, s := range b {
		if !setA[s] {
			return false
		}
	}

	return true
}

func buildMinimizedDFAFromClasses(
	automaton *model.Automaton,
	classes [][]model.State,
	acceptingSet map[model.State]bool,
	stateToClass map[model.State][]model.State,
) *model.Automaton {
	// 为每个类生成唯一名称
	classToName := make(map[string]model.State)
	var newStates []model.State
	var newAccepting []model.State

	for _, cls := range classes {
		sort.Slice(cls, func(i, j int) bool {
			return string(cls[i]) < string(cls[j])
		})
		nameStr := fmt.Sprintf("[%s]", joinStates(cls))
		name := model.State(nameStr)
		classToName[nameStr] = name
		newStates = append(newStates, name)

		// 检查是否为接受状态
		for _, s := range cls {
			if acceptingSet[s] {
				newAccepting = append(newAccepting, name)
				break
			}
		}
	}

	// 构建转移
	var newTransitions []model.Transition
	for _, cls := range classes {
		fromState := cls[0]
		fromName := model.State(fmt.Sprintf("[%s]", joinStates(cls)))

		for _, a := range automaton.Alphabet {
			next := getDFANextState(automaton, fromState, a)
			if next == "" {
				continue
			}
			nextCls := stateToClass[next]
			toName := model.State(fmt.Sprintf("[%s]", joinStates(nextCls)))

			newTransitions = append(newTransitions, model.Transition{
				FromState: fromName,
				Input:     a,
				ToStates:  []model.State{toName},
			})
		}
	}

	initialCls := stateToClass[automaton.InitialState]
	initialName := model.State(fmt.Sprintf("[%s]", joinStates(initialCls)))

	return &model.Automaton{
		States:          newStates,
		Alphabet:        automaton.Alphabet,
		Transitions:     newTransitions,
		InitialState:    initialName,
		AcceptingStates: newAccepting,
		Type:            model.DFA,
	}
}

func joinStates(states []model.State) string {
	strs := make([]string, len(states))
	for i, s := range states {
		strs[i] = string(s)
	}

	return strings.Join(strs, ",")
}

func getReachableStates(automaton *model.Automaton) map[model.State]bool {
	if automaton == nil {
		return nil
	}
	reachable := make(map[model.State]bool)
	queue := []model.State{automaton.InitialState}
	reachable[automaton.InitialState] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, t := range automaton.Transitions {
			if t.FromState == current && !reachable[t.ToStates[0]] {
				reachable[t.ToStates[0]] = true
				queue = append(queue, t.ToStates[0])
			}
		}
	}

	return reachable
}

func findEquivalenceClasses(states []model.State, distinguishable map[model.State]map[model.State]bool) [][]model.State {
	// 初始化：每个状态自成一类
	parent := make(map[model.State]model.State)
	for _, s := range states {
		parent[s] = s
	}

	// Union-Find 的 find 操作
	var find func(x model.State) model.State
	find = func(x model.State) model.State {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}

		return parent[x]
	}

	// Union 操作
	union := func(x, y model.State) {
		px, py := find(x), find(y)
		if px != py {
			// 字典序小的作为根
			if px < py {
				parent[py] = px
			} else {
				parent[px] = py
			}
		}
	}

	// 所有“不可区分”的状态对进行 union
	for i, p := range states {
		for _, q := range states[i+1:] {
			if !distinguishable[p][q] {
				union(p, q)
			}
		}
	}

	// 收集等价类
	classesMap := make(map[model.State][]model.State)
	for _, s := range states {
		root := find(s)
		classesMap[root] = append(classesMap[root], s)
	}

	var classes [][]model.State
	for _, cls := range classesMap {
		sort.Slice(cls, func(i, j int) bool {
			return string(cls[i]) < string(cls[j])
		}) // 排序，确保确定性，但不是必须的
		classes = append(classes, cls)
	}

	return classes
}

func buildMinimizedDFA(
	automaton *model.Automaton,
	classes [][]model.State,
	acceptingSet map[model.State]bool,
) *model.Automaton {
	// 映射：状态 → 所属类
	stateToClass := make(map[model.State][]model.State)
	for _, cls := range classes {
		for _, s := range cls {
			stateToClass[s] = cls
		}
	}

	// 新状态名：[q0,q1]
	var newStates []model.State
	var newAccepting []model.State
	var newTransitions []model.Transition

	for _, cls := range classes {
		name := fmt.Sprintf("%v", cls) // 简单表示，如 "[q0 q1]"
		newStates = append(newStates, model.State(name))

		// 如果类中任一状态是接受状态，则新状态是接受状态
		for _, s := range cls {
			if acceptingSet[s] {
				newAccepting = append(newAccepting, model.State(name))
				break
			}
		}
	}

	// 构建转移
	for _, cls := range classes {
		fromState := cls[0] // 任取一个代表状态
		fromName := fmt.Sprintf("%v", cls)

		for _, a := range automaton.Alphabet {
			next := getDFANextState(automaton, fromState, a)
			if next == "" {
				continue // 无转移（理论上不应发生）
			}
			nextCls := stateToClass[next]
			toName := fmt.Sprintf("%v", nextCls)

			newTransitions = append(newTransitions, model.Transition{
				FromState: model.State(fromName),
				Input:     a,
				ToStates:  []model.State{model.State(toName)},
			})
		}
	}

	// 确定初始状态
	initialCls := stateToClass[automaton.InitialState]
	initialName := fmt.Sprintf("%v", initialCls)

	return &model.Automaton{
		States:          newStates,
		Alphabet:        automaton.Alphabet,
		Transitions:     newTransitions,
		InitialState:    model.State(initialName),
		AcceptingStates: newAccepting,
		Type:            model.DFA,
	}
}
