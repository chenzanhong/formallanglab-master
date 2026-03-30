package automaton_s

import (
	"fmt"
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestFAToGrammar(t *testing.T) {
	tests := []struct {
		name      string
		automaton *model.Automaton
		checkFunc func(*model.Grammar) error
	}{
		{
			name: "simple DFA",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			checkFunc: func(g *model.Grammar) error {
				if len(g.NonTerminals) != 2 {
					return fmt.Errorf("FAToGrammar() should have 2 non-terminals, got %d", len(g.NonTerminals))
				}
				if len(g.Terminals) != 2 {
					return fmt.Errorf("FAToGrammar() should have 2 terminals, got %d", len(g.Terminals))
				}
				if g.StartSymbol != model.Symbol("q0") {
					return fmt.Errorf("FAToGrammar() start symbol should be 'q0', got %v", g.StartSymbol)
				}
				if len(g.Productions) != 2 {
					return fmt.Errorf("FAToGrammar() should have 2 productions, got %d", len(g.Productions))
				}

				return nil
			},
		},
		{
			name: "DFA with epsilon production",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q0", "q1"},
				Type:            model.DFA,
			},
			checkFunc: func(g *model.Grammar) error {
				hasEpsilon := false
				for _, p := range g.Productions {
					if len(p.Right) == 1 && p.Right[0] == model.Epsilon {
						hasEpsilon = true
						break
					}
				}
				if !hasEpsilon {
					return fmt.Errorf("FAToGrammar() should have epsilon production for accepting state")
				}

				return nil
			},
		},
		{
			name: "NFA with multiple transitions",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1", "q2"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.NFA,
			},
			checkFunc: func(g *model.Grammar) error {
				if len(g.Productions) != 3 {
					return fmt.Errorf("FAToGrammar() should have 3 productions for NFA with multiple targets, got %d", len(g.Productions))
				}

				return nil
			},
		},
		{
			name: "automaton with epsilon transition",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a", model.Epsilon},
				Transitions:     []model.Transition{{FromState: "q0", Input: model.Epsilon, ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.EpsilonNFA,
			},
			checkFunc: func(g *model.Grammar) error {
				hasEpsilonProd := false
				for _, p := range g.Productions {
					if len(p.Right) == 1 && p.Right[0] == model.Epsilon {
						hasEpsilonProd = true
						break
					}
				}
				if !hasEpsilonProd {
					return fmt.Errorf("FAToGrammar() should have epsilon production for epsilon transition")
				}

				return nil
			},
		},
		{
			name: "single state DFA",
			automaton: &model.Automaton{
				States:          []model.State{"q0"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q0"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q0"},
				Type:            model.DFA,
			},
			checkFunc: func(g *model.Grammar) error {
				if len(g.Productions) != 2 {
					return fmt.Errorf("FAToGrammar() should have 2 productions for single state, got %d", len(g.Productions))
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automaton := tt.automaton.Clone()
			grammar := FAToGrammar(automaton)

			if err := tt.checkFunc(grammar); err != nil {
				t.Errorf("FAToGrammar() failed: %v", err)
			}

			if grammar == nil {
				t.Error("FAToGrammar() should not return nil")
			}

			if len(grammar.NonTerminals) == 0 {
				t.Error("FAToGrammar() should have non-terminals")
			}

			if len(grammar.Terminals) == 0 {
				t.Error("FAToGrammar() should have terminals")
			}

			if grammar.StartSymbol == "" {
				t.Error("FAToGrammar() should have start symbol")
			}

			if len(grammar.Productions) == 0 {
				t.Error("FAToGrammar() should have productions")
			}
		})
	}
}

func TestFAToGrammarWithNil(t *testing.T) {
	grammar := FAToGrammar(nil)
	if grammar == nil {
		t.Error("FAToGrammar(nil) should not return nil")
	}
}

func TestFAToGrammarWithEmpty(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{},
		Alphabet:        []model.Symbol{},
		Transitions:     []model.Transition{},
		InitialState:    "",
		AcceptingStates: []model.State{},
		Type:            model.DFA,
	}

	grammar := FAToGrammar(automaton)
	if grammar == nil {
		t.Error("FAToGrammar(empty) should not return nil")
	}
}

func TestFAToGrammarProductionStructure(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2"},
		Alphabet:        []model.Symbol{"a", "b"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q2"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q2"},
		Type:            model.DFA,
	}

	grammar := FAToGrammar(automaton)

	for _, p := range grammar.Productions {
		if len(p.Left) != 1 {
			t.Errorf("Production left should have exactly 1 symbol, got %d", len(p.Left))
		}

		if len(p.Right) != 1 && len(p.Right) != 2 {
			t.Errorf("Production right should have 1 or 2 symbols, got %d", len(p.Right))
		}

		if len(p.Right) == 2 {
			if !isTerminal(grammar.Terminals, p.Right[0]) && !isNonTerminal(grammar.NonTerminals, p.Right[1]) {
				t.Errorf("Production right should be terminal followed by non-terminal, got %v", p.Right)
			}
		} else if len(p.Right) == 1 {
			if p.Right[0] != model.Epsilon && !isNonTerminal(grammar.NonTerminals, p.Right[0]) {
				t.Errorf("Production right with 1 symbol should be epsilon or non-terminal, got %v", p.Right[0])
			}
		}
	}
}

func TestFAToGrammarEpsilonProduction(t *testing.T) {
	automaton := &model.Automaton{
		States:          []model.State{"q0", "q1"},
		Alphabet:        []model.Symbol{"a"},
		Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
		InitialState:    "q0",
		AcceptingStates: []model.State{"q0"},
		Type:            model.DFA,
	}

	grammar := FAToGrammar(automaton)

	hasEpsilon := false
	for _, p := range grammar.Productions {
		if len(p.Right) == 1 && p.Right[0] == model.Epsilon {
			hasEpsilon = true
			break
		}
	}

	if !hasEpsilon {
		t.Error("FAToGrammar() should have epsilon production for accepting state q0")
	}
}

func isNonTerminal(nonTerminals []model.Symbol, sym model.Symbol) bool {
	for _, n := range nonTerminals {
		if n == sym {
			return true
		}
	}

	return false
}

func isTerminal(terminals []model.Symbol, sym model.Symbol) bool {
	for _, t := range terminals {
		if t == sym {
			return true
		}
	}

	return false
}
