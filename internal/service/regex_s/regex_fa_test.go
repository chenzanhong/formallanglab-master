package regex_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegexToFA(t *testing.T) {
	tests := []struct {
		name    string
		regex   model.Regex
		wantErr bool
	}{
		{"aa", model.Regex("aa"), false},
		{"a(∅|ε)a", model.Regex("a(∅|ε)a"), false},         // 应该化简为 aa
		{"(a|b|∅)*c∅?", model.Regex("(a|b|∅)*c∅?"), false}, // 应该化简为 (a|b)*c
		{"∅", model.Regex("∅"), false},                     // 空语言自动机
		{"ε", model.Regex("ε"), false},                     // ε 化简后返回接受空串的自动机
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa, err := RegexToFA(tt.regex)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegexToFA(%q) error = %v, wantErr %v", tt.regex, err, tt.wantErr)
				return
			}
			if !tt.wantErr && fa == nil {
				t.Errorf("RegexToFA(%q) fa = nil, want non-nil", tt.regex)
			}
		})
	}
}
