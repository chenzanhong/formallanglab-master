# 字符串识别

## 概述
判断输入字符串是否能被自动机接受。

## 算法分类

### 1. DFA 识别 (确定性)
**算法**: 直接模拟
```
1. 从初始状态开始
2. 依次读取输入符号
3. 根据转移函数转移到下一个状态
4. 读取完所有符号后,检查是否为接受状态
```

**时间复杂度**: O(n) (n 为字符串长度)
**空间复杂度**: O(1)

### 2. NFA 识别 (非确定性)
**算法**: DFS 或 BFS
```
1. 维护当前状态集合
2. 对每个输入符号,计算所有可能的转移
3. 合并所有可能的目标状态
4. 读取完所有符号后,检查是否有接受状态
```

**时间复杂度**: O(n * m) (m 为状态数)
**空间复杂度**: O(m)

### 3. ε-NFA 识别 (带 ε 转移)
**算法**: DFS + ε-闭包
```
1. 计算初始状态的 ε-闭包
2. 对每个输入符号,先计算 ε-闭包,再计算转移
3. 读取完所有符号后,检查是否能通过 ε 转移到接受状态
```

**时间复杂度**: O(n * m^2)
**空间复杂度**: O(m^2)

## 核心函数

### recognizeDFA
```go
func recognizeDFA(automaton *Automaton, str string) (*RecognitionResult, error)
```
- 输入: DFA 和字符串
- 输出: 识别结果和步骤

### recognizeNFA
```go
func recognizeNFA(automaton *Automaton, str string) (*RecognitionResult, error)
```
- 使用 DFS 遍历所有可能路径
- 记录最深路径用于失败时的反馈

### recognizeEpsilonNFA
```go
func recognizeEpsilonNFA(automaton *Automaton, str string) (*RecognitionResult, error)
```
- 处理 ε-转移
- 支持空字符串识别

## ε-闭包计算
```go
func computeEpsilonClosure(states []State, automaton *Automaton) []State
```
- 使用 BFS 或 DFS
- 避免循环

## 应用场景
- 正则表达式匹配
- 词法分析
- 模式匹配

## 参考实现
- 文件: `automaton_s/recognize.go`
- 函数: `Recognize()`, `recognizeDFA()`, `recognizeNFA()`, `recognizeEpsilonNFA()`

## 示例

### DFA 识别示例
```
DFA: 识别 a*b+
字符串: "aabbb"
步骤:
  q0 --a--> q0
  q0 --a--> q0
  q0 --b--> q1
  q1 --b--> q1
  q1 --b--> q1
结果: 接受 (在接受状态 q1)
```

### NFA 识别示例
```
NFA: (a|b)*a
字符串: "aba"
路径1: q0 --a--> q0 --b--> q0 --a--> q1 ✓
路径2: q0 --a--> q1 (失败)
结果: 接受 (存在成功路径)
```
