// Package builtin provides built-in validator functions.
package builtin

import (
	"github.com/go-playground/validator/v10"
)

// Rule 定义验证规则
type Rule struct {
	Tag          string
	Func         validator.Func
	Translations map[string]string
}

// Rules 返回所有内置验证规则
func Rules() []Rule {
	return []Rule{
		{
			Tag:          "snowflake_id",
			Func:         ValidateSnowflakeID,
			Translations: snowflakeTranslations,
		},
		// 未来扩展：
		// {
		//     Tag:          "your_rule",
		//     Func:         ValidateYourRule,
		//     Translations: yourRuleTranslations,
		// },
	}
}

// Register 注册所有内置验证规则
func Register(v *validator.Validate) error {
	for _, rule := range Rules() {
		if err := v.RegisterValidation(rule.Tag, rule.Func); err != nil {
			return err
		}
	}
	return nil
}
