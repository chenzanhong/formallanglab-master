package main

import (
	"backend/internal/domain/model"
	"backend/internal/service/convert"
	"fmt"
)

func main() {
	// 测试 FAToRegex
	automaton := &model.Automaton{
		States:          []model.State{"q0", "q1", "q2", "q3"},
		InitialState:    "q0",
		Alphabet:        []model.Symbol{"a", "b"},
		AcceptingStates: []model.State{"q3"},
		Transitions: []model.Transition{
			{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}},
			{FromState: "q1", Input: "b", ToStates: []model.State{"q2"}},
			{FromState: "q1", Input: "a", ToStates: []model.State{"q2"}},
			{FromState: "q2", Input: "a", ToStates: []model.State{"q2"}},
			{FromState: "q2", Input: "a", ToStates: []model.State{"q3"}},
			{FromState: "q3", Input: "a", ToStates: []model.State{"q3"}},
		},
	}

	regex := convert_s.FAToRegex(automaton)
	fmt.Println(regex)
}
