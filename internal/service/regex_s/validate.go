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

var alphaNumRegexp = regexp.MustCompile(`^[a-zA-Z0-9]*$`) // 或 +，根据需求

func RegexValidString(str string) (bool, error) {
	if !alphaNumRegexp.MatchString(str) {
		return false, fmt.Errorf("字符串只能包含字母和数字")
	}
	return true, nil
}
