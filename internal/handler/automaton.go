/*
有限状态自动机
AutomatonValidate                  // 是否有效
AutomatonCleanup                   // 去无效符号、不可达状态
DFAMinimize                        // DFA 最小化
AutomatonStringRecognize           // 字符串识别
NFADeterminization                 // NFA 转 DFA
AutomatonEquivalenceCheck          // 判断两个有限自动机是否等价
AutomatonGenerateExampleString     // 生成字符串示例，含可识别和不可识别的
*/
package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/dto"
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/metrics"
	"github.com/chenzanhong/formallanglab-master/internal/service/automaton_s"
	"github.com/chenzanhong/zlog"
	"github.com/gin-gonic/gin"
)

func AutomatonValidate(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "validate", time.Since(start).Seconds())
	}()

	var req dto.AutomatonValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AutomatonValidateResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("automaton", "validate", "failure: parameter parsing error")

		return
	}

	if err := automaton_s.AutomatonValidate(&req.Automaton); err != nil {
		zlog.Warnw("Automaton validation failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.AutomatonValidateResponse{
			Msg:    fmt.Sprintf("无效的自动机：%v", err),
			Result: false,
		})
		metrics.IncOperation("automaton", "validate", "failure: invalid Automaton")

		return
	}

	metrics.IncOperation("automaton", "validate", "success")
	c.JSON(http.StatusOK, dto.AutomatonValidateResponse{
		Msg:           "有效",
		Result:        true,
		Automaton:     &req.Automaton,
		AutomatonFlow: req.Automaton.ToReactFlow(),
	})
}

func AutomatonCleanup(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "cleanup", time.Since(start).Seconds())
	}()

	var req dto.AutomatonCleanupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AutomatonCleanupResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("automaton", "cleanup", "failure: parameter parsing error")

		return
	}

	err := automaton_s.AutomatonValidate(&req.Automaton)
	if err != nil {
		zlog.Warnw("Automaton cleanup failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.AutomatonCleanupResponse{
			Msg:    fmt.Sprintf("无效的自动机：%v", err),
			Result: false,
		})
		metrics.IncOperation("automaton", "cleanup", "failure: invalid Automaton")

		return
	}

	automaton_s.Cleanup(&req.Automaton)
	metrics.IncOperation("automaton", "cleanup", "success")
	c.JSON(http.StatusOK, dto.AutomatonCleanupResponse{
		Msg:           "自动机清理成功",
		Result:        true,
		Automaton:     &req.Automaton,
		AutomatonFlow: req.Automaton.ToReactFlow(),
	})
}

func DFAMinimize(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "minimize", time.Since(start).Seconds())
	}()

	var req dto.DFAMinimizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: parameter parsing error")

		return
	}

	err := automaton_s.AutomatonValidate(&req.Automaton)
	if err != nil {
		zlog.Warnw("DFA minimization failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    fmt.Sprintf("无效的自动机：%v", err),
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: invalid Automaton")

		return
	}

	if req.Automaton.Type != model.DFA {
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    "不是 DFA",
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: not dfa")

		return
	}

	newAutomaton, process := automaton_s.DFAMinimizeWithProcess(&req.Automaton)
	metrics.IncOperation("automaton", "minimize", "success")
	c.JSON(http.StatusOK, dto.DFAMinimizeResponse{
		Msg:           "最小化成功",
		Result:        true,
		Automaton:     newAutomaton,
		AutomatonFlow: newAutomaton.ToReactFlow(),
		Process:       process,
	})
}

func AutomatonStringRecognize(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "string_recognize", time.Since(start).Seconds())
	}()

	var req dto.AutomatonStringRecognizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Warnw("自动机字符串识别失败", "detail", "参数解析错误")
		metrics.IncOperation("automaton", "string_recognize", "failure: parameter parsing error")
		c.JSON(http.StatusBadRequest, dto.AutomatonStringRecognizeResponse{
			Msg:    "参数解析错误" + err.Error(),
			Result: false,
		})

		return
	}

	err := automaton_s.AutomatonValidate(&req.Automaton)
	if err != nil {
		zlog.Warnw("Automaton string recognition failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.AutomatonStringRecognizeResponse{
			Msg:    fmt.Sprintf("自动机验证失败：%v", err),
			Result: false,
		})
		metrics.IncOperation("automaton", "string_recognize", "failure: invalid Automaton")

		return
	}

	result, err := automaton_s.Recognize(&req.Automaton, req.Str)
	if err != nil || !result.IsAccepted {
		c.JSON(http.StatusOK, dto.AutomatonStringRecognizeResponse{
			Msg:               "识别成功，该字符串未被自动机接收：" + err.Error(),
			RecognitionResult: result,
			Result:            true,
		})
		metrics.IncOperation("automaton", "string_recognize", "failure: recognition failed")

		return
	}

	metrics.IncOperation("automaton", "string_recognize", "success")
	c.JSON(http.StatusOK, dto.AutomatonStringRecognizeResponse{
		Msg:               "识别成功",
		RecognitionResult: result,
		Result:            true,
	})
}

func NFADeterminization(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "nfa_to_dfa", time.Since(start).Seconds())
	}()

	var req dto.NFADeterminizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NFADeterminizationResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "failure: parameter parsing error")

		return
	}

	err := automaton_s.AutomatonValidate(&req.Automaton)
	if err != nil {
		zlog.Warnw("NFA determinization failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.NFADeterminizationResponse{
			Msg:    fmt.Sprintf("无效的自动机：%v", err),
			Result: false,
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "failure: invalid Automaton")

		return
	}

	if req.Automaton.Type == model.DFA {
		c.JSON(http.StatusOK, dto.NFADeterminizationResponse{
			Msg:           "该自动机已经是 DFA",
			Result:        true,
			Automaton:     &req.Automaton,
			AutomatonFlow: req.Automaton.ToReactFlow(),
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "success")

		return
	}

	process := automaton_s.NFAToDFAWithProcess(&req.Automaton)
	if process == nil {
		c.JSON(http.StatusBadRequest, dto.NFADeterminizationResponse{
			Msg:    "NFA 转换为 DFA 失败",
			Result: false,
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "failure: nfa_to_dfa failed")

		return
	}

	metrics.IncOperation("automaton", "nfa_to_dfa", "success")
	c.JSON(http.StatusOK, dto.NFADeterminizationResponse{
		Msg:           "NFA 转换为 DFA 成功",
		Automaton:     process.FinalAutomaton,
		AutomatonFlow: process.FinalAutomaton.ToReactFlow(),
		Result:        true,
		Process:       process,
	})
}

func AutomatonEquivalenceCheck(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "equivalence_check", time.Since(start).Seconds())
	}()

	var req dto.AutomatonEquivalenceCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("automaton", "equivalence_check", "failure: parameter parsing error")
		c.JSON(http.StatusBadRequest, dto.AutomatonEquivalenceCheckResponse{
			Msg:    "参数解析错误",
			Result: false,
		})

		return
	}

	err := automaton_s.AutomatonValidate(&req.Automaton1)
	if err != nil {
		zlog.Warnw("Automaton equivalence check failed", "automaton", "1", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.AutomatonEquivalenceCheckResponse{
			Msg:    fmt.Sprintf("无效的自动机 1：%v", err),
			Result: false,
		})

		return
	}

	err = automaton_s.AutomatonValidate(&req.Automaton2)
	if err != nil {
		zlog.Warnw("Automaton equivalence check failed", "automaton", "2", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.AutomatonEquivalenceCheckResponse{
			Msg:    fmt.Sprintf("无效的自动机 2：%v", err),
			Result: false,
		})

		return
	}

	dfa1, dfa2, isEquivalent := automaton_s.AutomatonEquivalenceCheck(&req.Automaton1, &req.Automaton2)
	if isEquivalent {
		metrics.IncOperation("automaton", "equivalence_check", "success")
	} else {
		metrics.IncOperation("automaton", "equivalence_check", "failure: not equivalent")
	}
	c.JSON(http.StatusOK, dto.AutomatonEquivalenceCheckResponse{
		Msg:           "自动机等价检查完成",
		Result:        true,
		IsEquivalent:  isEquivalent,
		MinimizedDFA1: dfa1,
		MinimizedDFA2: dfa2,
	})
}

func AutomatonGenerateExampleString(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "generate_example", time.Since(start).Seconds())
	}()

	var req dto.AutomatonGenerateExampleStringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AutomatonGenerateExampleStringResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("automaton", "generate_example", "failure: parameter parsing error")

		return
	}

	err := automaton_s.AutomatonValidate(&req.Automaton)
	if err != nil {
		zlog.Warnw("Automaton example generation failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.AutomatonGenerateExampleStringResponse{
			Msg:    fmt.Sprintf("无效的自动机：%v", err),
			Result: false,
		})
		metrics.IncOperation("automaton", "generate_example", "failure: invalid Automaton")

		return
	}

	accept, reject := automaton_s.AutomatonGenerateExampleString(&req.Automaton)
	metrics.IncOperation("automaton", "generate_example", "success")
	c.JSON(http.StatusOK, dto.AutomatonGenerateExampleStringResponse{
		Msg:            "示例字符串生成成功",
		Result:         true,
		AcceptExamples: accept,
		RejectExamples: reject,
	})
}
