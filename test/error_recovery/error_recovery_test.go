// test/error_recovery_test.go
package test

import (
	"backend/internal/domain/model"
	"backend/internal/service/grammar_s"
	"fmt"
	"strings"
	"testing"
)

func TestErrorRecovery(t *testing.T) {
	// 测试错误恢复功能
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

	fmt.Println("=== LL(1) 错误恢复测试 ===")
	fmt.Println("文法:")
	fmt.Println("S -> a A b")
	fmt.Println("A -> c | ε")
	fmt.Println()

	// 测试用例：包含错误的输入
	testCases := []struct {
		input       []model.Symbol
		description string
	}{
		{[]model.Symbol{"a", "x", "b"}, "错误输入: a x b (x是无效符号)"},
		{[]model.Symbol{"a", "c", "x", "b"}, "错误输入: a c x b (多余符号x)"},
		{[]model.Symbol{"x", "a", "c", "b"}, "错误输入: x a c b (开头有无效符号)"},
		{[]model.Symbol{"a", "c", "c", "b"}, "错误输入: a c c b (重复符号)"},
	}

	for i, tc := range testCases {
		fmt.Printf("测试 %d: %s\n", i+1, tc.description)

		// 普通LL(1)分析
		fmt.Println("\n普通LL(1)分析:")
		result1 := grammar_s.ParseStringWithMode(grammar, tc.input, "ll1", false)
		fmt.Printf("结果: %s\n", getResultString(result1.Accepted))
		if result1.Error != "" {
			fmt.Printf("错误: %s\n", result1.Error)
		}

		// 错误恢复LL(1)分析
		fmt.Println("\nLL(1)错误恢复分析:")
		result2 := grammar_s.ParseStringWithMode(grammar, tc.input, "ll1_recovery", true)
		fmt.Printf("结果: %s\n", getResultString(result2.Accepted))
		if result2.Error != "" {
			fmt.Printf("错误: %s\n", result2.Error)
		}
		if result2.Message != "" {
			fmt.Printf("消息: %s\n", result2.Message)
		}

		// 显示恢复步骤
		if len(result2.Steps) > 0 {
			fmt.Println("\n错误恢复步骤:")
			errorSteps := 0
			for j, step := range result2.Steps {
				if step.StepType == "error" || step.StepType == "accept" {
					fmt.Printf("%d. %s - %s\n", j+1, step.StepType, step.Description)
					if step.StepType == "error" {
						errorSteps++
					}
				}
			}
			fmt.Printf("总共恢复了 %d 个错误\n", errorSteps)
		}

		fmt.Println(strings.Repeat("=", 50))
	}

	fmt.Println("\n错误恢复测试完成!")
}

func getResultString(accepted bool) string {
	if accepted {
		return "接受"
	}
	return "拒绝"
}
