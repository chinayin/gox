package builtin

import (
	"errors"
	"testing"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func TestValidateSnowflakeID_Kinds(t *testing.T) {
	v := validator.New()
	if err := Register(v); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	type intCase struct {
		ID int `validate:"snowflake_id"`
	}
	type int8Case struct {
		ID int8 `validate:"snowflake_id"`
	}
	type int32Case struct {
		ID int32 `validate:"snowflake_id"`
	}
	type int64Case struct {
		ID int64 `validate:"snowflake_id"`
	}
	type uintCase struct {
		ID uint `validate:"snowflake_id"`
	}
	type uint64Case struct {
		ID uint64 `validate:"snowflake_id"`
	}
	type stringCase struct {
		ID string `validate:"snowflake_id"`
	}
	type floatCase struct {
		ID float64 `validate:"snowflake_id"`
	}
	type boolCase struct {
		ID bool `validate:"snowflake_id"`
	}

	tests := []struct {
		name    string
		data    any
		wantErr bool
	}{
		{"int positive", intCase{ID: 123}, false},
		{"int zero", intCase{ID: 0}, true},
		{"int negative", intCase{ID: -1}, true},
		{"int8 positive", int8Case{ID: 1}, false},
		{"int32 positive", int32Case{ID: 100}, false},
		{"int64 positive", int64Case{ID: 1234567890123456789}, false},
		{"int64 zero", int64Case{ID: 0}, true},
		{"uint positive", uintCase{ID: 1}, false},
		{"uint zero", uintCase{ID: 0}, true},
		{"uint64 positive", uint64Case{ID: 9876543210}, false},
		{"string numeric", stringCase{ID: "1234567890123456789"}, false},
		{"string empty", stringCase{ID: ""}, true},
		{"string non-numeric", stringCase{ID: "abc"}, true},
		{"string zero", stringCase{ID: "0"}, true},
		{"string negative", stringCase{ID: "-5"}, true},
		{"string overflow int64", stringCase{ID: "99999999999999999999999"}, true},
		{"float unsupported kind", floatCase{ID: 1.0}, true},
		{"bool unsupported kind", boolCase{ID: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Struct(%+v) error = %v, wantErr %v", tt.data, err, tt.wantErr)
			}
		})
	}
}

func TestRegisterTranslations(t *testing.T) {
	tests := []struct {
		name        string
		locale      string
		wantMessage string
	}{
		{"en translation", "en", "ID must be a valid Snowflake ID"},
		{"zh translation", "zh", "ID必须是有效的Snowflake ID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			if err := Register(v); err != nil {
				t.Fatalf("Register() error = %v", err)
			}

			enLocale := en.New()
			zhLocale := zh.New()
			uni := ut.New(enLocale, enLocale, zhLocale)
			trans, found := uni.GetTranslator(tt.locale)
			if !found {
				t.Fatalf("translator %q not found", tt.locale)
			}

			if err := RegisterTranslations(v, trans, tt.locale); err != nil {
				t.Fatalf("RegisterTranslations() error = %v", err)
			}

			type testStruct struct {
				ID int64 `validate:"snowflake_id"`
			}
			err := v.Struct(testStruct{ID: 0})
			if err == nil {
				t.Fatal("expected validation error")
			}

			var verrs validator.ValidationErrors
			if !errors.As(err, &verrs) {
				t.Fatalf("expected ValidationErrors, got %T", err)
			}
			if got := verrs[0].Translate(trans); got != tt.wantMessage {
				t.Errorf("translated message = %q, want %q", got, tt.wantMessage)
			}
		})
	}
}

func TestRegisterTranslations_DuplicateRegistrationFails(t *testing.T) {
	v := validator.New()
	if err := Register(v); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	if err := RegisterTranslations(v, trans, "en"); err != nil {
		t.Fatalf("first RegisterTranslations() error = %v", err)
	}
	// 同一翻译器上重复注册（override=false）应返回错误而非吞掉
	if err := RegisterTranslations(v, trans, "en"); err == nil {
		t.Error("duplicate RegisterTranslations() should return error")
	}
}

func TestRegisterTranslations_UnknownLocaleIsNoop(t *testing.T) {
	v := validator.New()
	if err := Register(v); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	// 规则没有 fr 翻译，应跳过且不报错
	if err := RegisterTranslations(v, trans, "fr"); err != nil {
		t.Errorf("RegisterTranslations() with unknown locale should be no-op, got error: %v", err)
	}
}

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
