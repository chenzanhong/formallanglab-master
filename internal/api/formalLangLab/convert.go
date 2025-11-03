/*
		文法与自动机的转换
		GrammarToNFA	// 正则文法转成 NFA，自然
	    FAToGrammar    // DFA 转成文法
		RegexToNFA		// 正则表达式转为NFA，Thompson 构造法
*/
package api

import (
	"backend/internal/metrics"
	"backend/internal/service/convert_s"
	re "backend/internal/service/regex_s"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GrammarToNFA(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "grammar_to_nfa", time.Since(start).Seconds())
	}()

	metrics.IncOperation("convert", "grammar_to_nfa", "success")
}

func FAToGrammar(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "fa_to_grammar", time.Since(start).Seconds())
	}()
	// 1.请求参数解析

	// 2.判断是否为有效自动机

	// 3.调用convert_s提供的方法进行转换

	// 4.响应结果

	metrics.IncOperation("convert", "fa_to_grammar", "success")
}

func RegexToNFA(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "regex_to_nfa", time.Since(start).Seconds())
	}()
	type Req struct {
		Pattern string `json:"pattern" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: parameter parsing error")
		return
	}

	// 先检验是否为有效的正则表达式
	_, err := re.RegexValidate(req.Pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid regular expression", "valid": false, "error": err})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: invalid regex")
		return
	}

	nfa, err := convert_s.RegexToNFA(req.Pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "转换失败，输入为空或无效的正则表达式", "result": false, "error": err.Error()})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: conversion failed")
		return
	}
	metrics.IncOperation("convert", "regex_to_nfa", "success")
	c.JSON(http.StatusOK, gin.H{"msg": "转换成功", "result": true, "fsm": nfa.ToReactFlow()})
}

// 自动机到正则表达式的转换
func FAToRegex(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "fa_to_regex", time.Since(start).Seconds())
	}()

	// 1.请求参数解析

	// 2.判断是否为有效自动机以及类型

	// 3.如果为NFA，先确定化为DFA

	// 4.调用convert_s提供的方法进行转换

	// 5.响应结果

	metrics.IncOperation("convert", "fa_to_regex", "success")
}
