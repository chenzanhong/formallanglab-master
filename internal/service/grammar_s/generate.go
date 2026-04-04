package grammar_s

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/service/automaton_s"
	"github.com/chenzanhong/formallanglab-master/pkg/util"
	"github.com/chenzanhong/zlog"
)

const (
	maxExampleNum = 6
	maxTimes      = 1000
)

type derivationNode struct {
	sententialForm []model.Symbol // 当前句型（可能含非终结符）
	depth          int
}

func isTerminals(ss []model.Symbol, set map[model.Symbol]bool) bool {
	for _, s := range ss {
		if !set[s] {
			return false
		}
	}

	return true
}

// 生成文法可推导和不可推导的字符串示例
func GrammarGenerateExampleString(g *model.Grammar) (accept, reject []string) {
	// 转DFA+补集
	if g.GrammarType == model.RegularGrammar {
		if acc, rej := grammarGenerateExampleStringByCompletedDFAAndBFS(g); len(acc) > 0 && len(rej) > 0 {
			return acc, rej
		}
	}
	// 枚举+验证
	return grammarGenerateExampleStringByEnumAndVerify(g)
}

func grammarGenerateExampleStringByCompletedDFAAndBFS(g *model.Grammar) (accept, reject []string) {
	DetermineLinearity(g)
	a, err := RegularGrammarToFA(g)
	if err != nil { // 非线性，但是理论上不会出现，前面已经确认是正则文法了
		return grammarGenerateExampleStringByEnumAndVerify(g)
	}
	if a.Type != model.DFA {
		a = automaton_s.NFAToDFA(a)
	}
	if err := a.CompleteDFA(); err != nil { // 完备化失败
		return grammarGenerateExampleStringByEnumAndVerify(g)
	}

	return automaton_s.GenerateExampleStringsFromCompletedDFA(a)
}

func grammarGenerateExampleStringByEnumAndVerify(g *model.Grammar) (accept, reject []string) {
	// 检查文法类型：只支持 Regular 和 CFG（即 2/3 型）
	if g.GrammarType != model.RegularGrammar && g.GrammarType != model.ContextFreeGrammar {
		msg := "(仅支持 2 型/3 型文法示例生成)"
		return []string{msg}, []string{msg}
	}
	accept = generateAcceptExampleString(g)
	seenAccept := make(map[string]bool)
	for _, s := range accept {
		seenAccept[s] = true
	}
	/*
		// 根据文法类型生成 reject
		if g.GrammarType == model.RegularGrammar {
			// 正则文法 → NFA
			// NFA → DFA（子集构造）
			// DFA → 补集 DFA（翻转接受状态）
			// 在补集 DFA 上 BFS 搜索最短路径 → 得到 reject 字符串
			reject = generateRejectViaComplementDFA(g)
		} else {
			reject = generateRejectByExampleEnumeration(g, seenAccept)
		}
	*/
	reject = generateRejectExampleString(g, seenAccept)

	return accept, reject
}

func generateAcceptExampleString(g *model.Grammar) (accept []string) {
	zlog.Infow("开始生成文法示例字符串",
		"startSymbol", g.StartSymbol,
		"productions", len(g.Productions),
		"terminals", len(g.Terminals),
		"nonTerminals", len(g.NonTerminals))

	cfgView, err := g.ToCFGView()
	if err != nil {
		zlog.Warnw("文法转换为 CFG 视图失败", "error", err)
		// 如果不是 CFG（左部非单个非终结符），暂不支持
		return []string{"(仅支持上下文无关文法示例生成)"}
	}

	queue := []derivationNode{{sententialForm: []model.Symbol{g.StartSymbol}, depth: 0}}
	seenAccept := make(map[string]bool)
	seenForms := make(map[string]bool) // 避免重复处理相同的句型

	maxDerivationDepth := max(5, min(len(g.Productions), 8)) // 降低深度限制
	maxSentenceLength := max(6, min(len(g.Terminals)*2, 8))  // 降低长度限制
	maxQueueSize := 100                                      // 严格限制队列大小，避免 OOM

	terminalsSet := make(map[model.Symbol]bool)
	for _, s := range g.Terminals {
		terminalsSet[s] = true
	}

	nonTerminalsSet := make(map[model.Symbol]bool)
	for _, s := range g.NonTerminals {
		nonTerminalsSet[s] = true
	}

	zlog.Infow("开始 BFS 推导",
		"maxDerivationDepth", maxDerivationDepth,
		"maxSentenceLength", maxSentenceLength,
		"maxQueueSize", maxQueueSize)

	times := 0 // 避免左递归导致无限步
	for len(queue) > 0 && len(accept) < maxExampleNum && times < maxTimes && len(queue) < maxQueueSize {
		times++
		node := queue[0]
		queue = queue[1:]
		s := util.SymbolsToString(node.sententialForm)

		// 每 50 次打印一次进度
		if times%50 == 0 {
			zlog.Infow("生成进度",
				"times", times,
				"queueSize", len(queue),
				"acceptCount", len(accept),
				"currentForm", s,
				"depth", node.depth)
		}

		// 跳过已处理过的句型
		if seenForms[s] {
			continue
		}
		seenForms[s] = true

		// 限制句型长度，避免无限增长
		if len(node.sententialForm) > maxSentenceLength {
			continue
		}

		// 如果是句子（全终结符），加入 accept
		if isTerminals(node.sententialForm, terminalsSet) {
			if !seenAccept[s] && len(s) <= maxSentenceLength {
				accept = append(accept, s)
				seenAccept[s] = true
				zlog.Infow("找到接受字符串",
					"acceptCount", len(accept),
					"string", s,
					"depth", node.depth)
				// 达到目标数量后立即返回，不再继续生成
				if len(accept) >= maxExampleNum {
					zlog.Infow("接受字符串已达到目标数量，停止生成", "acceptCount", len(accept))
					return accept
				}
			}
			continue // 句子不能再推导
		}

		// 超过深度限制，跳过
		if node.depth >= maxDerivationDepth {
			if times%100 == 0 {
				zlog.Debugw("跳过 - 超过深度限制", "depth", node.depth, "maxDepth", maxDerivationDepth)
			}
			continue
		}

		// 找第一个非终结符进行替换（最左推导）
		found := false
		for i, sym := range node.sententialForm {
			if nonTerminalsSet[sym] {
				// 应用所有以 sym 为左部的产生式
				if rhsList, ok := cfgView.Productions[sym]; ok {
					// 优先尝试包含终结符的产生式和 ε 产生式
					rankedProductions := rankProductions(rhsList, terminalsSet, nonTerminalsSet)

					addedCount := 0
					skippedByLength := 0
					skippedBySeen := 0

					for _, rhs := range rankedProductions {
						var newForm []model.Symbol
						if len(rhs) == 1 && rhs[0] == model.Epsilon {
							// 空产生式
							newForm = make([]model.Symbol, 0, len(node.sententialForm)-1)
							newForm = append(newForm, node.sententialForm[:i]...)
							newForm = append(newForm, node.sententialForm[i+1:]...)
						} else {
							newForm = make([]model.Symbol, 0, len(node.sententialForm)-1+len(rhs))
							newForm = append(newForm, node.sententialForm[:i]...)
							newForm = append(newForm, rhs...)
							newForm = append(newForm, node.sententialForm[i+1:]...)
						}

						// 检查新句型是否过长
						if len(newForm) > maxSentenceLength {
							skippedByLength++
							continue
						}

						newFormStr := util.SymbolsToString(newForm)
						if !seenForms[newFormStr] {
							queue = append(queue, derivationNode{
								sententialForm: newForm,
								depth:          node.depth + 1,
							})
							addedCount++
						} else {
							skippedBySeen++
						}
					}

					if times%100 == 0 {
						zlog.Debugw("应用产生式",
							"symbol", sym,
							"productions", len(rhsList),
							"added", addedCount,
							"skippedByLength", skippedByLength,
							"skippedBySeen", skippedBySeen,
							"newQueueSize", len(queue))
					}
				}
				found = true
				break // 最左推导，只替换第一个非终结符
			}
		}

		// 如果没找到非终结符（理论上不会发生），跳过
		if !found {
			continue
		}
	}

	if len(accept) == 0 {
		zlog.Warnw("未找到短接受字符串", "times", times, "finalQueueSize", len(queue))
		// 尝试直接生成空字符串
		if g.StartSymbol == "S" {
			for _, prod := range g.Productions {
				if len(prod.Left) == 1 && prod.Left[0] == "S" && len(prod.Right) == 1 && prod.Right[0] == model.Epsilon {
					zlog.Infow("找到空产生式，返回ε")
					return []string{"ε"}
				}
			}
		}
		return []string{"(未找到短接受字符串)"}
	}

	zlog.Infow("文法示例生成完成",
		"acceptCount", len(accept),
		"totalIterations", times,
		"finalQueueSize", len(queue))

	// 确保返回的字符串数量不超过 maxExampleNum
	if len(accept) > maxExampleNum {
		accept = accept[:maxExampleNum]
	}

	return accept
}

// 对产生式进行排序，优先选择能快速终结的产生式
func rankProductions(rhsList [][]model.Symbol, terminalsSet, nonTerminalsSet map[model.Symbol]bool) [][]model.Symbol {
	type productionInfo struct {
		rhs   []model.Symbol
		score int // 分数越高，优先级越高
	}

	infos := make([]productionInfo, len(rhsList))
	for i, rhs := range rhsList {
		score := 0

		// ε 产生式优先级最高（能快速终结）
		if len(rhs) == 1 && rhs[0] == model.Epsilon {
			score = 1000
		} else {
			// 计算终结符数量和非终结符数量
			terminalCount := 0
			nonTerminalCount := 0
			for _, sym := range rhs {
				if terminalsSet[sym] {
					terminalCount++
				} else if nonTerminalsSet[sym] {
					nonTerminalCount++
				}
			}

			// 优先选择：终结符多、非终结符少的产生式
			// 对于 S → SS 这种会导致爆炸的产生式，给予极低分数
			if nonTerminalCount >= 2 {
				score = -100 // 严重惩罚会增加非终结符的产生式
			} else {
				// 终结符越多分数越高，非终结符越少分数越高
				score = terminalCount*50 - nonTerminalCount*30
			}

			// 额外奖励：纯终结符的产生式（能直接完成推导）
			if nonTerminalCount == 0 && terminalCount > 0 {
				score += 500
			}
		}

		infos[i] = productionInfo{rhs: rhs, score: score}
	}

	// 按分数排序
	for i := 0; i < len(infos)-1; i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].score < infos[j].score {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}

	// 提取排序后的产生式
	sorted := make([][]model.Symbol, len(infos))
	for i, info := range infos {
		sorted[i] = info.rhs
	}

	return sorted
}

// 通过枚举 + 识别，找寻拒绝字符串
func generateRejectExampleString(g *model.Grammar, seenAccept map[string]bool) (reject []string) {
	zlog.Infow("开始生成拒绝字符串",
		"terminals", len(g.Terminals),
		"productions", len(g.Productions))

	var pureTerminals []string
	for _, s := range g.Terminals {
		if s != model.Epsilon {
			pureTerminals = append(pureTerminals, string(s))
		}
	}

	if len(pureTerminals) == 0 {
		zlog.Warnw("没有可用的终结符生成拒绝字符串")
		return []string{"(无法生成拒绝字符串：无可用终结符)"}
	}

	seenReject := make(map[string]bool)

	// 策略 1：优先生成奇数长度的字符串（很多文法只接受偶数长度）
	// 策略 2：生成单一字符重复的字符串（如 "aaa", "bbb"）
	// 策略 3：随机生成字符串
	maxValidation := 30 // 最多验证 30 次
	validationCount := 0
	validatedCount := 0
	rejectedCount := 0

	// 生成候选字符串的辅助函数
	generateCandidates := func() []string {
		var candidates []string

		// 策略 1：奇数长度字符串（1, 3, 5）
		for _, term := range pureTerminals {
			for len := 1; len <= 5; len += 2 {
				candidates = append(candidates, strings.Repeat(term, len))
			}
		}

		// 策略 2：混合但长度不平衡
		if len(pureTerminals) >= 2 {
			// 如 "aab", "abb", "aaab", "abbb"
			for i := 1; i <= 3; i++ {
				for j := 1; j <= 3; j++ {
					if i != j { // 只生成不平衡的组合
						candidates = append(candidates,
							strings.Repeat(pureTerminals[0], i)+
								strings.Repeat(pureTerminals[1], j))
					}
				}
			}
		}

		// 策略 3：随机生成一些字符串
		rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
		for i := 0; i < 10; i++ {
			length := rng.IntN(5) + 1 // 1-5
			s := ""
			for j := 0; j < length; j++ {
				s += pureTerminals[rng.IntN(len(pureTerminals))]
			}
			candidates = append(candidates, s)
		}

		return candidates
	}

	candidates := generateCandidates()
	zlog.Infow("生成候选字符串",
		"strategy", "smart",
		"totalCandidates", len(candidates),
		"candidates", candidates[:min(10, len(candidates))])

	// 筛选出拒绝字符串
	for _, s := range candidates {
		if len(reject) >= maxExampleNum || validationCount >= maxValidation {
			break
		}
		if seenAccept[s] {
			continue
		}
		if !seenReject[s] {
			validationCount++
			validatedCount++

			zlog.Infow("判断绝字符串", "s", s)
			// 使用快速模式验证（带超时限制）
			// 对于简单的 membership test，不需要完整的解析分析
			result := quickValidateString(g, s)
			if !result.Accepted {
				reject = append(reject, s)
				seenReject[s] = true
				rejectedCount++
				zlog.Infow("找到拒绝字符串",
					"rejectCount", len(reject),
					"string", s,
					"parseResult", result.Error)
			}
		}
	}

	zlog.Infow("拒绝字符串生成完成",
		"validated", validatedCount,
		"rejected", rejectedCount,
		"finalCount", len(reject))

	if len(reject) == 0 {
		zlog.Warnw("未找到拒绝字符串")
		return []string{"(在终结符集中未找到短拒绝字符串)"}
	}

	// 确保 reject 数量不超过 maxExampleNum
	if len(reject) > maxExampleNum {
		zlog.Infow("拒绝字符串过多，随机选择", "total", len(reject), "select", maxExampleNum)
		// 随机选择 maxExampleNum 个
		rand.Shuffle(len(reject), func(i, j int) { reject[i], reject[j] = reject[j], reject[i] })
		reject = reject[:maxExampleNum]
	}

	return reject
}

// quickValidateString 快速验证字符串是否被文法接受
// 使用带严格限制的 BFS，避免复杂文法导致性能问题
func quickValidateString(g *model.Grammar, input string) *model.ParseResult {
	// 转换输入为符号数组
	inputSymbols, err := g.StringToSymbols(input)
	if err != nil {
		return &model.ParseResult{
			Accepted: false,
			Error:    fmt.Sprintf("输入格式错误：%s", err.Error()),
		}
	}

	// 使用严格的 BFS 限制：最多 50 步，宽度 20
	// 这样即使对于复杂文法也能在几毫秒内完成
	accepted, _, err := RecognizeString(g, inputSymbols, 50, 20)
	if err != nil {
		return &model.ParseResult{
			Accepted: false,
			Error:    err.Error(),
		}
	}

	return &model.ParseResult{
		Accepted: accepted,
		Method:   "BFS (受限)",
	}
}
