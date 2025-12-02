# CLI

为 Go 应用提供专业的命令行输出格式化。

## 特性

- 简洁、专业的启动横幅
- 自动颜色检测（遵循 NO_COLOR 环境变量）
- 命令适配器支持自动提取参数
- 流式 API，易于配置
- 可自定义输出流，便于测试

## 快速开始

```go
import "github.com/chinayin/gox/cli"

func main() {
    cli.NewStartup("MyApp", "v1.0.0").
        AddSection(
            cli.NewSection("配置").
                Add("端口", 8080).
                Add("工作线程", 4),
        ).
        AddEndpoint("健康检查", "http://localhost:8080/health").
        Print()
}
```

输出：

```
MyApp v1.0.0
--------------------------------------------------------------------------------

配置
  端口:                8080
  工作线程:            4

服务端点
  健康检查:            http://localhost:8080/health

--------------------------------------------------------------------------------
✓ 服务启动成功
  按 Ctrl+C 优雅关闭
```

## Cobra 适配器

自动提取应用信息和参数：

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
    
    // 自动提取名称、版本和已修改的参数
    cli.NewStartupWithAdapter(adapter).
        AutoAddFlags("help", "version").
        AddSection(
            cli.NewSection("配置").
                Add("端口", port).
                Add("调试模式", debug),
        ).
        Print()
    
    return nil
}
```

## API 参考

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

## 文档

完整文档和高级用法请查看 [doc.go](./doc.go) 或 [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/cli)

## 相关包

- [cli/cobra](./cobra) - Cobra 框架适配器
