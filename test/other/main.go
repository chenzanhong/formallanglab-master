package main

import (
	"fmt"

	"backend/internal/domain/model"
	"backend/internal/service/automaton_s"
	"backend/internal/service/grammar_s"
	"backend/internal/service/regex_s"
)

var grammar = &model.Grammar{ // a+|ba*
	StartSymbol:  model.Symbol("S"),
	NonTerminals: []model.Symbol{"S", "A", "B"},
	Terminals:    []model.Symbol{"a", "b"},
	Productions: []model.Production{
		{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
		{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b", "B"}},
		{Left: []model.Symbol{"A"}, Right: []model.Symbol{"B"}},
		{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a", "A"}},
		{Left: []model.Symbol{"B"}, Right: []model.Symbol{model.Epsilon}},
	},
	GrammarType: model.RegularGrammar,
}

var automaton = &model.Automaton{
	States:          []model.State{"q0", "q1", "q2", "q3"},
	InitialState:    "q0",
	Alphabet:        []model.Symbol{"a", "b", "c"},
	AcceptingStates: []model.State{"q3"},
	Transitions: []model.Transition{
		{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}},
		{FromState: "q1", Input: "b", ToStates: []model.State{"q2"}},
		{FromState: "q1", Input: model.Epsilon, ToStates: []model.State{"q2"}},
		{FromState: "q2", Input: "a", ToStates: []model.State{"q2"}},
		{FromState: "q2", Input: "c", ToStates: []model.State{"q3"}},
		{FromState: "q3", Input: "a", ToStates: []model.State{"q3"}},
	},
}

var regex = model.Regex("a(b|c)*a+")

func main() {
	testGrammarRecognizeString()
	// testGrammarGenerateExampleString()
	// testAutomatonGenerateExampleString()
	// testRegexGenerateExampleString()
}

func testGrammarRecognizeString() {
	// 定义要测试的字符串和模式
	testString := "ba"
	modes := []struct {
		mode grammar_s.Mode
		name string
	}{
		{grammar_s.RecursiveDescentMode, "递归下降分析"},
		{grammar_s.BFSMode, "BFS分析"},
		{grammar_s.LL1Mode, "LL(1)分析"},
		{grammar_s.LR0Mode, "LR(0)分析"},
	}

	fmt.Println("===== 语法识别测试 =====")
	fmt.Printf("测试字符串: %s\n\n", testString)

	// 使用不同模式进行测试
	for _, m := range modes {
		fmt.Printf("【%s】\n", m.name)
		result := grammar_s.ParseStringWithMode(grammar, testString, m.mode, true)

		// 格式化输出结果
		fmt.Printf("是否接受: %v\n", result.Accepted)
		fmt.Printf("分析方法: %s\n", result.Method)
		if result.Message != "" {
			fmt.Printf("信息: %s\n", result.Message)
		}
		if result.Error != "" {
			fmt.Printf("错误信息: %s\n", result.Error)
		}

		// 输出解析步骤（如果有）
		if len(result.Steps) > 0 {
			fmt.Println("解析步骤:")
			for i, step := range result.Steps {
				fmt.Printf("  %2d. 符号栈: %v,  输入串: %v, 产生式：%v，动作: %s\n",
					i+1, step.Stack, step.Input, step.Production, step.Action)
			}
		}
		fmt.Println()
	}
	fmt.Println("========================")
}

func testGrammarGenerateExampleString() {
	accept, reject := grammar_s.GrammarGenerateExampleString(grammar)
	for _, s := range accept {
		fmt.Println(s)
	}
	fmt.Println("---")
	for _, s := range reject {
		fmt.Println(s)
	}
}

func testAutomatonGenerateExampleString() {
	ss, sj := automaton_s.AutomatonGenerateExampleString(automaton)
	for _, s := range ss {
		fmt.Println(s)
	}
	fmt.Println("---")
	for _, s := range sj {
		fmt.Println(s)
	}
}

func testRegexGenerateExampleString() {
	accept, reject := regex_s.RegexGenerateExampleString(regex)
	for _, s := range accept {
		fmt.Println(s)
	}
	fmt.Println("---")
	for _, s := range reject {
		fmt.Println(s)
	}
}
