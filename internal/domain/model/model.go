package model

import "fmt"

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

// const Epsilon Symbol = "ε" // 定义ε作为特殊输入符号
const Epsilon Symbol = "" // 定义""作为特殊输入符号，表示ε，前端传输时用""表示就行，毕竟ε不属于ASCII，属于Unicode

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
	States          []State      `json:"states"`          // 状态集合
	Alphabet        []Symbol     `json:"alphabet"`        // 符号表
	Transitions     []Transition `json:"transitions"`     // 状态转移规则集合
	InitialState    State        `json:"initialState"`    // 初始状态
	AcceptingStates []State      `json:"acceptingStates"` // 接受状态集合
	IsDFA           bool         `json:"isDFA"`           // 是否为DFA，否则为NFA
}

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

type ReactFlowNode struct {
	ID   string `json:"id"`
	Type string `json:"type"` // "initial" | "default"
	Data struct {
		Label       string `json:"label"`
		IsAccepting bool   `json:"isAccepting"`
	} `json:"data"`
}

type ReactFlowEdge struct {
	ID        string            `json:"id"`
	Source    string            `json:"source"`
	Target    string            `json:"target"`
	Label     string            `json:"label"`
	Style     map[string]string `json:"style,omitempty"`
	MarkerEnd struct {
		Type string `json:"type"` // "arrow"
	} `json:"markerEnd"`
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
				MarkerEnd: struct {
					Type string `json:"type"`
				}{Type: "arrow"},
			}
			if t.Input == Epsilon {
				edge.Style = map[string]string{"stroke": "#1890ff"}
			}
			edges = append(edges, edge)
			edgeID++
		}
	}

	return &ReactFlowAutomaton{Nodes: nodes, Edges: edges}
}

/*
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
        "isAccepting": false // 显式标记是否为接受状态，接收状态需要前端用双鱼圈表示，其他节点均用圆形节点表示（节点内部展示label）
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
      "markerEnd": { "type": "arrow" } // ReactFlow 的箭头标记
    },
    {
      "id": "e2",
      "source": "q1",
      "target": "q0",
      "label": "ε",
      "style": { "stroke": "#1890ff" }, // 蓝色边
      "markerEnd": { "type": "arrow" }
    }
  ]
}
*/

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
  "isDFA": false
}

文法前端传输格式：
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
*/
