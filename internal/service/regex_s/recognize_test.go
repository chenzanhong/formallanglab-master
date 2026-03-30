package regex_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegexRecognize(t *testing.T) {
	tests := []struct {
		name    string
		pattern model.Regex
		str     string
		want    bool
		wantErr bool
	}{
		{
			name:    "simple match: a with a",
			pattern: "a",
			str:     "a",
			want:    true,
			wantErr: false,
		},
		{
			name:    "simple no match: a with b",
			pattern: "a",
			str:     "b",
			want:    false,
			wantErr: false,
		},
		{
			name:    "star match: a* with aaa",
			pattern: "a*",
			str:     "aaa",
			want:    true,
			wantErr: false,
		},
		{
			name:    "star match: a* with empty",
			pattern: "a*",
			str:     "",
			want:    true,
			wantErr: false,
		},
		{
			name:    "union match: a|b with a",
			pattern: "a|b",
			str:     "a",
			want:    true,
			wantErr: false,
		},
		{
			name:    "union match: a|b with b",
			pattern: "a|b",
			str:     "b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "concat match: ab with ab",
			pattern: "ab",
			str:     "ab",
			want:    true,
			wantErr: false,
		},
		{
			name:    "epsilon match: ε with empty",
			pattern: "ε",
			str:     "",
			want:    true,
			wantErr: false,
		},
		{
			name:    "epsilon no match: ε with a",
			pattern: "ε",
			str:     "a",
			want:    false,
			wantErr: false,
		},
		{
			name:    "empty language: ∅ with empty",
			pattern: model.EmptyLanguageToken,
			str:     "",
			want:    false,
			wantErr: false,
		},
		{
			name:    "empty language: ∅ with a",
			pattern: model.EmptyLanguageToken,
			str:     "a",
			want:    false,
			wantErr: false,
		},
		{
			name:    "complex pattern: (a|b)*c with abc",
			pattern: "(a|b)*c",
			str:     "abc",
			want:    true,
			wantErr: false,
		},
		{
			name:    "plus match: a+ with a",
			pattern: "a+",
			str:     "a",
			want:    true,
			wantErr: false,
		},
		{
			name:    "plus no match: a+ with empty",
			pattern: "a+",
			str:     "",
			want:    false,
			wantErr: false,
		},
		{
			name:    "question match: a? with empty",
			pattern: "a?",
			str:     "",
			want:    true,
			wantErr: false,
		},
		{
			name:    "question match: a? with a",
			pattern: "a?",
			str:     "a",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegexRecognize(tt.pattern, tt.str)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegexRecognize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("RegexRecognize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexRecognizeWithInvalidRegex(t *testing.T) {
	_, err := RegexRecognize("a|", "a")
	if err == nil {
		t.Error("RegexRecognize() should return error for invalid regex")
	}
}

func TestRegexRecognizeWithNil(t *testing.T) {
	// 空模式
	got, err := RegexRecognize("", "a")
	if err != nil {
		t.Errorf("RegexRecognize(empty, a) error = %v", err)
	}
	if got {
		t.Error("RegexRecognize(empty, a) should return false")
	}
}
