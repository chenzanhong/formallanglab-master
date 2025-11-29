package model

// =================== 自动机识别字符串的过程记录 ===================
type RecognitionStep struct {
	Step      int    `json:"step"`
	State     State  `json:"state"`
	Input     Symbol `json:"input"`
	NextState State  `json:"nextState"`
}

type RecognitionResult struct {
	IsAccepted bool              `json:"isAccepted"`
	Steps      []RecognitionStep `json:"steps"`
}

// ===================	  DFA最小化的过程记录    ===================
// MinimizationStep 表示DFA最小化过程中的一个步骤
type MinimizationStep struct {
	Step          int                 `json:"step"`      // 步骤编号
	Partition     [][]State           `json:"partition"` // 当前的状态划分
	Actions       []string            `json:"actions"`   // 执行的操作描述
	AutomatonFlow *ReactFlowAutomaton `json:"automatonFlow"`
}

// MinimizationProcess 表示DFA最小化的完整过程记录
type MinimizationProcess struct {
	Steps []MinimizationStep `json:"steps"` // 最小化步骤列表
}

// ===================	  NFA最小化的过程记录    ===================
type NFADeterminizationStep struct {
	Step          int                 `json:"step"`        // 步骤编号
	Description   string              `json:"description"` // 步骤描述
	AutomatonFlow *ReactFlowAutomaton `json:"automatonFlow"`
}
type NFADeterminizationProcess struct {
	FinalAutomaton *Automaton               `json:"finalAutomaton"`
	Steps          []NFADeterminizationStep `json:"steps"` // 最小化步骤列表
}
