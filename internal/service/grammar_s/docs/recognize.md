# 文法识别与分析

## 概述
判断字符串是否能被文法生成,并提供多种分析算法。

## 分析算法

### 1. BFS 推导模拟
- 从起始符号开始 BFS
- 模拟所有可能的推导过程
- 找到匹配的推导路径

### 2. LL(1) 分析
- 预测分析
- 需要 FIRST 和 FOLLOW 集
- 线性时间分析

### 3. LR(0) 分析
- 自底向上分析
- 构造 LR(0) 项目集族
- 构建分析表

### 4. 递归下降分析
- 从上而下分析
- 为每个非终结符编写递归函数
- 需要消除左递归和提取左公因子

## FIRST/FOLLOW 集计算

### FIRST 集
```
FIRST(α) = { a ∈ Σ | α ⇒* aβ } ∪ { ε | α ⇒* ε }
```

### FOLLOW 集
```
FOLLOW(A) = { a ∈ Σ | S ⇒* αAaβ } ∪ { # | S ⇒* αA }
```

## 时间复杂度
- BFS: O(n^m) (指数级)
- LL(1): O(n)
- LR(0): O(n^2)

## 空间复杂度
- BFS: O(n^m)
- LL(1): O(n)
- LR(0): O(n^2)

## 应用场景
- 字符串识别
- 语法分析
- 编译器实现

## 参考实现
- 文件: `grammar_s/recognize.go`
- 函数: `RecognizeString()`, `CalculateFirst()`, `CalculateFollow()`, `LL1Parse()`, `LR0Parse()`
