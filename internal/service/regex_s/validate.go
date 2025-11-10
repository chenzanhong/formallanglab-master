package regex_s

import (
	"backend/internal/domain/model"
	"fmt"
	"regexp"
)

func RegexValidate(regex model.Regex) error {
	// 调用 Regex 对象的 IsValid 方法进行验证
	return regex.IsValid()
}

func RegexValidString(str string) (bool, error) {
	// 查看字符串是否只包含 a-zA-Z0-9 范围内的字符
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, str); !matched {
		return false, fmt.Errorf("字符串只能包含字母和数字")
	}
	return true, nil
}
