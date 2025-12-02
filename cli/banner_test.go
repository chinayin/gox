package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestStartup(t *testing.T) {
	var buf bytes.Buffer

	startup := NewStartup("TestApp", "v1.0.0").
		WithWriter(&buf).
		AddSection(
			NewSection("Configuration").
				Add("Total", 10).
				Add("Enabled", 8),
		).
		AddEndpoint("Health", "http://localhost:8080/health").
		AddEndpoint("Metrics", "http://localhost:8080/metrics")

	startup.Print()

	output := buf.String()

	// Check banner
	if !strings.Contains(output, "TestApp (v1.0.0)") {
		t.Errorf("expected app name and version in output, got: %s", output)
	}
	if !strings.Contains(output, "Command:") {
		t.Errorf("expected Command line in output")
	}

	// Check section
	if !strings.Contains(output, "Configuration") {
		t.Errorf("expected Configuration section in output")
	}
	if !strings.Contains(output, "Total") {
		t.Errorf("expected Total in output")
	}

	// Check endpoints
	if !strings.Contains(output, "Server Endpoints") {
		t.Errorf("expected Server Endpoints in output")
	}
	if !strings.Contains(output, "http://localhost:8080/health") {
		t.Errorf("expected health endpoint in output")
	}

	// Check footer
	if !strings.Contains(output, "Server started successfully") {
		t.Errorf("expected success message in output")
	}
	if !strings.Contains(output, "Ctrl+C") {
		t.Errorf("expected shutdown instruction in output")
	}
}

func TestSection(t *testing.T) {
	section := NewSection("Configuration").
		Add("Total", 10).
		Add("Enabled", 8).
		Add("Disabled", 2)

	if section.Title != "Configuration" {
		t.Errorf("expected title 'Configuration', got: %s", section.Title)
	}
	if len(section.Items) != 3 {
		t.Errorf("expected 3 items, got: %d", len(section.Items))
	}
	if section.Items[0].Key != "Total" {
		t.Errorf("expected first key 'Total', got: %s", section.Items[0].Key)
	}
	if section.Items[0].Value != 10 {
		t.Errorf("expected first value 10, got: %v", section.Items[0].Value)
	}
}

func TestEndpoint(t *testing.T) {
	ep := Endpoint{Name: "Health", URL: "http://localhost:8080/health"}

	if ep.Name != "Health" {
		t.Errorf("expected name 'Health', got: %s", ep.Name)
	}
	if ep.URL != "http://localhost:8080/health" {
		t.Errorf("expected URL 'http://localhost:8080/health', got: %s", ep.URL)
	}
}

func TestFormatFlagValue(t *testing.T) {
	tests := []struct {
		name     string
		info     FlagInfo
		expected string
	}{
		{
			name: "bool true",
			info: FlagInfo{
				Type:  "bool",
				Value: "true",
			},
			expected: "enabled",
		},
		{
			name: "bool false",
			info: FlagInfo{
				Type:  "bool",
				Value: "false",
			},
			expected: "disabled",
		},
		{
			name: "string with default",
			info: FlagInfo{
				Type:         "string",
				Value:        "test",
				DefaultValue: "test",
			},
			expected: "test (default)",
		},
		{
			name: "string changed",
			info: FlagInfo{
				Type:         "string",
				Value:        "changed",
				DefaultValue: "default",
			},
			expected: "changed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFlagValue(tt.info)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSupportsColor(t *testing.T) {
	// Test NO_COLOR environment variable
	t.Setenv("NO_COLOR", "1")
	if supportsColor() {
		t.Error("expected supportsColor to return false when NO_COLOR is set")
	}
}

func TestNewStartupWithAdapter(t *testing.T) {
	adapter := &mockAdapter{
		name:    "TestApp",
		version: "2.0.0",
		flags:   map[string]FlagInfo{},
	}

	var buf bytes.Buffer
	startup := NewStartupWithAdapter(adapter).
		WithWriter(&buf).
		AddSection(
			NewSection("Config").
				Add("Port", 8080),
		)

	startup.Print()

	output := buf.String()

	// 验证从 adapter 获取的名称和版本
	if !strings.Contains(output, "TestApp (2.0.0)") {
		t.Errorf("expected 'TestApp (2.0.0)' in output, got: %s", output)
	}

	// 验证 adapter 已设置
	if startup.adapter != adapter {
		t.Error("expected adapter to be set")
	}
}

// mockAdapter 用于测试
type mockAdapter struct {
	name    string
	version string
	flags   map[string]FlagInfo
}

func (m *mockAdapter) GetName() string {
	return m.name
}

func (m *mockAdapter) GetVersion() string {
	return m.version
}

func (m *mockAdapter) GetFlags() map[string]FlagInfo {
	return m.flags
}
