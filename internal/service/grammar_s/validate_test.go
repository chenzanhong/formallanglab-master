package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestGrammarCheckValidity(t *testing.T) {
	tests := []struct {
		name    string
		grammar *model.Grammar
		wantErr bool
	}{
		{
			name:    "nil grammar",
			grammar: nil,
			wantErr: true,
		},
		{
			name: "valid simple grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
				},
			},
			wantErr: false,
		},
		{
			name: "empty non-terminals",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
				},
			},
			wantErr: true,
		},
		{
			name: "empty terminals",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			wantErr: true,
		},
		{
			name: "empty start symbol",
			grammar: &model.Grammar{
				StartSymbol:  "",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
				},
			},
			wantErr: true,
		},
		{
			name: "start symbol not in non-terminals",
			grammar: &model.Grammar{
				StartSymbol:  "X",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
				},
			},
			wantErr: true,
		},
		{
			name: "empty productions",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions:  []model.Production{},
			},
			wantErr: true,
		},
		{
			name: "production with empty left",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{}, Right: []model.Symbol{"a"}},
				},
			},
			wantErr: true,
		},
		{
			name: "production left without non-terminal",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"a"}, Right: []model.Symbol{"a"}},
				},
			},
			wantErr: true,
		},
		{
			name: "valid grammar with epsilon",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GrammarCheckValidity(tt.grammar)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrammarCheckValidity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
