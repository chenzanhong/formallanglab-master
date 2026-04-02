# 生成示例字符串

## 概述
为自动机生成可接受和不可接受的字符串示例。

## 算法思路

### 可接受字符串 (BFS 最短路径)
```
1. 从初始状态开始 BFS
2. 记录到达每个状态的最短字符串
3. 收集最多 k 个最短的接受字符串
```

### 不可接受字符串
```
1. 构造补集 DFA
2. 在补集 DFA 上 BFS
3. 收集最多 k 个最短的接受字符串(即原 DFA 的拒绝字符串)
```

## 算法步骤

### 1. DFA 转换
```
if automaton.Type != DFA:
    automaton = NFAToDFA(automaton)
```

### 2. 完备化
```
automaton.CompleteDFA()
```

### 3. BFS 生成接受字符串
```
队列: (状态, 已构建字符串)
visited: (状态, 字符串长度)

1. 如果初始状态是接受状态,添加空字符串
2. BFS 遍历,记录路径
3. 收集最多 k 个最短接受字符串
4. 限制最大长度: max(len(states), min(len(states)*2, 12))
```

### 4. 生成拒绝字符串
```
1. 构造补集 DFA (翻转接受状态)
2. 在补集 DFA 上 BFS
3. 收集最多 k 个最短字符串
```

### 5. 随机采样
如果生成的字符串过多,进行随机采样。

## 时间复杂度
- BFS: O(n + m)
- 采样: O(k)

## 空间复杂度
- O(n * maxLen) (存储路径)

## 应用场景
- 教学演示
- 单元测试
- 用户示例

## 参考实现
- 文件: `automaton_s/generate.go`
- 函数: `AutomatonGenerateExampleString()`, `BFSShortestAcceptedStringsForDFA()`

## 示例
```
DFA: 匹配 a*b+
States: {q0, q1}
Alphabet: {a, b}
Transitions:
  q0 --a--> q0
  q0 --b--> q1
  q1 --b--> q1
Initial: q0
Accepting: {q1}

接受字符串: ["b", "ab", "aab", "bb", "abb", "aabb"]
拒绝字符串: ["", "a", "aa", "ba", "aba", "aaba"]
```

## 优化策略

### 1. 限制深度
避免生成过长字符串,限制最大长度。

### 2. 去重
使用 visited 集合避免重复字符串。

### 3. 随机采样
当生成字符串过多时,随机选择 k 个作为示例。

### 4. 兜底策略
如果 BFS 失败,回退到枚举验证方法。
