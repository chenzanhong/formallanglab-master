# 自动机示例
## dfa ： 识别以ab结尾的字符串
```json
{
  "states": ["q0", "q1", "q2"],
  "alphabet": ["a", "b"],
  "transitions": [
    {
      "fromState": "q0",
      "input": "a",
      "toStates": ["q1"]
    },
    {
      "fromState": "q0",
      "input": "b",
      "toStates": ["q0"]
    },
    {
      "fromState": "q1",
      "input": "a",
      "toStates": ["q1"]
    },
    {
      "fromState": "q1",
      "input": "b",
      "toStates": ["q2"]
    },
    {
      "fromState": "q2",
      "input": "a",
      "toStates": ["q1"]
    },
    {
      "fromState": "q2",
      "input": "b",
      "toStates": ["q0"]
    }
  ],
  "initialState": "q0",
  "acceptingStates": ["q2"],
  "type": 0
}
```

## nfa，包含aba
```json

```

## nfa，接收空串和a
```json

```