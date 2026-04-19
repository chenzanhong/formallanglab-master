package regex_s

import (
	"fmt"
	"regexp"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/service/automaton_s"
	"github.com/chenzanhong/formallanglab-master/pkg/util"
)

/*
实际实现策略（三层渐进式判定）：
1. 先化简 SimplifyRegex，比较是否完全相同
   - 相同：直接返回 true（等价）
2. 使用日常方法（有限测试集）
   - 生成测试字符串集，用 regexp.MatchString 分别匹配
   - 若存在反例（匹配结果不同）：返回 false（不等价）
   - 若全部相同：进入下一步
3. 使用理论方法（基于自动机，确保判定结果）
   - 将正则转换为 NFA（Thompson 构造法）
   - NFA 确定化为 DFA（子集构造法）
   - DFA 最小化（Hopcroft 算法）
   - 比较两个最小 DFA 是否同构
   - 返回最终判定结果（一定判断出结果）
*/

// RegexEquivalenceCheck 检查两个正则表达式是否等价（三层渐进式判定）
func RegexEquivalenceCheck(pattern1 model.Regex, pattern2 model.Regex) (bool, error) {
	// Step 0: 处理空集特殊情况
	if (pattern1 == model.EmptyLanguageToken) != (pattern2 == model.EmptyLanguageToken) {
		// 一个为空集，一个不为空集 → 不等价
		return false, nil
	} else if pattern1 == model.EmptyLanguageToken && pattern2 == model.EmptyLanguageToken {
		// 都为空集 → 等价
		return true, nil
	}

	// Step 1: 先化简，比较是否完全相同
	simplified1, _ := SimplifyRegex(pattern1.String())
	simplified2, _ := SimplifyRegex(pattern2.String())
	if simplified1 == simplified2 {
		// 化简后相同 → 等价
		return true, nil
	}

	// Step 2: 使用日常方法（有限测试集）
	equivalent, err := checkByTesting(pattern1, pattern2)
	if err != nil {
		// 找到反例 → 不等价
		return false, nil
	}
	if !equivalent {
		// 测试发现不等价
		return false, nil
	}

	// Step 3: 使用理论方法（基于自动机，确保判定结果）
	return checkByAutomata(pattern1, pattern2)
}

// checkByTesting 通过有限测试集检查正则等价性（日常方法）
func checkByTesting(pattern1 model.Regex, pattern2 model.Regex) (bool, error) {
	// 提取字母表（使用专门用于等价性检查的版本）
	alphabet1 := extractAlphabetForCheck(string(pattern1))
	alphabet2 := extractAlphabetForCheck(string(pattern2))
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

	// 生成测试字符串集
	maxLen := max(min(len(pattern1)*2, len(pattern1)+3), min(len(pattern2)*2, len(pattern2)+3))
	candidates := util.GenerateStrings(alphabet, maxLen)

	for _, candidate := range candidates {
		m1 := re1.MatchString(candidate)
		m2 := re2.MatchString(candidate)

		if m1 != m2 {
			return false, fmt.Errorf("找到反例：%s。%s匹配结果为%v，但是%s匹配结果为%v", candidate, pattern1, m1, pattern2, m2)
		}
	}

	return true, nil
}

// checkByAutomata 通过自动机方法检查正则等价性（理论方法）
func checkByAutomata(pattern1 model.Regex, pattern2 model.Regex) (bool, error) {
	fa1, err := RegexToFA(pattern1)
	if err != nil {
		return false, err
	}
	fa2, err := RegexToFA(pattern2)
	if err != nil {
		return false, err
	}

	_, _, isEquivalent := automaton_s.AutomatonEquivalenceCheck(fa1, fa2)

	return isEquivalent, nil
}

// extractAlphabetForCheck 从正则表达式中提取字母表（用于等价性检查）
// 与 generate.go 中的 extractAlphabet 不同，这个版本处理更一般的字符集
func extractAlphabetForCheck(regex string) []string {
	seen := make(map[string]bool)
	i := 0
	for i < len(regex) {
		b := regex[i]
		// 跳过元字符
		if b == '(' || b == ')' || b == '|' || b == '*' || b == '+' || b == '?' || b == '.' || b == '[' || b == ']' || b == '^' || b == '$' {
			i++
			continue
		}
		// 处理转义字符
		if b == '\\' && i+1 < len(regex) {
			i += 2
			continue
		}
		// 添加字符到字母表
		seen[string(b)] = true
		i++
	}

	alphabet := make([]string, 0, len(seen))
	for b := range seen {
		alphabet = append(alphabet, b)
	}
	return alphabet
}
