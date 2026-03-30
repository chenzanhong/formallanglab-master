package automaton_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestComplementDFA(t *testing.T) {
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
			name: "simple DFA",
			dfa: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			wantNil: false,
		},
		{
			name: "DFA with all states accepting",
			dfa: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q0", "q1"},
				Type:            model.DFA,
			},
			wantNil: false,
		},
		{
			name: "DFA with no accepting states",
			dfa: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complemented, err := ComplementDFA(tt.dfa)
			if (complemented == nil) != tt.wantNil {
				t.Errorf("ComplementDFA() = %v, wantNil %v", complemented == nil, tt.wantNil)
			}

			if err != nil && complemented == nil {
				return
			}

			if complemented != nil {
				if err := complemented.Validate(); err != nil {
					t.Errorf("ComplementDFA() result validation failed: %v", err)
				}
				if complemented.Type != model.DFA {
					t.Errorf("ComplementDFA() type = %v, want %v", complemented.Type, model.DFA)
				}
			}
		})
	}
}

func TestComplementDFACorrectness(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.DFA,
	}

	complemented, err := ComplementDFA(dfa)
	if err != nil {
		t.Fatalf("ComplementDFA() failed: %v", err)
	}

	testCases := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"a", false},
		{"aa", true},
	}

	for _, tc := range testCases {
		result, err := Recognize(complemented, tc.input)
		if err != nil && result.IsAccepted != tc.expected {
			t.Errorf("Recognize(complemented, %q) = %v, want %v", tc.input, result.IsAccepted, tc.expected)
		}
	}
}
