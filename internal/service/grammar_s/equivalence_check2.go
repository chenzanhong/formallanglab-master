package grammar_s

import (
	"backend/internal/domain/model"
	"backend/internal/service/automaton_s"
	"backend/internal/service/convert_s"
)

// 采用理论方法：转为最小DFA，判断是否同构。
// 默认输入的文法为正则文法
func GrammarIsEquivalent(g1, g2 *model.Grammar) bool {
	if g1 == nil || g2 == nil {
		return false
	}
	// 判断左/右线性
	linear1 := DetermineLinearity(g1)
	linear2 := DetermineLinearity(g2)

	// 转自动机
	fa1 := convert_s.RegularGrammarToFA(g1, linear1)
	fa2 := convert_s.RegularGrammarToFA(g2, linear2)

	// FA 简化
	automaton_s.Cleanup(fa1)
	automaton_s.Cleanup(fa2)

	// 判断FA类型
	fa1.ISValidate()
	fa2.ISValidate()

	// 转DFA
	if fa1.Type != model.DFA {
		fa1 = automaton_s.NFAToDFA(fa1)
	}
	if fa2.Type != model.DFA {
		fa2 = automaton_s.NFAToDFA(fa2)
	}

	// DFA最小化
	fa1 = automaton_s.DFAMinimize(fa1)
	fa2 = automaton_s.DFAMinimize(fa2)

	// DFA同构判断
	return automaton_s.AreDFAsIsomorphic(fa1, fa2)
}
