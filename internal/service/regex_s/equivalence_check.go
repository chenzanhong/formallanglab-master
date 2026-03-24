package regex_s

import (
	"fmt"
	"regexp"

	"backend/internal/domain/model"
	"backend/pkg/util"
)

/*
理论：
	将两个正则表达式 R₁ 和 R₂ 分别转换为等价的 NFA（使用 Thompson 构造法）
	将 NFA 确定化（子集构造法）
	对 DFA 进行最小化（Hopcroft 算法等）
	比较两个最小 DFA 是否同构（结构相同
日常：
	通过“有限测试集 + 补集空性检查”近似判定（实用但不完备）
	原理：
		如果两个正则不等价，则存在一个最短反例字符串 w
	实现步骤：
		提取字母表 Σ（从 r1 和 r2 中提取 [a-zA-Z0-9] 字符）
		生成所有长度 0 到 maxLen 的字符串（maxLen = 10 足够大多数情况）
		对每个字符串 s，用 regexp.MatchString 分别匹配 r1 和 r2
		若存在 s 使得匹配结果不同 → 不等价
		若全部相同 → 认为等价（有极小概率漏判，但实践中可靠）
*/
// RegexEquivalenceCheck 检查两个正则表达式是否等价
func RegexEquivalenceCheck(pattern1 model.Regex, pattern2 model.Regex) (bool, error) {
	if (pattern1 == model.EmptyLanguageToken) != (pattern2 == model.EmptyLanguageToken) { // 一个为空集
		return false, nil
	} else if pattern1 == model.EmptyLanguageToken && pattern2 == model.EmptyLanguageToken { // 都为空集
		return true, nil
	}
	// 提取字母表
	alphabet1 := extractAlphabet(string(pattern1))
	alphabet2 := extractAlphabet(string(pattern2))
	// 合并字母表
	seen := make(map[string]bool)
	for _, b := range alphabet1 {
		seen[b] = true
	}
	for _, b := range alphabet2 {
		seen[b] = true
	}
	alphabet := make([]string, 0, len(seen))
	for b := range seen {
		alphabet = append(alphabet, b)
	}

	if len(alphabet) == 0 {
		return true, nil
	}

	// 编译正则表达式为完整匹配模式（^...$）
	re1, err := regexp.Compile("^(?:" + string(pattern1) + ")$")
	if err != nil {
		return false, err
	}
	re2, err := regexp.Compile("^(?:" + string(pattern2) + ")$")
	if err != nil {
		return false, err
	}

	candidates := util.GenerateStrings(alphabet, max(min(len(pattern1)*2, len(pattern1)+3), min(len(pattern2)*2, len(pattern2)+3)))
	m1, m2 := false, false
	for _, candidate := range candidates {
		m1 = re1.MatchString(candidate)
		m2 = re2.MatchString(candidate)
		if m1 != m2 {
			return false, fmt.Errorf("找到反例: %s。%s匹配结果为%v，但是%s匹配结果为%v", candidate, pattern1, m1, pattern2, m2)
		}
	}

	return true, nil
}
