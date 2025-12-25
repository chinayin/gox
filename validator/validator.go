package validator

import (
	"errors"
	"fmt"
	"sync"

	"github.com/chinayin/gox/validator/builtin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	// 全局单例验证器
	globalValidator *Validator
	once            sync.Once
)

// Validator 封装 go-playground/validator
type Validator struct {
	validate   *validator.Validate
	uni        *ut.UniversalTranslator
	translator ut.Translator
	locale     string
}

// New 创建新的验证器实例
// 自动注册所有内置验证规则和翻译
func New(opts ...Option) *Validator {
	// 初始化 universal translator
	enLocale := en.New()
	zhLocale := zh.New()
	uni := ut.New(enLocale, enLocale, zhLocale)

	v := &Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
		uni:      uni,
		locale:   "en", // 默认英文
	}

	// 获取默认翻译器（英文）
	trans, _ := uni.GetTranslator("en")
	v.translator = trans

	// 注册官方英文翻译
	_ = en_translations.RegisterDefaultTranslations(v.validate, trans)

	// 注册所有内置验证规则
	_ = builtin.Register(v.validate)

	// 注册内置规则的英文翻译
	_ = builtin.RegisterTranslations(v.validate, trans, "en")

	// 应用选项
	for _, opt := range opts {
		_ = opt(v)
	}

	return v
}

// Global 获取全局验证器实例（懒加载）
func Global() *Validator {
	once.Do(func() {
		globalValidator = New()
	})
	return globalValidator
}

// Validate 验证结构体（使用全局验证器）
func Validate(v any) error {
	return Global().Validate(v)
}

// Validate 验证结构体
// 返回的错误已经翻译为当前 locale
func (v *Validator) Validate(data any) error {
	err := v.validate.Struct(data)
	if err == nil {
		return nil
	}

	// 如果是验证错误，返回翻译后的错误
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return &TranslatedError{
			errors:     validationErrors,
			translator: v.translator,
		}
	}

	return err
}

// RegisterValidation 注册自定义验证规则
func (v *Validator) RegisterValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

// SetLocale 设置验证器的 locale
// 支持的 locale：
//   - "en" - English
//   - "zh" - 简体中文
func (v *Validator) SetLocale(locale string) error {
	trans, found := v.uni.GetTranslator(locale)
	if !found {
		return fmt.Errorf("%w: %s", ErrLocaleNotFound, locale)
	}

	v.translator = trans
	v.locale = locale

	// 根据 locale 注册相应的翻译
	switch locale {
	case "zh":
		_ = zh_translations.RegisterDefaultTranslations(v.validate, trans)
		_ = builtin.RegisterTranslations(v.validate, trans, "zh")
	case "en":
		_ = en_translations.RegisterDefaultTranslations(v.validate, trans)
		_ = builtin.RegisterTranslations(v.validate, trans, "en")
	}

	return nil
}

// Locale 获取当前 locale
func (v *Validator) Locale() string {
	return v.locale
}

// Translator 获取当前翻译器
func (v *Validator) Translator() ut.Translator {
	return v.translator
}

// TranslatedError 翻译后的验证错误
type TranslatedError struct {
	errors     validator.ValidationErrors
	translator ut.Translator
}

// Error 返回翻译后的错误消息
func (e *TranslatedError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}

	// 返回第一个错误的翻译
	return e.errors[0].Translate(e.translator)
}

// Errors 返回所有翻译后的错误消息
func (e *TranslatedError) Errors() []string {
	messages := make([]string, len(e.errors))
	for i, err := range e.errors {
		messages[i] = err.Translate(e.translator)
	}
	return messages
}

// ValidationErrors 返回原始的验证错误（用于高级用法）
func (e *TranslatedError) ValidationErrors() validator.ValidationErrors {
	return e.errors
}
