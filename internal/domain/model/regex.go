package model

import (
	"fmt"
	"unicode/utf8"
)

type Regex string

const (
	EmptyLanguageToken Regex = "∅" // 空集
	EmptyStringToken   Regex = Regex(Epsilon)
)

// 正则表达式支持的字符，包括 0~1，a~z，A~Z，|，（，），*，？，·，ε，∅
var ValidCSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789*+?|()ε∅")

// IsValid 验证正则表达式的结构和字符合法性
func (r Regex) IsValid() error {
	if r == "" {
		return fmt.Errorf("pattern cannot be empty")
	}
	if r == EmptyLanguageToken {
		return nil // ∅ 是合法的
	}

	s := string(r)

	// 1. 定义合法字符集
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789*+?|()"

	// 2. 状态机变量
	parenCount := 0         // 括号计数
	lastWasOperand := false // 上一个是否是操作数（字母、数字、ε、∅、)）

	i := 0
	for i < len(s) {
		// 读取一个 rune (支持 ε 和 ∅ 等多字节字符)
		char, size := utf8.DecodeRuneInString(s[i:])
		if char == utf8.RuneError {
			return fmt.Errorf("invalid character encoding")
		}

		strChar := string(char)

		// --- 检查操作数 ---
		if isOperand(char) {
			lastWasOperand = true
			i += size
			continue
		}

		// --- 检查操作符 ---
		switch char {
		case '(':
			parenCount++
			lastWasOperand = false // ( 后面必须跟操作数
		case ')':
			parenCount--
			if parenCount < 0 {
				return fmt.Errorf("unmatched closing parenthesis ')'")
			}
			// ) 后面可以是操作符，也可以是操作数（隐式连接）
			lastWasOperand = true
		case '|':
			if !lastWasOperand {
				return fmt.Errorf("invalid position for alternation '|': cannot follow nothing or another operator")
			}
			lastWasOperand = false // | 后面必须跟操作数
		case '*', '+', '?':
			if !lastWasOperand {
				return fmt.Errorf("invalid position for quantifier '%c': cannot follow nothing or operator", char)
			}
			// 量词后面可以是操作数（隐式连接）或操作符
			// 量词整体算作一个操作数单元
		default:
			// 检查是否是合法的基础字符（a-z, A-Z, 0-9）
			if !containsRune(validChars, char) {
				return fmt.Errorf("pattern contains invalid character: '%s'", strChar)
			}
		}

		i += size
	}

	// 3. 最终检查
	if parenCount != 0 {
		return fmt.Errorf("unmatched opening parenthesis '('")
	}

	// 如果最后一个字符是 | 或 (，也是非法的（除非是 ε 或 ∅）
	if !lastWasOperand && len(s) > 0 {
		return fmt.Errorf("expression cannot end with an operator")
	}

	return nil
}

// isOperand 判断是否是原子操作数
func isOperand(r rune) bool {
	// 字母数字
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
		return true
	}
	// 特殊符号
	if r == 'ε' || r == '∅' {
		return true
	}
	return false
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
