// grammar_s/simplify.go
package grammar_s

import (
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/pkg/util"
)

// Simplify 对文法进行化简：去不可派生、不可达、空产生式、单一产生式
// 注意：会修改原文法结构，建议传入副本
func Simplify(grammar *model.Grammar) *model.Grammar {
	if grammar == nil || len(grammar.Productions) == 0 {
		return grammar
	}

	// 步骤1: 去除不可派生的变量（Non-generating variables）
	RemoveNonGenerating(grammar)
	util.PrintGrammar(grammar)

	// 步骤2: 去除不可达符号（Unreachable terminals/nonterminals）
	RemoveUnreachable(grammar)
	util.PrintGrammar(grammar)

	// 步骤3: 去除空产生式（ε-productions）
	RemoveEpsilonProductions(grammar)
	util.PrintGrammar(grammar)

	// 步骤4: 去除单一产生式（Unit productions）
	RemoveUnitProductions(grammar)
	util.PrintGrammar(grammar)

	return grammar
}

// removeNonGenerating 删除不可派生的变量（即不能推导出终结符串的变量）
func RemoveNonGenerating(g *model.Grammar) {
	var oldV, newV []model.Symbol

	// NEWV = { A | A → w ∈ P, w ∈ T* }
	for _, p := range g.Productions {
		if isAllTerminal(p.Right, g) {
			if !containsSymbol(newV, p.Left[0]) {
				newV = append(newV, p.Left[0])
			}
		}
	}

	for {
		oldV = copySlice(newV)
		// NEWV = OLDV ∪ { A | A → α ∈ P, α ∈ (T ∪ OLDV)* }
		for _, p := range g.Productions {
			if canDeriveUsingSet(p.Right, oldV, g) {
				if !containsSymbol(newV, p.Left[0]) {
					newV = append(newV, p.Left[0])
				}
			}
		}

		if slicesEqual(oldV, newV) {
			break
		}
	}

	// 过滤产生式：只保留左部 ∈ newV 且右部符号 ∈ (T ∪ newV)
	var filtered []model.Production
	for _, p := range g.Productions {
		if containsSymbol(newV, p.Left[0]) && allInSet(p.Right, newV, g.Terminals) {
			filtered = append(filtered, p)
		}
	}

	g.NonTerminals = newV
	g.Productions = filtered
}

// removeUnreachable 删除不可达的符号（从起始符号无法到达的变量和终结符）
func RemoveUnreachable(g *model.Grammar) {
	if len(g.Productions) == 0 {
		return
	}

	var oldV, newV []model.Symbol
	var oldT, newT []model.Symbol

	newV = append(newV, g.StartSymbol)

	// 初始：从 S 出发能直接到达的符号
	for _, p := range g.Productions {
		if p.Left[0] == g.StartSymbol {
			for _, sym := range p.Right {
				if g.CheckIsTerminal(sym) && !containsSymbol(newT, sym) {
					newT = append(newT, sym)
				} else if g.CheckIsNonTerminal(sym) && !containsSymbol(newV, sym) {
					newV = append(newV, sym)
				}
			}
		}
	}

	for {
		oldV = copySlice(newV)
		oldT = copySlice(newT)

		for _, p := range g.Productions {
			if containsSymbol(oldV, p.Left[0]) {
				for _, sym := range p.Right {
					if g.CheckIsTerminal(sym) && !containsSymbol(newT, sym) {
						newT = append(newT, sym)
					} else if g.CheckIsNonTerminal(sym) && !containsSymbol(newV, sym) {
						newV = append(newV, sym)
					}
				}
			}
		}

		if slicesEqual(oldV, newV) && slicesEqual(oldT, newT) {
			break
		}
	}

	// 过滤产生式
	var filtered []model.Production
	for _, p := range g.Productions {
		if containsSymbol(newV, p.Left[0]) && allInSet(p.Right, newV, newT) {
			filtered = append(filtered, p)
		}
	}

	g.NonTerminals = newV
	g.Terminals = newT
	g.Productions = filtered
}

// RemoveEpsilonProductions 去除空产生式（但若S可空，则保留S→ε）
func RemoveEpsilonProductions(g *model.Grammar) {
	U := getNullableVariables(g) // 可空变量集

	var newProductions []model.Production

	for _, p := range g.Productions {
		right := p.Right
		// 跳过空产生式（ε产生式）
		if len(right) == 1 && right[0] == model.Epsilon {
			continue
		}

		nullableIndices := []int{}
		for i, sym := range right {
			if containsSymbol(U, sym) {
				nullableIndices = append(nullableIndices, i)
			}
		}

		// 生成所有子集（除了全删）
		n := len(nullableIndices)
		for mask := 0; mask < (1 << n); mask++ {
			var newRight []model.Symbol
			skipAll := true
			for i, sym := range right {
				inNullable := -1
				for j, pos := range nullableIndices {
					if pos == i {
						inNullable = j
						break
					}
				}
				if inNullable == -1 || (mask&(1<<inNullable)) == 0 {
					newRight = append(newRight, sym)
					skipAll = false
				}
			}
			if skipAll {
				continue // 不允许全删
			}
			newProductions = append(newProductions, model.Production{
				Left:  []model.Symbol{p.Left[0]},
				Right: newRight,
			})
		}
	}

	// 如果 S ∈ U，且没有 S→ε，添加
	if containsSymbol(U, g.StartSymbol) {
		hasEpsilon := false
		for _, p := range newProductions {
			if p.Left[0] == g.StartSymbol && len(p.Right) == 1 && p.Right[0] == model.Epsilon {
				hasEpsilon = true
				break
			}
		}
		if !hasEpsilon {
			newProductions = append(newProductions, model.Production{
				Left:  []model.Symbol{g.StartSymbol},
				Right: []model.Symbol{model.Epsilon},
			})
		}
	}

	g.Productions = newProductions
}

// getNullableVariables 求可空变量集 U
func getNullableVariables(g *model.Grammar) []model.Symbol {
	var oldU, newU []model.Symbol

	// A → ε
	for _, p := range g.Productions {
		if len(p.Right) == 1 && p.Right[0] == model.Epsilon {
			if !containsSymbol(newU, p.Left[0]) {
				newU = append(newU, p.Left[0])
			}
		}
	}

	for {
		oldU = copySlice(newU)
		for _, p := range g.Productions {
			// 跳过空产生式
			if len(p.Right) == 1 && p.Right[0] == model.Epsilon {
				continue
			}
			if allInSet(p.Right, oldU, nil) {
				if !containsSymbol(newU, p.Left[0]) {
					newU = append(newU, p.Left[0])
				}
			}
		}
		if slicesEqual(oldU, newU) {
			break
		}
	}

	return newU
}

// removeUnitProductions 去除单一产生式 A → B
func RemoveUnitProductions(g *model.Grammar) {
	unitClosure := make(map[model.Symbol][]model.Symbol)
	nonTerminals := g.NonTerminals

	// 初始化：Ai ⇒* Ai
	for _, A := range nonTerminals {
		unitClosure[A] = []model.Symbol{A}
	}

	// 迭代扩展闭包
	for {
		changed := false
		for _, p := range g.Productions {
			// A → B 是单一产生式
			if len(p.Right) == 1 && g.CheckIsNonTerminal(p.Right[0]) {
				A := p.Left[0]
				B := p.Right[0]
				for _, C := range unitClosure[B] {
					if !containsSymbol(unitClosure[A], C) {
						unitClosure[A] = append(unitClosure[A], C)
						changed = true
					}
				}
			}
		}
		if !changed {
			break
		}
	}

	// 构建新产生式：Ai → α，其中 B → α 是非单一产生式，B ∈ closure[Ai]
	var newProductions []model.Production
	for _, A := range nonTerminals {
		for _, B := range unitClosure[A] {
			for _, p := range g.Productions {
				if p.Left[0] == B {
					// 非单一产生式：右部长度 ≠ 1 或不是非终结符
					if len(p.Right) != 1 || !g.CheckIsNonTerminal(p.Right[0]) {
						newProductions = append(newProductions, model.Production{
							Left:  []model.Symbol{A},
							Right: p.Right,
						})
					}
				}
			}
		}
	}

	g.Productions = newProductions
}

func containsSymbol(slice []model.Symbol, item model.Symbol) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

func copySlice(slice []model.Symbol) []model.Symbol {
	cpy := make([]model.Symbol, len(slice))
	copy(cpy, slice)

	return cpy
}

func slicesEqual(a, b []model.Symbol) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func isAllTerminal(symbols []model.Symbol, g *model.Grammar) bool {
	for _, s := range symbols {
		if !g.CheckIsTerminal(s) && s != model.Epsilon {
			return false
		}
	}

	return true
}

func canDeriveUsingSet(symbols []model.Symbol, vars []model.Symbol, g *model.Grammar) bool {
	for _, s := range symbols {
		if g.CheckIsNonTerminal(s) {
			if !containsSymbol(vars, s) {
				return false
			}
		} else if s != model.Epsilon && !g.CheckIsTerminal(s) {
			return false
		}
	}

	return true
}

func allInSet(symbols []model.Symbol, vars []model.Symbol, terms []model.Symbol) bool {
	var termSet []model.Symbol
	if terms != nil {
		termSet = terms
	} else {
		termSet = []model.Symbol{}
	}

	for _, s := range symbols {
		if g := containsSymbol(termSet, s); g || s == model.Epsilon {
			continue
		}
		if containsSymbol(vars, s) {
			continue
		}

		return false
	}

	return true
}
