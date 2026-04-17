package regex_s

import (
	"regexp"
	"slices"
	"strings"

	"github.com/chenzanhong/zlog"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/service/automaton_s"
	"github.com/chenzanhong/formallanglab-master/pkg/util"
)

const (
	maxExampleNum = 6
)

// 生成正则表达式可匹配和不可匹配的字符串示例。
func RegexGenerateExampleString(regex model.Regex) (accept, reject []string) {
	if regex == model.EmptyLanguageToken {
		return []string{}, []string{}
	}
	// 转 DFA+补集
	return regexGenerateExampleStringByCompletedDFAAndBFS(regex)
}

// 转 DFA+ 补集
func regexGenerateExampleStringByCompletedDFAAndBFS(regex model.Regex) (accept, reject []string) {
	// 先化简正则表达式
	pattern := string(regex)
	if strings.Contains(pattern, "∅") || strings.Contains(pattern, "ε") {
		simplified, isEmptyLanguage := SimplifyRegex(pattern)
		if isEmptyLanguage {
			return []string{}, []string{}
		}
		// 使用化简后的正则表达式
		regex = model.Regex(simplified)
	}

	a, err := RegexToFA(regex)
	if err != nil {
		zlog.Info("转为 DFA 失败")
		return regexGenerateExampleStringByEnumAndVerify(regex)
	}
	// 验证自动机并设置 Type 字段
	if err := a.Validate(); err != nil {
		zlog.Info("验证自动机失败：" + err.Error())
		return regexGenerateExampleStringByEnumAndVerify(regex)
	}
	if a.Type != model.DFA {
		a = automaton_s.NFAToDFA(a)
	}
	if err := a.CompleteDFA(); err != nil {
		zlog.Info("完备化失败：" + err.Error())
		return regexGenerateExampleStringByEnumAndVerify(regex)
	}

	return automaton_s.GenerateExampleStringsFromCompletedDFA(a)
}

// 枚举+验证
func regexGenerateExampleStringByEnumAndVerify(regex model.Regex) (accept, reject []string) {
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
	candidates := util.GenerateStrings(alphabet, min(len(pattern)*2, len(pattern)+3))

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

	// 随机采样
	if len(accept) > maxExampleNum {
		accept = util.SamplingExampleStrings(accept, maxExampleNum)
	}
	if len(reject) > maxExampleNum {
		reject = util.SamplingExampleStrings(reject, maxExampleNum)
	}

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
