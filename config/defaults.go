package config

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

// ApplyDefaults 应用默认值（分层处理）
// 1. 先应用 struct tag 的默认值（使用 creasty/defaults）
// 2. 再应用 SetDefaults() 方法的默认值（会覆盖 struct tag）
func ApplyDefaults(v *viper.Viper, config interface{}) error {
	// 第一层：应用 struct tag 的默认值
	if err := defaults.Set(config); err != nil {
		return fmt.Errorf("failed to apply struct tag defaults: %w", err)
	}

	// 第二层：应用 SetDefaults() 方法的默认值
	if defaultable, ok := config.(Defaultable); ok {
		defaultable.SetDefaults(func(key string, value interface{}) {
			v.SetDefault(key, value)
		})
	}

	return nil
}
