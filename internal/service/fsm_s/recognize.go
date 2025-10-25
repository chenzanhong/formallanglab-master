package fsm_s

import (
	"backend/internal/domain/model"
	"fmt"
)

// Recognize 判断自动机是否接受输入字符串 str
// 根据自动机类型（DFA或NFA）选择相应的识别算法
// 
// 算法说明：
// 对于DFA（确定性有限自动机）：
//   1. 从初始状态开始，依次读取输入字符串中的每个符号
//   2. 根据当前状态和输入符号，按照转移函数转移到下一个状态
//   3. 读取完所有符号后，检查当前状态是否为接受状态
//
// 对于NFA（非确定性有限自动机）：
//   1. 从初始状态开始，维护一个可能的当前状态集合
//   2. 依次读取输入字符串中的每个符号
//   3. 对于当前状态集合中的每个状态，计算在该符号下的所有可能转移
//   4. 将所有可能的转移目标合并为新的当前状态集合
//   5. 读取完所有符号后，检查当前状态集合中是否包含任何接受状态
//
// Symbol 是 string 类型，使用 map + 固定长度集合进行高效匹配
func Recognize(fsm *model.Automaton, str string) (bool, error) {
	if fsm == nil || len(fsm.States) == 0 {
		return false, fmt.Errorf("自动机状态为空")
	}

	// 初始化 TransMap 以提高识别效率
	fsm.InitTransMap()

	if fsm.IsDFA {
		return recognizeDFA(fsm, str)
	} else {
		return recognizeNFA(fsm, str)
	}
}

func recognizeDFA(fsm *model.Automaton, str string) (bool, error) {
	symbols, err := fsm.SplitString(str)
	if err != nil {
		return false, fmt.Errorf("DFA 识别失败：无法将输入字符串分词 -> %w", err)
	}
	currentState := fsm.InitialState
	for _, sym := range symbols {
		nextState := getDFANextState2(fsm, currentState, sym)
		if nextState == "" {
			return false, fmt.Errorf("当前状态 %s 下输入 %s 没有相关的有效转移", currentState, sym)
		}
		currentState = nextState
	}
	for _, st := range fsm.AcceptingStates {
		if st == currentState {
			return true, nil
		}
	}
	return false, fmt.Errorf("输入字符串 %s 被完整识别，但是未到达接收状态", str)
}

func recognizeNFA(fsm *model.Automaton, str string) (bool, error) {
	// 第一步：将字符串分词为符号序列
	symbols, err := fsm.SplitString(str)
	if err != nil {
		return false, fmt.Errorf("NFA 识别失败：无法将输入字符串分词 -> %w", err)
	}

	// 当前可能处于的状态集合（NFA 的核心：状态集合）
	currentStates := []model.State{fsm.InitialState}

	// 逐个处理每个符号
	for _, sym := range symbols {
		var nextStates []model.State
		found := false

		// 对当前每个可能的状态，查找该符号的转移
		for _, state := range currentStates {
			targets := GetNFANextStates(fsm, state, sym)
			if len(targets) > 0 {
				found = true
				// 将目标状态加入 nextStates（去重）
				for _, t := range targets {
					if !containsState(nextStates, t) {
						nextStates = append(nextStates, t)
					}
				}
			}
		}

		// 如果没有任何状态能处理当前符号
		if !found {
			return false, fmt.Errorf("NFA 识别失败：在输入符号 '%s' 时，当前状态集合中没有状态可以转移", sym)
		}

		// 更新当前状态集合
		currentStates = nextStates
	}

	// 最终：检查当前状态集合中是否有任意一个接受状态
	for _, s := range currentStates {
		if containsState(fsm.AcceptingStates, s) {
			return true, nil
		}
	}

	return false, fmt.Errorf("输入字符串 %s 被完整识别，但未到达任何接受状态", str)
}

func getDFANextState2(fsm *model.Automaton, from model.State, input model.Symbol) model.State {
	// 使用 TransMap 提高查找效率
	if stateMap, exists := fsm.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists && len(targets) == 1 {
			return targets[0]
		}
	}
	return ""
}

func GetNFANextStates(fsm *model.Automaton, from model.State, input model.Symbol) []model.State {
	// 使用 TransMap 提高查找效率
	if stateMap, exists := fsm.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists {
			return targets
		}
	}
	return nil
}

func containsState(states []model.State, s model.State) bool {
	for _, st := range states {
		if st == s {
			return true
		}
	}
	return false
}
