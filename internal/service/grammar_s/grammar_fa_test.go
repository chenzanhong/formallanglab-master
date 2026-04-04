package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRegularGrammarToFA(t *testing.T) {
	tests := []struct {
		name      string
		grammar   *model.Grammar
		wantErr   bool
		checkFunc func(*model.Automaton) bool
	}{
		{
			name:    "nil grammar",
			grammar: nil,
			wantErr: true,
		},
		{
			name: "right linear grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b", "S"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
				},
				GrammarType:      model.RegularGrammar,
				GrammarLinearity: model.RightLinear,
			},
			wantErr: false,
			checkFunc: func(a *model.Automaton) bool {
				return a != nil && len(a.States) > 0 && len(a.Transitions) > 0
			},
		},
		{
			name: "left linear grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A", "a"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"S", "b"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
				},
				GrammarType:      model.RegularGrammar,
				GrammarLinearity: model.LeftLinear,
			},
			wantErr: false,
			checkFunc: func(a *model.Automaton) bool {
				return a != nil && len(a.States) > 0 && len(a.Transitions) > 0
			},
		},
		{
			name: "grammar with epsilon production",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType:      model.RegularGrammar,
				GrammarLinearity: model.RightLinear,
			},
			wantErr: false,
			checkFunc: func(a *model.Automaton) bool {
				return a != nil && len(a.AcceptingStates) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegularGrammarToFA(tt.grammar)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegularGrammarToFA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc != nil && !tt.checkFunc(got) {
				t.Errorf("RegularGrammarToFA() check failed for result: %v", got)
			}
		})
	}
}

func TestRightLinearGrammarToFA(t *testing.T) {
	g := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S", "A"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
		},
		GrammarType:      model.RegularGrammar,
		GrammarLinearity: model.RightLinear,
	}

	automaton := rightLinearGrammarToFA(g)

	if automaton == nil {
		t.Fatal("rightLinearGrammarToFA() returned nil")
	}

	if automaton.InitialState != model.State("S") {
		t.Errorf("Initial state = %v, want S", automaton.InitialState)
	}

	if len(automaton.Transitions) == 0 {
		t.Error("rightLinearGrammarToFA() should create transitions")
	}
}

func TestLeftLinearGrammarToFA(t *testing.T) {
	g := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S", "A"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A", "a"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
		},
		GrammarType:      model.RegularGrammar,
		GrammarLinearity: model.LeftLinear,
	}

	automaton := leftLinearGrammarToFA(g)

	if automaton == nil {
		t.Fatal("leftLinearGrammarToFA() returned nil")
	}

	if automaton.InitialState != model.UniqueInitialState {
		t.Errorf("Initial state = %v, want %v", automaton.InitialState, model.UniqueInitialState)
	}

	if len(automaton.AcceptingStates) == 0 {
		t.Error("leftLinearGrammarToFA() should have accepting states")
	}
}
