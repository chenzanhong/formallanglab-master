package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestGrammarIsEquivalent(t *testing.T) {
	tests := []struct {
		name      string
		g1        *model.Grammar
		g2        *model.Grammar
		wantEquiv bool
		wantErr   bool
	}{
		{
			name: "nil grammar",
			g1:   nil,
			g2: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
				},
				GrammarType: model.RegularGrammar,
			},
			wantEquiv: false,
			wantErr:   true,
		},
		{
			name: "equivalent regular grammars",
			g1: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.RegularGrammar,
			},
			g2: &model.Grammar{
				StartSymbol:  "A",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.RegularGrammar,
			},
			wantEquiv: true,
			wantErr:   false,
		},
		{
			name: "non-equivalent regular grammars",
			g1: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.RegularGrammar,
			},
			g2: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.RegularGrammar,
			},
			wantEquiv: false,
			wantErr:   false,
		},
		{
			name: "same grammar",
			g1: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
				},
				GrammarType: model.RegularGrammar,
			},
			g2: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
				},
				GrammarType: model.RegularGrammar,
			},
			wantEquiv: true,
			wantErr:   false,
		},
		{
			name: "non-regular grammar (CFG)",
			g1: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S", "b"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.ContextFreeGrammar,
			},
			g2: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S", "b"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.ContextFreeGrammar,
			},
			wantEquiv: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GrammarIsEquivalent(tt.g1, tt.g2)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrammarIsEquivalent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantEquiv {
				t.Errorf("GrammarIsEquivalent() = %v, want %v", got, tt.wantEquiv)
			}
		})
	}
}
