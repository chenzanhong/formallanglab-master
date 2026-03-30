package regex_s

import (
	"fmt"
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegexToFA(t *testing.T) {
	tests := []struct {
		name      string
		regex     model.Regex
		checkFunc func(*model.Automaton, error) error
	}{
		{
			name:  "simple pattern: a",
			regex: "a",
			checkFunc: func(a *model.Automaton, err error) error {
				if err != nil {
					return err
				}
				if a == nil {
					return fmt.Errorf("RegexToFA() should not return nil")
				}
				if a.Type != model.EpsilonNFA {
					return fmt.Errorf("RegexToFA() should return EpsilonNFA, got %v", a.Type)
				}
				if len(a.States) < 2 {
					return fmt.Errorf("RegexToFA() should have at least 2 states, got %d", len(a.States))
				}
				if len(a.Transitions) == 0 {
					return fmt.Errorf("RegexToFA() should have at least 1 transition, got %d", len(a.Transitions))
				}

				return nil
			},
		},
		{
			name:  "pattern with star: a*",
			regex: "a*",
			checkFunc: func(a *model.Automaton, err error) error {
				if err != nil {
					return err
				}
				if a == nil {
					return fmt.Errorf("RegexToFA() should not return nil")
				}

				return nil
			},
		},
		{
			name:  "pattern with union: a|b",
			regex: "a|b",
			checkFunc: func(a *model.Automaton, err error) error {
				if err != nil {
					return err
				}
				if a == nil {
					return fmt.Errorf("RegexToFA() should not return nil")
				}

				return nil
			},
		},
		{
			name:  "epsilon pattern",
			regex: "ε",
			checkFunc: func(a *model.Automaton, err error) error {
				if err != nil {
					return err
				}
				if a == nil {
					return fmt.Errorf("RegexToFA() should not return nil for epsilon")
				}

				return nil
			},
		},
		{
			name:  "empty language",
			regex: model.EmptyLanguageToken,
			checkFunc: func(a *model.Automaton, err error) error {
				if err != nil {
					return err
				}
				if a == nil {
					return fmt.Errorf("RegexToFA() should not return nil for empty language")
				}

				return nil
			},
		},
		{
			name:  "complex pattern: (a|b)*c",
			regex: "(a|b)*c",
			checkFunc: func(a *model.Automaton, err error) error {
				if err != nil {
					return err
				}
				if a == nil {
					return fmt.Errorf("RegexToFA() should not return nil for complex pattern")
				}

				return nil
			},
		},
		{
			name:  "empty pattern",
			regex: "",
			checkFunc: func(a *model.Automaton, err error) error {
				if err == nil {
					return fmt.Errorf("RegexToFA() should return error for empty pattern")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RegexToFA(tt.regex)

			if err := tt.checkFunc(result, err); err != nil {
				t.Errorf("RegexToFA() failed: %v", err)
			}
		})
	}
}

func TestRegexToFAWithNil(t *testing.T) {
	result, err := RegexToFA("")
	if err == nil {
		t.Error("RegexToFA(empty) should return error")
	}
	if result != nil {
		t.Error("RegexToFA(empty) should return nil result")
	}
}

func TestRegexToFAWithInvalidRegex(t *testing.T) {
	result, err := RegexToFA("a|")
	if err == nil {
		t.Error("RegexToFA() should return error for invalid regex")
	}
	if result != nil {
		t.Error("RegexToFA() should return nil result for invalid regex")
	}
}

func TestRegexToFAWithPlus(t *testing.T) {
	result, err := RegexToFA("a+")
	if err != nil {
		t.Errorf("RegexToFA(a+) error = %v", err)
	}
	if result == nil {
		t.Error("RegexToFA(a+) should not return nil")
	}
}

func TestRegexToFAWithQuestion(t *testing.T) {
	result, err := RegexToFA("a?")
	if err != nil {
		t.Errorf("RegexToFA(a?) error = %v", err)
	}
	if result == nil {
		t.Error("RegexToFA(a?) should not return nil")
	}
}
