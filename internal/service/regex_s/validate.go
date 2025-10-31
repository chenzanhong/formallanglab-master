package regex_s

import (
	"fmt"
	"regexp"
)

func RegexValidate(pattern string) (bool, error) {
	// 检查是否包含非法字符，限定所使用的符号为后端实际功能实现所支持的：仅包含：[a-zA-Z0-9]、*、+、？、|、（）
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9\*\+\?\|\(\)\s]+$`, pattern); !matched {
		return false, fmt.Errorf("Pattern contains invalid characters. Only letters, digits, *, +, ?, |, (, ) are allowed.")
	}

	// 最终验证：尝试编译正则表达式
	if _, err := regexp.Compile(pattern); err != nil {
		return false, fmt.Errorf("Invalid regular expression: " + err.Error())
	}
	return true, nil
}

func RegexValidString(str string) (bool, error) {
	// 查看字符串是否只包含 a-zA-Z0-9 范围内的字符
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, str); !matched {
		return false, fmt.Errorf("字符串只能包含字母和数字")
	}
	return true, nil
}
