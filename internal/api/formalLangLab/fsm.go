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
	var fsm model.Automaton
	if err := c.ShouldBindJSON(&fsm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "result": false})
		metrics.IncOperation("fsm", "validate", "failure: parameter parsing error")
		return
	}

	ok, err := fsm.ISValidate()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的自动机", "result": ok, "error": err.Error()})
		metrics.IncOperation("fsm", "validate", "failure: invalid fsm")
		return
	}

	metrics.IncOperation("fsm", "validate", "success")
	c.JSON(http.StatusOK, gin.H{"msg": "有效", "result": ok, "fsm": fsm.ToReactFlow()})
}

func FSMCleanup(c *gin.Context) { // 去无效符号、无效状态以及相关的转移函数
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "cleanup", time.Since(start).Seconds())
	}()
	var fsm model.Automaton
	if err := c.ShouldBindJSON(&fsm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "result": false})
		metrics.IncOperation("fsm", "cleanup", "failure: parameter parsing error")
		return
	}

	ok, err := fsm_s.FSMValidate(&fsm) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的自动机", "result": ok, "error": err.Error()})
		metrics.IncOperation("fsm", "minimize", "failure: invalid fsm")
		return
	}

	fsm_s.Cleanup(&fsm)
	metrics.IncOperation("fsm", "cleanup", "success")
	c.JSON(http.StatusOK, gin.H{"msg": "简化失败", "result": true, "fsm": fsm.ToReactFlow()})
}

func DFAMinimize(c *gin.Context) { // DFA 最小化
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "minimize", time.Since(start).Seconds())
	}()
	var fsm model.Automaton
	if err := c.ShouldBindJSON(&fsm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "result": false})
		metrics.IncOperation("fsm", "minimize", "failure: parameter parsing error")
		return
	}
	fmt.Printf("要最小化的fsm：%+v", fsm)
	ok, err := fsm_s.FSMValidate(&fsm) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的自动机", "result": ok, "error": err.Error()})
		metrics.IncOperation("fsm", "minimize", "failure: invalid fsm")
		return
	}
	if !fsm.IsDFA {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "不是DFA", "result": false})
		metrics.IncOperation("fsm", "minimize", "failure: not dfa")
		return
	}
	new_fsm := fsm_s.DFAMinimize(&fsm) // reduce
	metrics.IncOperation("fsm", "minimize", "success")
	c.JSON(http.StatusOK, gin.H{"msg": "最小化成功", "result": true, "fsm": new_fsm.ToReactFlow()})
}

func FSMStringRecognize(c *gin.Context) { // 字符串识别，
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "string_recognize", time.Since(start).Seconds())
	}()
	type Req struct {
		FSM model.Automaton
		Str string
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "accepted": false})
		metrics.IncOperation("fsm", "string_recognize", "failure: parameter parsing error")
		return
	}
	ok, err := fsm_s.FSMValidate(&req.FSM) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的自动机", "accepted": ok, "error": err.Error()})
		metrics.IncOperation("fsm", "minimize", "failure: invalid fsm")
		return
	}
	ok, err = fsm_s.Recognize(&req.FSM, req.Str) // 先对Str分词，再模拟状态转移
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "识别失败", "accepted": ok, "error": err.Error()})
		metrics.IncOperation("fsm", "string_recognize", "failure: recognition failed")
		return
	}
	metrics.IncOperation("fsm", "string_recognize", "success")
	c.JSON(http.StatusOK, gin.H{"msg": "识别成功", "accepted": ok})
}

func NFADeterminization(c *gin.Context) { // NFA 转 DFA，子集构造法
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("fsm", "nfa_to_dfa", time.Since(start).Seconds())
	}()
	var fsm model.Automaton
	if err := c.ShouldBindJSON(&fsm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "result": false})
		metrics.IncOperation("fsm", "nfa_to_dfa", "failure: parameter parsing error")
		return
	}
	ok, err := fsm_s.FSMValidate(&fsm) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的自动机", "result": false, "error": err.Error()})
		metrics.IncOperation("fsm", "nfa_to_dfa", "failure: invalid fsm")
		return
	}
	if fsm.IsDFA {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "该自动机已经是DFA", "result": false})
		metrics.IncOperation("fsm", "nfa_to_dfa", "failure: not nfa")
		return
	}
	new_fsm := fsm_s.NFAToDFA(&fsm)
	metrics.IncOperation("fsm", "nfa_to_dfa", "success")
	c.JSON(http.StatusOK, gin.H{"msg": "NFA转换为DFA成功", "dfa": new_fsm.ToReactFlow(), "result": true})
}
