package config

import (
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	Port     int    `default:"8080" mapstructure:"port" validate:"required,min=1,max=65535"`
	LogLevel string `default:"info" mapstructure:"log_level" validate:"oneof=debug info warn error"`
	Name     string `mapstructure:"name"`
}

func (c *testConfig) SetDefaults(_ DefaultOption) {
	// 可以覆盖 struct tag 的默认值
	// 这里保持默认值不变，仅作为示例
}

func (c *testConfig) Validate() error {
	// 简单验证逻辑（实际项目中应使用 validator 包）
	if c.Port < 1 || c.Port > 65535 {
		return os.ErrInvalid
	}
	return nil
}

func TestNewLoader(t *testing.T) {
	loader := NewLoader()
	if loader == nil {
		t.Fatal("NewLoader() returned nil")
	}
	if loader.v == nil {
		t.Error("Loader.v is nil")
	}
	if loader.disableEnv {
		t.Error("disableEnv should be false by default")
	}
}

func TestNewLoader_WithOptions(t *testing.T) {
	t.Run("WithoutEnv", func(t *testing.T) {
		loader := NewLoader(WithoutEnv())
		if !loader.disableEnv {
			t.Error("disableEnv should be true")
		}
	})

	t.Run("WithEnvPrefix", func(t *testing.T) {
		loader := NewLoader(WithEnvPrefix("TEST"))
		if loader.envPrefix != "TEST" {
			t.Errorf("envPrefix = %s, want TEST", loader.envPrefix)
		}
	})

	t.Run("Multiple options", func(t *testing.T) {
		loader := NewLoader(WithoutEnv(), WithEnvPrefix("APP"))
		if !loader.disableEnv {
			t.Error("disableEnv should be true")
		}
		if loader.envPrefix != "APP" {
			t.Errorf("envPrefix = %s, want APP", loader.envPrefix)
		}
	})
}

func TestLoader_Load(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
port: 9090
log_level: debug
name: test-app
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}

	loader := NewLoader()
	var cfg testConfig
	if err := loader.Load(configPath, &cfg); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %s, want debug", cfg.LogLevel)
	}
	if cfg.Name != "test-app" {
		t.Errorf("Name = %s, want test-app", cfg.Name)
	}
}

func TestLoader_Load_WithDefaults(t *testing.T) {
	// 创建临时配置文件（不包含 port 和 log_level）
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
name: test-app
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}

	loader := NewLoader()
	var cfg testConfig
	if err := loader.Load(configPath, &cfg); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// 应该使用默认值
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080 (default)", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %s, want info (default)", cfg.LogLevel)
	}
}

func TestLoader_Load_WithLocalConfig(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	localConfigPath := filepath.Join(tmpDir, "config.local.yaml")

	configContent := `
port: 8080
log_level: info
name: prod-app
`
	localConfigContent := `
port: 9090
name: dev-app
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}
	if err := os.WriteFile(localConfigPath, []byte(localConfigContent), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}

	loader := NewLoader()
	var cfg testConfig
	if err := loader.Load(configPath, &cfg); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// local 配置应该覆盖主配置
	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090 (from local)", cfg.Port)
	}
	if cfg.Name != "dev-app" {
		t.Errorf("Name = %s, want dev-app (from local)", cfg.Name)
	}
	// log_level 没有在 local 中定义，应该使用主配置的值
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %s, want info (from main)", cfg.LogLevel)
	}
}

func TestLoader_LoadDirectory(t *testing.T) {
	// 创建临时目录和多个配置文件
	tmpDir := t.TempDir()

	configs := map[string]string{
		"app1.yaml": `
port: 8081
log_level: debug
name: app1
`,
		"app2.yaml": `
port: 8082
log_level: info
name: app2
`,
		"app3.local.yaml": `
port: 8083
`, // 这个文件应该被忽略
	}

	for name, content := range configs {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil { //nolint:gosec
			t.Fatal(err)
		}
	}

	loader := NewLoader()
	results, err := loader.LoadDirectory(tmpDir, &testConfig{})
	if err != nil {
		t.Fatalf("LoadDirectory() error = %v", err)
	}

	// 应该只加载 2 个配置文件（忽略 .local.yaml）
	if len(results) != 2 {
		t.Errorf("LoadDirectory() returned %d configs, want 2", len(results))
	}

	// 验证配置内容
	names := make(map[string]bool)
	for _, result := range results {
		cfg, _ := result.(*testConfig)
		names[cfg.Name] = true
	}

	if !names["app1"] || !names["app2"] {
		t.Error("LoadDirectory() did not load expected configs")
	}
}

func TestLoader_GetMethods(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
port: 9090
log_level: debug
enabled: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}

	loader := NewLoader()
	var cfg testConfig
	if err := loader.Load(configPath, &cfg); err != nil {
		t.Fatal(err)
	}

	if loader.GetInt("port") != 9090 {
		t.Errorf("GetInt(port) = %d, want 9090", loader.GetInt("port"))
	}
	if loader.GetString("log_level") != "debug" {
		t.Errorf("GetString(log_level) = %s, want debug", loader.GetString("log_level"))
	}
	if !loader.GetBool("enabled") {
		t.Error("GetBool(enabled) = false, want true")
	}
}
