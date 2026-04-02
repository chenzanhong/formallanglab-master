# 自动机转正则表达式 - 状态消除法

## 概述
将 NFA/DFA 转换为等价的正则表达式。

## 算法: 状态消除法 (State Elimination)

### 核心思想
逐步消除中间状态,维护状态间的正则表达式,最终得到初态到终态的正则表达式。

## 算法步骤

### 1. 预处理
- 添加唯一初态 initial
- 添加唯一终态 final
- initial --ε--> 原初态
- 原终态 --ε--> final

### 2. 构建邻接矩阵
```
regexMap[i][j] = 从状态 i 到状态 j 的正则表达式
```

初始化:
- 如果 i 有自环,regexMap[i][i] = ∅ (空语言)
- 如果 i --a--> j,regexMap[i][j] = a
- 否则 regexMap[i][j] = ∅

### 3. 逐个消除状态
对于每个要消除的状态 r:

```
对每对状态 (i, j):
    regexMap[i][j] = regexMap[i][j] | 
                     regexMap[i][r] · regexMap[r][r]* · regexMap[r][j]
```

其中:
- `regexMap[r][r]*` 是自环的 Kleene 闭包
- `·` 是正则表达式连接
- `|` 是正则表达式并(选择)

### 4. 结果
最终 `regexMap[initial][final]` 即为所求正则表达式。

## 正则表达式运算

### 连接 (concat)
```
r1 · r2
```

### 并 (union)
```
r1 | r2
```

### Kleene 闭包 (star)
```
r*
```

## 时间复杂度
- O(n! * k) (n 为状态数, k 为字母表大小)
- 实际表现通常更好

## 空间复杂度
- O(n^2) (存储正则表达式矩阵)

## 注意事项
- 正则表达式可能很大
- 可能包含冗余括号
- 结果不唯一

## 应用场景
- 自动机分析
- 正则表达式生成
- 教学演示

## 参考实现
- 文件: `automaton_s/fa_regex.go`
- 函数: `FAToRegex()`, `FAToRegexWithProcess()`

## 示例
```
NFA:
  q0 --a--> q1
  q1 --b--> q1
  q1 --ε--> q2

步骤1: 添加初态和终态
  initial --ε--> q0
  q2 --ε--> final

步骤2: 消除 q1
  regexMap[q0][q1] = a
  regexMap[q1][q1] = b
  regexMap[q1][q2] = ε
  
  regexMap[q0][q2] = ∅ | a · b* · ε = ab*

步骤3: 消除 q0
  regexMap[initial][q0] = ε
  regexMap[q0][q2] = ab*
  
  regexMap[initial][q2] = ∅ | ε · (ab*) = ab*

步骤4: 消除 final
  结果: ab*

正则表达式: ab*
```
