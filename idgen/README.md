# idgen

Unified ID generation package for Go applications.

## Features

- **Unified Interface**: Single `Generator` interface for all ID algorithms
- **Snowflake**: Twitter Snowflake algorithm (vendored locally)
- **Type Access**: Get underlying type via `Unwrap()` for algorithm-specific features
- **Thread-Safe**: All generators are safe for concurrent use

## Installation

```go
import "github.com/chinayin/gox/idgen"
```

## Quick Start

### Basic Usage

```go
// Single-node deployment (default nodeID=1)
snowflake, err := idgen.NewSnowflake()
if err != nil {
    log.Fatal(err)
}

// Distributed deployment (each node MUST have unique ID 0-1023)
snowflake, err := idgen.NewSnowflake(cfg.NodeID)
if err != nil {
    log.Fatal(err)
}

id := snowflake.Generate()
fmt.Println(id.Int64())   // 1234567890123456789
fmt.Println(id.String())  // "1234567890123456789"
```

### Dependency Injection (Recommended)

Use the `Generator` interface for dependency injection:

```go
// Define service with interface dependency
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

// Inject concrete implementation at startup
func main() {
    snowflake, _ := idgen.NewSnowflake()
    svc := NewOrderService(snowflake)
    orderID := svc.CreateOrder()
    fmt.Println(orderID)
}
```

### Unit Testing with Mock

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

### Logging with slog

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

### Accessing Snowflake-Specific Features

```go
if sf, ok := id.Unwrap().(snowflake.ID); ok {
    fmt.Println(sf.Time())   // timestamp in milliseconds
    fmt.Println(sf.Node())   // node ID
    fmt.Println(sf.Step())   // sequence number
    fmt.Println(sf.Base64()) // base64 encoding
}
```

### Global Generator (Use with Caution)

> ⚠️ **Not recommended for production.** Prefer dependency injection.

Problems:
- Makes unit testing harder (global state)
- Hides dependencies in function signatures
- Risk of duplicate IDs if nodeID is not unique

```go
// At startup (ONCE only)
snowflake, _ := idgen.NewSnowflake(cfg.NodeID)
idgen.SetDefault(snowflake)

// Anywhere in your code
id := idgen.Generate()
```

## API

### Types

```go
type ID struct { ... }
func (id ID) Int64() int64
func (id ID) String() string
func (id ID) IsZero() bool
func (id ID) Unwrap() any
```

### Interfaces

```go
type Generator interface {
    Generate() ID
}
```

### Functions

```go
func NewSnowflake(nodeID ...int64) (*Snowflake, error)  // optional nodeID, default=1
func SetDefault(g Generator) error                      // set global (use with caution)
func Default() Generator
func Generate() ID
```

## Snowflake ID Structure

```
+--------------------------------------------------------------------------+
| 1 Bit Unused | 41 Bit Timestamp | 10 Bit NodeID | 12 Bit Sequence ID     |
+--------------------------------------------------------------------------+
```

- **Timestamp**: Milliseconds since epoch (Nov 04 2010), ~69 years capacity
- **NodeID**: 0-1023 (10 bits) - **MUST be unique per node in distributed systems**
- **Sequence**: 0-4095 per millisecond per node (12 bits)

## License

This package uses the Snowflake implementation based on [github.com/bwmarrin/snowflake](https://github.com/bwmarrin/snowflake) (BSD 2-Clause License).
