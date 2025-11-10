package model

import (
	"fmt"
	"regexp"
)

var ValidCSet []byte // 正则表达式支持的符合，包括0~1，a~z，A~Z，|，（，），*，？，·，

type Regex string

func (r Regex) IsValid() error {
	// 检查是否包含非法字符，限定所使用的符号为后端实际功能实现所支持的：仅包含：[a-zA-Z0-9]、*、+、？、|、（）
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9\*\+\?\|\(\)\s]+$`, string(r)); !matched {
		return fmt.Errorf("Pattern contains invalid characters. Only letters, digits, *, +, ?, |, (, ) are allowed.")
	}

	// 最终验证：尝试编译正则表达式
	if _, err := regexp.Compile(string(r)); err != nil {
		return fmt.Errorf("Invalid regular expression: " + err.Error())
	}
	return nil
}
