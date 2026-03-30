package automaton_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestAutomatonGenerateExampleString(t *testing.T) {
	tests := []struct {
		name      string
		automaton *model.Automaton
	}{
		{
			name: "DFA with simple language",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
		},
		{
			name: "NFA",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0", "q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.NFA,
			},
		},
		{
			name: "DFA with no accepting states",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{},
				Type:            model.DFA,
			},
		},
		{
			name: "DFA with all states accepting",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q0", "q1"},
				Type:            model.DFA,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accept, reject := AutomatonGenerateExampleString(tt.automaton)

			if accept == nil {
				t.Error("AutomatonGenerateExampleString() accept should not be nil")
			}

			if reject == nil {
				t.Error("AutomatonGenerateExampleString() reject should not be nil")
			}

			if len(accept) > 6 {
				t.Errorf("AutomatonGenerateExampleString() accept should have at most 6 examples, got %d", len(accept))
			}

			if len(reject) > 6 {
				t.Errorf("AutomatonGenerateExampleString() reject should have at most 6 examples, got %d", len(reject))
			}
		})
	}
}

func TestBFSShortestAcceptedStringsForDFA(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a", "b"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q2"},
		Type:            model.DFA,
	}

	results := BFSShortestAcceptedStringsForDFA(dfa, 5)

	if results == nil {
		t.Error("BFSShortestAcceptedStringsForDFA() should not return nil")
	}

	if len(results) == 0 {
		t.Error("BFSShortestAcceptedStringsForDFA() should find some accepted strings")
	}

	for _, s := range results {
		result, err := Recognize(dfa, s)
		if err != nil {
			t.Errorf("String %q should be accepted but got error: %v", s, err)
		}
		if !result.IsAccepted {
			t.Errorf("String %q should be accepted", s)
		}
	}
}

func TestBFSShortestAcceptedStringsForDFAWithNil(t *testing.T) {
	results := BFSShortestAcceptedStringsForDFA(nil, 5)
	if results != nil {
		t.Error("BFSShortestAcceptedStringsForDFA(nil, 5) should return nil")
	}
}

func TestBFSShortestAcceptedStringsForDFAWithEmpty(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{},
		Alphabet:        []model.Symbol{},
		Transitions:     []model.Transition{},
		InitialState:    "",
		AcceptingStates: []model.State{},
		Type:            model.DFA,
	}

	results := BFSShortestAcceptedStringsForDFA(dfa, 5)
	if results != nil {
		t.Error("BFSShortestAcceptedStringsForDFA(empty, 5) should return nil")
	}
}

func TestBFSShortestAcceptedStringsForDFAWithKZero(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q0"},
		Type:            model.DFA,
	}

	results := BFSShortestAcceptedStringsForDFA(dfa, 0)
	if results != nil {
		t.Error("BFSShortestAcceptedStringsForDFA(dfa, 0) should return nil")
	}
}

func TestGenerateExampleStringsFromCompletedDFA(t *testing.T) {
	dfa := &model.Automaton{
		States:          []model.State{"q0", "q1"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.DFA,
	}

	accept, reject := GenerateExampleStringsFromCompletedDFA(dfa)

	if accept == nil {
		t.Error("GenerateExampleStringsFromCompletedDFA() accept should not be nil")
	}

	if reject == nil {
		t.Error("GenerateExampleStringsFromCompletedDFA() reject should not be nil")
	}
}

func TestGenerateExampleStringsFromCompletedDFAWithNFA(t *testing.T) {
	nfa := &model.Automaton{
		States:          []model.State{"q0", "q1"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0", "q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q1"},
		Type:            model.NFA,
	}

	accept, reject := GenerateExampleStringsFromCompletedDFA(nfa)

	if accept == nil {
		t.Error("GenerateExampleStringsFromCompletedDFA() accept should not be nil")
	}

	if reject == nil {
		t.Error("GenerateExampleStringsFromCompletedDFA() reject should not be nil")
	}
}
