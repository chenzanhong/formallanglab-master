/*
		文法与自动机的转换
		GrammarToFA	// 正则文法转成 FA，自然
	    FAToGrammar    // FA 转成文法
		RegexToFA		// 正则表达式转为 FA，Thompson 构造法
		FAToRegex		// FA 转正则表达式
*/
package api

import (
	"backend/internal/domain/dto"
	"backend/internal/metrics"
	"backend/internal/service/convert_s"
	"backend/internal/service/grammar_s"
	re "backend/internal/service/regex_s"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	// 检查是否为有效文法
	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		c.JSON(http.StatusBadRequest, dto.GrammarToFAResponse{
			Msg:    "无效的文法，请检查文法规则" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("convert", "grammar_to_nfa", "failure: invalid grammar")
		return
	}

	// 检查是否为三型文法（正则文法）
	isRegular, isRightLinear := grammar_s.IsRegular(&req.Grammar)
	if !isRegular {
		c.JSON(http.StatusBadRequest, dto.GrammarToFAResponse{
			Msg:    "非正则文法，无法转换为NFA",
			Result: false,
		})
		metrics.IncOperation("convert", "grammar_to_nfa", "failure: not regular grammar")
		return
	}
	fmt.Printf("文法转自动机，文法：%v", req.Grammar)
	automaton := convert_s.RegularGrammarToFA(&req.Grammar, isRightLinear)

	metrics.IncOperation("convert", "grammar_to_nfa", "success")
	// 响应结果
	c.JSON(http.StatusOK, dto.GrammarToFAResponse{
		Msg:           "正则文法转自动机成功",
		Result:        true,
		Automaton:     automaton,
		AutomatonFlow: automaton.ToReactFlow(),
	})
}

func FAToGrammar(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "fa_to_grammar", time.Since(start).Seconds())
	}()
	// 1.请求参数解析
	var req dto.FAToGrammarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FAToGrammarResponse{
			Msg:    "参数解析失败，请检查参数格式",
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_grammar", "failure: parameter parsing error")
		return
	}
	// 2.判断是否为有效自动机
	if req.Automaton.ISValidate() != nil {
		c.JSON(http.StatusBadRequest, dto.FAToGrammarResponse{
			Msg:    "无效的自动机，请检查自动机状态",
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_grammar", "failure: invalid automaton")
		return
	}
	fmt.Printf("自动机转文法：%v", req.Automaton)

	// 3.调用convert_s提供的方法进行转换
	grammar := convert_s.FAToGrammar(&req.Automaton)

	// 4.响应结果
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
			Msg:    "Invalid request: " + err.Error(),
			Result: false,
		})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: parameter parsing error")
		return
	}

	// 先检验是否为有效的正则表达式
	if err := re.RegexValidate(req.Pattern); err != nil {
		c.JSON(http.StatusBadRequest, dto.RegexToFAResponse{
			Msg:    "invalid regular expression",
			Result: false,
		})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: invalid regex")
		return
	}

	nfa, err := convert_s.RegexToFA(string(req.Pattern))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.RegexToFAResponse{
			Msg:    "转换失败：" + err.Error(),
			Result: false,
		})
		metrics.IncOperation("convert", "regex_to_nfa", "failure: conversion failed")
		return
	}
	metrics.IncOperation("convert", "regex_to_nfa", "success")
	c.JSON(http.StatusOK, dto.RegexToFAResponse{
		Msg:           "正则表达式转NFA成功",
		Result:        true,
		Automaton:     nfa,
		AutomatonFlow: nfa.ToReactFlow(),
	})
}

// 自动机到正则表达式的转换
//
// ✅ 核心思想：
// 任何有限自动机（DFA/NFA/ε-NFA）识别的语言都是正则语言，因此一定存在一个等价的正则表达式
// 我们的目标是：通过算法，从自动机构造出这个正则表达式
//
// ✅ 主流方法：状态消去法（State Elimination Method）
// 这是最通用、最直观的方法，适用于含 ε 转移的自动机
// 📌 基本步骤：
// 1. 标准化自动机（可选但推荐）：
//   - 添加一个新的唯一初态 q_start，通过 ε 转移到原初态
//   - 添加一个新的唯一终态 q_final，所有原终态通过 ε 转移到它
//   - 确保初态无入边，终态无出边
//
// 2. 逐步删除中间状态（非初非终），每次删除一个状态时，更新其前驱到后继的转移标签，用正则表达式合并路径
// 3. 最后只剩初态和终态，它们之间的转移标签就是所求的正则表达式
//
// ✅ 处理空转移（ε-transitions）：
// - ε 就是正则表达式中的 ε（或写作 λ），在合并路径时按正则表达式规则处理
// - 在状态消去过程中，ε 和其他符号一样参与运算
// - 最终结果中通常会自动"吸收"掉 ε（因为 r + ε 或 rε = r）
// - 无需预先消除 ε 转移！状态消去法天然支持 ε
//
// ✅ 状态消去的规则（重点！）
// 假设要删除状态 q，它有：
// - 入边：从 p ->R-> q
// - 出边：从 q ->S-> r
// - 自环：q ->T-> q
// 那么，删除 q 后，需为每对 (p, r) 添加（或更新）转移：
// p ->(R·T*·S)-> r
// 其中 R, S, T 是正则表达式，若有多条入边或出边，先用 + 合并
func FAToRegex(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("convert", "fa_to_regex", time.Since(start).Seconds())
	}()

	// 1.请求参数解析
	var req dto.FAToRegexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FAToRegexResponse{
			Msg:    "参数解析失败，请检查参数格式",
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_regex", "failure: parameter parsing error")
		return
	}

	// 2.判断是否为有效自动机
	if req.Automaton.ISValidate() != nil {
		c.JSON(http.StatusBadRequest, dto.FAToRegexResponse{
			Msg:    "无效的自动机，请检查自动机状态",
			Result: false,
		})
		metrics.IncOperation("convert", "fa_to_regex", "failure: invalid automaton")
		return
	}

	// 简化自动机，删除不可达状态和不可派生状态

	// 3.调用convert_s提供的方法进行转换
	regex := convert_s.FAToRegex(&req.Automaton)

	// 4.响应结果
	c.JSON(http.StatusOK, dto.FAToRegexResponse{
		Msg:     "自动机转正则表达式成功",
		Result:  true,
		Pattern: regex,
	})
	metrics.IncOperation("convert", "fa_to_regex", "success")
}
