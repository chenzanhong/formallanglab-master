/*
GrammarValidate        	// 文法校验——是否有效
GrammarAmbiguityCheck   // 正则文法二义性判断，其他型的文法的二义性是不可判定问题
GrammarStringRecognize  // 字符串识别——是否被指定文法所接受；（可选）扩展：返回递归下降分析、LL(1)分析、LR(0)分析或LR(1)分析的过程
GrammarTypeDetermine    // 判断所给文法的类型
GrammarEquivalenceCheck // 判断所给的两个正则文法是否等价
GrammarSimplify			// 文法的化简——去无用符号（不可派生、不可达）、单一产生式、空产生式
GrammarTree			// 生成树（是否实现待确定）
GrammarFirstSet			// 计算文法的First集
GrammarFollowSet		// 计算文法的Follow集
*/
package api

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	"backend/internal/metrics"
	"backend/internal/service/grammar_s"
	"backend/logs"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 文法校验——是否有效
func GrammarValidate(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("grammar", "validate", time.Since(start).Seconds())
	}()
	var req dto.GrammarValidateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("grammar", "validate", "failure: parameter parsing error")
		logs.Sugar.Warnw("文法校验失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.GrammarValidateResponse{
			Valid: false,
			Msg:   "参数解析失败",
			// Error:  err.Error(),
			Result: false,
		})
		return
	}

	// 先统一 ε 表示
	fmt.Printf("%+v\n", req.Grammar)
	normalizeGrammar(&req.Grammar)
	fmt.Printf("%+v\n", req.Grammar)
	// 如果文法结构不完整，先完善结构
	// req.Grammar = *grammar_s.CompleteGrammarStructure(&req.Grammar)

	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		metrics.IncOperation("grammar", "validate", "failure: invalid grammar")
		logs.Sugar.Warnw("文法校验失败", "detail", "文法格式无效")
		c.JSON(http.StatusBadRequest, dto.GrammarValidateResponse{
			Valid:  false,
			Msg:    "invalid grammar：" + err.Error(),
			Result: false,
		})
		return
	}

	metrics.IncOperation("grammar", "validate", "success")
	logs.Sugar.Infow("文法校验成功")
	c.JSON(http.StatusOK, dto.GrammarValidateResponse{
		Valid:  true,
		Msg:    "grammar is valid",
		Type:   req.Grammar.GrammarType,
		Result: true,
	})
}

// 正则文法二义性判断，其他型的文法的二义性是不可判定问题
func GrammarAmbiguityCheck(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("grammar", "ambiguity_check", time.Since(start).Seconds())
	}()
	var req dto.GrammarAmbiguityCheckRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("grammar", "ambiguity_check", "failure: parameter parsing error")
		logs.Sugar.Warnw("文法二义性检查失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.GrammarAmbiguityCheckResponse{
			Msg: "参数解析失败",
			// Error:  err.Error(),
			Result: false,
		})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&req.Grammar)
	// 如果文法结构不完整，先完善结构
	// req.Grammar = *grammar_s.CompleteGrammarStructure(&req.Grammar)

	// 可选：先校验文法有效性
	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		metrics.IncOperation("grammar", "ambiguity_check", "failure: invalid grammar")
		logs.Sugar.Warnw("文法二义性检查失败", "detail", "文法格式无效")
		c.JSON(http.StatusOK, dto.GrammarAmbiguityCheckResponse{
			Msg:  "无效文法：" + err.Error(),
			Type: model.InvalidGrammar,
		})
		return
	}

	// === 1. 先判断文法类型 ===
	grammarType := grammar_s.TypeDetermine(&req.Grammar)

	if grammarType != model.RegularGrammar {
		metrics.IncOperation("grammar", "ambiguity_check", "failure: not regular grammar")
		logs.Sugar.Warnw("文法二义性检查失败", "detail", "文法类型不是正则文法，二义性不可判定")
		c.JSON(http.StatusOK, dto.GrammarAmbiguityCheckResponse{
			Msg:         "文法类型不是正则文法（3型），其二义性为不可判定问题",
			Type:        grammarType,
			IsAmbiguous: nil, // 无法判断
		})
		return
	}

	// === 2. 检查正则文法是否二义 ===
	isAmbiguous, err := grammar_s.IsAmbiguousRegular(&req.Grammar)
	if err != nil {
		metrics.IncOperation("grammar", "ambiguity_check", "failure: ambiguity check error")
		logs.Sugar.Warnw("文法二义性检查失败", "detail", "二义性判断过程中发生错误")
		c.JSON(http.StatusInternalServerError, dto.GrammarAmbiguityCheckResponse{
			Msg:  "二义性判断出错：" + err.Error(),
			Type: model.RegularGrammar,
			// Error:  err.Error(),
			Result: false,
		})
		return
	}

	// === 3. 返回结果 ===
	var resultMsg string
	if isAmbiguous {
		resultMsg = "该正则文法是二义的"
	} else {
		resultMsg = "该正则文法是无二义的"
	}

	metrics.IncOperation("grammar", "ambiguity_check", "success")
	logs.Sugar.Infow("文法二义性检查成功")
	c.JSON(http.StatusOK, dto.GrammarAmbiguityCheckResponse{
		Msg:         resultMsg,
		IsAmbiguous: &isAmbiguous,
		Type:        model.RegularGrammar,
		Result:      true,
	})
}

// 字符串识别——是否被指定文法所接受；（可选）扩展：返回递归下降分析、LL(1)分析、LR(0)分析或LR(1)分析的过程
func GrammarStringRecognize(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("grammar", "string_recognize", time.Since(start).Seconds())
	}()
	var req dto.GrammarStringRecognizeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("grammar", "string_recognize", "failure: parameter parsing error")
		logs.Sugar.Warnw("文法字符串识别失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.GrammarStringRecognizeResponse{
			Msg: "参数解析失败",
			// Error:  err.Error(),
			Result: false,
		})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&req.Grammar)
	// 如果文法结构不完整，先完善结构
	// req.Grammar = *grammar_s.CompleteGrammarStructure(&req.Grammar)

	// 可选：先校验文法有效性
	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		metrics.IncOperation("grammar", "string_recognize", "failure: invalid grammar")
		logs.Sugar.Warnw("文法字符串识别失败", "detail", "文法格式无效")
		c.JSON(http.StatusBadRequest, dto.GrammarStringRecognizeResponse{
			Accepted: false,
			Msg:      "无效文法: " + err.Error(),
			Result:   false,
		})
		return
	}

	// 默认模式为 auto
	if req.Mode == "" {
		req.Mode = "auto"
	}

	// 直接将输入字符串传递给分析函数，不再提前转换为符号
	result := grammar_s.ParseStringWithMode(&req.Grammar, req.Input, req.Mode, req.ShowSteps)

	metrics.IncOperation("grammar", "string_recognize", "success")
	logs.Sugar.Infow("文法字符串识别成功")

	// 返回结果
	response := dto.GrammarStringRecognizeResponse{
		Accepted: result.Accepted,
		Method:   result.Method,
	}

	if result.Message != "" {
		response.Msg = result.Message
	}

	if req.ShowSteps && len(result.Steps) > 0 {
		response.Steps = result.Steps
	}

	if result.Error != "" {
		response.Msg += result.Error
	}

	c.JSON(http.StatusOK, dto.GrammarStringRecognizeResponse{
		Accepted: response.Accepted,
		Method:   response.Method,
		Msg:      response.Msg,
		Steps:    response.Steps,
		Result:   true,
	})
}

// 判断所给文法的类型
func GrammarTypeDetermine(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("grammar", "type_determine", time.Since(start).Seconds())
	}()
	var req dto.GrammarTypeDetermineRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("grammar", "type_determine", "failure: parameter parsing error")
		logs.Sugar.Warnw("文法类型判断失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.GrammarTypeDetermineResponse{
			Type: model.InvalidGrammar,
			// Error:  err.Error(),
			Result: false,
			Msg:    "参数解析失败",
		})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&req.Grammar)
	// 如果文法结构不完整，先完善结构
	// req.Grammar = *grammar_s.CompleteGrammarStructure(&req.Grammar)

	// 可选：先校验文法有效性
	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		metrics.IncOperation("grammar", "type_determine", "failure: invalid grammar")
		logs.Sugar.Warnw("文法类型判断失败", "detail", "文法格式无效")
		c.JSON(http.StatusOK, dto.GrammarTypeDetermineResponse{
			Type:     -1,
			TypeName: "无效文法：" + err.Error(),
			Result:   false,
		})
		return
	}

	grammarType := grammar_s.TypeDetermine(&req.Grammar)

	typeName := model.GrammarTypeNameMap[grammarType]

	metrics.IncOperation("grammar", "type_determine", "success")
	logs.Sugar.Infow("文法类型判断成功")
	c.JSON(http.StatusOK, dto.GrammarTypeDetermineResponse{
		Type:     grammarType,
		TypeName: typeName,
		Result:   true,
	})
}

// 文法的化简——去无用符号（不可派生、不可达）、单一产生式、空产生式
func GrammarSimplify(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("grammar", "simplify", time.Since(start).Seconds())
	}()
	var req dto.GrammarSimplifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("grammar", "simplify", "failure: parameter parsing error")
		logs.Sugar.Warnw("文法化简失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.GrammarSimplifyResponse{
			Msg: "参数解析失败",
			// Error:  err.Error(),
			Result: false,
		})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&req.Grammar)
	// 如果文法结构不完整，先完善结构
	// grammar = *grammar_s.CompleteGrammarStructure(&grammar)

	// 可选：先校验文法有效性
	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		metrics.IncOperation("grammar", "simplify", "failure: invalid grammar")
		logs.Sugar.Warnw("文法化简失败", "detail", "文法格式无效")
		c.JSON(http.StatusOK, dto.GrammarSimplifyResponse{
			Msg:    "无效文法：" + err.Error(),
			Result: false,
		})
		return
	}

	new_grammar := grammar_s.Simplify(&req.Grammar)
	metrics.IncOperation("grammar", "simplify", "success")
	logs.Sugar.Infow("文法化简成功")
	c.JSON(http.StatusOK, dto.GrammarSimplifyResponse{
		Msg:     "简化成功",
		Grammar: new_grammar,
		Result:  true,
	})
}

// 判断所给的两个正则文法是否等价
func GrammarEquivalenceCheck(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("grammar", "equivalence_check", time.Since(start).Seconds())
	}()
	type Req dto.GrammarEquivalenceCheckRequest
	var req Req

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("grammar", "equivalence_check", "failure: parameter parsing error")
		logs.Sugar.Warnw("文法等价性检查失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.GrammarEquivalenceCheckResponse{
			Msg: "参数解析失败",
			// Error:        err.Error(),
			IsEquivalent: false,
			Result:       false,
		})
		return
	}
	fmt.Printf("%+v\n", req.Grammar1)
	fmt.Printf("%+v\n", req.Grammar2)
	// 先统一 ε 表示
	normalizeGrammar(&req.Grammar1)
	normalizeGrammar(&req.Grammar2)
	fmt.Printf("统一空转移符号后%+v\n", req.Grammar1)
	fmt.Printf("统一空转移符号后%+v\n", req.Grammar2)
	// 如果文法结构不完整，先完善结构
	// req.Grammar1 = *grammar_s.CompleteGrammarStructure(&req.Grammar1)
	// req.Grammar2 = *grammar_s.CompleteGrammarStructure(&req.Grammar2)

	// 先判断两个文法是否有效，且为正则文法
	if grammar_s.GrammarCheckValidity(&req.Grammar1) != nil {
		metrics.IncOperation("grammar", "equivalence_check", "failure: invalid grammar1")
		logs.Sugar.Warnw("文法等价性检查失败", "detail", "第一个文法格式无效")
		c.JSON(http.StatusBadRequest, dto.GrammarEquivalenceCheckResponse{
			Msg:          "invalid grammar1",
			IsEquivalent: false,
			Result:       false,
		})
		return
	}
	if grammar_s.TypeDetermine(&req.Grammar1) != model.RegularGrammar { // 非正则文法
		metrics.IncOperation("grammar", "equivalence_check", "failure: grammar1 not regular")
		logs.Sugar.Warnw("文法等价性检查失败", "detail", "第一个文法不是正则文法")
		c.JSON(http.StatusBadRequest, dto.GrammarEquivalenceCheckResponse{
			Msg:          "grammar1 is not regular grammar",
			IsEquivalent: false,
			Result:       false,
		})
		return
	}
	if grammar_s.GrammarCheckValidity(&req.Grammar2) != nil {
		metrics.IncOperation("grammar", "equivalence_check", "failure: invalid grammar2")
		logs.Sugar.Warnw("文法等价性检查失败", "detail", "第二个文法格式无效")
		c.JSON(http.StatusBadRequest, dto.GrammarEquivalenceCheckResponse{
			Msg:          "invalid grammar2",
			IsEquivalent: false,
			Result:       false,
		})
		return
	}
	if grammar_s.TypeDetermine(&req.Grammar2) != model.RegularGrammar { // 非正则文法
		metrics.IncOperation("grammar", "equivalence_check", "failure: grammar2 not regular")
		logs.Sugar.Warnw("文法等价性检查失败", "detail", "第二个文法不是正则文法")
		c.JSON(http.StatusBadRequest, dto.GrammarEquivalenceCheckResponse{
			Msg:          "grammar2 is not regular grammar",
			IsEquivalent: false,
			Result:       false,
		})
		return
	}

	// 再判断是否等价
	if !grammar_s.IsEquivalent(&req.Grammar1, &req.Grammar2) {
		metrics.IncOperation("grammar", "equivalence_check", "success: not equivalent")
		logs.Sugar.Infow("文法等价性检查成功", "detail", "两个文法不等价")
		c.JSON(http.StatusOK, dto.GrammarEquivalenceCheckResponse{
			Msg:          "this two grammars are not equivalent",
			IsEquivalent: false,
			Result:       true,
		})
		return
	}
	metrics.IncOperation("grammar", "equivalence_check", "success: equivalent")
	logs.Sugar.Infow("文法等价性检查成功", "detail", "两个文法等价")
	c.JSON(http.StatusOK, dto.GrammarEquivalenceCheckResponse{
		Msg:          "this two grammars are equivalent",
		IsEquivalent: true,
		Result:       true,
	})
}

// 计算文法的First集
func GrammarFirstSet(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("grammar", "first_set", time.Since(start).Seconds())
	}()
	var req dto.GrammarFirstSetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("grammar", "first_set", "failure: parameter parsing error")
		logs.Sugar.Warnw("计算First集失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.GrammarFirstSetResponse{
			Msg:    "参数解析失败",
			Result: false,
		})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&req.Grammar)

	// 校验文法有效性
	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		metrics.IncOperation("grammar", "first_set", "failure: invalid grammar")
		logs.Sugar.Warnw("计算First集失败", "detail", "文法格式无效")
		c.JSON(http.StatusBadRequest, dto.GrammarFirstSetResponse{
			Msg:    "无效文法：" + err.Error(),
			Result: false,
		})
		return
	}

	// 计算First集
	firstSet := grammar_s.CalculateFirst(&req.Grammar)

	// 转换为响应格式
	result := convertMapToStringSlice(firstSet)

	metrics.IncOperation("grammar", "first_set", "success")
	logs.Sugar.Infow("计算First集成功")
	c.JSON(http.StatusOK, dto.GrammarFirstSetResponse{
		FirstSet: result,
		Msg:      "计算First集成功",
		Result:   true,
	})
}

// 计算文法的Follow集
func GrammarFollowSet(c *gin.Context) {
	start := time.Now()
	defer func() {
		metrics.ObserveOperationDuration("grammar", "follow_set", time.Since(start).Seconds())
	}()
	var req dto.GrammarFollowSetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.IncOperation("grammar", "follow_set", "failure: parameter parsing error")
		logs.Sugar.Warnw("计算Follow集失败", "detail", "参数解析失败，请检查请求格式是否正确")
		c.JSON(http.StatusBadRequest, dto.GrammarFollowSetResponse{
			Msg:    "参数解析失败",
			Result: false,
		})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&req.Grammar)

	// 校验文法有效性
	if err := grammar_s.GrammarCheckValidity(&req.Grammar); err != nil {
		metrics.IncOperation("grammar", "follow_set", "failure: invalid grammar")
		logs.Sugar.Warnw("计算Follow集失败", "detail", "文法格式无效")
		c.JSON(http.StatusBadRequest, dto.GrammarFollowSetResponse{
			Msg:    "无效文法：" + err.Error(),
			Result: false,
		})
		return
	}

	// 计算First集（用于Follow集计算）
	firstSet := grammar_s.CalculateFirst(&req.Grammar)

	// 计算Follow集
	followSet := grammar_s.CalculateFollow(&req.Grammar, firstSet)

	// 转换为响应格式
	result := convertMapToStringSlice(followSet)

	metrics.IncOperation("grammar", "follow_set", "success")
	logs.Sugar.Infow("计算Follow集成功")
	c.JSON(http.StatusOK, dto.GrammarFollowSetResponse{
		FollowSet: result,
		Msg:       "计算Follow集成功",
		Result:    true,
	})
}

// normalizeGrammar 将前端传入的文法中的 "ε" 或空字符串统一规范为后端使用的 model.Epsilon("ε") 表示
func normalizeGrammar(g *model.Grammar) {
	// 终结符集合规范化
	for i, t := range g.Terminals {
		if string(t) == "ε" || string(t) == "" {
			g.Terminals[i] = model.Epsilon
		}
	}
	// 产生式左右部规范化
	for pi := range g.Productions {
		// 左部兜底
		for li, ls := range g.Productions[pi].Left {
			if string(ls) == "ε" {
				g.Productions[pi].Left[li] = model.Epsilon
			}
		}
		// 右部：将 "ε"/"" 统一为单符号 [Epsilon]
		right := g.Productions[pi].Right
		if len(right) == 0 {
			g.Productions[pi].Right = []model.Symbol{model.Epsilon}
			continue
		}
		allEmpty := len(right) > 0
		for ri, rs := range right {
			if string(rs) == "ε" || string(rs) == "" {
				right[ri] = model.Epsilon
			} else {
				allEmpty = false
			}
		}
		if allEmpty {
			g.Productions[pi].Right = []model.Symbol{model.Epsilon}
		} else {
			g.Productions[pi].Right = right
		}
	}
}

// convertMapToStringSlice 将 map[model.Symbol]map[model.Symbol]struct{} 转换为 map[string][]string
func convertMapToStringSlice(original map[model.Symbol]map[model.Symbol]struct{}) map[string][]string {
	result := make(map[string][]string)
	for symbol, symbolsMap := range original {
		symbols := make([]string, 0, len(symbolsMap))
		for s := range symbolsMap {
			symbols = append(symbols, string(s))
		}
		result[string(symbol)] = symbols
	}
	return result
}
