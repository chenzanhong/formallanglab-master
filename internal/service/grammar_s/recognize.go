// grammar_s/recognize.go - 文法识别和分析核心实现
// 包含多种文法分析算法：LL(1)、LR(0)、递归下降等

package grammar_s

/*
===================================================================================
                                文件结构说明
===================================================================================

本文件包含以下主要功能模块：

1. 【FIRST/FOLLOW集计算】(第30-200行)
   - CalculateFirst()          计算FIRST集
   - CalculateFirstCached()    带缓存的FIRST集计算
   - CalculateFollow()         计算FOLLOW集
   - CalculateFollowCached()   带缓存的FOLLOW集计算
   - PrintFirst(), PrintFollow() 打印FIRST/FOLLOW集

2. 【工具函数】(第200-300行)
   - symbolsToStringJoinSepFast()    符号序列转字符串
   - copySymbols()            符号序列复制
   - symbolsEqual()           符号序列比较
   - canDeriveEpsilon()       判断是否可推导出ε

3. 【字符串识别与BFS推导】(第300-500行)
   - RecognizeString()        BFS模拟推导识别字符串
   - BFSParseDetailed()       详细BFS推导分析

4. 【LL(1)分析算法】(第500-800行)
   - IsLL1()                  判断是否为LL(1)文法
   - LL1ParseDetailed()       详细LL(1)分析
   - LL1Parse()               简单LL(1)分析
   - LL1ParseWithRecovery()   带错误恢复的LL(1)分析

5. 【分析模式统一接口】(第800-900行)
   - ParseStringWithMode()    根据模式选择分析器

6. 【LR分析数据结构定义】(第900-1000行)
   - LR0Item, LR1Item         LR项目结构
   - LRState, LRAction        LR状态和动作结构
   - LRTable                  LR分析表结构

7. 【LR(0)分析算法】(第1000-1600行)
   - IsLR0Grammar()           判断是否为LR(0)文法
   - augmentGrammar()         增广文法
   - BuildLR0Table()          构造LR(0)分析表
   - constructLR0ItemSets()   构造LR(0)项目集族
   - closure()                计算项目集闭包
   - gotoItemSet()            计算GOTO函数
   - LR0ParseDetailed()       详细LR(0)分析

8. 【其他LR分析算法占位】(第1600-1800行)
   - IsLR1Grammar()           LR(1)文法判断
   - LR1ParseDetailed()       LR(1)分析

9. 【递归下降分析算法】(第1800-2206行)
   - RecursiveDescentParser   递归下降分析器结构
   - canUseRecursiveDescent() 判断是否适用递归下降
   - hasLeftRecursion()       检查左递归
   - eliminateLeftRecursion() 消除左递归
   - needsLeftFactoring()     检查是否需要左公因子提取
   - leftFactor()             提取左公因子
   - RecursiveDescentParse()  递归下降分析主函数
   - parseSymbol()            解析符号
   - parseNonTerminal()       解析非终结符

===================================================================================
*/

import (
	"backend/internal/domain/model"
	"fmt"
	"strings"
	"sync"
)

// ===================================================================================
//                          1. FIRST/FOLLOW集计算模块
// ===================================================================================

// 性能优化：FIRST/FOLLOW集缓存
var (
	firstSetCache  = make(map[string]map[model.Symbol]map[model.Symbol]struct{})
	followSetCache = make(map[string]map[model.Symbol]map[model.Symbol]struct{})
	cacheMutex     sync.RWMutex
)

// 生成文法的缓存键
func grammarCacheKey(grammar *model.Grammar) string {
	var key strings.Builder
	key.WriteString(string(grammar.StartSymbol))
	key.WriteRune('|')

	for _, prod := range grammar.Productions {
		for _, left := range prod.Left {
			key.WriteString(string(left))
		}
		key.WriteString("->")
		for _, right := range prod.Right {
			key.WriteString(string(right))
		}
		key.WriteRune(';')
	}
	return key.String()
}

// 高性能符号切片复制
func copySymbols(src []model.Symbol) []model.Symbol {
	if src == nil {
		return nil
	}
	dst := make([]model.Symbol, len(src))
	copy(dst, src)
	return dst
}

// 高性能字符串转换
func symbolsToStringJoinSepFast(symbols []model.Symbol) string {
	if len(symbols) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(symbols) * 2) // 预分配容量

	for i, sym := range symbols {
		if i > 0 {
			builder.WriteRune(' ')
		}
		if sym == model.Epsilon {
			builder.WriteString("ε")
		} else {
			builder.WriteString(string(sym))
		}
	}
	return builder.String()
}

// CalculateFirstCached 计算每个符号的 FIRST 集（带缓存）
func CalculateFirstCached(grammar *model.Grammar) map[model.Symbol]map[model.Symbol]struct{} {
	cacheKey := grammarCacheKey(grammar)

	cacheMutex.RLock()
	if cached, exists := firstSetCache[cacheKey]; exists {
		cacheMutex.RUnlock()
		return cached
	}
	cacheMutex.RUnlock()

	// 计算FIRST集
	firstSet := CalculateFirst(grammar)

	// 存入缓存
	cacheMutex.Lock()
	firstSetCache[cacheKey] = firstSet
	cacheMutex.Unlock()

	return firstSet
}

// CalculateFirst 计算每个符号的 FIRST 集
func CalculateFirst(grammar *model.Grammar) map[model.Symbol]map[model.Symbol]struct{} {
	firstSet := make(map[model.Symbol]map[model.Symbol]struct{})

	// 初始化所有终结符的 FIRST 集
	for _, terminal := range grammar.Terminals {
		firstSet[terminal] = map[model.Symbol]struct{}{terminal: {}}
	}

	// 初始化所有非终结符的 FIRST 集为空
	for _, nonTerminal := range grammar.NonTerminals {
		firstSet[nonTerminal] = make(map[model.Symbol]struct{})
	}

	changed := true
	for changed {
		changed = false
		for _, production := range grammar.Productions {
			left := production.Left[0]
			right := production.Right

			// 处理右部符号序列
			for _, symbol := range right {
				if _, exists := firstSet[symbol]; !exists {
					continue
				}
				for t := range firstSet[symbol] {
					if t != model.Epsilon {
						if _, in := firstSet[left][t]; !in {
							firstSet[left][t] = struct{}{}
							changed = true
						}
					}
				}
				// 如果 symbol 不能推出 ε，则停止
				if _, hasEpsilon := firstSet[symbol][model.Epsilon]; !hasEpsilon {
					break
				}
			}

			// 检查整个右部是否都能推出 ε
			allCanDeriveEpsilon := true
			for _, symbol := range right {
				if _, can := firstSet[symbol][model.Epsilon]; !can {
					allCanDeriveEpsilon = false
					break
				}
			}
			if allCanDeriveEpsilon {
				if _, has := firstSet[left][model.Epsilon]; !has {
					firstSet[left][model.Epsilon] = struct{}{}
					changed = true
				}
			}
		}
	}

	return firstSet
}

// PrintFirst 打印 FIRST 集
func PrintFirst(firstSet map[model.Symbol]map[model.Symbol]struct{}) {
	for symbol, first := range firstSet {
		items := []string{}
		for s := range first {
			if s == model.Epsilon {
				items = append(items, "ε")
			} else {
				items = append(items, string(s))
			}
		}
		fmt.Printf("FIRST(%s) = {%s}\n", symbol, strings.Join(items, ", "))
	}
}

// CalculateFollowCached 计算 FOLLOW 集（带缓存）
func CalculateFollowCached(grammar *model.Grammar, firstSet map[model.Symbol]map[model.Symbol]struct{}) map[model.Symbol]map[model.Symbol]struct{} {
	cacheKey := grammarCacheKey(grammar) + "_follow"

	cacheMutex.RLock()
	if cached, exists := followSetCache[cacheKey]; exists {
		cacheMutex.RUnlock()
		return cached
	}
	cacheMutex.RUnlock()

	// 计算FOLLOW集
	followSet := CalculateFollow(grammar, firstSet)

	// 存入缓存
	cacheMutex.Lock()
	followSetCache[cacheKey] = followSet
	cacheMutex.Unlock()

	return followSet
}

// CalculateFollow 计算 FOLLOW 集
func CalculateFollow(grammar *model.Grammar, firstSet map[model.Symbol]map[model.Symbol]struct{}) map[model.Symbol]map[model.Symbol]struct{} {
	followSet := make(map[model.Symbol]map[model.Symbol]struct{})

	// 初始化
	for _, nonTerminal := range grammar.NonTerminals {
		followSet[nonTerminal] = make(map[model.Symbol]struct{})
	}

	// 起始符号的 FOLLOW 包含 #
	followSet[grammar.StartSymbol]["#"] = struct{}{}

	changed := true
	for changed {
		changed = false
		for _, production := range grammar.Productions {
			left := production.Left[0]
			right := production.Right

			for i := 0; i < len(right); i++ {
				symbol := right[i]
				if !grammar.CheckIsNonTerminal(symbol) {
					continue
				}

				// 后继符号：right[i+1:]
				if i < len(right)-1 {
					nextSymbol := right[i+1]
					for t := range firstSet[nextSymbol] {
						if t != model.Epsilon {
							if _, in := followSet[symbol][t]; !in {
								followSet[symbol][t] = struct{}{}
								changed = true
							}
						}
					}
					if _, hasEpsilon := firstSet[nextSymbol][model.Epsilon]; hasEpsilon {
						for t := range followSet[left] {
							if _, in := followSet[symbol][t]; !in {
								followSet[symbol][t] = struct{}{}
								changed = true
							}
						}
					}
				} else {
					// 是最后一个符号
					for t := range followSet[left] {
						if _, in := followSet[symbol][t]; !in {
							followSet[symbol][t] = struct{}{}
							changed = true
						}
					}
				}
			}
		}
	}

	return followSet
}

// PrintFollow 打印 FOLLOW 集
func PrintFollow(followSet map[model.Symbol]map[model.Symbol]struct{}) {
	for symbol, follow := range followSet {
		items := []string{}
		for s := range follow {
			items = append(items, string(s))
		}
		fmt.Printf("FOLLOW(%s) = {%s}\n", symbol, strings.Join(items, ", "))
	}
}

// ===================================================================================
//                              2. 工具函数模块
// ===================================================================================

// symbolsEqual 比较两个符号切片是否相等
func symbolsEqual(a, b []model.Symbol) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// canDeriveEpsilon 判断非终结符是否可推出 ε
func canDeriveEpsilon(grammar *model.Grammar) bool {
	firstSet := CalculateFirst(grammar)
	start := grammar.StartSymbol
	_, can := firstSet[start][model.Epsilon]
	return can
}

// ===================================================================================
//                         3. 字符串识别与BFS推导模块
// ===================================================================================

// RecognizeString 判断字符串是否被文法生成（BFS 模拟推导）
func RecognizeString(g *model.Grammar, input []model.Symbol, maxSteps, maxWidth int) (bool, []string, error) {
	// 支持两种空串表示：[] 和 [ε]
	if len(input) == 0 || (len(input) == 1 && input[0] == model.Epsilon) {
		// 处理空串
		if canDeriveEpsilon(g) {
			for _, prod := range g.Productions {
				if prod.Left[0] == g.StartSymbol && len(prod.Right) == 1 && prod.Right[0] == model.Epsilon {
					return true, []string{fmt.Sprintf("%s → ε", string(g.StartSymbol))}, nil
				}
			}
			return true, []string{fmt.Sprintf("%s → ε", string(g.StartSymbol))}, nil
		}
		return false, nil, nil
	}

	// BFS 队列：记录当前符号串和推导路径
	type state struct {
		symbols []model.Symbol
		steps   int
		path    []string // 推导历史，如 "S → aA"
	}

	queue := []state{
		{symbols: []model.Symbol{g.StartSymbol}, steps: 0, path: []string{string(g.StartSymbol)}},
	}

	visited := make(map[string]bool)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// 去重
		key := symbolsToStringJoinSep(curr.symbols)
		if visited[key] || curr.steps >= maxSteps {
			continue
		}
		visited[key] = true

		// 检查是否匹配输入
		if symbolsEqual(curr.symbols, input) {
			return true, curr.path, nil
		}

		// 剪枝：长度超过输入太多
		if len(curr.symbols) > len(input)+10 {
			continue
		}

		// 尝试应用每个产生式
		for _, prod := range g.Productions {
			leftSymbol := prod.Left[0]
			for i := 0; i < len(curr.symbols); i++ {
				if curr.symbols[i] == leftSymbol {
					// 替换第 i 个符号（若产生式为 ε，则不添加任何符号）
					newSymbols := append([]model.Symbol{}, curr.symbols[:i]...)
					if !(len(prod.Right) == 0 || prod.Right[0] == model.Epsilon) {
						for _, r := range prod.Right {
							if r != model.Epsilon {
								newSymbols = append(newSymbols, r)
							}
						}
					}
					newSymbols = append(newSymbols, curr.symbols[i+1:]...)

					if len(newSymbols) > maxWidth {
						continue
					}

					// 构建新推导路径
					step := fmt.Sprintf("%s → %s", string(leftSymbol), symbolsToStringJoinSep(prod.Right))
					newPath := append([]string{}, curr.path...)
					newPath = append(newPath, step)

					queue = append(queue, state{
						symbols: newSymbols,
						steps:   curr.steps + 1,
						path:    newPath,
					})
				}
			}
		}
	}

	return false, nil, nil
}

// ===================================================================================
//                            4. LL(1)分析算法模块
// ===================================================================================

// IsLL1 判断文法是否为 LL(1)
func IsLL1(grammar *model.Grammar) (bool, string) {
	firstSet := CalculateFirst(grammar)
	followSet := CalculateFollow(grammar, firstSet)

	table := make(map[model.Symbol]map[model.Symbol]int)

	for _, nonTerminal := range grammar.NonTerminals {
		table[nonTerminal] = make(map[model.Symbol]int)
	}

	for i, prod := range grammar.Productions {
		left := prod.Left[0]
		right := prod.Right

		var firstOfRight map[model.Symbol]struct{}
		if len(right) > 0 {
			firstOfRight = firstSet[right[0]]
		} else {
			firstOfRight = map[model.Symbol]struct{}{model.Epsilon: {}}
		}

		for t := range firstOfRight {
			if t != model.Epsilon {
				if other, conflict := table[left][t]; conflict {
					return false, fmt.Sprintf("冲突: 产生式 %d 和 %d 在 %s/%s", other, i, string(left), string(t))
				}
				table[left][t] = i
			}
		}

		if _, hasEpsilon := firstOfRight[model.Epsilon]; hasEpsilon {
			for t := range followSet[left] {
				if other, conflict := table[left][t]; conflict {
					return false, fmt.Sprintf("FOLLOW 冲突: 产生式 %d 和 %d 在 %s/%s", other, i, string(left), string(t))
				}
				table[left][t] = i
			}
		}
	}

	return true, ""
}

// LL1ParseDetailed 使用 LL(1) 分析表解析输入串，返回详细的分析结果
func LL1ParseDetailed(grammar *model.Grammar, input []model.Symbol) *model.ParseResult {
	result := &model.ParseResult{
		Method: "LL(1) 分析",
		Steps:  make([]model.ParseStep, 0, len(input)*2+10), // 性能优化：预分配
	}

	// 检查是否为 LL(1) 文法
	if isLL1, msg := IsLL1(grammar); !isLL1 {
		result.Accepted = false
		result.Error = fmt.Sprintf("文法不是 LL(1) 文法: %s", msg)
		return result
	}

	firstSet := CalculateFirstCached(grammar)
	followSet := CalculateFollowCached(grammar, firstSet)

	// 构建预测分析表
	table := make(map[model.Symbol]map[model.Symbol]int)
	for _, nt := range grammar.NonTerminals {
		table[nt] = make(map[model.Symbol]int)
	}

	for i, prod := range grammar.Productions {
		left := prod.Left[0]
		right := prod.Right

		firstOfRight := map[model.Symbol]struct{}{}
		if len(right) > 0 {
			for k := range firstSet[right[0]] {
				firstOfRight[k] = struct{}{}
			}
		} else {
			firstOfRight[model.Epsilon] = struct{}{}
		}

		for t := range firstOfRight {
			if t != model.Epsilon {
				if _, exists := table[left][t]; exists {
					result.Accepted = false
					result.Error = fmt.Sprintf("LL(1) 冲突: %s → ... 和另一产生式在 %s", string(left), string(t))
					return result
				}
				table[left][t] = i
			}
		}

		if _, hasEpsilon := firstOfRight[model.Epsilon]; hasEpsilon {
			for t := range followSet[left] {
				if _, exists := table[left][t]; exists {
					result.Accepted = false
					result.Error = fmt.Sprintf("LL(1) FOLLOW 冲突: %s → ... 和另一产生式在 %s", string(left), string(t))
					return result
				}
				table[left][t] = i
			}
		}
	}

	// 初始化分析栈和输入
	stack := []model.Symbol{"#", grammar.StartSymbol}
	inputWithEnd := append(input, "#")
	ip := 0

	// 记录初始状态
	result.Steps = append(result.Steps, model.ParseStep{
		StepType:    "init",
		Description: "初始化分析栈和输入",
		Stack:       copySymbols(stack),
		Input:       copySymbols(inputWithEnd),
		InputPos:    ip,
		Action:      "分析开始",
	})

	for len(stack) > 1 { // 栈中至少有 "#" 和另一个符号
		if ip >= len(inputWithEnd) {
			result.Accepted = false
			result.Error = "输入串太短"
			return result
		}

		top := stack[len(stack)-1]
		current := inputWithEnd[ip]

		step := model.ParseStep{
			Stack:    append([]model.Symbol{}, stack...),
			Input:    append([]model.Symbol{}, inputWithEnd[ip:]...),
			InputPos: ip,
		}

		// 如果栈顶符号与当前输入符号相同
		if top == current {
			stack = stack[:len(stack)-1] // 出栈
			ip++                         // 移动输入指针

			step.StepType = "match"
			step.Description = fmt.Sprintf("匹配终结符 %s", string(top))
			step.Action = fmt.Sprintf("出栈 %s，移动输入指针", string(top))

			result.Steps = append(result.Steps, step)
			continue
		}

		// 如果栈顶是终结符
		if grammar.CheckIsTerminal(top) {
			step.StepType = "error"
			step.Description = fmt.Sprintf("匹配失败: 栈顶 %s ≠ 输入 %s", string(top), string(current))
			step.Action = "分析失败"
			result.Steps = append(result.Steps, step)

			result.Accepted = false
			result.Error = step.Description
			return result
		}

		// 栈顶是非终结符，查找预测分析表
		if productionIndex, ok := table[top][current]; ok {
			prod := grammar.Productions[productionIndex]

			// 出栈非终结符
			stack = stack[:len(stack)-1]

			// 反向压栈产生式右部（除了 ε）
			for i := len(prod.Right) - 1; i >= 0; i-- {
				if prod.Right[i] != model.Epsilon {
					stack = append(stack, prod.Right[i])
				}
			}

			step.StepType = "predict"
			step.Description = fmt.Sprintf("使用产生式: %s → %s", string(prod.Left[0]), symbolsToStringJoinSep(prod.Right))
			step.Action = fmt.Sprintf("出栈 %s，压入 %s", string(top), symbolsToStringJoinSep(prod.Right))
			step.Production = &prod

			result.Steps = append(result.Steps, step)
		} else {
			step.StepType = "error"
			step.Description = fmt.Sprintf("无匹配产生式: %s / %s", string(top), string(current))
			step.Action = "分析失败"
			result.Steps = append(result.Steps, step)

			result.Accepted = false
			result.Error = step.Description
			return result
		}
	}

	// 检查是否成功接受
	if len(stack) == 1 && stack[0] == "#" && ip < len(inputWithEnd) && inputWithEnd[ip] == "#" {
		result.Steps = append(result.Steps, model.ParseStep{
			StepType:    "accept",
			Description: "分析成功: 输入被接受",
			Stack:       []model.Symbol{"#"},
			Input:       []model.Symbol{"#"},
			InputPos:    ip,
			Action:      "接受输入串",
		})

		result.Accepted = true
		result.Message = "输入串被 LL(1) 文法接受"
	} else {
		result.Accepted = false
		result.Error = "分析结束但状态不正确"
	}

	return result
}

// LL1Parse 使用 LL(1) 分析表解析输入串（简单版本）
func LL1Parse(grammar *model.Grammar, input []model.Symbol) (bool, []string, error) {
	result := LL1ParseDetailed(grammar, input)

	if result.Error != "" {
		return false, nil, fmt.Errorf(result.Error)
	}

	// 转换成简单的字符串描述
	var steps []string
	for _, step := range result.Steps {
		steps = append(steps, step.Description)
	}

	return result.Accepted, steps, nil
}

// LL1ParseWithRecovery LL(1)分析器带错误恢复机制
func LL1ParseWithRecovery(grammar *model.Grammar, input []model.Symbol) *model.ParseResult {
	result := &model.ParseResult{
		Method: "LL(1) 分析（带错误恢复）",
		Steps:  make([]model.ParseStep, 0, len(input)*2+10),
	}

	// 检查是否为 LL(1) 文法
	if isLL1, msg := IsLL1(grammar); !isLL1 {
		result.Accepted = false
		result.Error = fmt.Sprintf("文法不是 LL(1) 文法: %s", msg)
		return result
	}

	firstSet := CalculateFirstCached(grammar)
	followSet := CalculateFollowCached(grammar, firstSet)

	// 构建预测分析表
	table := make(map[model.Symbol]map[model.Symbol]int)
	for _, nt := range grammar.NonTerminals {
		table[nt] = make(map[model.Symbol]int)
	}

	for i, prod := range grammar.Productions {
		left := prod.Left[0]
		right := prod.Right

		firstOfRight := map[model.Symbol]struct{}{}
		if len(right) > 0 {
			for k := range firstSet[right[0]] {
				firstOfRight[k] = struct{}{}
			}
		} else {
			firstOfRight[model.Epsilon] = struct{}{}
		}

		for t := range firstOfRight {
			if t != model.Epsilon {
				if _, exists := table[left][t]; exists {
					result.Accepted = false
					result.Error = fmt.Sprintf("LL(1) 冲突: %s → ... 和另一产生式在 %s", string(left), string(t))
					return result
				}
				table[left][t] = i
			}
		}

		if _, hasEpsilon := firstOfRight[model.Epsilon]; hasEpsilon {
			for t := range followSet[left] {
				if _, exists := table[left][t]; exists {
					result.Accepted = false
					result.Error = fmt.Sprintf("LL(1) FOLLOW 冲突: %s → ... 和另一产生式在 %s", string(left), string(t))
					return result
				}
				table[left][t] = i
			}
		}
	}

	// 初始化分析栈和输入
	stack := []model.Symbol{"#", grammar.StartSymbol}
	inputWithEnd := append(copySymbols(input), "#")
	ip := 0
	errorCount := 0
	maxErrors := 5 // 最多允许的错误数

	// 记录初始状态
	result.Steps = append(result.Steps, model.ParseStep{
		StepType:    "init",
		Description: "初始化分析栈和输入（带错误恢复）",
		Stack:       copySymbols(stack),
		Input:       copySymbols(inputWithEnd),
		InputPos:    ip,
		Action:      "分析开始",
	})

	for len(stack) > 1 && errorCount < maxErrors {
		if ip >= len(inputWithEnd) {
			result.Accepted = false
			result.Error = "输入串太短"
			return result
		}

		top := stack[len(stack)-1]
		current := inputWithEnd[ip]

		step := model.ParseStep{
			Stack:    copySymbols(stack),
			Input:    inputWithEnd[ip:],
			InputPos: ip,
		}

		// 如果栈顶符号与当前输入符号相同
		if top == current {
			stack = stack[:len(stack)-1]
			ip++

			step.StepType = "match"
			step.Description = fmt.Sprintf("匹配终结符 %s", string(top))
			step.Action = fmt.Sprintf("出栈 %s，移动输入指针", string(top))

			result.Steps = append(result.Steps, step)
			continue
		}

		// 如果栈顶是终结符
		if grammar.CheckIsTerminal(top) {
			// 错误恢复：跳过当前输入符号
			step.StepType = "error"
			step.Description = fmt.Sprintf("错误: 栈顶 %s ≠ 输入 %s，尝试恢复", string(top), string(current))
			step.Action = fmt.Sprintf("跳过输入符号 %s", string(current))
			result.Steps = append(result.Steps, step)

			ip++ // 跳过错误的输入符号
			errorCount++
			continue
		}

		// 栈顶是非终结符，查找预测分析表
		if productionIndex, ok := table[top][current]; ok {
			prod := grammar.Productions[productionIndex]

			// 出栈非终结符
			stack = stack[:len(stack)-1]

			// 反向压栈产生式右部（除了 ε）
			for i := len(prod.Right) - 1; i >= 0; i-- {
				if prod.Right[i] != model.Epsilon {
					stack = append(stack, prod.Right[i])
				}
			}

			step.StepType = "predict"
			step.Description = fmt.Sprintf("使用产生式: %s → %s", string(prod.Left[0]), symbolsToStringJoinSepFast(prod.Right))
			step.Action = fmt.Sprintf("出栈 %s，压入 %s", string(top), symbolsToStringJoinSepFast(prod.Right))
			step.Production = &prod

			result.Steps = append(result.Steps, step)
		} else {
			// 错误恢复：弹出栈顶非终结符
			step.StepType = "error"
			step.Description = fmt.Sprintf("错误: 无匹配产生式 %s / %s，尝试恢复", string(top), string(current))
			step.Action = fmt.Sprintf("弹出非终结符 %s", string(top))
			result.Steps = append(result.Steps, step)

			stack = stack[:len(stack)-1] // 弹出栈顶非终结符
			errorCount++
		}
	}

	// 检查是否成功接受
	if len(stack) == 1 && stack[0] == "#" && ip < len(inputWithEnd) && inputWithEnd[ip] == "#" {
		result.Steps = append(result.Steps, model.ParseStep{
			StepType:    "accept",
			Description: fmt.Sprintf("分析成功: 输入被接受（含 %d 个错误）", errorCount),
			Stack:       []model.Symbol{"#"},
			Input:       []model.Symbol{"#"},
			InputPos:    ip,
			Action:      "接受输入串",
		})

		result.Accepted = true
		if errorCount > 0 {
			result.Message = fmt.Sprintf("输入串被接受，但有 %d 个错误被恢复", errorCount)
		} else {
			result.Message = "输入串被 LL(1) 文法完美接受"
		}
	} else {
		result.Accepted = false
		if errorCount >= maxErrors {
			result.Error = fmt.Sprintf("错误太多（%d），停止分析", errorCount)
		} else {
			result.Error = "分析结束但状态不正确"
		}
	}

	return result
}

// BFSParseDetailed BFS推导详细分析版本
func BFSParseDetailed(g *model.Grammar, input []model.Symbol, maxSteps, maxWidth int) *model.ParseResult {
	result := &model.ParseResult{
		Method: "BFS 推导模拟",
		Steps:  []model.ParseStep{},
	}

	if len(input) == 0 {
		// 处理空串
		if canDeriveEpsilon(g) {
			for _, prod := range g.Productions {
				if prod.Left[0] == g.StartSymbol && len(prod.Right) == 1 && prod.Right[0] == model.Epsilon {
					result.Accepted = true
					result.Steps = append(result.Steps, model.ParseStep{
						StepType:    "accept",
						Description: fmt.Sprintf("空串被接受: %s → ε", string(g.StartSymbol)),
						Stack:       []model.Symbol{g.StartSymbol},
						Input:       []model.Symbol{model.Epsilon},
						InputPos:    0,
						Action:      "接受空串",
						Production:  &prod,
					})
					result.Message = "空串被文法接受"
					return result
				}
			}
			result.Accepted = true
			result.Message = "空串被文法接受"
			return result
		}
		result.Accepted = false
		result.Error = "文法不能推导出空串"
		return result
	}

	// BFS 队列：记录当前符号串和推导路径
	type state struct {
		symbols []model.Symbol
		steps   int
		path    []string // 推导历史，如 "S → aA"
	}

	queue := []state{
		{symbols: []model.Symbol{g.StartSymbol}, steps: 0, path: []string{string(g.StartSymbol)}},
	}

	visited := make(map[string]bool)
	stepCount := 0

	// 记录初始状态
	result.Steps = append(result.Steps, model.ParseStep{
		StepType:    "init",
		Description: "初始化BFS推导，从起始符号开始",
		Stack:       []model.Symbol{g.StartSymbol},
		Input:       input,
		InputPos:    0,
		Action:      "开始推导",
	})

	for len(queue) > 0 && stepCount < maxSteps {
		curr := queue[0]
		queue = queue[1:]
		stepCount++

		// 去重
		key := symbolsToStringJoinSep(curr.symbols)
		if visited[key] || curr.steps >= maxSteps {
			continue
		}
		visited[key] = true

		// 检查是否匹配输入
		if symbolsEqual(curr.symbols, input) {
			result.Steps = append(result.Steps, model.ParseStep{
				StepType:    "accept",
				Description: "找到匹配的推导序列",
				Stack:       curr.symbols,
				Input:       input,
				InputPos:    len(input),
				Action:      "接受输入串",
			})

			result.Accepted = true
			result.Message = fmt.Sprintf("输入串通过 %d 步推导被接受", curr.steps)
			return result
		}

		// 剪枝：长度超过输入太多
		if len(curr.symbols) > len(input)+10 {
			continue
		}

		// 尝试应用每个产生式
		for _, prod := range g.Productions {
			leftSymbol := prod.Left[0]
			for i := 0; i < len(curr.symbols); i++ {
				if curr.symbols[i] == leftSymbol {
					// 替换第 i 个符号
					newSymbols := append([]model.Symbol{}, curr.symbols[:i]...)
					newSymbols = append(newSymbols, prod.Right...)
					newSymbols = append(newSymbols, curr.symbols[i+1:]...)

					if len(newSymbols) > maxWidth {
						continue
					}

					// 构建新推导路径
					step := fmt.Sprintf("%s → %s", string(leftSymbol), symbolsToStringJoinSep(prod.Right))
					newPath := append([]string{}, curr.path...)
					newPath = append(newPath, step)

					// 记录推导步骤
					result.Steps = append(result.Steps, model.ParseStep{
						StepType:    "predict",
						Description: fmt.Sprintf("应用产生式: %s → %s", string(leftSymbol), symbolsToStringJoinSep(prod.Right)),
						Stack:       newSymbols,
						Input:       input,
						InputPos:    0,
						Action:      fmt.Sprintf("替换 %s 为 %s", string(leftSymbol), symbolsToStringJoinSep(prod.Right)),
						Production:  &prod,
					})

					queue = append(queue, state{
						symbols: newSymbols,
						steps:   curr.steps + 1,
						path:    newPath,
					})
				}
			}
		}
	}

	result.Accepted = false
	result.Error = fmt.Sprintf("在 %d 步内未找到匹配的推导", maxSteps)
	return result
}

// ===================================================================================
//                         5. 分析模式统一接口模块
// ===================================================================================

// ParseStringWithMode 根据指定模式分析输入串
func ParseStringWithMode(grammar *model.Grammar, input []model.Symbol, mode string, showSteps bool) *model.ParseResult {
	switch mode {
	case "ll1":
		result := LL1ParseDetailed(grammar, input)
		if !showSteps {
			result.Steps = nil
		}
		return result

	case "ll1_recovery":
		result := LL1ParseWithRecovery(grammar, input)
		if !showSteps {
			result.Steps = nil
		}
		return result

	case "recursive_descent":
		result := RecursiveDescentParse(grammar, input)
		if !showSteps {
			result.Steps = nil
		}
		return result

	case "lr0":
		result := LR0ParseDetailed(grammar, input)
		if !showSteps {
			result.Steps = nil
		}
		return result

	case "lr1":
		result := LR1ParseDetailed(grammar, input)
		if !showSteps {
			result.Steps = nil
		}
		return result

	case "bfs":
		result := BFSParseDetailed(grammar, input, 100, 100)
		if !showSteps {
			result.Steps = nil
		}
		return result

	case "auto":
		fallthrough // 忽略下一个 case 的条件判断，直接执行它的代码块
	default:
		// 自动选择最适合的分析方法

		// 先尝试 LL(1)
		if isLL1, _ := IsLL1(grammar); isLL1 {
			result := LL1ParseDetailed(grammar, input)
			if !showSteps {
				result.Steps = nil
			}
			return result
		}

		// 再尝试 LR(1)
		if isLR1, _ := IsLR1Grammar(grammar); isLR1 {
			result := LR1ParseDetailed(grammar, input)
			if !showSteps {
				result.Steps = nil
			}
			return result
		}

		// 尝试递归下降（适用于特定结构的文法）
		if canUseRecursiveDescent(grammar) {
			result := RecursiveDescentParse(grammar, input)
			if !showSteps {
				result.Steps = nil
			}
			return result
		}

		// 默认使用 BFS 推导
		result := BFSParseDetailed(grammar, input, 100, 100)
		if !showSteps {
			result.Steps = nil
		}
		return result
	}
}

// ===================================================================================
//                         6. LR分析数据结构定义模块
// ===================================================================================

// LR0Item 表示LR(0)项目
type LR0Item struct {
	Production *model.Production // 对应的产生式
	DotPos     int               // 点的位置
}

// LR1Item 表示LR(1)项目
type LR1Item struct {
	Production *model.Production // 对应的产生式
	DotPos     int               // 点的位置
	Lookahead  model.Symbol      // 向前看符号
}

// LRState 表示LR状态
type LRState struct {
	ID     int       // 状态ID
	Items  []LR0Item // LR(0)项目集
	Items1 []LR1Item // LR(1)项目集（仅LR(1)使用）
}

// LRAction 表示LR分析表中的动作
type LRAction struct {
	Type       string            // "shift", "reduce", "accept", "error"
	Target     int               // 目标状态号（shift时）或产生式号（reduce时）
	Production *model.Production // 归约使用的产生式
}

// LRTable 表示LR分析表
type LRTable struct {
	Action map[int]map[model.Symbol]*LRAction // ACTION表 [state][terminal] -> action
	Goto   map[int]map[model.Symbol]int       // GOTO表 [state][nonterminal] -> state
	States []LRState                          // 状态集合
}

// ===================================================================================
//                           7. LR(0)分析算法模块
// ===================================================================================

// IsLR0Grammar 判断文法是否为LR(0)文法
func IsLR0Grammar(grammar *model.Grammar) (bool, string) {
	// 增广文法
	augmentedGrammar := augmentGrammar(grammar)

	// 构造LR(0)项目集族
	table, err := BuildLR0Table(augmentedGrammar)
	if err != nil {
		return false, err.Error()
	}

	// 检查是否存在冲突
	for stateID, actions := range table.Action {
		for symbol, action := range actions {
			if action.Type == "conflict" {
				return false, fmt.Sprintf("状态 %d 在符号 %s 上存在冲突", stateID, string(symbol))
			}
		}
	}

	return true, ""
}

// augmentGrammar 增广文法
func augmentGrammar(grammar *model.Grammar) *model.Grammar {
	// 创建新的起始符号 S'
	newStartSymbol := model.Symbol(string(grammar.StartSymbol) + "'")

	// 复制原文法
	augmented := &model.Grammar{
		StartSymbol:  newStartSymbol,
		Terminals:    append([]model.Symbol{}, grammar.Terminals...),
		NonTerminals: append([]model.Symbol{newStartSymbol}, grammar.NonTerminals...),
		Productions:  append([]model.Production{}, grammar.Productions...),
	}

	// 添加新的起始产生式 S' -> S
	augmented.Productions = append([]model.Production{{
		Left:  []model.Symbol{newStartSymbol},
		Right: []model.Symbol{grammar.StartSymbol},
	}}, augmented.Productions...)

	return augmented
}

// BuildLR0Table 构造LR(0)分析表
func BuildLR0Table(grammar *model.Grammar) (*LRTable, error) {
	table := &LRTable{
		Action: make(map[int]map[model.Symbol]*LRAction),
		Goto:   make(map[int]map[model.Symbol]int),
		States: []LRState{},
	}

	// 构造LR(0)项目集族
	itemSets := constructLR0ItemSets(grammar)
	table.States = itemSets

	// 初始化ACTION和GOTO表
	for i := range itemSets {
		table.Action[i] = make(map[model.Symbol]*LRAction)
		table.Goto[i] = make(map[model.Symbol]int)
	}

	// 填充分析表
	for i, itemSet := range itemSets {
		for _, item := range itemSet.Items {
			prod := item.Production

			// 情况1：移进项目 A -> α • a β (a是终结符)
			if item.DotPos < len(prod.Right) {
				nextSymbol := prod.Right[item.DotPos]
				if grammar.CheckIsTerminal(nextSymbol) {
					// 查找GOTO(I, a)
					gotoState := findGotoState(itemSets, i, nextSymbol, grammar)
					if gotoState != -1 {
						// 检查冲突
						if _, exists := table.Action[i][nextSymbol]; exists {
							return nil, fmt.Errorf("移进-归约冲突：状态 %d，符号 %s", i, string(nextSymbol))
						}
						table.Action[i][nextSymbol] = &LRAction{
							Type:   "shift",
							Target: gotoState,
						}
					}
				} else if grammar.CheckIsNonTerminal(nextSymbol) {
					// GOTO项目
					gotoState := findGotoState(itemSets, i, nextSymbol, grammar)
					if gotoState != -1 {
						table.Goto[i][nextSymbol] = gotoState
					}
				}
			} else {
				// 情况2：归约项目 A -> α •
				if len(prod.Left) > 0 && prod.Left[0] == grammar.StartSymbol && len(prod.Right) == 1 {
					// 接受项目 S' -> S •
					table.Action[i]["#"] = &LRAction{
						Type: "accept",
					}
				} else {
					// 归约项目
					prodIndex := findProductionIndex(grammar, prod)
					for _, terminal := range grammar.Terminals {
						// 检查冲突
						if existing, exists := table.Action[i][terminal]; exists {
							if existing.Type == "shift" {
								return nil, fmt.Errorf("移进-归约冲突：状态 %d，符号 %s", i, string(terminal))
							} else if existing.Type == "reduce" {
								return nil, fmt.Errorf("归约-归约冲突：状态 %d，符号 %s", i, string(terminal))
							}
						}
						table.Action[i][terminal] = &LRAction{
							Type:       "reduce",
							Target:     prodIndex,
							Production: prod,
						}
					}
					// 也要考虑$符号
					if existing, exists := table.Action[i]["#"]; exists {
						if existing.Type != "accept" {
							return nil, fmt.Errorf("归约冲突：状态 %d，符号 #", i)
						}
					} else {
						table.Action[i]["#"] = &LRAction{
							Type:       "reduce",
							Target:     prodIndex,
							Production: prod,
						}
					}
				}
			}
		}
	}

	return table, nil
}

// constructLR0ItemSets 构造LR(0)项目集族
func constructLR0ItemSets(grammar *model.Grammar) []LRState {
	var states []LRState
	stateMap := make(map[string]int) // 项目集字符串 -> 状态ID

	// 初始项目集 I0
	startProd := &grammar.Productions[0] // S' -> S
	initialItem := LR0Item{
		Production: startProd,
		DotPos:     0,
	}

	initialItemSet := closure([]LR0Item{initialItem}, grammar)
	states = append(states, LRState{
		ID:    0,
		Items: initialItemSet,
	})
	stateMap[itemSetToString(initialItemSet)] = 0

	// 工作队列
	queue := []int{0}

	for len(queue) > 0 {
		currentStateID := queue[0]
		queue = queue[1:]
		currentState := states[currentStateID]

		// 收集可转移的符号
		symbols := make(map[model.Symbol]bool)
		for _, item := range currentState.Items {
			if item.DotPos < len(item.Production.Right) {
				nextSymbol := item.Production.Right[item.DotPos]
				symbols[nextSymbol] = true
			}
		}

		// 对每个符号计算GOTO
		for symbol := range symbols {
			newItemSet := gotoItemSet(currentState.Items, symbol, grammar)
			if len(newItemSet) == 0 {
				continue
			}

			itemSetStr := itemSetToString(newItemSet)
			if _, exists := stateMap[itemSetStr]; exists {
				// 状态已存在，只需要记录转移
				continue
			} else {
				// 新状态
				newStateID := len(states)
				states = append(states, LRState{
					ID:    newStateID,
					Items: newItemSet,
				})
				stateMap[itemSetStr] = newStateID
				queue = append(queue, newStateID)
			}
		}
	}

	return states
}

// closure 计算项目集的闭包
func closure(items []LR0Item, grammar *model.Grammar) []LR0Item {
	result := append([]LR0Item{}, items...)
	added := make(map[string]bool)

	// 初始化已添加的项目
	for _, item := range items {
		added[itemToString(item)] = true
	}

	changed := true
	for changed {
		changed = false
		for _, item := range result {
			// 如果点不在最后，且点后面是非终结符
			if item.DotPos < len(item.Production.Right) {
				nextSymbol := item.Production.Right[item.DotPos]
				if grammar.CheckIsNonTerminal(nextSymbol) {
					// 添加所有 B -> • γ 的项目
					for _, prod := range grammar.Productions {
						if len(prod.Left) > 0 && prod.Left[0] == nextSymbol {
							newItem := LR0Item{
								Production: &prod,
								DotPos:     0,
							}
							itemStr := itemToString(newItem)
							if !added[itemStr] {
								result = append(result, newItem)
								added[itemStr] = true
								changed = true
							}
						}
					}
				}
			}
		}
	}

	return result
}

// gotoItemSet 计算GOTO(I, X)
func gotoItemSet(items []LR0Item, symbol model.Symbol, grammar *model.Grammar) []LR0Item {
	var newItems []LR0Item

	for _, item := range items {
		// 找到所有形如 A -> α • X β 的项目
		if item.DotPos < len(item.Production.Right) && item.Production.Right[item.DotPos] == symbol {
			// 移动点的位置：A -> α X • β
			newItem := LR0Item{
				Production: item.Production,
				DotPos:     item.DotPos + 1,
			}
			newItems = append(newItems, newItem)
		}
	}

	if len(newItems) == 0 {
		return nil
	}

	// 返回闭包
	return closure(newItems, grammar)
}

// 辅助函数
func itemToString(item LR0Item) string {
	prod := item.Production
	var parts []string

	// 左部
	leftStr := ""
	for _, sym := range prod.Left {
		leftStr += string(sym)
	}
	parts = append(parts, leftStr, "->")

	// 右部（插入点）
	for i, sym := range prod.Right {
		if i == item.DotPos {
			parts = append(parts, "•")
		}
		if sym == model.Epsilon {
			parts = append(parts, "ε")
		} else {
			parts = append(parts, string(sym))
		}
	}
	if item.DotPos == len(prod.Right) {
		parts = append(parts, "•")
	}

	return strings.Join(parts, " ")
}

func itemSetToString(items []LR0Item) string {
	var strs []string
	for _, item := range items {
		strs = append(strs, itemToString(item))
	}
	return strings.Join(strs, "|")
}

func findGotoState(states []LRState, fromState int, symbol model.Symbol, grammar *model.Grammar) int {
	if fromState >= len(states) {
		return -1
	}

	gotoItems := gotoItemSet(states[fromState].Items, symbol, grammar)
	if len(gotoItems) == 0 {
		return -1
	}

	gotoStr := itemSetToString(gotoItems)
	for i, state := range states {
		if itemSetToString(state.Items) == gotoStr {
			return i
		}
	}

	return -1
}

func findProductionIndex(grammar *model.Grammar, target *model.Production) int {
	for i, prod := range grammar.Productions {
		if productionsEqual(&prod, target) {
			return i
		}
	}
	return -1
}

func productionsEqual(p1, p2 *model.Production) bool {
	if len(p1.Left) != len(p2.Left) || len(p1.Right) != len(p2.Right) {
		return false
	}

	for i, sym := range p1.Left {
		if sym != p2.Left[i] {
			return false
		}
	}

	for i, sym := range p1.Right {
		if sym != p2.Right[i] {
			return false
		}
	}

	return true
}

// ===================================================================================
//                        8. 其他LR分析算法占位模块
// ===================================================================================

// IsLR1Grammar 判断文法是否为LR(1)文法
func IsLR1Grammar(grammar *model.Grammar) (bool, string) {
	// TODO: 实现LR(1)文法判断逻辑
	// 1. 构造LR(1)项目集族
	// 2. 检查向前看符号是否能解决所有冲突
	return false, "LR(1)文法判断尚未实现"
}

// BuildLR1Table 构造LR(1)分析表
func BuildLR1Table(grammar *model.Grammar) (*LRTable, error) {
	// TODO: 实现LR(1)分析表构造
	// 1. 计算LR(1)项目集族
	// 2. 构造ACTION和GOTO表（包含向前看符号）
	return nil, fmt.Errorf("LR(1)分析表构造尚未实现")
}
func LR0ParseDetailed(grammar *model.Grammar, input []model.Symbol) *model.ParseResult {
	result := &model.ParseResult{
		Method: "LR(0) 分析",
		Steps:  []model.ParseStep{},
	}

	// 检查是否为LR(0)文法
	if isLR0, msg := IsLR0Grammar(grammar); !isLR0 {
		result.Accepted = false
		result.Error = fmt.Sprintf("文法不是LR(0)文法: %s", msg)
		return result
	}

	// 增广文法并构造分析表
	augmentedGrammar := augmentGrammar(grammar)
	table, err := BuildLR0Table(augmentedGrammar)
	if err != nil {
		result.Accepted = false
		result.Error = fmt.Sprintf("构造LR(0)分析表失败: %s", err.Error())
		return result
	}

	// 初始化分析栈和输入
	stateStack := []int{0}             // 状态栈，初始状态为0
	symbolStack := []model.Symbol{"#"} // 符号栈，底部为#
	inputWithEnd := append(copySymbols(input), "#")
	ip := 0 // 输入指针

	// 记录初始状态
	result.Steps = append(result.Steps, model.ParseStep{
		StepType:    "init",
		Description: "初始化LR(0)分析栈和输入",
		Stack:       copySymbols(symbolStack),
		Input:       copySymbols(inputWithEnd),
		InputPos:    ip,
		Action:      "分析开始，状态栈: [0]",
	})

	for {
		if ip >= len(inputWithEnd) {
			result.Accepted = false
			result.Error = "输入串意外结束"
			return result
		}

		currentState := stateStack[len(stateStack)-1]
		currentSymbol := inputWithEnd[ip]

		step := model.ParseStep{
			Stack:    append([]model.Symbol{}, symbolStack...),
			Input:    inputWithEnd[ip:],
			InputPos: ip,
		}

		// 查找ACTION表
		action, exists := table.Action[currentState][currentSymbol]
		if !exists {
			step.StepType = "error"
			step.Description = fmt.Sprintf("在状态 %d 无法处理符号 %s", currentState, string(currentSymbol))
			step.Action = "分析失败"
			result.Steps = append(result.Steps, step)

			result.Accepted = false
			result.Error = step.Description
			return result
		}

		switch action.Type {
		case "shift":
			// 移进动作
			symbolStack = append(symbolStack, currentSymbol)
			stateStack = append(stateStack, action.Target)
			ip++

			step.StepType = "match"
			step.Description = fmt.Sprintf("移进: 将 %s 压入栈，转到状态 %d", string(currentSymbol), action.Target)
			step.Action = fmt.Sprintf("移进 %s，状态栈: %v", string(currentSymbol), stateStack)
			result.Steps = append(result.Steps, step)

		case "reduce":
			// 归约动作
			prod := action.Production
			popCount := len(prod.Right)

			// 处理ε产生式
			if popCount == 1 && len(prod.Right) > 0 && prod.Right[0] == model.Epsilon {
				popCount = 0
			}

			// 弹出符号和状态
			if popCount > 0 {
				symbolStack = symbolStack[:len(symbolStack)-popCount]
				stateStack = stateStack[:len(stateStack)-popCount]
			}

			// 压入归约后的非终结符
			leftSymbol := prod.Left[0]
			symbolStack = append(symbolStack, leftSymbol)

			// 查找GOTO表
			newState := stateStack[len(stateStack)-1]
			gotoState, gotoExists := table.Goto[newState][leftSymbol]
			if !gotoExists {
				step.StepType = "error"
				step.Description = fmt.Sprintf("GOTO(%d, %s) 未定义", newState, string(leftSymbol))
				step.Action = "分析失败"
				result.Steps = append(result.Steps, step)

				result.Accepted = false
				result.Error = step.Description
				return result
			}

			stateStack = append(stateStack, gotoState)

			step.StepType = "predict"
			step.Description = fmt.Sprintf("归约: 使用产生式 %s → %s", string(leftSymbol), symbolsToStringJoinSepFast(prod.Right))
			step.Action = fmt.Sprintf("归约，弹出 %d 个符号，转到状态 %d", popCount, gotoState)
			step.Production = prod
			result.Steps = append(result.Steps, step)

		case "accept":
			// 接受动作
			step.StepType = "accept"
			step.Description = "LR(0)分析成功，接受输入串"
			step.Action = "接受输入串"
			result.Steps = append(result.Steps, step)

			result.Accepted = true
			result.Message = "输入串被LR(0)分析器接受"
			return result

		default:
			step.StepType = "error"
			step.Description = fmt.Sprintf("未知的动作类型: %s", action.Type)
			step.Action = "分析失败"
			result.Steps = append(result.Steps, step)

			result.Accepted = false
			result.Error = step.Description
			return result
		}
	}
}

// LR0ParseDetailed LR(0)分析器
func LR1ParseDetailed(grammar *model.Grammar, input []model.Symbol) *model.ParseResult {
	result := &model.ParseResult{
		Method: "LR(1) 分析",
		Steps:  []model.ParseStep{},
	}

	// 检查是否为LR(1)文法
	if isLR1, msg := IsLR1Grammar(grammar); !isLR1 {
		result.Accepted = false
		result.Error = fmt.Sprintf("文法不是LR(1)文法: %s", msg)
		return result
	}

	// TODO: 实现LR(1)分析逻辑
	// 1. 构造LR(1)分析表
	// 2. 使用栈进行分析
	// 3. 记录详细步骤
	result.Accepted = false
	result.Error = "LR(1)分析器尚未实现"
	return result
} // ===================================================================================
//                         9. 递归下降分析算法模块
// ===================================================================================

// RecursiveDescentParser 递归下降分析器结构
type RecursiveDescentParser struct {
	Grammar *model.Grammar
	Input   []model.Symbol
	Pos     int
	Steps   []model.ParseStep
}

// canUseRecursiveDescent 判断文法是否适用于递归下降分析
func canUseRecursiveDescent(grammar *model.Grammar) bool {
	// 检查是否有左递归
	if hasLeftRecursion(grammar) {
		return false
	}

	// 检查是否需要提取左公因子
	if needsLeftFactoring(grammar) {
		return false
	}

	// 检查是否为LL(1)
	if isLL1, _ := IsLL1(grammar); isLL1 {
		return true
	}

	return false
}

// needsLeftFactoring 检查是否需要左公因子提取
func needsLeftFactoring(grammar *model.Grammar) bool {
	// 按非终结符分组产生式
	productionMap := make(map[model.Symbol][]model.Production)
	for _, prod := range grammar.Productions {
		if len(prod.Left) > 0 {
			left := prod.Left[0]
			productionMap[left] = append(productionMap[left], prod)
		}
	}

	// 检查每个非终结符的产生式是否有公共前缀
	for _, prods := range productionMap {
		if len(prods) <= 1 {
			continue
		}

		for i := 0; i < len(prods); i++ {
			for j := i + 1; j < len(prods); j++ {
				if len(findCommonPrefix(prods[i].Right, prods[j].Right)) > 0 {
					return true
				}
			}
		}
	}

	return false
}

// hasLeftRecursion 检查文法是否有左递归
func hasLeftRecursion(grammar *model.Grammar) bool {
	// 检查直接左递归：A -> A α
	for _, prod := range grammar.Productions {
		if len(prod.Right) > 0 && len(prod.Left) > 0 {
			left := prod.Left[0]
			first := prod.Right[0]
			// 直接左递归
			if left == first {
				return true
			}
		}
	}

	// 检查间接左递归：使用可达性分析
	// 构建非终结符的可达关系图
	reachable := make(map[model.Symbol]map[model.Symbol]bool)
	for _, nt := range grammar.NonTerminals {
		reachable[nt] = make(map[model.Symbol]bool)
	}

	// 第一步：直接可达关系
	for _, prod := range grammar.Productions {
		if len(prod.Right) > 0 && len(prod.Left) > 0 {
			left := prod.Left[0]
			first := prod.Right[0]
			// 如果右部第一个符号是非终结符，建立可达关系
			if grammar.CheckIsNonTerminal(first) {
				reachable[left][first] = true
			}
		}
	}

	// Floyd-Warshall算法计算传递闭包
	for _, k := range grammar.NonTerminals {
		for _, i := range grammar.NonTerminals {
			for _, j := range grammar.NonTerminals {
				if reachable[i][k] && reachable[k][j] {
					reachable[i][j] = true
				}
			}
		}
	}

	// 检查是否存在 A *=> A（即自己可达自己）
	for _, nt := range grammar.NonTerminals {
		if reachable[nt][nt] {
			return true
		}
	}

	return false
}

// HasLeftRecursion 检查文法是否有左递归（导出用于测试）
func HasLeftRecursion(grammar *model.Grammar) bool {
	return hasLeftRecursion(grammar)
}

// EliminateLeftRecursion 消除左递归（导出用于测试）
func EliminateLeftRecursion(grammar *model.Grammar) *model.Grammar {
	return eliminateLeftRecursion(grammar)
}

// eliminateLeftRecursion 消除左递归
func eliminateLeftRecursion(grammar *model.Grammar) *model.Grammar {
	// 创建新文法的副本
	newGrammar := &model.Grammar{
		StartSymbol:  grammar.StartSymbol,
		Terminals:    append([]model.Symbol{}, grammar.Terminals...),
		NonTerminals: append([]model.Symbol{}, grammar.NonTerminals...),
		Productions:  []model.Production{},
	}

	// 按非终结符分组产生式
	productionMap := make(map[model.Symbol][]model.Production)
	for _, prod := range grammar.Productions {
		if len(prod.Left) > 0 {
			left := prod.Left[0]
			productionMap[left] = append(productionMap[left], prod)
		}
	}

	// 为每个非终结符消除直接左递归
	for _, nt := range grammar.NonTerminals {
		prods := productionMap[nt]

		// 分离左递归和非左递归产生式
		var leftRecursive []model.Production
		var nonLeftRecursive []model.Production

		for _, prod := range prods {
			if len(prod.Right) > 0 && prod.Right[0] == nt {
				// 左递归产生式：A -> A α
				leftRecursive = append(leftRecursive, prod)
			} else {
				// 非左递归产生式：A -> β
				nonLeftRecursive = append(nonLeftRecursive, prod)
			}
		}

		if len(leftRecursive) > 0 {
			// 需要消除左递归
			// 创建新的非终结符 A'
			newNT := model.Symbol(string(nt) + "'")
			newGrammar.NonTerminals = append(newGrammar.NonTerminals, newNT)

			// 替换原有产生式
			for _, prod := range nonLeftRecursive {
				// A -> β 变为 A -> β A'
				newRight := append(prod.Right, newNT)
				newGrammar.Productions = append(newGrammar.Productions, model.Production{
					Left:  prod.Left,
					Right: newRight,
				})
			}

			// 添加新产生式
			for _, prod := range leftRecursive {
				// A -> A α 变为 A' -> α A'
				alpha := prod.Right[1:] // 去掉第一个A
				newRight := append(alpha, newNT)
				newGrammar.Productions = append(newGrammar.Productions, model.Production{
					Left:  []model.Symbol{newNT},
					Right: newRight,
				})
			}

			// 添加 A' -> ε
			newGrammar.Productions = append(newGrammar.Productions, model.Production{
				Left:  []model.Symbol{newNT},
				Right: []model.Symbol{model.Epsilon},
			})
		} else {
			// 无左递归，直接复制
			for _, prod := range nonLeftRecursive {
				newGrammar.Productions = append(newGrammar.Productions, prod)
			}
		}
	}

	return newGrammar
}

// leftFactor 提取左公因子
func leftFactor(grammar *model.Grammar) *model.Grammar {
	// 创建新文法的副本
	newGrammar := &model.Grammar{
		StartSymbol:  grammar.StartSymbol,
		Terminals:    append([]model.Symbol{}, grammar.Terminals...),
		NonTerminals: append([]model.Symbol{}, grammar.NonTerminals...),
		Productions:  []model.Production{},
	}

	// 按非终结符分组产生式
	productionMap := make(map[model.Symbol][]model.Production)
	for _, prod := range grammar.Productions {
		if len(prod.Left) > 0 {
			left := prod.Left[0]
			productionMap[left] = append(productionMap[left], prod)
		}
	}

	counter := 1

	// 为每个非终结符提取左公因子
	for _, nt := range grammar.NonTerminals {
		prods := productionMap[nt]
		if len(prods) <= 1 {
			// 只有一个或零个产生式，无需处理
			for _, prod := range prods {
				newGrammar.Productions = append(newGrammar.Productions, prod)
			}
			continue
		}

		// 寻找公共前缀
		processed := make(map[int]bool)
		for i, prod1 := range prods {
			if processed[i] {
				continue
			}

			// 找到具有相同前缀的产生式组
			var group []int
			var commonPrefix []model.Symbol

			for j, prod2 := range prods {
				if i != j && !processed[j] {
					prefix := findCommonPrefix(prod1.Right, prod2.Right)
					if len(prefix) > 0 {
						if len(group) == 0 {
							group = append(group, i)
							commonPrefix = prefix
						}
						group = append(group, j)
					}
				}
			}

			if len(group) > 1 {
				// 需要提取左公因子
				// 创建新的非终结符
				newNT := model.Symbol(fmt.Sprintf("%s%d", string(nt), counter))
				counter++
				newGrammar.NonTerminals = append(newGrammar.NonTerminals, newNT)

				// 添加新产生式：A -> α A'
				newRight := append(commonPrefix, newNT)
				newGrammar.Productions = append(newGrammar.Productions, model.Production{
					Left:  []model.Symbol{nt},
					Right: newRight,
				})

				// 为组中每个产生式创建 A' -> β
				for _, idx := range group {
					processed[idx] = true
					prod := prods[idx]
					remainder := prod.Right[len(commonPrefix):]
					if len(remainder) == 0 {
						remainder = []model.Symbol{model.Epsilon}
					}
					newGrammar.Productions = append(newGrammar.Productions, model.Production{
						Left:  []model.Symbol{newNT},
						Right: remainder,
					})
				}
			} else {
				// 无公因子，直接复制
				newGrammar.Productions = append(newGrammar.Productions, prod1)
				processed[i] = true
			}
		}
	}

	return newGrammar
}

// findCommonPrefix 找到两个符号序列的公共前缀
func findCommonPrefix(seq1, seq2 []model.Symbol) []model.Symbol {
	var prefix []model.Symbol
	minLen := len(seq1)
	if len(seq2) < minLen {
		minLen = len(seq2)
	}

	for i := 0; i < minLen; i++ {
		if seq1[i] == seq2[i] {
			prefix = append(prefix, seq1[i])
		} else {
			break
		}
	}

	return prefix
}

// RecursiveDescentParse 递归下降分析器
func RecursiveDescentParse(grammar *model.Grammar, input []model.Symbol) *model.ParseResult {
	result := &model.ParseResult{
		Method: "递归下降分析",
		Steps:  []model.ParseStep{},
	}

	// 预处理文法：消除左递归和提取左公因子
	processedGrammar := grammar
	if hasLeftRecursion(grammar) {
		processedGrammar = eliminateLeftRecursion(processedGrammar)
		result.Message = "文法包含左递归，已自动消除"
	}

	if needsLeftFactoring(processedGrammar) {
		processedGrammar = leftFactor(processedGrammar)
		if result.Message != "" {
			result.Message += "; "
		}
		result.Message += "文法包含左公因子，已自动提取"
	}

	// 检查处理后的文法是否适用递归下降
	if !canUseRecursiveDescent(processedGrammar) {
		result.Accepted = false
		result.Error = "文法处理后仍不适用于递归下降分析"
		return result
	}

	// 创建解析器实例
	parser := &RecursiveDescentParser{
		Grammar: processedGrammar,
		Input:   append(copySymbols(input), "#"), // 添加输入结束标记
		Pos:     0,
		Steps:   []model.ParseStep{},
	}

	// 记录初始状态
	parser.addStep(model.ParseStep{
		StepType:    "init",
		Description: "开始递归下降分析",
		Stack:       []model.Symbol{processedGrammar.StartSymbol},
		Input:       parser.Input,
		InputPos:    parser.Pos,
		Action:      "从起始符号开始分析",
	})

	// 开始递归下降分析
	success := parser.parseSymbol(processedGrammar.StartSymbol)

	if success && parser.getCurrentSymbol() == "#" {
		parser.addStep(model.ParseStep{
			StepType:    "accept",
			Description: "递归下降分析成功",
			Stack:       []model.Symbol{},
			Input:       []model.Symbol{"#"},
			InputPos:    parser.Pos,
			Action:      "接受输入串",
		})
		result.Accepted = true
		if result.Message == "" {
			result.Message = "输入串被递归下降分析器接受"
		}
	} else {
		parser.addStep(model.ParseStep{
			StepType:    "error",
			Description: "递归下降分析失败",
			Stack:       []model.Symbol{},
			Input:       parser.Input[parser.Pos:],
			InputPos:    parser.Pos,
			Action:      "拒绝输入串",
		})
		result.Accepted = false
		result.Error = "输入串不能被文法接受"
	}

	result.Steps = parser.Steps
	return result
}

// addStep 为递归下降分析器添加步骤
func (p *RecursiveDescentParser) addStep(step model.ParseStep) {
	p.Steps = append(p.Steps, step)
}

// getCurrentSymbol 获取当前输入符号
func (p *RecursiveDescentParser) getCurrentSymbol() model.Symbol {
	if p.Pos >= len(p.Input) {
		return "#" // 输入结束
	}
	return p.Input[p.Pos]
}

// consume 消费一个输入符号
func (p *RecursiveDescentParser) consume(expected model.Symbol) bool {
	if p.getCurrentSymbol() == expected {
		p.Pos++
		p.addStep(model.ParseStep{
			StepType:    "match",
			Description: fmt.Sprintf("匹配终结符 %s", string(expected)),
			Stack:       []model.Symbol{}, // 递归下降不使用显式栈
			Input:       p.Input[p.Pos:],
			InputPos:    p.Pos,
			Action:      fmt.Sprintf("消费符号 %s", string(expected)),
		})
		return true
	}
	return false
}

// parseSymbol 解析一个符号（终结符或非终结符）
func (p *RecursiveDescentParser) parseSymbol(symbol model.Symbol) bool {
	// 如果是终结符，直接匹配
	if p.Grammar.CheckIsTerminal(symbol) {
		return p.consume(symbol)
	}

	// 如果是非终结符，查找适合的产生式
	return p.parseNonTerminal(symbol)
}

// parseNonTerminal 解析非终结符
func (p *RecursiveDescentParser) parseNonTerminal(nt model.Symbol) bool {
	current := p.getCurrentSymbol()
	firstSet := CalculateFirstCached(p.Grammar)

	// 找到所有以 nt 为左部的产生式
	var candidates []model.Production
	for _, prod := range p.Grammar.Productions {
		if len(prod.Left) > 0 && prod.Left[0] == nt {
			candidates = append(candidates, prod)
		}
	}

	if len(candidates) == 0 {
		p.addStep(model.ParseStep{
			StepType:    "error",
			Description: fmt.Sprintf("未找到非终结符 %s 的产生式", string(nt)),
			Stack:       []model.Symbol{nt},
			Input:       p.Input[p.Pos:],
			InputPos:    p.Pos,
			Action:      "解析失败",
		})
		return false
	}

	// 选择适合的产生式（基于FIRST集）
	var selectedProd *model.Production
	for _, prod := range candidates {
		if len(prod.Right) == 0 || (len(prod.Right) == 1 && prod.Right[0] == model.Epsilon) {
			// 空产生式，检查FOLLOW集
			followSet := CalculateFollowCached(p.Grammar, firstSet)
			if _, inFollow := followSet[nt][current]; inFollow || current == "#" {
				selectedProd = &prod
				break
			}
		} else {
			// 非空产生式，检查FIRST集
			firstSymbol := prod.Right[0]
			if firstSymbol == current {
				// 直接匹配
				selectedProd = &prod
				break
			} else if p.Grammar.CheckIsNonTerminal(firstSymbol) {
				// 检查FIRST集
				if _, inFirst := firstSet[firstSymbol][current]; inFirst {
					selectedProd = &prod
					break
				}
				// 检查是否可以推出ε
				if _, hasEpsilon := firstSet[firstSymbol][model.Epsilon]; hasEpsilon {
					// 递归检查后续符号
					if p.canDerive(prod.Right, current, firstSet) {
						selectedProd = &prod
						break
					}
				}
			}
		}
	}

	if selectedProd == nil {
		p.addStep(model.ParseStep{
			StepType:    "error",
			Description: fmt.Sprintf("没有适合的产生式匹配当前输入 %s", string(current)),
			Stack:       []model.Symbol{nt},
			Input:       p.Input[p.Pos:],
			InputPos:    p.Pos,
			Action:      "解析失败",
		})
		return false
	}

	// 记录使用的产生式
	p.addStep(model.ParseStep{
		StepType:    "predict",
		Description: fmt.Sprintf("使用产生式: %s → %s", string(nt), symbolsToStringJoinSepFast(selectedProd.Right)),
		Stack:       []model.Symbol{nt},
		Input:       p.Input[p.Pos:],
		InputPos:    p.Pos,
		Action:      fmt.Sprintf("选择产生式并展开 %s", string(nt)),
		Production:  selectedProd,
	})

	// 递归解析产生式右部
	for _, symbol := range selectedProd.Right {
		if symbol == model.Epsilon {
			// 空产生式，不需要处理
			continue
		}
		if !p.parseSymbol(symbol) {
			return false
		}
	}

	return true
}

// canDerive 检查符号序列是否能推出特定符号
func (p *RecursiveDescentParser) canDerive(symbols []model.Symbol, target model.Symbol, firstSet map[model.Symbol]map[model.Symbol]struct{}) bool {
	for i, symbol := range symbols {
		if symbol == target {
			return true
		}
		if p.Grammar.CheckIsTerminal(symbol) {
			return symbol == target
		}
		if _, inFirst := firstSet[symbol][target]; inFirst {
			return true
		}
		if _, hasEpsilon := firstSet[symbol][model.Epsilon]; !hasEpsilon {
			return false
		}
		// 如果能推出ε，继续检查下一个符号
		if i == len(symbols)-1 {
			// 所有符号都能推出ε
			return false
		}
	}
	return false
}
