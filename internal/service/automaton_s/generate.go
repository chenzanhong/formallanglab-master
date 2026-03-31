package automaton_s

import (
	"fmt"

	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/pkg/util"
)

const maxExampleNum = 6

// 生成自动机可识别和不可识别的字符串示例。
func AutomatonGenerateExampleString(a *model.Automaton) (accept, reject []string) {
	// 转DFA+补集
	return automatonGenerateExampleStringByDFAAndBFS(a)
	// 枚举+验证
	// return automatonGenerateExampleStringByEnumAndVerify(a)
}

func automatonGenerateExampleStringByDFAAndBFS(a *model.Automaton) (accept, reject []string) {
	if a.Type != model.DFA {
		a = NFAToDFA(a)
	}
	if err := a.CompleteDFA(); err != nil { // 完备化失败
		fmt.Println("BFSShortestAcceptedStringsForDFA")
		return automatonGenerateExampleStringByEnumAndVerify(a)
	}

	return GenerateExampleStringsFromCompletedDFA(a)
}

// 根据一个完备的自动机，生成其可识别和不可识别的字符串示例
func GenerateExampleStringsFromCompletedDFA(a *model.Automaton) (accept, reject []string) {
	if a.Type != model.DFA {
		return automatonGenerateExampleStringByEnumAndVerify(a) // 备选
	}
	// accept: 在原 DFA 上 BFS 最多 maxNum 个
	accept = BFSShortestAcceptedStringsForDFA(a, 3*maxExampleNum)
	if len(accept) == 0 {
		accept = []string{"(no short accepted string found)"}
	}

	// reject: 在补集 DFA 上 BFS
	complementDFA := a.Clone()
	complementDFA.AcceptingStates = GetNonAcceptingStates(a)
	if err := complementDFA.CompleteDFA(); err != nil {
		reject = []string{"(failed to generate reject example)"}
	} else {
		reject = BFSShortestAcceptedStringsForDFA(complementDFA, 3*maxExampleNum)
		if len(reject) == 0 {
			reject = []string{"(all strings are accepted)"}
		}
	}

	if len(reject) > 1 && reject[0] == "" { // 对于拒绝的字符串优先展示非空的
		reject = reject[1:]
	}

	// 随机采样
	if len(accept) > maxExampleNum {
		accept = util.SamplingExampleStrings(accept, maxExampleNum)
	}
	if len(reject) > maxExampleNum {
		reject = util.SamplingExampleStrings(reject, maxExampleNum)
	}

	// 兜底
	if len(accept) == 0 {
		accept = []string{"(no short accepted string found)"}
	}
	if len(reject) == 0 {
		reject = []string{"(all short strings are accepted)"}
	}

	return accept, reject
}

// BFSShortestAcceptedStringsForDFA 在完备 DFA 上 BFS，返回最多 k 个最短的被接受字符串
func BFSShortestAcceptedStringsForDFA(dfa *model.Automaton, k int) []string {
	if dfa == nil || len(dfa.States) == 0 || k <= 0 {
		return nil
	}
	fmt.Println("BFSShortestAcceptedStringsForDFA")
	dfa.InitTransMap()

	acceptingSet := make(map[model.State]bool)
	for _, s := range dfa.AcceptingStates {
		acceptingSet[s] = true
	}

	var results []string
	seenStrings := make(map[string]bool)

	// 如果初始状态就是接受态
	if acceptingSet[dfa.InitialState] {
		results = append(results, "")
		seenStrings[""] = true
		if len(results) >= k {
			return results
		}
	}

	// BFS 队列
	type bfsNode struct {
		state model.State
		path  string
	}

	queue := []bfsNode{{state: dfa.InitialState, path: ""}}

	// 构建字母表（不含 ε）
	alphabet := []model.Symbol{}
	for _, sym := range dfa.Alphabet {
		if sym != model.Epsilon {
			alphabet = append(alphabet, sym)
		}
	}

	// 限制最大长度：避免无限循环，同时覆盖足够多短字符串
	maxLen := max(len(dfa.States), min(len(dfa.States)*2, 12)) // 经验值：不超过 12 位

	for len(queue) > 0 && len(results) < k {
		node := queue[0]
		queue = queue[1:]

		if len(node.path) >= maxLen {
			continue
		}

		for _, sym := range alphabet {
			transitions, ok := dfa.TransMap[node.state][sym]
			if !ok || len(transitions) == 0 {
				continue
			}
			nextState := transitions[0] // DFA

			newPath := node.path + string(sym)

			if seenStrings[newPath] {
				continue
			}
			seenStrings[newPath] = true

			if acceptingSet[nextState] {
				results = append(results, newPath)
				if len(results) >= k {
					return results
				}
			}

			// 即使 nextState 已经“见过”，也要继续扩展！
			// 因为 newPath 不同，后续可能产生新接受串
			queue = append(queue, bfsNode{
				state: nextState,
				path:  newPath,
			})
		}
	}

	return results
}

func automatonGenerateExampleStringByEnumAndVerify(a *model.Automaton) (accept, reject []string) {
	if a == nil || len(a.AcceptingStates) == 0 {
		return []string{"(no accepting states)"}, []string{"(any string is rejected)"}
	}

	a.InitTransMap()

	// 1. 提取纯终结符字母表（排除 ε）
	var alphabet []string
	for _, sym := range a.Alphabet {
		if sym != model.Epsilon {
			alphabet = append(alphabet, string(sym))
		}
	}

	// 2. 生成候选字符串（长度 0 ~ maxLen）
	maxLen := min(len(a.States)*2, len(a.States)+4)
	candidates := util.GenerateStrings(alphabet, maxLen)

	acceptMap := make(map[string]bool)
	rejectMap := make(map[string]bool)
	acceptByLen := make(map[int][]string)
	rejectByLen := make(map[int][]string)

	// 3. 用 Recognize 判断每个字符串
	for _, s := range candidates {
		result, _ := Recognize(a, s)
		if result.IsAccepted {
			if !acceptMap[s] {
				acceptByLen[len(s)] = append(acceptByLen[len(s)], s)
				acceptMap[s] = true
			}
		} else {
			if !rejectMap[s] {
				rejectByLen[len(s)] = append(rejectByLen[len(s)], s)
				rejectMap[s] = true
			}
		}
	}

	// 随机采样
	if len(accept) > maxExampleNum {
		accept = util.SamplingExampleStrings(accept, maxExampleNum)
	}
	if len(reject) > maxExampleNum {
		reject = util.SamplingExampleStrings(reject, maxExampleNum)
	}

	// 兜底
	if len(accept) == 0 {
		accept = []string{"(no short accepted string found)"}
	}
	if len(reject) == 0 {
		reject = []string{"(all short strings are accepted)"}
	}

	return accept, reject
}
