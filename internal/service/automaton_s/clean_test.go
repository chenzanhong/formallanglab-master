package automaton_s

import (
	"fmt"
	"testing"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// automatonEquals 判断两个自动机是否相等（不考虑状态顺序、转移顺序等）
func automatonEquals(a, b *model.Automaton) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// 比较基本信息
	if a.InitialState != b.InitialState || a.Type != b.Type {
		return false
	}

	// 比较状态集合（不考虑顺序）
	if !stateSetEqual(a.States, b.States) {
		return false
	}

	// 比较字母表（不考虑顺序）
	if !symbolSetEqual(a.Alphabet, b.Alphabet) {
		return false
	}

	// 比较接受状态集合（不考虑顺序）
	if !stateSetEqual(a.AcceptingStates, b.AcceptingStates) {
		return false
	}

	// 比较转移集合（不考虑顺序）
	return transitionSetEqual(a.Transitions, b.Transitions)
}

// stateSetEqual 判断两个状态集合是否相等（不考虑顺序）
func stateSetEqual(a, b []model.State) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[model.State]bool)
	for _, s := range a {
		set[s] = true
	}
	for _, s := range b {
		if !set[s] {
			return false
		}
	}
	return true
}

// symbolSetEqual 判断两个符号集合是否相等（不考虑顺序）
func symbolSetEqual(a, b []model.Symbol) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[model.Symbol]bool)
	for _, s := range a {
		set[s] = true
	}
	for _, s := range b {
		if !set[s] {
			return false
		}
	}
	return true
}

// transitionSetEqual 判断两个转移集合是否相等（不考虑顺序）
func transitionSetEqual(a, b []model.Transition) bool {
	if len(a) != len(b) {
		return false
	}

	// 将转移转换为字符串进行比较
	toStr := func(t model.Transition) string {
		return fmt.Sprintf("%s-%s-%v", t.FromState, t.Input, t.ToStates)
	}

	set := make(map[string]bool)
	for _, t := range a {
		set[toStr(t)] = true
	}
	for _, t := range b {
		if !set[toStr(t)] {
			return false
		}
	}
	return true
}

func TestCleanup(t *testing.T) {
	tests := []struct {
		name      string
		automaton *model.Automaton
		want      *model.Automaton
	}{
		{
			name: "automaton with unreachable states",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2", "q3"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
			want: &model.Automaton{
				States:          []model.State{"q0", "q1"},
				Alphabet:        []model.Symbol{"a"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q1"},
				Type:            model.DFA,
			},
		},
		{
			name: "automaton with all states reachable",
			automaton: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q2"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.DFA,
			},
			want: &model.Automaton{
				States:          []model.State{"q0", "q1", "q2"},
				Alphabet:        []model.Symbol{"a", "b"},
				Transitions:     []model.Transition{{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}}, {FromState: "q1", Input: "b", ToStates: []model.State{"q2"}}},
				InitialState:    "q0",
				AcceptingStates: []model.State{"q2"},
				Type:            model.DFA,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automaton := tt.automaton.Clone()
			Cleanup(automaton)

			if !automatonEquals(automaton, tt.want) {
				t.Errorf("Cleanup() result mismatch\n got: %v\nwant: %v", automaton, tt.want)
			}

			if err := automaton.Validate(); err != nil {
				t.Errorf("Cleanup() result validation failed: %v", err)
			}
		})
	}
}

func TestCleanupWithNil(t *testing.T) {
	Cleanup(nil)
}
