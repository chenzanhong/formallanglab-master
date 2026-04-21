package grammar_s

import (
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

func TestDeduplicateProductions(t *testing.T) {
	tests := []struct {
		name     string
		input    []model.Production
		expected int // 期望的去重后产生式数量
	}{
		{
			name: "无重复产生式",
			input: []model.Production{
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b"}},
				{Left: []model.Symbol{"A"}, Right: []model.Symbol{"c"}},
			},
			expected: 3,
		},
		{
			name: "有重复产生式",
			input: []model.Production{
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}}, // 重复
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b"}},
				{Left: []model.Symbol{"A"}, Right: []model.Symbol{"c"}},
				{Left: []model.Symbol{"A"}, Right: []model.Symbol{"c"}}, // 重复
			},
			expected: 3,
		},
		{
			name: "多个重复产生式",
			input: []model.Production{
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}}, // 重复
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}}, // 重复
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b"}},
			},
			expected: 2,
		},
		{
			name: "空产生式",
			input: []model.Production{
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}}, // 重复
				{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
			},
			expected: 2,
		},
		{
			name: "空输入",
			input: []model.Production{},
			expected: 0,
		},
		{
			name: "左部相同但右部不同",
			input: []model.Production{
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b"}},
				{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
			},
			expected: 3, // 这些都不同，不应去重
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicateProductions(tt.input)
			if len(result) != tt.expected {
				t.Errorf("期望产生式数量为 %d, 实际为 %d", tt.expected, len(result))
			}
		})
	}
}

func TestProductionToKey(t *testing.T) {
	tests := []struct {
		name     string
		prod1    model.Production
		prod2    model.Production
		expected bool // 是否应该生成相同的 key
	}{
		{
			name:     "相同产生式",
			prod1:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
			prod2:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
			expected: true,
		},
		{
			name:     "不同左部",
			prod1:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
			prod2:    model.Production{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
			expected: false,
		},
		{
			name:     "不同右部",
			prod1:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
			prod2:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{"b"}},
			expected: false,
		},
		{
			name:     "相同空产生式",
			prod1:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
			prod2:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{model.Epsilon}},
			expected: true,
		},
		{
			name:     "右部长度不同",
			prod1:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
			prod2:    model.Production{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a", "A"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key1 := productionToKey(tt.prod1)
			key2 := productionToKey(tt.prod2)

			if tt.expected && key1 != key2 {
				t.Errorf("期望 key 相同，但得到 %q 和 %q", key1, key2)
			}
			if !tt.expected && key1 == key2 {
				t.Errorf("期望 key 不同，但都得到 %q", key1)
			}
		})
	}
}

func TestNormalizeGrammar(t *testing.T) {
	tests := []struct {
		name          string
		input         model.Grammar
		expectedCount int
	}{
		{
			name: "ε 符号规范化",
			input: model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"epsilon"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"ε"}}, // 重复
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"a"}},
				},
			},
			expectedCount: 2, // 两个 ε 产生式应该被去重
		},
		{
			name: "混合重复和λ符号",
			input: model.Grammar{
				StartSymbol:  "S",
				Terminals:    []model.Symbol{"a", "b"},
				NonTerminals: []model.Symbol{"S", "A"},
				Productions: []model.Production{
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"λ"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}},
					{Left: []model.Symbol{"S"}, Right: []model.Symbol{"a"}}, // 重复
					{Left: []model.Symbol{"A"}, Right: []model.Symbol{"b"}},
				},
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NormalizeGrammar(&tt.input)
			if len(tt.input.Productions) != tt.expectedCount {
				t.Errorf("期望产生式数量为 %d, 实际为 %d", tt.expectedCount, len(tt.input.Productions))
			}
		})
	}
}
