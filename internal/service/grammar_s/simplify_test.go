package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestSimplify(t *testing.T) {
	tests := []struct {
		name    string
		grammar *model.Grammar
		wantNil bool
	}{
		{
			name:    "nil grammar",
			grammar: nil,
			wantNil: true,
		},
		{
			name: "grammar with non-generating variables",
			grammar: &model.Grammar{
				StartSymbol: "S",
				Terminals:   []model.Symbol{"a"},
				NonTerminals: []model.Symbol{
					"S", "A", "B",
				},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A", "a"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"B"}},
					{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a"}},
				},
			},
			wantNil: false,
		},
		{
			name: "grammar with unreachable symbols",
			grammar: &model.Grammar{
				StartSymbol: "S",
				Terminals:   []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{
					"S", "A", "B",
				},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
					{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a"}},
				},
			},
			wantNil: false,
		},
		{
			name: "grammar with epsilon productions",
			grammar: &model.Grammar{
				StartSymbol: "S",
				Terminals:   []model.Symbol{"a"},
				NonTerminals: []model.Symbol{
					"S", "A",
				},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
				},
			},
			wantNil: false,
		},
		{
			name: "grammar with unit productions",
			grammar: &model.Grammar{
				StartSymbol: "S",
				Terminals:   []model.Symbol{"a"},
				NonTerminals: []model.Symbol{
					"S", "A", "B",
				},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"B"}},
					{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a"}},
				},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var grammar *model.Grammar
			if tt.grammar != nil {
				grammar = cloneGrammar(tt.grammar)
			}

			result := Simplify(grammar)

			if (result == nil) != tt.wantNil {
				t.Errorf("Simplify() = %v, wantNil %v", result == nil, tt.wantNil)
			}

			// 有效文法化简后应该仍然是有效文法
			if result != nil && tt.wantNil == false {
				err := GrammarCheckValidity(result)
				if err != nil {
					t.Errorf("Simplify() validation error = %v", err)
				}
			}
		})
	}
}

// TestSimplifyInvalidGrammar 测试无效文法的化简（只需要一个测试）
func TestSimplifyInvalidGrammar(t *testing.T) {
	// 空产生式的文法是无效文法
	grammar := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a"},
		NonTerminals: []model.Symbol{"S"},
		Productions:  []model.Production{},
	}

	result := Simplify(grammar)

	// 无效文法应返回自身，但验证会失败
	if result == nil {
		t.Error("Simplify() should return the grammar itself for invalid grammar")
	}

	err := GrammarCheckValidity(result)
	if err == nil {
		t.Error("Simplify() result should fail validation for invalid grammar")
	}
}

func TestRemoveNonGenerating(t *testing.T) {
	g := &model.Grammar{
		StartSymbol: "S",
		Terminals:   []model.Symbol{"a"},
		NonTerminals: []model.Symbol{
			"S", "A", "B",
		},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A", "a"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"B"}},
			{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a"}},
		},
	}

	RemoveNonGenerating(g)

	if len(g.NonTerminals) == 0 {
		t.Error("RemoveNonGenerating() should not remove all non-terminals")
	}

	for _, p := range g.Productions {
		if !containsSymbol(g.NonTerminals, p.Left[0]) {
			t.Errorf("Production left side %v not in NonTerminals", p.Left[0])
		}
	}
}

func TestRemoveUnreachable(t *testing.T) {
	g := &model.Grammar{
		StartSymbol: "S",
		Terminals:   []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{
			"S", "A", "B",
		},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
			{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a"}},
		},
	}

	RemoveUnreachable(g)

	if !containsSymbol(g.NonTerminals, g.StartSymbol) {
		t.Error("RemoveUnreachable() should keep start symbol")
	}
}

func TestRemoveEpsilonProductions(t *testing.T) {
	g := &model.Grammar{
		StartSymbol: "S",
		Terminals:   []model.Symbol{"a"},
		NonTerminals: []model.Symbol{
			"S", "A",
		},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
		},
	}

	RemoveEpsilonProductions(g)

	hasEpsilon := false
	for _, p := range g.Productions {
		if len(p.Right) == 1 && p.Right[0] == model.Epsilon {
			hasEpsilon = true
			t.Logf("Found epsilon production: %v", p)
		}
	}

	// S 是可空的,所以应该保留 S->ε
	if !hasEpsilon {
		t.Error("RemoveEpsilonProductions() should keep S->ε if S is nullable")
	}
}

func TestRemoveUnitProductions(t *testing.T) {
	g := &model.Grammar{
		StartSymbol: "S",
		Terminals:   []model.Symbol{"a"},
		NonTerminals: []model.Symbol{
			"S", "A", "B",
		},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"B"}},
			{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a"}},
		},
	}

	RemoveUnitProductions(g)

	for _, p := range g.Productions {
		if len(p.Right) == 1 && g.CheckIsNonTerminal(p.Right[0]) {
			t.Errorf("Unit production still exists: %v -> %v", p.Left, p.Right)
		}
	}
}

func TestSimplifyCompleteGrammar(t *testing.T) {
	g := &model.Grammar{
		StartSymbol: "S",
		Terminals:   []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{
			"S", "A", "B", "C",
		},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A", "B"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
			{Left: []model.Symbol{"B"}, Right: []model.Symbol{"b"}},
			{Left: []model.Symbol{"B"}, Right: []model.Symbol{model.Epsilon}},
			{Left: []model.Symbol{"C"}, Right: []model.Symbol{"a"}},
		},
	}

	originalLen := len(g.Productions)
	result := Simplify(g)

	if result == nil {
		t.Fatal("Simplify() should not return nil")
	}

	if len(result.Productions) > originalLen {
		t.Error("Simplify() should not increase number of productions")
	}

	err := GrammarCheckValidity(result)
	if err != nil {
		t.Errorf("Simplified grammar is invalid: %v", err)
	}
}

func TestSimplifyWithEpsilonInStart(t *testing.T) {
	g := &model.Grammar{
		StartSymbol: "S",
		Terminals:   []model.Symbol{"a"},
		NonTerminals: []model.Symbol{
			"S", "A",
		},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
		},
	}

	result := Simplify(g)

	if result == nil {
		t.Fatal("Simplify() should not return nil")
	}

	hasStartEpsilon := false
	for _, p := range result.Productions {
		if p.Left[0] == g.StartSymbol && len(p.Right) == 1 && p.Right[0] == model.Epsilon {
			hasStartEpsilon = true
			break
		}
	}

	if !hasStartEpsilon {
		t.Error("Simplify() should keep S->ε if S is nullable")
	}
}

func cloneGrammar(g *model.Grammar) *model.Grammar {
	if g == nil {
		return nil
	}

	clone := &model.Grammar{
		StartSymbol:  g.StartSymbol,
		Terminals:    make([]model.Symbol, len(g.Terminals)),
		NonTerminals: make([]model.Symbol, len(g.NonTerminals)),
		Productions:  make([]model.Production, len(g.Productions)),
		GrammarType:  g.GrammarType,
	}

	copy(clone.Terminals, g.Terminals)
	copy(clone.NonTerminals, g.NonTerminals)

	for i, p := range g.Productions {
		clone.Productions[i] = model.Production{
			Left:  make([]model.Symbol, len(p.Left)),
			Right: make([]model.Symbol, len(p.Right)),
		}
		copy(clone.Productions[i].Left, p.Left)
		copy(clone.Productions[i].Right, p.Right)
	}

	return clone
}
