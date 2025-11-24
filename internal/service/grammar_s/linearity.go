package grammar_s

import (
	"backend/internal/domain/model"
)

// DetermineLinearity 判断一个正则文法是右线性还是左线性。
// 前提：调用者必须确保 g 是有效的正则文法（即 IsRegular(g) == true）。
// 返回值：
//   - model.RightLinear: 右线性
//   - model.LeftLinear: 左线性
//   - model.InvalidGrammar: 非正则或无效
//
// 默认输入的文法为正则文法
func DetermineLinearity(g *model.Grammar) model.GrammarLinearity {
	if g == nil {
		return model.InvalidLinearity
	}

	hasRightLinear := false
	hasLeftLinear := false

	for _, p := range g.Productions {
		right := p.Right

		// 空产生式或单一终结符/ε：两者都允许
		if len(right) == 0 || (len(right) == 1 && (right[0] == model.Epsilon || g.CheckIsTerminal(right[0]))) {
			continue
		}

		if isRightLinearProduction(right, g) {
			hasRightLinear = true
		} else if isLeftLinearProduction(right, g) {
			hasLeftLinear = true
		} else {
			// 理论上不会走到这里（因为前提是正则文法）
			return model.InvalidLinearity
		}
	}

	// 必须全部一致：不能混合
	if hasRightLinear && !hasLeftLinear {
		return model.RightLinear
	}
	if hasLeftLinear && !hasRightLinear {
		return model.LeftLinear
	}
	// 全是 ε/终结符产生式（无非终结符出现在右部），视为右线性（惯例）
	if !hasRightLinear && !hasLeftLinear {
		return model.RightLinear
	}
	// 混合了左右线性 → 不合法（但 IsRegular 应已排除）
	return model.InvalidLinearity
}
