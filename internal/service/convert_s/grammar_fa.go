package convert_s

import (
	"backend/internal/domain/model"
	"strconv"
)

// 正则文法转自动机，需要区分左线性和右线性
// 默认已经检查过类型检查，是正则文法
func RegularGrammarToFA(g *model.Grammar, isRightLinear bool) *model.Automaton {
	// 检查文法为左线性还是右线性
	if isRightLinear {
		return RightLinearGrammarToFA(g)
	}
	return LeftLinearGrammarToFA(g)
}

// 右线性文法转自动机
// 算法思路：
// 1. 右线性文法的产生式形式为 A → aB 或 A → a（A,B为非终结符，a为终结符串）
// 2. 将每个非终结符映射为自动机的一个状态
// 3. 文法的开始符号对应自动机的初始状态
// 4. 对于形如 A → aB 的产生式，创建一条从状态A到状态B的转移边，标记为a
// 5. 对于形如 A → a 的产生式，创建一条从状态A到接受状态的转移边，标记为a
// 6. 如果有多个直接产生终结符串的产生式，可能需要创建一个额外的接受状态
func RightLinearGrammarToFA(g *model.Grammar) *model.Automaton {
	var automaton model.Automaton

	// 符号集
	automaton.Alphabet = append(automaton.Alphabet, g.Terminals...)

	// 状态集
	for _, sym := range g.NonTerminals {
		automaton.States = append(automaton.States, model.State(sym))
	}

	// 初始状态
	automaton.InitialState = model.State(g.StartSymbol)

	// 接受状态，只有一个
	automaton.AcceptingStates = make([]model.State, 0)
	automaton.AcceptingStates = append(automaton.AcceptingStates, model.State("accept"))

	// 记录遇到多符号产生式右部时，需要添加的中间状态的序号
	// 中间状态索引，形如 A -> ab……B 时，需要添加中间状态
	intermediateStateIndex := -1
	getNextIntermediateStateIndex := func() int {
		intermediateStateIndex++
		return intermediateStateIndex
	}

	// 状态转移
	automaton.Transitions = make([]model.Transition, 0)
	for _, production := range g.Productions {
		var transition model.Transition
		// 左部
		transition.FromState = model.State(production.Left[0])
		// 右部
		if len(production.Right) == 1 { // 单一符号产生式右部
			// 如果是终结符，直接转移到接受状态
			if g.CheckIsTerminal(production.Right[0]) {
				transition.Input = production.Right[0] // 终结符
				transition.ToStates = append(transition.ToStates, automaton.AcceptingStates[0])
				automaton.Transitions = append(automaton.Transitions, transition)
			} else {
				// 如果是非终结符，空转移，转移到下一个状态
				transition.Input = model.Epsilon // 空转移
				transition.ToStates = append(transition.ToStates, model.State(production.Right[0]))
				automaton.Transitions = append(automaton.Transitions, transition)
			}
		} else if len(production.Right) == 2 { // 形如 A → aB 的产生式
			transition.Input = production.Right[0] // 终结符
			transition.ToStates = append(transition.ToStates, model.State(production.Right[1]))
			automaton.Transitions = append(automaton.Transitions, transition)
		} else { // 多符号产生式右部，形如 A → aa……B，需添加中间状态
			// 第 0 个符号（终结符）
			transition.Input = production.Right[0]
			newState := model.State(transition.FromState + model.State("_") + model.State(strconv.Itoa(getNextIntermediateStateIndex()))) // 中间状态
			transition.ToStates = append(transition.ToStates, newState)
			newTransition := model.Transition{ // 新创建
				FromState: newState,
				Input:     model.Epsilon,
				ToStates:  []model.State{newState},
			}
			automaton.States = append(automaton.States, newState)
			automaton.Transitions = append(automaton.Transitions, newTransition)
			// 前第 1 到 n-3 个符号（终结符）
			for _, sym := range production.Right[1 : len(production.Right)-2] {
				var newTransition model.Transition // 新创建
				newTransition.FromState = newState
				newTransition.Input = sym
				newState = model.State(transition.FromState + model.State("_") + model.State(strconv.Itoa(getNextIntermediateStateIndex())))
				newTransition.ToStates = append(newTransition.ToStates, newState)
				automaton.States = append(automaton.States, newState)
				automaton.Transitions = append(automaton.Transitions, newTransition)
			}
			// 最后的形如 aB 的部分
			automaton.Transitions = append(automaton.Transitions, model.Transition{
				FromState: newState,
				Input:     production.Right[len(production.Right)-2],
				ToStates:  []model.State{model.State(production.Right[len(production.Right)-1])},
			})
		}
	}

	return &automaton
}

// 左线性文法转自动机
// 算法思路（两种正确方法）：

// ✅ 方法一：通过语言反转（间接法）
// 1. 反转语言：左线性文法G生成语言L，其反转语言L^R可由右线性文法G^R生成（将每条产生式右侧字符串反转）
// 2. 转右线性为NFA：用标准方法将G^R转换为NFA M'，它识别L^R
// 3. 反转自动机：将M'的起始状态与所有终态互换，所有转移边方向反转，得到新NFA M，它识别原语言L
// 核心思想："左线性 = 右线性描述的语言的反转"

// ✅ 方法二：直接构造（反向建图）
// 1. 状态设计：每个非终结符作为一个状态；新增一个初始状态q0
// 2. 转移规则：
//   - 对产生式A → Ba：添加转移 B →a→ A（注意方向是从B到A）
//   - 对产生式A → a：添加转移 q0 →a→ A
//   - 若开始符号S → ε，则q0也是终态
//
// 3. 终态设定：开始符号S对应的状态设为唯一终态（或根据ε产生式调整）
// 核心思想：自动机从左读输入，而左线性文法从右生成字符串，因此转移方向要"反过来"建模
func LeftLinearGrammarToFA(g *model.Grammar) *model.Automaton {
	return leftLinearGrammarToFAByBuild(g) // 默认使用直接构造法
	// return LeftLinearGrammarToFAByReverse(g) // 默认使用间接法
}

// 左线性文法转自动机：通过语言反转（间接法）
// 1. 反转语言：左线性文法G生成语言L，其反转语言L^R可由右线性文法G^R生成（将每条产生式右侧字符串反转）
// 2. 转右线性为NFA：用标准方法将G^R转换为NFA M'，它识别L^R
// 3. 反转自动机：将M'的起始状态与所有终态互换，所有转移边方向反转，得到新NFA M，它识别原语言L
func leftLinearGrammarToFAByReverse(g *model.Grammar) *model.Automaton {
	// 实现左线性文法转自动机的逻辑
	return nil
}

// 左线性文法转自动机：直接构造（反向建图）
// 1. 状态设计：每个非终结符作为一个状态；新增一个初始状态q0
// 2. 转移规则：
//   - 对产生式A → Ba：添加转移 B →a→ A（注意方向是从B到A）
//   - 对产生式A → a：添加转移 q0 →a→ A
//   - 若开始符号S → ε，则q0也是终态
//
// 3. 终态设定：开始符号S对应的状态设为唯一终态（或根据ε产生式调整）
func leftLinearGrammarToFAByBuild(g *model.Grammar) *model.Automaton {
	var a model.Automaton
	// 1. 状态：每个非终结符作为一个状态；新增一个初始状态q0
	a.States = append(a.States, model.State("initial"))
	a.InitialState = model.State("initial")
	for _, nonTerm := range g.NonTerminals {
		a.States = append(a.States, model.State(nonTerm))
	}

	// 2. 符号
	for _, term := range g.Terminals {
		a.Alphabet = append(a.Alphabet, model.Symbol(term))
	}

	// 中间状态索引，形如 A -> ab……B 时，需要添加中间状态
	intermediateStateIndex := -1
	getNextIntermediateStateIndex := func() int {
		intermediateStateIndex++
		return intermediateStateIndex
	}

	// 3. 转移规则
	for _, production := range g.Productions {
		// 对产生式A → Ba：添加转移 B →a→ A（注意方向是从B到A）
		if g.CheckIsNonTerminal(production.Right[0]) { // 形如A -> Bw，w =a1a2a3...an，需展开为 B →a1→ A1 →a2→ A2 →...→ An-1 →an→ A
			wLength := len(production.Right) - 1
			switch wLength {
			case 0: // 形如A -> B
				transition := model.Transition{
					FromState: model.State(production.Right[0]),
					Input:     model.Epsilon,
					ToStates:  []model.State{model.State(production.Left[0])},
				}
				a.Transitions = append(a.Transitions, transition)
			case 1: // 形如A -> Ba
				transition := model.Transition{
					FromState: model.State(production.Right[0]),
					Input:     model.Symbol(production.Right[1]),
					ToStates:  []model.State{model.State(production.Left[0])},
				}
				a.Transitions = append(a.Transitions, transition)
			default: // 形如 A -> Baa……
				// 新增中间状态
				// 新增转移 B →a1→ A1
				newTransition := model.Transition{
					FromState: model.State(production.Right[0]),
					Input:     model.Symbol(production.Right[1]),
				}
				newState := model.State(string(production.Left[0]) + "_" + strconv.Itoa(getNextIntermediateStateIndex()))
				newTransition.ToStates = []model.State{newState}
				a.States = append(a.States, newState)
				a.Transitions = append(a.Transitions, newTransition)
				// 新增转移 A1 →a2→ A2 →...→ An-1
				for i := 2; i < wLength; i++ {
					newTransition := model.Transition{
						FromState: newState,
						Input:     model.Symbol(production.Right[i]),
					}
					newState = model.State(string(production.Left[0]) + "_" + strconv.Itoa(getNextIntermediateStateIndex()))
					newTransition.ToStates = []model.State{newState}
					a.States = append(a.States, newState)
					a.Transitions = append(a.Transitions, newTransition)
				}
				// 新增转移 An-1 →an→ A
				newTransition = model.Transition{
					FromState: newState,
					Input:     model.Symbol(production.Right[wLength]),
					ToStates:  []model.State{model.State(production.Left[0])},
				}
				a.Transitions = append(a.Transitions, newTransition)
			}
		} else { // 形如A -> w
			wLength := len(production.Right)
			if wLength == 1 { // 形如A -> ε 或 A -> a, 新增转移 q0 →a→ A
				transition := model.Transition{
					FromState: a.InitialState,
					Input:     model.Symbol(production.Right[0]),
					ToStates:  []model.State{model.State(production.Left[0])},
				}
				a.Transitions = append(a.Transitions, transition)
			} else { // 形如A -> w，w =a1a2a3...an，需展开为 q0 →a1→ A1 →a2→ A2 →...→ An-1 →an→ A
				// 新增转移 q0 →a1→ A1
				newTransition := model.Transition{
					FromState: a.InitialState,
					Input:     model.Symbol(production.Right[0]),
				}
				newState := model.State(string(production.Left[0]) + "_" + strconv.Itoa(getNextIntermediateStateIndex()))
				newTransition.ToStates = []model.State{newState}
				a.States = append(a.States, newState)
				a.Transitions = append(a.Transitions, newTransition)
				// 新增转移 A1 →a2→ A2 →...→ An-1
				for i := 1; i < wLength-1; i++ {
					newTransition := model.Transition{ // 新创建
						FromState: newState,
						Input:     model.Symbol(production.Right[i]),
					}
					newState = model.State(string(production.Left[0]) + "_" + strconv.Itoa(getNextIntermediateStateIndex()))
					newTransition.ToStates = append(newTransition.ToStates, newState)
					a.States = append(a.States, newState)
					a.Transitions = append(a.Transitions, newTransition)
				}
				// 新增转移 An-1 →an→ A
				newTransition = model.Transition{
					FromState: newState,
					Input:     model.Symbol(production.Right[wLength-1]),
					ToStates:  []model.State{model.State(production.Left[0])},
				}
				a.Transitions = append(a.Transitions, newTransition)
			}
		}
	}

	return &a
}
