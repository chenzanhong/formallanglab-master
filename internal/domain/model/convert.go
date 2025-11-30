package model

import (
	"fmt"
	"strings"
)

// 转换功能相关的模型定义

// ========== 文法转自动机 ==========

type GrammarToFAStep struct {
	Production          *Production  `json:"production,omitempty"`
	Action              string       `json:"action"`
	Description         string       `json:"description"`
	NewStates           []State      `json:"newStates,omitempty"`
	NewTransitions      []Transition `json:"newTransitions,omitempty"`
	AcceptingStateAdded *State       `json:"acceptingStateAdded,omitempty"`
}

type GrammarToFAProcess struct {
	Linear          GrammarLinearity  `json:"linear"`
	OriginalGrammar Grammar           `json:"originalGrammar"`
	Steps           []GrammarToFAStep `json:"steps"`
	FinalAutomaton  *Automaton        `json:"finalAutomaton"`
}

// ========== FA 转 Regex ===========
type PathUpdate struct {
	From     State  `json:"from"`
	To       State  `json:"to"`
	OldRegex Regex `json:"oldRegex"`
	NewPart  Regex `json:"newPart"`
	NewRegex Regex `json:"newRegex"`
}

type ConversionStep struct {
	EliminatedState State               `json:"eliminatedState"`
	UpdatedPaths    []PathUpdate        `json:"updatedPaths"`
	AutomatonFlow   *ReactFlowAutomaton `json:"automatonFlow"`
}

type ConversionProcess struct {
	Steps        []ConversionStep `json:"steps"`
	FinalRegex   Regex            `json:"finalRegex"`
	InitialState State            `json:"initialState"`
	FinalState   State            `json:"finalState"`
}

// ========== Regex 转 FA ===========
type RegexToFAStep struct {
	Expr          string              `json:"expr"` // 当前子表达式字符串（如 "a", "(b|c)*"）
	StartState    State               `json:"startState"`
	EndState      State               `json:"endState"`
	AutomatonFlow *ReactFlowAutomaton `json:"automatonFlow"` // 当前完整的 FA 快照
}

type RegexToFAProcess struct {
	Regex          Regex           `json:"regex"`
	Steps          []RegexToFAStep `json:"steps"`
	FinalAutomaton *Automaton      `json:"finalAutomaton"`
}

// 使用指针版本：FA *Automaton + Steps []*RegexToFAStep，但在生成每一步时 显式深拷贝 自动机。

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
	// fmt.Printf("  状态: %v\n", cs.CurrentAutomaton.States)
	// fmt.Println("  转移:")
	// for _, t := range cs.CurrentAutomaton.Transitions {
	// 	fmt.Printf("    %s ──%s──> %s\n", t.FromState, t.Input, t.ToStates[0])
	// }
	fmt.Println(strings.Repeat("─", 50))
}
