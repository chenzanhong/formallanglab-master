package dto

import "github.com/chenzanhong/formallanglab-master/internal/domain/model"

// 文法验证请求
type GrammarValidateRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法验证响应
type GrammarValidateResponse struct {
	Valid bool   `json:"valid"`
	Msg   string `json:"msg"`
	// Error  string `json:"error,omitempty"`
	Type   model.GrammarType `json:"type"`
	Result bool              `json:"result"`
}

// 文法二义性检查请求
type GrammarAmbiguityCheckRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法二义性检查响应
type GrammarAmbiguityCheckResponse struct {
	Msg         string            `json:"msg"`
	IsAmbiguous *bool             `json:"isAmbiguous,omitempty"` // 使用指针类型以支持 nil 值（无法判断）
	Type        model.GrammarType `json:"type"`
	// Error  string `json:"error,omitempty"`
	Result bool `json:"result"`
}

// 字符串识别请求
type GrammarStringRecognizeRequest struct {
	Grammar   model.Grammar `json:"grammar" binding:"required"`
	Str       string        `json:"str" binding:"required"`
	ShowSteps bool          `json:"showSteps"` // 是否返回分析步骤
	Mode      string        `json:"mode"`      // 分析模式
}

// 字符串识别响应
type GrammarStringRecognizeResponse struct {
	Accepted bool              `json:"accepted"`
	Method   string            `json:"method"`
	Msg      string            `json:"msg,omitempty"`
	Steps    []model.ParseStep `json:"steps,omitempty"`
	// Error    string            `json:"error,omitempty"`
	Result bool `json:"result"`
}

type GrammarTypeDetermineRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法类型判断响应
type GrammarTypeDetermineResponse struct {
	Type     model.GrammarType `json:"type"`
	TypeName string            `json:"typeName"`
	// Error    string `json:"error,omitempty"`
	Result bool   `json:"result"`
	Msg    string `json:"msg,omitempty"`
}

// 文法等价性检查请求
type GrammarEquivalenceCheckRequest struct {
	Grammar1 model.Grammar `json:"grammar1" binding:"required"`
	Grammar2 model.Grammar `json:"grammar2" binding:"required"`
}

// 文法等价性检查响应
type GrammarEquivalenceCheckResponse struct {
	Msg          string `json:"msg"`
	IsEquivalent bool   `json:"isEquivalent"`
	// Error        string `json:"error,omitempty"`
	Result bool `json:"result"`
}

// 文法化简请求
type GrammarSimplifyRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法化简响应
type GrammarSimplifyResponse struct {
	Msg     string         `json:"msg"`
	Grammar *model.Grammar `json:"grammar,omitempty"`
	// Error   string       `json:"error,omitempty"`
	Result bool `json:"result"`
}

// 文法 First 集请求
type GrammarFirstSetRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法 First 集响应
type GrammarFirstSetResponse struct {
	FirstSet map[string][]string `json:"firstSet"`
	Msg      string              `json:"msg"`
	// Error    string              `json:"error,omitempty"`
	Result bool `json:"result"`
}

// 文法 Follow 集请求
type GrammarFollowSetRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法 Follow 集响应
type GrammarFollowSetResponse struct {
	FollowSet map[string][]string `json:"followSet"`
	Msg       string              `json:"msg"`
	// Error     string              `json:"error,omitempty"`
	Result bool `json:"result"`
}

type GrammarGenerateExampleStringRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

type GrammarGenerateExampleStringResponse struct {
	Msg            string   `json:"msg"`
	Result         bool     `json:"result"`
	AcceptExamples []string `json:"accept"`
	RejectExamples []string `json:"reject"`
}
