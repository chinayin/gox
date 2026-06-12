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
	validate *validator.Validate
	uni      *ut.UniversalTranslator

	mu         sync.RWMutex // 保护 translator 和 locale
	translator ut.Translator
	locale     string
}

// New 创建新的验证器实例，自动注册所有内置验证规则和中英文翻译。
//
// 内置规则与翻译的注册失败属于库自身不变量被破坏，
// 选项配置错误（如不支持的 locale）属于启动期配置错误，
// 两者均直接 panic（fail-fast），不会发生在正常运行期。
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

	enTrans, _ := uni.GetTranslator("en")
	zhTrans, _ := uni.GetTranslator("zh")
	v.translator = enTrans

	// 一次性注册内置规则与全部支持语言的翻译，
	// SetLocale 时只切换翻译器，不再重复注册
	if err := builtin.Register(v.validate); err != nil {
		panic(fmt.Errorf("validator: register builtin rules: %w", err))
	}
	if err := en_translations.RegisterDefaultTranslations(v.validate, enTrans); err != nil {
		panic(fmt.Errorf("validator: register en translations: %w", err))
	}
	if err := zh_translations.RegisterDefaultTranslations(v.validate, zhTrans); err != nil {
		panic(fmt.Errorf("validator: register zh translations: %w", err))
	}
	if err := builtin.RegisterTranslations(v.validate, enTrans, "en"); err != nil {
		panic(fmt.Errorf("validator: register builtin en translations: %w", err))
	}
	if err := builtin.RegisterTranslations(v.validate, zhTrans, "zh"); err != nil {
		panic(fmt.Errorf("validator: register builtin zh translations: %w", err))
	}

	// 应用选项，配置错误立即暴露
	for _, opt := range opts {
		if err := opt(v); err != nil {
			panic(fmt.Errorf("validator: apply option: %w", err))
		}
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
			translator: v.Translator(),
		}
	}

	return err
}

// RegisterValidation 注册自定义验证规则
func (v *Validator) RegisterValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

// SetLocale 设置验证器的 locale（并发安全）
// 支持的 locale：
//   - "en" - English
//   - "zh" - 简体中文
//
// 注意：locale 是实例级状态，对 Global() 全局实例调用会影响所有使用方。
// 需要并发使用多种语言时，请为每种 locale 创建独立实例：New(WithLocale("zh"))。
func (v *Validator) SetLocale(locale string) error {
	trans, found := v.uni.GetTranslator(locale)
	if !found {
		return fmt.Errorf("%w: %s", ErrLocaleNotFound, locale)
	}

	// 所有语言的翻译已在 New 中注册完毕，这里只切换翻译器
	v.mu.Lock()
	v.translator = trans
	v.locale = locale
	v.mu.Unlock()

	return nil
}

// Locale 获取当前 locale
func (v *Validator) Locale() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.locale
}

// Translator 获取当前翻译器
func (v *Validator) Translator() ut.Translator {
	v.mu.RLock()
	defer v.mu.RUnlock()
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
