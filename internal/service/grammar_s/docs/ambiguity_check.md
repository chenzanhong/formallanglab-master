# 文法二义性检查

## 概述
检查正则文法是否存在二义性。

## 二义性定义
存在至少一个句子可以通过两种或以上不同的最左推导得到。

## 算法思路 (BFS)

### 1. BFS 生成句子
- 从起始符号开始 BFS
- 限制搜索深度防止无限循环
- 记录每个句子被推导出的次数

### 2. 检测二义性
- 一旦发现某个句子被推导出超过一次,立即返回 true
- 若遍历完所有可能推导仍未发现重复句子,返回 false

## 算法步骤

### Step 1: 初始化
```
maxDepth = 2 * len(NonTerminals) (最大不超过 10)
sentenceCount = {}
visited = {}
```

### Step 2: BFS 遍历
```
queue = [(StartSymbol, 0)]

while queue not empty:
    curr = queue.pop()
    
    if curr 是句子:
        sentenceCount[sentence]++
        if sentenceCount[sentence] > 1:
            return true
    
    for each production:
        new_state = apply production
        if not visited:
            visited.add(new_state)
            queue.push(new_state)
```

### Step 3: 返回结果
```
return false
```

## 时间复杂度
- O(k * n^d) (k 为产生式数量, n 为非终结符数量, d 为搜索深度)

## 空间复杂度
- O(m) (m 为生成的不同句子数量)

## 应用场景
- 文法验证
- 编译器设计
- 形式语言理论

## 参考实现
- 文件: `grammar_s/ambiguity_check.go`
- 函数: `IsAmbiguousRegular()`

## 示例
```
二义文法: S → aS | S | ε
句子 "a" 可以通过两种方式推导:
  1. S ⇒ aS ⇒ a
  2. S ⇒ S ⇒ aS ⇒ a

结果: 二义 ✓
```
