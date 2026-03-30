package automaton_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestUnCompleteDFA(t *testing.T) {
	tests := []struct {
		name        string
		automaton   *model.Automaton
		sink        model.State
		wantHasSink bool
	}{
		{
			name: "automaton without sink state",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			sink:        model.SinkState,
			wantHasSink: false,
		},
		{
			name: "automaton with sink state",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", model.SinkState},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q0", Input: "b", ToStates: []model.State{model.SinkState}}, {FromState: "q1", Input: "b", ToStates: []model.State{model.SinkState}}, {FromState: model.SinkState, Input: "a", ToStates: []model.State{model.SinkState}}, {FromState: model.SinkState, Input: "b", ToStates: []model.State{model.SinkState}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			sink:        model.SinkState,
			wantHasSink: false,
		},
		{
			name: "custom sink state",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "dead"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q0", Input: "b", ToStates: []model.State{"dead"}}, {FromState: "dead", Input: "a", ToStates: []model.State{"dead"}}, {FromState: "dead", Input: "b", ToStates: []model.State{"dead"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			sink:        "dead",
			wantHasSink: false,
		},
		{
			name: "empty automaton",
			automaton: &model.Automaton{
				States:          []model.State{},
				Alphabet:        []model.Symbol{},
				Transitions:     []model.Transition{},
				InitialState:    "",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
			sink:        model.SinkState,
			wantHasSink: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automaton := tt.automaton.Clone()
			result := UnCompleteDFA(automaton, tt.sink)

			if result == nil {
				t.Error("UnCompleteDFA() returned nil")
			}

			hasSink := false
			for _, s := range result.States {
				if s == tt.sink {
					hasSink = true
					break
				}
			}

			if hasSink != tt.wantHasSink {
				t.Errorf("UnCompleteDFA() hasSink = %v, want %v", hasSink, tt.wantHasSink)
			}

			if hasSink {
				for _, tr := range result.Transitions {
					if tr.ToStates[0] == tt.sink {
						t.Errorf("UnCompleteDFA() should remove transitions to sink state")
					}
				}
			}
		})
	}
}

func TestUnCompleteDFAWithNil(t *testing.T) {
	result := UnCompleteDFA(nil, model.SinkState)
	if result != nil {
		t.Error("UnCompleteDFA(nil) should return nil")
	}
}

func TestCompleteDFAIntegration(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.DFA,
	}

	err := dfa.CompleteDFA()
	if err != nil {
		t.Errorf("CompleteDFA() failed: %v", err)
	}

	hasSink := false
	for _, s := range dfa.States {
		if s == model.SinkState {
			hasSink = true
			break
		}
	}

	if !hasSink {
		t.Error("CompleteDFA() should add sink state for incomplete DFA")
	}

	result := UnCompleteDFA(dfa, model.SinkState)
	hasSinkAfter := false
	for _, s := range result.States {
		if s == model.SinkState {
			hasSinkAfter = true
			break
		}
	}

	if hasSinkAfter {
		t.Error("UnCompleteDFA() should remove sink state")
	}
}
