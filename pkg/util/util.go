package util

import (
	"backend/internal/domain/model"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// []model.Symbol转string，空切片时返回string(model.Epsilon)
func SymbolsToString(syms []model.Symbol) string {
	if len(syms) == 0 {
		return string(model.Epsilon)
	}
	var s strings.Builder
	for _, sym := range syms {
		s.WriteString(string(sym))
	}
	return s.String()
}

func GenerateStrings(terminals []string, maxLen int) []string {
	if len(terminals) == 0 {
		return []string{""}
	}

	// 生成所有可能的字符串（包括空字符串）
	var result []string
	queue := []string{terminals[0]}

	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		result = append(result, s)

		if len(s) < maxLen {
			for _, c := range terminals {
				queue = append(queue, s+c)
			}
		}
	}

	// 按长度排序
	sort.Slice(result, func(i, j int) bool {
		if len(result[i]) != len(result[j]) {
			return len(result[i]) < len(result[j])
		}
		return result[i] < result[j]
	})
	fmt.Println("len:", len(result), " maxlen:", maxLen)
	return result
}

func SamplingExampleStrings(samplingByLen map[int][]string, maxExampleNum int) []string {
	var sampling []string
	lengths := make([]int, 0, len(samplingByLen))
	for l := range samplingByLen {
		lengths = append(lengths, l)
	}
	sort.Ints(lengths)

	// 对每一层的 sampling 字符串列表进行随机打乱
	for l := range samplingByLen {
		rand.Shuffle(len(samplingByLen[l]), func(i, j int) {
			samplingByLen[l][i], samplingByLen[l][j] = samplingByLen[l][j], samplingByLen[l][i]
		})
	}

	var has bool
	for i, j := 0, 0; i < maxExampleNum; j++ {
		has = false
		for _, l := range lengths {
			if j < len(samplingByLen[l]) {
				sampling = append(sampling, samplingByLen[l][j])
				i++
				has = true
				if i >= maxExampleNum {
					break
				}
			}
		}
		if !has {
			break
		}
	}
	return sampling
}

// PrintAutomaton 打印自动机的详细信息，便于调试
func PrintAutomaton(a *model.Automaton) {
	if a == nil {
		fmt.Println("Automaton is nil")
		return
	}

	fmt.Println("=== Automaton ===")

	// 类型
	typeStr := map[model.AutomatonType]string{
		model.DFA:        "DFA",
		model.NFA:        "NFA",
		model.EpsilonNFA: "ε-NFA",
	}[a.Type]
	fmt.Printf("Type: %s\n", typeStr)

	// 状态集合
	fmt.Printf("States: [%s]\n", strings.Join(quoteStates(a.States), ", "))

	// 初始状态
	fmt.Printf("Initial State: %q\n", a.InitialState)

	// 接受状态
	fmt.Printf("Accepting States: [%s]\n", strings.Join(quoteStates(a.AcceptingStates), ", "))

	// 字母表（隐藏 ε，因为它不应在 Alphabet 中；但若存在也显示）
	alphabetToShow := make([]string, 0, len(a.Alphabet))
	for _, sym := range a.Alphabet {
		if sym == model.Epsilon {
			alphabetToShow = append(alphabetToShow, "ε")
		} else {
			alphabetToShow = append(alphabetToShow, fmt.Sprintf("%q", string(sym)))
		}
	}
	fmt.Printf("Alphabet: [%s]\n", strings.Join(alphabetToShow, ", "))

	// 转移规则
	fmt.Println("Transitions:")
	if len(a.Transitions) == 0 {
		fmt.Println("  (none)")
	} else {
		for i, t := range a.Transitions {
			input := func() string {
				if t.Input == model.Epsilon {
					return "ε"
				}
				return fmt.Sprintf("%q", string(t.Input))
			}()
			toStrs := quoteStates(t.ToStates)
			fmt.Printf("  %d. %q --%s--> [%s]\n", i+1, t.FromState, input, strings.Join(toStrs, ", "))
		}
	}
}

// 辅助函数：将 []State 转为带引号的字符串切片
func quoteStates(states []model.State) []string {
	result := make([]string, len(states))
	for i, s := range states {
		result[i] = fmt.Sprintf("%q", s)
	}
	return result
}
