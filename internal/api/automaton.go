/*
有限状态自动机
AutomatonValidate                             // 是否有效
AutomatonCleanup                              // 去无效符号、不可达状态
DFAMinimize                             // DFA 最小化
AutomatonStringRecognize						// 字符串识别
NFADeterminization                     	// NFA 转 DFA
AutomatonEquivalenceCheck				// 判断两个有限自动机是否等价
AutomatonGenerateExampleString					// 生成字符串示例，含可识别和不可识别的
*/
package api

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	"backend/internal/metrics"
	"fmt"
	"time"

	"backend/internal/service/automaton_s"
	"net/http"

	"github.com/chenzanhong/zlog"
	"github.com/gin-gonic/gin"
)

// 是否有效
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
	fmt.Printf("%+v", req.Automaton)
	if err := automaton_s.AutomatonValidate(&req.Automaton); err != nil {
		c.JSON(http.StatusBadRequest, dto.AutomatonValidateResponse{
			Msg:    "无效的自动机" + err.Error(),
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

// 去无效符号、无效状态以及相关的转移函数
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
	// fmt.Printf("%+v\n", req.Automaton)
	err := automaton_s.AutomatonValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AutomatonCleanupResponse{
			Msg:    "无效的自动机" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: invalid Automaton")
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

// DFA 最小化
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
	// fmt.Printf("要最小化的Automaton：%+v", req.Automaton)
	err := automaton_s.AutomatonValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    "无效的自动机" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: invalid Automaton")
		return
	}
	if req.Automaton.Type != model.DFA {
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    "不是DFA",
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: not dfa")
		return
	}
	// 使用带有过程记录的最小化方法
	new_Automaton, process := automaton_s.DFAMinimizeWithProcess(&req.Automaton) // reduce
	metrics.IncOperation("automaton", "minimize", "success")
	c.JSON(http.StatusOK, dto.DFAMinimizeResponse{
		Msg:           "最小化成功",
		Result:        true,
		Automaton:     new_Automaton,
		AutomatonFlow: new_Automaton.ToReactFlow(),
		Process:       process,
	})
}

// 字符串识别，
func AutomatonStringRecognize(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "string_recognize", time.Since(start).Seconds())
	}()
	var req dto.AutomatonStringRecognizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Warnf("自动机字符串识别失败", "detail", "参数解析错误")
		metrics.IncOperation("automaton", "string_recognize", "failure: parameter parsing error")
		c.JSON(http.StatusBadRequest, dto.AutomatonStringRecognizeResponse{
			Msg:    "参数解析错误" + err.Error(),
			Result: false,
		})
		return
	}
	err := automaton_s.AutomatonValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AutomatonStringRecognizeResponse{
			Msg:    "自动机验证失败" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: invalid Automaton")
		return
	}
	result, err := automaton_s.Recognize(&req.Automaton, req.Str) // 先对Str分词，再模拟状态转移
	if err != nil || !result.IsAccepted {
		c.JSON(http.StatusOK, dto.AutomatonStringRecognizeResponse{
			Msg:               "识别成功，该字符串未被自动机接收：" + err.Error(),
			RecognitionResult: result, // 通过result.IsAccepted判断是否接收
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

// NFA 转 DFA，子集构造法
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
	err := automaton_s.AutomatonValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NFADeterminizationResponse{
			Msg:    "无效的自动机：" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "failure: invalid Automaton")
		return
	}
	if req.Automaton.Type == model.DFA {
		c.JSON(http.StatusOK, dto.NFADeterminizationResponse{
			Msg:           "该自动机已经是DFA",
			Result:        true,
			Automaton:     &req.Automaton,
			AutomatonFlow: req.Automaton.ToReactFlow(),
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "success")
		return
	}

	// new_Automaton := automaton_s.NFAToDFA(&req.Automaton)
	// 使用带有过程记录的最小化方法
	process := automaton_s.NFAToDFAWithProcess(&req.Automaton) // reduce
	if process == nil {
		c.JSON(http.StatusBadRequest, dto.NFADeterminizationResponse{
			Msg:    "NFA转换为DFA失败",
			Result: false,
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "failure: nfa_to_dfa failed")
		return
	}

	metrics.IncOperation("automaton", "nfa_to_dfa", "success")
	c.JSON(http.StatusOK, dto.NFADeterminizationResponse{
		Msg:           "NFA转换为DFA成功",
		Automaton:     process.FinalAutomaton,
		AutomatonFlow: process.FinalAutomaton.ToReactFlow(),
		Result:        true,
		Process:       process,
	})
}

// 判断两个有限自动机是否等价
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
	err := automaton_s.AutomatonValidate(&req.Automaton1) // 包含了DFA还是NFA的判断
	if err != nil {
		metrics.IncOperation("automaton", "equivalence_check", "failure: invalid Automaton1")
		c.JSON(http.StatusBadRequest, dto.AutomatonEquivalenceCheckResponse{
			Msg:    "无效的自动机1：" + err.Error(),
			Result: false,
		})
		return
	}
	err = automaton_s.AutomatonValidate(&req.Automaton2) // 包含了DFA还是NFA的判断
	if err != nil {
		metrics.IncOperation("automaton", "equivalence_check", "failure: invalid Automaton2")
		c.JSON(http.StatusBadRequest, dto.AutomatonEquivalenceCheckResponse{
			Msg:    "无效的自动机2：" + err.Error(),
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

// 生成自动机字符串示例
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
	err := automaton_s.AutomatonValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AutomatonGenerateExampleStringResponse{
			Msg:    "无效的自动机：" + err.Error(),
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
