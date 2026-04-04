package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestDetermineLinearity(t *testing.T) {
	tests := []struct {
		name       string
		grammar    *model.Grammar
		wantLinear model.GrammarLinearity
	}{
		{
			name:       "nil grammar",
			grammar:    nil,
			wantLinear: model.InvalidLinearity,
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
			},
			wantLinear: model.RightLinear,
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
			wantLinear: model.LeftLinear,
		},
		{
			name: "epsilon only grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			wantLinear: model.RightLinear,
		},
		{
			name: "single terminal production",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
				},
			},
			wantLinear: model.RightLinear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DetermineLinearity(tt.grammar)
			if tt.grammar != nil && tt.grammar.GrammarLinearity != tt.wantLinear {
				t.Errorf("DetermineLinearity() = %v, want %v", tt.grammar.GrammarLinearity, tt.wantLinear)
			}
		})
	}
}

func TestIsRightLinearProduction(t *testing.T) {
	g := &model.Grammar{
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S", "A"},
	}

	tests := []struct {
		name     string
		right    []model.Symbol
		expected bool
	}{
		{
			name:     "single terminal",
			right:    []model.Symbol{"a"},
			expected: true,
		},
		{
			name:     "terminal followed by non-terminal",
			right:    []model.Symbol{"a", "S"},
			expected: true,
		},
		{
			name:     "non-terminal at start",
			right:    []model.Symbol{"S", "a"},
			expected: false,
		},
		{
			name:     "epsilon",
			right:    []model.Symbol{model.Epsilon},
			expected: true,
		},
		{
			name:     "multiple terminals then non-terminal",
			right:    []model.Symbol{"a", "b", "S"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRightLinearProduction(tt.right, g)
			if got != tt.expected {
				t.Errorf("isRightLinearProduction() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsLeftLinearProduction(t *testing.T) {
	g := &model.Grammar{
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S", "A"},
	}

	tests := []struct {
		name     string
		right    []model.Symbol
		expected bool
	}{
		{
			name:     "single terminal",
			right:    []model.Symbol{"a"},
			expected: true,
		},
		{
			name:     "non-terminal followed by terminal",
			right:    []model.Symbol{"S", "a"},
			expected: true,
		},
		{
			name:     "terminal at start",
			right:    []model.Symbol{"a", "S"},
			expected: false,
		},
		{
			name:     "epsilon",
			right:    []model.Symbol{model.Epsilon},
			expected: true,
		},
		{
			name:     "non-terminal followed by multiple terminals",
			right:    []model.Symbol{"S", "a", "b"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLeftLinearProduction(tt.right, g)
			if got != tt.expected {
				t.Errorf("isLeftLinearProduction() = %v, want %v", got, tt.expected)
			}
		})
	}
}
