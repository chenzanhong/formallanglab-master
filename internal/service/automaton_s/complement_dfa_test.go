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
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "a", ToStates: []model.State{"q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.DFA,
	}

	complemented, err := ComplementDFA(dfa)
	if err != nil {
		t.Fatalf("ComplementDFA() failed: %v", err)
	}

	// 调试信息
	t.Logf("Original DFA accepting states: %v", dfa.AcceptingStates)
	t.Logf("Complemented DFA accepting states: %v", complemented.AcceptingStates)
	t.Logf("Complemented DFA initial state: %v", complemented.InitialState)

	testCases := []struct {
		input    string
		expected bool
	}{
		{"", true},       // 空串应该被接受（原 DFA 不接受空串）
		{"a", false},     // "a"不应该被接受（原 DFA 接受"a"）
		{"aa", false},    // "aa"不应该被接受（原 DFA 接受"aa"）
	}

	for _, tc := range testCases {
		result, err := Recognize(complemented, tc.input)
		if err != nil {
			t.Logf("Recognize(complemented, %q) error: %v", tc.input, err)
		}
		if result.IsAccepted != tc.expected {
			t.Errorf("Recognize(complemented, %q) = %v, want %v", tc.input, result.IsAccepted, tc.expected)
		}
	}
}
