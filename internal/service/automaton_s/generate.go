package automaton_s

import (
	"backend/internal/domain/model"
	"backend/pkg/util"
)

const maxExampleNum = 6

func AutomatonGenerateExampleString(a *model.Automaton) (accept, reject []string) {
	a.InitTransMap()

	if a == nil || len(a.AcceptingStates) == 0 {
		return []string{"(no accepting states)"}, []string{"(any string is rejected)"}
	}

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

	// 确保数量不超过maxExampleNum
	accept = util.SamplingExampleStrings(acceptByLen, maxExampleNum)
	reject = util.SamplingExampleStrings(rejectByLen, maxExampleNum)

	// 兜底
	if len(accept) == 0 {
		accept = []string{"(no short accepted string found)"}
	}
	if len(reject) == 0 {
		reject = []string{"(all short strings are accepted)"}
	}

	return accept, reject
}
