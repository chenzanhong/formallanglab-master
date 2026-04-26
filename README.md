# Master 服务

## 模块定位

Master 服务是 FormalLangLab 项目的核心业务模块，提供形式语言与自动机理论的所有核心算法功能。作为系统的计算引擎，Master 服务负责文法、自动机、正则表达式的解析、验证、转换和分析，为前端提供强大的理论计算支持。

## 功能特性

- **文法处理**：文法验证、类型判定、化简、First/Follow 集计算、二义性检测
- **自动机处理**：DFA/NFA 验证、最小化、确定化、等价性检查
- **正则表达式处理**：正则验证、等价性检查
- **形式转换**：文法↔自动机、正则表达式↔自动机相互转换
- **字符串识别**：判断字符串是否能被文法/自动机/正则表达式接受
- **示例生成**：生成可接受和不可接受的字符串示例
- **用户存储**：支持用户保存和管理自己的文法、自动机、正则表达式
- **学习资源**：提供形式语言相关的学习资源和教程

## 技术栈

| 类别 | 技术 |
|------|------|
| 开发语言 | Go 1.23+ |
| Web 框架 | Gin |
| 数据库 | PostgreSQL |
| 对象存储 | 阿里云 OSS |
| 消息队列 | Kafka（日志收集） |
| 认证机制 | JWT |
| 日志系统 | 结构化日志（zlog） |
| 配置管理 | YAML 配置文件 |
| 监控指标 | Prometheus + Grafana |

## 项目结构

```
backend/master/
├── cmd/                        # 程序入口
│   └── main.go                 # 主服务入口
├── configs/                    # 配置管理
│   ├── config.go               # 配置结构定义与加载
│   └── config.yaml.example     # 配置文件示例
├── internal/                   # 核心业务代码
│   ├── domain/                 # 领域模型
│   │   ├── dto/                # 数据传输对象
│   │   │   ├── grammar.go      # 文法相关 DTO
│   │   │   ├── automaton.go    # 自动机相关 DTO
│   │   │   ├── regex.go        # 正则表达式相关 DTO
│   │   │   ├── convert.go      # 转换相关 DTO
│   │   │   ├── learn.go        # 学习资源 DTO
│   │   │   ├── store.go        # 存储相关 DTO
│   │   │   └── base.go         # 基础 DTO
│   │   ├── model/              # 领域实体模型
│   │   │   ├── grammar.go      # 文法模型
│   │   │   ├── automaton.go    # 自动机模型
│   │   │   ├── regex.go        # 正则表达式模型
│   │   │   ├── convert.go      # 转换模型
│   │   │   ├── learn.go        # 学习资源模型
│   │   │   └── fa_process.go   # 自动机处理流程模型
│   │   └── storage/            # 存储实体
│   │       ├── grammar.go      # 文法存储
│   │       ├── automaton.go    # 自动机存储
│   │       └── regex.go        # 正则表达式存储
│   ├── handler/                # HTTP 接口层
│   │   ├── router.go           # 路由定义
│   │   ├── grammar.go          # 文法处理器
│   │   ├── automaton.go        # 自动机处理器
│   │   ├── regex.go            # 正则表达式处理器
│   │   ├── convert.go          # 转换处理器
│   │   ├── store.go            # 存储处理器
│   │   └── learn.go            # 学习资源处理器
│   ├── service/                # 业务逻辑层
│   │   ├── grammar_s/          # 文法服务
│   │   │   ├── validate.go     # 文法验证
│   │   │   ├── type_determine.go # 文法类型判定
│   │   │   ├── simplify.go     # 文法化简
│   │   │   ├── ambiguity_check.go # 二义性检测
│   │   │   ├── equivalence_check.go # 等价性检查
│   │   │   ├── recognize.go    # 字符串识别
│   │   │   ├── generate.go     # 示例生成
│   │   │   ├── first_follow.go # First/Follow 集
│   │   │   └── grammar_fa.go   # 文法转自动机
│   │   ├── automaton_s/        # 自动机服务
│   │   │   ├── validate.go     # 自动机验证
│   │   │   ├── clean.go        # 自动机清理
│   │   │   ├── minimize.go     # DFA 最小化
│   │   │   ├── nfatodfa.go     # NFA 确定化
│   │   │   ├── equivalence_check.go # 等价性检查
│   │   │   ├── recognize.go    # 字符串识别
│   │   │   ├── generate.go     # 示例生成
│   │   │   ├── complete.go     # DFA 补全
│   │   │   ├── complement_dfa.go # DFA 补集
│   │   │   └── fa_grammar.go   # 自动机转文法
│   │   ├── regex_s/            # 正则表达式服务
│   │   │   ├── validate.go     # 正则验证
│   │   │   ├── equivalence_check.go # 等价性检查
│   │   │   ├── recognize.go    # 字符串匹配
│   │   │   ├── generate.go     # 示例生成
│   │   │   └── regex_fa.go     # 正则转自动机
│   │   ├── convert_s/          # 转换服务
│   │   │   ├── grammar_to_nfa.go # 文法转 NFA
│   │   │   ├── fa_to_grammar.go  # 自动机转文法
│   │   │   ├── regex_to_nfa.go   # 正则转 NFA
│   │   │   └── fa_to_regex.go    # 自动机转正则
│   │   ├── store_s/            # 存储服务
│   │   │   └── store_service.go # 存储服务实现
│   │   ├── learn_s/            # 学习资源服务
│   │   │   └── learn.go        # 学习资源管理
│   │   └── kafka_s/            # Kafka 服务
│   │       └── producer.go     # Kafka 生产者（日志收集）
│   ├── repository/             # 数据访问层
│   │   ├── init.go             # 存储库初始化
│   │   ├── store_repo.go       # 存储数据访问
│   │   └── learn_repo.go       # 学习资源数据访问
│   ├── middleware/             # 中间件
│   │   ├── cors/               # 跨域中间件
│   │   ├── jwt/                # JWT 认证中间件
│   │   ├── metrics/            # 指标采集中间件
│   │   ├── rate/               # 速率限制中间件
│   │   └── requestid/          # 请求 ID 中间件
│   └── errors/                 # 错误定义
├── pkg/                        # 公共包
│   ├── oss/                    # 阿里云 OSS 客户端
│   │   └── client.go           # OSS 客户端封装
│   ├── binding/                # 参数绑定与验证
│   │   └── bind.go             # 自定义验证器
│   ├── contextx/               # 上下文工具
│   │   └── contextx.go         # 上下文扩展
│   ├── cryptoutil/             # 加密工具
│   │   └── crypto.go           # 加密函数
│   ├── token/                  # 令牌工具
│   │   └── token.go           # JWT 令牌处理
│   ├── http_client/            # HTTP 客户端
│   │   └── client.go           # HTTP 客户端封装
│   └── util/                   # 通用工具
│       └── util.go             # 工具函数
├── migrations/                 # 数据库迁移脚本
├── scripts/                    # 脚本工具
│   └── init_chroma.go          # 初始化 Chroma 向量库
├── test/                       # 测试文件
│   └── performance/            # 性能测试
│       ├── algorithm_perf_test.go # 算法性能测试
│       ├── automaton_test.go   # 自动机测试
│       └── concurrency_test.go # 并发测试
├── logs/                       # 日志文件目录
├── .gitignore                  # Git 忽略配置
├── .golangci.yml               # Go 代码检查配置
├── Dockerfile                  # Docker 镜像构建文件
├── Makefile                    # Make 命令配置
├── go.mod                      # Go 模块依赖
└── go.sum                      # Go 依赖校验文件
```

## 中间件

| 中间件 | 文件路径 | 功能描述 |
|--------|---------|---------|
| RequestID | `middleware/requestid/requestid.go` | 为每个请求生成唯一 ID，便于日志追踪和调试 |
| CORS | `middleware/cors/cors.go` | 处理跨域请求，支持预检请求和凭证传递 |
| JWT | `middleware/jwt/jwt.go` | JWT Token 认证，验证请求头中的 Token 有效性 |
| UserRateLimit | `middleware/rate/rate.go` | 用户级速率限制，防止恶意请求和资源滥用 |
| Metrics | `middleware/metrics/metrics.go` | Prometheus 指标采集，监控请求延迟、QPS 和错误率 |

## 核心服务

### 文法服务（Grammar Service）

**文件**：`internal/service/grammar_s/`

**功能模块**：

#### 文法基础操作
- **验证（Validate）**：验证文法是否符合形式语言定义
- **类型判定（Type Determine）**：判断文法的 Chomsky 类型（0、1、2、3 型）
- **化简（Simplify）**：去除无用符号、单一产生式、空产生式
- **First 集计算**：计算文法的 First 集
- **Follow 集计算**：计算文法的 Follow 集

#### 文法分析功能
- **二义性检测（Ambiguity Check）**：判断正则文法是否存在二义性
- **等价性检查（Equivalence Check）**：判断两个正则文法是否等价
- **字符串识别（Recognize）**：判断字符串是否能被指定文法接受
- **示例生成（Generate）**：生成可推导和不可推导的字符串示例

#### 文法与自动机转换
- **文法转自动机（Grammar to FA）**：将右线性文法转换为等价的 NFA

### 自动机服务（Automaton Service）

**文件**：`internal/service/automaton_s/`

**功能模块**：

#### 自动机基础操作
- **验证（Validate）**：检查自动机定义是否正确
- **清理（Clean）**：去除无效状态和不可达状态
- **最小化（Minimize）**：使用 Hopcroft 算法将 DFA 转换为最小状态数的等价 DFA
- **确定化（NFA to DFA）**：使用子集构造法将 NFA 转换为等价的 DFA
- **补全（Complete）**：补全 DFA 的转移函数
- **补集（Complement）**：构造 DFA 的补集

#### 自动机分析功能
- **等价性检查（Equivalence Check）**：判断两个自动机是否等价
- **字符串识别（Recognize）**：判断字符串是否能被自动机接受
- **示例生成（Generate）**：生成可接受和不可接受的字符串示例

#### 自动机与其他形式转换
- **自动机转文法（FA to Grammar）**：将 FA 转换为等价的文法
- **自动机转正则（FA to Regex）**：将 FA 转换为等价的正则表达式

### 正则表达式服务（Regex Service）

**文件**：`internal/service/regex_s/`

**功能模块**：

#### 正则表达式基础操作
- **验证（Validate）**：检查正则表达式语法是否正确

#### 正则表达式分析功能
- **等价性检查（Equivalence Check）**：判断两个正则表达式是否等价
- **字符串匹配（Recognize）**：判断字符串是否匹配正则表达式
- **示例生成（Generate）**：生成可匹配和不可匹配的字符串示例

#### 正则与自动机转换
- **正则转自动机（Regex to FA）**：使用 Thompson 构造法将正则表达式转换为 NFA

### 转换服务（Convert Service）

**文件**：`internal/service/convert_s/`

**功能**：
- 文法转 NFA
- NFA 转文法
- 正则转 NFA
- NFA 转正则
- 提供统一的形式转换接口

### 存储服务（Store Service）

**文件**：`internal/service/store_s/`

**功能**：
- 用户文法的创建、查询、删除
- 用户自动机的创建、查询、删除
- 用户正则表达式的创建、查询、删除
- 支持分页查询和筛选

### 学习资源服务（Learn Service）

**文件**：`internal/service/learn_s/`

**功能**：
- 学习资源列表获取
- 学习资源详情查询
- 学习资源上传（管理员）
- 学习资源同步（从 OSS 同步到数据库）
- 学习资源删除（管理员）

## 数据访问层

### 存储存储库（Store Repository）

**文件**：`internal/repository/store_repo.go`

**功能**：
- 用户文法数据的 CRUD 操作
- 用户自动机数据的 CRUD 操作
- 用户正则表达式数据的 CRUD 操作
- 支持分页查询和条件筛选

### 学习资源存储库（Learn Repository）

**文件**：`internal/repository/learn_repo.go`

**功能**：
- 学习资源数据的 CRUD 操作
- 资源分类和标签管理
- 资源热度统计

## 领域模型

### 数据传输对象（DTO）

**文件**：`internal/domain/dto/`

**主要 DTO**：
- `GrammarRequest`：文法请求体
- `AutomatonRequest`：自动机请求体
- `RegexRequest`：正则表达式请求体
- `ConvertRequest`：转换请求体
- `StoreRequest`：存储请求体
- `LearnResourceRequest`：学习资源请求体

### 领域实体（Model）

**文件**：`internal/domain/model/`

**主要模型**：
- `Grammar`：文法实体（产生式、终结符、非终结符、开始符号）
- `Automaton`：自动机实体（状态集、字母表、转移函数、开始状态、接受状态集）
- `Regex`：正则表达式实体（表达式字符串、解析树）
- `ConversionResult`：转换结果实体
- `LearnResource`：学习资源实体

## API 接口

### 文法接口

| 方法 | 路由 | 功能 |
|------|------|------|
| POST | `/formallanglab/master/grammar/validate` | 文法验证 |
| POST | `/formallanglab/master/grammar/type` | 文法类型判定 |
| POST | `/formallanglab/master/grammar/simplify` | 文法化简 |
| POST | `/formallanglab/master/grammar/first` | First 集计算 |
| POST | `/formallanglab/master/grammar/follow` | Follow 集计算 |
| POST | `/formallanglab/master/grammar/ambiguity` | 二义性检测 |
| POST | `/formallanglab/master/grammar/equivalence` | 文法等价性检查 |
| POST | `/formallanglab/master/grammar/recognize` | 字符串识别 |
| POST | `/formallanglab/master/grammar/generate` | 示例字符串生成 |

### 自动机接口

| 方法 | 路由 | 功能 |
|------|------|------|
| POST | `/formallanglab/master/automaton/validate` | 自动机验证 |
| POST | `/formallanglab/master/automaton/cleanup` | 自动机清理 |
| POST | `/formallanglab/master/automaton/minimize` | DFA 最小化 |
| POST | `/formallanglab/master/automaton/nfatodfa` | NFA 确定化 |
| POST | `/formallanglab/master/automaton/equivalence` | 自动机等价性检查 |
| POST | `/formallanglab/master/automaton/recognize` | 字符串识别 |
| POST | `/formallanglab/master/automaton/generate` | 示例字符串生成 |
| POST | `/formallanglab/master/automaton/complete` | DFA 补全 |
| POST | `/formallanglab/master/automaton/complement` | DFA 补集 |

### 正则表达式接口

| 方法 | 路由 | 功能 |
|------|------|------|
| POST | `/formallanglab/master/regex/validate` | 正则表达式验证 |
| POST | `/formallanglab/master/regex/equivalence` | 正则等价性检查 |
| POST | `/formallanglab/master/regex/recognize` | 字符串匹配 |
| POST | `/formallanglab/master/regex/generate` | 示例字符串生成 |

### 转换接口

| 方法 | 路由 | 功能 |
|------|------|------|
| POST | `/formallanglab/master/convert/grammar-to-nfa` | 文法转 NFA |
| POST | `/formallanglab/master/convert/fa-to-grammar` | 自动机转文法 |
| POST | `/formallanglab/master/convert/regex-to-nfa` | 正则转 NFA |
| POST | `/formallanglab/master/convert/fa-to-regex` | 自动机转正则 |

### 存储接口

| 方法 | 路由 | 功能 |
|------|------|------|
| POST | `/formallanglab/master/store/automaton` | 保存自动机 |
| GET | `/formallanglab/master/store/automatons` | 获取自动机列表 |
| DELETE | `/formallanglab/master/store/automaton/:id` | 删除自动机 |
| POST | `/formallanglab/master/store/grammar` | 保存文法 |
| GET | `/formallanglab/master/store/grammars` | 获取文法列表 |
| DELETE | `/formallanglab/master/store/grammar/:id` | 删除文法 |
| POST | `/formallanglab/master/store/regex` | 保存正则 |
| GET | `/formallanglab/master/store/regexes` | 获取正则列表 |
| DELETE | `/formallanglab/master/store/regex/:id` | 删除正则 |

### 学习资源接口

| 方法 | 路由 | 功能 |
|------|------|------|
| GET | `/formallanglab/master/learn/` | 获取学习资源列表 |
| GET | `/formallanglab/master/learn/:id` | 获取学习资源详情 |
| POST | `/formallanglab/master/learn/` | 添加学习资源（管理员） |
| POST | `/formallanglab/master/learn/sync` | 同步 OSS 文件到数据库 |
| DELETE | `/formallanglab/master/learn/:id` | 删除学习资源（管理员） |

### 系统接口

| 方法 | 路由 | 功能 |
|------|------|------|
| GET | `/formallanglab/master/metrics` | Prometheus 指标采集 |
| GET | `/formallanglab/master/health` | 健康检查 |

## 核心算法

### NFA 确定化（子集构造法）

**文件**：`internal/service/automaton_s/nfatodfa.go`

**算法步骤**：
1. 计算初始状态的 ε-闭包作为 DFA 初始状态
2. BFS 遍历所有状态集合，对每个输入符号计算转移
3. 对转移结果取 ε-闭包
4. 生成新的 DFA 状态和转移
5. 标记包含原 NFA 接受状态的状态集为 DFA 接受状态

### DFA 最小化（Hopcroft 算法）

**文件**：`internal/service/automaton_s/minimize.go`

**算法步骤**：
1. 移除不可达状态
2. 初始划分：接受状态 vs 非接受状态
3. 构建反向转移图
4. 迭代划分直到收敛（无法再细分）
5. 从每个划分中选取代表状态
6. 构建最小化 DFA 的转移函数

### 正则表达式转 NFA（Thompson 构造法）

**文件**：`internal/service/regex_s/regex_fa.go`

**支持的运算**：
- 基本符号：单个字符
- 连接运算：ab
- 选择运算：a|b
- 闭包运算：a*
- 正闭包运算：a+
- 可选运算：a?

**构造规则**：
- 基本符号：创建两个状态，添加带标签的转移
- 连接：串联两个 NFA
- 选择：添加新的开始和接受状态，并行连接两个 NFA
- 闭包：添加ε转移形成循环

### 文法类型判定

**文件**：`internal/service/grammar_s/type_determine.go`

**判定规则**：
- **3 型（正则文法）**：右线性或左线性文法
- **2 型（上下文无关文法）**：产生式左部为单个非终结符
- **1 型（上下文有关文法）**：|α| ≤ |β|（α→β）
- **0 型（短语结构文法）**：无限制

### 二义性检测

**文件**：`internal/service/grammar_s/ambiguity_check.go`

**检测方法**：
- 构造文法的分析树
- 检查是否存在句子对应多棵不同的分析树
- 使用 LR(0) 项目集规范族检测冲突

## 配置说明

**文件**：`configs/config.yaml.example`

**主要配置项**：

### 服务配置
- 服务端口配置
- 服务名称和版本

### 数据库配置
- PostgreSQL 连接信息（主机、端口、用户名、密码、数据库名）
- 连接池配置（最大连接数、空闲连接数）

### JWT 配置
- JWT 密钥
- Access Token 过期时间
- Refresh Token 过期时间

### OSS 配置
- 阿里云 OSS 访问密钥（AccessKey ID/Secret）
- OSS Bucket 名称
- OSS Endpoint

### Kafka 配置
- Kafka 集群地址
- Kafka 主题配置
- 生产者配置

### 日志配置
- 日志级别
- 日志格式
- 日志输出路径

### 速率限制配置
- 用户级速率限制
- 接口级速率限制

## 监控指标

### Prometheus 指标

**主要指标**：
- `api_request_total`：API 请求总数（Counter）
- `api_request_duration_seconds`：API 请求延迟直方图（Histogram）
- `algorithm_execution_total`：算法执行次数（Counter）
- `algorithm_execution_duration_seconds`：算法执行延迟（Histogram）
- `database_connection_pool_size`：数据库连接池大小（Gauge）
- `oss_operation_total`：OSS 操作次数（Counter）

## 测试

### 性能测试

**文件**：`test/performance/`

**测试内容**：
- 算法性能测试：测试各核心算法在不同输入规模下的性能
- 自动机测试：测试自动机相关功能的正确性
- 并发测试：测试高并发场景下的系统表现

### 单元测试

每个核心算法模块都配有详细的单元测试，确保算法正确性。

## Docker 与 Makefile 使用说明

### Dockerfile 说明

Dockerfile 采用多阶段构建，分为构建阶段和运行阶段：

**构建阶段**：
- 使用 Go 1.24 Alpine 镜像作为构建环境
- 配置国内 Go 代理镜像源加速依赖下载
- 先下载依赖再复制源代码，优化 Docker 缓存层
- 构建静态二进制文件（CGO_ENABLED=0）

**运行阶段**：
- 使用 Alpine 3.20 精简镜像
- 创建非 root 用户（masteruser）运行应用，提升安全性
- 配置 Asia/Shanghai 时区
- 复制二进制文件、配置文件和迁移脚本
- 暴露服务端口（默认 8080）

### Makefile 命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `make lint` | 运行代码质量检查 | `make lint` |
| `make lint-fix` | 运行代码质量检查并自动修复 | `make lint-fix` |
| `make build` | 构建并推送 Docker 镜像 | `make build` |
| `make build T=false` | 仅构建镜像，不推送 | `make build T=false` |
| `make deploy` | 构建镜像并部署服务 | `make deploy` |

**快速开始**：
```bash
# 1. 代码质量检查
make lint

# 2. 构建镜像（本地测试，不推送）
make build T=false

# 3. 构建并推送镜像
make build

# 4. 部署服务
make deploy
```
