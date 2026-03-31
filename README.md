# FormalLangLab

FormalLangLab 是一个基于 Web 的交互式形式语言与自动机学习系统，旨在帮助计算机科学专业的学生更直观、高效地理解文法、正则表达式、有限自动机（DFA/NFA）等核心概念。

# Master 服务

## 技术栈

| 类别 | 技术 |
|------|------|
| 框架 | Golang (Gin) |
| 数据库 | PostgreSQL |
| 认证 | JWT |
| 消息队列 | Kafka（用于日志收集） |
| 对象存储 | 阿里云 OSS |
| 日志 | 结构化日志系统 (zlog) |
| 配置管理 | YAML 配置文件 |
| 监控 | Prometheus + Grafana |

## 项目结构

```
backend/master/
├── cmd/                        # 程序入口
│   └── main.go                 # 主服务入口
├── configs/                    # 配置文件
│   ├── config.go               # 配置结构定义与加载
│   └── config.yaml             # 配置文件
├── internal/                   # 核心业务代码
│   ├── api/                    # HTTP 接口层
│   │   ├── router.go           # 路由定义
│   │   ├── grammar.go          # 文法接口
│   │   ├── automaton.go        # 自动机接口
│   │   ├── regex.go            # 正则表达式接口
│   │   ├── convert.go          # 转换接口
│   │   ├── learn.go            # 学习资源接口
│   │   └── store.go            # 存储接口
│   ├── domain/                 # 领域模型
│   │   ├── dto/                # 数据传输对象
│   │   └── model/              # 领域模型定义
│   ├── service/                # 业务逻辑层
│   │   ├── grammar_s/          # 文法相关服务
│   │   ├── automaton_s/        # 自动机相关服务
│   │   ├── regex_s/            # 正则表达式相关服务
│   │   ├── learn_s/            # 学习资源服务
│   │   ├── store_s/            # 存储服务
│   │   └── kafka_s/            # Kafka 生产者
│   ├── repository/             # 数据访问层
│   ├── middleware/             # 中间件
│   ├── metrics/                # Prometheus 指标
│   └── errors/                 # 错误定义
├── pkg/                        # 公共包
│   ├── oss/                    # 阿里云 OSS 客户端
│   └── binding/                # Gin 参数绑定验证
├── logs/                       # 日志文件
├── migrations/                 # 数据库迁移脚本
└── test/                       # 测试文件
```

## 服务模块详解

### 1. 文法服务（Grammar Service）

#### 1.1 文法基础操作

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| POST | `/gdesign/master/grammar/validate` | 文法校验：验证文法是否有效 |
| POST | `/gdesign/master/grammar/type` | 文法类型判定：判断文法的 Chomsky 类型（0、1、2、3型） |
| POST | `/gdesign/master/grammar/simplify` | 文法化简：去除无用符号（不可派生、不可达）、单一产生式、空产生式 |
| POST | `/gdesign/master/grammar/first` | First 集计算：计算文法的 First 集 |
| POST | `/gdesign/master/grammar/follow` | Follow 集计算：计算文法的 Follow 集 |

#### 1.2 文法分析功能

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| POST | `/gdesign/master/grammar/ambiguity` | 二义性检测：判断正则文法是否存在二义性 |
| POST | `/gdesign/master/grammar/equivalence` | 文法等价性检查：判断两个正则文法是否等价 |
| POST | `/gdesign/master/grammar/recognize` | 字符串识别：判断字符串是否能被指定文法接受 |
| POST | `/gdesign/master/grammar/generate` | 示例字符串生成：生成可推导和不可推导的字符串 |

### 2. 自动机服务（Automaton Service）

#### 2.1 自动机基础操作

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| POST | `/gdesign/master/automaton/validate` | 自动机有效性验证：检查自动机定义是否正确 |
| POST | `/gdesign/master/automaton/cleanup` | 自动机清理：去除无效状态和不可达状态 |
| POST | `/gdesign/master/automaton/minimize` | DFA 最小化：使用 Hopcroft 算法将 DFA 转换为最小状态数的等价 DFA |
| POST | `/gdesign/master/automaton/nfatodfa` | NFA 确定化：使用子集构造法将 NFA 转换为等价的 DFA |

#### 2.2 自动机分析功能

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| POST | `/gdesign/master/automaton/equivalence` | 自动机等价性检查：判断两个自动机是否等价 |
| POST | `/gdesign/master/automaton/recognize` | 字符串识别：判断字符串是否能被自动机接受 |
| POST | `/gdesign/master/automaton/generate` | 示例字符串生成：生成可接受和不可接受的字符串 |

### 3. 正则表达式服务（Regex Service）

#### 3.1 正则表达式基础操作

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| POST | `/gdesign/master/regex/validate` | 正则表达式有效性验证：检查正则表达式语法是否正确 |

#### 3.2 正则表达式分析功能

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| POST | `/gdesign/master/regex/equivalence` | 正则表达式等价性检查：判断两个正则表达式是否等价 |
| POST | `/gdesign/master/regex/recognize` | 字符串识别：判断字符串是否匹配正则表达式 |
| POST | `/gdesign/master/regex/generate` | 示例字符串生成：生成可匹配和不可匹配的字符串 |

### 4. 转换服务（Convert Service）

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| POST | `/gdesign/master/convert/grammar-to-nfa` | 文法转自动机：将右线性文法转换为等价的 NFA |
| POST | `/gdesign/master/convert/fa-to-grammar` | 自动机转文法：将 FA 转换为等价的文法 |
| POST | `/gdesign/master/convert/regex-to-nfa` | 正则表达式转自动机：使用 Thompson 构造法将正则表达式转换为 NFA |
| POST | `/gdesign/master/convert/fa-to-regex` | 自动机转正则表达式：将 FA 转换为等价的正则表达式 |

### 5. 存储服务（Store Service）

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| POST | `/gdesign/master/store/automaton` | 自动机创建与保存 |
| GET | `/gdesign/master/store/automatons` | 自动机列表查询与分页 |
| DELETE | `/gdesign/master/store/automaton/:id` | 自动机删除 |
| POST | `/gdesign/master/store/grammar` | 文法创建与保存 |
| GET | `/gdesign/master/store/grammars` | 文法列表查询与分页 |
| DELETE | `/gdesign/master/store/grammar/:id` | 文法删除 |
| POST | `/gdesign/master/store/regex` | 正则表达式创建与保存 |
| GET | `/gdesign/master/store/regexes` | 正则表达式列表查询与分页 |
| DELETE | `/gdesign/master/store/regex/:id` | 正则表达式删除 |

### 6. 学习资源服务（Learn Service）

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| GET | `/gdesign/master/learn/` | 学习资源列表获取 |
| POST | `/gdesign/master/learn/` | 添加学习资源（管理员在 OSS 上传后调用） |
| POST | `/gdesign/master/learn/sync` | 同步 OSS 文件到数据库（自动识别新增文件） |
| GET | `/gdesign/master/learn/:id` | 学习资源详情查询 |
| DELETE | `/gdesign/master/learn/:id` | 学习资源删除 |

### 7. 系统监控接口

| 接口 | 路由 | 功能描述 |
|------|------|----------|
| GET | `/gdesign/master/metrics` | Prometheus 指标采集端点 |
| GET | `/gdesign/master/health` | 健康检查接口 |

## 核心算法实现

### NFA 确定化（子集构造法）

核心文件：[nfatodfa.go](internal/service/automaton_s/nfatodfa.go)

算法步骤：
1. 计算初始状态的 ε-闭包作为 DFA 初始状态
2. BFS 遍历所有状态集合，对每个输入符号计算转移
3. 对转移结果取 ε-闭包
4. 生成新的 DFA 状态和转移

### DFA 最小化（Hopcroft 算法）

核心文件：[minimize.go](internal/service/automaton_s/minimize.go)

算法步骤：
1. 移除不可达状态
2. 初始划分：接受状态 vs 非接受状态
3. 构建反向转移图
4. 迭代划分直到收敛
5. 构建最小化 DFA

### 正则表达式转 NFA（Thompson 构造法）

核心文件：[regex_fa.go](internal/service/regex_s/regex_fa.go)

支持的运算：
- 基本符号
- 连接运算
- 选择运算 (|)
- 闭包运算 (*)
- 正闭包 (+)
- 可选运算 (?)

## 中间件

| 中间件 | 功能 |
|--------|------|
| RequestID | 为每个请求生成唯一 ID |
| CORS | 跨域请求处理 |
| JWT | JWT Token 认证 |
| UserRateLimit | 用户级速率限制 |
| HTTPMiddleware | Prometheus 指标收集 |

## 配置说明

### 环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `SERVER_PORT` | 服务端口 | `8081` |
| `METRICS_PORT` | 指标端口 | `9091` |
| `PPROF_PORT` | pprof 端口 | `6060` |
| `JWT_KEY` | JWT 密钥 | `your-secret-key` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PORT` | 数据库端口 | `5432` |
| `DB_NAME` | 数据库名称 | `formallanglab` |
| `DB_USER` | 数据库用户 | `postgres` |
| `DB_PASSWORD` | 数据库密码 | `password` |
| `KAFKA_BROKERS` | Kafka 地址 | `localhost:9092` |
| `KAFKA_TOPIC` | Kafka 主题 | `logs` |
| `OSS_ENDPOINT` | OSS 端点 | `oss-cn-hangzhou.aliyuncs.com` |
| `OSS_ACCESS_KEY_ID` | OSS AccessKey ID | - |
| `OSS_ACCESS_KEY_SECRET` | OSS AccessKey Secret | - |
| `OSS_BUCKET_NAME` | OSS Bucket 名称 | - |
| `RATE_USER_RATE` | 用户速率限制 | `100` |
| `RATE_USER_BURST` | 用户突发限制 | `200` |

### 配置文件示例 (config.yaml)

```yaml
server:
  port: 8081
  metrics_port: 9091
  pprof_port: 6060
  enable_trace: false

jwt:
  key: "your-jwt-secret-key"

pg:
  host: localhost
  port: 5432
  name: formallanglab
  user: postgres
  password: password

rate:
  user_rate: 100
  user_burst: 200

kafka:
  brokers: localhost:9092
  topic: logs

oss:
  endpoint: oss-cn-hangzhou.aliyuncs.com
  access_key_id: your-access-key-id
  access_key_secret: your-access-key-secret
  bucket_name: your-bucket-name
  region: cn-hangzhou
  base_url: https://your-bucket-name.oss-cn-hangzhou.aliyuncs.com

log:
  level: info
  output: console
  format: json
  file_path: logs/app.log
  max_size: 100
  max_backups: 5
  max_age: 30
  compress: true
  sampling: false
```

## 部署指南

### 环境要求

- Go 1.21+
- PostgreSQL 13+
- Kafka（可选，用于日志收集）

### 本地启动

```bash
# 进入项目目录
cd backend/master

# 设置环境变量
export SERVER_PORT=8081
export DB_HOST=localhost
export DB_PORT=5432
# ... 其他环境变量

# 启动服务
go run cmd/main.go
```

### Docker 部署

```bash
# 构建镜像
docker build -t formallanglab-master .

# 运行容器
docker run -d \
  -p 8081:8081 \
  -e SERVER_PORT=8081 \
  -e DB_HOST=postgres \
  formallanglab-master
```

## 关键技术点

| 模块 | 技术点 | 实现文件 |
|------|--------|----------|
| 文法合法性判断 | 正则文法、上下文无关文法识别 | `grammar_s/validate.go` |
| NFA → DFA 转换 | 子集构造法、ε-闭包 | `automaton_s/nfatodfa.go` |
| DFA 最小化 | Hopcroft 算法 | `automaton_s/minimize.go` |
| 自动机等价性检查 | 状态等价判断 | `automaton_s/equivalence_check.go` |
| 正则表达式转 NFA | Thompson 构造法 | `regex_s/regex_fa.go` |
| 自动机转正则表达式 | 状态消去法 | `automaton_s/fa_regex.go` |
| 文法与自动机转换 | 右线性文法转换 | `grammar_s/grammar_fa.go` |

## 服务启动流程

```
main.go
    │
    ├── 1. 加载配置 (configs.LoadConfig)
    │
    ├── 2. 同步配置到环境变量 (configs.SyncConfigToEnv)
    │
    ├── 3. 初始化 JWT (jwtx.InitWithHS256)
    │
    ├── 4. 初始化日志 (zlog.InitLogger)
    │
    ├── 5. 初始化数据库 (rep.Init)
    │
    ├── 6. 初始化 OSS 客户端 (oss.NewAliyunOSSClient)
    │
    ├── 7. 组装服务层
    │       ├── StoreRepository
    │       ├── LearnRepository
    │       ├── KafkaProducer
    │       ├── StoreService
    │       └── LearnService
    │
    ├── 8. 初始化 Handler
    │       ├── StoreHandler
    │       └── LearnHandler
    │
    ├── 9. 注册路由 (api.SetupRouter)
    │
    ├── 10. 启动 HTTP 服务
    │
    └── 11. 优雅关闭处理
            ├── 关闭 HTTP 服务
            ├── 关闭数据库连接
            └── 关闭 Kafka 生产者
```

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (React)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway / Nginx                       │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌─────────┐     ┌─────────┐     ┌─────────┐
        │  Auth   │     │ Master  │     │   AI    │
        │ Service │     │ Service │     │ Service │
        └─────────┘     └─────────┘     └─────────┘
              │               │               │
              └───────────────┼───────────────┘
                              ▼
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    └─────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌─────────┐     ┌─────────┐     ┌─────────┐
        │  Kafka  │     │   OSS   │     │  Redis  │
        │ (Logs)  │     │ (Files) │     │ (Cache) │
        └─────────┘     └─────────┘     └─────────┘
```

## 总结

Master 服务是形式语言学习系统的核心服务，提供了：

- **完整的文法处理**：验证、类型判断、化简、First/Follow 集计算
- **强大的自动机操作**：验证、NFA 确定化、DFA 最小化、等价性检查
- **正则表达式支持**：验证、匹配、转换
- **形式语言互转**：文法 ↔ 自动机 ↔ 正则表达式
- **数据持久化**：自动机、文法、正则表达式的存储管理
- **学习资源管理**：基于 OSS 的学习资料管理
