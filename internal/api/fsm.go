/*
有限状态自动机
FSMValidate                             // 是否有效
FSMCleanup                              // 去无效符号、不可达状态
DFAMinimize                             // DFA 最小化
FSMStringRecognize						// 字符串识别
NFADeterminization                     	// NFA 转 DFA
*/
package api

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	"backend/internal/metrics"
	"fmt"
	"time"

	"backend/internal/service/fsm_s"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FSMValidate(c *gin.Context) { // 是否有效
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "validate", time.Since(start).Seconds())
	}()
	var req dto.FSMValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FSMValidateResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("fsm", "validate", "failure: parameter parsing error")
		return
	}

	ok, err := fsm_s.FSMValidate(&req.Automaton)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.FSMValidateResponse{
			Msg:    "无效的自动机",
			Result: ok,
			Error:  err.Error(),
		})
		metrics.IncOperation("fsm", "validate", "failure: invalid fsm")
		return
	}

	metrics.IncOperation("fsm", "validate", "success")
	c.JSON(http.StatusOK, dto.FSMValidateResponse{
		Msg:     "有效",
		Result:  ok,
		FSM:     req.Automaton,
		FSMFlow: req.Automaton.ToReactFlow(),
	})
}

func FSMCleanup(c *gin.Context) { // 去无效符号、无效状态以及相关的转移函数
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "cleanup", time.Since(start).Seconds())
	}()
	var req dto.FSMCleanupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FSMCleanupResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("fsm", "cleanup", "failure: parameter parsing error")
		return
	}
	fmt.Printf("%+v\n", req.Automaton)
	ok, err := fsm_s.FSMValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, dto.FSMCleanupResponse{
			Msg:    "无效的自动机",
			Result: ok,
			Error:  err.Error(),
		})
		metrics.IncOperation("fsm", "minimize", "failure: invalid fsm")
		return
	}

	fsm_s.Cleanup(&req.Automaton)
	metrics.IncOperation("fsm", "cleanup", "success")
	c.JSON(http.StatusOK, dto.FSMCleanupResponse{
		Msg:     "简化失败",
		Result:  true,
		FSM:     req.Automaton,
		FSMFlow: req.Automaton.ToReactFlow(),
	})
}

func DFAMinimize(c *gin.Context) { // DFA 最小化
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "minimize", time.Since(start).Seconds())
	}()
	var req dto.DFAMinimizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("fsm", "minimize", "failure: parameter parsing error")
		return
	}
	fmt.Printf("要最小化的fsm：%+v", req.Automaton)
	ok, err := fsm_s.FSMValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    "无效的自动机",
			Result: ok,
			Error:  err.Error(),
		})
		metrics.IncOperation("fsm", "minimize", "failure: invalid fsm")
		return
	}
	if !req.Automaton.IsDFA {
		c.JSON(http.StatusBadRequest, dto.DFAMinimizeResponse{
			Msg:    "不是DFA",
			Result: false,
		})
		metrics.IncOperation("fsm", "minimize", "failure: not dfa")
		return
	}
	new_fsm := fsm_s.DFAMinimize(&req.Automaton) // reduce
	metrics.IncOperation("fsm", "minimize", "success")
	c.JSON(http.StatusOK, dto.DFAMinimizeResponse{
		Msg:     "最小化成功",
		Result:  true,
		FSM:     *new_fsm,
		FSMFlow: new_fsm.ToReactFlow(),
	})
}

func FSMStringRecognize(c *gin.Context) { // 字符串识别，
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "string_recognize", time.Since(start).Seconds())
	}()
	var req dto.FSMStringRecognizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FSMStringRecognizeResponse{
			Msg:    "参数解析错误",
			Result: model.RecognitionResult{},
		})
		metrics.IncOperation("fsm", "string_recognize", "failure: parameter parsing error")
		return
	}
	ok, err := fsm_s.FSMValidate(&req.FSM) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, dto.FSMStringRecognizeResponse{
			Msg:    "无效的自动机",
			Result: model.RecognitionResult{},
			Error:  err.Error(),
		})
		metrics.IncOperation("fsm", "minimize", "failure: invalid fsm")
		return
	}
	result, err := fsm_s.Recognize(&req.FSM, req.Str) // 先对Str分词，再模拟状态转移
	if result == nil || !result.IsAccepted {
		c.JSON(http.StatusBadRequest, dto.FSMStringRecognizeResponse{
			Msg:    "识别失败",
			Result: *result,
			Error:  err.Error(),
		})
		metrics.IncOperation("fsm", "string_recognize", "failure: recognition failed")
		return
	}
	metrics.IncOperation("fsm", "string_recognize", "success")
	c.JSON(http.StatusOK, dto.FSMStringRecognizeResponse{
		Msg:    "识别成功",
		Result: *result,
	})
}

func NFADeterminization(c *gin.Context) { // NFA 转 DFA，子集构造法
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "nfa_to_dfa", time.Since(start).Seconds())
	}()
	var req dto.NFADeterminizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NFADeterminizationResponse{
			Msg:    "参数解析错误",
			Result: false,
		})
		metrics.IncOperation("fsm", "nfa_to_dfa", "failure: parameter parsing error")
		return
	}
	ok, err := fsm_s.FSMValidate(&req.Automaton) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, dto.NFADeterminizationResponse{
			Msg:    "无效的自动机",
			Result: false,
			Error:  err.Error(),
		})
		metrics.IncOperation("fsm", "nfa_to_dfa", "failure: invalid fsm")
		return
	}
	if req.Automaton.IsDFA {
		c.JSON(http.StatusBadRequest, dto.NFADeterminizationResponse{
			Msg:    "该自动机已经是DFA",
			Result: false,
		})
		metrics.IncOperation("fsm", "nfa_to_dfa", "failure: not nfa")
		return
	}
	new_fsm := fsm_s.NFAToDFA(&req.Automaton)
	metrics.IncOperation("fsm", "nfa_to_dfa", "success")
	c.JSON(http.StatusOK, dto.NFADeterminizationResponse{
		Msg:     "NFA转换为DFA成功",
		FSM:     *new_fsm,
		FSMFlow: new_fsm.ToReactFlow(),
		Result:  true,
	})
}
