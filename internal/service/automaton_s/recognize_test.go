package automaton_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRecognizeDFA(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a", "b"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q2"},
		Type:            model.DFA,
	}

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"single a", "a", false},
		{"ab", "ab", true},
		{"abc", "abc", false},
		{"aba", "aba", false},
		{"abb", "abb", false},
		{"abab", "abab", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Recognize(dfa, tc.input)
			if err != nil && result.IsAccepted != tc.expected {
				t.Errorf("Recognize(dfa, %q) = %v, want %v", tc.input, result.IsAccepted, tc.expected)
			}
		})
	}
}

func TestRecognizeNFA(t *testing.T) {
	nfa := &model.Automaton{
		States:          []model.State{"q0", "q1"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0", "q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.NFA,
	}

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"single a", "a", true},
		{"aa", "aa", true},
		{"aaa", "aaa", true},
		{"b", "b", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Recognize(nfa, tc.input)
			if err != nil && result.IsAccepted != tc.expected {
				t.Errorf("Recognize(nfa, %q) = %v, want %v", tc.input, result.IsAccepted, tc.expected)
			}
		})
	}
}

func TestRecognizeEpsilonNFA(t *testing.T) {
	epsilonNFA := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a", "b", model.Epsilon},
		Transitions:     []model.Transition{{FromState: "q0", Input: model.Epsilon, ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "a", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q2"},
		Type:            model.EpsilonNFA,
	}

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"a", "a", true},
		{"aa", "aa", false},
		{"b", "b", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Recognize(epsilonNFA, tc.input)
			if err != nil && result.IsAccepted != tc.expected {
				t.Errorf("Recognize(epsilonNFA, %q) = %v, want %v", tc.input, result.IsAccepted, tc.expected)
			}
		})
	}
}

func TestRecognizeNilAutomaton(t *testing.T) {
	result, err := Recognize(nil, "a")
	if err == nil {
		t.Error("Recognize(nil, \"a\") should return error")
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
	if result.IsAccepted {
		t.Error("IsAccepted should be false for nil automaton")
	}
}

func TestRecognizeEmptyAutomaton(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{},
		InitialState:    "",
		AcceptingStates: []model.State{},
		Type:            model.DFA,
	}

	result, err := Recognize(automaton, "a")
	if err == nil {
		t.Error("Recognize(empty automaton, \"a\") should return error")
	}
	if result.IsAccepted {
		t.Error("IsAccepted should be false for empty automaton")
	}
}
