package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestIsAmbiguousRegular(t *testing.T) {
	tests := []struct {
		name          string
		grammar       *model.Grammar
		wantAmbiguous bool
		wantErr       bool
	}{
		{
			name: "unambiguous regular grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
				},
				GrammarType: model.RegularGrammar,
			},
			wantAmbiguous: false,
			wantErr:       false,
		},
		{
			name: "ambiguous regular grammar - same prefix",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S", "A", "B"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "B"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
					{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a"}},
				},
				GrammarType:      model.RegularGrammar,
				GrammarLinearity: model.RightLinear,
			},
			wantAmbiguous: false,
			wantErr:       false,
		},
		{
			name: "simple unambiguous grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.RegularGrammar,
			},
			wantAmbiguous: false,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsAmbiguousRegular(tt.grammar)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAmbiguousRegular() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantAmbiguous {
				t.Errorf("IsAmbiguousRegular() = %v, want %v", got, tt.wantAmbiguous)
			}
		})
	}
}
