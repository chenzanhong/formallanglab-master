package grammar_s

import (
    "backend/internal/domain/model"
    "strings"
)

// symbolsToString 将符号切片转为字符串，用于 map 的 key
// 注意：ε 用特殊符号表示，如 "<eps>"
func symbolsToString(symbols []model.Symbol) string {
    parts := make([]string, len(symbols))
    for i, sym := range symbols {
        if sym == model.Epsilon {
            parts[i] = "<eps>"
        } else {
            parts[i] = string(sym)
        }
    }
    return strings.Join(parts, ",")
}

// IsAmbiguousRegular 检查一个正则文法是否二义
// 假设输入文法已经是 Type3（正则文法）
func IsAmbiguousRegular(g *model.Grammar) (bool, error) {
    // 设置搜索深度上限：2 * 非终结符数量
    maxDepth := 2 * len(g.NonTerminals)
    if maxDepth > 10 { // 防止爆炸，设个软上限
        maxDepth = 10
    }

    // 记录每个句子被推导出的次数
    sentenceCount := make(map[string]int)

    // 使用 BFS 生成所有可能推导
    type state struct {
        symbols []model.Symbol
    }

    queue := []state{{symbols: []model.Symbol{g.StartSymbol}}}
    visited := make(map[string]bool) // 防止重复状态

    for len(queue) > 0 {
        curr := queue[0]
        queue = queue[1:]

        key := symbolsToString(curr.symbols)
        if visited[key] {
            continue
        }
        visited[key] = true

        // 如果全是终结符，得到一个句子
        if isAllTerminals(curr.symbols, g) {
            sentence := symbolsToString(curr.symbols)
            sentenceCount[sentence]++
            if sentenceCount[sentence] > 1 {
                return true, nil
            }
            continue
        }

        // 尝试应用每个产生式
        for _, p := range g.Productions {
            if len(curr.symbols) > 0 && curr.symbols[0] == p.Left[0] {
                newSymbols := append(p.Right, curr.symbols[1:]...)
                // 直接使用 flattenSymbols 返回的 int
                if flattenSymbols(newSymbols) <= maxDepth {
                    queue = append(queue, state{symbols: newSymbols})
                }
            }
        }
    }

    // 再次检查（理论上前面已检查，但保险）
    for _, count := range sentenceCount {
        if count > 1 {
            return true, nil
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
        if !isInSet(g.Terminals, sym) {
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