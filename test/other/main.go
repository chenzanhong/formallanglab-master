package main

import (
	"backend/internal/domain/model"
	"backend/internal/service/automaton_s"
	"backend/internal/service/grammar_s"
	"backend/internal/service/regex_s"
	"fmt"
)

var grammar = &model.Grammar{
	StartSymbol:  model.Symbol("S"),
	NonTerminals: []model.Symbol{"S", "A", "B"},
	Terminals:    []model.Symbol{"a", "b"},
	Productions: []model.Production{
		{[]model.Symbol{"S"}, []model.Symbol{"a", "A"}},
		{[]model.Symbol{"S"}, []model.Symbol{"b", "B"}},
		{[]model.Symbol{"A"}, []model.Symbol{"B"}},
		{Left: []model.Symbol{"B"}, Right: []model.Symbol{"a", "A"}},
		{[]model.Symbol{"B"}, []model.Symbol{model.Epsilon}},
	},
	GrammarType: model.RegularGrammar,
}

var automaton = &model.Automaton{
	States:          []model.State{"q0", "q1", "q2", "q3"},
	InitialState:    "q0",
	Alphabet:        []model.Symbol{"a", "b", "c"},
	AcceptingStates: []model.State{"q3"},
	Transitions: []model.Transition{
		{FromState: "q0", Input: "a", ToStates: []model.State{"q1"}},
		{FromState: "q1", Input: "b", ToStates: []model.State{"q2"}},
		{FromState: "q1", Input: model.Epsilon, ToStates: []model.State{"q2"}},
		{FromState: "q2", Input: "a", ToStates: []model.State{"q2"}},
		{FromState: "q2", Input: "c", ToStates: []model.State{"q3"}},
		{FromState: "q3", Input: "a", ToStates: []model.State{"q3"}},
	},
}

var regex = model.Regex("a(b|c)*a+")

func main() {
	// testGrammarGenerateExampleString()
	// testAutomatonGenerateExampleString()
	testRegexGenerateExampleString()
}

func testGrammarGenerateExampleString() {
	accept, reject := grammar_s.GrammarGenerateExampleString(grammar)
	for _, s := range accept {
		fmt.Println(s)
	}
	fmt.Println("---")
	for _, s := range reject {
		fmt.Println(s)
	}
}

func testAutomatonGenerateExampleString() {
	ss, sj := automaton_s.AutomatonGenerateExampleString(automaton)
	for _, s := range ss {
		fmt.Println(s)
	}
	fmt.Println("---")
	for _, s := range sj {
		fmt.Println(s)
	}
}

func testRegexGenerateExampleString() {
	accept, reject := regex_s.RegexGenerateExampleString(regex)
	for _, s := range accept {
		fmt.Println(s)
	}
	fmt.Println("---")
	for _, s := range reject {
		fmt.Println(s)
	}
}
