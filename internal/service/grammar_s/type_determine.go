package grammar_s

import "backend/internal/domain/model"

const (
	Type3       = 3
	Type2       = 2
	Type1       = 1
	Type0       = 0
	TypeInvalid = -1
)

// TypeDetermine 判断文法的Chomsky类型（0、1、2、3型）
// 算法思路：
// 1. 首先检查文法是否有效
// 2. 判断是否为上下文无关文法（CFG，2型）：
//    - 所有产生式的左部都必须是单个非终结符
// 3. 如果是CFG，进一步判断是否为正则文法（3型）：
//    - 所有产生式的右部必须符合正则文法的模式（右线性或左线性）
// 4. 如果不是CFG，判断是否为上下文有关文法（CSG，1型）：
//    - 除S→ε外，所有产生式都满足|左部|≤|右部|
// 5. 如果以上都不满足，则为0型文法（短语结构文法）
//
// 返回值：
// Type3: 正则文法（右线性或左线性）
// Type2: 上下文无关文法
// Type1: 上下文有关文法
// Type0: 短语结构文法
// TypeInvalid: 无效文法
//
// 时间复杂度：O(n*m)，其中n为产生式数量，m为产生式平均长度
// 空间复杂度：O(1)
func TypeDetermine(g *model.Grammar) int {
	if g == nil || !IsValidGrammar(g) {
		return TypeInvalid
	}

	// === 第一步：判断是否为 2型或3型（即：所有产生式左部是否都只有一个非终结符）===
	if isContextFreeForm(g) {
		// 左部都是单个非终结符 → 可能是 Type2 或 Type3
		if isRegular(g) {
			return Type3
		}
		return Type2
	}

	// === 第二步：不是 CFG → 判断是 Type1 还是 Type0 ===
	if isContextSensitive(g) {
		return Type1
	}
	return Type0
}

// isContextFreeForm 检查是否所有产生式左部都是“单个非终结符”
func isContextFreeForm(g *model.Grammar) bool {
	for _, p := range g.Productions {
		if len(p.Left) != 1 || !g.CheckIsNonTerminal(p.Left[0]) {
			return false
		}
	}
	return true
}

// isRegular 检查是否为正则文法（右线性或左线性）
func isRegular(g *model.Grammar) bool {
	hasRightLinear := false
	hasLeftLinear := false

	for _, p := range g.Productions {
		// left := p.Left[0]
		right := p.Right

		// 允许空产生式
		if len(right) == 1 && right[0] == model.Epsilon {
			continue
		}

		// 检测是否右线性
		if isRightLinearProduction(right, g) {
			hasRightLinear = true
			continue
		}

		// 检测是否左线性
		if isLeftLinearProduction(right, g) {
			hasLeftLinear = true
			continue
		}

		// 都不是，非正则
		return false
	}

	// 必须全部右线性或全部左线性
	return hasRightLinear != hasLeftLinear
}

func isRightLinearProduction(right []model.Symbol, g *model.Grammar) bool {
	// 情况1：全为终结符（A → w）
	allTerminals := true
	for _, sym := range right {
		if sym != model.Epsilon && !g.CheckIsTerminal(sym) {
			allTerminals = false
			break
		}
	}
	if allTerminals {
		return true
	}

	// 情况2：非终结符在最右端（A → wX）
	return checkWithNonTerminalAt(right, g, len(right)-1)
}

func isLeftLinearProduction(right []model.Symbol, g *model.Grammar) bool {
	// 情况1：全为终结符（A → w）
	allTerminals := true
	for _, sym := range right {
		if sym != model.Epsilon && !g.CheckIsTerminal(sym) {
			allTerminals = false
			break
		}
	}
	if allTerminals {
		return true
	}

	// 情况2：非终结符在最左端（A → Xw）
	return checkWithNonTerminalAt(right, g, 0)
}

// checkWithNonTerminalAt 检查：在 index 位置是非终结符，其余位置都是终结符
func checkWithNonTerminalAt(right []model.Symbol, g *model.Grammar, index int) bool {
	n := len(right)
	if index < 0 || index >= n {
		return false
	}

	// 检查 index 位置：必须是非终结符
	if !g.CheckIsNonTerminal(right[index]) {
		return false
	}

	// 检查其他位置：必须是终结符（允许 ε）
	for i := 0; i < n; i++ {
		if i == index {
			continue
		}
		sym := right[i]
		if sym != model.Epsilon && !g.CheckIsTerminal(sym) {
			return false
		}
	}
	return true
}

// isContextSensitive 检查是否为上下文有关文法（CSG）
// 条件：所有产生式满足 len(Left) <= len(Right)，除非是 S → ε
func isContextSensitive(g *model.Grammar) bool {
	for _, p := range g.Productions {
		leftLen := len(p.Left)
		rightLen := len(p.Right)

		// 对于空产生式，由于用model.Epsilon表示，长度为1，所以不用专门判断
		// 一般规则：|左部| <= |右部|
		if leftLen > rightLen {
			return false
		}
	}
	return true
}
