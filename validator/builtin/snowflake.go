package builtin

import (
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// ValidateSnowflakeID 验证 Snowflake ID 格式
func ValidateSnowflakeID(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() > 0
	case reflect.String:
		idStr := field.String()
		if idStr == "" {
			return false
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		return err == nil && id > 0
	default:
		return false
	}
}

// snowflakeTranslations 翻译（包内私有）
var snowflakeTranslations = map[string]string{
	"en": "{0} must be a valid Snowflake ID",
	"zh": "{0}必须是有效的Snowflake ID",
}
