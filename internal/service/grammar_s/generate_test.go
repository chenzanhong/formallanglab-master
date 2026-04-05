package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestGrammarGenerateExampleString(t *testing.T) {
	tests := []struct {
		name          string
		grammar       *model.Grammar
		wantAcceptLen int
		wantRejectLen int
	}{
		{
			name: "recursive CFG with epsilon",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"S", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S", "b"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b", "S", "a"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.ContextFreeGrammar,
			},
			wantAcceptLen: 1,
			wantRejectLen: 1,
		},
		{
			name: "simple regular grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b", "S"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
				},
				GrammarType: model.RegularGrammar,
			},
			wantAcceptLen: 1,
			wantRejectLen: 1,
		},
		{
			name: "CFG with epsilon production",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"A"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a", "A", "b"}},
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
				},
				GrammarType: model.ContextFreeGrammar,
			},
			wantAcceptLen: 1,
			wantRejectLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accept, reject := GrammarGenerateExampleString(tt.grammar)
			if len(accept) < tt.wantAcceptLen {
				t.Errorf("GrammarGenerateExampleString() accept = %v, want at least %d", accept, tt.wantAcceptLen)
			}
			if len(reject) < tt.wantRejectLen {
				t.Errorf("GrammarGenerateExampleString() reject = %v, want at least %d", reject, tt.wantRejectLen)
			}
		})
	}
}

func TestGenerateAcceptExampleString(t *testing.T) {
	g := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a"},
		NonTerminals: []model.Symbol{"S"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
		},
		GrammarType: model.ContextFreeGrammar,
	}

	accept := generateAcceptExampleString(g)

	if len(accept) == 0 {
		t.Error("generateAcceptExampleString() should return at least one string")
	}

	// Check that empty string (epsilon) is in accept
	hasEpsilon := false
	for _, s := range accept {
		if s == "ε" || s == "" {
			hasEpsilon = true
			break
		}
	}
	if !hasEpsilon {
		t.Error("generateAcceptExampleString() should include epsilon for this grammar")
	}
}

func TestGenerateRejectExampleString(t *testing.T) {
	g := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
		},
		GrammarType: model.ContextFreeGrammar,
	}

	seenAccept := make(map[string]bool)
	seenAccept["a"] = true
	seenAccept["aa"] = true
	seenAccept[""] = true

	reject := generateRejectExampleString(g, seenAccept)

	if len(reject) == 0 {
		t.Error("generateRejectExampleString() should return at least one string")
	}

	// Check that reject strings are not in accept set
	for _, s := range reject {
		if seenAccept[s] {
			t.Errorf("generateRejectExampleString() returned accept string: %s", s)
		}
	}
}
