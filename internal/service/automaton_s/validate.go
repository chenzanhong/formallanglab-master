package automaton_s

import (
	"errors"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// 根据你的项目路径调整

/*
1. 状态集合非空	len(States) > 0 且无空状态名	NFA/DFA
2. 初始状态合法	InitialState ∈ States	NFA/DFA
3. 接受状态合法	∀acc ∈ AcceptingStates: acc ∈ States	NFA/DFA
4. 输入符号合法	∀t.Input ∈ Alphabet	NFA/DFA
5. 转移状态合法	FromState ∈ States, ToStates[i] ∈ States	NFA/DFA
6. DFA 确定性	len(ToStates) == 1 且无重复 (from, input)	仅 DFA
*/

// automatonValidate 验证一个自动机是否有效
// 返回：是否有效，以及错误信息（如果无效）
func AutomatonValidate(automaton *model.Automaton) error {
	// 1. 检查 automaton 是否为 nil
	if automaton == nil {
		return errors.New("automaton cannot be nil")
	}

	// 2. 调用 Automaton 对象的 Validate 方法进行验证
	return automaton.Validate()
}
