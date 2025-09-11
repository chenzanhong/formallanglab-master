package fsm_s

import (
	"backend/internal/domain/model"
	"fmt"
	"sort"
)

// DFAMinimize 对给定的 DFA 进行极小化 reduce
// 假设输入 fsm 是一个有效的 DFA（IsDFA == true）
func DFAMinimize(fsm *model.Automaton) *model.Automaton {

	// 步骤 1: 去除不可达状态
	reachable := getReachableStates(fsm)
	var states []model.State
	for s := range reachable {
		states = append(states, s)
	}

	// 过滤状态、转移、接受状态
	newStates := states
	newAccepting := []model.State{}
	for _, s := range fsm.AcceptingStates {
		if reachable[s] {
			newAccepting = append(newAccepting, s)
		}
	}

	// 创建新自动机（仅含可达状态）
	reduced := &model.Automaton{
		States:          newStates,
		Alphabet:        fsm.Alphabet,
		Transitions:     []model.Transition{},
		InitialState:    fsm.InitialState,
		AcceptingStates: newAccepting,
		IsDFA:           true,
	}

	// 如果只剩一个状态，直接返回
	if len(newStates) <= 1 {
		return reduced
	}

	// 步骤 2: 标记可区分状态对（表格填充法）
	// 使用 map[[2]State]bool，true 表示可区分
	distinguishable := make(map[[2]model.State]bool)

	// 初始化：所有 (accepting, non-accepting) 对标记为可区分
	acceptingSet := make(map[model.State]bool)
	for _, s := range newAccepting {
		acceptingSet[s] = true
	}

	for i, p := range newStates {
		for _, q := range newStates[i+1:] {
			_, pIsAccept := acceptingSet[p]
			_, qIsAccept := acceptingSet[q]

			if pIsAccept != qIsAccept {
				// 一个是接受，一个不是 → 可区分
				distinguishable[[2]model.State{p, q}] = true
			}
		}
	}

	// 步骤 3: 迭代标记
	changed := true
	for changed {
		changed = false
		for i, p := range newStates {
			for _, q := range newStates[i+1:] {
				pq := [2]model.State{p, q}
				if distinguishable[pq] {
					continue // 已标记
				}

				// 检查所有输入符号
				for _, a := range fsm.Alphabet {
					// 获取 p 和 q 在 a 上的转移目标
					nextP := getDFANextState(reduced, p, a)
					nextQ := getDFANextState(reduced, q, a)

					// 如果任一无法转移，跳过（理论上 DFA 应完备）
					if nextP == "" || nextQ == "" {
						continue
					}

					// 确保有序：小的在前
					var nextPair [2]model.State
					if nextP <= nextQ {
						nextPair = [2]model.State{nextP, nextQ}
					} else {
						nextPair = [2]model.State{nextQ, nextP}
					}

					// 如果 (nextP, nextQ) 可区分，则 (p,q) 也可区分
					if distinguishable[nextPair] {
						distinguishable[pq] = true
						changed = true
						break // 当前符号已可区分
					}
				}
			}
		}
	}

	// 步骤 4: 合并不可区分状态
	// 使用并查集或 DFS 找连通分量（不可区分状态构成等价类）
	equivalenceClasses := findEquivalenceClasses(newStates, distinguishable)

	// 步骤 5: 构建新的极小 DFA
	return buildMinimizedDFA(reduced, equivalenceClasses, acceptingSet)
}

func buildMinimizedDFA(
	fsm *model.Automaton,
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

		for _, a := range fsm.Alphabet {
			next := getDFANextState(fsm, fromState, a)
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
	initialCls := stateToClass[fsm.InitialState]
	initialName := fmt.Sprintf("%v", initialCls)

	return &model.Automaton{
		States:          newStates,
		Alphabet:        fsm.Alphabet,
		Transitions:     newTransitions,
		InitialState:    model.State(initialName),
		AcceptingStates: newAccepting,
		IsDFA:           true,
	}
}

func getReachableStates(fsm *model.Automaton) map[model.State]bool {
	reachable := make(map[model.State]bool)
	queue := []model.State{fsm.InitialState}
	reachable[fsm.InitialState] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, t := range fsm.Transitions {
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
