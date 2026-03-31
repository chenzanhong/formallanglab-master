package regex_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegexEquivalenceCheck(t *testing.T) {
	tests := []struct {
		name     string
		pattern1 model.Regex
		pattern2 model.Regex
		want     bool
		wantErr  bool
	}{
		{
			name:     "identical patterns",
			pattern1: "a",
			pattern2: "a",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "equivalent patterns: a|b and b|a",
			pattern1: "a|b",
			pattern2: "b|a",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "equivalent patterns: a* and (a*)*",
			pattern1: "a*",
			pattern2: "(a*)*",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "equivalent patterns: ab and a.b",
			pattern1: "ab",
			pattern2: "a.b",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "not equivalent: a and b",
			pattern1: "a",
			pattern2: "b",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "not equivalent: a* and b*",
			pattern1: "a*",
			pattern2: "b*",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "empty language",
			pattern1: model.EmptyLanguageToken,
			pattern2: model.EmptyLanguageToken,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "one empty language",
			pattern1: model.EmptyLanguageToken,
			pattern2: "a",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "epsilon patterns",
			pattern1: "ε",
			pattern2: "ε",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "complex equivalent patterns",
			pattern1: "a(b|c)*",
			pattern2: "a(b*|c*)*",
			want:     true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegexEquivalenceCheck(tt.pattern1, tt.pattern2)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegexEquivalenceCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("RegexEquivalenceCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexEquivalenceCheckWithEmpty(t *testing.T) {
	// 空字符串正则
	got, err := RegexEquivalenceCheck("", "")
	if err != nil {
		t.Errorf("RegexEquivalenceCheck(empty, empty) error = %v", err)
	}
	if !got {
		t.Error("RegexEquivalenceCheck(empty, empty) should return true")
	}
}

func TestRegexEquivalenceCheckWithInvalidRegex(t *testing.T) {
	// 无效的正则表达式
	_, err := RegexEquivalenceCheck("a|", "b")
	if err == nil {
		t.Error("RegexEquivalenceCheck() should return error for invalid regex")
	}
}
