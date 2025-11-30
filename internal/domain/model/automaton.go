package model

import (
	"errors"
	"fmt"
)

/* 自动机 */
// Symbol 表示自动机中的输入符号
// type Symbol string

// State 表示自动机中的一个状态
type State string

// const Epsilon Symbol = "ε" // 定义ε作为特殊输入符号，表示空转移符号

const (
	SinkState          State = "__sink__"  // 陷阱状态
	UniqueInitialState State = "_initial_" // 唯一初始状态
	UniqueFinalState   State = "_final_"   // 唯一接受状态
)

// Transition 表示一个状态转移规则
type Transition struct {
	FromState State   `json:"fromState"` // 起始状态
	Input     Symbol  `json:"input"`     // 输入符号
	ToStates  []State `json:"toStates"`  // 目标状态（对于NFA可以有多个，但是目前大部分还是分开来的，不合并相同FromState+Input的产生式）
}

type AutomatonType int

const (
	DFA        AutomatonType = 0 // 确定有限自动机
	NFA        AutomatonType = 1 // 非确定有限自动机
	EpsilonNFA AutomatonType = 2 // 非确定有限自动机（允许 ε-转移）
)

// Automaton 基础自动机结构
type Automaton struct {
	States          []State                      `json:"states"`                    // 状态集合
	Alphabet        []Symbol                     `json:"alphabet"`                  // 符号表
	Transitions     []Transition                 `json:"transitions"`               // 状态转移规则集合
	InitialState    State                        `json:"initialState,omitempty"`    // 初始状态
	AcceptingStates []State                      `json:"acceptingStates,omitempty"` // 接受状态集合
	Type            AutomatonType                `json:"type"`                      // 类型
	TransMap        map[State]map[Symbol][]State `json:"-"`                         // Map存储状态转移规则，识别字符串时效率高；不参与 JSON 序列化
}

func NewEmptyLanguageAutomaton() *Automaton {
	return &Automaton{
		States:          []State{"q0"},
		Alphabet:        []Symbol{},
		Transitions:     []Transition{},
		InitialState:    "q0",
		AcceptingStates: []State{},
		Type:            DFA,
	}
}

// 按字符集以及最长匹配原则切分字符串，返回可被正确切分的符号序列
func (a *Automaton) SplitString(s string) ([]Symbol, error) {
	if len(s) == 0 || s == "" {
		return nil, errors.New("输入字符串为空")
	}

	// 构建字母表集合 + 计算最大符号长度（用于剪枝）
	symSet := make(map[string]bool, len(a.Alphabet))
	maxLen := 0
	for _, sym := range a.Alphabet {
		strSym := string(sym)
		symSet[strSym] = true
		if len(strSym) > maxLen {
			maxLen = len(strSym)
		}
	}

	var res []Symbol
	i := 0
	n := len(s)

	for i < n {
		// 确定本次最多尝试到哪：不能超过字符串末尾，也不能超过 maxLen
		end := i + maxLen
		if end > n {
			end = n
		}

		matched := false
		// 从最长可能子串开始，向短尝试（贪心最长匹配）
		for j := end; j > i; j-- {
			sub := s[i:j]
			if symSet[sub] {
				res = append(res, Symbol(sub))
				i = j // 跳到匹配结束位置
				matched = true
				break
			}
		}

		if !matched {
			// 找不到任何以 s[i] 开头的合法符号
			return res, fmt.Errorf("在位置 %d 无法匹配任何符号，剩余字符串: '%s'", i, s[i:])
		}
	}

	return res, nil
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

// 简单判断自动机结构是否正确，并进行不完整的类型判断（未考虑不可达状态）
func (a *Automaton) Validate() error {
	// 使用 map 提高查找效率
	stateSet := make(map[State]bool)
	alphabetSet := make(map[Symbol]bool)
	// 1. 检查状态集合不能为空
	if len(a.States) == 0 {
		return errors.New("状态集合不能为空")
	}

	// 构建状态集合
	for _, s := range a.States {
		if s == "" {
			return errors.New("状态名不能为空字符串")
		}
		stateSet[s] = true
	}

	// 2. 检查初始状态是否在状态集合中
	if !stateSet[a.InitialState] {
		return fmt.Errorf("初始状态 '%s' 不在状态集合中", a.InitialState)
	}

	// 3. 检查接受状态是否都是合法状态
	for _, acc := range a.AcceptingStates {
		if !stateSet[acc] {
			return fmt.Errorf("接受状态 '%s' 不在状态集合中", acc)
		}
	}

	// 4. 构建字母表集合
	for _, sym := range a.Alphabet {
		alphabetSet[sym] = true
	}
	alphabetSet[Epsilon] = true // 允许空转移

	// 5. 验证所有转移规则，检查所有转移符合和状态是否有定义，顺带确定type
	a.Type = DFA
	transitionMap := make(map[string]State) // key: state|symbol
	for _, t := range a.Transitions {
		// 5.1 检查起始状态是否已定义
		if !stateSet[t.FromState] {
			return fmt.Errorf("转移规则'%v'中起始状态 '%s' 未定义", t, t.FromState)
		}
		// 5.2 检查输入符号是否属于字母表（允许 ε）
		if t.Input != Epsilon && !alphabetSet[t.Input] {
			return fmt.Errorf("转移规则'%v'中输入符号 '%s' 不属于字母表", t, t.Input)
		}
		// 5.3 检查目标状态是否已定义
		for _, to := range t.ToStates {
			if !stateSet[to] {
				return fmt.Errorf("转移规则'%v'中目标状态 '%s' 未定义", t, to)
			}
		}
		// 5.4 判断类型
		if t.Input == Epsilon {
			a.Type = EpsilonNFA
		} else if a.Type == DFA {
			// 检查是否重复定义了同一 (fromState, input)
			if len(t.ToStates) > 1 {
				a.Type = NFA
				continue
			}
			key := string(t.FromState) + "|" + string(t.Input)
			if to, exists := transitionMap[key]; exists && to != t.ToStates[0] {
				a.Type = NFA
				continue
			}
			transitionMap[key] = t.ToStates[0]
		}
	}
	return nil
}

// Clone 返回 Automaton 的深拷贝
func (a *Automaton) Clone() *Automaton {
	if a == nil {
		return nil
	}

	clone := &Automaton{
		States:          make([]State, len(a.States)),
		Alphabet:        make([]Symbol, len(a.Alphabet)),
		Transitions:     make([]Transition, len(a.Transitions)),
		InitialState:    a.InitialState,
		AcceptingStates: make([]State, len(a.AcceptingStates)),
		Type:            a.Type,
		// TransMap 不复制，因为它是运行时缓存
	}

	copy(clone.States, a.States)
	copy(clone.Alphabet, a.Alphabet)
	copy(clone.AcceptingStates, a.AcceptingStates)

	for i, t := range a.Transitions {
		toStates := make([]State, len(t.ToStates))
		copy(toStates, t.ToStates)
		clone.Transitions[i] = Transition{
			FromState: t.FromState,
			Input:     t.Input,
			ToStates:  toStates,
		}
	}

	return clone
}

// CompleteDFA 将自动机原地转换为等价的完备 DFA
// 前提：调用者应确保 a 是一个有效的 DFA（可通过 a.Validate() 验证）
// 步骤：
//  1. 移除不可达状态
//  2. 清理字母表（移除 ε）
//  3. 添加陷阱状态（若需要）
//  4. 补全所有缺失的转移
func (a *Automaton) CompleteDFA() error {
	if err := a.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if a.Type != DFA {
		return fmt.Errorf("only DFA can be completed")
	}

	// Step 1: 找出可达状态并过滤
	reachable := a.FindReachableStates()
	reachableSet := make(map[State]bool)
	for _, s := range reachable {
		reachableSet[s] = true
	}

	// 过滤接受状态
	newAccepting := []State{}
	for _, s := range a.AcceptingStates {
		if reachableSet[s] {
			newAccepting = append(newAccepting, s)
		}
	}

	// 过滤转移（同时排除 ε）
	newTransitions := []Transition{}
	for _, t := range a.Transitions {
		if reachableSet[t.FromState] && t.Input != Epsilon {
			newTransitions = append(newTransitions, t)
		}
	}

	// 清理字母表：移除 ε（DFA 不应包含）
	cleanAlphabet := []Symbol{}
	for _, sym := range a.Alphabet {
		if sym != Epsilon {
			cleanAlphabet = append(cleanAlphabet, sym)
		}
	}

	// Step 2: 构建转移映射（基于过滤后的转移）
	transMap := make(map[State]map[Symbol]State)
	for _, t := range newTransitions {
		if len(t.ToStates) != 1 {
			return fmt.Errorf("DFA transition must have exactly one target state")
		}
		if _, ok := transMap[t.FromState]; !ok {
			transMap[t.FromState] = make(map[Symbol]State)
		}
		transMap[t.FromState][t.Input] = t.ToStates[0]
	}

	// Step 3: 检查是否需要 sink 并收集缺失转移
	sinkState := SinkState
	hasSink := false
	maxMissing := len(reachable) * len(cleanAlphabet)
	missingTransitions := make([]Transition, 0, maxMissing)

	for _, state := range reachable {
		for _, sym := range cleanAlphabet {
			if _, exists := transMap[state][sym]; !exists {
				missingTransitions = append(missingTransitions, Transition{
					FromState: state,
					Input:     sym,
					ToStates:  []State{sinkState},
				})
				hasSink = true
			}
		}
	}

	// Step 4: 更新 a 的字段
	a.States = reachable
	a.AcceptingStates = newAccepting
	a.Alphabet = cleanAlphabet
	a.Transitions = newTransitions

	if hasSink {
		a.States = append(a.States, sinkState)
		// 添加 sink 自环
		for _, sym := range cleanAlphabet {
			a.Transitions = append(a.Transitions, Transition{
				FromState: sinkState,
				Input:     sym,
				ToStates:  []State{sinkState},
			})
		}
		a.Transitions = append(a.Transitions, missingTransitions...)
	}

	// 重建内部转移映射（如用于模拟或识别）
	a.InitTransMap()
	a.Type = DFA // 明确标记为 DFA（尽管本来就是）

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

// GNFA 表示广义非确定有限自动机，边标签为正则表达式
type GNFA struct {
	States       []State
	Transitions  []GNFATransition
	InitialState State
	FinalState   State
}

type GNFATransition struct {
	FromState State
	ToState   State
	Label     Regex // 注意：这里是 Regex，不是 Symbol
}

// =================== 自动机转换为 ReactFlow 格式 ===================
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

// ToReactFlow 将 GNFA 转换为 ReactFlow 可视化格式
func (g *GNFA) ToReactFlow() *ReactFlowAutomaton {
	// 构建节点
	nodeSet := make(map[State]bool)
	for _, s := range g.States {
		nodeSet[s] = true
	}

	nodes := make([]ReactFlowNode, 0, len(g.States))
	for _, state := range g.States {
		isInitial := state == g.InitialState
		isAccepting := state == g.FinalState

		node := ReactFlowNode{
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
		nodes = append(nodes, node)
	}

	// 构建边
	edges := make([]ReactFlowEdge, 0, len(g.Transitions))
	for idx, t := range g.Transitions {
		edge := ReactFlowEdge{
			ID:     fmt.Sprintf("gnfa_e%d", idx+1),
			Source: string(t.FromState),
			Target: string(t.ToState),
			Label:  string(t.Label), // 直接使用 Regex 字符串
		}
		edges = append(edges, edge)
	}

	return &ReactFlowAutomaton{
		Nodes: nodes,
		Edges: edges,
	}
}

// Automaton 转为 *ReactFlowAutomaton
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

// ReactFlowAutomaton 转为 *Automaton
func (r *ReactFlowAutomaton) ToAutomaton() *Automaton {
	var automaton Automaton

	// 处理节点集合
	for _, node := range r.Nodes {
		automaton.States = append(automaton.States, State(node.Data.Label)) // 状态集
		if node.Type == "initial" {                                         // 起始节点
			automaton.InitialState = State(node.Data.Label)
		}
		if node.Data.IsAccepting { // 接收状态
			automaton.AcceptingStates = append(automaton.AcceptingStates, State(node.Data.Label))
		}
	}

	var symbolMap = make(map[string]bool)
	//
	for _, edge := range r.Edges {
		if !symbolMap[edge.Label] && edge.Label != string(Epsilon) {
			automaton.Alphabet = append(automaton.Alphabet, Symbol(edge.Label)) // 符号集
			symbolMap[edge.Label] = true
		}
		automaton.Transitions = append(automaton.Transitions, Transition{ // 状态转移
			FromState: State(edge.Source),
			Input:     Symbol(edge.Label),
			ToStates:  []State{State(edge.Target)},
		}) // 优先单符号单边，不合并起始节点和终止节点相同的边
	}

	return &automaton
}

/*
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
  "type": 1
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
