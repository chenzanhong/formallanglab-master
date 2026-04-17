package regex_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegexGenerateExampleString(t *testing.T) {
	tests := []struct {
		name  string
		regex string
	}{
		{"aa", "aa"},
		{"a(∅|ε)a", "a(∅|ε)a"},         // 应该化简为 aa
		{"(a|b|∅)*c∅?", "(a|b|∅)*c∅?"}, // 应该化简为 (a|b)*c
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accept, reject := RegexGenerateExampleString(model.Regex(tt.regex))
			if len(accept) == 0 {
				t.Errorf("RegexGenerateExampleString(%q) accept is empty", tt.regex)
			}
			if len(reject) == 0 {
				t.Errorf("RegexGenerateExampleString(%q) reject is empty", tt.regex)
			}

			// 特别检查 a(∅|ε)a 是否化简为 aa
			if tt.regex == "a(∅|ε)a" {
				// aa 应该能匹配
				foundAA := false
				for _, s := range accept {
					if s == "aa" {
						foundAA = true
						break
					}
				}
				if !foundAA {
					t.Errorf("RegexGenerateExampleString(%q) should accept 'aa', but got %v", tt.regex, accept)
				}
			}
		})
	}
}
