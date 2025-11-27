package config

// Option 配置加载器选项
type Option func(*Loader)

// WithoutEnv 禁用环境变量自动读取
// 用于特殊场景，避免环境变量污染配置
func WithoutEnv() Option {
	return func(l *Loader) {
		l.disableEnv = true
	}
}

// WithEnvPrefix 设置环境变量前缀
// 例如：WithEnvPrefix("APP") 会将 app.port 映射到 APP_APP_PORT
func WithEnvPrefix(prefix string) Option {
	return func(l *Loader) {
		l.envPrefix = prefix
	}
}
