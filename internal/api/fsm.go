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
		return
	}

	ok, err := fsm.ISValidate()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的自动机", "result": ok, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "有效", "result": ok, "fsm": fsm.ToReactFlow()})
}

func FSMCleanup(c *gin.Context) { // 去无效符号、无效状态以及相关的转移函数
	var fsm model.Automaton
	if err := c.ShouldBindJSON(&fsm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "result": false})
		return
	}

	fsm_s.Cleanup(&fsm)
	c.JSON(http.StatusOK, gin.H{"msg": "简化失败", "result": true, "fsm": fsm.ToReactFlow()})
}

func DFAMinimize(c *gin.Context) { // DFA 最小化
	var fsm model.Automaton
	if err := c.ShouldBindJSON(&fsm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "result": false})
		return
	}
	fmt.Printf("要最小化的fsm：%+v", fsm)
	ok, err := fsm_s.FSMValidate(&fsm) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的自动机", "result": ok, "error": err.Error()})
		return
	}
	if !fsm.IsDFA {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "不是DFA", "result": false})
		return
	}
	new_fsm := fsm_s.DFAMinimize(&fsm) // reduce
	c.JSON(http.StatusOK, gin.H{"msg": "最小化成功", "result": true, "fsm": new_fsm.ToReactFlow()})
}

func FSMStringRecognize(c *gin.Context) { // 字符串识别，
	type Req struct {
		FSM model.Automaton
		Str string
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "result": false})
		return
	}
	req.FSM.CheckIsDFA()                          // 先判断类型，确定isDFA
	ok, err := fsm_s.Recognize(&req.FSM, req.Str) // 先对Str分词，再模拟状态转移
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "识别失败", "result": ok, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "识别成功", "result": ok})
}

func NFADeterminization(c *gin.Context) { // NFA 转 DFA，子集构造法
	var fsm model.Automaton
	if err := c.ShouldBindJSON(&fsm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析错误", "result": false})
		return
	}
	ok, err := fsm_s.FSMValidate(&fsm) // 包含了DFA还是NFA的判断
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的自动机", "result": false, "error": err.Error()})
		return
	}
	if fsm.IsDFA {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "不是NFA", "result": false})
		return
	}
	new_fsm := fsm_s.NFAToDFA(&fsm)
	c.JSON(http.StatusOK, gin.H{"msg": "NFA转换为DFA成功", "dfa": new_fsm.ToReactFlow(), "result": true})
}
