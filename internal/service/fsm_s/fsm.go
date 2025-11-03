package fsm_s

import "backend/internal/domain/model"

func getDFANextStateFromMap(fsm *model.Automaton, from model.State, input model.Symbol) model.State {
	// 使用 TransMap 提高查找效率
	if stateMap, exists := fsm.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists && len(targets) == 1 {
			return targets[0] //
		}
	}
	return ""
}

func getNFANextStates(fsm *model.Automaton, from model.State, input model.Symbol) []model.State {
	// 使用 TransMap 提高查找效率
	if stateMap, exists := fsm.TransMap[from]; exists {
		if targets, exists := stateMap[input]; exists {
			return targets
		}
	}
	return nil
}

// containsState 检查状态是否在状态列表中
func containsState(states []model.State, s model.State) bool {
	for _, st := range states {
		if st == s {
			return true
		}
	}
	return false
}
