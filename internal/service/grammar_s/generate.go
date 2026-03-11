package grammar_s

import (
	"backend/internal/domain/model"
	"backend/internal/service/automaton_s"
	"backend/pkg/util"
	"math/rand/v2"
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
		grammarGenerateExampleStringByCompletedDFAAndBFS(g)
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
	cfgView, err := g.ToCFGView()
	if err != nil {
		// 如果不是 CFG（左部非单个非终结符），暂不支持
		return []string{"(仅支持上下文无关文法示例生成)"}
	}

	queue := []derivationNode{{sententialForm: []model.Symbol{g.StartSymbol}, depth: 0}}
	seenAccept := make(map[string]bool)
	seenForms := make(map[string]bool) // 避免重复处理相同的句型

	maxDerivationDepth := len(g.Productions) * 3
	if maxDerivationDepth < 10 {
		maxDerivationDepth = 10
	}

	terminalsSet := make(map[model.Symbol]bool)
	for _, s := range g.Terminals {
		terminalsSet[s] = true
	}

	nonTerminalsSet := make(map[model.Symbol]bool)
	for _, s := range g.NonTerminals {
		nonTerminalsSet[s] = true
	}

	times := 0           // 避免左递归导致无限步
	maxQueueSize := 1000 // 限制队列大小，避免内存溢出
	for len(queue) > 0 && len(accept) < maxExampleNum && times < maxTimes && len(queue) < maxQueueSize {
		times++
		node := queue[0]
		queue = queue[1:]
		s := util.SymbolsToString(node.sententialForm)

		// 跳过已处理过的句型
		if seenForms[s] {
			continue
		}
		seenForms[s] = true

		// 如果是句子（全终结符），加入 accept
		if isTerminals(node.sententialForm, terminalsSet) {
			if !seenAccept[s] {
				accept = append(accept, s)
				seenAccept[s] = true
			}
			continue // 句子不能再推导
		}

		// 超过深度限制，跳过
		if node.depth >= maxDerivationDepth {
			continue
		}

		// 找第一个非终结符进行替换（最左推导）
		found := false
		for i, sym := range node.sententialForm {
			if nonTerminalsSet[sym] {
				// 应用所有以 sym 为左部的产生式
				if rhsList, ok := cfgView.Productions[sym]; ok {
					for _, rhs := range rhsList {
						var newForm []model.Symbol
						if len(rhs) == 1 && rhs[0] == model.Epsilon {
							newForm = make([]model.Symbol, 0, len(node.sententialForm)-1)
							newForm = append(newForm, node.sententialForm[:i]...)
							newForm = append(newForm, node.sententialForm[i+1:]...)
						} else {
							newForm = make([]model.Symbol, 0, len(node.sententialForm)-1+len(rhs))
							newForm = append(newForm, node.sententialForm[:i]...)
							newForm = append(newForm, rhs...)
							newForm = append(newForm, node.sententialForm[i+1:]...)
						}

						newFormStr := util.SymbolsToString(newForm)
						if !seenForms[newFormStr] {
							queue = append(queue, derivationNode{
								sententialForm: newForm,
								depth:          node.depth + 1,
							})
						}
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
		return []string{"(未找到短接受字符串)"}
	}
	return accept
}

// 通过枚举+识别，找寻拒绝字符串
func generateRejectExampleString(g *model.Grammar, seenAccept map[string]bool) (reject []string) {
	var pureTerminals []string
	for _, s := range g.Terminals {
		if s != model.Epsilon {
			pureTerminals = append(pureTerminals, string(s))
		}
	}

	// 生成采样字符串，限制最大长度
	maxLen := min(len(g.Productions)*2, len(g.Productions)+4)
	if maxLen > 5 { // 限制最大长度，避免生成过长字符串
		maxLen = 5
	}
	ss := util.GenerateStrings(pureTerminals, maxLen)
	seenReject := make(map[string]bool)

	// 筛选出拒绝字符串，限制验证次数
	maxValidation := 100 // 最大验证次数，避免性能问题
	validationCount := 0
	for _, s := range ss {
		if len(reject) >= maxExampleNum || validationCount >= maxValidation {
			break
		}
		if seenAccept[s] {
			continue
		}
		if !seenReject[s] {
			validationCount++
			if !ParseStringWithMode(g, s, AutoMode, false).Accepted {
				reject = append(reject, s)
				seenReject[s] = true
			}
		}
	}

	if len(reject) == 0 {
		return []string{"(在终结符集中未找到短拒绝字符串)"}
	}

	// 确保reject数量不超过maxExampleNum
	if len(reject) > maxExampleNum {
		// 随机选择maxExampleNum个
		rand.Shuffle(len(reject), func(i, j int) { reject[i], reject[j] = reject[j], reject[i] })
		reject = reject[:maxExampleNum]
	}
	return reject
}
