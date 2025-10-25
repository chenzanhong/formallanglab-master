package regexs

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
