package dto

import "backend/internal/domain/model"

// Automaton相关请求响应结构
// 不使用 Error 字段，具体错误在 Msg 字段说明

// AutomatonValidateRequest Automaton验证请求结构
type AutomatonValidateRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// AutomatonValidateResponse Automaton验证响应结构
type AutomatonValidateResponse struct {
	Msg           string                    `json:"msg"`
	Result        bool                      `json:"result"`
	Automaton     *model.Automaton          `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton `json:"automatonFlow,omitempty"`
	// Error         string                    `json:"error,omitempty"`
}

// AutomatonCleanupRequest Automaton清理请求结构
type AutomatonCleanupRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// AutomatonCleanupResponse Automaton清理响应结构
type AutomatonCleanupResponse struct {
	Msg           string                    `json:"msg"`
	Result        bool                      `json:"result"`
	Automaton     *model.Automaton          `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton `json:"automatonFlow,omitempty"`
	// Error         string                    `json:"error,omitempty"`
}

// DFAMinimizeRequest DFA最小化请求结构
type DFAMinimizeRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// DFAMinimizeResponse DFA最小化响应结构
type DFAMinimizeResponse struct {
	Msg           string                    `json:"msg"`
	Result        bool                      `json:"result"`
	Automaton     *model.Automaton          `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton `json:"automatonFlow,omitempty"`
	// Error         string                    `json:"error,omitempty"`
}

// AutomatonStringRecognizeRequest Automaton字符串识别请求结构
type AutomatonStringRecognizeRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
	Str       string          `json:"str" binding:"required"`
}

// AutomatonStringRecognizeResponse Automaton字符串识别响应结构
type AutomatonStringRecognizeResponse struct {
	Msg string `json:"msg"`
	// Error           string                  `json:"error,omitempty"`
	RecognitionResult *model.RecognitionResult `json:"recognitionResult,omitempty"`
	Result            bool                     `json:"result"`
}

// NFADeterminizationRequest NFA确定化请求结构
type NFADeterminizationRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// NFADeterminizationResponse NFA确定化响应结构
type NFADeterminizationResponse struct {
	Msg           string                    `json:"msg"`
	Result        bool                      `json:"result"`
	Automaton     *model.Automaton          `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton `json:"automatonFlow,omitempty"`
	// Error         string                    `json:"error,omitempty"`
}
