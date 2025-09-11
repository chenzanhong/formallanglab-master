package fsm_s

import (
	"backend/internal/domain/model" // 根据你的项目路径调整
	"errors"
	"fmt"
)

/*
1. 状态集合非空	len(States) > 0 且无空状态名	NFA/DFA
2. 初始状态合法	InitialState ∈ States	NFA/DFA
3. 接受状态合法	∀acc ∈ AcceptingStates: acc ∈ States	NFA/DFA
4. 输入符号合法	∀t.Input ∈ Alphabet	NFA/DFA
5. 转移状态合法	FromState ∈ States, ToStates[i] ∈ States	NFA/DFA
6. DFA 确定性	len(ToStates) == 1 且无重复 (from, input)	仅 DFA
*/

// FSMValidate 验证一个自动机是否有效
// 返回：是否有效，以及错误信息（如果无效）
func FSMValidate(fsm *model.Automaton) (bool, error) {
	// 使用 map 提高查找效率
	stateSet := make(map[model.State]bool)
	alphabetSet := make(map[model.Symbol]bool)

	// 1. 检查状态集合不能为空
	if len(fsm.States) == 0 {
		return false, errors.New("状态集合不能为空")
	}

	// 构建状态集合
	for _, s := range fsm.States {
		if s == "" {
			return false, errors.New("状态名不能为空字符串")
		}
		stateSet[s] = true
	}

	// 2. 检查初始状态是否在状态集合中
	if !stateSet[fsm.InitialState] {
		return false, fmt.Errorf("初始状态 '%s' 不在状态集合中", fsm.InitialState)
	}

	// 3. 检查接受状态是否都是合法状态
	for _, acc := range fsm.AcceptingStates {
		if !stateSet[acc] {
			return false, fmt.Errorf("接受状态 '%s' 不在状态集合中", acc)
		}
	}

	// 4. 构建字母表集合
	for _, sym := range fsm.Alphabet {
		alphabetSet[sym] = true
	}
	alphabetSet[""] = true // 允许空转移

	// 5. 验证所有转移规则
	// 如果是 DFA，我们需要检查：每个 (fromState, input) 只能有一个 toState
	_, err := fsm.CheckIsDFA()
	if err != nil {
		return false, err
	}

	// 6. （可选）如果是 DFA，检查是否每个状态对每个输入都有定义（完备性）
	// if fsm.IsDFA {
	// 	for _, state := range fsm.States {
	// 		for _, input := range fsm.Alphabet {
	// 			key := string(state) + "|" + string(input)
	// 			if _, exists := dfaTransitionMap[key]; !exists {
	// 				return fmt.Errorf("DFA 不完备：状态 '%s' 对输入 '%s' 缺少转移", state, input)
	// 			}
	// 		}
	// 	}
	// }

	return true,nil // 有效
}
