package grammar_s

import (
	"strings"

	"backend/internal/domain/model"
)

// symbolsToStringJoinSep 将符号切片转为字符串，用于 map 的 key
// 注意：正确处理ε符号，确保[ a ]和[ a, ε ]被正确区分
func symbolsToStringJoinSep(symbols []model.Symbol) string {
	parts := make([]string, len(symbols))
	for i, sym := range symbols {
		if sym == model.Epsilon {
			parts[i] = "ε"
		} else {
			parts[i] = string(sym)
		}
	}

	return strings.Join(parts, "|")
}

// IsAmbiguousRegular 检查一个正则文法是否二义
// 二义性定义：存在至少一个句子（终结符串）可以通过两种或以上不同的最左推导得到
// 算法思路：
// 1. 使用广度优先搜索(BFS)生成文法能推导出的所有句子
// 2. 限制搜索深度防止无限循环（设置为2倍非终结符数量，最大不超过10）
// 3. 记录每个句子被推导出的次数
// 4. 一旦发现某个句子被推导出超过一次，立即返回true
// 5. 若遍历完所有可能推导仍未发现重复句子，则返回false
//
// 时间复杂度：O(k * n^d)，其中k为产生式数量，n为非终结符数量，d为搜索深度
// 空间复杂度：O(m)，其中m为生成的不同句子数量
//
// 假设输入文法已经是 Type3（正则文法）
func IsAmbiguousRegular(g *model.Grammar) (bool, error) {
	maxDepth := 2 * len(g.NonTerminals)
	if maxDepth > 10 {
		maxDepth = 10
	}

	sentenceCount := make(map[string]int)
	visited := make(map[string]bool) // 记录已处理的状态，避免重复处理

	type state struct {
		symbols []model.Symbol
		depth   int // 记录实际推导深度
	}

	queue := []state{{symbols: []model.Symbol{g.StartSymbol}, depth: 0}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// 深度检查（防止无限）
		if curr.depth > maxDepth {
			continue
		}

		// 检查是否为句子（全终结符）
		if isAllTerminals(curr.symbols, g) {
			// 对于句子，将ε替换为空字符串，因为ε在句子中表示空
			sentence := ""
			for _, sym := range curr.symbols {
				if sym != model.Epsilon {
					sentence += string(sym)
				}
			}
			sentenceCount[sentence]++
			if sentenceCount[sentence] > 1 {
				return true, nil
			}

			continue
		}

		// 生成当前状态的唯一标识
		stateKey := symbolsToStringJoinSep(curr.symbols)
		if visited[stateKey] {
			continue
		}
		visited[stateKey] = true

		// 应用产生式
		for _, p := range g.Productions {
			if len(curr.symbols) > 0 && curr.symbols[0] == p.Left[0] {
				newSymbols := append(p.Right, curr.symbols[1:]...)
				queue = append(queue, state{
					symbols: newSymbols,
					depth:   curr.depth + 1,
				})
			}
		}
	}

	return false, nil
}

// === 辅助函数 ===
// isAllTerminals 判断符号序列是否全为终结符（允许 ε）
func isAllTerminals(symbols []model.Symbol, g *model.Grammar) bool {
	for _, sym := range symbols {
		if sym == model.Epsilon {
			continue
		}
		if !g.CheckIsTerminal(sym) {
			return false
		}
	}

	return true
}

// flattenSymbols 计算符号序列的实际长度（跳过 ε）
func flattenSymbols(symbols []model.Symbol) int {
	count := 0
	for _, sym := range symbols {
		if sym != model.Epsilon {
			count++
		}
	}

	return count
}
