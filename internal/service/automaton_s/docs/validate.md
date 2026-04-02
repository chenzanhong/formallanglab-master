# 自动机验证

## 验证内容

### 1. 基本有效性检查
- 自动机对象非空
- 状态集合非空
- 初始状态在状态集合中
- 接受状态都是状态集合的元素

### 2. 输入符号合法性
- 所有转移的输入符号都在字母表中

### 3. 转移状态合法性
- 转移的起点和终点都在状态集合中

### 4. DFA 确定性检查
- 每个转移只有一个目标状态
- 不存在重复的 (状态, 输入符号) 对

## 验证规则

```
1. len(States) > 0
2. InitialState ∈ States
3. ∀acc ∈ AcceptingStates: acc ∈ States
4. ∀t ∈ Transitions: t.Input ∈ Alphabet
5. ∀t ∈ Transitions: t.FromState ∈ States ∧ t.ToStates[i] ∈ States
6. DFA: len(ToStates) == 1 ∧ 无重复 (from, input)
```

## 时间复杂度
- O(n + m) (n 为状态数, m 为转移数)

## 空间复杂度
- O(1) (不需要额外存储)

## 应用场景
- 用户输入验证
- 文件加载校验
- 算法执行前的预检查

## 参考实现
- 文件: `automaton_s/validate.go`
- 函数: `AutomatonValidate()`

## 错误类型
- `automaton cannot be nil`: 自动机为空
- `invalid initial state`: 初始状态无效
- `invalid accepting state`: 接受状态无效
- `invalid alphabet`: 输入符号无效
- `invalid transition`: 转移无效
- `DFA non-determinism`: DFA 不确定性
