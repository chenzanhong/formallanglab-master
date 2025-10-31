package model

import (
	"errors"
	"fmt"
)

/* 用户 */
type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Email    string `json:"email"`
}

// Symbol 表示一个符号，可以是终结符、非终结符或者自动机所识别的一个符号
type Symbol string

const Epsilon Symbol = "ε" // 定义ε作为特殊输入符号，表示空转移符号

// 工具函数：将用户输入映射为标准 ε
func NormalizeSymbol(s string) Symbol {
	switch s {
	case "ε", "epsilon", "e", "E", "", "λ", "eps":
		return Epsilon
	default:
		return Symbol(s)
	}
}

/* 文法 */
// Production 规定了一个产生式
type Production struct {
	Left  []Symbol `json:"left"`  // 左部，通常是单个非终结符，但也可以是多个符号
	Right []Symbol `json:"right"` // 右部，可以包含多个符号
}

// Grammar 表示整个文法
type Grammar struct {
	StartSymbol  Symbol       `json:"startSymbol"`  // 起始符号
	Terminals    []Symbol     `json:"terminals"`    // 终结符集合
	NonTerminals []Symbol     `json:"nonTerminals"` // 非终结符集合
	Productions  []Production `json:"productions"`  // 产生式集合
}

func (g *Grammar) CheckIsTerminal(s Symbol) bool {
	for _, v := range g.Terminals {
		if v == s {
			return true
		}
	}
	return false
}

func (g *Grammar) CheckIsNonTerminal(s Symbol) bool {
	for _, v := range g.NonTerminals {
		if v == s {
			return true
		}
	}
	return false
}

/* 自动机 */
// State 表示自动机中的一个状态
type State string

// Transition 表示一个状态转移规则
type Transition struct {
	FromState State   `json:"fromState"` // 起始状态
	Input     Symbol  `json:"input"`     // 输入符号
	ToStates  []State `json:"toStates"`  // 目标状态（对于NFA可以有多个）
}

// Automaton 基础自动机结构
type Automaton struct {
	States          []State                      `json:"states"`          // 状态集合
	Alphabet        []Symbol                     `json:"alphabet"`        // 符号表
	Transitions     []Transition                 `json:"transitions"`     // 状态转移规则集合
	InitialState    State                        `json:"initialState"`    // 初始状态
	AcceptingStates []State                      `json:"acceptingStates"` // 接受状态集合
	IsDFA           bool                         `json:"isDFA"`           // 是否为DFA，否则为NFA
	TransMap        map[State]map[Symbol][]State // Map存储状态转移规则，识别字符串时效率高
}

var ValidCSet []byte // 正则表达式支持的符合，包括0~1，a~z，A~Z，|，（，），*，？，·，

type Regex struct {
	Patten string
	CSet   []byte
}

// 按字符集以及最长匹配原则切分字符串
func (a *Automaton) SplitString(s string) ([]Symbol, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("输入字符串为空")
	}
	var mp = make(map[string]bool, len(a.Alphabet))
	for _, sym := range a.Alphabet {
		mp[string(sym)] = true
	}
	var res []Symbol
	for i, j := 0, 1; i < len(s) && j <= len(s); {
		if mp[s[i:j]] {
			j++ // 尝试更长的匹配
		} else {
			if i == j-1 { // 连一个字符都匹配不上
				return nil, fmt.Errorf("输入字符串包含无效字符")
			} else if j > i+1 {
				res = append(res, Symbol(s[i:j-1]))
				i = j - 1
				j++
			}
		}
	}
	return res, nil
}

// CheckIsDFA 检查当前自动机是否真的是一个有效的 DFA
// 如果是，则返回 true，并将 IsDFA 设为 true
// 否则返回 false，并将 IsDFA 设为 false
func (a *Automaton) CheckIsDFA() (bool, error) {
	// 构建转移映射：fromState + input -> toState（用于检查完备性和唯一性）
	transitionMap := make(map[string]State) // key: state|symbol

	// 遍历所有转移
	for _, t := range a.Transitions {
		if t.Input == Epsilon {
			return false, fmt.Errorf("包含空转移，不是DFA")
		}
		// 检查是否重复定义了同一 (fromState, input)
		key := string(t.FromState) + "|" + string(t.Input)
		if prev, exists := transitionMap[key]; exists {
			a.IsDFA = false
			return false, fmt.Errorf("DFA 中状态 '%s' 对输入 '%s' 定义了多个转移（已存在: %s）", t.FromState, t.Input, prev) // 重复定义
		}
		transitionMap[key] = t.ToStates[0]
	}
	// 所有条件满足，是 DFA
	a.IsDFA = true
	return true, nil
}

func (a *Automaton) InitTransMap() {
	a.TransMap = make(map[State]map[Symbol][]State)
	for _, t := range a.Transitions {
		if _, exists := a.TransMap[t.FromState]; !exists {
			a.TransMap[t.FromState] = make(map[Symbol][]State)
		}
		if len(t.ToStates) > 0 {
			a.TransMap[t.FromState][t.Input] = append(a.TransMap[t.FromState][t.Input], t.ToStates...)
		}
	}
}

func (a *Automaton) ISValidate() (bool, error) {
	// 使用 map 提高查找效率
	stateSet := make(map[State]bool)
	alphabetSet := make(map[Symbol]bool)

	// 1. 检查状态集合不能为空
	if len(a.States) == 0 {
		return false, errors.New("状态集合不能为空")
	}

	// 构建状态集合
	for _, s := range a.States {
		if s == "" {
			return false, errors.New("状态名不能为空字符串")
		}
		stateSet[s] = true
	}

	// 2. 检查初始状态是否在状态集合中
	if !stateSet[a.InitialState] {
		return false, fmt.Errorf("初始状态 '%s' 不在状态集合中", a.InitialState)
	}

	// 3. 检查接受状态是否都是合法状态
	for _, acc := range a.AcceptingStates {
		if !stateSet[acc] {
			return false, fmt.Errorf("接受状态 '%s' 不在状态集合中", acc)
		}
	}

	// 4. 构建字母表集合
	for _, sym := range a.Alphabet {
		alphabetSet[sym] = true
	}
	alphabetSet[Epsilon] = true // 允许空转移

	// 5. 验证所有转移规则，检查所有转移符合和状态是否有定义，顺带确定isDFA
	a.IsDFA = true
	transitionMap := make(map[string]State) // key: state|symbol
	for _, t := range a.Transitions {
		// 5.1 检查起始状态是否已定义
		if !stateSet[t.FromState] {
			return false, fmt.Errorf("转移规则'%v'中起始状态 '%s' 未定义", t, t.FromState)
		}
		// 5.2 检查输入符号是否属于字母表（允许 ε）
		if t.Input != Epsilon && !alphabetSet[t.Input] {
			return false, fmt.Errorf("转移规则'%v'中输入符号 '%s' 不属于字母表", t, t.Input)
		}
		// 5.3 检查目标状态是否已定义
		for _, to := range t.ToStates {
			if !stateSet[to] {
				return false, fmt.Errorf("转移规则'%v'中目标状态 '%s' 未定义", t, to)
			}
		}
		if t.Input == Epsilon {
			a.IsDFA = false
		} else if a.IsDFA {
			// 检查是否重复定义了同一 (fromState, input)
			if len(t.ToStates) > 1 {
				a.IsDFA = false
				continue
			}
			key := string(t.FromState) + "|" + string(t.Input)
			if to, exists := transitionMap[key]; exists && to != t.ToStates[0] {
				a.IsDFA = false
				continue
			}
			transitionMap[key] = t.ToStates[0]
		}
	}
	return true, nil
}

// CompleteDFA 将自动机转换为等价的完备 DFA
// 前提：调用者应确保 a 是一个有效的 DFA（可通过 a.CheckIsDFA() 验证）
// 步骤：
//  1. 移除不可达状态
//  2. 添加陷阱状态（若需要）
//  3. 补全所有缺失的转移
func (a *Automaton) CompleteDFA() error {
	// Step 0: 确保是 DFA
	if !a.IsDFA {
		if ok, err := a.CheckIsDFA(); !ok {
			return fmt.Errorf("cannot complete non-DFA: %w", err)
		}
	}

	// Step 1: 找出所有可达状态
	reachable := a.FindReachableStates()

	// 构建可达状态集合（用于快速查找）
	reachableSet := make(map[State]bool)
	for _, s := range reachable {
		reachableSet[s] = true
	}

	// 过滤状态、接受状态、转移
	newStates := reachable
	newAccepting := []State{}
	for _, s := range a.AcceptingStates {
		if reachableSet[s] {
			newAccepting = append(newAccepting, s)
		}
	}

	newTransitions := []Transition{}
	for _, t := range a.Transitions {
		if reachableSet[t.FromState] {
			// 只保留起点可达的转移
			newTransitions = append(newTransitions, t)
		}
	}

	// 更新自动机
	a.States = newStates
	a.AcceptingStates = newAccepting
	a.Transitions = newTransitions

	// Step 2: 构建当前转移映射（用于检查缺失）
	transMap := make(map[State]map[Symbol]State)
	for _, t := range a.Transitions {
		if t.Input == Epsilon {
			return fmt.Errorf("unexpected epsilon in DFA")
		}
		if _, ok := transMap[t.FromState]; !ok {
			transMap[t.FromState] = make(map[Symbol]State)
		}
		transMap[t.FromState][t.Input] = t.ToStates[0] // DFA only
	}

	// Step 3: 创建陷阱状态（仅当需要时）
	sinkState := State("__sink__")
	hasSink := false
	missingTransitions := []Transition{}

	// 遍历每个可达状态和每个字母表符号
	for _, state := range a.States {
		for _, sym := range a.Alphabet {
			if sym == Epsilon {
				continue // DFA 不应有 ε
			}
			if _, exists := transMap[state][sym]; !exists {
				// 缺失转移：指向 sink
				missingTransitions = append(missingTransitions, Transition{
					FromState: state,
					Input:     sym,
					ToStates:  []State{sinkState},
				})
				hasSink = true
			}
		}
	}

	// 如果有缺失转移，添加 sink 状态及其自环
	if hasSink {
		// 添加 sink 状态
		a.States = append(a.States, sinkState)

		// sink 对所有输入自环
		for _, sym := range a.Alphabet {
			if sym == Epsilon {
				continue
			}
			a.Transitions = append(a.Transitions, Transition{
				FromState: sinkState,
				Input:     sym,
				ToStates:  []State{sinkState},
			})
		}

		// 添加缺失的转移
		a.Transitions = append(a.Transitions, missingTransitions...)
	}

	// 重新构建 TransMap（供后续识别使用）
	a.InitTransMap()
	a.IsDFA = true // 仍为 DFA
	return nil
}

// findReachableStates 使用 BFS 找出从初始状态可达的所有状态
func (a *Automaton) FindReachableStates() []State {
	visited := make(map[State]bool)
	queue := []State{a.InitialState}
	visited[a.InitialState] = true

	// 构建邻接表
	adj := make(map[State][]State)
	for _, t := range a.Transitions {
		if t.Input == Epsilon {
			// 虽然 DFA 不应有 ε，但为健壮性考虑
			for _, to := range t.ToStates {
				adj[t.FromState] = append(adj[t.FromState], to)
			}
		} else {
			for _, to := range t.ToStates {
				adj[t.FromState] = append(adj[t.FromState], to)
			}
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, next := range adj[current] {
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}

	// 收集所有 visited 状态
	var reachable []State
	for _, state := range a.States {
		if visited[state] {
			reachable = append(reachable, state)
		}
	}
	return reachable
}

type ReactFlowNode struct {
	ID   string `json:"id"`
	Type string `json:"type"` // "initial" | "default"
	Data struct {
		Label       string `json:"label"`
		IsAccepting bool   `json:"isAccepting"`
	} `json:"data"`
}

type ReactFlowEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
	// Style     map[string]string `json:"style,omitempty"`
	// MarkerEnd struct {
	// 	Type string `json:"type"` // "arrow"
	// } `json:"markerEnd"`
}

type ReactFlowAutomaton struct {
	Nodes []ReactFlowNode `json:"nodes"`
	Edges []ReactFlowEdge `json:"edges"`
}

// 转换函数
func (a *Automaton) ToReactFlow() *ReactFlowAutomaton {
	nodes := make([]ReactFlowNode, len(a.States))
	for i, state := range a.States {
		isInitial := state == a.InitialState
		isAccepting := false
		for _, accState := range a.AcceptingStates {
			if state == accState {
				isAccepting = true
				break
			}
		}
		nodes[i] = ReactFlowNode{
			ID: string(state),
			Type: func() string {
				if isInitial {
					return "initial"
				}
				return "default"
			}(),
			Data: struct {
				Label       string `json:"label"`
				IsAccepting bool   `json:"isAccepting"`
			}{
				Label:       string(state),
				IsAccepting: isAccepting,
			},
		}
	}

	edges := make([]ReactFlowEdge, 0)
	edgeID := 1
	for _, t := range a.Transitions {
		for _, toState := range t.ToStates {
			label := string(t.Input)
			if t.Input == Epsilon {
				label = "ε"
			}
			edge := ReactFlowEdge{
				ID:     fmt.Sprintf("e%d", edgeID),
				Source: string(t.FromState),
				Target: string(toState),
				Label:  label,
				// MarkerEnd: struct {
				// 	Type string `json:"type"`
				// }{Type: "arrow"},
			}
			// if t.Input == Epsilon {
			// 	edge.Style = map[string]string{"stroke": "#1890ff"}
			// }
			edges = append(edges, edge)
			edgeID++
		}
	}

	return &ReactFlowAutomaton{Nodes: nodes, Edges: edges}
}

// ParseStep 表示文法分析过程中的一个步骤
type ParseStep struct {
	StepType    string      `json:"stepType"`    // "init", "match", "predict", "accept", "error"
	Description string      `json:"description"` // 步骤描述
	Stack       []Symbol    `json:"stack"`       // 当前栈状态
	Input       []Symbol    `json:"input"`       // 剩余输入
	InputPos    int         `json:"inputPos"`    // 输入指针位置
	Action      string      `json:"action"`      // 执行的动作描述
	Production  *Production `json:"production"`  // 使用的产生式（如果有）
}

// ParseResult 表示文法分析的完整结果
type ParseResult struct {
	Accepted bool        `json:"accepted"` // 是否接受输入串
	Method   string      `json:"method"`   // 使用的分析方法
	Steps    []ParseStep `json:"steps"`    // 分析步骤（可选）
	Error    string      `json:"error"`    // 错误信息（如果有）
	Message  string      `json:"message"`  // 额外信息
}

// ToReactFlow 返回格式化后的产生式字符串，用于前端展示
func (g *Grammar) ToReactFlow() []string {
	productions := make([]string, 0, len(g.Productions))

	for _, prod := range g.Productions {
		// 格式化左部
		left := string(prod.Left[0]) // 通常左部只有一个符号

		// 格式化右部
		right := ""
		if len(prod.Right) == 0 || (len(prod.Right) == 1 && prod.Right[0] == Epsilon) {
			right = string(Epsilon)
		} else {
			for i, symbol := range prod.Right {
				if i > 0 {
					right += " "
				}
				if symbol == Epsilon {
					right += string(Epsilon)
				} else {
					right += string(symbol)
				}
			}
		}

		// 组合成产生式字符串
		productions = append(productions, left+" -> "+right)
	}

	return productions
}

/*

前端传给后端的文法格式：：
{
  "startSymbol": "S",
  "terminals": ["a", "b", "ε"],
  "nonTerminals": ["S", "A", "B"],
  "productions": [
    {
      "left": "S",
      "right": ["a", "A", "b"]
    },
    {
      "left": "A",
      "right": ["b", "B"]
    },
    {
      "left": "B",
      "right": ["ε"]
    }
  ]
}

前端传给后端的自动机格式：
{
  "states": ["q0", "q1", "q2"],
  "alphabet": ["a", "b", "ε"],
  "transitions": [
    {
      "fromState": "q0",
      "input": "a",
      "toStates": ["q1"]
    },
    {
      "fromState": "q1",
      "input": "b",
      "toStates": ["q2"]
    },
    {
      "fromState": "q1",
      "input": "ε",
      "toStates": ["q0"]
    }
  ],
  "initialState": "q0",
  "acceptingStates": ["q2"],
  "isDFA": false
}


返回给前端的自动机数据，JSON格式，点和边：
    1. id对应状态
    2. "type": "input",对应初始状态
    3. "type": "default",对应非开始状态
    4. label为点上的状态名称或边上的符号
{
  "nodes": [
    {
    	"id": "q0",
      	"type": "initial", // 明确初始状态类型
      	"data": {
        "label": "q₀",
        "isAccepting": false // 显式标记是否为接受状态，接收状态需要前端用双圆圈表示，其他节点均用圆形节点表示（节点内部展示label）
      }
    },
    {
      "id": "q1",
      "type": "default",
      "data": {
        "label": "q₁",
        "isAccepting": true
      }
    }
  ],
  "edges": [
    {
      "id": "e1", // 由后端生成唯一ID
      "source": "q0",
      "target": "q1",
      "label": "a",
    //   "markerEnd": { "type": "arrow" } // ReactFlow 的箭头标记
    },
    {
      "id": "e2",
      "source": "q1",
      "target": "q0",
      "label": "ε",
    //   "style": { "stroke": "#1890ff" }, // 空转移，使用蓝色边
    //   "markerEnd": { "type": "arrow" }
    }
  ]
}
*/
