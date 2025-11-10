/*
有限状态自动机
AutomatonValidate                             // 是否有效
AutomatonCleanup                              // 去无效符号、不可达状态
DFAMinimize                             // DFA 最小化
AutomatonStringRecognize						// 字符串识别
NFADeterminization                     	// NFA 转 DFA
*/
package api

import (
	"backend/internal/domain/dto"
	"backend/internal/metrics"
	"backend/logs"
	"time"

	"backend/internal/service/automaton_s"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AutomatonValidate(c *gin.Context) { // 是否有效
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
		c.JSON(http.StatusBadRequest, dto.AutomatonValidateResponse{
			Msg:    "无效的自动机",
			Result: false,
			Error:  err.Error(),
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

func AutomatonCleanup(c *gin.Context) { // 去无效符号、无效状态以及相关的转移函数
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
			Msg:    "无效的自动机",
			Result: false,
			Error:  err.Error(),
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

func DFAMinimize(c *gin.Context) { // DFA 最小化
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
			Msg:    "无效的自动机",
			Result: false,
			Error:  err.Error(),
		})
		metrics.IncOperation("automaton", "minimize", "failure: invalid Automaton")
		return
	}
	if !req.Automaton.IsDFA {
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    "不是DFA",
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: not dfa")
		return
	}
	new_Automaton := automaton_s.DFAMinimize(&req.Automaton) // reduce
	metrics.IncOperation("automaton", "minimize", "success")
	c.JSON(http.StatusOK, dto.DFAMinimizeResponse{
		Msg:           "最小化成功",
		Result:        true,
		Automaton:     new_Automaton,
		AutomatonFlow: new_Automaton.ToReactFlow(),
	})
}

func AutomatonStringRecognize(c *gin.Context) { // 字符串识别，
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("automaton", "string_recognize", time.Since(start).Seconds())
	}()
	var req dto.AutomatonStringRecognizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logs.Sugar.Warnf("自动机字符串识别失败", "detail", "参数解析错误")
		metrics.IncOperation("automaton", "string_recognize", "failure: parameter parsing error")
		c.JSON(http.StatusBadRequest, dto.AutomatonStringRecognizeResponse{
			Msg:    "参数解析错误",
			Error:  err.Error(),
			Result: false,
		})
		return
	}
	err := automaton_s.AutomatonValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AutomatonStringRecognizeResponse{
			Msg:    "自动机验证失败",
			Error:  err.Error(),
			Result: false,
		})
		metrics.IncOperation("automaton", "minimize", "failure: invalid Automaton")
		return
	}
	result, err := automaton_s.Recognize(&req.Automaton, req.Str) // 先对Str分词，再模拟状态转移
	if err != nil || !result.IsAccepted {
		c.JSON(http.StatusOK, dto.AutomatonStringRecognizeResponse{
			Msg:               "识别成功，该字符串未被自动机接收",
			RecognitionResult: result, // 通过result.IsAccepted判断是否接收
			Error:             err.Error(),
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

func NFADeterminization(c *gin.Context) { // NFA 转 DFA，子集构造法
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
			Msg:    "无效的自动机",
			Result: false,
			Error:  err.Error(),
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "failure: invalid Automaton")
		return
	}
	if req.Automaton.IsDFA {
		c.JSON(http.StatusOK, dto.NFADeterminizationResponse{
			Msg:           "该自动机已经是DFA",
			Result:        true,
			Automaton:     &req.Automaton,
			AutomatonFlow: req.Automaton.ToReactFlow(),
		})
		metrics.IncOperation("automaton", "nfa_to_dfa", "success")
		return
	}

	new_Automaton := automaton_s.NFAToDFA(&req.Automaton)

	metrics.IncOperation("automaton", "nfa_to_dfa", "success")
	c.JSON(http.StatusOK, dto.NFADeterminizationResponse{
		Msg:           "NFA转换为DFA成功",
		Automaton:     new_Automaton,
		AutomatonFlow: new_Automaton.ToReactFlow(),
		Result:        true,
	})
}
