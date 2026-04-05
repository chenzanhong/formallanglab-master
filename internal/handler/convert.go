/*
文法与自动机的转换
GrammarToFA		// 正则文法转成 FA，自然
FAToGrammar     // FA 转成文法
RegexToFA		// 正则表达式转为 FA，Thompson 构造法
FAToRegex		// FA 转正则表达式
*/
package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chenzanhong/zlog"
	"github.com/gin-gonic/gin"

	"github.com/chenzanhong/formallanglab-master/internal/domain/dto"
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/metrics"
	"github.com/chenzanhong/formallanglab-master/internal/service/automaton_s"
	"github.com/chenzanhong/formallanglab-master/internal/service/grammar_s"
	"github.com/chenzanhong/formallanglab-master/internal/service/regex_s"
)

func GrammarToFA(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "grammar_to_nfa", time.Since(start).Seconds())
	}()

	var req dto.GrammarToFARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.GrammarToFAResponse{
			Msg:    "参数解析失败，请检查参数格式",
			Result: false,
		})
		metrics.IncOperation("convert", "grammar_to_nfa", "failure: parameter parsing error")

		return
	}

	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		c.JSON(http.StatusBadRequest, dto.GrammarToFAResponse{
			Msg:    "无效的文法，请检查文法规则" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("convert", "grammar_to_nfa", "failure: invalid grammar")

		return
	}

	isRegular, isRightLinear := grammar_s.IsRegular(&req.Grammar)
	if !isRegular {
		c.JSON(http.StatusBadRequest, dto.GrammarToFAResponse{
			Msg:    "非正则文法，无法转换为 NFA",
			Result: false,
		})
		metrics.IncOperation("convert", "grammar_to_nfa", "failure: not regular grammar")

		return
	}

	process := grammar_s.RegularGrammarToFAWithProcess(&req.Grammar, isRightLinear)

	metrics.IncOperation("convert", "grammar_to_nfa", "success")

	c.JSON(http.StatusOK, dto.GrammarToFAResponse{
		Msg:     "正则文法转自动机成功",
		Result:  true,
		Process: process,
	})
}

func FAToGrammar(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "fa_to_grammar", time.Since(start).Seconds())
	}()

	var req dto.FAToGrammarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FAToGrammarResponse{
			Msg:    "参数解析失败，请检查参数格式",
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_grammar", "failure: parameter parsing error")

		return
	}

	if req.Automaton.Validate() != nil {
		c.JSON(http.StatusBadRequest, dto.FAToGrammarResponse{
			Msg:    "无效的自动机，请检查自动机状态",
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_grammar", "failure: invalid automaton")

		return
	}

	grammar := automaton_s.FAToGrammar(&req.Automaton)

	c.JSON(http.StatusOK, dto.FAToGrammarResponse{
		Msg:     "自动机转文法成功",
		Result:  true,
		Grammar: grammar,
	})

	metrics.IncOperation("convert", "fa_to_grammar", "success")
}

func RegexToFA(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "regex_to_nfa", time.Since(start).Seconds())
	}()

	var req dto.RegexToFARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.RegexToFAResponse{
			Msg:    "无效请求：" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: parameter parsing error")

		return
	}

	if err := regex_s.RegexValidate(req.Pattern); err != nil {
		zlog.Warnw("Regex to FA conversion failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.RegexToFAResponse{
			Msg:    fmt.Sprintf("无效的正则表达式：%v", err),
			Result: false,
		})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: invalid regex")

		return
	}

	process, err := regex_s.RegexToFAWithSteps(req.Pattern)
	if err != nil {
		zlog.Errorw("Regex to FA conversion failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.RegexToFAResponse{
			Msg:    fmt.Sprintf("转换失败：%v", err),
			Result: false,
		})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: conversion failed")

		return
	}

	metrics.IncOperation("convert", "regex_to_nfa", "success")

	c.JSON(http.StatusOK, dto.RegexToFAResponse{
		Msg:    "正则表达式转 NFA 成功",
		Result: true,
		Process: &model.RegexToFAProcess{
			Regex:          req.Pattern,
			Steps:          process.Steps,
			FinalAutomaton: process.FinalAutomaton,
		},
	})
}

func FAToRegex(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "fa_to_regex", time.Since(start).Seconds())
	}()

	var req dto.FAToRegexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FAToRegexResponse{
			Msg:    "参数解析失败，请检查参数格式",
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_regex", "failure: parameter parsing error")

		return
	}

	if req.Automaton.Validate() != nil {
		c.JSON(http.StatusBadRequest, dto.FAToRegexResponse{
			Msg:    "无效的自动机，请检查自动机状态",
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_regex", "failure: invalid automaton")

		return
	}

	process, err := automaton_s.FAToRegexWithProcess(&req.Automaton)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.FAToRegexResponse{
			Msg:    "转换失败：" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_regex", "failure: conversion failed")

		return
	}

	c.JSON(http.StatusOK, dto.FAToRegexResponse{
		Msg:     "自动机转正则表达式成功",
		Result:  true,
		Pattern: process.FinalRegex,
		Process: process,
	})

	metrics.IncOperation("convert", "fa_to_regex", "success")
}
