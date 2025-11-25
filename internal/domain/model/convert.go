package model

import (
	"fmt"
	"strings"
)

// 转换功能相关的模型定义

// ========== FA 转 Regex ===========
type PathUpdate struct {
	From     State  `json:"from"`
	To       State  `json:"to"`
	OldRegex Symbol `json:"oldRegex"`
	NewPart  Symbol `json:"newPart"`
	NewRegex Symbol `json:"newRegex"`
}

type ConversionStep struct {
	EliminatedState  State        `json:"eliminatedState"`
	UpdatedPaths     []PathUpdate `json:"updatedPaths"`
	CurrentAutomaton Automaton    `json:"currentAutomaton"`
}

type ConversionProcess struct {
	Steps        []ConversionStep `json:"steps"`
	FinalRegex   Regex            `json:"finalRegex"`
	InitialState State            `json:"initialState"`
	FinalState   State            `json:"finalState"`
}

// ========== Regex 转 FA ===========

// ========== 一些Print函数 （仅用于调试）==========
func (cs *ConversionStep) Print() {
	fmt.Printf("🔄 消除状态: %s\n", cs.EliminatedState)
	fmt.Println("📝 更新的路径:")
	if len(cs.UpdatedPaths) == 0 {
		fmt.Println("  (无)")
	}
	for _, up := range cs.UpdatedPaths {
		fmt.Printf("  %s ──[%s]──> %s\n", up.From, up.NewRegex, up.To)
		fmt.Printf("    原表达式: %s\n", up.OldRegex)
		fmt.Printf("    新增部分: %s\n", up.NewPart)
	}
	fmt.Println("📊 当前自动机状态:")
	fmt.Printf("  状态: %v\n", cs.CurrentAutomaton.States)
	fmt.Println("  转移:")
	for _, t := range cs.CurrentAutomaton.Transitions {
		fmt.Printf("    %s ──%s──> %s\n", t.FromState, t.Input, t.ToStates[0])
	}
	fmt.Println(strings.Repeat("─", 50))
}
