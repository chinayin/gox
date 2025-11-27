# Log

Logging package based on `log/slog` interface, implemented with `zap` for high performance.

## Features

- Standard interface: Uses Go standard library `log/slog` API
- High performance: Implemented with `zap`, zero allocation
- Adapter pattern: Connects slog and zap through `ZapHandler`
- Supports colored console output and JSON format

## Quick Start

```go
import "github.com/chinayin/gox/log"

func main() {
    logger, err := log.New(log.Options{
        Level:  "info",
        Format: "console",
    })
    if err != nil {
        panic(err)
    }

    logger.Info("application started", "port", 8080)
    logger.Error("error occurred", "error", err)
}
```

## Architecture

```
Business Code → *slog.Logger (stdlib) → ZapHandler (adapter) → *zap.Logger (high-perf)
```

## Documentation

For complete documentation and advanced usage, see [doc.go](./doc.go) or [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/log)

中文文档请查看 [README_zh.md](./README_zh.md)
