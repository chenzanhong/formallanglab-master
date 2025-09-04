package grammar_s

import (
	"backend/internal/domain/model"
)

// IsValidGrammar 判断一个 Grammar 是否为有效文法
func IsValidGrammar(g *model.Grammar) bool {
	if g == nil {
		return false
	}

	// 1. 非终结符和终结符都不能为空
	if len(g.NonTerminals) == 0 {
		return false
	}
	if len(g.Terminals) == 0 {
		return false
	}

	// 2. 起始符号必须是非终结符
	if g.StartSymbol == "" {
		return false
	}
	if !isInSet(g.NonTerminals, g.StartSymbol) {
		return false
	}

	// 3. 非终结符和终结符不能有交集
	for sym := range g.NonTerminals {
		if isInSet(g.Terminals, sym) {
			return false
		}
	}

	// 4. 各自集合内部无重复符号（map 天然保证 key 不重复，所以这一步由 JSON 解析保证）
	// 因此无需额外检查重复，Go 的 map[string]struct{} 天然去重

	// 5. 至少有一个产生式
	if len(g.Productions) == 0 {
		return false
	}

	// 5. 每个产生式的左部至少包含一个非终结符
	for _, p := range g.Productions {
		if len(p.Left) == 0 {
			return false // 左部不能为空
		}
		hasNonTerminal := false
		for _, sym := range p.Left {
			if isInSet(g.NonTerminals, sym) {
				hasNonTerminal = true
				break
			}
		}
		if !hasNonTerminal {
			return false // 左部没有非终结符
		}
	}

	// 全部通过
	return true
}

// 工具函数：判断符号是否在集合中
func isInSet(set map[model.Symbol]struct{}, sym model.Symbol) bool {
	_, exists := set[sym]
	return exists
}
