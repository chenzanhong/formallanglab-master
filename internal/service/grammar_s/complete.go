package grammar_s

import (
	"backend/internal/domain/model"
)

// CompleteGrammarStructure 完善文法结构，如果某些部分缺失则根据默认规则填充
func CompleteGrammarStructure(grammar *model.Grammar) *model.Grammar {
	// 如果文法结构完整，直接返回
	if grammar.StartSymbol != "" && len(grammar.Terminals) > 0 && len(grammar.NonTerminals) > 0 {
		return grammar
	}

	// 创建新的文法结构
	completeGrammar := &model.Grammar{
		Productions: grammar.Productions,
	}

	// 收集所有符号
	allSymbols := make(map[model.Symbol]bool)
	terminals := make(map[model.Symbol]bool)
	nonTerminals := make(map[model.Symbol]bool)

	// 遍历所有产生式，收集符号
	for _, production := range grammar.Productions {
		// 处理左部
		for _, symbol := range production.Left {
			allSymbols[symbol] = true
			// 判断是否为非终结符（大写字母开头）
			if isNonTerminal(symbol) {
				nonTerminals[symbol] = true
			} else {
				terminals[symbol] = true
			}
		}

		// 处理右部
		for _, symbol := range production.Right {
			if symbol == model.Epsilon || symbol == "" {
				continue // 空符号不计入
			}
			allSymbols[symbol] = true
			// 判断是否为非终结符（大写字母开头）
			if isNonTerminal(symbol) {
				nonTerminals[symbol] = true
			} else {
				terminals[symbol] = true
			}
		}
	}

	// 设置起始符号（默认为第一个产生式的左部，且必须是单个大写字母）
	completeGrammar.StartSymbol = "S" // 默认起始符号
	if len(grammar.Productions) > 0 && len(grammar.Productions[0].Left) > 0 {
		firstSymbol := grammar.Productions[0].Left[0]
		if isNonTerminal(firstSymbol) && len(string(firstSymbol)) == 1 {
			completeGrammar.StartSymbol = firstSymbol
		}
	}

	// 转换为切片
	for symbol := range terminals {
		completeGrammar.Terminals = append(completeGrammar.Terminals, symbol)
	}
	for symbol := range nonTerminals {
		completeGrammar.NonTerminals = append(completeGrammar.NonTerminals, symbol)
	}

	return completeGrammar
}

// isNonTerminal 判断符号是否为非终结符（大写字母开头）
func isNonTerminal(symbol model.Symbol) bool {
	if len(symbol) == 0 {
		return false
	}
	// 大写字母开头认为是非终结符
	firstChar := string(symbol)[0]
	return firstChar >= 'A' && firstChar <= 'Z'
}
