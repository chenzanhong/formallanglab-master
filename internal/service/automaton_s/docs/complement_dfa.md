# DFA 补集

## 概述
给定一个 DFA,构造接受其补集语言的 DFA。

## 理论基础

### 语言补集
```
¬L(M) = {w ∈ Σ* | w ∉ L(M)}
```

### 构造方法
1. 确保 DFA 是完备的
2. 翻转接受状态: 接受 ↔ 非接受

## 算法步骤

### Step 1: 完备化
```
if not complete:
    添加陷阱状态
    补充缺失转移
```

### Step 2: 翻转接受状态
```
newAccepting = States - AcceptingStates
```

### Step 3: 验证
检查新自动机的有效性。

## 正确性证明

### 完备性保证
完备 DFA 确保对任何输入字符串都有确定的接受/拒绝决策。

### 翻转正确性
```
w ∈ L(M)  ⇔  M 接受 w
w ∈ ¬L(M) ⇔  M 不接受 w ⇔ M' 接受 w
```

其中 M' 是翻转接受状态后的 DFA。

## 时间复杂度
- 完备化: O(n * k)
- 翻转: O(n)
- 总计: O(n * k)

## 空间复杂度
- O(1) (只添加一个陷阱状态)

## 应用场景
- 正则表达式差集
- 语言包含性检查
- 自动机运算

## 注意事项
- 必须先完备化,否则翻转不正确
- 陷阱状态永远不是接受状态

## 参考实现
- 文件: `automaton_s/complement_dfa.go`
- 函数: `ComplementDFA()`

## 示例
```
DFA M: 匹配偶数个 a
States: {q0, q1}
Alphabet: {a}
Transitions:
  q0 --a--> q1
  q1 --a--> q0
Initial: q0
Accepting: {q0}

L(M) = {ε, aa, aaaa, ...}

补集 DFA M':
Accepting: {q1}

L(M') = {a, aaa, aaaaa, ...}
```

## 扩展: 其他集合运算

### 交集 (L1 ∩ L2)
```
构造乘积自动机
接受状态: (q1, q2) 其中 q1 ∈ F1 且 q2 ∈ F2
```

### 并集 (L1 ∪ L2)
```
接受状态: (q1, q2) 其中 q1 ∈ F1 或 q2 ∈ F2
```

### 差集 (L1 - L2)
```
接受状态: (q1, q2) 其中 q1 ∈ F1 且 q2 ∉ F2
```
