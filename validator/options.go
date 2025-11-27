package validator

// Option 配置选项
type Option func(*Validator) error

// WithLocale 设置验证器的 locale（语言）
// 支持的 locale：
//   - "en" - English (default)
//   - "zh" - 简体中文
func WithLocale(locale string) Option {
	return func(v *Validator) error {
		return v.SetLocale(locale)
	}
}
