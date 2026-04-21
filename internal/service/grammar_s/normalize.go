package grammar_s

import (
	"strings"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func NormalizeGrammar(grammar *model.Grammar) {
	// 1. ε 符号规范化
	for _, prod := range grammar.Productions {
		for i, sym := range prod.Right {
			if sym == "ε" || sym == "λ" || sym == "epsilon" {
				prod.Right[i] = model.Epsilon
			}
		}
	}

	// 2. 去除重复产生式
	grammar.Productions = deduplicateProductions(grammar.Productions)
}

// deduplicateProductions 去除重复的产生式
// 两个产生式相同当且仅当：左部相同且右部符号序列完全相同
func deduplicateProductions(prods []model.Production) []model.Production {
	if len(prods) == 0 {
		return prods
	}

	// 使用 map 记录已出现的产生式
	seen := make(map[string]bool)
	result := make([]model.Production, 0, len(prods))

	for _, prod := range prods {
		// 生成产生式的唯一标识：左部 -> 右部符号序列
		key := productionToKey(prod)
		if !seen[key] {
			seen[key] = true
			result = append(result, prod)
		}
	}

	return result
}

// productionToKey 将产生式转换为字符串 key，用于去重比较
func productionToKey(prod model.Production) string {
	// 格式："Left->Right1,Right2,..."
	var key strings.Builder

	// 左部
	for _, sym := range prod.Left {
		key.WriteString(string(sym))
	}
	key.WriteString("->")

	// 右部
	for i, sym := range prod.Right {
		if i > 0 {
			key.WriteString(",")
		}
		key.WriteString(string(sym))
	}

	return key.String()
}
