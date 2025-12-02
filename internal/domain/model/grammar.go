package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Symbol 表示一个符号，可以是终结符、非终结符或者自动机所识别的一个符号
type Symbol string

const (
	Epsilon Symbol = "ε" // 定义ε作为特殊输入符号，表示空转移符号
)

type GrammarType int

const (
	InvalidGrammar          GrammarType = iota - 1 // iota=0 → 0-1 = -1
	PhraseStructureGrammar                         // iota=1 → 1-1 = 0
	ContextSensitiveGrammar                        // iota=2 → 2-1 = 1
	ContextFreeGrammar                             // iota=3 → 3-1 = 2
	RegularGrammar                                 // iota=4 → 4-1 = 3
)

type GrammarLinearity string

const (
	InvalidLinearity = "非线性"
	LeftLinear       = "左线性"
	RightLinear      = "右线性"
)

var GrammarTypeNameMap = map[GrammarType]string{
	RegularGrammar:          "3型文法（正则文法）",
	ContextFreeGrammar:      "2型文法（上下文无关文法）",
	ContextSensitiveGrammar: "1型文法（上下文有关文法）",
	PhraseStructureGrammar:  "0型文法（短语结构文法）",
	InvalidGrammar:          "无效文法",
}

// 工具函数：将用户输入映射为标准 ε
func NormalizeSymbol(s string) Symbol {
	switch s {
	case "ε", "epsilon", "e", "E", "", "λ", "eps":
		return Epsilon
	default:
		return Symbol(s)
	}
}

/* 文法 */
// Production 规定了一个产生式
type Production struct {
	Left  []Symbol `json:"left"`  // 左部，通常是单个非终结符，但也可以是多个符号
	Right []Symbol `json:"right"` // 右部，可以包含多个符号
}

func (p *Production) String() string {
	var s strings.Builder
	for _, sym := range p.Left {
		s.WriteString(string(sym))
	}
	s.WriteString("->")
	for _, sym := range p.Right {
		s.WriteString(string(sym))
	}
	return s.String()
}

// Grammar 表示整个文法
type Grammar struct {
	StartSymbol      Symbol           `json:"startSymbol"`                           // 起始符号
	Terminals        []Symbol         `json:"terminals"`                             // 终结符集合
	NonTerminals     []Symbol         `json:"nonTerminals"`                          // 非终结符集合
	Productions      []Production     `json:"productions"`                           // 产生式集合
	GrammarType      GrammarType      `json:"grammarType" default:"PhraseStructure"` // 文法类型
	GrammarLinearity GrammarLinearity `json:"-"`
}

// CFGView 表示上下文无关文法/正则文法的视图结构
type CFGView struct {
	StartSymbol  Symbol                `json:"startSymbol"`  // 起始符号
	Terminals    []Symbol              `json:"terminals"`    // 终结符集合
	NonTerminals []Symbol              `json:"nonTerminals"` // 非终结符集合
	Productions  map[Symbol][][]Symbol `json:"productions"`  // 产生式集合
}

func (g *Grammar) ToCFGView() (*CFGView, error) {
	if g.GrammarType < PhraseStructureGrammar {
		return nil, errors.New("无效文法类型")
	}

	// 初始化 CFGView
	cfgView := &CFGView{
		StartSymbol:  g.StartSymbol,
		Terminals:    g.Terminals,
		NonTerminals: g.NonTerminals,
		Productions:  make(map[Symbol][][]Symbol),
	}

	// 填充产生式
	for _, prod := range g.Productions {
		if len(prod.Left) != 1 || !g.CheckIsNonTerminal(prod.Left[0]) {
			return nil, errors.New("左部必须是单个非终结符")
		}
		cfgView.Productions[prod.Left[0]] = append(cfgView.Productions[prod.Left[0]], prod.Right)
	}

	return cfgView, nil
}

// 辅助函数：根据文法中定义的符号来分割字符串
// 使用最长匹配原则，优先匹配较长的终结符
// 例如：如果文法中有 "if" 和 "i" 两个终结符，则 "if" 会被优先匹配为一个符号，而不是 "i" + "f"
func (g *Grammar) StringToSymbols(s string) ([]Symbol, error) {
	if s == "" || s == "ε" {
		return []Symbol{Epsilon}, nil
	}

	// 构建终结符列表和集合
	var terminals []string
	terminalSet := make(map[string]bool) // 可选：用于后续验证
	for _, t := range g.Terminals {
		ts := string(t)
		if ts != "" && ts != "ε" {
			terminals = append(terminals, ts)
			terminalSet[ts] = true
		}
	}

	if len(terminals) == 0 {
		return nil, fmt.Errorf("grammar has no defined terminals")
	}

	// 按长度降序排序（最长优先）
	sort.Slice(terminals, func(i, j int) bool {
		return len(terminals[i]) > len(terminals[j])
	})

	var symbols []Symbol
	i := 0
	for i < len(s) {
		matched := false
		for _, term := range terminals {
			if i+len(term) <= len(s) && s[i:i+len(term)] == term {
				symbols = append(symbols, Symbol(term))
				i += len(term)
				matched = true
				break
			}
		}
		if !matched {
			// 返回已切分部分
			return symbols, fmt.Errorf("no terminal matches substring starting at position %d in input %q", i, s)
		}
	}
	return symbols, nil
}

func (g *Grammar) CheckIsTerminal(s Symbol) bool {
	for _, v := range g.Terminals {
		if v == s {
			return true
		}
	}
	return false
}

func (g *Grammar) CheckIsNonTerminal(s Symbol) bool {
	for _, v := range g.NonTerminals {
		if v == s {
			return true
		}
	}
	return false
}

// ParseStep 表示文法分析过程中的一个步骤
type ParseStep struct {
	StepType    string      `json:"stepType"`    // "init", "match", "predict", "accept", "error"
	Description string      `json:"description"` // 步骤描述
	Stack       []Symbol    `json:"stack"`       // 当前栈状态
	Input       []Symbol    `json:"input"`       // 剩余输入
	InputPos    int         `json:"inputPos"`    // 输入指针位置
	Action      string      `json:"action"`      // 执行的动作描述
	Production  *Production `json:"production"`  // 使用的产生式（如果有）
}

// ParseResult 表示文法分析的完整结果
type ParseResult struct {
	Accepted bool        `json:"accepted"` // 是否接受输入串
	Method   string      `json:"method"`   // 使用的分析方法
	Steps    []ParseStep `json:"steps"`    // 分析步骤（可选）
	Error    string      `json:"error"`    // 错误信息（如果有）
	Message  string      `json:"message"`  // 额外信息
}

// ToReactFlow 返回格式化后的产生式字符串，用于前端展示
func (g *Grammar) ToReactFlow() []string {
	productions := make([]string, 0, len(g.Productions))

	for _, prod := range g.Productions {
		// 格式化左部
		left := string(prod.Left[0]) // 通常左部只有一个符号

		// 格式化右部
		right := ""
		if len(prod.Right) == 0 || (len(prod.Right) == 1 && prod.Right[0] == Epsilon) {
			right = string(Epsilon)
		} else {
			for i, symbol := range prod.Right {
				if i > 0 {
					right += " "
				}
				if symbol == Epsilon {
					right += string(Epsilon)
				} else {
					right += string(symbol)
				}
			}
		}

		// 组合成产生式字符串
		productions = append(productions, left+" -> "+right)
	}

	return productions
}

/*

前端传给后端的文法格式：：
{
  "startSymbol": "S",
  "terminals": ["a", "b", "ε"],
  "nonTerminals": ["S", "A", "B"],
  "productions": [
    {
      "left": "S",
      "right": ["a", "A", "b"]
    },
    {
      "left": "A",
      "right": ["b", "B"]
    },
    {
      "left": "B",
      "right": ["ε"]
    }
  ]
}
*/
