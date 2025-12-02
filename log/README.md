# Log

Logging package based on `log/slog` interface.

## Features

- Standard interface: Uses Go standard library `log/slog` API
- Zero dependencies: Default implementation uses stdlib
- High performance: Optional zap adapter
- Simple configuration: Unified Options with constants
- K8s ready: JSON to stdout for log collection

## Quick Start

### Using stdlib (no dependencies)

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
    
    logger.Info("application started", "port", 8080)
}
```

### Using zap (high performance)

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
    
    logger.Info("application started", "port", 8080)
}
```

## Configuration

### Level Constants

```go
log.LevelDebug  // Debug level
log.LevelInfo   // Info level (default)
log.LevelWarn   // Warning level
log.LevelError  // Error level
```

### Format Constants

```go
log.FormatJSON     // JSON format (K8s standard)
log.FormatConsole  // Console format (human-readable)
```

### Output Constants

```go
log.OutputStdout   // Standard output (K8s standard)
log.OutputStderr   // Standard error
"/path/to/file"    // File path
```

## Usage Examples

### Kubernetes Deployment

```go
// Production: JSON to stdout for log collection
logger, _ := log.New(log.Options{
    Level:  log.LevelInfo,
    Format: log.FormatJSON,
    Output: log.OutputStdout,
})
```

### Local Development

```go
// Development: Console format for readability
logger, _ := log.New(log.Options{
    Level:  log.LevelDebug,
    Format: log.FormatConsole,
    Output: log.OutputStdout,
})
```

### File Output

```go
// Traditional deployment: Write to file
logger, _ := log.New(log.Options{
    Level:  log.LevelInfo,
    Format: log.FormatJSON,
    Output: "/var/log/app.log",
})
```

## API Reference

### Core Functions

```go
func New(opts Options) (*slog.Logger, error)
func NewWithHandler(handler slog.Handler) *slog.Logger
func DefaultOptions() Options
```

### Options

```go
type Options struct {
    Level  string // debug, info, warn, error
    Format string // json, console
    Output string // stdout, stderr, /path/to/file
}
```

## Architecture

```
Default:  log.New() → slog (stdlib) → no dependencies
With Zap: zaplog.New() → zapslog → zap (high-perf)
```

## Documentation

For complete documentation and advanced usage, see [doc.go](./doc.go) or [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/log)

中文文档请查看 [README_zh.md](./README_zh.md)

## Related Packages

- [log/zap](./zap) - Zap adapter for high-performance logging
