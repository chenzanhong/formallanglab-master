package grammar_s

import (
	"fmt"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/service/automaton_s"
)

// 采用理论方法：转为最小 DFA，判断是否同构。
// 默认输入的文法为正则文法
func GrammarIsEquivalent(g1, g2 *model.Grammar) (bool, error) {
	if g1 == nil || g2 == nil {
		return false, fmt.Errorf("有自动机为空，请检查输入")
	}

	// 转自动机
	fa1, err := RegularGrammarToFA(g1)
	if err != nil {
		return false, err
	}
	fa2, err := RegularGrammarToFA(g2)
	if err != nil {
		return false, err
	}
	// FA 简化
	automaton_s.Cleanup(fa1)
	automaton_s.Cleanup(fa2)

	// 判断 FA 类型
	fa1.Validate()
	fa2.Validate()

	// 转 DFA
	if fa1.Type != model.DFA {
		fa1 = automaton_s.NFAToDFA(fa1)
	}
	if fa2.Type != model.DFA {
		fa2 = automaton_s.NFAToDFA(fa2)
	}

	// DFA 最小化
	fa1 = automaton_s.DFAMinimize(fa1)
	fa2 = automaton_s.DFAMinimize(fa2)

	// DFA 同构判断
	return automaton_s.AreDFAsIsomorphic(fa1, fa2), nil
}
