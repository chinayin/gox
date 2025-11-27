# Builtin Validation Rules

Built-in validation rules for gox validator. Extremely simple and clean.

## Structure

```
builtin/
├── builtin.go          # Rule registration
├── snowflake.go        # Snowflake ID validation + translations
├── translations.go     # One file, one for loop, done!
└── README.md
```

## Design Philosophy

**极简主义 (Minimalism)**:
- No complex directory structures
- No unnecessary abstractions
- One translation file, one for loop
- Clean and maintainable

## Adding New Rules

### Step 1: Create rule file

```go
// builtin/yourrule.go
package builtin

import "github.com/go-playground/validator/v10"

func ValidateYourRule(fl validator.FieldLevel) bool {
    return fl.Field().String() != ""
}

var yourRuleTranslations = map[string]string{
    "en": "{0} must be valid",
    "zh": "{0}必须有效",
}
```

### Step 2: Register in builtin.go

```go
func Rules() []Rule {
    return []Rule{
        {Tag: "snowflake_id", Func: ValidateSnowflakeID, Translations: snowflakeTranslations},
        {Tag: "your_rule", Func: ValidateYourRule, Translations: yourRuleTranslations},
    }
}
```

### Done!

`translations.go` automatically handles all languages with one for loop.

## Adding New Languages

Just add to the translation map:

```go
var snowflakeTranslations = map[string]string{
    "en": "{0} must be a valid Snowflake ID",
    "zh": "{0}必须是有效的Snowflake ID",
    "ja": "{0}は有効なSnowflake IDである必要があります",  // Add here
    "ko": "{0}은(는) 유효한 Snowflake ID여야 합니다",      // Add here
}
```

Then update `validator.go` to support the new locale.

## Current Rules

- `snowflake_id` - Validates Snowflake ID (int64/string, must be > 0)

## Example

```go
type User struct {
    ID int64 `validate:"snowflake_id"`
}

v := validator.New(validator.WithLocale("zh"))
err := v.Validate(&User{ID: 0})  // "ID必须是有效的Snowflake ID"
```

## Why This Design?

1. **Simple**: No complex directory structures
2. **Clean**: One translation file handles everything
3. **Maintainable**: Easy to understand and modify
4. **Scalable**: Easy to add rules and languages
5. **No Duplication**: Translations defined once per rule
