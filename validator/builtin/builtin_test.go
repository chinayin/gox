package builtin

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestRules(t *testing.T) {
	rules := Rules()

	if len(rules) == 0 {
		t.Error("Rules() should return at least one rule")
	}

	// 验证每个规则都有 Tag、Func 和 Translations
	for _, rule := range rules {
		if rule.Tag == "" {
			t.Error("Rule tag should not be empty")
		}
		if rule.Func == nil {
			t.Errorf("Rule %s func should not be nil", rule.Tag)
		}
		if rule.Translations == nil {
			t.Errorf("Rule %s Translations should not be nil", rule.Tag)
		}
	}
}

func TestRegister(t *testing.T) {
	v := validator.New()

	if err := Register(v); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// 验证规则已注册（通过尝试使用）
	type TestStruct struct {
		SnowflakeID int64 `validate:"snowflake_id"`
	}

	// 有效数据应该通过
	valid := TestStruct{
		SnowflakeID: 1234567890123456789,
	}
	if err := v.Struct(valid); err != nil {
		t.Errorf("Valid data should pass validation: %v", err)
	}

	// 无效数据应该失败
	invalid := TestStruct{
		SnowflakeID: 0,
	}
	if err := v.Struct(invalid); err == nil {
		t.Error("Invalid data should fail validation")
	}
}
