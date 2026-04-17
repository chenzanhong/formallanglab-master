package automaton_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestNFAToDFA(t *testing.T) {
	tests := []struct {
		name    string
		nfa     *model.Automaton
		wantNil bool
	}{
		{
			name:    "nil NFA",
			nfa:     nil,
			wantNil: true,
		},
		{
			name: "simple NFA without epsilon",
			nfa: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0", "q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.NFA,
			},
			wantNil: false,
		},
		{
			name: "NFA with epsilon",
			nfa: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b", model.Epsilon},
				Transitions:     []model.Transition{{FromState: "q0", Input: model.Epsilon, ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.EpsilonNFA,
			},
			wantNil: false,
		},
		{
			name: "DFA input (should return nil)",
			nfa: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			wantNil: true,
		},
		{
			name: "single state NFA",
			nfa: &model.Automaton{
				States:          []model.State{"q0"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{}, // 没有转移
				InitialState:    "q0",
				AcceptingStates: []model.State{"q0"},
				Type:            model.NFA,
			},
			wantNil: true, // 没有转移的 NFA 实际上已经是 DFA，应该返回 nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调试：检查输入 NFA 的类型
			t.Logf("Input NFA type before NFAToDFA: %v", tt.nfa.Type)
			dfa := NFAToDFA(tt.nfa)
			t.Logf("NFAToDFA() result: %v", dfa)
			if (dfa == nil) != tt.wantNil {
				t.Errorf("NFAToDFA() = %v, wantNil %v", dfa == nil, tt.wantNil)
			}

			if dfa != nil && dfa.Type != model.DFA {
				t.Errorf("NFAToDFA() type = %v, want %v", dfa.Type, model.DFA)
			}

			if dfa != nil {
				if err := dfa.Validate(); err != nil {
					t.Errorf("NFAToDFA() result validation failed: %v", err)
				}
			}
		})
	}
}

func TestNFAToDFAWithProcess(t *testing.T) {
	nfa := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a", "b", model.Epsilon},
		Transitions:     []model.Transition{{FromState: "q0", Input: model.Epsilon, ToStates: []model.State{"q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q2"},
		Type:            model.EpsilonNFA,
	}

	process := NFAToDFAWithProcess(nfa)

	if process == nil {
		t.Error("NFAToDFAWithProcess() returned nil")
	}

	if process.FinalAutomaton == nil {
		t.Error("FinalAutomaton is nil")
	}

	if process.FinalAutomaton.Type != model.DFA {
		t.Errorf("FinalAutomaton type = %v, want %v", process.FinalAutomaton.Type, model.DFA)
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

func TestNFAToDFACorrectness(t *testing.T) {
	nfa := &model.Automaton{
		States:          []model.State{"q0", "q1"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0", "q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.NFA,
	}

	dfa := NFAToDFA(nfa)

	if dfa == nil {
		t.Fatal("DFA should not be nil")
	}

	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"a", true},
		{"aa", true},
		{"aaa", true},
		{"b", false},
	}

	for _, tc := range testCases {
		result, err := Recognize(dfa, tc.input)
		if err != nil && result.IsAccepted != tc.expected {
			t.Errorf("Recognize(dfa, %q) = %v, want %v", tc.input, result.IsAccepted, tc.expected)
		}
	}
}
