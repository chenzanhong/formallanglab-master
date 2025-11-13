package convert_s

import (
	"backend/internal/domain/model"
	"fmt"
	"strings"
)

// concatRegex 拼接两个正则表达式，注意空串和 ε 的处理
func concatRegex(r1, r2 model.Symbol) model.Symbol {
	// 不可达
	if r1 == "" || r1 == "∅" {
		return "∅"
	}
	if r2 == "" || r2 == "∅" {
		return "∅"
	}
	// ε 连接不改变
	if r1 == "ε" {
		return r2
	}
	if r2 == "ε" {
		return r1
	}
	// 加括号避免歧义，对包含 '|' 且 '|' 未被括号包裹的表达式加括号
	wrap := func(s model.Symbol) string {
		str := string(s)
		if str[0] != '(' && needWrap(str) {
			return "(" + str + ")"
		}
		return str
	}
	return model.Symbol(wrap(r1) + wrap(r2))
}

// 检查字符串的|是否需要加括号，分别从左到右和从右到遍历字符串，记录每个位置的(与)的数量
func needWrap(s string) bool {
	leftCount := 0
	rightCount := 0
	has := false
	for i := 0; i < len(s); i++ {
		if s[i] == '|' {
			has = true
			break
		}
		if s[i] == '(' {
			leftCount++
		}
	}
	// 从右到左遍历
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '|' {
			has = true
			break
		}
		if s[i] == ')' {
			rightCount++
		}
	}
	fmt.Println(leftCount, " ", rightCount)
	return has && (leftCount != rightCount || (leftCount == 0 && rightCount == 0))
}

// unionRegex 并联两个正则表达式
func unionRegex(r1, r2 model.Symbol) model.Symbol {
	if r1 == "∅" {
		return r2
	}
	if r2 == "∅" {
		return r1
	}
	if r1 == r2 {
		return r1
	}
	if (r1 == "ε" && r2 == "") || (r2 == "ε" && r1 == "") {
		return "ε"
	}
	return model.Symbol(string(r1) + "|" + string(r2))
}

// starRegex Kleene 星
func starRegex(a *model.Automaton, r model.Symbol) model.Symbol {
	if r == "∅" || r == "ε" {
		return "ε"
	}
	str := string(r)
	// 单一符号或括号包裹的符号
	if containSymbol(a, r) || (strings.HasPrefix(str, "(") && strings.HasSuffix(str, ")")) {
		return model.Symbol(str + "*")
	}
	return model.Symbol("(" + str + ")*")
}

func containSymbol(a *model.Automaton, sym model.Symbol) bool {
	for _, s := range a.Alphabet {
		fmt.Println("sym:", sym, " s:", s)
		if s == sym {
			return true
		}
	}
	return false
}

// FAToRegex 将 NFA/DFA 转换为等价的正则表达式（状态消除法）
func FAToRegex(a *model.Automaton) model.Regex {
	// Step 1: 预处理 —— 添加唯一初态和唯一终态
	newInitialState := model.State("initial")
	a.States = append(a.States, newInitialState)
	a.Transitions = append(a.Transitions, model.Transition{
		FromState: newInitialState,
		Input:     model.Epsilon,
		ToStates:  []model.State{a.InitialState},
	})
	a.InitialState = newInitialState

	newAcceptingState := model.State("accept")
	a.States = append(a.States, newAcceptingState)
	for _, acceptState := range a.AcceptingStates {
		a.Transitions = append(a.Transitions, model.Transition{
			FromState: acceptState,
			Input:     model.Epsilon,
			ToStates:  []model.State{newAcceptingState},
		})
	}
	a.AcceptingStates = []model.State{newAcceptingState}

	finalState := a.AcceptingStates[0]

	// Step 2: 构建状态间正则表达式的邻接矩阵（map of map）
	// regexMap[i][j] = 从 i 到 j 的当前正则表达式（初始为 ∅）
	regexMap := make(map[model.State]map[model.State]model.Symbol)
	for _, s1 := range a.States {
		regexMap[s1] = make(map[model.State]model.Symbol)
		for _, s2 := range a.States {
			regexMap[s1][s2] = "∅"
		}
		regexMap[s1][s1] = model.Epsilon // 自环初始为 ε（允许不走）
	}

	// 初始化转移
	for _, t := range a.Transitions {
		for _, to := range t.ToStates {
			input := t.Input
			current := regexMap[t.FromState][to]
			if current == "∅" || t.FromState == to { // 自环或空转移
				regexMap[t.FromState][to] = input
			} else {
				regexMap[t.FromState][to] = unionRegex(current, input)
			}

			fmt.Println("from:", t.FromState, "input:", input, "to:", to, "regex:", regexMap[t.FromState][to])
		}
	}

	// Step 3: 获取要消除的状态列表（排除 initial 和 accept）
	var statesToEliminate []model.State
	for _, s := range a.States {
		if s != a.InitialState && s != finalState {
			statesToEliminate = append(statesToEliminate, s)
		}
	}
	// 记录被消除的状态
	var eliminatedStates map[model.State]bool = make(map[model.State]bool)

	// Step 4: 逐个消除状态
	for _, r := range statesToEliminate {
		// R_rr*
		loop := regexMap[r][r]
		starLoop := starRegex(a, loop)

		// 对每对 (i, j)，更新路径 i -> r -> j 为 i -> j
		for _, i := range a.States {
			if i == r || eliminatedStates[i] {
				continue
			}
			for _, j := range a.States {
				if j == r || eliminatedStates[j] {
					continue
				}
				// i -> r
				ir := regexMap[i][r]
				if ir == "∅" {
					continue
				}
				// r -> j
				rj := regexMap[r][j]
				if rj == "∅" {
					continue
				}

				// 新路径：ir · (loop)* · rj
				var newPath model.Symbol
				if loop == model.Epsilon {
					newPath = concatRegex(ir, rj)
					fmt.Println(1, " ", i, " ", r, " ", j)
				} else {
					newPath = concatRegex(concatRegex(ir, starLoop), rj)
					fmt.Println(2, " ", i, " ", r, " ", j)
				}
				oldPath := regexMap[i][j]
				regexMap[i][j] = unionRegex(oldPath, newPath)
				fmt.Println("oldPath:", oldPath, " newPath", newPath, "union:", regexMap[i][j])
			}
		}

		eliminatedStates[r] = true
	}

	// Step 5: 结果在 initial -> accept 之间
	result := regexMap[a.InitialState][finalState]
	if result == "∅" {
		return model.Regex("∅")
	}
	if result == "ε" {
		return model.Regex("ε")
	}
	return model.Regex(result)
}
