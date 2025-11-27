# Config

基于 Viper 的统一配置加载包，支持分层默认值、环境变量和本地配置覆盖。

## 核心特性

- 分层默认值：struct tag + SetDefaults 方法
- 自动合并 `.local.yaml` 本地配置
- 环境变量自动映射
- 零依赖泄漏（业务代码不依赖 viper）

## 快速开始

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

## 文档

完整文档和高级用法请参考 [doc.go](./doc.go) 或 [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/config)

English documentation: [README.md](./README.md)

## 相关包

- [validator](../validator) - 配置验证
