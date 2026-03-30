package regex_s

import (
	"fmt"
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegexGenerateExampleString(t *testing.T) {
	tests := []struct {
		name      string
		regex     model.Regex
		checkFunc func(accept, reject []string) error
	}{
		{
			name:  "simple pattern: a",
			regex: "a",
			checkFunc: func(accept, reject []string) error {
				if len(accept) == 0 {
					return fmt.Errorf("should have at least one accepted string")
				}
				if len(reject) == 0 {
					return fmt.Errorf("should have at least one rejected string")
				}

				return nil
			},
		},
		{
			name:  "pattern with star: a*",
			regex: "a*",
			checkFunc: func(accept, reject []string) error {
				if len(accept) == 0 {
					return fmt.Errorf("should have at least one accepted string")
				}

				return nil
			},
		},
		{
			name:  "pattern with union: a|b",
			regex: "a|b",
			checkFunc: func(accept, reject []string) error {
				if len(accept) == 0 {
					return fmt.Errorf("should have at least one accepted string")
				}

				return nil
			},
		},
		{
			name:  "empty language",
			regex: model.EmptyLanguageToken,
			checkFunc: func(accept, reject []string) error {
				if len(accept) != 0 {
					return fmt.Errorf("empty language should have no accepted strings")
				}
				if len(reject) != 0 {
					return fmt.Errorf("empty language should have no rejected strings")
				}

				return nil
			},
		},
		{
			name:  "epsilon",
			regex: "ε",
			checkFunc: func(accept, reject []string) error {
				if len(accept) == 0 {
					return fmt.Errorf("epsilon should have at least one accepted string (empty string)")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accept, reject := RegexGenerateExampleString(tt.regex)

			if accept == nil {
				t.Error("RegexGenerateExampleString() accept should not be nil")
			}

			if reject == nil {
				t.Error("RegexGenerateExampleString() reject should not be nil")
			}

			if err := tt.checkFunc(accept, reject); err != nil {
				t.Errorf("RegexGenerateExampleString() failed: %v", err)
			}
		})
	}
}

func TestRegexGenerateExampleStringWithInvalidRegex(t *testing.T) {
	// 无效的正则表达式
	accept, reject := RegexGenerateExampleString("a|")
	if len(accept) == 0 || len(reject) == 0 {
		t.Error("RegexGenerateExampleString() should return results even for invalid regex")
	}
}
