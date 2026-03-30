package regex_s

import (
	"errors"
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegexValidate(t *testing.T) {
	tests := []struct {
		name  string
		regex model.Regex
		want  error
	}{
		{
			name:  "valid pattern: a",
			regex: "a",
			want:  nil,
		},
		{
			name:  "valid pattern: a*",
			regex: "a*",
			want:  nil,
		},
		{
			name:  "valid pattern: a|b",
			regex: "a|b",
			want:  nil,
		},
		{
			name:  "valid pattern: (a|b)*c",
			regex: "(a|b)*c",
			want:  nil,
		},
		{
			name:  "valid pattern: ε",
			regex: "ε",
			want:  nil,
		},
		{
			name:  "valid pattern: ∅",
			regex: model.EmptyLanguageToken,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegexValidate(tt.regex)

			if !errors.Is(err, tt.want) {
				t.Errorf("RegexValidate() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestRegexValidString(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    bool
		wantErr bool
	}{
		{
			name:    "valid string: abc123",
			str:     "abc123",
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid string: empty",
			str:     "",
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid string: ε",
			str:     "ε",
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid string: ∅",
			str:     "∅",
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid string: a b",
			str:     "a b",
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid string: a+b",
			str:     "a+b",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegexValidString(tt.str)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegexValidString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("RegexValidString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexValidStringWithSpecialChars(t *testing.T) {
	_, err := RegexValidString("a*b")
	if err == nil {
		t.Error("RegexValidString() should return error for string with *")
	}
}
