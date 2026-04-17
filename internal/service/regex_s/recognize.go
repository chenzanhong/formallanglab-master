package regex_s

import (
	"regexp"
	"strings"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// SimplifyRegex 化简正则表达式，处理 ∅ 和 ε
// 返回化简后的正则表达式和是否为空语言
func SimplifyRegex(regexStr string) (simplified string, isEmptyLanguage bool) {
	// 空集不接受任何字符串
	if regexStr == string(model.EmptyLanguageToken) {
		return "", true
	}

	// 处理 ε 符号：如果整个模式就是 ε，只匹配空字符串
	if regexStr == string(model.Epsilon) {
		return "", false
	}

	// 预处理：根据形式语言理论进行代数化简
	processedRegex := regexStr

	// ========== 第 1 步：处理 ∅ 的闭包 ==========
	// ∅* = ε, ∅+ = ∅, ∅? = ε
	// ∅* 和 ∅? 替换为匹配空串的特殊标记
	for {
		newRegex := regexp.MustCompile(`∅\+`).ReplaceAllString(processedRegex, "∅")
		newRegex = regexp.MustCompile(`∅\*`).ReplaceAllString(newRegex, "\x01") // 特殊标记
		newRegex = regexp.MustCompile(`∅\?`).ReplaceAllString(newRegex, "\x01") // 特殊标记
		if newRegex == processedRegex {
			break
		}
		processedRegex = newRegex
	}

	// ========== 第 2 步：处理 ε 的闭包 ==========
	// ε* = ε, ε+ = ε, ε? = ε → 替换为匹配空串的特殊标记
	for {
		newRegex := regexp.MustCompile(`ε[*+?]+`).ReplaceAllString(processedRegex, "\x01")
		if newRegex == processedRegex {
			break
		}
		processedRegex = newRegex
	}

	// ========== 第 3 步：处理并集中的 ∅ ==========
	// a|∅ = a, ∅|a = a
	// 特殊处理：(∅|ε) = ε, (ε|∅) = ε
	for {
		changed := false

		// 先处理 (∅|ε) 和 (ε|∅) 这种情况
		newRegex := regexp.MustCompile(`\(([^()]*?)∅\|ε([^()]*?)\)`).ReplaceAllString(processedRegex, `$1$2`)
		if newRegex != processedRegex {
			changed = true
			processedRegex = newRegex
			continue
		}

		newRegex = regexp.MustCompile(`\(([^()]*?)ε\|∅([^()]*?)\)`).ReplaceAllString(processedRegex, `$1$2`)
		if newRegex != processedRegex {
			changed = true
			processedRegex = newRegex
			continue
		}

		// 处理括号内的 |∅ 和 ∅|
		newRegex = regexp.MustCompile(`(\([^()]*?)\|∅`).ReplaceAllString(processedRegex, `$1`)
		if newRegex != processedRegex {
			changed = true
			processedRegex = newRegex
			continue
		}

		newRegex = regexp.MustCompile(`∅\|(\([^()]*?)`).ReplaceAllString(processedRegex, `$1`)
		if newRegex != processedRegex {
			changed = true
			processedRegex = newRegex
			continue
		}

		// 处理顶层的 |∅ 和 ∅|
		newRegex = regexp.MustCompile(`^([^|]+?)\|∅$`).ReplaceAllString(processedRegex, `$1`)
		if newRegex != processedRegex {
			changed = true
			processedRegex = newRegex
			continue
		}

		newRegex = regexp.MustCompile(`^∅\|(.+?)$`).ReplaceAllString(processedRegex, `$1`)
		if newRegex != processedRegex {
			changed = true
			processedRegex = newRegex
			continue
		}

		if !changed {
			break
		}
	}

	// ========== 第 3.8 步：处理单独的 ∅ ==========
	// 如果化简后只剩下 ∅，说明是空语言
	if processedRegex == "∅" {
		return "", true
	}

	// ========== 第 4 步：处理连接中的 ∅ ==========
	// a∅ = ∅, ∅a = ∅
	if regexp.MustCompile(`[^|(*+?]\s*∅|∅\s*[^|)*+?]`).MatchString(processedRegex) {
		return "", true
	}

	// ========== 第 5 步：处理 ε 的连接 ==========
	// aε = a, εa = a（直接移除）
	processedRegex = regexp.MustCompile(`ε`).ReplaceAllString(processedRegex, "")

	// ========== 第 6 步：处理剩余的 ∅ ==========
	// 替换为不匹配模式
	processedRegex = regexp.MustCompile(`∅`).ReplaceAllString(processedRegex, "(?!x)x")

	// ========== 第 7 步：处理特殊标记 ==========
	// \x01 代表 ε（空串），在并集中使用
	if processedRegex == "" || processedRegex == "\x01" {
		return "", false
	}

	// 将 \x01 替换为匹配空串的模式
	// (a|\x01) -> (a|) 在 Go regexp 中不合法，需要特殊处理
	// 策略：将 (X|\x01) 转换为 (X)?
	processedRegex = regexp.MustCompile(`\(([^()|]+)\|\x01\)`).ReplaceAllString(processedRegex, `($1)?`)
	processedRegex = regexp.MustCompile(`\(\x01\|([^()|]+)\)`).ReplaceAllString(processedRegex, `($1)?`)
	// 处理剩余的 \x01
	processedRegex = regexp.MustCompile(`\x01`).ReplaceAllString(processedRegex, "")

	// 如果处理后为空
	if processedRegex == "" {
		return "", false
	}

	return processedRegex, false
}

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
