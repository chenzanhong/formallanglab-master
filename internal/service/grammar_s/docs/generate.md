# 文法示例字符串生成

## 概述
为文法生成可推导和不可推导的字符串示例。

## 算法思路

### 可推导字符串 (BFS 推导)
```
1. 从起始符号开始 BFS
2. 模拟所有可能的推导过程
3. 收集最多 k 个最短的句子

### 不可推导字符串 (枚举+验证)
```
1. 生成候选字符串 (长度 0 ~ maxLen)
2. 对每个字符串,使用识别算法验证
3. 收集最多 k 个不可推导的字符串

## 算法步骤

### Step 1: 生成可推导字符串
```
queue = [(StartSymbol, depth=0)]
seen = {}

while queue not empty and len(accept) < maxExampleNum:
    curr = queue.pop()
    
    if curr 是句子:
        accept.add(curr)
        continue
    
    if depth > maxDepth:
        continue
    
    for each production:
        new_form = apply production
        if not seen:
            seen.add(new_form)
            queue.push(new_form)
```

### Step 2: 生成不可推导字符串
```
for length in 0 to maxLen:
    for each string of this length:
        if not recognized:
            reject.add(string)
            if len(reject) >= maxExampleNum:
                break
```

## 时间复杂度
- BFS: O(n^m) (指数级)
- 枚举: O(k^m) (k 为字母表大小, m 为最大长度)

## 空间复杂度
- O(n^m) (存储推导状态)

## 应用场景
- 教学演示
- 单元测试
- 用户示例

## 参考实现
- 文件: `grammar_s/generate.go`
- 函数: `GrammarGenerateExampleString()`, `generateAcceptExampleString()`, `generateRejectExampleString()`

## 示例
```
文法: S → aS | b

可推导字符串: ["b", "ab", "aab", "aaab", ...]
不可推导字符串: ["", "a", "aa", "ba", "aba", ...]
```
