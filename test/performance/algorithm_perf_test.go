package performance

import (
	"fmt"
	"testing"
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/service/automaton_s"
	regex_s "github.com/chenzanhong/formallanglab-master/internal/service/regex_s"
)

// TestCoreAlgorithmPerformance 核心算法性能测试
// 直接调用内部函数，不使用 HTTP
// 测试流程：正则表达式 → NFA → DFA → 最小化 DFA → 正则表达式
func TestCoreAlgorithmPerformance(t *testing.T) {
	// 测试次数
	testRuns := 1000

	// 测试用的正则表达式
	// testPattern := model.Regex("(a|b)*(c|ε)")
	// testPattern := model.Regex("(a|b)*(c|ε)(a|b+)?")
	testPattern := model.Regex("(a|b)*(c|ε)(a|b+)?|(a?b?c?)+")

	t.Run("核心算法端到端性能测试", func(t *testing.T) {
		var (
			regexToNfaLatencies  []time.Duration
			nfaToDfaLatencies    []time.Duration
			dfaMinimizeLatencies []time.Duration
			faToRegexLatencies   []time.Duration
		)

		for i := 0; i < testRuns; i++ {
			var step1Start, step2Start, step3Start, step4Start time.Time

			// 步骤 1：正则表达式 → NFA
			step1Start = time.Now()
			nfa, err := regex_s.RegexToFA(testPattern)
			if err != nil || nfa == nil {
				t.Logf("第%d次测试：正则表达式转 NFA 失败：%v", i+1, err)
				continue
			}
			regexToNfaLatencies = append(regexToNfaLatencies, time.Since(step1Start))

			// 步骤 2：NFA → DFA（确定化）
			step2Start = time.Now()
			dfa := automaton_s.NFAToDFA(nfa)
			if dfa == nil {
				t.Logf("第%d次测试：NFA 确定化失败", i+1)
				continue
			}
			nfaToDfaLatencies = append(nfaToDfaLatencies, time.Since(step2Start))

			// 步骤 3：DFA 最小化
			step3Start = time.Now()
			minimizedDfa := automaton_s.DFAMinimize(dfa)
			if minimizedDfa == nil {
				t.Logf("第%d次测试：DFA 最小化失败", i+1)
				continue
			}
			dfaMinimizeLatencies = append(dfaMinimizeLatencies, time.Since(step3Start))

			// 步骤 4：最小化 DFA → 正则表达式
			step4Start = time.Now()
			regex := automaton_s.FAToRegex(minimizedDfa)
			if regex == "" {
				t.Logf("第%d次测试：DFA 转正则表达式失败", i+1)
				continue
			}
			faToRegexLatencies = append(faToRegexLatencies, time.Since(step4Start))
		}

		// 打印统计结果
		printAvgStats(t, "正则表达式到 NFA 转换", regexToNfaLatencies)
		printAvgStats(t, "NFA 确定化", nfaToDfaLatencies)
		printAvgStats(t, "DFA 最小化", dfaMinimizeLatencies)
		printAvgStats(t, "DFA 到正则表达式转换", faToRegexLatencies)
	})
}

// printAvgStats 打印平均耗时统计（微秒级别）
func printAvgStats(t *testing.T, testName string, latencies []time.Duration) {
	if len(latencies) == 0 {
		t.Logf("%s: 无成功请求\n", testName)
		return
	}

	var totalNanoseconds int64
	for _, lat := range latencies {
		totalNanoseconds += lat.Nanoseconds()
	}

	avgNanoseconds := float64(totalNanoseconds) / float64(len(latencies))
	avgMicroseconds := avgNanoseconds / 1000.0

	fmt.Printf("%s: 平均耗时 %.2f μs\n", testName, avgMicroseconds)
}
