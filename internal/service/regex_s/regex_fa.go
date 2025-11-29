// regex_to_nfa.go
package regex_s

import (
	"backend/internal/domain/model"
	"fmt"
	"sort"
	"strings"

	"github.com/chenzanhong/zlog"
	"go.uber.org/zap"
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

func (node *astNode) String() string {
	if node == nil {
		return string(model.Epsilon)
	}

	switch node.typ {
	case "char":
		if node.val == model.Epsilon {
			return string(model.Epsilon)
		}
		return string(node.val)
	case "union":
		left := node.left.String()
		right := node.right.String()
		return fmt.Sprintf("(%s|%s)", left, right)
	case "concat":
		left := node.left.String()
		right := node.right.String()
		// // 对 union 子表达式加括号
		// if node.left.typ == "union" {
		//     left = "(" + left + ")"
		// }
		// if node.right.typ == "union" {
		//     right = "(" + right + ")"
		// }
		return left + right
	case "star":
		childStr := node.child.String()
		if node.child.typ == "char" || node.child.typ == "union" {
			return "(" + childStr + ")*"
		}
		return childStr + "*"
	default:
		return "?"
	}

}

// ===== 1. 词法分析器（支持 + ?）=====
func lex(pattern string) ([]token, error) {
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
				// panic(fmt.Sprintf("invalid character in regex: %c", r))
				// 不panic，而是返回错误
				return nil, fmt.Errorf("invalid character in regex: %c", r)
			}
		}
		i++
	}
	tokens = append(tokens, token{typ: tokEOF})
	return tokens, nil
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

func parseRegex(pattern string) (*astNode, error) {
	if strings.TrimSpace(pattern) == "" {
		return nil, nil
	}
	tokens, err := lex(pattern)
	if err != nil {
		return nil, err
	}
	p := &parser{tokens: tokens, pos: 0}
	return p.parseUnion(), nil
}

// ===== 3. Thompson 构造器 =====
// 记录每一步的中间自动机
type BuildStep = model.RegexToFAStep
type thompsonBuilder struct {
	states   []model.State
	trans    []model.Transition
	nextID   int
	alphabet map[model.Symbol]bool

	steps              []BuildStep
	uniqueInitialState model.State
	uniqueAcceptState  model.State
}

func newThompsonBuilder() *thompsonBuilder {
	return &thompsonBuilder{
		states:             append(make([]model.State, 0), model.AcceptState, model.InitialState),
		trans:              make([]model.Transition, 0),
		nextID:             0,
		alphabet:           make(map[model.Symbol]bool),
		steps:              make([]BuildStep, 0),
		uniqueInitialState: model.InitialState,
		uniqueAcceptState:  model.AcceptState,
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
		tb.recordStep(node, s, a)
		return s, a
	}

	var start, end model.State

	switch node.typ {
	case "char":
		if node.val == model.Epsilon {
			s := tb.newState()
			a := tb.newState()
			tb.addTransition(s, a, model.Epsilon)
			start, end = s, a
		} else {
			s := tb.newState()
			a := tb.newState()
			tb.addTransition(s, a, node.val)
			start, end = s, a
		}

	case "union":
		lStart, lEnd := tb.build(node.left)
		rStart, rEnd := tb.build(node.right)
		s := tb.newState()
		tb.addTransition(s, lStart, model.Epsilon)
		tb.addTransition(s, rStart, model.Epsilon)
		a := tb.newState()
		tb.addTransition(lEnd, a, model.Epsilon)
		tb.addTransition(rEnd, a, model.Epsilon)
		start, end = s, a

	case "concat":
		lStart, lEnd := tb.build(node.left)
		rStart, rEnd := tb.build(node.right)
		tb.addTransition(lEnd, rStart, model.Epsilon)
		start, end = lStart, rEnd

	case "star":
		s := tb.newState()
		a := tb.newState()
		cStart, cEnd := tb.build(node.child)
		tb.addTransition(s, cStart, model.Epsilon)
		tb.addTransition(s, a, model.Epsilon)
		tb.addTransition(cEnd, cStart, model.Epsilon)
		tb.addTransition(cEnd, a, model.Epsilon)
		start, end = s, a

	default:
		zlog.Panic("regex to fa: unknown node type")
	}

	tb.recordStep(node, start, end)
	return start, end
}

func (tb *thompsonBuilder) recordStep(node *astNode, start, end model.State) {
	// 构建 alphabet（不含 ε）
	alphabet := make([]model.Symbol, 0)
	for sym := range tb.alphabet {
		alphabet = append(alphabet, sym)
	}
	sort.Slice(alphabet, func(i, j int) bool {
		return string(alphabet[i]) < string(alphabet[j])
	})

	// 深拷贝 slices，后续会修改 tb.states / tb.trans
	statesCopy := make([]model.State, len(tb.states))
	copy(statesCopy, tb.states)

	transCopy := make([]model.Transition, len(tb.trans))
	copy(transCopy, tb.trans)

	alphabetCopy := make([]model.Symbol, len(alphabet))
	copy(alphabetCopy, alphabet)

	fa := &model.Automaton{
		States:      statesCopy,
		Alphabet:    alphabetCopy,
		Transitions: transCopy,
		// InitialState:    start, //tb.uniqueInitialState,
		// AcceptingStates: []model.State{end},
		Type: model.EpsilonNFA,
	}

	expr := node.String()
	tb.steps = append(tb.steps, BuildStep{
		Expr:          expr,
		StartState:    start,
		EndState:      end,
		AutomatonFlow: fa.ToReactFlow(),
	})
}

// 正则表达式转FA，不带转换过程
func RegexToFA(regex model.Regex) (*model.Automaton, error) {
	pattern := string(regex)
	ast, err := parseRegex(pattern)
	if err != nil {
		return nil, err
	} else if ast == nil {
		return nil, fmt.Errorf("empty pattern")
	}

	builder := newThompsonBuilder()
	start, end := builder.build(ast)

	// 构建 alphabet（不含 ε）
	var alphabet []model.Symbol
	for sym := range builder.alphabet {
		alphabet = append(alphabet, sym)
	}
	sort.Slice(alphabet, func(i, j int) bool {
		return string(alphabet[i]) < string(alphabet[j])
	})

	return &model.Automaton{
		States:          builder.states,
		Alphabet:        alphabet,
		Transitions:     builder.trans,
		InitialState:    start,
		AcceptingStates: []model.State{end},
		Type:            model.EpsilonNFA,
	}, nil
}

// RegexToFAWithSteps 正则表达式转FA，返回完整的转换步骤序列
func RegexToFAWithSteps(regex model.Regex) (result *model.RegexToFAProcess, err error) {
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Errorf("regex to fa with steps: recovered from panic, %v", r)
			zlog.Warnf("regex to fa with steps: recovered from panic, %v", zap.Any("panic", r))
			result = nil
			err = errMsg
		}
	}()
	pattern := string(regex)
	if strings.TrimSpace(pattern) == "" {
		return nil, fmt.Errorf("empty pattern")
	}

	ast, err := parseRegex(pattern)
	if err != nil {
		return nil, err
	}
	if ast == nil {
		return nil, fmt.Errorf("parsed to nil AST")
	}

	builder := &thompsonBuilder{
		states:   make([]model.State, 0),
		trans:    make([]model.Transition, 0),
		nextID:   0,
		alphabet: make(map[model.Symbol]bool),
		steps:    make([]BuildStep, 0),
	}

	start, end := builder.build(ast) // 忽略返回值，我们只关心 steps
	finalAutomaton := builder.steps[len(builder.steps)-1].AutomatonFlow.ToAutomaton()
	finalAutomaton.InitialState = start
	finalAutomaton.AcceptingStates = append(finalAutomaton.AcceptingStates, end)
	return &model.RegexToFAProcess{
		Regex:          regex,
		Steps:          builder.steps,
		FinalAutomaton: finalAutomaton,
	}, nil
}
