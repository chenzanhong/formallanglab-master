package grammar_s

import (
	"fmt"

	"github.com/chenzanhong/zlog"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/pkg/util"
)

// RegularGrammarToFA 正则文法转自动机，需要区分左线性和右线性
// 默认已经检查过类型检查，是正则文法
func RegularGrammarToFA(g *model.Grammar) (*model.Automaton, error) {
	if g == nil {
		return nil, fmt.Errorf("文法为空")
	}

	DetermineLinearity(g) // 确定线性

	// 检查文法为左线性还是右线性
	switch g.GrammarLinearity {
	case model.RightLinear:
		return rightLinearGrammarToFA(g), nil
	case model.LeftLinear:
		return leftLinearGrammarToFA(g), nil
	}

	return nil, fmt.Errorf("非线性文法，暂不支持转为有限自动机")
}

func RegularGrammarToFAWithProcess(g *model.Grammar, grammarLinearity model.GrammarLinearity) *model.GrammarToFAProcess {
	if g == nil {
		return nil
	}
	// 检查文法为左线性还是右线性
	switch grammarLinearity {
	case model.RightLinear:
		return rightLinearGrammarToFAWithProcess(g)
	case model.LeftLinear:
		return leftLinearGrammarToFAWithProcess(g)
	}

	return nil
}

// rightLinearGrammarToFA 右线性文法转自动机
// 算法思路：
// 1. 右线性文法的产生式形式为 A → aB 或 A → a（A,B 为非终结符，a 为终结符串）
// 2. 将每个非终结符映射为自动机的一个状态
// 3. 文法的开始符号对应自动机的初始状态
// 4. 对于形如 A → aB 的产生式，创建一条从状态 A 到状态 B 的转移边，标记为 a
// 5. 对于形如 A → a 的产生式，创建一条从状态 A 到接受状态的转移边，标记为 a，如果 a 为 model.Epsilon，也可以单纯只把 A 纳入接受态
// 6. 如果有多个直接产生终结符串的产生式，可能需要创建一个额外的接受状态
func rightLinearGrammarToFA(g *model.Grammar) *model.Automaton {
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

	// 记录遇到多符号产生式右部时，需要添加的中间状态的序号
	// 中间状态索引，形如 A -> ab……B 时，需要添加中间状态
	intermediateStateIndex := -1
	getNextIntermediateStateIndex := func(base model.State) model.State {
		intermediateStateIndex++
		return model.State(fmt.Sprintf("%s_%d", base, intermediateStateIndex))
	}

	// 状态转移
	automaton.Transitions = make([]model.Transition, 0)
	hasModelAccept := false
	for _, production := range g.Productions {
		var transition model.Transition
		// 左部
		transition.FromState = model.State(production.Left[0])
		// 右部
		if len(production.Right) == 1 { // 单一符号产生式右部
			// 如果是终结符，直接转移到接受状态
			if g.CheckIsTerminal(production.Right[0]) {
				if !hasModelAccept {
					hasModelAccept = true
				}
				transition.Input = production.Right[0] // 终结符
				transition.ToStates = append(transition.ToStates, model.UniqueFinalState)
				automaton.Transitions = append(automaton.Transitions, transition)
			} else if production.Right[0] == model.Epsilon {
				// 空转移，加入接受态
				automaton.AcceptingStates = append(automaton.AcceptingStates, model.State(production.Left[0]))
			} else {
				// 如果是非终结符，转移到下一个状态
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
			newState := getNextIntermediateStateIndex(transition.FromState) // 中间状态
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
				newState = getNextIntermediateStateIndex(transition.FromState)
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

	if hasModelAccept {
		automaton.States = append(automaton.States, model.UniqueFinalState)
		automaton.AcceptingStates = append(automaton.AcceptingStates, model.UniqueFinalState)
	}

	return &automaton
}

func rightLinearGrammarToFAWithProcess(g *model.Grammar) *model.GrammarToFAProcess {
	var automaton model.Automaton
	process := &model.GrammarToFAProcess{
		Linear:          model.RightLinear,
		OriginalGrammar: *g,
		Steps:           make([]model.GrammarToFAStep, 0),
	}

	// 初始化基础结构
	automaton.Alphabet = append(automaton.Alphabet, g.Terminals...)
	for _, sym := range g.NonTerminals {
		automaton.States = append(automaton.States, model.State(sym))
	}
	automaton.InitialState = model.State(g.StartSymbol)

	process.Steps = append(process.Steps, model.GrammarToFAStep{
		Action:      "init",
		Description: "初始化自动机：设置字母表、状态机、初始状态",
		NewStates:   util.SymbolsToStates(g.NonTerminals),
	})

	// process.Steps = append(process.Steps, model.GrammarToFAStep{
	// 	Action:      "create_accept_state",
	// 	Description: "创建唯一接受状态‘accept’",
	// 	NewStates:   []model.State{model.UniqueFinalState},
	// })

	// 中间状态计数器
	intermediateStateIndex := 0
	getNextIntermediateStateIndex := func(from model.State) model.State {
		intermediateStateIndex++
		return model.State(fmt.Sprintf("%s_%d", from, intermediateStateIndex))
	}

	hasModelAccept := false
	newStateFmt := fmt.Sprintf("，新增加接受状态：%s", model.UniqueFinalState)
	// 遍历所有产生式
	for _, prod := range g.Productions {
		step := model.GrammarToFAStep{
			Production: &prod,
			Action:     "handle_production",
		}

		leftSym := prod.Left[0]
		fromState := model.State(leftSym)

		rightLen := len(prod.Right)
		switch {
		// 1) 形如 A -> ε，将A加入接收状态
		case rightLen == 1 && prod.Right[0] == model.Epsilon:
			// 将A加入接收状态
			automaton.AcceptingStates = append(automaton.AcceptingStates, fromState)
			step.Description = fmt.Sprintf("因产生式 %s → ε，将状态 %s 加入接受状态集合", leftSym, fromState)

		// 2) 形如 A -> a，新增转移 A →a→ accept
		case rightLen == 1 && g.CheckIsTerminal(prod.Right[0]):
			input := prod.Right[0]
			toState := model.UniqueFinalState
			trans := model.Transition{
				FromState: fromState,
				Input:     input,
				ToStates:  []model.State{toState},
			}
			automaton.Transitions = append(automaton.Transitions, trans)
			step.Description = fmt.Sprintf("因产生式 %s → %s，新增转移 %s → %s → %s", leftSym, input, fromState, input, toState)
			if !hasModelAccept {
				hasModelAccept = true
				step.Description += newStateFmt
			}

		// 3) 形如 A -> B，新增转移 A →ε→ B
		case rightLen == 1 && g.CheckIsNonTerminal(prod.Right[0]):
			toState := model.State(prod.Right[0])
			trans := model.Transition{
				FromState: fromState,
				Input:     model.Epsilon,
				ToStates:  []model.State{toState},
			}
			automaton.Transitions = append(automaton.Transitions, trans)
			step.NewTransitions = []model.Transition{trans}
			step.Description = fmt.Sprintf("添加 ε-转移: %s --ε--> %s", fromState, toState)

		// 4) 形如 A -> aB，新增转移 A →a→ B
		case rightLen == 2 && g.CheckIsTerminal(prod.Right[0]) && g.CheckIsNonTerminal(prod.Right[1]):
			input := prod.Right[0]
			toState := model.State(prod.Right[1])
			trans := model.Transition{
				FromState: fromState,
				Input:     input,
				ToStates:  []model.State{toState},
			}
			automaton.Transitions = append(automaton.Transitions, trans)
			step.NewTransitions = []model.Transition{trans}
			step.Description = fmt.Sprintf("添加 转移: %s --%s--> %s", fromState, input, toState)

		// 5) 形如A -> w，w =a1a2a3...an，需展开为 A →a1→ A_1 →a2→ A_2 →...→ A_n-1 →an→ accept
		case rightLen >= 2 && g.CheckIsTerminal(prod.Right[rightLen-1]):
			currentState := fromState
			newStates := make([]model.State, 0)
			newTransitions := make([]model.Transition, 0)
			// 处理前n个符号（终结符）
			for i := 0; i < rightLen; i++ {
				input := prod.Right[i]
				if i == rightLen-1 {
					// 最后一步，转移到接受状态
					newTransitions = append(newTransitions, model.Transition{
						FromState: currentState,
						Input:     input,
						ToStates:  []model.State{model.UniqueFinalState},
					})
				} else {
					// 创建中间状态
					newState := getNextIntermediateStateIndex(fromState)
					newStates = append(newStates, newState)
					newTransitions = append(newTransitions, model.Transition{
						FromState: currentState,
						Input:     input,
						ToStates:  []model.State{newState},
					})
					currentState = newState
				}
			}

			// 更新自动机
			automaton.States = append(automaton.States, newStates...)
			automaton.Transitions = append(automaton.Transitions, newTransitions...)
			step.NewStates = newStates
			step.NewTransitions = newTransitions
			step.Description = fmt.Sprintf("处理终结符串产生式 %s → %s：引入 %d 个中间状态，最终到 %s",
				leftSym,
				util.SymbolsToString(prod.Right),
				len(newStates),
				model.UniqueFinalState)
			if !hasModelAccept {
				hasModelAccept = true
				step.Description += newStateFmt
			}

		// 6） 形如 A -> wB，w =a1a2a3...an，n>1，需展开为 A →a1→ A_1 →a2→ A_2 →...→ A_n-1 →an→ B
		case rightLen > 2 && g.CheckIsNonTerminal(prod.Right[rightLen-1]):
			// 最后一个是非终结符
			currentState := fromState
			newStates := make([]model.State, 0)
			newTransitions := make([]model.Transition, 0)

			// 处理前n-1个符号（终结符）
			for i := 0; i < rightLen-1; i++ {
				input := prod.Right[i]
				if i == rightLen-2 {
					// 最后一个终结符，转移到非终结符状态
					toState := model.State(prod.Right[rightLen-1])
					newTransitions = append(newTransitions, model.Transition{
						FromState: currentState,
						Input:     input,
						ToStates:  []model.State{toState},
					})
				} else {
					// 创建中间状态
					newState := getNextIntermediateStateIndex(fromState)
					newStates = append(newStates, newState)
					newTransitions = append(newTransitions, model.Transition{
						FromState: currentState,
						Input:     input,
						ToStates:  []model.State{newState},
					})
					currentState = newState
				}
			}

			// 更新自动机
			automaton.States = append(automaton.States, newStates...)
			automaton.Transitions = append(automaton.Transitions, newTransitions...)
			step.NewStates = newStates
			step.NewTransitions = newTransitions
			step.Description = fmt.Sprintf("处理长产生式 %s → %s：引入 %d 个中间状态",
				leftSym,
				util.SymbolsToString(prod.Right),
				len(newStates))

		default:
			// 理论上不应该发生
			step.Description = fmt.Sprintf("跳过非法产生式（应已通过类型检查）: %s → %s",
				util.SymbolsToString(prod.Left),
				util.SymbolsToString(prod.Right))
			zlog.Warn(step.Description)
		}
		process.Steps = append(process.Steps, step)
	}

	if hasModelAccept {
		automaton.States = append(automaton.States, model.UniqueFinalState)
		automaton.AcceptingStates = append(automaton.AcceptingStates, model.UniqueFinalState)
	}
	process.FinalAutomaton = &automaton

	return process
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
func leftLinearGrammarToFA(g *model.Grammar) *model.Automaton {
	return leftLinearGrammarToFAByBuild(g) // 默认使用直接构造法
	// return leftLinearGrammarToFAByReverse(g) // 默认使用间接法
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
//   - 对产生式A → a：添加转移 q0 →a→ A，这里的a可以是model.Epsilon
//   - 若开始符号S → ε，则q0也是终态
//
// δ(A, a) = {B|B→Aa∈P}
// δ(Z, a) = {B|B→a∈P}
// δ(Z, ε) = {B|B→ε∈P}
// 3. 终态设定：开始符号S对应的状态设为唯一终态（或根据ε产生式调整）
func leftLinearGrammarToFAByBuild(g *model.Grammar) *model.Automaton {
	var a model.Automaton
	// 1. 状态：每个非终结符作为一个状态；新增一个初始状态q0
	a.States = append(a.States, model.UniqueInitialState)
	a.InitialState = model.UniqueInitialState
	for _, nonTerm := range g.NonTerminals {
		a.States = append(a.States, model.State(nonTerm))
	}

	// 2. 符号
	for _, term := range g.Terminals {
		a.Alphabet = append(a.Alphabet, model.Symbol(term))
	}

	// 中间状态索引，形如 A -> ab……B 时，需要添加中间状态
	intermediateStateIndex := -1
	getNextIntermediateStateIndex := func(base model.Symbol) model.State {
		intermediateStateIndex++
		return model.State(fmt.Sprintf("%s_%d", base, intermediateStateIndex))
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
				newState := getNextIntermediateStateIndex(production.Left[0])
				newTransition.ToStates = []model.State{newState}
				a.States = append(a.States, newState)
				a.Transitions = append(a.Transitions, newTransition)
				// 新增转移 A1 →a2→ A2 →...→ An-1
				for i := 2; i < wLength; i++ {
					newTransition = model.Transition{
						FromState: newState,
						Input:     model.Symbol(production.Right[i]),
					}
					newState = getNextIntermediateStateIndex(production.Left[0])
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
				input := model.Symbol(production.Right[0])
				transition := model.Transition{
					FromState: a.InitialState,
					Input:     input,
					ToStates:  []model.State{model.State(production.Left[0])},
				}
				a.Transitions = append(a.Transitions, transition)
			} else { // 形如A -> w，w =a1a2a3...an，需展开为 q0 →a1→ A1 →a2→ A2 →...→ An-1 →an→ A
				// 新增转移 q0 →a1→ A1
				newTransition := model.Transition{
					FromState: a.InitialState,
					Input:     model.Symbol(production.Right[0]),
				}
				newState := getNextIntermediateStateIndex(production.Left[0])
				newTransition.ToStates = []model.State{newState}
				a.States = append(a.States, newState)
				a.Transitions = append(a.Transitions, newTransition)
				// 新增转移 A1 →a2→ A2 →...→ An-1
				for i := 1; i < wLength-1; i++ {
					newTransition = model.Transition{ // 新创建
						FromState: newState,
						Input:     model.Symbol(production.Right[i]),
					}
					newState = getNextIntermediateStateIndex(production.Left[0])
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

	// 在构建完所有状态和转移后，添加：
	a.AcceptingStates = []model.State{model.State(g.StartSymbol)}

	// 如果存在 S → ε 的产生式，则 q0 也应是接受状态
	hasEpsilonFromStart := false
	for _, prod := range g.Productions {
		if len(prod.Left) == 1 && prod.Left[0] == g.StartSymbol &&
			len(prod.Right) == 1 && prod.Right[0] == model.Epsilon {
			hasEpsilonFromStart = true
			break
		}
	}
	if hasEpsilonFromStart {
		a.AcceptingStates = append(a.AcceptingStates, model.UniqueInitialState)
	}

	return &a
}

// 左线性文法转自动机：直接构造（反向建图）
// 1. 状态设计：每个非终结符作为一个状态；新增一个初始状态q0
// 2. 转移规则：
//   - 对产生式A → Ba：添加转移 B →a→ A（注意方向是从B到A）
//   - 对产生式A → a：添加转移 q0 →a→ A，这里的a可以是model.Epsilon
//   - 若开始符号S → ε，则q0也是终态
//
// δ(A, a) = {B|B→Aa∈P}
// δ(Z, a) = {B|B→a∈P}
// δ(Z, ε) = {B|B→ε∈P}
// 3. 终态设定：开始符号S对应的状态设为唯一终态（或根据ε产生式调整）
func leftLinearGrammarToFAWithProcess(g *model.Grammar) *model.GrammarToFAProcess {
	var automaton model.Automaton
	process := &model.GrammarToFAProcess{
		Linear:          model.LeftLinear,
		OriginalGrammar: *g,
		Steps:           make([]model.GrammarToFAStep, 0),
	}

	// 左线性文法不能直接使用开始符号作为初始状态，
	// 因为其产生式从右向左生成字符串，而自动机从左向右读取输入。
	// 因此需引入新的初始状态 "initial" 来启动识别过程。
	// 创建唯一的初始状态
	initialState := model.UniqueInitialState
	automaton.InitialState = initialState
	automaton.States = append(automaton.States, initialState)

	// 添加非终结符为状态
	for _, nonTerm := range g.NonTerminals {
		automaton.States = append(automaton.States, model.State(nonTerm))
	}

	// 字母表
	automaton.Alphabet = append(automaton.Alphabet, g.Terminals...)

	// 记录步骤
	process.Steps = append(process.Steps, model.GrammarToFAStep{
		Action:      "init",
		Description: "初始化自动机：添加初始状态 'initial' 和所有非终结符状态",
		NewStates:   append([]model.State{initialState}, util.SymbolsToStates(g.NonTerminals)...),
	})

	// 中间状态计数器
	intermediateStateIndex := 0
	// 新增状态索引
	getNextIntermediateStateIndex := func(base model.State) model.State {
		intermediateStateIndex++
		return model.State(fmt.Sprintf("%s_%d", base, intermediateStateIndex))
	}

	// 遍历所有产生式
	for _, prod := range g.Productions {
		leftSym := prod.Left[0]
		fromNonTerminal := model.State(leftSym)
		right := prod.Right
		rightLength := len(right)

		step := model.GrammarToFAStep{
			Production: &prod,
			Action:     "handle_production",
		}

		switch {
		// 1) 形如 A -> ε 新增转移 q0 →ε→ A
		case rightLength == 1 && right[0] == model.Epsilon:
			input := model.Epsilon
			transition := model.Transition{
				FromState: initialState,
				Input:     input,
				ToStates:  []model.State{fromNonTerminal},
			}
			automaton.Transitions = append(automaton.Transitions, transition)
			step.NewTransitions = append(step.NewTransitions, transition)
			step.Description = fmt.Sprintf("处理产生式 %s → %s：新增转移 %s →%s→ %s", leftSym, right[0], initialState, input, fromNonTerminal)

		// 2) 形如 A -> a，新增转移 q0 →a→ B
		case rightLength == 1 && g.CheckIsTerminal(right[0]):
			input := right[0]
			transition := model.Transition{
				FromState: initialState,
				Input:     input,
				ToStates:  []model.State{fromNonTerminal},
			}
			automaton.Transitions = append(automaton.Transitions, transition)
			step.NewTransitions = append(step.NewTransitions, transition)
			step.Description = fmt.Sprintf("处理产生式 %s → %s：新增转移 %s →%s→ %s", leftSym, right[0], initialState, input, fromNonTerminal)

		// 3) 形如 A -> B，新增转移 B →ε→ A
		case rightLength == 1 && g.CheckIsNonTerminal(right[0]):
			from := model.State(right[0])
			transition := model.Transition{
				FromState: from,
				Input:     model.Epsilon,
				ToStates:  []model.State{fromNonTerminal},
			}
			automaton.Transitions = append(automaton.Transitions, transition)
			step.NewTransitions = append(step.NewTransitions, transition)
			step.Description = fmt.Sprintf("处理产生式 %s → %s：添加 ε-转移 %s --ε--> %s", leftSym, right[0], from, fromNonTerminal)

		// 4）形如 A -> Ba，新增转移 B →a→ A
		case rightLength == 2 && g.CheckIsNonTerminal(right[0]):
			input := right[1]
			toState := model.State(right[0])
			transition := model.Transition{
				FromState: toState,
				Input:     input,
				ToStates:  []model.State{fromNonTerminal},
			}
			automaton.Transitions = append(automaton.Transitions, transition)
			step.NewTransitions = append(step.NewTransitions, transition)
			step.Description = fmt.Sprintf("处理产生式 %s → %s：新增转移 %s →%s→ %s", leftSym, util.SymbolsToString(right), toState, input, fromNonTerminal)

		// 5）形如 A -> w，w =a1a2a3...an，需展开为 q0 →a1→ A1 →a2→ A2 →...→ An-1 →an→ A
		case rightLength >= 2 && g.CheckIsTerminal(right[0]):
			currentState := initialState
			newStates := make([]model.State, 0, rightLength-1)
			newTransitions := make([]model.Transition, 0, rightLength)

			for i, sym := range right {
				// 处理非最后一个符号（终结符）
				if i == rightLength-1 {
					// 转移到 A
					transition := model.Transition{
						FromState: currentState,
						Input:     sym,
						ToStates:  []model.State{fromNonTerminal},
					}
					newTransitions = append(newTransitions, transition)
				} else {
					// 新增中间状态
					newState := getNextIntermediateStateIndex(fromNonTerminal)
					newStates = append(newStates, newState)
					// 新增转移
					transition := model.Transition{
						FromState: currentState,
						Input:     sym,
						ToStates:  []model.State{newState},
					}
					newTransitions = append(newTransitions, transition)
					// 更新当前状态
					currentState = newState
				}
			}

			// 更新自动机
			automaton.States = append(automaton.States, newStates...)
			automaton.Transitions = append(automaton.Transitions, newTransitions...)

			step.NewStates = append(step.NewStates, newStates...)
			step.NewTransitions = append(step.NewTransitions, newTransitions...)
			step.Description = fmt.Sprintf("处理产生式 %s → %s：新增 %d 个中间状态和 %d 个转移",
				leftSym,
				util.SymbolsToString(right),
				len(newStates),
				len(newTransitions))

		// 6）形如 A -> Bw，w =a1a2a3...an，n>1，需展开为 B →a1→ A1 →a2→ A2 →...→ An-1 →an→ A
		case rightLength > 2 && g.CheckIsNonTerminal(right[0]):
			// 第一个符号是非终结符
			currentState := model.State(right[0])
			newStates := make([]model.State, 0, rightLength-1)
			newTransitions := make([]model.Transition, 0, rightLength)

			// 添加B →a1→ A1
			newState := getNextIntermediateStateIndex(fromNonTerminal)
			newStates = append(newStates, newState)
			transition := model.Transition{
				FromState: currentState,
				Input:     right[1],
				ToStates:  []model.State{newState},
			}
			newTransitions = append(newTransitions, transition)

			// 新增转移 A1 →a2→ A2 →...→ An-1 →an→ A
			for i := 1; i < rightLength; i++ {
				if i == rightLength-1 {
					// 转移到 A
					transition := model.Transition{
						FromState: currentState,
						Input:     right[i],
						ToStates:  []model.State{fromNonTerminal},
					}
					newTransitions = append(newTransitions, transition)
				} else {
					// 新增中间状态
					newState := getNextIntermediateStateIndex(fromNonTerminal)
					newStates = append(newStates, newState)
					// 新增转移
					transition := model.Transition{
						FromState: currentState,
						Input:     right[i],
						ToStates:  []model.State{newState},
					}
					newTransitions = append(newTransitions, transition)
					// 更新当前状态
					currentState = newState
				}
			}

			// 更新自动机
			automaton.States = append(automaton.States, newStates...)
			automaton.Transitions = append(automaton.Transitions, newTransitions...)

			step.NewStates = append(step.NewStates, newStates...)
			step.NewTransitions = append(step.NewTransitions, newTransitions...)
			step.Description = fmt.Sprintf("处理产生式 %s → %s：新增 %d 个中间状态和 %d 个转移",
				leftSym,
				util.SymbolsToString(right),
				len(newStates),
				len(newTransitions))
		default:
			// 理论上不应该发生
			step.Description = fmt.Sprintf("跳过非法产生式（应已通过类型检查）: %s → %s",
				util.SymbolsToString(prod.Left),
				util.SymbolsToString(prod.Right))
			zlog.Warn(step.Description)
		}
		process.Steps = append(process.Steps, step)
	}

	// 在构建完所有状态和转移后，添加：
	automaton.AcceptingStates = []model.State{model.State(g.StartSymbol)}

	// 如果存在 S → ε 的产生式，则 q0 也应是接受状态
	hasEpsilonFromStart := false
	for _, prod := range g.Productions {
		if len(prod.Left) == 1 && prod.Left[0] == g.StartSymbol &&
			len(prod.Right) == 1 && prod.Right[0] == model.Epsilon {
			hasEpsilonFromStart = true
			break
		}
	}
	if hasEpsilonFromStart {
		automaton.AcceptingStates = append(automaton.AcceptingStates, model.UniqueInitialState)
	}

	process.FinalAutomaton = &automaton

	return process
}

// 特性		右线性文法 → FA					左线性文法 → FA
// 初始状态	开始符号 S						新增状态 q₀
// 接受状态	新增 accept（或含 ε 时加 S）	开始符号 S（或含 ε 时加 q₀）
// 转移方向	A → aB ⇒ A --a--> B				A → Ba ⇒ B --a--> A
// 直观性	与输入顺序一致					需“反向建图”
