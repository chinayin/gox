# Log

基于 `log/slog` 接口的日志包。

## 特性

- 标准接口：使用 Go 标准库 `log/slog` API
- 零依赖：默认实现使用标准库
- 高性能：可选的 zap 适配器
- 简单配置：统一的 Options 结构和常量
- K8s 就绪：JSON 输出到 stdout 用于日志收集

## 快速开始

### 使用标准库（无依赖）

```go
import "github.com/chinayin/gox/log"

func main() {
    logger, err := log.New(log.Options{
        Level:  log.LevelInfo,
        Format: log.FormatConsole,
        Output: log.OutputStdout,
    })
    if err != nil {
        panic(err)
    }
    
    logger.Info("应用启动", "端口", 8080)
}
```

### 使用 zap（高性能）

```go
import zaplog "github.com/chinayin/gox/log/zap"

func main() {
    logger, err := zaplog.New(log.Options{
        Level:  log.LevelInfo,
        Format: log.FormatJSON,
        Output: log.OutputStdout,
    })
    if err != nil {
        panic(err)
    }
    
    logger.Info("应用启动", "端口", 8080)
}
```

## 配置

### 日志级别常量

```go
log.LevelDebug  // 调试级别
log.LevelInfo   // 信息级别（默认）
log.LevelWarn   // 警告级别
log.LevelError  // 错误级别
```

### 格式常量

```go
log.FormatJSON     // JSON 格式（K8s 标准）
log.FormatConsole  // 控制台格式（人类可读）
```

### 输出常量

```go
log.OutputStdout   // 标准输出（K8s 标准）
log.OutputStderr   // 标准错误
"/path/to/file"    // 文件路径
```

## 使用示例

### Kubernetes 部署

```go
// 生产环境：JSON 输出到 stdout 用于日志收集
logger, _ := log.New(log.Options{
    Level:  log.LevelInfo,
    Format: log.FormatJSON,
    Output: log.OutputStdout,
})
```

### 本地开发

```go
// 开发环境：控制台格式便于阅读
logger, _ := log.New(log.Options{
    Level:  log.LevelDebug,
    Format: log.FormatConsole,
    Output: log.OutputStdout,
})
```

### 文件输出

```go
// 传统部署：写入文件
logger, _ := log.New(log.Options{
    Level:  log.LevelInfo,
    Format: log.FormatJSON,
    Output: "/var/log/app.log",
})
```

## API 参考

### 核心函数

```go
func New(opts Options) (*slog.Logger, error)
func NewWithHandler(handler slog.Handler) *slog.Logger
func DefaultOptions() Options
```

### 配置选项

```go
type Options struct {
    Level  string // debug, info, warn, error
    Format string // json, console
    Output string // stdout, stderr, /path/to/file
}
```

## 架构

```
默认：log.New() → slog (标准库) → 无依赖
Zap：zaplog.New() → zapslog → zap (高性能)
```

## 文档

完整文档和高级用法请查看 [doc.go](./doc.go) 或 [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/log)

## 相关包

- [log/zap](./zap) - 高性能日志的 Zap 适配器
