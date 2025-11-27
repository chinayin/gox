# Validator

Data validation package based on go-playground/validator, providing unified validation interface and custom validation rules with i18n support.

## Features

- Wraps official validator/v10, supports all official validation rules
- Global singleton and custom instances
- gox custom rules: `snowflake_id` (example template)
- Multi-language support (en, zh) - extremely simple design
- Thread-safe
- Minimalist architecture - easy to extend

## Quick Start

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

## Multi-language Support

```go
// Create validator with Chinese locale
v := validator.New(validator.WithLocale("zh"))

// Or switch locale at runtime
v.SetLocale("zh")

// Validation errors will be automatically translated
err := v.Validate(&user)
if err != nil {
    fmt.Println(err.Error())  // Output in Chinese
}
```

## Official Validation Rules

Complete list of official validation rules: https://pkg.go.dev/github.com/go-playground/validator

Common rules: `required`, `email`, `url`, `min`, `max`, `oneof`, `datetime`, `ip`, `uuid`, etc.

## Documentation

For complete documentation and advanced usage, see [doc.go](./doc.go) or [pkg.go.dev](https://pkg.go.dev/github.com/chinayin/gox/validator)

中文文档请查看 [README_zh.md](./README_zh.md)

## Related Packages

- [config](../config) - Configuration loading
