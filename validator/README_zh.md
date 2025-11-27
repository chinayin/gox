# Validator

基于 go-playground/validator 的数据验证包，提供统一的验证接口和自定义验证规则，支持多语言。

## 核心特性

- 封装官方 validator/v10，支持所有官方验证规则
- 全局单例和自定义实例
- gox 自定义规则：`snowflake_id`（示例模板）
- 多语言支持（en, zh）- 极简设计
- 线程安全
- 极简架构 - 易于扩展

## 快速开始

```go
import "github.com/chinayin/gox/validator"

type User struct {
    Name  string `validate:"required,min=2,max=50"`
    Email string `validate:"required,email"`
    Age   int    `validate:"required,min=1,max=120"`
}

func main() {
    user := User{Name: "Alice", Email: "alice@example.com", Age: 25}
    if err := validator.Validate(&user); err != nil {
        panic(err)
    }
}
```

## 多语言支持

```go
// 创建中文验证器
v := validator.New(validator.WithLocale("zh"))

// 或运行时切换语言
v.SetLocale("zh")

// 验证错误会自动翻译为中文
err := v.Validate(&user)
if err != nil {
    fmt.Println(err.Error())  // 输出中文错误消息
}
```

## 官方验证规则

完整的官方验证规则列表：https://pkg.go.dev/github.com/go-playground/validator

常用规则：`required`、`email`、`url`、`min`、`max`、`oneof`、`datetime`、`ip`、`uuid` 等

## 文档

完整文档和高级用法请参考 [doc.go](./doc.go) 或 [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/validator)

English documentation: [README.md](./README.md)

## 相关包

- [config](../config) - 配置加载
