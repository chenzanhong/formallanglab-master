// regex_to_nfa.go
package convert_s

import (
	"backend/internal/domain/model" // 👈 替换为你的实际路径
	"fmt"
	"sort"
	"strings"
)

// Token 类型
type tokenType int

const (
	tokChar tokenType = iota
	tokStar
	tokPlus
	tokQuestion
	tokPipe
	tokLParen
	tokRParen
	tokEOF
)

type token struct {
	typ   tokenType
	value model.Symbol
}

// AST 节点（只用基本操作）
type astNode struct {
	typ   string // "char", "union", "concat", "star"
	val   model.Symbol
	left  *astNode
	right *astNode
	child *astNode
}

// ===== 1. 词法分析器（支持 + ?）=====
func lex(pattern string) []token {
	var tokens []token
	runes := []rune(pattern)
	i := 0
	for i < len(runes) {
		r := runes[i]
		switch r {
		case '*':
			tokens = append(tokens, token{typ: tokStar})
		case '+':
			tokens = append(tokens, token{typ: tokPlus})
		case '?':
			tokens = append(tokens, token{typ: tokQuestion})
		case '|':
			tokens = append(tokens, token{typ: tokPipe})
		case '(':
			tokens = append(tokens, token{typ: tokLParen})
		case ')':
			tokens = append(tokens, token{typ: tokRParen})
		case ' ', '\t', '\n':
			// 忽略空白
		default:
			// 只允许字母、数字作为字符
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				tokens = append(tokens, token{typ: tokChar, value: model.Symbol(string(r))})
			} else {
				panic(fmt.Sprintf("invalid character in regex: %c", r))
			}
		}
		i++
	}
	tokens = append(tokens, token{typ: tokEOF})
	return tokens
}

// ===== 2. 递归下降解析器（处理 + ?）=====
type parser struct {
	tokens []token
	pos    int
}

func (p *parser) peek() token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token{typ: tokEOF}
}

func (p *parser) consume() token {
	t := p.peek()
	p.pos++
	return t
}

func (p *parser) parseUnion() *astNode {
	left := p.parseConcat()
	for p.peek().typ == tokPipe {
		p.consume()
		right := p.parseConcat()
		left = &astNode{typ: "union", left: left, right: right}
	}
	return left
}

func (p *parser) parseConcat() *astNode {
	var nodes []*astNode
	for {
		t := p.peek()
		if t.typ == tokEOF || t.typ == tokRParen || t.typ == tokPipe {
			break
		}
		nodes = append(nodes, p.parseAtom())
	}
	if len(nodes) == 0 {
		return nil
	}
	if len(nodes) == 1 {
		return nodes[0]
	}
	result := nodes[0]
	for i := 1; i < len(nodes); i++ {
		result = &astNode{typ: "concat", left: result, right: nodes[i]}
	}
	return result
}

// parseAtom 解析原子表达式（字符或括号），并处理后缀 *, +, ?
func (p *parser) parseAtom() *astNode {
	t := p.peek()
	var node *astNode
	switch t.typ {
	case tokChar:
		p.consume()
		node = &astNode{typ: "char", val: t.value}
	case tokLParen:
		p.consume()
		node = p.parseUnion()
		if p.peek().typ != tokRParen {
			panic("expected )")
		}
		p.consume()
	default:
		panic("unexpected token")
	}

	// 处理后缀操作符（从左到右，但通常只有一个）
	for {
		switch p.peek().typ {
		case tokStar:
			p.consume()
			node = &astNode{typ: "star", child: node}
		case tokPlus:
			p.consume()
			// r+ ≡ r r*
			starNode := &astNode{typ: "star", child: node}
			node = &astNode{typ: "concat", left: node, right: starNode}
		case tokQuestion:
			p.consume()
			// r? ≡ (r | ε)
			epsilonNode := &astNode{typ: "char", val: model.Epsilon}
			node = &astNode{typ: "union", left: node, right: epsilonNode}
		default:
			return node
		}
	}
}

func parseRegex(pattern string) *astNode {
	if strings.TrimSpace(pattern) == "" {
		return nil
	}
	tokens := lex(pattern)
	p := &parser{tokens: tokens, pos: 0}
	return p.parseUnion()
}

// ===== 3. Thompson 构造器（不变）=====
type thompsonBuilder struct {
	states   []model.State
	trans    []model.Transition
	nextID   int
	alphabet map[model.Symbol]bool
}

func newThompsonBuilder() *thompsonBuilder {
	return &thompsonBuilder{
		states:   make([]model.State, 0),
		trans:    make([]model.Transition, 0),
		nextID:   0,
		alphabet: make(map[model.Symbol]bool),
	}
}

func (tb *thompsonBuilder) newState() model.State {
	id := fmt.Sprintf("q%d", tb.nextID)
	tb.nextID++
	tb.states = append(tb.states, model.State(id))
	return model.State(id)
}

func (tb *thompsonBuilder) addTransition(from, to model.State, input model.Symbol) {
	if input != model.Epsilon {
		tb.alphabet[input] = true
	}
	tb.trans = append(tb.trans, model.Transition{
		FromState: from,
		Input:     input,
		ToStates:  []model.State{to},
	})
}

func (tb *thompsonBuilder) build(node *astNode) (model.State, model.State) {
	if node == nil {
		s := tb.newState()
		a := tb.newState()
		tb.addTransition(s, a, model.Epsilon)
		return s, a
	}

	switch node.typ {
	case "char":
		if node.val == model.Epsilon {
			s := tb.newState()
			a := tb.newState()
			tb.addTransition(s, a, model.Epsilon)
			return s, a
		} else {
			s := tb.newState()
			a := tb.newState()
			tb.addTransition(s, a, node.val)
			return s, a
		}

	case "union":
		s := tb.newState()
		a := tb.newState()
		lStart, lAccept := tb.build(node.left)
		rStart, rAccept := tb.build(node.right)
		tb.addTransition(s, lStart, model.Epsilon)
		tb.addTransition(s, rStart, model.Epsilon)
		tb.addTransition(lAccept, a, model.Epsilon)
		tb.addTransition(rAccept, a, model.Epsilon)
		return s, a

	case "concat":
		lStart, lAccept := tb.build(node.left)
		rStart, rAccept := tb.build(node.right)
		tb.addTransition(lAccept, rStart, model.Epsilon)
		return lStart, rAccept

	case "star":
		s := tb.newState()
		a := tb.newState()
		cStart, cAccept := tb.build(node.child)
		tb.addTransition(s, cStart, model.Epsilon)
		tb.addTransition(s, a, model.Epsilon)
		tb.addTransition(cAccept, cStart, model.Epsilon)
		tb.addTransition(cAccept, a, model.Epsilon)
		return s, a

	default:
		panic("unknown node type")
	}
}

func RegexToNFA(pattern string) (*model.Automaton, error) {
	ast := parseRegex(pattern)
	if ast == nil {
		return nil, fmt.Errorf("empty or invalid pattern")
	}

	builder := newThompsonBuilder()
	start, accept := builder.build(ast)

	// 构建 alphabet（含 ε）
	var alphabet []model.Symbol
	for sym := range builder.alphabet {
		alphabet = append(alphabet, sym)
	}
	alphabet = append(alphabet, model.Epsilon)
	sort.Slice(alphabet, func(i, j int) bool {
		return string(alphabet[i]) < string(alphabet[j])
	})

	return &model.Automaton{
		States:          builder.states,
		Alphabet:        alphabet,
		Transitions:     builder.trans,
		InitialState:    start,
		AcceptingStates: []model.State{accept},
		IsDFA:           false,
	}, nil
}
