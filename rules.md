# Go 微服务项目开发规范

你是一个 Go 语言专家,帮助开发者编写符合团队规范的代码。

## 核心原则

- 使用 Go 1.25+ 特性
- 统一使用 slog 日志接口,禁止直接使用 zap
- JSON 字段名统一使用 snake_case
- Protobuf 字段名使用 snake_case
- 所有外部调用必须设置超时

## 代码规范

### 强制规范 (MUST)

#### 包命名
- 包名使用小写字母,不使用下划线或驼峰
- 包名简短且有意义,避免泛化词汇

#### 变量命名
- 局部变量使用短名称: ctx, err, cfg, buf, req, resp
- 包级变量使用完整名称: DefaultTimeout, MaxRetryCount
- 常见缩写保持大写: ID, URL, HTTP, JSON, XML
- 避免类型前缀: userID 而非 intUserID

#### 错误处理
- 必须检查并处理所有 error 返回值
- 包级错误使用 Err 前缀: ErrNotFound, ErrInvalidInput
- 使用 errors.Is 和 errors.As 进行错误判断

#### JSON 标签
- struct tag 必须使用 snake_case
- 示例: `json:"user_id"` 而非 `json:"userId"`

#### 日志接口
- 统一使用 log/slog 标准库接口
- 禁止直接依赖 zap 或其他日志库
- 使用 github.com/chinayin/gox/log 封装

#### 代码质量
- 圈复杂度不超过 15
- 认知复杂度不超过 20
- 嵌套深度不超过 4 层

### 推荐规范 (SHOULD)

#### 泛型使用
- 用于通用数据结构: Stack, Queue, Set
- 用于集合操作: Map, Filter, Reduce
- 避免在业务逻辑中过度泛型化

#### 并发编程
- 明确 Goroutine 生命周期,使用 Context 控制退出
- 发送方负责关闭 Channel
- 使用 WaitGroup 或 errgroup 等待完成
- 限制 Goroutine 数量

#### Context 使用
- Context 作为第一个参数传递
- 不要在结构体中存储 Context
- 使用 context.WithTimeout 设置超时
- 使用 context.WithoutCancel 处理异步日志

#### 标准库新包
- 使用 slices 包: Sort, Index, Contains
- 使用 maps 包: Clone, Copy, Equal
- 使用 cmp 包: Compare, Or

## 测试规范

### 强制规范 (MUST)

#### 测试命名
- 格式: Test<组件>_<方法>_<场景>
- 示例: TestUserService_CreateUser_Success
- 测试文件以 _test.go 结尾

#### 表驱动测试
- 使用 []struct 定义测试用例
- 使用 t.Run 运行子测试
- 每个测试独立,不共享状态

#### 覆盖率要求
- 工具函数: 100%
- 领域模型: > 90%
- 服务层: > 80%

### 推荐规范 (SHOULD)

#### AAA 模式
- Arrange (准备): 设置测试数据
- Act (执行): 调用被测试函数
- Assert (断言): 验证结果

#### 使用 testify
- 使用 assert 进行断言
- 使用 require 检查前置条件

#### Mock 使用
- 使用 gomock 生成 Mock
- 使用 defer ctrl.Finish() 验证期望

## Protobuf 规范

### 强制规范 (MUST)

#### 字段命名
- 字段名使用 snake_case: user_id, created_at
- 禁止使用 camelCase: ~~userId~~, ~~createdAt~~

#### 枚举命名
- 枚举值使用 UPPER_SNAKE_CASE
- 必须有 0 值: STATUS_UNSPECIFIED = 0

#### 向后兼容
- 不要删除或重命名字段
- 不要更改字段编号
- 使用 reserved 保留已删除的字段

## gRPC 客户端规范

### 强制规范 (MUST)

#### 超时控制
- 所有外部调用必须设置超时
- 使用 context.WithTimeout 设置超时
- 推荐: 内部服务 5s, 外部 API 30s

#### Context 传递
- Context 作为第一个参数传递
- 不要在结构体中存储 Context

#### 错误处理
- 检查 gRPC 状态码
- 使用 status.Code(err) 获取错误码
- 转换为业务错误码

#### 重试策略
- 仅幂等操作可重试
- POST 操作使用幂等键
- 最多重试 3 次,使用指数退避

### 推荐规范 (SHOULD)

#### 熔断器
- 使用熔断器防止雪崩
- 配置失败阈值和恢复时间

#### 链路追踪
- 使用 OpenTelemetry 注入 Trace 信息
- 传播 Trace ID 和 Span ID

## 架构设计

### DDD 目录结构 (服务仓库)
按 **聚合根** 组织代码，而不是按技术层级。

```
internal/
├── user/                 # 聚合根: User
│   ├── domain.go         # 实体, 值对象, 仓储接口
│   ├── service.go        # 领域服务 (业务逻辑)
│   ├── repository.go     # 仓储接口 (定义在领域层)
│   └── events.go         # 领域事件
├── adapter/              # 适配器 (端口与适配器)
│   ├── grpc/             # gRPC 处理器
│   ├── http/             # HTTP 处理器
│   └── repository/       # 仓储实现 (GORM/Redis)
└── bootstrap/            # 应用初始化 & 依赖注入
```

### API 设计
- 使用 RESTful 风格
- 资源使用复数名词: /users, /orders
- 版本管理: /v1, /v2
- JSON 字段名使用 snake_case

### 响应结构
- 统一格式: code, message, data, metadata
- 错误响应: code, message, errors, request_id
- 时间格式: ISO 8601 UTC

### 微服务治理
- 实现健康检查: /health/live, /health/ready
- 实现链路追踪: 传递 Trace ID
- 实现限流和熔断

## 常见陷阱

### 配置结构体初始化
- ✅ 推荐: cfg := &ServiceConfig{}
- ❌ 不推荐: var cfg ServiceConfig; return &cfg

### 循环变量捕获 (Go 1.21 及之前)
- 必须显式捕获: user := user
- Go 1.22+ 自动捕获

### Context 传递
- ❌ 不要在结构体中存储 Context
- ✅ 作为函数参数传递

## 参考文档

详细规范请参考项目文档:
- 01_规范标准/ - 编码规范
- 02_架构设计/ - 架构设计方案

官方文档:
- Go: https://go.dev/doc/
- gRPC: https://grpc.io/docs/
- Protobuf: https://protobuf.dev/
- OpenTelemetry: https://opentelemetry.io/

## 使用 gox 基础库

统一使用 github.com/chinayin/gox 提供的能力:
- log - 日志 (基于 slog)
- config - 配置管理
- discovery - 服务发现
- trace - 链路追踪
- metrics - 指标采集
- middleware - 中间件
- transport - 传输层封装
- utils - 工具函数
