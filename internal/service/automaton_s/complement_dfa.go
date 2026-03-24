package automaton_s

import (
	"fmt"

	"backend/internal/domain/model"
)

// ComplementDFA 返回原 DFA 语言的补集 DFA
// 要求：输入 a 必须是一个有效的 DFA（可以是非完备的）
// 内部会先调用 CompleteDFA 将其完备化，再翻转接受状态
func ComplementDFA(a *model.Automaton) (*model.Automaton, error) {
	if a == nil {
		return nil, fmt.Errorf("automaton is nil")
	}

	// 深拷贝自动机，避免修改原始对象
	aCopy := a.Clone()
	if aCopy == nil {
		return nil, fmt.Errorf("failed to clone automaton")
	}

	// Step 1: 验证并完备化为 DFA
	if err := aCopy.CompleteDFA(); err != nil {
		return nil, fmt.Errorf("failed to complete DFA: %w", err)
	}

	// Step 2: 构建当前接受状态集合（用于快速查找）
	acceptSet := make(map[model.State]bool)
	for _, s := range aCopy.AcceptingStates {
		acceptSet[s] = true
	}

	// Step 3: 翻转接受状态 —— 所有非接受状态变为接受状态
	newAccepting := []model.State{}
	for _, state := range aCopy.States {
		if !acceptSet[state] {
			newAccepting = append(newAccepting, state)
		}
	}

	aCopy.AcceptingStates = newAccepting

	// 可选：重新验证（调试用）
	if err := aCopy.Validate(); err != nil {
		return nil, fmt.Errorf("complemented automaton is invalid: %w", err)
	}

	return aCopy, nil
}
