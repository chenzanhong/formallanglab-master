// grammar_s/first_follow.go

package grammar_s

import (
    "backend/internal/domain/model"
    "fmt"
    "strings"
)


// CalculateFirst 计算每个符号的 FIRST 集
func CalculateFirst(grammar *model.Grammar) map[model.Symbol]map[model.Symbol]struct{} {
    firstSet := make(map[model.Symbol]map[model.Symbol]struct{})

    // 初始化所有终结符的 FIRST 集
    for _, terminal := range grammar.Terminals {
        firstSet[terminal] = map[model.Symbol]struct{}{terminal: {}}
    }

    // 初始化所有非终结符的 FIRST 集为空
    for _, nonTerminal := range grammar.NonTerminals {
        firstSet[nonTerminal] = make(map[model.Symbol]struct{})
    }

    changed := true
    for changed {
        changed = false
        for _, production := range grammar.Productions {
            left := production.Left[0]
            right := production.Right

            // 处理右部符号序列
            for _, symbol := range right {
                if _, exists := firstSet[symbol]; !exists {
                    continue
                }
                for t := range firstSet[symbol] {
                    if t != model.Epsilon {
                        if _, in := firstSet[left][t]; !in {
                            firstSet[left][t] = struct{}{}
                            changed = true
                        }
                    }
                }
                // 如果 symbol 不能推出 ε，则停止
                if _, hasEpsilon := firstSet[symbol][model.Epsilon]; !hasEpsilon {
                    break
                }
            }

            // 检查整个右部是否都能推出 ε
            allCanDeriveEpsilon := true
            for _, symbol := range right {
                if _, can := firstSet[symbol][model.Epsilon]; !can {
                    allCanDeriveEpsilon = false
                    break
                }
            }
            if allCanDeriveEpsilon {
                if _, has := firstSet[left][model.Epsilon]; !has {
                    firstSet[left][model.Epsilon] = struct{}{}
                    changed = true
                }
            }
        }
    }

    return firstSet
}

// PrintFirst 打印 FIRST 集
func PrintFirst(firstSet map[model.Symbol]map[model.Symbol]struct{}) {
    for symbol, first := range firstSet {
        items := []string{}
        for s := range first {
            if s == model.Epsilon {
                items = append(items, "ε")
            } else {
                items = append(items, string(s))
            }
        }
        fmt.Printf("FIRST(%s) = {%s}\n", symbol, strings.Join(items, ", "))
    }
}

// CalculateFollow 计算 FOLLOW 集
func CalculateFollow(grammar *model.Grammar, firstSet map[model.Symbol]map[model.Symbol]struct{}) map[model.Symbol]map[model.Symbol]struct{} {
    followSet := make(map[model.Symbol]map[model.Symbol]struct{})

    // 初始化
    for _,nonTerminal := range grammar.NonTerminals {
        followSet[nonTerminal] = make(map[model.Symbol]struct{})
    }

    // 起始符号的 FOLLOW 包含 #
    followSet[grammar.StartSymbol]["#"] = struct{}{}

    changed := true
    for changed {
        changed = false
        for _, production := range grammar.Productions {
            left := production.Left[0]
            right := production.Right

            for i := 0; i < len(right); i++ {
                symbol := right[i]
                if !grammar.CheckIsNonTerminal(symbol) {
                    continue
                }

                // 后继符号：right[i+1:]
                if i < len(right)-1 {
                    nextSymbol := right[i+1]
                    for t := range firstSet[nextSymbol] {
                        if t != model.Epsilon {
                            if _, in := followSet[symbol][t]; !in {
                                followSet[symbol][t] = struct{}{}
                                changed = true
                            }
                        }
                    }
                    if _, hasEpsilon := firstSet[nextSymbol][model.Epsilon]; hasEpsilon {
                        for t := range followSet[left] {
                            if _, in := followSet[symbol][t]; !in {
                                followSet[symbol][t] = struct{}{}
                                changed = true
                            }
                        }
                    }
                } else {
                    // 是最后一个符号
                    for t := range followSet[left] {
                        if _, in := followSet[symbol][t]; !in {
                            followSet[symbol][t] = struct{}{}
                            changed = true
                        }
                    }
                }
            }
        }
    }

    return followSet
}

// PrintFollow 打印 FOLLOW 集
func PrintFollow(followSet map[model.Symbol]map[model.Symbol]struct{}) {
    for symbol, follow := range followSet {
        items := []string{}
        for s := range follow {
            items = append(items, string(s))
        }
        fmt.Printf("FOLLOW(%s) = {%s}\n", symbol, strings.Join(items, ", "))
    }
}


// symbolsEqual 比较两个符号切片是否相等
func symbolsEqual(a, b []model.Symbol) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

// canDeriveEpsilon 判断非终结符是否可推出 ε
func canDeriveEpsilon(grammar *model.Grammar) bool {
    firstSet := CalculateFirst(grammar)
    start := grammar.StartSymbol
    _, can := firstSet[start][model.Epsilon]
    return can
}

// RecognizeString 判断字符串是否被文法生成（BFS 模拟推导）
func RecognizeString(g *model.Grammar, input []model.Symbol, maxSteps, maxWidth int) (bool, []string, error) {
    if len(input) == 0 {
        // 处理空串
        if canDeriveEpsilon(g) {
            for _, prod := range g.Productions {
                if prod.Left[0] == g.StartSymbol && len(prod.Right) == 1 && prod.Right[0] == model.Epsilon {
                    return true, []string{fmt.Sprintf("%s → ε", string(g.StartSymbol))}, nil
                }
            }
            return true, []string{"S → ε"}, nil
        }
        return false, nil, nil
    }

    // BFS 队列：记录当前符号串和推导路径
    type state struct {
        symbols []model.Symbol
        steps   int
        path    []string // 推导历史，如 "S → aA"
    }

    queue := []state{
        {symbols: []model.Symbol{g.StartSymbol}, steps: 0, path: []string{string(g.StartSymbol)}},
    }

    visited := make(map[string]bool)

    for len(queue) > 0 {
        curr := queue[0]
        queue = queue[1:]

        // 去重
        key := symbolsToString(curr.symbols)
        if visited[key] || curr.steps >= maxSteps {
            continue
        }
        visited[key] = true

        // 检查是否匹配输入
        if symbolsEqual(curr.symbols, input) {
            return true, curr.path, nil
        }

        // 剪枝：长度超过输入太多
        if len(curr.symbols) > len(input)+10 {
            continue
        }

        // 尝试应用每个产生式
        for _, prod := range g.Productions {
            leftSymbol := prod.Left[0]
            for i := 0; i < len(curr.symbols); i++ {
                if curr.symbols[i] == leftSymbol {
                    // 替换第 i 个符号
                    newSymbols := append([]model.Symbol{}, curr.symbols[:i]...)
                    newSymbols = append(newSymbols, prod.Right...)
                    newSymbols = append(newSymbols, curr.symbols[i+1:]...)

                    if len(newSymbols) > maxWidth {
                        continue
                    }

                    // 构建新推导路径
                    step := fmt.Sprintf("%s → %s", string(leftSymbol), symbolsToString(prod.Right))
                    newPath := append([]string{}, curr.path...)
                    newPath = append(newPath, step)

                    queue = append(queue, state{
                        symbols: newSymbols,
                        steps:   curr.steps + 1,
                        path:    newPath,
                    })
                }
            }
        }
    }

    return false, nil, nil
}

// IsLL1 判断文法是否为 LL(1)
func IsLL1(grammar *model.Grammar) (bool, string) {
    firstSet := CalculateFirst(grammar)
    followSet := CalculateFollow(grammar, firstSet)

    table := make(map[model.Symbol]map[model.Symbol]int)

    for _, nonTerminal := range grammar.NonTerminals {
        table[nonTerminal] = make(map[model.Symbol]int)
    }

    for i, prod := range grammar.Productions {
        left := prod.Left[0]
        right := prod.Right

        var firstOfRight map[model.Symbol]struct{}
        if len(right) > 0 {
            firstOfRight = firstSet[right[0]]
        } else {
            firstOfRight = map[model.Symbol]struct{}{model.Epsilon: {}}
        }

        for t := range firstOfRight {
            if t != model.Epsilon {
                if other, conflict := table[left][t]; conflict {
                    return false, fmt.Sprintf("冲突: 产生式 %d 和 %d 在 %s/%s", other, i, string(left), string(t))
                }
                table[left][t] = i
            }
        }

        if _, hasEpsilon := firstOfRight[model.Epsilon]; hasEpsilon {
            for t := range followSet[left] {
                if other, conflict := table[left][t]; conflict {
                    return false, fmt.Sprintf("FOLLOW 冲突: 产生式 %d 和 %d 在 %s/%s", other, i, string(left), string(t))
                }
                table[left][t] = i
            }
        }
    }

    return true, ""
}

// LL1Parse 使用 LL(1) 分析表解析输入串
func LL1Parse(grammar *model.Grammar, input []model.Symbol) (bool, []string, error) {
    firstSet := CalculateFirst(grammar)
    followSet := CalculateFollow(grammar, firstSet)

    // 构建预测分析表
    table := make(map[model.Symbol]map[model.Symbol]int)
    for _, nt := range grammar.NonTerminals {
        table[nt] = make(map[model.Symbol]int)
    }

    for i, prod := range grammar.Productions {
        left := prod.Left[0]
        right := prod.Right

        firstOfRight := map[model.Symbol]struct{}{}
        if len(right) > 0 {
            for k := range firstSet[right[0]] {
                firstOfRight[k] = struct{}{}
            }
        } else {
            firstOfRight[model.Epsilon] = struct{}{}
        }

        for t := range firstOfRight {
            if t != model.Epsilon {
                if _, exists := table[left][t]; exists {
                    return false, nil, fmt.Errorf("LL(1) 冲突: %s → ... 和另一产生式在 %s", string(left), string(t))
                }
                table[left][t] = i
            }
        }

        if _, hasEpsilon := firstOfRight[model.Epsilon]; hasEpsilon {
            for t := range followSet[left] {
                if _, exists := table[left][t]; exists {
                    return false, nil, fmt.Errorf("LL(1) FOLLOW 冲突: %s → ... 和另一产生式在 %s", string(left), string(t))
                }
                table[left][t] = i
            }
        }
    }

    // 模拟分析栈
    stack := []model.Symbol{"#"}
    stack = append(stack, grammar.StartSymbol)

    input = append(input, "#")
    ip := 0
    var steps []string

    for len(stack) > 0 {
        top := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if ip >= len(input) {
            return false, steps, nil
        }
        current := input[ip]

        steps = append(steps, fmt.Sprintf("栈: [%s], 输入: %s, 指针: %d", 
            symbolsToString(stack), ",", 
            symbolsToString(input[ip:]), ip))

        if top == current && top == "#" {
            steps = append(steps, "分析成功: 输入被接受")
            return true, steps, nil
        }

        if top == current {
            ip++
            continue
        }

        if grammar.CheckIsTerminal(top) {
            steps = append(steps, fmt.Sprintf("匹配失败: 栈顶 %s ≠ 输入 %s", string(top), string(current)))
            return false, steps, nil
        }

        if productionIndex, ok := table[top][current]; ok {
            prod := grammar.Productions[productionIndex]
            // 反向压栈
            for i := len(prod.Right) - 1; i >= 0; i-- {
                if prod.Right[i] != model.Epsilon {
                    stack = append(stack, prod.Right[i])
                }
            }
            steps = append(steps, fmt.Sprintf("使用产生式: %s → %s", 
                string(top), symbolsToString(prod.Right)))
        } else {
            steps = append(steps, fmt.Sprintf("无匹配产生式: %s / %s", string(top), string(current)))
            return false, steps, nil
        }
    }

    return false, steps, nil
}