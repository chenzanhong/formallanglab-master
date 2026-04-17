package grammar_s

import (
	"errors"
	"fmt"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// 判断一个 Grammar 是否为有效文法
// 验证规则：
// 1. 文法对象不能为空
// 2. 非终结符集和终结符集都不能为空
// 3. 起始符号必须存在且为非终结符
// 4. 非终结符和终结符集合不能有交集
// 5. 必须至少有一个产生式
// 6. 每个产生式的左部至少包含一个非终结符
//
// 时间复杂度：O(n*m + k)，其中 n 为产生式数量，m 为产生式平均长度，k 为符号集合大小
// 空间复杂度：O(1)
func GrammarCheckValidity(g *model.Grammar) error {
	if g == nil {
		return errors.New("文法对象为 nil")
	}

	// 1. 非终结符和终结符都不能为空
	if len(g.NonTerminals) == 0 {
		g.GrammarType = model.InvalidGrammar
		return errors.New("非终结符集合不能为空")
	}
	if len(g.Terminals) == 0 {
		g.GrammarType = model.InvalidGrammar
		return errors.New("终结符集合不能为空")
	}

	// 2. 起始符号必须是非终结符
	if g.StartSymbol == "" {
		g.GrammarType = model.InvalidGrammar
		return errors.New("起始符号不能为空")
	}
	if !g.CheckIsNonTerminal(g.StartSymbol) {
		g.GrammarType = model.InvalidGrammar
		return fmt.Errorf("起始符号'%s'必须是非终结符", g.StartSymbol)
	}

	symbolMap := make(map[model.Symbol]bool)
	for _, sym := range g.NonTerminals {
		symbolMap[sym] = true
	}
	for _, sym := range g.Terminals {
		symbolMap[sym] = true
	}
	symbolMap[model.Epsilon] = true

	// 3. 非终结符和终结符不能有交集
	for _, sym := range g.NonTerminals {
		if g.CheckIsTerminal(sym) {
			g.GrammarType = model.InvalidGrammar
			return fmt.Errorf("符号'%s'同时存在于非终结符和终结符集合中", sym)
		}
	}

	// 4. 各自集合内部无重复符号（map 天然保证 key 不重复，所以这一步由 JSON 解析保证）

	// 5. 至少有一个产生式
	if len(g.Productions) == 0 {
		g.GrammarType = model.InvalidGrammar
		return errors.New("文法必须至少包含一个产生式")
	}

	// 6. 每个产生式的左部至少包含一个非终结符，而且所使用的符号都在符号集里面
	for _, p := range g.Productions {
		if len(p.Left) == 0 {
			g.GrammarType = model.InvalidGrammar
			return fmt.Errorf("产生式 %s 的左部不能为空", p.String())
		}

		hasNonTerminal := false
		for _, sym := range p.Left {
			if g.CheckIsNonTerminal(sym) {
				hasNonTerminal = true
				break
			}
			if !symbolMap[sym] {
				g.GrammarType = model.InvalidGrammar
				return fmt.Errorf("产生式 %s 左部中的符号'%s'不在符号集合中", p.String(), sym)
			}
		}
		if !hasNonTerminal {
			g.GrammarType = model.InvalidGrammar
			return fmt.Errorf("产生式 %s 的左部必须至少包含一个非终结符", p.String())
		}
		for _, sym := range p.Right {
			if !symbolMap[sym] {
				g.GrammarType = model.InvalidGrammar
				return fmt.Errorf("产生式 %s 右部中的符号'%s'不在符号集合中", p.String(), sym)
			}
		}
	}

	return nil
}
