package regex_s

import (
	"backend/internal/domain/model"
	"backend/pkg/util"
	"regexp"
	"slices"
)

const (
	maxExampleNum = 6
)

// 理论：转DFA+补集
// 实用：枚举+验证
func RegexGenerateExampleString(regex model.Regex) (accept, reject []string) {
	pattern := string(regex)

	// 编译正则表达式为完整匹配模式（^...$）
	re, err := regexp.Compile("^(?:" + pattern + ")$")
	if err != nil {
		msg := "(编译失败: " + err.Error() + ")"
		return []string{msg}, []string{msg}
	}

	// 提取字母表
	alphabet := extractAlphabet(pattern)

	// 生成候选字符串
	candidates := util.GenerateStrings(alphabet, min(len(regex)*2, len(regex)+3))

	acceptMap := make(map[string]bool)
	rejectMap := make(map[string]bool)
	acceptByLen := make(map[int][]string)
	rejectByLen := make(map[int][]string)

	// 筛选出可匹配的字符串
	for _, candidate := range candidates {
		if re.MatchString(candidate) {
			if !acceptMap[candidate] {
				acceptByLen[len(candidate)] = append(acceptByLen[len(candidate)], candidate)
				acceptMap[candidate] = true
			}
		} else {
			if !rejectMap[candidate] {
				rejectByLen[len(candidate)] = append(rejectByLen[len(candidate)], candidate)
				rejectMap[candidate] = true
			}
		}
	}

	// 确保数量不超过maxExampleNum
	accept = util.SamplingExampleStrings(acceptByLen, maxExampleNum)
	reject = util.SamplingExampleStrings(rejectByLen, maxExampleNum)

	// 兜底
	if len(accept) == 0 {
		accept = []string{"(no short accepted string found)"}
	}
	if len(reject) == 0 {
		reject = []string{"(all short strings are accepted)"}
	}

	return accept, reject
}

// extractAlphabet 从正则表达式字符串中提取字面量字符（a-zA-Z0-9）
func extractAlphabet(pattern string) []string {
	seen := make(map[rune]bool)
	var alphabet []string
	for _, char := range pattern {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			if !seen[char] {
				seen[char] = true
				alphabet = append(alphabet, string(char))
			}
		}
	}
	slices.Sort(alphabet)
	return alphabet
}
