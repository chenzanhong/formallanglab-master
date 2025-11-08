package dto

import "backend/internal/domain/model"

// 正则表达式转NFA请求
type RegexToNFARequest struct {
	Pattern model.Regex `json:"pattern" binding:"required"`
}

// 正则表达式转NFA响应
type RegexToNFAResponse struct {
	Msg           string          `json:"msg"`
	Result        bool            `json:"result"`
	Automaton     *model.Automaton `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton  `json:"automatonFlow,omitempty"` // ToReactFlow的结果
	Error         string          `json:"error,omitempty"`
}

// 文法转NFA请求
type GrammarToNFARequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法转NFA响应
type GrammarToNFAResponse struct {
	Msg           string                    `json:"msg"`
	Result        bool                      `json:"result"`
	Automaton     *model.Automaton           `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton `json:"automatonFlow,omitempty"` // ToReactFlow的结果
	Error         string                    `json:"error,omitempty"`
}

// 自动机转文法请求
type FAToGrammarRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// 自动机转文法响应
type FAToGrammarResponse struct {
	Msg     string      `json:"msg"`
	Result  bool        `json:"result"`
	Grammar *model.Grammar `json:"grammar,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 自动机转正则表达式请求
type FAToRegexRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// 自动机转正则表达式响应
type FAToRegexResponse struct {
	Msg     string      `json:"msg"`
	Result  bool        `json:"result"`
	Pattern model.Regex `json:"pattern,omitempty"`
	Error   string      `json:"error,omitempty"`
}
