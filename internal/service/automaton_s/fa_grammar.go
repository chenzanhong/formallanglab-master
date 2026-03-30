package automaton_s

import (
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// 自动机转正则文法，默认自动机有效，且转为右线性文法
// 标准转换方法（以右线性文法为例）
// 给定一个 NFA M = (Q, Σ, δ, q₀, F)，构造右线性文法 G = (V, Σ, P, S)：
// 1. 非终结符集合 V = Q（每个状态是一个非终结符）
// 2. 终结符集合 Σ 不变
// 3. 开始符号 S = q₀
// 4. 产生式规则 P：
//   - 对每条转移 δ(qᵢ, a) ∋ qⱼ（即 qᵢ →a→ qⱼ），添加产生式：qᵢ → a qⱼ
//   - 对每个终态 q_f ∈ F，添加：q_f → ε
//
// 对于 ε 转移
// 方法一：先消除 ε 转移，再转文法（将 ε-NFA → NFA（无 ε）→ 正则文法）得到干净、标准的正则文法
// 方法二：直接转文法，再消除单位产生式（ε 转移 → 单位产生式 → 后处理消除）
func FAToGrammar(automaton *model.Automaton) *model.Grammar {
	if automaton == nil {
		return &model.Grammar{}
	}
	var g model.Grammar
	// 非终结符集合 V = Q（每个状态是一个非终结符）
	g.NonTerminals = make([]model.Symbol, len(automaton.States))
	for i, s := range automaton.States {
		g.NonTerminals[i] = model.Symbol(s)
	}

	// 终结符集合 Σ 不变
	g.Terminals = automaton.Alphabet

	// 开始符号 S = q₀
	g.StartSymbol = model.Symbol(automaton.InitialState)

	// 产生式规则 P，采用方法二，不消除单一产生式，用户前端可以手动调用自动机简化接口
	g.Productions = make([]model.Production, 0, len(automaton.Transitions))
	for _, t := range automaton.Transitions {
		for _, toState := range t.ToStates { // 遍历所有目标，允许多个转移状态（NFA）
			var right []model.Symbol
			if t.Input == model.Epsilon {
				right = []model.Symbol{model.Symbol(toState)} // A → B
			} else {
				right = []model.Symbol{model.Symbol(t.Input), model.Symbol(toState)} // A → a B
			}
			g.Productions = append(g.Productions, model.Production{
				Left:  []model.Symbol{model.Symbol(t.FromState)},
				Right: right,
			})
		}
	}

	// 为每个接受状态添加 ε 产生式（如果未存在）
	for _, state := range automaton.AcceptingStates {
		hasEpsilon := false
		symState := model.Symbol(state)
		for _, p := range g.Productions {
			if p.Left[0] == symState && p.Right[0] == model.Epsilon {
				hasEpsilon = true
				break
			}
		}
		if !hasEpsilon {
			g.Productions = append(g.Productions, model.Production{
				Left:  []model.Symbol{symState},
				Right: []model.Symbol{model.Epsilon},
			})
		}
	}

	return &g
}
