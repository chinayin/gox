package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Loader 统一的配置加载器（基于 Viper）
type Loader struct {
	v          *viper.Viper
	disableEnv bool
	envPrefix  string
}

// NewLoader 创建配置加载器
// 默认支持环境变量自动读取，使用 WithoutEnv() 选项可禁用
func NewLoader(opts ...Option) *Loader {
	l := &Loader{
		v:          viper.New(),
		disableEnv: false,
		envPrefix:  "",
	}

	// 应用选项
	for _, opt := range opts {
		opt(l)
	}

	// 配置 viper
	l.v.SetConfigType("yaml")

	// 根据选项配置环境变量
	if !l.disableEnv {
		l.v.AutomaticEnv()
		l.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		if l.envPrefix != "" {
			l.v.SetEnvPrefix(l.envPrefix)
		}
	}

	return l
}

// Load 加载配置文件
// 自动处理：默认值、环境变量、.local.yaml 合并、验证
func (l *Loader) Load(path string, config any) error {
	// 1. 应用默认值（struct tag + SetDefaults）
	if err := ApplyDefaults(l.v, config); err != nil {
		return err
	}

	// 2. 加载主配置文件
	l.v.SetConfigFile(path)
	if err := l.v.ReadInConfig(); err != nil {
		return fmt.Errorf("%w: %s (%w)", ErrReadFailed, path, err)
	}

	// 3. 尝试加载 .local.yaml 覆盖配置
	if err := l.loadLocalConfig(path); err != nil {
		return err
	}

	// 4. 解析到结构体
	if err := l.v.Unmarshal(config); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshalFailed, err)
	}

	// 5. 验证配置（如果配置实现了 Validatable 接口）
	if validatable, ok := config.(Validatable); ok {
		if err := validatable.Validate(); err != nil {
			return fmt.Errorf("%w: %w", ErrValidationFailed, err)
		}
	}

	return nil
}

// LoadDirectory 加载目录下所有配置文件
// configType 应该是配置结构体的指针（例如：&Config{}）
func (l *Loader) LoadDirectory(dir string, configType any) ([]any, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w: %s (%w)", ErrReadFailed, dir, err)
	}

	configs := make([]any, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !isConfigFile(name) || isLocalConfig(name) {
			continue
		}

		cfg := createInstance(configType)
		configPath := filepath.Join(dir, name)

		// 创建新的 loader 实例，继承当前 loader 的选项
		loader := l.clone()
		if err := loader.Load(configPath, cfg); err != nil {
			return nil, fmt.Errorf("%w: %s (%w)", ErrReadFailed, configPath, err)
		}

		configs = append(configs, cfg)
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, dir)
	}

	return configs, nil
}

// Get 获取配置值
func (l *Loader) Get(key string) any {
	return l.v.Get(key)
}

// GetString 获取字符串配置
func (l *Loader) GetString(key string) string {
	return l.v.GetString(key)
}

// GetInt 获取整数配置
func (l *Loader) GetInt(key string) int {
	return l.v.GetInt(key)
}

// GetBool 获取布尔配置
func (l *Loader) GetBool(key string) bool {
	return l.v.GetBool(key)
}

// GetViper 获取内部 Viper 实例（用于高级用法）
func (l *Loader) GetViper() *viper.Viper {
	return l.v
}

// clone 克隆 loader（用于 LoadDirectory）
func (l *Loader) clone() *Loader {
	opts := []Option{}
	if l.disableEnv {
		opts = append(opts, WithoutEnv())
	}
	if l.envPrefix != "" {
		opts = append(opts, WithEnvPrefix(l.envPrefix))
	}
	return NewLoader(opts...)
}

// loadLocalConfig 加载 .local.yaml 覆盖配置
func (l *Loader) loadLocalConfig(path string) error {
	localPath := getLocalConfigPath(path)
	if !fileExists(localPath) {
		return nil
	}

	localViper := viper.New()
	localViper.SetConfigFile(localPath)
	if err := localViper.ReadInConfig(); err != nil {
		return fmt.Errorf("%w: %s (%w)", ErrReadFailed, localPath, err)
	}

	if err := l.v.MergeConfigMap(localViper.AllSettings()); err != nil {
		return fmt.Errorf("%w: %w", ErrMergeFailed, err)
	}

	return nil
}

// getLocalConfigPath 获取 local 配置文件路径
func getLocalConfigPath(configPath string) string {
	ext := filepath.Ext(configPath)
	base := strings.TrimSuffix(configPath, ext)
	return base + ".local" + ext
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// isConfigFile 检查是否是配置文件
func isConfigFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".yaml" || ext == ".yml"
}

// isLocalConfig 检查是否是 local 配置文件
func isLocalConfig(filename string) bool {
	return strings.Contains(filename, ".local.yaml") || strings.Contains(filename, ".local.yml")
}

// createInstance 创建配置实例
func createInstance(template any) any {
	t := reflect.TypeOf(template)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
