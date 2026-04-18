package regex_s

import (
	"fmt"
	"strings"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// simplifyASTNode AST 节点类型
type simplifyASTNode struct {
	typ      string // "char", "union", "concat", "star", "plus", "question"
	val      rune
	children []*simplifyASTNode
}

// String 将 AST 转为字符串
func (node *simplifyASTNode) String() string {
	if node == nil {
		return ""
	}

	switch node.typ {
	case "char":
		return string(node.val)
	case "union":
		if len(node.children) == 2 {
			return fmt.Sprintf("(%s|%s)", node.children[0].String(), node.children[1].String())
		}
		// 多个分支的并集
		parts := make([]string, len(node.children))
		for i, child := range node.children {
			parts[i] = child.String()
		}
		return "(" + strings.Join(parts, "|") + ")"
	case "concat":
		parts := make([]string, len(node.children))
		for i, child := range node.children {
			parts[i] = child.String()
		}
		return strings.Join(parts, "")
	case "star":
		childStr := node.children[0].String()
		if node.children[0].typ == "union" || node.children[0].typ == "concat" {
			return "(" + childStr + ")*"
		}
		return childStr + "*"
	case "plus":
		childStr := node.children[0].String()
		if node.children[0].typ == "union" || node.children[0].typ == "concat" {
			return "(" + childStr + ")+"
		}
		return childStr + "+"
	case "question":
		childStr := node.children[0].String()
		if node.children[0].typ == "union" || node.children[0].typ == "concat" {
			return "(" + childStr + ")?"
		}
		return childStr + "?"
	}

	return "?"
}

// SimplifyRegex 使用递归下降解析器化简正则表达式
func SimplifyRegex(regexStr string) (simplified string, isEmptyLanguage bool) {
	// 词法分析
	tokens, err := simplifyLex(regexStr)
	if err != nil {
		return regexStr, false
	}

	// 语法分析构建 AST
	parser := &simplifyParser{tokens: tokens, pos: 0}
	ast := parser.parseUnion()

	// 化简 AST
	simplifiedAST := simplifyAST(ast)

	// 检查是否为空语言
	if isZeroLanguage(simplifiedAST) {
		return "", true
	}

	// 将化简后的 AST 转回字符串
	return simplifiedAST.String(), false
}

// simplifyLex 词法分析
func simplifyLex(pattern string) ([]token, error) {
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
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				tokens = append(tokens, token{typ: tokChar, value: model.Symbol(string(r))})
			} else if r == 'ε' {
				tokens = append(tokens, token{typ: tokChar, value: model.Epsilon})
			} else if r == '∅' {
				tokens = append(tokens, token{typ: tokChar, value: model.Symbol("∅")})
			} else {
				return nil, fmt.Errorf("invalid character: %c", r)
			}
		}
		i++
	}
	tokens = append(tokens, token{typ: tokEOF})
	return tokens, nil
}

// simplifyParser 递归下降解析器
type simplifyParser struct {
	tokens []token
	pos    int
}

func (p *simplifyParser) peek() token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token{typ: tokEOF}
}

func (p *simplifyParser) consume() token {
	t := p.peek()
	p.pos++
	return t
}

func (p *simplifyParser) parseUnion() *simplifyASTNode {
	left := p.parseConcat()
	for p.peek().typ == tokPipe {
		p.consume()
		right := p.parseConcat()
		left = &simplifyASTNode{typ: "union", children: []*simplifyASTNode{left, right}}
	}
	return left
}

func (p *simplifyParser) parseConcat() *simplifyASTNode {
	var nodes []*simplifyASTNode
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
		result = &simplifyASTNode{typ: "concat", children: []*simplifyASTNode{result, nodes[i]}}
	}
	return result
}

func (p *simplifyParser) parseAtom() *simplifyASTNode {
	t := p.peek()
	var node *simplifyASTNode
	switch t.typ {
	case tokChar:
		p.consume()
		node = &simplifyASTNode{typ: "char", val: []rune(t.value)[0]}
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

	// 处理后缀操作符
	for {
		switch p.peek().typ {
		case tokStar:
			p.consume()
			node = &simplifyASTNode{typ: "star", children: []*simplifyASTNode{node}}
		case tokPlus:
			p.consume()
			node = &simplifyASTNode{typ: "plus", children: []*simplifyASTNode{node}}
		case tokQuestion:
			p.consume()
			node = &simplifyASTNode{typ: "question", children: []*simplifyASTNode{node}}
		default:
			return node
		}
	}
}

// simplifyAST 化简 AST
func simplifyAST(node *simplifyASTNode) *simplifyASTNode {
	if node == nil {
		return nil
	}

	// 递归化简子节点
	for i, child := range node.children {
		node.children[i] = simplifyAST(child)
	}

	switch node.typ {
	case "char":
		// 字符节点，不需要化简
		return node

	case "star":
		return simplifyStar(node)

	case "plus":
		return simplifyPlus(node)

	case "question":
		return simplifyQuestion(node)

	case "concat":
		return simplifyConcat(node)

	case "union":
		return simplifyUnion(node)
	}

	return node
}

// simplifyStar 化简 * 操作
func simplifyStar(node *simplifyASTNode) *simplifyASTNode {
	child := node.children[0]

	// ∅* = ε
	if child.typ == "char" && child.val == '∅' {
		return &simplifyASTNode{typ: "char", val: 'ε'}
	}

	// ε* = ε
	if child.typ == "char" && child.val == 'ε' {
		return &simplifyASTNode{typ: "char", val: 'ε'}
	}

	return node
}

// simplifyPlus 化简 + 操作
func simplifyPlus(node *simplifyASTNode) *simplifyASTNode {
	child := node.children[0]

	// ∅+ = ∅
	if child.typ == "char" && child.val == '∅' {
		return &simplifyASTNode{typ: "char", val: '∅'}
	}

	// ε+ = ε
	if child.typ == "char" && child.val == 'ε' {
		return &simplifyASTNode{typ: "char", val: 'ε'}
	}

	return node
}

// simplifyQuestion 化简 ? 操作
func simplifyQuestion(node *simplifyASTNode) *simplifyASTNode {
	child := node.children[0]

	// ∅? = ε
	if child.typ == "char" && child.val == '∅' {
		return &simplifyASTNode{typ: "char", val: 'ε'}
	}

	// ε? = ε
	if child.typ == "char" && child.val == 'ε' {
		return &simplifyASTNode{typ: "char", val: 'ε'}
	}

	return node
}

// simplifyConcat 化简连接操作
func simplifyConcat(node *simplifyASTNode) *simplifyASTNode {
	if len(node.children) < 2 {
		return node
	}

	left := node.children[0]
	right := node.children[1]

	// εA = A
	if left.typ == "char" && left.val == 'ε' {
		return right
	}

	// Aε = A
	if right.typ == "char" && right.val == 'ε' {
		return left
	}

	// ∅A = ∅
	if left.typ == "char" && left.val == '∅' {
		return &simplifyASTNode{typ: "char", val: '∅'}
	}

	// A∅ = ∅
	if right.typ == "char" && right.val == '∅' {
		return &simplifyASTNode{typ: "char", val: '∅'}
	}

	return node
}

// simplifyUnion 化简并集操作
func simplifyUnion(node *simplifyASTNode) *simplifyASTNode {
	if len(node.children) < 2 {
		return node
	}

	left := node.children[0]
	right := node.children[1]

	// A|∅ = A
	if left.typ == "char" && left.val == '∅' {
		return right
	}

	// ∅|A = A
	if right.typ == "char" && right.val == '∅' {
		return left
	}

	// A|A = A (可选优化)
	if left.typ == right.typ && left.val == right.val && left.typ == "char" {
		return left
	}

	return node
}

// isZeroLanguage 检查 AST 是否表示空语言
func isZeroLanguage(node *simplifyASTNode) bool {
	if node == nil {
		return false
	}

	switch node.typ {
	case "char":
		return node.val == '∅'
	case "concat":
		// 连接中只要有一个是 ∅，结果就是 ∅
		for _, child := range node.children {
			if isZeroLanguage(child) {
				return true
			}
		}
		return false
	case "union":
		// 并集所有分支都是 ∅，结果才是 ∅
		for _, child := range node.children {
			if !isZeroLanguage(child) {
				return false
			}
		}
		return true
	case "star", "plus", "question":
		// * 和 ? 永远不会是空语言（至少接受 ε）
		// + 如果子节点是 ∅，则是 ∅
		if node.typ == "plus" {
			return isZeroLanguage(node.children[0])
		}
		return false
	}

	return false
}
