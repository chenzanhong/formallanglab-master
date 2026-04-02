# DFA 完备化

## 算法概述
将一个非完备的 DFA 转换为完备的 DFA,确保每个状态对每个输入符号都有确定的转移。

## 完备化定义
DFA 是完备的,当且仅当:
```
∀q ∈ Q, ∀a ∈ Σ, ∃! q' ∈ Q, δ(q, a) = q'
```

## 算法步骤

### 1. 检查完备性
遍历所有状态和输入符号,检查是否都有转移。

### 2. 添加陷阱状态
如果存在缺失转移,添加一个陷阱状态("dead_state" 或 "_sink_")。

### 3. 补充转移
对每个缺失的转移 (q, a),添加转移 q --a--> sink。

### 4. 陷阱状态自环
陷阱状态对所有输入符号自环:
```
sink --a--> sink, ∀a ∈ Σ
```

## 时间复杂度
- O(n * k) (n 为状态数, k 为字母表大小)

## 空间复杂度
- O(1) (只添加一个陷阱状态)

## 应用场景
- DFA 最小化前的预处理
- 正则表达式匹配
- 编译器语法分析

## 注意事项
- 陷阱状态不是接受状态
- 陷阱状态不可达性不影响语言接受
- 完备化不改变 DFA 的语言

## 参考实现
- 文件: `automaton_s/complete.go`
- 函数: `CompleteDFA()`

## 示例
输入非完备 DFA:
```
States: {q0, q1, q2}
Alphabet: {a, b}
Transitions:
  q0 --a--> q1
  q0 --b--> q2
  q1 --a--> q0
  // q1 缺失 b 转移
  // q2 缺失 a 和 b 转移
Initial: q0
Accepting: {q0}
```

输出完备 DFA:
```
States: {q0, q1, q2, _sink_}
Alphabet: {a, b}
Transitions:
  q0 --a--> q1
  q0 --b--> q2
  q1 --a--> q0
  q1 --b--> _sink_
  q2 --a--> _sink_
  q2 --b--> _sink_
  _sink_ --a--> _sink_
  _sink_ --b--> _sink_
Initial: q0
Accepting: {q0}
```
