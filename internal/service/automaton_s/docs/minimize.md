# DFA 最小化 - Hopcroft 算法

## 算法概述
将一个 DFA 转换为具有最少状态数的等价 DFA。

## 核心思想
通过迭代划分状态集合,将不可区分的状态合并,直到划分不再变化。

## 算法步骤

### 1. 移除不可达状态
从初始状态出发,通过 BFS 或 DFS 找到所有可达状态,移除不可达状态。

### 2. 初始划分
将状态分为两类:
- 接受状态集合 F
- 非接受状态集合 Q-F

### 3. 迭代划分
对于每个划分块 A 和每个输入符号 a:
1. 找到所有能通过 a 转移到 A 中状态的前驱状态集合 X
2. 对每个划分块 Y:
   - 如果 Y ∩ X ≠ ∅ 且 Y - X ≠ ∅,则分裂 Y 为 Y ∩ X 和 Y - X
3. 更新划分 P

### 4. 构建最小 DFA
- 每个划分块作为一个新状态
- 初始状态: 包含原初始状态的块
- 接受状态: 包含原接受状态的块
- 转移: 如果原状态 q ∈ block1 通过 a 转移到 block2,则 block1 --a--> block2

## 时间复杂度
- O(n * k * log(n)) (n 为状态数, k 为字母表大小)

## 空间复杂度
- O(n) (存储划分和反向转移图)

## 关键数据结构
- **Partition**: 当前的状态划分
- **Reverse Transitions**: 反向转移图,用于快速查找前驱
- **Worklist W**: 待处理的划分块

## 应用场景
- 压缩自动机存储空间
- 提高匹配效率
- 自动机等价性检查

## 参考实现
- 文件: `automaton_s/minimize.go`
- 函数: `DFAMinimize()`, `minimizeByHopcroft()`

## 示例
输入 DFA:
```
States: {q0, q1, q2, q3, q4}
Alphabet: {a, b}
Transitions:
  q0 --a--> q1
  q0 --b--> q2
  q1 --a--> q3
  q1 --b--> q4
  q2 --a--> q3
  q2 --b--> q4
  q3 --a--> q3
  q3 --b--> q3
  q4 --a--> q4
  q4 --b--> q4
Initial: q0
Accepting: {q3, q4}
```

输出最小 DFA:
```
States: {[q0], [q1,q2], [q3], [q4]}
Alphabet: {a, b}
Transitions:
  [q0] --a--> [q1,q2]
  [q0] --b--> [q1,q2]
  [q1,q2] --a--> [q3]
  [q1,q2] --b--> [q4]
  [q3] --a--> [q3]
  [q3] --b--> [q3]
  [q4] --a--> [q4]
  [q4] --b--> [q4]
Initial: [q0]
Accepting: {[q3], [q4]}
```
