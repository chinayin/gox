# CLI

Professional command-line output formatting for Go applications.

## Features

- Clean, professional startup banners
- Automatic color detection (respects NO_COLOR)
- Command adapter support for auto-extracting parameters
- Fluent API for easy configuration
- Customizable output writer for testing

## Quick Start

```go
import "github.com/chinayin/gox/cli"

func main() {
    cli.NewStartup("MyApp", "v1.0.0").
        AddSection(
            cli.NewSection("Configuration").
                Add("Port", 8080).
                Add("Workers", 4),
        ).
        AddEndpoint("Health", "http://localhost:8080/health").
        Print()
}
```

Output:

```
MyApp v1.0.0
--------------------------------------------------------------------------------

Configuration
  Port:                8080
  Workers:             4

Server Endpoints
  Health:              http://localhost:8080/health

--------------------------------------------------------------------------------
✓ Server started successfully
  Press Ctrl+C to shutdown gracefully
```

## Cobra Adapter

Automatically extract application info and parameters:

```go
import (
    "github.com/chinayin/gox/cli"
    clicobra "github.com/chinayin/gox/cli/cobra"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:     "myapp",
    Version: "1.0.0",
    RunE:    run,
}

func run(cmd *cobra.Command, args []string) error {
    adapter := clicobra.NewAdapter(cmd)
    
    // Automatically extracts name, version, and changed flags
    cli.NewStartupWithAdapter(adapter).
        AutoAddFlags("help", "version").
        AddSection(
            cli.NewSection("Configuration").
                Add("Port", port).
                Add("Debug", debug),
        ).
        Print()
    
    return nil
}
```

## API Reference

### Startup

```go
func NewStartup(name, version string) *Startup
func NewStartupWithAdapter(adapter CommandAdapter) *Startup
func (s *Startup) WithWriter(w io.Writer) *Startup
func (s *Startup) WithAdapter(adapter CommandAdapter) *Startup
func (s *Startup) AutoAddFlags(excludeNames ...string) *Startup
func (s *Startup) AddSection(section *Section) *Startup
func (s *Startup) AddEndpoint(name, url string) *Startup
func (s *Startup) Print()
```

### Section

```go
func NewSection(title string) *Section
func (s *Section) Add(key string, value any) *Section
```

### CommandAdapter

```go
type CommandAdapter interface {
    GetName() string
    GetVersion() string
    GetFlags() map[string]FlagInfo
}
```

## Documentation

For complete documentation and advanced usage, see [doc.go](./doc.go) or [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/cli)

中文文档请查看 [README_zh.md](./README_zh.md)

## Related Packages

- [cli/cobra](./cobra) - Cobra framework adapter
