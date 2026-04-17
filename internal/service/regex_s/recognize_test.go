package regex_s

import (
	"testing"
)

func TestSimplifyRegex(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantSimplified string
		wantEmpty     bool
	}{
		{"空集", "∅", "", true},
		{"空串", "ε", "", false},
		{"a(∅|ε)a", "a(∅|ε)a", "aa", false},
		{"(a|b|∅)*c∅?", "(a|b|∅)*c∅?", "(a|b)*c", false},
		{"a|∅", "a|∅", "a", false},
		{"∅|a", "∅|a", "a", false},
		{"a∅", "a∅", "", true},
		{"∅a", "∅a", "", true},
		{"ε*", "ε*", "", false},
		{"ε+", "ε+", "", false},
		{"ε?", "ε?", "", false},
		{"∅*", "∅*", "", false},
		{"∅+", "∅+", "", true},  // ∅+ = ∅，是空语言
		{"∅?", "∅?", "", false}, // ∅? = ε，匹配空串
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSimplified, gotEmpty := SimplifyRegex(tt.input)
			if gotSimplified != tt.wantSimplified {
				t.Errorf("SimplifyRegex(%q) gotSimplified = %q, want %q", tt.input, gotSimplified, tt.wantSimplified)
			}
			if gotEmpty != tt.wantEmpty {
				t.Errorf("SimplifyRegex(%q) gotEmpty = %v, want %v", tt.input, gotEmpty, tt.wantEmpty)
			}
		})
	}
}
