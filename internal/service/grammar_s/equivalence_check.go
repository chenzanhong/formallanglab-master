package grammar_s

import (
	"backend/internal/domain/model"
)

// IsEquivalent 判断两个正则文法是否等价
// 通过比较两个文法生成的语言是否相同来判断等价性
// IsEquivalent 判断两个文法是否等价（生成相同的语言）
// 算法思路：
// 1. 默认传入的两个文法都是有效的正则文法（Type3）
// 2. 通过BFS生成两个文法在有限深度内能推导出的所有句子
// 3. 比较两个文法生成的语言集合是否完全相同
//
// 注意：由于文法可能生成无限语言，这里只比较有限深度内的句子
// 限制：最大深度设为8，可能会有误判（不等价的文法可能在深度8内生成相同语言）
//
// 时间复杂度：O(k * n^d)，其中k为产生式数量，n为非终结符数量，d为搜索深度
// 空间复杂度：O(m)，其中m为生成的不同句子数量
/*
	目前还有点问题：
	S -> a A
	A -> c
	与
	S -> a | a c
	被判定为不等价
*/
func IsEquivalent(g1, g2 *model.Grammar) bool {
	// // 先验证两个文法都是有效的正则文法
	// if !IsValidGrammar(g1) || !IsValidGrammar(g2) {
	// 	return false
	// }

	// if TypeDetermine(g1) != Type3 || TypeDetermine(g2) != Type3 {
	// 	return false
	// }

	// 生成两个文法的语言（有限深度）
	maxDepth := 8 // 限制生成字符串的最大长度
	lang1 := generateLanguage(g1, maxDepth)
	lang2 := generateLanguage(g2, maxDepth)

	// 比较两个语言是否相同
	return compareSets(lang1, lang2)
}

// HasAmbiguityOptimized 优化的二义性检查算法
func HasAmbiguityOptimized(g *model.Grammar, maxLen int) (bool, []model.Symbol) {
	// 为每个长度检查所有可能的字符串
	for length := 0; length <= maxLen; length++ {
		strings := generateStringsOfLength(g, length)
		for _, w := range strings {
			count := countDerivationsOptimized(g, w)
			if count > 1 {
				return true, w
			}
		}
	}
	return false, nil
}

// generateStringsOfLength 生成指定长度的所有终结符字符串
func generateStringsOfLength(g *model.Grammar, length int) [][]model.Symbol {
	if length == 0 {
		return [][]model.Symbol{{}}
	}

	var result [][]model.Symbol
	shorter := generateStringsOfLength(g, length-1)

	for _, str := range shorter {
		for _, terminal := range g.Terminals {
			if terminal != model.Epsilon { // 跳过ε
				newStr := make([]model.Symbol, length)
				copy(newStr[:length-1], str)
				newStr[length-1] = terminal
				result = append(result, newStr)
			}
		}
	}
	return result
}

// countDerivationsOptimized 改进的推导计数算法
func countDerivationsOptimized(g *model.Grammar, w []model.Symbol) int {
	n := len(w)

	// dp[A][i][j] = 非终结符A推导出w[i:j]的方式数量
	dp := make(map[model.Symbol][][]int)
	for _, A := range g.NonTerminals {
		table := make([][]int, n+1)
		for i := 0; i <= n; i++ {
			table[i] = make([]int, n+1)
		}
		dp[A] = table
	}

	// 处理空串推导 A → ε（修正：ε可以在任意位置推导）
	for _, p := range g.Productions {
		if len(p.Right) == 0 || (len(p.Right) == 1 && p.Right[0] == model.Epsilon) {
			A := p.Left[0]
			for i := 0; i <= n; i++ {
				dp[A][i][i]++ // ε可以在任意位置推导
			}
		}
	}

	// 处理长度为1的推导 A → a
	for i := 0; i < n; i++ {
		for _, p := range g.Productions {
			if len(p.Right) == 1 && p.Right[0] == w[i] {
				A := p.Left[0]
				dp[A][i][i+1]++
			}
		}
	}

	// 处理长度>1的推导 A → aB
	for length := 2; length <= n; length++ {
		for i := 0; i <= n-length; i++ {
			j := i + length
			for _, p := range g.Productions {
				if len(p.Right) == 2 {
					A := p.Left[0]
					a := p.Right[0]
					B := p.Right[1]

					if i < n && a == w[i] && g.CheckIsNonTerminal(B) {
						dp[A][i][j] += dp[B][i+1][j]
					}
				}
			}
		}
	}

	return dp[g.StartSymbol][0][n]
}

// generateLanguage 生成文法的语言（有限深度）
func generateLanguage(g *model.Grammar, maxDepth int) map[string]bool {
	language := make(map[string]bool)

	// 使用BFS生成所有可能的字符串
	type state struct {
		symbols []model.Symbol
		depth   int
	}

	queue := []state{{symbols: []model.Symbol{g.StartSymbol}, depth: 0}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		key := symbolsToStringJoinSep(curr.symbols)
		if visited[key] {
			continue
		}
		visited[key] = true

		// 如果深度超限，跳过
		if curr.depth > maxDepth {
			continue
		}

		// 如果全是终结符，加入语言
		if isAllTerminals(curr.symbols, g) {
			language[symbolsToStringJoinSep(curr.symbols)] = true
			continue
		}

		// 应用产生式
		for i, sym := range curr.symbols {
			if g.CheckIsNonTerminal(sym) {
				for _, p := range g.Productions {
					if len(p.Left) == 1 && p.Left[0] == sym {
						newSymbols := make([]model.Symbol, 0)
						newSymbols = append(newSymbols, curr.symbols[:i]...)
						newSymbols = append(newSymbols, p.Right...)
						newSymbols = append(newSymbols, curr.symbols[i+1:]...)

						queue = append(queue, state{
							symbols: newSymbols,
							depth:   curr.depth + 1,
						})
					}
				}
				break // 只替换第一个非终结符
			}
		}
	}

	return language
}

// compareSets 比较两个字符串集合是否相同
func compareSets(set1, set2 map[string]bool) bool {
	if len(set1) != len(set2) {
		return false
	}

	for str := range set1 {
		if !set2[str] {
			return false
		}
	}

	return true
}
