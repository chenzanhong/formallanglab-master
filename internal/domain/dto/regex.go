package dto

import "github.com/chenzanhong/formallanglab-master/internal/domain/model"

// 正则表达式验证请求
type RegexValidateRequest struct {
	Pattern model.Regex `json:"pattern" binding:"required"`
}

// 正则表达式验证响应
type RegexValidateResponse struct {
	Valid   bool        `json:"valid"`
	Msg     string      `json:"msg"`
	Pattern model.Regex `json:"pattern,omitempty"`
	// Error   interface{} `json:"error,omitempty"`
	Result bool `json:"result"`
}

// 正则表达式识别请求
type RegexRecognizeRequest struct {
	Pattern model.Regex `json:"pattern" binding:"required"`
	Str     string      `json:"str" binding:"required"`
}

// 正则表达式识别响应
type RegexRecognizeResponse struct {
	Matched bool   `json:"matched"`
	Msg     string `json:"msg"`
	Result  bool   `json:"result"`
	// Error   interface{} `json:"error,omitempty"`
}

type RegexEquivalenceCheckRequest struct {
	Pattern1 model.Regex `json:"pattern1" binding:"required"`
	Pattern2 model.Regex `json:"pattern2" binding:"required"`
}

type RegexEquivalenceCheckResponse struct {
	Msg          string `json:"msg"`
	IsEquivalent bool   `json:"isEquivalent"`
	Result       bool   `json:"result"`
}

type RegexGenerateExampleStringRequest struct {
	Pattern model.Regex `json:"pattern" binding:"required"`
}

type RegexGenerateExampleStringResponse struct {
	Msg            string   `json:"msg"`
	Result         bool     `json:"result"`
	AcceptExamples []string `json:"accept"`
	RejectExamples []string `json:"reject"`
}

// 正则表达式化简请求
type RegexSimplifyRequest struct {
	Pattern model.Regex `json:"pattern" binding:"required"`
}

// 正则表达式化简响应
type RegexSimplifyResponse struct {
	Msg             string `json:"msg"`
	Result          bool   `json:"result"`
	Original        string `json:"original"`
	Simplified      string `json:"simplified"`
	IsEmptyLanguage bool   `json:"isEmptyLanguage"`
}
