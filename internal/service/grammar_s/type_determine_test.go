package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestTypeDetermine(t *testing.T) {
	tests := []struct {
		name      string
		grammar   *model.Grammar
		wantType  model.GrammarType
	}{
		{
			name: "nil grammar",
			grammar: nil,
			wantType: model.InvalidGrammar,
		},
		{
			name: "right linear regular grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b", "S"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
				},
			},
			wantType: model.RegularGrammar,
		},
		{
			name: "context free grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S", "b"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			wantType: model.ContextFreeGrammar,
		},
		{
			name: "context sensitive grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A", "B"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A", "B", "C"}},
					{Left: []model.Symbol{"A", "B"}, Right: []model.Symbol{"b", "A"}},
				},
			},
			wantType: model.ContextSensitiveGrammar,
		},
		{
			name: "simple CFG with epsilon",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			wantType: model.ContextFreeGrammar,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType := TypeDetermine(tt.grammar)
			if gotType != tt.wantType {
				t.Errorf("TypeDetermine() = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}

func TestIsRegular(t *testing.T) {
	tests := []struct {
		name         string
		grammar      *model.Grammar
		wantRegular  bool
		wantLinear   model.GrammarLinearity
	}{
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
			},
			wantRegular: true,
			wantLinear:  model.RightLinear,
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
			},
			wantRegular: true,
			wantLinear:  model.LeftLinear,
		},
		{
			name: "non-regular CFG",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S", "b"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			wantRegular: false,
			wantLinear:  model.InvalidLinearity,
		},
		{
			name: "mixed linear grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"A", "b"}},
				},
			},
			wantRegular: false,
			wantLinear:  model.InvalidLinearity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isRegular, linear := IsRegular(tt.grammar)
			if isRegular != tt.wantRegular {
				t.Errorf("IsRegular() = %v, want %v", isRegular, tt.wantRegular)
			}
			if isRegular && linear != tt.wantLinear {
				t.Errorf("IsRegular() linear = %v, want %v", linear, tt.wantLinear)
			}
		})
	}
}
