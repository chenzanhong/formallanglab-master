package model

import (
	"fmt"
	"regexp"
)

// 正则表达式支持的符合，包括0~1，a~z，A~Z，|，（，），*，？，·，
var ValidCSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789*+?|()")

// 允许：字母、数字、*, +, ?, |, (, )
const ValidRegex = `^[a-zA-Z0-9*+?|()]+$`

type Regex string

const (
	EmptyLanguageToken Regex = "∅" // 空集
)

func (r Regex) IsValid() error {
	if r == EmptyLanguageToken {
		return nil
	}
	// 检查是否包含非法字符，限定所使用的符号为后端实际功能实现所支持的：仅包含：[a-zA-Z0-9]、*、+、？、|、（）
	matched, err := regexp.MatchString(ValidRegex, string(r))
	if err != nil {
		// 这是程序 bug，不是用户输入问题
		panic(fmt.Sprintf("internal regex compile error: %v", err))
	}
	if !matched {
		return fmt.Errorf("pattern contains invalid characters: only letters, digits, *, +, ?, |, (, ) are allowed")
	}

	// 最终验证：尝试编译正则表达式
	if _, err := regexp.Compile(string(r)); err != nil {
		return fmt.Errorf("Invalid regular expression: " + err.Error())
	}
	return nil
}
