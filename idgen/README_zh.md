# idgen

统一的 ID 生成器包。

## 特性

- **统一接口**：所有 ID 算法使用同一个 `Generator` 接口
- **Snowflake**：Twitter 雪花算法（本地落地）
- **类型访问**：通过 `Unwrap()` 获取底层类型的完整能力
- **线程安全**：所有生成器都是并发安全的

## 安装

```go
import "github.com/chinayin/gox/idgen"
```

## 快速开始

### 基础用法

```go
// 单机部署（默认 nodeID=1）
snowflake, err := idgen.NewSnowflake()
if err != nil {
    log.Fatal(err)
}

// 分布式部署（每个节点必须有唯一的 ID 0-1023）
snowflake, err := idgen.NewSnowflake(cfg.NodeID)
if err != nil {
    log.Fatal(err)
}

id := snowflake.Generate()
fmt.Println(id.Int64())   // 1234567890123456789
fmt.Println(id.String())  // "1234567890123456789"
```

### 依赖注入（推荐）

使用 `Generator` 接口进行依赖注入，方便切换实现和单元测试：

```go
// 定义服务，依赖接口而非具体实现
type OrderService struct {
    idgen idgen.Generator
}

func NewOrderService(g idgen.Generator) *OrderService {
    return &OrderService{idgen: g}
}

func (s *OrderService) CreateOrder() int64 {
    id := s.idgen.Generate()
    return id.Int64()
}

// 初始化时注入具体实现
func main() {
    snowflake, _ := idgen.NewSnowflake()
    svc := NewOrderService(snowflake)
    orderID := svc.CreateOrder()
    fmt.Println(orderID)
}
```

### 单元测试 Mock

```go
// mock_test.go
type MockGenerator struct {
    id int64
}

func (m *MockGenerator) Generate() idgen.ID {
    m.id++
    return idgen.NewID(m.id, fmt.Sprintf("%d", m.id), nil)
}

func TestOrderService_CreateOrder(t *testing.T) {
    mock := &MockGenerator{id: 1000}
    svc := NewOrderService(mock)
    
    id1 := svc.CreateOrder()
    id2 := svc.CreateOrder()
    
    if id1 != 1001 || id2 != 1002 {
        t.Errorf("unexpected IDs: %d, %d", id1, id2)
    }
}
```

### 配合 slog 日志

```go
import "log/slog"

func (s *OrderService) CreateOrder(ctx context.Context) (int64, error) {
    id := s.idgen.Generate()
    
    slog.InfoContext(ctx, "order created",
        "order_id", id.Int64(),
        "order_id_str", id.String(),
    )
    
    return id.Int64(), nil
}
```

### 访问 Snowflake 特有功能

```go
if sf, ok := id.Unwrap().(snowflake.ID); ok {
    fmt.Println(sf.Time())   // 时间戳（毫秒）
    fmt.Println(sf.Node())   // 节点 ID
    fmt.Println(sf.Step())   // 序列号
    fmt.Println(sf.Base64()) // Base64 编码
}
```

### 全局生成器（谨慎使用）

> ⚠️ **不推荐在生产环境使用。** 建议使用依赖注入。

问题：
- 使单元测试变得困难（全局状态）
- 隐藏了函数签名中的依赖
- 分布式系统中 nodeID 不唯一会导致重复 ID

```go
// 应用启动时（仅一次）
snowflake, _ := idgen.NewSnowflake(cfg.NodeID)
idgen.SetDefault(snowflake)

// 任意位置使用
id := idgen.Generate()
```

## API

### 类型

```go
type ID struct { ... }
func (id ID) Int64() int64
func (id ID) String() string
func (id ID) IsZero() bool
func (id ID) Unwrap() any
```

### 接口

```go
type Generator interface {
    Generate() ID
}
```

### 函数

```go
func NewSnowflake(nodeID ...int64) (*Snowflake, error)  // 可选 nodeID，默认为 1
func SetDefault(g Generator) error                      // 设置全局（谨慎使用）
func Default() Generator
func Generate() ID
```

## Snowflake ID 结构

```
+--------------------------------------------------------------------------+
| 1 Bit 未使用 | 41 Bit 时间戳 | 10 Bit 节点ID | 12 Bit 序列号              |
+--------------------------------------------------------------------------+
```

- **时间戳**：从 epoch（2010年11月4日）开始的毫秒数，可用约 69 年
- **节点ID**：0-1023（10 位）- **分布式系统中每个节点必须唯一**
- **序列号**：每毫秒每节点 0-4095（12 位）

## 许可证

本包使用的 Snowflake 实现基于 [github.com/bwmarrin/snowflake](https://github.com/bwmarrin/snowflake)（BSD 2-Clause 许可证）。
