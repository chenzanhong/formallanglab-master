# 后端 Master 性能测试方案

## 分析结果

### 1. API 路由结构

#### 文法相关接口 (10个)
- `/gdesign/master/grammar/validate` - 文法验证
- `/gdesign/master/grammar/type` - 文法类型判定
- `/gdesign/master/grammar/simplify` - 文法化简
- `/gdesign/master/grammar/ambiguity` - 文法二义性检查
- `/gdesign/master/grammar/equivalence` - 文法等价性检查
- `/gdesign/master/grammar/recognize` - 文法字符串识别
- `/gdesign/master/grammar/generate` - 文法示例字符串生成
- `/gdesign/master/grammar/first` - 文法 First 集计算
- `/gdesign/master/grammar/follow` - 文法 Follow 集计算

#### 自动机相关接口 (7个)
- `/gdesign/master/automaton/validate` - 自动机验证
- `/gdesign/master/automaton/recognize` - 字符串识别
- `/gdesign/master/automaton/cleanup` - 自动机清理
- `/gdesign/master/automaton/minimize` - DFA 最小化
- `/gdesign/master/automaton/nfatodfa` - NFA 确定化
- `/gdesign/master/automaton/equivalence` - 自动机等价性检查
- `/gdesign/master/automaton/generate` - 示例字符串生成

#### 正则表达式相关接口 (4个)
- `/gdesign/master/regex/validate` - 正则表达式验证
- `/gdesign/master/regex/recognize` - 字符串匹配
- `/gdesign/master/regex/equivalence` - 等价性检查
- `/gdesign/master/regex/generate` - 示例字符串生成

#### 转换接口 (4个)
- `/gdesign/master/convert/grammar-to-nfa` - 文法转 NFA
- `/gdesign/master/convert/fa-to-grammar` - 自动机转文法
- `/gdesign/master/convert/regex-to-nfa` - 正则转 NFA
- `/gdesign/master/convert/fa-to-regex` - 自动机转正则

#### 学习模块接口 (4个)
- `GET /gdesign/master/learn/` - 获取学习资料列表
- `POST /gdesign/master/learn/` - 添加学习资源
- `POST /gdesign/master/learn/sync` - 同步 OSS 文件
- `GET /gdesign/master/learn/:id` - 获取单个资源
- `DELETE /gdesign/master/learn/:id` - 删除学习资源

#### 存储模块接口 (9个)
- `POST /gdesign/master/store/automaton` - 创建自动机
- `GET /gdesign/master/store/automatons` - 查询自动机列表
- `DELETE /gdesign/master/store/automaton/:id` - 删除自动机
- `POST /gdesign/master/store/grammar` - 创建文法
- `GET /gdesign/master/store/grammars` - 查询文法列表
- `DELETE /gdesign/master/store/grammar/:id` - 删除文法
- `POST /gdesign/master/store/regex` - 创建正则
- `GET /gdesign/master/store/regexes` - 查询正则列表
- `DELETE /gdesign/master/store/regex/:id` - 删除正则

### 2. 指标收集系统

已集成 Prometheus 指标收集：
- `operation_result_total` - 业务操作计数
- `operation_duration_seconds` - 操作执行时间
- `http_requests_total` - HTTP 请求计数
- `http_request_duration_seconds` - HTTP 响应时间
