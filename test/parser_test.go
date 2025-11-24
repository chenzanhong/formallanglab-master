// test/parser_test.go
package test

import (
	"backend/internal/domain/model"
	"backend/internal/service/grammar_s"
	"fmt"
	"testing"
)

// TestRecursiveDescentParser 测试递归下降分析器
func TestRecursiveDescentParser(t *testing.T) {
	fmt.Println("=== 递归下降分析器测试 ===")

	// 简单LL(1)文法：S -> aAb, A -> c | ε
	grammar := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b", "c"},
		NonTerminals: []model.Symbol{"S", "A"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A", "b"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"c"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
		},
	}

	testCases := []struct {
		input    []model.Symbol
		expected bool
		desc     string
	}{
		{[]model.Symbol{"a", "c", "b"}, true, "输入 'acb' 应该被接受"},
		{[]model.Symbol{"a", "b"}, true, "输入 'ab' 应该被接受（A -> ε）"},
		{[]model.Symbol{"a", "c", "c", "b"}, false, "输入 'accb' 应该被拒绝"},
		{[]model.Symbol{"c", "b"}, false, "输入 'cb' 应该被拒绝"},
	}

	for _, tc := range testCases {
		result := grammar_s.RecursiveDescentParse(grammar, tc.input)
		if result.Accepted != tc.expected {
			t.Errorf("%s：期望 %v，实际 %v", tc.desc, tc.expected, result.Accepted)
		} else {
			fmt.Printf("✓ %s\n", tc.desc)
		}

		if result.Error != "" {
			fmt.Printf("  错误: %s\n", result.Error)
		}
		if result.Message != "" {
			fmt.Printf("  消息: %s\n", result.Message)
		}
	}
}

// TestLR0Parser 测试LR(0)分析器
func TestLR0Parser(t *testing.T) {
	fmt.Println("\n=== LR(0)分析器测试 ===")

	// 简单LR(0)文法：S' -> S, S -> aS | b
	grammar := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b"}},
		},
	}

	testCases := []struct {
		input    []model.Symbol
		expected bool
		desc     string
	}{
		{[]model.Symbol{"b"}, true, "输入 'b' 应该被接受"},
		{[]model.Symbol{"a", "b"}, true, "输入 'ab' 应该被接受"},
		{[]model.Symbol{"a", "a", "b"}, true, "输入 'aab' 应该被接受"},
		{[]model.Symbol{"a"}, false, "输入 'a' 应该被拒绝"},
		{[]model.Symbol{"a", "a"}, false, "输入 'aa' 应该被拒绝"},
	}

	for _, tc := range testCases {
		result := grammar_s.LR0ParseDetailed(grammar, tc.input)
		if result.Accepted != tc.expected {
			t.Errorf("%s：期望 %v，实际 %v", tc.desc, tc.expected, result.Accepted)
		} else {
			fmt.Printf("✓ %s\n", tc.desc)
		}

		if result.Error != "" {
			fmt.Printf("  错误: %s\n", result.Error)
		}
		if result.Message != "" {
			fmt.Printf("  消息: %s\n", result.Message)
		}
	}
}

// TestLeftRecursionElimination 测试左递归消除
func TestLeftRecursionElimination(t *testing.T) {
	fmt.Println("\n=== 左递归消除测试 ===")

	// 含左递归的文法：A -> Aa | b
	grammar := &model.Grammar{
		StartSymbol:  "A",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"A"},
		Productions: []model.Production{
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"A", "a"}}, // 左递归
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
		},
	}

	// 检查是否有左递归
	hasLR := grammar_s.HasLeftRecursion(grammar)
	if !hasLR {
		t.Error("应该检测到左递归")
	} else {
		fmt.Println("✓ 正确检测到左递归")
	}

	// 测试消除左递归
	processedGrammar := grammar_s.EliminateLeftRecursion(grammar)
	hasLRAfter := grammar_s.HasLeftRecursion(processedGrammar)
	if hasLRAfter {
		t.Error("左递归消除后仍然存在左递归")
	} else {
		fmt.Println("✓ 成功消除左递归")
	}

	// 验证处理后的文法仍能识别相同的语言
	testInput := []model.Symbol{"b", "a", "a"}
	originalResult := grammar_s.BFSParseDetailed(grammar, testInput, 50, 50)
	processedResult := grammar_s.RecursiveDescentParse(processedGrammar, testInput)

	if originalResult.Accepted != processedResult.Accepted {
		t.Errorf("左递归消除后语言识别能力发生变化：原文法 %v，处理后 %v",
			originalResult.Accepted, processedResult.Accepted)
	} else {
		fmt.Printf("✓ 消除左递归后语言等价性保持：输入 'baa' 识别结果 %v\n",
			processedResult.Accepted)
	}
}

// TestParseStringWithMode 测试模式选择功能
func TestParseStringWithMode(t *testing.T) {
	fmt.Println("\n=== 分析模式选择测试 ===")

	// LL(1)文法：S -> aAb, A -> c | ε
	grammar := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b", "c"},
		NonTerminals: []model.Symbol{"S", "A"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A", "b"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"c"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
		},
	}

	// 测试不同的分析模式
	modes := []string{"auto", "ll1", "recursive_descent", "bfs"}
	inputStr := "acb"

	for _, mode := range modes {
		result := grammar_s.ParseStringWithMode(grammar, inputStr, grammar_s.Mode(mode), false)
		fmt.Printf("模式 %s: %s (方法: %s)\n", mode,
			map[bool]string{true: "接受", false: "拒绝"}[result.Accepted],
			result.Method)

		if !result.Accepted {
			fmt.Printf("  错误: %s\n", result.Error)
		}
	}
}

// 辅助函数，需要导出一些内部函数用于测试
func init() {
	// 这里可以添加一些测试初始化代码
}
