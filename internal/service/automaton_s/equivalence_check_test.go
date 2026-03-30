package automaton_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestAutomatonEquivalenceCheck(t *testing.T) {
	tests := []struct {
		name           string
		dfa1           *model.Automaton
		dfa2           *model.Automaton
		wantEquivalent bool
	}{
		{
			name: "two identical DFAs",
			dfa1: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			dfa2: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			wantEquivalent: true,
		},
		{
			name: "two equivalent DFAs with different states",
			dfa1: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			dfa2: &model.Automaton{
				States:          []model.State{"p0", "p1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "p0", Input: "a", ToStates: []model.State{"p1"}}},
				InitialState:    "p0",
				AcceptingStates: []model.State{"p1"},
				Type:            model.DFA,
			},
			wantEquivalent: true,
		},
		{
			name: "two non-equivalent DFAs",
			dfa1: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			dfa2: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantEquivalent: false,
		},
		{
			name:           "nil inputs",
			dfa1:           nil,
			dfa2:           nil,
			wantEquivalent: true,
		},
		{
			name: "one nil input",
			dfa1: nil,
			dfa2: &model.Automaton{
				States:          []model.State{"q0"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{},
				InitialState:    "q0",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantEquivalent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minDFA1, minDFA2, isEquivalent := AutomatonEquivalenceCheck(tt.dfa1, tt.dfa2)

			if isEquivalent != tt.wantEquivalent {
				t.Errorf("AutomatonEquivalenceCheck() isEquivalent = %v, want %v", isEquivalent, tt.wantEquivalent)
			}

			if tt.dfa1 != nil && tt.dfa2 != nil {
				if minDFA1 == nil {
					t.Error("minDFA1 should not be nil")
				}
				if minDFA2 == nil {
					t.Error("minDFA2 should not be nil")
				}
			}
		})
	}
}

func TestAutomatonEquivalenceCheckWithNFA(t *testing.T) {
	nfa := &model.Automaton{
		States:          []model.State{"q0", "q1"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0", "q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.NFA,
	}

	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "a", ToStates: []model.State{"q2"}}, {FromState: "q2", Input: "a", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1", "q2"},
		Type:            model.DFA,
	}

	_, _, isEquivalent := AutomatonEquivalenceCheck(nfa, dfa)

	if !isEquivalent {
		t.Error("NFA and equivalent DFA should be considered equivalent")
	}
}
