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

### 3. 测试文件结构

```
backend/master/test/performance/
├── automaton_test.go     # 自动机性能测试
├── grammar_test.go       # 文法性能测试
├── regex_test.go         # 正则表达式性能测试
├── convert_test.go       # 转换接口性能测试
├── learn_test.go         # 学习模块性能测试
├── store_test.go         # 存储模块性能测试
└── run_all_tests.go      # 统一测试入口
```

### 4. 性能测试配置

```go
type PerformanceTestConfig struct {
    BaseURL         string // API 基础 URL
    Concurrency     int    // 并发用户数
    RequestsPerUser int    // 每用户请求数
    TotalUsers      int    // 总用户数
    Timeout         int    // 请求超时时间(秒)
}
```

### 5. 测试结果指标

- **TotalRequests** - 总请求数
- **Successful** - 成功请求数
- **Failed** - 失败请求数
- **TotalTime** - 总耗时
- **AvgResponseTime** - 平均响应时间
- **MinResponseTime** - 最小响应时间
- **MaxResponseTime** - 最大响应时间
- **P50/P95/P99ResponseTime** - 分位数响应时间
- **RequestsPerSec** - 吞吐量(每秒请求数)

### 6. 使用方法

```bash
# 运行所有测试
cd backend/master
go test -v ./test/performance/... -run TestAll

# 运行单个测试
go test -v ./test/performance -run TestAutomatonValidate

# 配置测试参数
export BASE_URL=http://localhost:8080
export CONCURRENCY=10
export REQUESTS_PER_USER=100
```

### 7. 下一步工作

1. ✅ 完成 automaton_test.go (已创建)
2. ⏳ 创建 grammar_test.go
3. ⏳ 创建 regex_test.go
4. ⏳ 创建 convert_test.go
5. ⏳ 创建 learn_test.go
6. ⏳ 创建 store_test.go
7. ⏳ 创建 run_all_tests.go
8. ⏳ 编写测试配置文件
9. ⏳ 编写测试报告模板

### 8. 注意事项

1. 所有接口都需要 JWT 认证
2. 需要启动后端服务后再运行测试
3. 测试数据需要覆盖各种边界情况
4. 建议先进行小规模测试，再逐步增加并发数
5. 监控服务器资源使用情况
