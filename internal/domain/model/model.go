package model

import "fmt"

/* 用户 */
type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Name     string `json:"name"`
	Password string `json:"password"`
}


// Symbol 表示一个符号，可以是终结符、非终结符或者自动机所识别的一个符号
type Symbol string
const Epsilon Symbol = "ε" // 定义ε作为特殊输入符号

/* // 工具函数：将用户输入映射为标准 ε
func NormalizeSymbol(s string) Symbol {
    switch s {
    case "ε", "epsilon", "e", "E", "", "λ", "eps":
        return Epsilon
    default:
        return Symbol(s)
    }
} */

/* 文法 */
// Production 规定了一个产生式
type Production struct {
	Left  []Symbol `json:"left"`   // 左部，通常是单个非终结符，但也可以是多个符号
	Right []Symbol `json:"right"` // 右部，可以包含多个符号
}

// Grammar 表示整个文法
type Grammar struct {
	StartSymbol  Symbol              `json:"startSymbol"`  // 起始符号
	Terminals    map[Symbol]struct{} `json:"terminals"`    // 终结符集合
	NonTerminals map[Symbol]struct{} `json:"nonTerminals"` // 非终结符集合
	Productions  []Production               `json:"productions"`  // 产生式集合
}

/* 自动机 */
// State 表示自动机中的一个状态
type State string

// Transition 表示一个状态转移规则
type Transition struct {
	FromState State           `json:"fromState"` // 起始状态
	Input     Symbol `json:"input"`     // 输入符号
	ToStates  []State         `json:"toStates"`  // 目标状态（对于NFA可以有多个）
}

// Automaton 基础自动机结构
type Automaton struct {
	States          []State           `json:"states"`          // 状态集合
	Alphabet        []Symbol `json:"alphabet"`        // 符号表
	Transitions     []Transition      `json:"transitions"`     // 状态转移规则集合
	InitialState    State             `json:"initialState"`    // 初始状态
	AcceptingStates []State           `json:"acceptingStates"` // 接受状态集合
	isDFA           bool              `json:"isDFA"`           // 是否为DFA
}

// 节点，前端所需格式
type Node struct {
	ID    string `json:"id"`
	Type  string `json:"type"` // "input" for initial, "default" otherwise
	Data  struct {
		Label string `json:"label"`
	} `json:"data"`
}

/// 边，前端所需格式
type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
	Style  map[string]string `json:"style,omitempty"` // 可选样式，如 ε 用蓝色
}

type ReactFlowNode struct {
	ID    string `json:"id"`
	Type  string `json:"type"` // "initial" | "default"
	Data  struct {
		Label       string `json:"label"`
		IsAccepting bool   `json:"isAccepting"`
	} `json:"data"`
}

type ReactFlowEdge struct {
	ID       string            `json:"id"`
	Source   string            `json:"source"`
	Target   string            `json:"target"`
	Label    string            `json:"label"`
	Style    map[string]string `json:"style,omitempty"`
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
			ID:   string(state),
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