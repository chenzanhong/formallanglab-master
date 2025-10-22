package fsm_s

import (
	"backend/internal/domain/model"
	"fmt"
	"sort"
)

// DFAMinimize_ 使用Hopcroft算法对DFA进行最小化
func DFAMinimize_(fsm *model.Automaton) *model.Automaton {
	// 步骤1: 去除不可达状态
	reachable := getReachableStates(fsm)
	var states []model.State
	for s := range reachable {
		states = append(states, s)
	}

	// 过滤状态、转移、接受状态
	newAccepting := []model.State{}
	for _, s := range fsm.AcceptingStates {
		if reachable[s] {
			newAccepting = append(newAccepting, s)
		}
	}

	// 创建新自动机（仅含可达状态）
	reduced := &model.Automaton{
		States:          states,
		Alphabet:        fsm.Alphabet,
		Transitions:     []model.Transition{},
		InitialState:    fsm.InitialState,
		AcceptingStates: newAccepting,
		IsDFA:           true,
	}

	// 如果只剩一个状态，直接返回
	if len(states) <= 1 {
		return reduced
	}

	// 步骤2: Hopcroft算法初始化
	// 初始化划分：接受状态和非接受状态
	acceptingSet := make(map[model.State]bool)
	for _, s := range newAccepting {
		acceptingSet[s] = true
	}

	// 初始划分：P = {F, Q-F}
	partitions := []map[model.State]bool{
		make(map[model.State]bool), // 非接受状态
		make(map[model.State]bool), // 接受状态
	}
	for _, s := range states {
		if acceptingSet[s] {
			partitions[1][s] = true
		} else {
			partitions[0][s] = true
		}
	}

	// 工作列表：存储需要细分的划分
	worklist := []map[model.State]bool{partitions[0], partitions[1]}

	// 步骤3: Hopcroft算法主循环
	for len(worklist) > 0 {
		// 取出当前划分
		current := worklist[0]
		worklist = worklist[1:]

		// 对每个输入符号尝试细分
		for _, a := range fsm.Alphabet {
			// 收集所有状态通过a转移后到达current的状态
			inverseMap := make(map[model.State][]model.State)
			for s := range current {
				next := getDFANextState_(reduced, s, a)
				if next != "" { // 忽略无转移的情况（理论上DFA应完备）
					inverseMap[next] = append(inverseMap[next], s)
				}
			}

			// 尝试细分current
			for _, group := range partitions {
				// 找出group中通过a转移到current中不同子集的状态
				toSplit := make(map[model.State]bool)
				for s := range group {
					next := getDFANextState_(reduced, s, a)
					if next != "" {
						// 检查next是否属于current（因为我们在处理current的逆映射）
						if _, exists := current[next]; exists {
							toSplit[s] = true
						}
					}
				}

				// 如果没有状态需要分裂，继续
				if len(toSplit) == 0 || len(toSplit) == len(group) {
					continue
				}

				// 分裂group
				newGroup := make(map[model.State]bool)
				for s := range toSplit {
					newGroup[s] = true
					delete(group, s)
				}

				// 添加新划分到工作列表
				worklist = append(worklist, newGroup)
				partitions = append(partitions, newGroup)
			}
		}
	}

	// 步骤4: 构建等价类
	equivalenceClasses := make([][]model.State, 0, len(partitions))
	for _, group := range partitions {
		if len(group) == 0 {
			continue
		}
		class := make([]model.State, 0, len(group))
		for s := range group {
			class = append(class, s)
		}
		sort.Slice(class, func(i, j int) bool {
			return string(class[i]) < string(class[j])
		})
		equivalenceClasses = append(equivalenceClasses, class)
	}

	// 步骤5: 构建最小化DFA
	return buildMinimizedDFA_(reduced, equivalenceClasses, acceptingSet)
}

// getDFANextState_ 获取DFA在状态s上输入a的转移目标
func getDFANextState_(fsm *model.Automaton, s model.State, a model.Symbol) model.State {
	for _, t := range fsm.Transitions {
		if t.FromState == s && t.Input == a {
			return t.ToStates[0] // DFA每个转移只有一个目标状态
		}
	}
	return "" // 无转移（理论上DFA应完备）
}

func buildMinimizedDFA_(
	fsm *model.Automaton,
	classes [][]model.State,
	acceptingSet map[model.State]bool,
) *model.Automaton {
	// 映射：状态 → 所属类
	stateToClass := make(map[model.State][]model.State)
	classRepresentatives := make(map[model.State]model.State) // 每个类的代表状态

	for _, cls := range classes {
		// 使用第一个状态作为代表
		rep := cls[0]
		for _, s := range cls {
			stateToClass[s] = cls
			classRepresentatives[s] = rep
		}
	}

	// 新状态名：[q0,q1]
	var newStates []model.State
	var newAccepting []model.State
	var newTransitions []model.Transition

	// 创建新状态和接受状态
	stateNames := make(map[model.State]model.State) // 原始状态到新状态名的映射
	for _, cls := range classes {
		name := model.State(fmt.Sprintf("%v", cls))
		newStates = append(newStates, name)
		stateNames[cls[0]] = name // 使用类中第一个状态作为代表

		// 检查是否为接受状态类
		for _, s := range cls {
			if acceptingSet[s] {
				newAccepting = append(newAccepting, name)
				break
			}
		}
	}

	// 构建转移
	for _, cls := range classes {
		fromRep := cls[0]
		fromName := stateNames[fromRep]

		for _, a := range fsm.Alphabet {
			next := getDFANextState_(fsm, fromRep, a)
			if next == "" {
				continue // 无转移（理论上不应发生）
			}
			nextName := stateNames[classRepresentatives[next]]

			newTransitions = append(newTransitions, model.Transition{
				FromState: fromName,
				Input:     a,
				ToStates:  []model.State{nextName},
			})
		}
	}

	// 确定初始状态
	initialRep := classRepresentatives[fsm.InitialState]
	initialName := stateNames[initialRep]

	return &model.Automaton{
		States:          newStates,
		Alphabet:        fsm.Alphabet,
		Transitions:     newTransitions,
		InitialState:    initialName,
		AcceptingStates: newAccepting,
		IsDFA:           true,
	}
}