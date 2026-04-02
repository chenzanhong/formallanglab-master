# 正则表达式示例字符串生成

## 概述
为正则表达式生成可匹配和不可匹配的字符串示例。

## 算法思路

### 可匹配字符串 (转 DFA + BFS)

1. 正则表达式 → NFA (Thompson 构造法)
2. NFA → DFA (子集构造法)
3. DFA 完备化
4. 在 DFA 上 BFS 生成最短匹配字符串

### 不可匹配字符串 (补集 DFA + BFS)

1. 构造补集 DFA (翻转接受状态)
2. 在补集 DFA 上 BFS
3. 生成最短的"匹配"字符串(即原正则的不匹配字符串)

## 算法步骤

### Step 1: 转换为 DFA
```
if regex == ∅:
    return [], []
    
nfa = RegexToFA(regex)
dfa = NFAToDFA(nfa)
dfa.CompleteDFA()
```

### Step 2: 生成匹配字符串
```go
accept = BFSShortestAcceptedStringsForDFA(dfa, k)
```

### Step 3: 生成不匹配字符串
```go
complement = complementDFA(dfa)
reject = BFSShortestAcceptedStringsForDFA(complement, k)
```

### Step 4: 随机采样
```go
if len(accept) > maxExampleNum:
    accept = sampling(accept, maxExampleNum)
if len(reject) > maxExampleNum:
    reject = sampling(reject, maxExampleNum)
```

## 时间复杂度
- 正则 → NFA: O(n)
- NFA → DFA: O(2^n)
- BFS: O(m)
- 总计: O(2^n) (最坏情况)

## 空间复杂度
- O(2^n) (DFA 状态数)

## 应用场景
- 教学演示
- 单元测试
- 用户示例

## 参考实现
- 文件: `regex_s/generate.go`
- 函数: `RegexGenerateExampleString()`, `regexGenerateExampleStringByCompletedDFAAndBFS()`

## 示例
```
正则表达式: a*b+

可匹配字符串: ["b", "ab", "aab", "bb", "abb", "aabb"]
不可匹配字符串: ["", "a", "aa", "ba", "aba", "aaba"]
```
