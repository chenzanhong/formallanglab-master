package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestRecognizeString(t *testing.T) {
	tests := []struct {
		name       string
		grammar    *model.Grammar
		input      []model.Symbol
		maxSteps   int
		maxWidth   int
		wantAccept bool
	}{
		{
			name: "empty string with epsilon grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			input:      []model.Symbol{},
			maxSteps:   100,
			maxWidth:   100,
			wantAccept: true,
		},
		{
			name: "simple string recognition",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "b"}},
				},
			},
			input:      []model.Symbol{"a", "b"},
			maxSteps:   100,
			maxWidth:   100,
			wantAccept: true,
		},
		{
			name: "reject non-derivable string",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			input:      []model.Symbol{"b"},
			maxSteps:   100,
			maxWidth:   100,
			wantAccept: false,
		},
		{
			name: "recursive grammar",
			grammar: &model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a"},
				NonTerminals: []model.Symbol{"S"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "S"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				},
			},
			input:      []model.Symbol{"a", "a", "a"},
			maxSteps:   100,
			maxWidth:   100,
			wantAccept: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accepted, _, err := RecognizeString(tt.grammar, tt.input, tt.maxSteps, tt.maxWidth)
			if err != nil {
				t.Errorf("RecognizeString() error = %v", err)
				return
			}
			if accepted != tt.wantAccept {
				t.Errorf("RecognizeString() = %v, want %v", accepted, tt.wantAccept)
			}
		})
	}
}

func TestCalculateFirst(t *testing.T) {
	g := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S", "A"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{model.Epsilon}},
		},
	}

	firstSet := CalculateFirst(g)

	if firstSet == nil {
		t.Fatal("CalculateFirst() returned nil")
	}

	if _, ok := firstSet["a"]; !ok {
		t.Error("CalculateFirst() should include terminal 'a'")
	}

	if _, ok := firstSet["S"]; !ok {
		t.Error("CalculateFirst() should include non-terminal 'S'")
	}
}

func TestCalculateFollow(t *testing.T) {
	g := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S", "A"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
			{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
		},
	}

	firstSet := CalculateFirst(g)
	followSet := CalculateFollow(g, firstSet)

	if followSet == nil {
		t.Fatal("CalculateFollow() returned nil")
	}

	if _, ok := followSet["S"]; !ok {
		t.Error("CalculateFollow() should include non-terminal 'S'")
	}

	if _, ok := followSet["#"]; !ok {
		t.Error("Follow set of start symbol should contain '#'")
	}
}

func TestBFSParseDetailed(t *testing.T) {
	g := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "b"}},
		},
	}

	result := BFSParseDetailed(g, []model.Symbol{"a", "b"}, 100, 100)

	if result == nil {
		t.Fatal("BFSParseDetailed() returned nil")
	}

	if !result.Accepted {
		t.Error("BFSParseDetailed() should accept 'ab'")
	}
}

func TestParseStringWithMode(t *testing.T) {
	g := &model.Grammar{
		StartSymbol:  "S",
		Terminals:    []model.Symbol{"a", "b"},
		NonTerminals: []model.Symbol{"S"},
		Productions: []model.Production{
			{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "b"}},
		},
	}

	tests := []struct {
		name       string
		input      string
		mode       ParseMode
		wantAccept bool
	}{
		{
			name:       "accept valid string with BFS mode",
			input:      "ab",
			mode:       BFSMode,
			wantAccept: true,
		},
		{
			name:       "reject invalid string with BFS mode",
			input:      "ba",
			mode:       BFSMode,
			wantAccept: false,
		},
		{
			name:       "accept valid string with Auto mode",
			input:      "ab",
			mode:       AutoMode,
			wantAccept: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseStringWithMode(g, tt.input, tt.mode, false)
			if result.Accepted != tt.wantAccept {
				t.Errorf("ParseStringWithMode() accepted = %v, want %v", result.Accepted, tt.wantAccept)
			}
		})
	}
}
