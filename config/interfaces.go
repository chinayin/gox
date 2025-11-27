package config

// DefaultOption 默认值设置函数
// 用于在 SetDefaults 方法中设置配置的默认值
type DefaultOption func(key string, value interface{})

// Defaultable 配置可以设置默认值
// 实现此接口的配置结构体可以通过 SetDefaults 方法提供默认值
// 默认值优先级：struct tag < SetDefaults < 配置文件 < 环境变量
type Defaultable interface {
	SetDefaults(set DefaultOption)
}

// Validatable 配置可以自定义验证逻辑
// 实现此接口的配置结构体会在加载后自动执行验证
// 推荐使用 github.com/chinayin/gox/validator 包进行验证
type Validatable interface {
	Validate() error
}
