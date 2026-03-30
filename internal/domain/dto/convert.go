package dto

import "github.com/chenzanhong/formallanglab-master/internal/domain/model"

// 文法/正则表达式转自动机的响应要带上

// 文法转FA请求
type GrammarToFARequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法转FA响应
type GrammarToFAResponse struct {
	Msg     string                    `json:"msg"`
	Result  bool                      `json:"result"`
	Process *model.GrammarToFAProcess `json:"process,omitempty"`
}

// 自动机转文法请求
type FAToGrammarRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// 自动机转文法响应
type FAToGrammarResponse struct {
	Msg     string         `json:"msg"`
	Result  bool           `json:"result"`
	Grammar *model.Grammar `json:"grammar,omitempty"`
}

// 自动机转正则表达式请求
type FAToRegexRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// 自动机转正则表达式响应
type FAToRegexResponse struct {
	Msg     string                   `json:"msg"`
	Result  bool                     `json:"result"`
	Pattern model.Regex              `json:"pattern,omitempty"`
	Process *model.ConversionProcess `json:"process,omitempty"`
}

// 正则表达式转FA请求
type RegexToFARequest struct {
	Pattern model.Regex `json:"pattern" binding:"required"`
}

// 正则表达式转FA响应
type RegexToFAResponse struct {
	Msg     string                  `json:"msg"`
	Result  bool                    `json:"result"`
	Process *model.RegexToFAProcess `json:"process,omitempty"`
}
