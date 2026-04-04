package model

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
	From     State `json:"from"`
	To       State `json:"to"`
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
