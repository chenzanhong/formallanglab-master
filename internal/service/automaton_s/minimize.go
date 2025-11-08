package automaton_s

import (
	"backend/internal/domain/model"
	"fmt"
	"sort"
	"strings"
)

// DFAMinimize 是 DFA 最小化的统一入口（默认使用表格填充法）
// 可根据需要切换为 Hopcroft 算法（如通过配置 flag）
func DFAMinimize(automaton *model.Automaton) *model.Automaton {
	// 可选：未来可加 algo := config.GetMinimizationAlgo()
	// return minimizeByHopcroft(Automaton)
	return minimizeByTableFilling(automaton)
}

// minimizeByTableFilling 使用表格填充法（Table-Filling Method）对 DFA 进行最小化
func minimizeByTableFilling(automaton *model.Automaton) *model.Automaton {
	// 步骤 1: 去除不可达状态
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
			IsDFA:           true,
		}
	}

	// 构建仅含可达状态的子自动机（用于后续转移查询）
	reduced := &model.Automaton{
		States:          states,
		Alphabet:        automaton.Alphabet,
		Transitions:     filterTransitions(automaton.Transitions, reachable),
		InitialState:    automaton.InitialState,
		AcceptingStates: newAccepting,
		IsDFA:           true,
	}

	// 步骤 2 & 3: 表格填充法标记可区分状态对
	distinguishable := markDistinguishablePairs(reduced, acceptingSet)

	// 步骤 4: 构建等价类
	classes := findEquivalenceClasses(reduced.States, distinguishable)

	// 步骤 5: 构建最小 DFA
	return buildMinimizedDFA(reduced, classes, acceptingSet)
}

// minimizeByHopcroft 使用 Hopcroft 算法对 DFA 进行最小化（备用实现）
func minimizeByHopcroft(automaton *model.Automaton) *model.Automaton {
	// Step 1: 移除不可达状态
	reachable := getReachableStates(automaton)
	reduced := &model.Automaton{
		States:          filterStates(automaton.States, reachable),
		Alphabet:        automaton.Alphabet,
		Transitions:     filterTransitions(automaton.Transitions, reachable),
		InitialState:    automaton.InitialState,
		AcceptingStates: filterStates(automaton.AcceptingStates, reachable),
		IsDFA:           true,
	}

	if len(reduced.States) <= 1 {
		return reduced
	}

	// 构建反向转移图（用于 Hopcroft）
	revTrans := buildReverseTransitions(reduced)

	// 初始划分：接受状态 vs 非接受状态
	acceptingSet := make(map[model.State]bool)
	for _, s := range reduced.AcceptingStates {
		acceptingSet[s] = true
	}

	var P [][]model.State
	P = append(P, reduced.AcceptingStates)
	nonAccepting := []model.State{}
	for _, s := range reduced.States {
		if !acceptingSet[s] {
			nonAccepting = append(nonAccepting, s)
		}
	}
	if len(nonAccepting) > 0 {
		P = append(P, nonAccepting)
	}

	// W 是待处理的划分块（初始为接受状态块）
	W := [][]model.State{reduced.AcceptingStates}

	// Hopcroft 主循环
	for len(W) > 0 {
		A := W[0]
		W = W[1:]

		for _, c := range reduced.Alphabet {
			// 找到所有能通过 c 转移到 A 中状态的前驱状态
			X := make(map[model.State]bool)
			for _, state := range A {
				if preds, ok := revTrans[state][c]; ok {
					for _, p := range preds {
						X[p] = true
					}
				}
			}

			// 对 P 中每个块 Y，检查是否需要分裂
			var newP [][]model.State
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

					// 更新 W
					found := false
					for i, w := range W {
						if equalSet(w, Y) {
							// 替换 W 中的 Y 为两个新块（较小者入 W）
							W = append(W[:i], W[i+1:]...)
							if len(YInX) <= len(YNotInX) {
								W = append(W, YInX)
							} else {
								W = append(W, YNotInX)
							}
							found = true
							break
						}
					}
					if !found {
						// 将较小的新块加入 W
						if len(YInX) <= len(YNotInX) {
							W = append(W, YInX)
						} else {
							W = append(W, YNotInX)
						}
					}
				}
			}
			P = newP
		}
	}

	// 构建等价类映射
	classes := P
	stateToClass := make(map[model.State][]model.State)
	for _, cls := range classes {
		for _, s := range cls {
			stateToClass[s] = cls
		}
	}

	// 构建新 DFA
	return buildMinimizedDFAFromClasses(reduced, classes, acceptingSet, stateToClass)
}

// ------------------ 辅助函数 ------------------
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
			return t.ToStates[0] // DFA only has one next state
		}
	}
	return "" // No transition (should not happen in complete DFA)
}

func markDistinguishablePairs(automaton *model.Automaton, acceptingSet map[model.State]bool) map[[2]model.State]bool {
	states := automaton.States
	distinguishable := make(map[[2]model.State]bool)

	// 初始化：接受 vs 非接受
	for i, p := range states {
		for _, q := range states[i+1:] {
			pAcc := acceptingSet[p]
			qAcc := acceptingSet[q]
			if pAcc != qAcc {
				if p <= q {
					distinguishable[[2]model.State{p, q}] = true
				} else {
					distinguishable[[2]model.State{q, p}] = true
				}
			}
		}
	}

	// 迭代标记
	changed := true
	for changed {
		changed = false
		for i, p := range states {
			for _, q := range states[i+1:] {
				pair := [2]model.State{p, q}
				if distinguishable[pair] {
					continue
				}
				for _, a := range automaton.Alphabet {
					nextP := getDFANextState(automaton, p, a)
					nextQ := getDFANextState(automaton, q, a)
					if nextP == "" || nextQ == "" {
						if nextP == "" && nextQ == "" { // 都没有转移，先跳过
							continue
						} else { // 有一个没有转移，标记为不同
							distinguishable[pair] = true
							changed = true
							break
						}
					}
					var nextPair [2]model.State
					if nextP <= nextQ {
						nextPair = [2]model.State{nextP, nextQ}
					} else {
						nextPair = [2]model.State{nextQ, nextP}
					}
					if distinguishable[nextPair] {
						distinguishable[pair] = true
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
		IsDFA:           true,
	}
}

func joinStates(states []model.State) string {
	strs := make([]string, len(states))
	for i, s := range states {
		strs[i] = string(s)
	}
	return fmt.Sprintf("%s", strings.Join(strs, ","))
}

func getReachableStates(automaton *model.Automaton) map[model.State]bool {
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

func findEquivalenceClasses(states []model.State, distinguishable map[[2]model.State]bool) [][]model.State {
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
			pair := [2]model.State{p, q}
			if !distinguishable[pair] {
				union(p, q)
			}
		}
	}

	// 收集等价类
	classes := make(map[model.State][]model.State)
	for _, s := range states {
		root := find(s)
		classes[root] = append(classes[root], s)
	}

	var result [][]model.State
	for _, cls := range classes {
		sort.Slice(cls, func(i, j int) bool {
			return string(cls[i]) < string(cls[j])
		}) // 排序，确保确定性，但不是必须的
		result = append(result, cls)
	}

	return result
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
		IsDFA:           true,
	}
}
