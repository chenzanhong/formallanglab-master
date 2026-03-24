package regex_s

import (
	"regexp"

	"backend/internal/domain/model"
)

func RegexRecognize(pattern model.Regex, str string) (bool, error) {
	// 空集不接受任何字符串
	if pattern == model.EmptyLanguageToken {
		return false, nil // 空集永远不匹配任何字符串，包括空串
	}

	// 处理 ε 符号
	regexStr := string(pattern)
	if regexStr == string(model.Epsilon) {
		// ε 只匹配空字符串
		return str == "", nil
	}

	// 编译正则表达式为完整匹配模式（^...$）
	re, err := regexp.Compile("^(?:" + regexStr + ")$")
	if err != nil {
		return false, err
	}
	match := re.MatchString(str)
	if !match {
		return match, nil
	}

	return match, nil
}
