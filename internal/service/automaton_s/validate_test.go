package automaton_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestAutomatonValidate(t *testing.T) {
	tests := []struct {
		name      string
		automaton *model.Automaton
		wantErr   bool
	}{
		{
			name:      "nil automaton",
			automaton: nil,
			wantErr:   true,
		},
		{
			name: "empty states",
			automaton: &model.Automaton{
				States:          []model.State{},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{},
				InitialState:    "q0",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantErr: true,
		},
		{
			name: "empty state name",
			automaton: &model.Automaton{
				States:          []model.State{""},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{},
				InitialState:    "",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantErr: true,
		},
		{
			name: "invalid initial state",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{},
				InitialState:    "q2",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantErr: true,
		},
		{
			name: "invalid accepting state",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.DFA,
			},
			wantErr: true,
		},
		{
			name: "invalid transition source",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q2", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantErr: true,
		},
		{
			name: "invalid transition symbol",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "b", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantErr: true,
		},
		{
			name: "invalid transition target",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q2"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			wantErr: true,
		},
		{
			name: "valid DFA",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.DFA,
			},
			wantErr: false,
		},
		{
			name: "valid NFA with multiple targets",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1", "q2"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.NFA,
			},
			wantErr: false,
		},
		{
			name: "valid epsilon NFA",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b", model.Epsilon},
				Transitions:     []model.Transition{{FromState: "q0", Input: model.Epsilon, ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.EpsilonNFA,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AutomatonValidate(tt.automaton)
			if (err != nil) != tt.wantErr {
				t.Errorf("AutomatonValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
