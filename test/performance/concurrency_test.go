package performance

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/service/automaton_s"
)

// TestConcurrencyPerformance 并发性能测试
func TestConcurrencyPerformance(t *testing.T) {
	// 并发配置
	numGoroutines := 60000      // 60000 个并发 goroutine
	iterationsPerGoroutine := 1 // 每个 goroutine 只执行 1 次识别

	// 测试字符串：3 个接受，3 个拒绝
	testStrings := []struct {
		input    string
		expected bool // true=接受，false=拒绝
	}{
		{"ab", true},              // 接受
		{"ababbab", true},         // 接受
		{"bbbaaabababab", true},   // 接受
		{"aaa", false},            // 拒绝
		{"ababbbb", false},        // 拒绝
		{"bbbaaababababb", false}, // 拒绝
	}

	t.Run("自动机识别字符串并发测试", func(t *testing.T) {
		var (
			successCount int64
			failCount    int64
			totalLatency int64
			wg           sync.WaitGroup
		)

		startTime := time.Now()
		testStartCtx, cancel := context.WithCancel(context.Background())

		// 启动 6000 个 goroutine，每个执行一次识别
		for goroutineID := 0; goroutineID < numGoroutines; goroutineID++ {
			wg.Add(1)

			go func(id int) {
				defer wg.Done()

				testStr := testStrings[id%len(testStrings)]
				expectedAccept := testStr.expected
				<-testStartCtx.Done()

				// 深拷贝自动机
				testAutomaton := &model.Automaton{
					States:          []model.State{"q0", "q1", "q2"},
					Alphabet:        []model.Symbol{"a", "b"},
					InitialState:    "q0",
					AcceptingStates: []model.State{"q2"},
					Transitions: []model.Transition{
						{FromState: "q0", Input: "b", ToStates: []model.State{"q0"}},
						{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}},
						{FromState: "q1", Input: "a", ToStates: []model.State{"q1"}},
						{FromState: "q1", Input: "b", ToStates: []model.State{"q2"}},
						{FromState: "q2", Input: "a", ToStates: []model.State{"q1"}},
						{FromState: "q2", Input: "b", ToStates: []model.State{"q0"}},
					},
					Type: model.DFA,
				}

				start := time.Now()

				// 调用验证函数
				validateErr := automaton_s.AutomatonValidate(testAutomaton)

				// 调用识别函数
				result, recognizeErr := automaton_s.Recognize(testAutomaton, testStr.input)

				latency := time.Since(start)
				atomic.AddInt64(&totalLatency, int64(latency))

				// 验证结果是否正确
				if validateErr == nil && recognizeErr == nil && result.IsAccepted == expectedAccept {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&failCount, 1)
					// if result.IsAccepted != expectedAccept {
					// 	t.Logf("结果不匹配：输入=%s, 期望=%v, 实际=%v", testStr.input, expectedAccept, result.IsAccepted)
					// }
				}
			}(goroutineID)
		}
		cancel()

		wg.Wait()
		totalDuration := time.Since(startTime)

		// 计算统计信息
		totalStrings := int64(numGoroutines * iterationsPerGoroutine)
		avgLatency := time.Duration(totalLatency / totalStrings)

		// 打印结果
		fmt.Printf("\n========== 并发测试报告 ==========\n")
		fmt.Printf("  Goroutine 数量：%d\n", numGoroutines)
		fmt.Printf("  每个 Goroutine 循环次数：%d\n", iterationsPerGoroutine)
		fmt.Printf("  总识别次数：%d\n", totalStrings)
		fmt.Printf("  平均响应时间：%v\n", avgLatency)
		fmt.Printf("  总耗时：%v\n", totalDuration)
		fmt.Printf("  吞吐量：%.2f ops/s\n", float64(totalStrings)/totalDuration.Seconds())
	})
}
