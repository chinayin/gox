# Log

基于 `log/slog` 接口的日志包，底层使用 `zap` 实现高性能日志。

## 核心特性

- 标准接口：使用 Go 标准库 `log/slog` API
- 高性能：底层使用 `zap` 实现，零分配
- 适配器模式：通过 `ZapHandler` 连接 slog 和 zap
- 支持彩色控制台输出和 JSON 格式

## 快速开始

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

## 架构

```
业务代码 → *slog.Logger (标准库) → ZapHandler (适配器) → *zap.Logger (高性能)
```

## 文档

完整文档和高级用法请参考 [doc.go](./doc.go) 或 [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/log)

English documentation: [README.md](./README.md)
