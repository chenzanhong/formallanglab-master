# 正则表达式等价性检查

## 概述
判断两个正则表达式是否等价(匹配相同的字符串集合)。

## 算法思路 (有限测试集)

### 1. 提取字母表
- 从两个正则表达式中提取所有字符

### 2. 生成测试字符串
- 生成所有长度 0 ~ maxLen 的字符串

### 3. 分别匹配
- 对每个字符串,分别用两个正则表达式匹配
- 如果结果不同,则不等价

### 4. 全部相同
- 如果所有测试字符串结果都相同,认为等价

## 算法步骤

### Step 1: 提取字母表
```
alphabet1 = extractAlphabet(r1)
alphabet2 = extractAlphabet(r2)
alphabet = alphabet1 ∪ alphabet2
```

### Step 2: 生成测试字符串
```
candidates = generateStrings(alphabet, maxLen)
```

### Step 3: 分别匹配
```
for each candidate:
    m1 = match(r1, candidate)
    m2 = match(r2, candidate)
    if m1 != m2:
        return false, candidate  // 反例
```

### Step 4: 返回结果
```
return true
```

## 时间复杂度
- O(k^m * n) (k 为字母表大小, m 为最大长度, n 为正则表达式长度)

## 空间复杂度
- O(k^m) (存储测试字符串)

## 注意事项
- 这是一种近似方法,理论上不完备
- 实践中对于大多数正则表达式是可靠的
- 有极小概率漏判(不等价的正则表达式在测试集中表现相同)

## 应用场景
- 正则表达式等价性验证
- 正则表达式优化
- 教学演示

## 参考实现
- 文件: `regex_s/equivalence_check.go`
- 函数: `RegexEquivalenceCheck()`, `extractAlphabet()`

## 示例
```
正则表达式1: (a|b)*a
正则表达式2: b*ab*a | a

测试字符串: ["", "a", "b", "aa", "ab", "ba", "bb", ...]

结果: 等价 ✓
```
