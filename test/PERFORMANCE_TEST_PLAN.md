# 后端 Master 性能测试方案

## 分析结果

### 1. API 路由结构

#### 文法相关接口 (10个)
- `/formallanglab/master/grammar/validate` - 文法验证
- `/formallanglab/master/grammar/type` - 文法类型判定
- `/formallanglab/master/grammar/simplify` - 文法化简
- `/formallanglab/master/grammar/ambiguity` - 文法二义性检查
- `/formallanglab/master/grammar/equivalence` - 文法等价性检查
- `/formallanglab/master/grammar/recognize` - 文法字符串识别
- `/formallanglab/master/grammar/generate` - 文法示例字符串生成
- `/formallanglab/master/grammar/first` - 文法 First 集计算
- `/formallanglab/master/grammar/follow` - 文法 Follow 集计算

#### 自动机相关接口 (7个)
- `/formallanglab/master/automaton/validate` - 自动机验证
- `/formallanglab/master/automaton/recognize` - 字符串识别
- `/formallanglab/master/automaton/cleanup` - 自动机清理
- `/formallanglab/master/automaton/minimize` - DFA 最小化
- `/formallanglab/master/automaton/nfatodfa` - NFA 确定化
- `/formallanglab/master/automaton/equivalence` - 自动机等价性检查
- `/formallanglab/master/automaton/generate` - 示例字符串生成

#### 正则表达式相关接口 (4个)
- `/formallanglab/master/regex/validate` - 正则表达式验证
- `/formallanglab/master/regex/recognize` - 字符串匹配
- `/formallanglab/master/regex/equivalence` - 等价性检查
- `/formallanglab/master/regex/generate` - 示例字符串生成

#### 转换接口 (4个)
- `/formallanglab/master/convert/grammar-to-nfa` - 文法转 NFA
- `/formallanglab/master/convert/fa-to-grammar` - 自动机转文法
- `/formallanglab/master/convert/regex-to-nfa` - 正则转 NFA
- `/formallanglab/master/convert/fa-to-regex` - 自动机转正则

#### 学习模块接口 (4个)
- `GET /formallanglab/master/learn/` - 获取学习资料列表
- `POST /formallanglab/master/learn/` - 添加学习资源
- `POST /formallanglab/master/learn/sync` - 同步 OSS 文件
- `GET /formallanglab/master/learn/:id` - 获取单个资源
- `DELETE /formallanglab/master/learn/:id` - 删除学习资源

#### 存储模块接口 (9个)
- `POST /formallanglab/master/store/automaton` - 创建自动机
- `GET /formallanglab/master/store/automatons` - 查询自动机列表
- `DELETE /formallanglab/master/store/automaton/:id` - 删除自动机
- `POST /formallanglab/master/store/grammar` - 创建文法
- `GET /formallanglab/master/store/grammars` - 查询文法列表
- `DELETE /formallanglab/master/store/grammar/:id` - 删除文法
- `POST /formallanglab/master/store/regex` - 创建正则
- `GET /formallanglab/master/store/regexes` - 查询正则列表
- `DELETE /formallanglab/master/store/regex/:id` - 删除正则

### 2. 指标收集系统

已集成 Prometheus 指标收集：
- `operation_result_total` - 业务操作计数
- `operation_duration_seconds` - 操作执行时间
- `http_requests_total` - HTTP 请求计数
- `http_request_duration_seconds` - HTTP 响应时间
