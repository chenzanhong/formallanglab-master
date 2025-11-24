package regex_s

import (
	"backend/internal/domain/model"
	"regexp"
)

func RegexRecognize(pattern model.Regex, str string) (bool, error) {
	// 编译正则表达式为完整匹配模式（^...$）
	re, err := regexp.Compile("^(?:" + string(pattern) + ")$")
	if err != nil {
		return false, err
	}
	match := re.MatchString(str)
	if !match {
		return match, nil
	}
	return match, nil
}
