package automaton_s

import (
	"errors"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// 验证一个自动机是否有效
func AutomatonValidate(automaton *model.Automaton) error {
	if automaton == nil {
		return errors.New("automaton cannot be nil")
	}

	return automaton.Validate()
}
