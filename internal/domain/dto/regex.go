package dto

import "backend/internal/domain/model"

// 正则表达式验证请求
type RegexValidateRequest struct {
	Pattern model.Regex `json:"pattern" binding:"required"`
}

// 正则表达式验证响应
type RegexValidateResponse struct {
	Valid   bool        `json:"valid"`
	Msg     string      `json:"msg"`
	Pattern model.Regex `json:"pattern,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  bool        `json:"result"`
}

// 正则表达式识别请求
type RegexRecognizeRequest struct {
	Pattern model.Regex `json:"pattern" binding:"required"`
	Str     string      `json:"str" binding:"required"`
}

// 正则表达式识别响应
type RegexRecognizeResponse struct {
	Matched bool        `json:"matched"`
	Msg     string      `json:"msg"`
	Error   interface{} `json:"error,omitempty"`
	Result  bool        `json:"result"`
}
