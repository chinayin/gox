# Config

Unified configuration loading package based on Viper, supporting layered defaults, environment variables, and local configuration overrides.

## Features

- Layered defaults: struct tag + SetDefaults method
- Automatic merging of `.local.yaml` local configurations
- Automatic environment variable mapping
- Zero dependency leakage (business code doesn't depend on viper)

## Quick Start

```go
import "github.com/chinayin/gox/config"

type AppConfig struct {
    Port     int    `default:"8080" validate:"required,min=1,max=65535"`
    LogLevel string `default:"info" validate:"oneof=debug info warn error"`
}

func main() {
    loader := config.NewLoader()
    var cfg AppConfig
    if err := loader.Load("config.yaml", &cfg); err != nil {
        panic(err)
    }
}
```

## Documentation

For complete documentation and advanced usage, see [doc.go](./doc.go) or [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/config)

中文文档请查看 [README_zh.md](./README_zh.md)

## Related Packages

- [validator](../validator) - Configuration validation
