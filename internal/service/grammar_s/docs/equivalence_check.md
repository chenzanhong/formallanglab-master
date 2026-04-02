# 文法等价性检查

## 概述
判断两个正则文法是否等价(生成相同的语言)。

## 算法思路 (有限测试集)

### 1. 生成语言 (有限深度)
- 通过 BFS 生成两个文法在有限深度内能推导出的所有句子
- 限制最大深度防止无限生成

### 2. 比较语言集合
- 比较两个文法生成的语言集合是否完全相同

## 算法步骤

### Step 1: 生成文法语言
```
language1 = generateLanguage(g1, maxDepth=8)
language2 = generateLanguage(g2, maxDepth=8)
```

### Step 2: 比较集合
```
return language1 == language2
```

## 时间复杂度
- O(k * n^d) (k 为产生式数量, n 为非终结符数量, d 为搜索深度)

## 空间复杂度
- O(m) (m 为生成的不同句子数量)

## 注意事项
- 由于文法可能生成无限语言,这里只比较有限深度内的句子
- 限制:最大深度设为 8,可能会有误判
- 理论上应该转为最小 DFA 判断是否同构

## 应用场景
- 文法等价性验证
- 文法优化
- 教学演示

## 参考实现
- 文件: `grammar_s/equivalence_check.go`
- 函数: `IsEquivalent()`, `generateLanguage()`, `compareSets()`

## 示例
```
文法1: S → aS | ε
文法2: S → S a | ε

两个文法都生成语言 {ε, a, aa, aaa, ...}

结果: 等价 ✓
```
