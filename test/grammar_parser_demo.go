// test/grammar_parser_demo.go
package test

import (
	"backend/internal/domain/model"
	"backend/internal/service/grammar_s"
	"fmt"
	"strings"
)

// DemoLL1Parser 演示LL(1)分析器的使用
func DemoLL1Parser() {
	// 测试用例1：简单的 LL(1) 文法
	testSimpleLL1Grammar()

	fmt.Println("\n" + strings.Repeat("=", 50))

	// 测试用例2：算术表达式文法
	testArithmeticGrammar()

	fmt.Println("\n" + strings.Repeat("=", 50))

	// 测试用例3：包含递归的复杂文法
	testRecursiveGrammar()

	fmt.Println("\n" + strings.Repeat("=", 50))

	// 测试用例4：错误处理测试
	testErrorHandling()

	fmt.Println("\n所有测试完成!")
}

func testSimpleLL1Grammar() {
	// 测试用例1：简单的 LL(1) 文法
	// S -> aAb
	// A -> c | ε
	grammar1 := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b", "c"},
		NonTerminals: []model.Symbol{"S", "A"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A", "b"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"c"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
		},
	}

	fmt.Println("=== 测试用例1: 简单LL(1)文法 ===")
	fmt.Println("文法:")
	fmt.Println("S -> aAb")
	fmt.Println("A -> c | ε")
	fmt.Println()

	// 测试输入串 "acb"
	input1 := []model.Symbol{"a", "c", "b"}
	fmt.Println("测试输入: acb")

	result1 := grammar_s.ParseStringWithMode(grammar1, input1, "ll1", true)

	fmt.Printf("结果: %s\n", getResultString(result1.Accepted))
	fmt.Printf("方法: %s\n", result1.Method)
	if result1.Error != "" {
		fmt.Printf("错误: %s\n", result1.Error)
	}
	if result1.Message != "" {
		fmt.Printf("消息: %s\n", result1.Message)
	}

	fmt.Println("\n分析步骤:")
	for i, step := range result1.Steps {
		fmt.Printf("%d. %s - %s\n", i+1, step.StepType, step.Description)
		if step.Production != nil {
			fmt.Printf("   使用产生式: %s -> %s\n",
				symbolsToString(step.Production.Left),
				symbolsToString(step.Production.Right))
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))

	// 测试输入串 "ab" (应该匹配 A -> ε)
	input2 := []model.Symbol{"a", "b"}
	fmt.Println("测试输入: ab")

	result2 := grammar_s.ParseStringWithMode(grammar1, input2, "ll1", true)

	fmt.Printf("结果: %s\n", getResultString(result2.Accepted))
	fmt.Printf("方法: %s\n", result2.Method)
	if result2.Error != "" {
		fmt.Printf("错误: %s\n", result2.Error)
	}
	if result2.Message != "" {
		fmt.Printf("消息: %s\n", result2.Message)
	}

	fmt.Println("\n简化步骤显示:")
	for i, step := range result2.Steps {
		if step.StepType == "predict" || step.StepType == "accept" {
			fmt.Printf("%d. %s - %s\n", i+1, step.StepType, step.Description)
		}
	}
}

func testArithmeticGrammar() {
	// 算术表达式文法 (LL(1))
	// E -> T E'
	// E' -> + T E' | ε
	// T -> F T'
	// T' -> * F T' | ε
	// F -> ( E ) | id
	grammar := &model.Grammar{
		StartSymbol:  "E",
		Terminals:    []model.Symbol{"+", "*", "(", ")", "id"},
		NonTerminals: []model.Symbol{"E", "E'", "T", "T'", "F"},
		Productions: []model.Production{
			{Left: []model.Symbol{"E"}, Right: []model.Symbol{"T", "E'"}},
			{Left: []model.Symbol{"E'"}, Right: []model.Symbol{"+", "T", "E'"}},
			{Left: []model.Symbol{"E'"}, Right: []model.Symbol{model.Epsilon}},
			{Left: []model.Symbol{"T"}, Right: []model.Symbol{"F", "T'"}},
			{Left: []model.Symbol{"T'"}, Right: []model.Symbol{"*", "F", "T'"}},
			{Left: []model.Symbol{"T'"}, Right: []model.Symbol{model.Epsilon}},
			{Left: []model.Symbol{"F"}, Right: []model.Symbol{"(", "E", ")"}},
			{Left: []model.Symbol{"F"}, Right: []model.Symbol{"id"}},
		},
	}

	fmt.Println("=== 测试用例2: 算术表达式文法 ===")
	fmt.Println("文法:")
	fmt.Println("E -> T E'")
	fmt.Println("E' -> + T E' | ε")
	fmt.Println("T -> F T'")
	fmt.Println("T' -> * F T' | ε")
	fmt.Println("F -> ( E ) | id")
	fmt.Println()

	// 测试输入: id + id * id
	input := []model.Symbol{"id", "+", "id", "*", "id"}
	fmt.Println("测试输入: id + id * id")

	result := grammar_s.ParseStringWithMode(grammar, input, "ll1", true)

	fmt.Printf("结果: %s\n", getResultString(result.Accepted))
	fmt.Printf("方法: %s\n", result.Method)
	if result.Error != "" {
		fmt.Printf("错误: %s\n", result.Error)
	}

	fmt.Println("\n分析步骤 (前10步):")
	for i, step := range result.Steps {
		if i >= 10 {
			fmt.Println("... (更多步骤)")
			break
		}
		fmt.Printf("%d. %s - %s\n", i+1, step.StepType, step.Description)
	}
}

func testRecursiveGrammar() {
	// 包含递归的复杂文法
	// S -> A B
	// A -> a A | ε
	// B -> b B c | d
	grammar := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b", "c", "d"},
		NonTerminals: []model.Symbol{"S", "A", "B"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A", "B"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a", "A"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
			{Left: []model.Symbol{"B"}, Right: []model.Symbol{"b", "B", "c"}},
			{Left: []model.Symbol{"B"}, Right: []model.Symbol{"d"}},
		},
	}

	fmt.Println("=== 测试用例3: 递归文法 ===")
	fmt.Println("文法:")
	fmt.Println("S -> A B")
	fmt.Println("A -> a A | ε")
	fmt.Println("B -> b B c | d")
	fmt.Println()

	// 测试输入: a a b b d c c
	input := []model.Symbol{"a", "a", "b", "b", "d", "c", "c"}
	fmt.Println("测试输入: a a b b d c c")

	result := grammar_s.ParseStringWithMode(grammar, input, "auto", true)

	fmt.Printf("结果: %s\n", getResultString(result.Accepted))
	fmt.Printf("方法: %s\n", result.Method)
	if result.Error != "" {
		fmt.Printf("错误: %s\n", result.Error)
	}

	fmt.Println("\n关键步骤:")
	for i, step := range result.Steps {
		if step.StepType == "predict" || step.StepType == "accept" || step.StepType == "error" {
			fmt.Printf("%d. %s - %s\n", i+1, step.StepType, step.Description)
			if step.Production != nil {
				fmt.Printf("   使用产生式: %s -> %s\n",
					symbolsToString(step.Production.Left),
					symbolsToString(step.Production.Right))
			}
		}
	}
}

func testErrorHandling() {
	// 简单文法用于测试错误处理
	grammar := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "b"}},
		},
	}

	fmt.Println("=== 测试用例4: 错误处理 ===")
	fmt.Println("文法: S -> a b")
	fmt.Println()

	// 测试错误输入
	errorInputs := [][]model.Symbol{
		{"a"},           // 输入不完整
		{"b", "a"},      // 顺序错误
		{"a", "b", "c"}, // 多余字符
		{"c"},           // 无效字符
	}

	errorDescriptions := []string{
		"输入不完整: a",
		"顺序错误: b a",
		"多余字符: a b c",
		"无效字符: c",
	}

	for i, input := range errorInputs {
		fmt.Printf("测试 %d: %s\n", i+1, errorDescriptions[i])

		result := grammar_s.ParseStringWithMode(grammar, input, "ll1", false)

		fmt.Printf("结果: %s\n", getResultString(result.Accepted))
		if result.Error != "" {
			fmt.Printf("错误: %s\n", result.Error)
		}
		fmt.Println()
	}
}


func symbolsToString(symbols []model.Symbol) string {
	if len(symbols) == 0 {
		return "ε"
	}

	var result string
	for i, sym := range symbols {
		if i > 0 {
			result += ""
		}
		if sym == model.Epsilon {
			result += "ε"
		} else {
			result += string(sym)
		}
	}
	return result
}

func getResultString(accepted bool) string {
	if accepted {
		return "接受"
	}
	return "拒绝"
}
