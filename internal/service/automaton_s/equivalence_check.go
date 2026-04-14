package automaton_s

import (
	"fmt"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// AutomatonEquivalenceCheck 自动机等价检查
// 检查两个自动机是否等价
// 通过转化为最小DFA来检查是否同构
func AutomatonEquivalenceCheck(a1, a2 *model.Automaton) (minDFA1, minDFA2 *model.Automaton, isEquivalent bool) {
	if a1 == nil || a2 == nil {
		if a1 == nil && a2 == nil {
			return nil, nil, true
		}

		return a1, a2, false
	}

	a1.Validate() // 设置 Type
	if a1.Type != model.DFA {
		a1 = NFAToDFA(a1)
	}
	minDFA1 = DFAMinimize(a1)

	a2.Validate() // 设置 Type
	if a2.Type != model.DFA {
		a2 = NFAToDFA(a2)
	}

	minDFA2 = DFAMinimize(a2)

	return minDFA1, minDFA2, AreDFAsIsomorphic(minDFA1, minDFA2)
}

// canonicalize 对DFA进行BFS编号映射（规范重命名）
func canonicalize(dfa *model.Automaton) *model.Automaton {
	// 创建一个新的Automaton实例来存储规范化的结果
	canonicalDFA := &model.Automaton{
		Type:            dfa.Type,
		States:          make([]model.State, 0),
		Alphabet:        dfa.Alphabet,
		InitialState:    "",
		AcceptingStates: make([]model.State, 0),
		Transitions:     make([]model.Transition, 0),
	}

	queue := []model.State{dfa.InitialState}
	stateMap := map[model.State]model.State{}
	newStateName := model.State("q0")

	for len(queue) > 0 {
		currentState := queue[0]
		queue = queue[1:]

		if _, exists := stateMap[currentState]; !exists {
			// 添加新状态名到map和新的Automaton中
			stateMap[currentState] = newStateName
			canonicalDFA.States = append(canonicalDFA.States, newStateName)

			// 设置初始状态
			if currentState == dfa.InitialState {
				canonicalDFA.InitialState = newStateName
			}

			// 添加接受状态
			if contains(dfa.AcceptingStates, currentState) {
				canonicalDFA.AcceptingStates = append(canonicalDFA.AcceptingStates, newStateName)
			}

			// 遍历所有转移
			for _, trans := range dfa.Transitions {
				if trans.FromState == currentState {
					nextState := stateMap[trans.ToStates[0]]
					if nextState == "" {
						queue = append(queue, trans.ToStates[0])
						nextState = model.State(fmt.Sprintf("q%d", len(stateMap)))
					}
					canonicalDFA.Transitions = append(canonicalDFA.Transitions, model.Transition{
						FromState: newStateName,
						Input:     trans.Input,
						ToStates:  []model.State{nextState},
					})
				}
			}
			newStateName = model.State(fmt.Sprintf("q%d", len(stateMap)))
		}
	}

	return canonicalDFA
}

// AreDFAsIsomorphic 检查两个最小DFA是否同构
func AreDFAsIsomorphic(dfa1, dfa2 *model.Automaton) bool {
	// 使用canonicalize函数规范化DFA
	dfa1 = canonicalize(dfa1)
	dfa2 = canonicalize(dfa2)

	// 比较状态数、初始状态、接受状态、转移
	return compareStates(dfa1.States, dfa2.States) &&
		dfa1.InitialState == dfa2.InitialState &&
		compareAcceptingStates(dfa1.AcceptingStates, dfa2.AcceptingStates) &&
		compareTransitions(dfa1.Transitions, dfa2.Transitions)
}

// 辅助函数：用于比较状态列表
func compareStates(s1, s2 []model.State) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

// 辅助函数：用于比较接受状态列表
func compareAcceptingStates(a1, a2 []model.State) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if !contains(a2, a1[i]) {
			return false
		}
	}

	return true
}

// 辅助函数：用于比较转移函数
func compareTransitions(t1, t2 []model.Transition) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i := range t1 {
		if t1[i].FromState != t2[i].FromState ||
			t1[i].Input != t2[i].Input ||
			t1[i].ToStates[0] != t2[i].ToStates[0] {
			return false
		}
	}

	return true
}

// 辅助函数：判断一个字符串是否在字符串数组中
func contains(arr []model.State, str model.State) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}

	return false
}
