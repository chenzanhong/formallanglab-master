package dto

import "github.com/chenzanhong/formallanglab-master/internal/domain/model"

// Automaton 相关请求响应结构
// 不使用 Error 字段，具体错误在 Msg 字段说明

// AutomatonValidateRequest Automaton 验证请求结构
type AutomatonValidateRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// AutomatonValidateResponse Automaton 验证响应结构
type AutomatonValidateResponse struct {
	Msg           string                    `json:"msg"`
	Result        bool                      `json:"result"`
	Automaton     *model.Automaton          `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton `json:"automatonFlow,omitempty"`
	// Error         string                    `json:"error,omitempty"`
}

// AutomatonCleanupRequest Automaton 清理请求结构
type AutomatonCleanupRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// AutomatonCleanupResponse Automaton 清理响应结构
type AutomatonCleanupResponse struct {
	Msg           string                    `json:"msg"`
	Result        bool                      `json:"result"`
	Automaton     *model.Automaton          `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton `json:"automatonFlow,omitempty"`
	// Error         string                    `json:"error,omitempty"`
}

// DFAMinimizeRequest DFA 最小化请求结构
type DFAMinimizeRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// DFAMinimizeResponse DFA 最小化响应结构
type DFAMinimizeResponse struct {
	Msg           string                     `json:"msg"`
	Result        bool                       `json:"result"`
	Automaton     *model.Automaton           `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton  `json:"automatonFlow,omitempty"`
	Process       *model.MinimizationProcess `json:"process,omitempty"`
	// Error         string                    `json:"error,omitempty"`
}

// AutomatonStringRecognizeRequest Automaton 字符串识别请求结构
type AutomatonStringRecognizeRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
	Str       string          `json:"str" binding:"required"`
}

// AutomatonStringRecognizeResponse Automaton 字符串识别响应结构
type AutomatonStringRecognizeResponse struct {
	Msg string `json:"msg"`
	// Error           string                  `json:"error,omitempty"`
	RecognitionResult *model.RecognitionResult `json:"recognitionResult,omitempty"`
	Result            bool                     `json:"result"`
}

// NFADeterminizationRequest NFA 确定化请求结构
type NFADeterminizationRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

// NFADeterminizationResponse NFA 确定化响应结构
type NFADeterminizationResponse struct {
	Msg           string                           `json:"msg"`
	Result        bool                             `json:"result"`
	Automaton     *model.Automaton                 `json:"automaton,omitempty"`
	AutomatonFlow *model.ReactFlowAutomaton        `json:"automatonFlow,omitempty"`
	Process       *model.NFADeterminizationProcess `json:"process,omitempty"`
}

// AutomatonEquivalenceCheckRequest 自动机等价检查请求结构
type AutomatonEquivalenceCheckRequest struct {
	Automaton1 model.Automaton `json:"automaton1" binding:"required"`
	Automaton2 model.Automaton `json:"automaton2" binding:"required"`
}

type AutomatonEquivalenceCheckResponse struct {
	Msg           string           `json:"msg"`
	Result        bool             `json:"result"`
	IsEquivalent  bool             `json:"isEquivalent"`
	MinimizedDFA1 *model.Automaton `json:"minimizedDFA1,omitempty"`
	MinimizedDFA2 *model.Automaton `json:"minimizedDFA2,omitempty"`
}

// AutomatonGenerateExampleStringRequest 自动机生成示例字符串请求结构
type AutomatonGenerateExampleStringRequest struct {
	Automaton model.Automaton `json:"automaton" binding:"required"`
}

type AutomatonGenerateExampleStringResponse struct {
	Msg            string   `json:"msg"`
	Result         bool     `json:"result"`
	AcceptExamples []string `json:"accept,omitempty"`
	RejectExamples []string `json:"reject,omitempty"`
}
