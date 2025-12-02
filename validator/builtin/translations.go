package builtin

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// RegisterTranslations 注册指定语言的翻译
// 一个函数，一个 for 循环，搞定所有语言
func RegisterTranslations(v *validator.Validate, trans ut.Translator, locale string) error {
	for _, rule := range Rules() {
		if translation, ok := rule.Translations[locale]; ok {
			if err := registerTranslation(v, trans, rule.Tag, translation); err != nil {
				return err
			}
		}
	}
	return nil
}

// registerTranslation 注册单个翻译
func registerTranslation(v *validator.Validate, trans ut.Translator, tag, translation string) error {
	return v.RegisterTranslation(
		tag,
		trans,
		func(ut ut.Translator) error {
			return ut.Add(tag, translation, false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(fe.Tag(), fe.Field())
			return t
		},
	)
}
