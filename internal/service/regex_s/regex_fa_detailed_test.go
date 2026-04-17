package regex_s

import (
	"strings"
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegexToFANoEmptySet(t *testing.T) {
	// 测试 a(∅|ε)a 转换后的 NFA 不包含 ∅
	regex := model.Regex("a(∅|ε)a")
	fa, err := RegexToFA(regex)
	if err != nil {
		t.Fatalf("RegexToFA failed: %v", err)
	}

	// 检查字母表中是否包含 ∅
	for _, sym := range fa.Alphabet {
		if strings.Contains(string(sym), "∅") {
			t.Errorf("NFA alphabet should not contain ∅, got: %v", fa.Alphabet)
		}
		if strings.Contains(string(sym), "ε") {
			t.Errorf("NFA alphabet should not contain ε, got: %v", fa.Alphabet)
		}
	}

	// 验证化简后的 NFA 应该只接受 "aa"
	// 字母表应该只包含 'a'
	if len(fa.Alphabet) != 1 || string(fa.Alphabet[0]) != "a" {
		t.Errorf("Expected alphabet [a], got: %v", fa.Alphabet)
	}
}

func TestRegexToFAWithStepsNoEmptySet(t *testing.T) {
	// 测试 (a|b|∅)*c∅? 转换后的 NFA 不包含 ∅
	regex := model.Regex("(a|b|∅)*c∅?")
	result, err := RegexToFAWithSteps(regex)
	if err != nil {
		t.Fatalf("RegexToFAWithSteps failed: %v", err)
	}

	// 检查最终自动机的字母表
	for _, sym := range result.FinalAutomaton.Alphabet {
		if strings.Contains(string(sym), "∅") {
			t.Errorf("NFA alphabet should not contain ∅, got: %v", result.FinalAutomaton.Alphabet)
		}
		if strings.Contains(string(sym), "ε") {
			t.Errorf("NFA alphabet should not contain ε, got: %v", result.FinalAutomaton.Alphabet)
		}
	}

	// 化简后应该是 (a|b)*c，字母表应该是 [a, b, c]
	expectedAlphabet := []model.Symbol{"a", "b", "c"}
	if len(result.FinalAutomaton.Alphabet) != len(expectedAlphabet) {
		t.Errorf("Expected alphabet length %d, got %d: %v", len(expectedAlphabet), len(result.FinalAutomaton.Alphabet), result.FinalAutomaton.Alphabet)
	}
}
