package automaton_s

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestFAToRegex(t *testing.T) {
	tests := []struct {
		name      string
		automaton *model.Automaton
		checkFunc func(model.Regex) error
	}{
		{
			name: "simple DFA accepting 'a'",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			checkFunc: func(regex model.Regex) error {
				if regex == "" {
					return fmt.Errorf("FAToRegex() should return non-empty regex")
				}

				return nil
			},
		},
		{
			name: "DFA accepting 'a*'",
			automaton: &model.Automaton{
				States:          []model.State{"q0"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q0"},
				Type:            model.DFA,
			},
			checkFunc: func(regex model.Regex) error {
				if regex == "" {
					return fmt.Errorf("FAToRegex() should return non-empty regex")
				}

				return nil
			},
		},
		{
			name: "NFA",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0", "q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.NFA,
			},
			checkFunc: func(regex model.Regex) error {
				if regex == "" {
					return fmt.Errorf("FAToRegex() should return non-empty regex")
				}

				return nil
			},
		},
		{
			name: "automaton with epsilon",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a", model.Epsilon},
				Transitions:     []model.Transition{{FromState: "q0", Input: model.Epsilon, ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.EpsilonNFA,
			},
			checkFunc: func(regex model.Regex) error {
				if regex == "" {
					return fmt.Errorf("FAToRegex() should return non-empty regex")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automaton := tt.automaton.Clone()
			regex := FAToRegex(automaton)

			if err := tt.checkFunc(regex); err != nil {
				t.Errorf("FAToRegex() failed: %v", err)
			}

			if len(regex) == 0 {
				t.Error("FAToRegex() should return non-empty regex")
			}
		})
	}
}

func TestFAToRegexWithNil(t *testing.T) {
	regex := FAToRegex(nil)
	if regex != "" {
		t.Errorf("FAToRegex(nil) should return empty string, got %q", regex)
	}
}

func TestFAToRegexWithEmpty(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{},
		Alphabet:        []model.Symbol{},
		Transitions:     []model.Transition{},
		InitialState:    "",
		AcceptingStates: []model.State{},
		Type:            model.DFA,
	}

	regex := FAToRegex(automaton)
	if regex == "" {
		t.Error("FAToRegex(empty) should return a regex, got empty string")
	}
}

func TestFAToRegexWithSingleState(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{"q0"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q0"},
		Type:            model.DFA,
	}

	regex := FAToRegex(automaton)
	if regex == "" {
		t.Error("FAToRegex() should return non-empty regex for single state DFA")
	}
}

func TestFAToRegexWithMultipleAccepting(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "a", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1", "q2"},
		Type:            model.DFA,
	}

	regex := FAToRegex(automaton)
	if regex == "" {
		t.Error("FAToRegex() should return non-empty regex for multiple accepting states")
	}
}

func TestFAToRegexWithNestedRegex(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a", "b"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q2"},
		Type:            model.DFA,
	}

	regex := FAToRegex(automaton)
	regexStr := string(regex)

	if !strings.Contains(regexStr, "a") {
		t.Error("FAToRegex() should contain 'a' in the regex")
	}

	if !strings.Contains(regexStr, "b") {
		t.Error("FAToRegex() should contain 'b' in the regex")
	}
}

func TestFAToRegexWithStar(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{"q0"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q0"},
		Type:            model.DFA,
	}

	regex := FAToRegex(automaton)
	regexStr := string(regex)

	if !strings.Contains(regexStr, "*") {
		t.Error("FAToRegex() should contain '*' for loop")
	}
}

func TestFAToRegexWithUnion(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a", "b"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q0", Input: "b", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1", "q2"},
		Type:            model.DFA,
	}

	regex := FAToRegex(automaton)
	regexStr := string(regex)

	if !strings.Contains(regexStr, "|") {
		t.Error("FAToRegex() should contain '|' for union")
	}
}
