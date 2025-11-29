package generate

import (
	"backend/internal/domain/model"
	"backend/internal/service/automaton_s"
	"backend/internal/service/grammar_s"
	"backend/internal/service/regex_s"
	"backend/pkg/util"
	"fmt"
)

// GenerateExampleStrings 统一生成可识别和不可识别的字符串示例
// 优先使用 DFA+BFS 方法；若不支持（如 CFG），则 fallback 到枚举
func GenerateExampleStrings(source interface{}) (accept, reject []string) {
	const maxNum = 6

	// 尝试构建原 DFA 和补集 DFA
	originalDFA, err := buildDFAFromSource(source)
	if err != nil || originalDFA == nil {
		// fallback 到各模块的枚举方法
		return fallbackGenerate(source)
	}

	// 确保是完备 DFA
	if err := originalDFA.CompleteDFA(); err != nil {
		return fallbackGenerate(source)
	}

	// accept: 在原 DFA 上 BFS 最多 maxNum 个
	accept = automaton_s.BFSShortestAcceptedStringsForDFA(originalDFA, 3*maxNum)
	if len(accept) == 0 {
		accept = []string{"(no short accepted string found)"}
	}

	// reject: 在补集 DFA 上 BFS
	complementDFA := originalDFA.Clone()
	complementDFA.AcceptingStates = automaton_s.GetNonAcceptingStates(originalDFA)
	if err := complementDFA.CompleteDFA(); err != nil {
		reject = []string{"(failed to generate reject example)"}
	} else {
		reject = automaton_s.BFSShortestAcceptedStringsForDFA(complementDFA, 3*maxNum)
		if len(reject) == 0 {
			reject = []string{"(all strings are accepted)"}
		}
	}

	// 随机采样
	if len(accept) > maxNum {
		accept = util.SamplingExampleStrings(accept, maxNum)
	}
	if len(reject) > maxNum {
		reject = util.SamplingExampleStrings(reject, maxNum)
	}

	return accept, reject
}

// buildDFAFromSource 将 source 转为 DFA（不补全）
func buildDFAFromSource(source interface{}) (*model.Automaton, error) {
	switch s := source.(type) {
	case model.Regex:
		if err := regex_s.RegexValidate(s); err != nil {
			return nil, fmt.Errorf("非法的正则表达式")
		}
		nfa, err := regex_s.RegexToFA(s)
		if err != nil {
			return nil, err
		}
		return automaton_s.NFAToDFA(nfa), nil

	case *model.Automaton:
		if s == nil {
			return nil, fmt.Errorf("automaton is nil")
		}
		if err := s.ISValidate(); err != nil {
			return nil, fmt.Errorf("无效的自动机")
		}
		if s.Type != model.DFA {
			return automaton_s.NFAToDFA(s), nil
		}
		return s.Clone(), nil

	case *model.Grammar:
		if s == nil {
			return nil, fmt.Errorf("grammar is nil")
		}
		if err := grammar_s.GrammarCheckValidity(s); err != nil {
			return nil, fmt.Errorf("无效文法：%v", err)
		}
		grammar_s.TypeDetermine(s)
		if s.GrammarType != model.RegularGrammar {
			return nil, fmt.Errorf("非正则文法，暂不支持") // 不支持，触发 fallback
		}
		nfa, err := grammar_s.RegularGrammarToFA(s)
		if err != nil { // 理论上不会出现，因为已经通过TypeDetermine，确认是正则文法了
			return nil, fmt.Errorf("非线性文法")
		}
		return automaton_s.NFAToDFA(nfa), nil

	default:
		return nil, fmt.Errorf("unsupported type: %T", source)
	}
}

// fallbackGenerate 回退到各模块的枚举验证方法
func fallbackGenerate(source interface{}) (accept, reject []string) {
	switch s := source.(type) {
	case model.Regex:
		return regex_s.RegexGenerateExampleString(s)
	case *model.Automaton:
		return automaton_s.AutomatonGenerateExampleString(s)
	case *model.Grammar:
		return grammar_s.GrammarGenerateExampleString(s)
	default:
		msg := "(unsupported type)"
		return []string{msg}, []string{msg}
	}
}
