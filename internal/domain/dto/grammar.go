package dto

import "backend/internal/domain/model"

// 文法验证请求
type GrammarValidateRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法验证响应
type GrammarValidateResponse struct {
	Valid bool   `json:"valid"`
	Msg   string `json:"msg"`
	Error string `json:"error,omitempty"`
	Type  int    `json:"type"`
}

// 文法二义性检查请求
type GrammarAmbiguityCheckRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法二义性检查响应
type GrammarAmbiguityCheckResponse struct {
	Msg         string `json:"msg"`
	IsAmbiguous *bool  `json:"isAmbiguous,omitempty"` // 使用指针类型以支持nil值（无法判断）
	// IsRegular   bool   `json:"isRegular"`
	Type  int    `json:"type"`
	Error string `json:"error,omitempty"`
}

// 字符串识别请求
type GrammarStringRecognizeRequest struct {
	Grammar   model.Grammar `json:"grammar" binding:"required"`
	Input     string        `json:"input" binding:"required"`
	ShowSteps bool          `json:"showSteps"` // 是否返回分析步骤
	Mode      string        `json:"mode"`      // 分析模式
}

// 字符串识别响应
type GrammarStringRecognizeResponse struct {
	Accepted bool              `json:"accepted"`
	Method   string            `json:"method"`
	Message  string            `json:"message,omitempty"`
	Steps    []model.ParseStep `json:"steps,omitempty"`
	Error    string            `json:"error,omitempty"`
}

type GrammarTypeDetermineRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法类型判断响应
type GrammarTypeDetermineResponse struct {
	Type     int    `json:"type"`
	TypeName string `json:"typeName"`
	Error    string `json:"error,omitempty"`
}

// 文法等价性检查请求
type GrammarEquivalenceCheckRequest struct {
	G1 model.Grammar `json:"g1" binding:"required"`
	G2 model.Grammar `json:"g2" binding:"required"`
}

// 文法等价性检查响应
type GrammarEquivalenceCheckResponse struct {
	Msg          string `json:"msg"`
	IsEquivalent bool   `json:"isEquivalent"`
	Error        string `json:"error,omitempty"`
}

// 文法化简请求
type GrammarSimplifyRequest struct {
	Grammar model.Grammar `json:"grammar" binding:"required"`
}

// 文法化简响应
type GrammarSimplifyResponse struct {
	Msg     string      `json:"msg"`
	Grammar model.Grammar `json:"grammar,omitempty"`
	// Type    int         `json:"type,omitempty"`
	Error   string      `json:"error,omitempty"`
}
