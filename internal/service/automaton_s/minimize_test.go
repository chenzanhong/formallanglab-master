package automaton_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestDFAMinimize(t *testing.T) {
	tests := []struct {
		name    string
		dfa     *model.Automaton
		wantNil bool
	}{
		{
			name:    "nil DFA",
			dfa:     nil,
			wantNil: true,
		},
		{
			name: "minimal DFA with 1 state",
			dfa: &model.Automaton{
				States:          []model.State{"q0"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q0"},
				Type:            model.DFA,
			},
			wantNil: false,
		},
		{
			name: "DFA with unreachable states",
			dfa: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2", "q3"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			wantNil: false,
		},
		{
			name: "DFA with equivalent states",
			dfa: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q0", Input: "b", ToStates: []model.State{"q2"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1", "q2"},
				Type:            model.DFA,
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minimized := DFAMinimize(tt.dfa)
			if (minimized == nil) != tt.wantNil {
				t.Errorf("DFAMinimize() = %v, wantNil %v", minimized == nil, tt.wantNil)
			}

			if minimized != nil {
				if err := minimized.Validate(); err != nil {
					t.Errorf("DFAMinimize() result validation failed: %v", err)
				}
				if minimized.Type != model.DFA {
					t.Errorf("DFAMinimize() type = %v, want %v", minimized.Type, model.DFA)
				}
			}
		})
	}
}

func TestDFAMinimizeWithProcess(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a", "b"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q0", Input: "b", ToStates: []model.State{"q0"}}, {FromState: "q1", Input: "a", ToStates: []model.State{"q2"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q1"}}, {FromState: "q2", Input: "a", ToStates: []model.State{"q2"}}, {FromState: "q2", Input: "b", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.DFA,
	}

	minimized, process := DFAMinimizeWithProcess(dfa)

	if minimized == nil {
		t.Error("Minimized DFA should not be nil")
	}

	if process == nil {
		t.Error("Process should not be nil")
	}

	if len(process.Steps) == 0 {
		t.Error("Steps should not be empty")
	}

	for _, step := range process.Steps {
		if step.AutomatonFlow == nil {
			t.Error("AutomatonFlow should not be nil")
		}
	}
}

func TestDFAMinimizeCorrectness(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2", "q3"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "a", ToStates: []model.State{"q2"}}, {FromState: "q2", Input: "a", ToStates: []model.State{"q3"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q3"},
		Type:            model.DFA,
	}

	minimized := DFAMinimize(dfa)

	if minimized == nil {
		t.Fatal("Minimized DFA should not be nil")
	}

	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"a", false},
		{"aa", false},
		{"aaa", true},
		{"aaaa", false},
	}

	for _, tc := range testCases {
		result, err := Recognize(minimized, tc.input)
		if err != nil && result.IsAccepted != tc.expected {
			t.Errorf("Recognize(minimized, %q) = %v, want %v", tc.input, result.IsAccepted, tc.expected)
		}
	}
}

func TestDFAMinimizeEquivalence(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2", "q3"},
		Alphabet:        []model.Symbol{"a", "b"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q0", Input: "b", ToStates: []model.State{"q0"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q2"}}, {FromState: "q1", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q2", Input: "a", ToStates: []model.State{"q3"}}, {FromState: "q2", Input: "b", ToStates: []model.State{"q2"}}, {FromState: "q3", Input: "b", ToStates: []model.State{"q3"}}, {FromState: "q3", Input: "a", ToStates: []model.State{"q3"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q3"},
		Type:            model.DFA,
	}

	minimized := DFAMinimize(dfa)

	if minimized == nil {
		t.Fatal("Minimized DFA should not be nil")
	}

	testCases := []struct {
		input    string
		expected bool
	}{
		{"ab", false},
		{"aba", false},
		{"abab", true},
		{"ababa", false},
		{"ababab", true},
	}

	for _, tc := range testCases {
		result, err := Recognize(minimized, tc.input)
		if err != nil && result.IsAccepted != tc.expected {
			t.Errorf("Recognize(minimized, %q) = %v, want %v", tc.input, result.IsAccepted, tc.expected)
		}
	}
}
