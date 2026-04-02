# 自动机等价性检查

## 概述
判断两个自动机是否接受相同的语言。

## 算法思路

### 方法1: 最小化 + 同构检查 (推荐)
```
1. 将两个自动机都转换为 DFA
2. 对两个 DFA 分别进行最小化
3. 检查两个最小 DFA 是否同构
```

**优点**: 理论上完全准确
**时间复杂度**: O(n^3) (主要在最小化)

### 方法2: 对称差检查
```
1. 构造 L1 ∩ ¬L2 的自动机
2. 检查是否接受空语言
3. 同理检查 L2 ∩ ¬L1
```

**优点**: 可以找到反例
**时间复杂度**: O(n^2)

## 算法步骤 (方法1)

### 1. 转换为 DFA
```
if a1.Type != DFA:
    a1 = NFAToDFA(a1)
if a2.Type != DFA:
    a2 = NFAToDFA(a2)
```

### 2. 最小化
```
minDFA1 = DFAMinimize(a1)
minDFA2 = DFAMinimize(a2)
```

### 3. 规范化
对两个最小 DFA 进行状态重命名,使状态名标准化。

### 4. 同构检查
比较:
- 状态数
- 初始状态
- 接受状态集合
- 转移函数

## 规范化算法 (canonicalize)
```go
func canonicalize(dfa *Automaton) *Automaton
```
- BFS 遍历状态
- 按访问顺序重新命名: q0, q1, q2, ...
- 保持转移结构不变

## 同构检查
```go
func AreDFAsIsomorphic(dfa1, dfa2 *Automaton) bool
```
比较规范化后的两个 DFA 是否完全相同。

## 时间复杂度
- 转 DFA: O(2^n)
- 最小化: O(n^2)
- 同构检查: O(n)
- 总计: O(2^n) (最坏情况)

## 空间复杂度
- O(2^n) (存储 DFA)

## 应用场景
- 正则表达式等价性检查
- 自动机优化验证
- 教学演示

## 参考实现
- 文件: `automaton_s/equivalence_check.go`
- 函数: `AutomatonEquivalenceCheck()`, `AreDFAsIsomorphic()`

## 示例
```
自动机1: (a|b)*a
自动机2: (a*a|b*a)*

转换为最小 DFA 后:
  状态数: 相同
  初始状态: 相同
  接受状态: 相同
  转移: 相同
结果: 等价 ✓
```
