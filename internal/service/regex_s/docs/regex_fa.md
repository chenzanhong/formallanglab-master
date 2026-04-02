# 正则表达式转自动机 - Thompson 构造法

## 概述
将正则表达式转换为等价的 NFA。

## 算法: Thompson 构造法

### 核心思想
递归地将正则表达式分解为基本操作,为每个操作构造小的 NFA,然后组合。

## 支持的运算

### 1. 基本符号
```
a → (start) --a--> (end)
```

### 2. 连接 (Concatenation)
```
r1 · r2 → (start_r1) --ε--> (start_r2) --ε--> (end_r2)
           (end_r1) --ε--> (start_r2)
```

### 3. 选择 (Union)
```
r1 | r2 → (start) --ε--> (start_r1)
          (start) --ε--> (start_r2)
          (end_r1) --ε--> (end)
          (end_r2) --ε--> (end)
```

### 4. Kleene 闭包 (Star)
```
r* → (start) --ε--> (start_r)
     (start) --ε--> (end)
     (end_r) --ε--> (start_r)
     (end_r) --ε--> (end)
```

### 5. 正闭包 (Plus)
```
r+ → r · r* (转换为连接和闭包)
```

### 6. 可选 (Question)
```
r? → r | ε (转换为选择)
```

## 算法步骤

### 1. 词法分析
```
tokens = lex(pattern)
```

### 2. 语法分析 (递归下降)
```
ast = parse(pattern)
```

### 3. Thompson 构造
```
(start, end) = build(ast)
```

### 4. 构建 NFA
```
nfa = {
    States: states,
    Alphabet: alphabet,
    Transitions: transitions,
    InitialState: start,
    AcceptingStates: [end],
    Type: EpsilonNFA
}
```

## 时间复杂度
- O(n) (n 为正则表达式长度)

## 空间复杂度
- O(n) (状态和转移数量与正则表达式长度成正比)

## 注意事项
- 生成的 NFA 是 ε-NFA (允许 ε-转移)
- 状态数和转移数与正则表达式长度成线性关系
- 结果不唯一(取决于构造顺序)

## 应用场景
- 正则表达式引擎
- 字符串匹配
- 编译器词法分析

## 参考实现
- 文件: `regex_s/regex_fa.go`
- 函数: `RegexToFA()`, `RegexToFAWithSteps()`, `parseRegex()`, `build()`

## 示例
```
正则表达式: (a|b)*c

Thompson 构造:
1. a → (q0) --a--> (q1)
2. b → (q2) --b--> (q3)
3. a|b → (q4) --ε--> (q0), (q4) --ε--> (q2), (q1) --ε--> (q5), (q3) --ε--> (q5)
4. (a|b)* → (q6) --ε--> (q4), (q6) --ε--> (q7), (q5) --ε--> (q4), (q5) --ε--> (q7)
5. c → (q8) --c--> (q9)
6. (a|b)*c → (q6) --ε--> (q8), (q7) --ε--> (q9)

最终 NFA:
  InitialState: q6
  AcceptingStates: [q9]
```
