package automaton_s

import (
	"fmt"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// 返回原 DFA 语言的补集 DFA
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

	if err := aCopy.CompleteDFA(); err != nil {
		return nil, fmt.Errorf("failed to complete DFA: %w", err)
	}

	acceptSet := make(map[model.State]bool)
	for _, s := range aCopy.AcceptingStates {
		acceptSet[s] = true
	}

	newAccepting := []model.State{}
	for _, state := range aCopy.States {
		if !acceptSet[state] {
			newAccepting = append(newAccepting, state)
		}
	}

	aCopy.AcceptingStates = newAccepting

	if err := aCopy.Validate(); err != nil {
		return nil, fmt.Errorf("complemented automaton is invalid: %w", err)
	}

	return aCopy, nil
}
