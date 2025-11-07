package dto

import "backend/internal/domain/model"

// FSM相关请求响应结构

// FSMValidateRequest FSM验证请求结构
type FSMValidateRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// FSMValidateResponse FSM验证响应结构
type FSMValidateResponse struct {
	Msg     string                    `json:"msg"`
	Result  bool                      `json:"result"`
	FSM     model.Automaton           `json:"fsm,omitempty"`
	FSMFlow *model.ReactFlowAutomaton `json:"fsmFlow,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

// FSMCleanupRequest FSM清理请求结构
type FSMCleanupRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// FSMCleanupResponse FSM清理响应结构
type FSMCleanupResponse struct {
	Msg     string                    `json:"msg"`
	Result  bool                      `json:"result"`
	FSM     model.Automaton           `json:"fsm,omitempty"`
	FSMFlow *model.ReactFlowAutomaton `json:"fsmFlow,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

// DFAMinimizeRequest DFA最小化请求结构
type DFAMinimizeRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// DFAMinimizeResponse DFA最小化响应结构
type DFAMinimizeResponse struct {
	Msg     string                    `json:"msg"`
	Result  bool                      `json:"result"`
	FSM     model.Automaton           `json:"fsm,omitempty"`
	FSMFlow *model.ReactFlowAutomaton `json:"fsmFlow,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

// FSMStringRecognizeRequest FSM字符串识别请求结构
type FSMStringRecognizeRequest struct {
	FSM model.Automaton `json:"fsm" binding:"required"`
	Str string          `json:"str" binding:"required"`
}

// FSMStringRecognizeResponse FSM字符串识别响应结构
type FSMStringRecognizeResponse struct {
	Msg    string                  `json:"msg"`
	Error  string                  `json:"error,omitempty"`
	Result model.RecognitionResult `json:"result,omitempty"`
}

// NFADeterminizationRequest NFA确定化请求结构
type NFADeterminizationRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// NFADeterminizationResponse NFA确定化响应结构
type NFADeterminizationResponse struct {
	Msg     string                    `json:"msg"`
	Result  bool                      `json:"result"`
	FSM     model.Automaton           `json:"fsm,omitempty"`
	FSMFlow *model.ReactFlowAutomaton `json:"fsmFlow,omitempty"`
	Error   string                    `json:"error,omitempty"`
}
