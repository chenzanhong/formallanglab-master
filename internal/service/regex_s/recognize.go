package regex_s

import (
	"regexp"
	"strings"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func RegexRecognize(pattern model.Regex, str string) (bool, error) {
	regexStr := string(pattern)

	// 空集不接受任何字符串
	if regexStr == string(model.EmptyLanguageToken) {
		return false, nil
	}

	// 处理 ε 符号：如果整个模式就是 ε，只匹配空字符串
	if regexStr == string(model.Epsilon) {
		return str == "", nil
	}

	// 优化：如果不包含 ∅ 或 ε，直接使用 Go regexp
	if !strings.Contains(regexStr, "∅") && !strings.Contains(regexStr, "ε") {
		re, err := regexp.Compile("^(?:" + regexStr + ")$")
		if err != nil {
			return false, err
		}

		return re.MatchString(str), nil
	}

	// 化简正则表达式
	processedRegex, isEmptyLanguage := SimplifyRegex(regexStr)
	if isEmptyLanguage {
		return false, nil
	}

	// 编译正则表达式为完整匹配模式（^...$）
	re, err := regexp.Compile("^(?:" + processedRegex + ")$")
	if err != nil {
		return false, err
	}

	return re.MatchString(str), nil
}
