package automaton_s

import "backend/internal/domain/model"

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
func AutomatonValidate(automaton *model.Automaton) (bool, error) {
	return automaton.ISValidate()
}
