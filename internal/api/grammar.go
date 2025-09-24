/*
GrammarValidate        	// 文法校验——是否有效
GrammarAmbiguityCheck   // 正则文法二义性判断，其他型的文法的二义性是不可判定问题
GrammarStringRecognize  // 字符串识别——是否被指定文法所接受；（可选）扩展：返回递归下降分析、LL(1)分析、LR(0)分析或LR(1)分析的过程
GrammarTypeDetermine    // 判断所给文法的类型
GrammarEquivalenceCheck // 判断所给的两个正则文法是否等价
GrammarSimplify			// 文法的化简——去无用符号（不可派生、不可达）、单一产生式、空产生式
*/
package api

import (
	"backend/internal/domain/model"
	"backend/internal/service/grammar_s"
	"net/http"

	"github.com/gin-gonic/gin"
)

// normalizeGrammar 将前端传入的文法中的 "ε" 或空字符串统一规范为后端使用的 model.Epsilon("") 表示
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

func GrammarValidate(c *gin.Context) { // 文法校验——是否有效
	var grammar model.Grammar

	if err := c.ShouldBindJSON(&grammar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析失败", "error": err.Error()})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&grammar)
	// 如果文法结构不完整，先完善结构
	grammar = *grammar_s.CompleteGrammarStructure(&grammar)

	if !grammar_s.IsValidGrammar(&grammar) {
		c.JSON(http.StatusBadRequest, gin.H{"valid": false, "msg": "invalid grammar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "msg": "grammar is valid"})
}

func GrammarAmbiguityCheck(c *gin.Context) { // 正则文法二义性判断，其他型的文法的二义性是不可判定问题
	var grammar model.Grammar

	if err := c.ShouldBindJSON(&grammar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "参数解析失败",
			"error": err.Error(),
		})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&grammar)
	// 如果文法结构不完整，先完善结构
	grammar = *grammar_s.CompleteGrammarStructure(&grammar)

	// 可选：先校验文法有效性
	if !grammar_s.IsValidGrammar(&grammar) {
		c.JSON(http.StatusOK, gin.H{
			"type":     -1,
			"typeName": "无效文法",
		})
		return
	}

	// === 1. 先判断文法类型 ===
	grammarType := grammar_s.TypeDetermine(&grammar)

	if grammarType != grammar_s.Type3 {
		c.JSON(http.StatusOK, gin.H{
			"msg":         "文法类型不是正则文法（3型），其二义性为不可判定问题",
			"type":        grammarType,
			"isAmbiguous": nil, // 无法判断
			"isRegular":   false,
		})
		return
	}

	// === 2. 检查正则文法是否二义 ===
	isAmbiguous, err := grammar_s.IsAmbiguousRegular(&grammar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "二义性判断出错",
			"error": err.Error(),
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

	c.JSON(http.StatusOK, gin.H{
		"msg":         resultMsg,
		"isAmbiguous": isAmbiguous,
		"isRegular":   true,
		"type":        3,
	})
}

func GrammarStringRecognize(c *gin.Context) { // 字符串识别——是否被指定文法所接受；（可选）扩展：返回递归下降分析、LL(1)分析、LR(0)分析或LR(1)分析的过程
	var req struct {
		Grammar   model.Grammar `json:"grammar" binding:"required"`
		Input     string        `json:"input" binding:"required"`
		ShowSteps bool          `json:"showSteps"` // 是否返回分析步骤
		Mode      string        `json:"mode"`      // 分析模式: "auto", "ll1", "ll1_recovery", "recursive_descent", "lr0", "lr1", "bfs"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "参数解析失败",
			"error": err.Error(),
		})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&req.Grammar)
	// 如果文法结构不完整，先完善结构
	req.Grammar = *grammar_s.CompleteGrammarStructure(&req.Grammar)

	// 可选：先校验文法有效性
	if !grammar_s.IsValidGrammar(&req.Grammar) {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":     -1,
			"typeName": "无效文法",
		})
		return
	}

	// 默认模式为 auto
	if req.Mode == "" {
		req.Mode = "auto"
	}

	input := req.Input
	if input == "" {
		input = "ε"
	}
	inputSymbols := stringToSymbols(input, &req.Grammar)

	// 使用新的统一分析接口
	result := grammar_s.ParseStringWithMode(&req.Grammar, inputSymbols, req.Mode, req.ShowSteps)

	// 返回结果
	response := gin.H{
		"accepted": result.Accepted,
		"method":   result.Method,
	}

	if result.Message != "" {
		response["message"] = result.Message
	}

	if req.ShowSteps && len(result.Steps) > 0 {
		response["steps"] = result.Steps
	}

	if result.Error != "" {
		response["error"] = result.Error
	}

	c.JSON(http.StatusOK, response)
}

func GrammarTypeDetermine(c *gin.Context) { // 判断所给文法的类型
	var grammar model.Grammar

	if err := c.ShouldBindJSON(&grammar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析失败", "error": err.Error()})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&grammar)
	// 如果文法结构不完整，先完善结构
	grammar = *grammar_s.CompleteGrammarStructure(&grammar)

	// 可选：先校验文法有效性
	if !grammar_s.IsValidGrammar(&grammar) {
		c.JSON(http.StatusOK, gin.H{
			"type":     -1,
			"typeName": "无效文法",
		})
		return
	}

	typ := grammar_s.TypeDetermine(&grammar)

	typeName := map[int]string{
		grammar_s.Type3: "3型文法（正则文法）",
		grammar_s.Type2: "2型文法（上下文无关文法）",
		grammar_s.Type1: "1型文法（上下文有关文法）",
		grammar_s.Type0: "0型文法（无限制文法）",
		-1:              "无效文法",
	}[typ]

	c.JSON(http.StatusOK, gin.H{
		"type":     typ,
		"typeName": typeName,
	})
}

func GrammarEquivalenceCheck(c *gin.Context) { // 判断所给的两个正则文法是否等价
	type Req struct {
		g1 model.Grammar
		g2 model.Grammar
	}
	var req Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析失败", "error": err.Error()})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&req.g1)
	normalizeGrammar(&req.g2)
	// 如果文法结构不完整，先完善结构
	req.g1 = *grammar_s.CompleteGrammarStructure(&req.g1)
	req.g2 = *grammar_s.CompleteGrammarStructure(&req.g2)

	// 先判断两个文法是否有效，且为正则文法
	if !grammar_s.IsValidGrammar(&req.g1) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid grammar1", "isEquivalent": false})
		return
	}
	if grammar_s.TypeDetermine(&req.g1) != 3 { // 非正则文法
		c.JSON(http.StatusBadRequest, gin.H{"msg": "grammar1 is not regular grammar", "isEquivalent": false})
		return
	}
	if !grammar_s.IsValidGrammar(&req.g2) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid grammar2", "isEquivalent": false})
		return
	}
	if grammar_s.TypeDetermine(&req.g2) != 3 { // 非正则文法
		c.JSON(http.StatusBadRequest, gin.H{"msg": "grammar2 is not regular grammar", "isEquivalent": false})
		return
	}

	// 再判断是否等价
	if !grammar_s.IsEquivalent(&req.g1, &req.g2) {
		c.JSON(http.StatusOK, gin.H{"msg": "this two grammars are not equivalent", "isEquivalent": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "this two grammars are equivalent", "isEquivalent": true})
}

func GrammarSimplify(c *gin.Context) { // 文法的化简——去无用符号（不可派生、不可达）、单一产生式、空产生式
	var grammar model.Grammar

	if err := c.ShouldBindJSON(&grammar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析失败", "error": err.Error()})
		return
	}

	// 先统一 ε 表示
	normalizeGrammar(&grammar)
	// 如果文法结构不完整，先完善结构
	grammar = *grammar_s.CompleteGrammarStructure(&grammar)

	// 可选：先校验文法有效性
	if !grammar_s.IsValidGrammar(&grammar) {
		c.JSON(http.StatusOK, gin.H{
			"type": -1,
			"msg":  "无效文法",
		})
		return
	}

	new_grammar := grammar_s.Simplify(&grammar)
	c.JSON(http.StatusOK, gin.H{"msg": "简化成功", "grammar": new_grammar.ToReactFlow()})
}

// 辅助函数：根据文法中定义的符号来分割字符串
// 辅助函数：根据文法中定义的符号来分割字符串
// 使用最长匹配原则，优先匹配较长的终结符
// 例如：如果文法中有 "if" 和 "i" 两个终结符，则 "if" 会被优先匹配为一个符号，而不是 "i" + "f"
func stringToSymbols(s string, grammar *model.Grammar) []model.Symbol {
	if s == "" || s == "ε" {
		return []model.Symbol{model.Epsilon}
	}

	// 收集文法中所有的终结符，并按长度排序（最长先匹配）
	var terminals []string

	// 从 Grammar.Terminals 获取定义的终结符
	for _, terminal := range grammar.Terminals {
		termStr := string(terminal)
		if termStr != "" && termStr != "ε" { // 排除空符号
			terminals = append(terminals, termStr)
		}
	}

	// 按长度降序排序，确保最长匹配
	for i := 0; i < len(terminals)-1; i++ {
		for j := i + 1; j < len(terminals); j++ {
			if len(terminals[i]) < len(terminals[j]) {
				terminals[i], terminals[j] = terminals[j], terminals[i]
			}
		}
	}

	// 按最长匹配原则分割字符串
	var symbols []model.Symbol
	i := 0
	for i < len(s) {
		matched := false
		// 从最长的符号开始尝试匹配
		for _, terminal := range terminals {
			if i+len(terminal) <= len(s) && s[i:i+len(terminal)] == terminal {
				symbols = append(symbols, model.Symbol(terminal))
				i += len(terminal)
				matched = true
				break
			}
		}
		// 如果没有匹配到任何已定义的终结符，按单字符处理
		if !matched {
			symbols = append(symbols, model.Symbol(string(s[i])))
			i++
		}
	}

	return symbols
}
