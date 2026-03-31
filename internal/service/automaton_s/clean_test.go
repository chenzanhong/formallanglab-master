package automaton_s

import (
	"fmt"
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestCleanup(t *testing.T) {
	tests := []struct {
		name      string
		automaton *model.Automaton
		checkFunc func(*model.Automaton) error
	}{
		{
			name: "automaton with unreachable states",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2", "q3"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			checkFunc: func(a *model.Automaton) error {
				if len(a.States) != 2 {
					return fmt.Errorf("Cleanup() should remove unreachable states, got %d states", len(a.States))
				}
				if len(a.Alphabet) != 1 {
					return fmt.Errorf("Cleanup() should remove unused symbols, got %d symbols", len(a.Alphabet))
				}

				return nil
			},
		},
		{
			name: "automaton with epsilon transitions",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", model.Epsilon},
				Transitions:     []model.Transition{{FromState: "q0", Input: model.Epsilon, ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "a", ToStates: []model.State{"q2"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.EpsilonNFA,
			},
			checkFunc: func(a *model.Automaton) error {
				if len(a.States) != 3 {
					return fmt.Errorf("Cleanup() should keep all reachable states with epsilon, got %d states", len(a.States))
				}

				return nil
			},
		},
		{
			name: "automaton with all states reachable",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q2"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.DFA,
			},
			checkFunc: func(a *model.Automaton) error {
				if len(a.States) != 3 {
					return fmt.Errorf("Cleanup() should keep all reachable states, got %d states", len(a.States))
				}
				if len(a.Alphabet) != 2 {
					return fmt.Errorf("Cleanup() should keep all used symbols, got %d symbols", len(a.Alphabet))
				}

				return nil
			},
		},
		{
			name: "nil automaton",
			automaton: &model.Automaton{
				States:          []model.State{},
				Alphabet:        []model.Symbol{},
				Transitions:     []model.Transition{},
				InitialState:    "",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			checkFunc: func(a *model.Automaton) error {
				if len(a.States) != 0 {
					return fmt.Errorf("Cleanup() should handle empty automaton, got %d states", len(a.States))
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automaton := tt.automaton.Clone()
			Cleanup(automaton)

			if err := tt.checkFunc(automaton); err != nil {
				t.Errorf("Cleanup() failed: %v", err)
			}

			if err := automaton.Validate(); err != nil {
				t.Errorf("Cleanup() result validation failed: %v", err)
			}
		})
	}
}

func TestCleanupWithNil(t *testing.T) {
	Cleanup(nil)
}
