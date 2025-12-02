package cobra

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAdapter_GetName(t *testing.T) {
	tests := []struct {
		name     string
		cmdUse   string
		expected string
	}{
		{
			name:     "with use",
			cmdUse:   "myapp",
			expected: "myapp",
		},
		{
			name:     "empty use",
			cmdUse:   "",
			expected: "App",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use: tt.cmdUse,
			}
			adapter := NewAdapter(cmd)

			result := adapter.GetName()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAdapter_GetVersion(t *testing.T) {
	tests := []struct {
		name       string
		cmdVersion string
		expected   string
	}{
		{
			name:       "with version",
			cmdVersion: "1.0.0",
			expected:   "1.0.0",
		},
		{
			name:       "empty version",
			cmdVersion: "",
			expected:   "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Version: tt.cmdVersion,
			}
			adapter := NewAdapter(cmd)

			result := adapter.GetVersion()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAdapter_GetFlags(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	var port int
	var debug bool
	var host string

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Server port")
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	cmd.Flags().StringVarP(&host, "host", "H", "localhost", "Server host")

	// 模拟用户修改了 port 和 debug
	cmd.Flags().Set("port", "9000")
	cmd.Flags().Set("debug", "true")

	adapter := NewAdapter(cmd)
	flags := adapter.GetFlags()

	// 检查 flag 数量
	if len(flags) != 3 {
		t.Errorf("expected 3 flags, got %d", len(flags))
	}

	// 检查 port flag
	if portFlag, ok := flags["port"]; ok {
		if portFlag.Name != "port" {
			t.Errorf("expected name 'port', got %q", portFlag.Name)
		}
		if portFlag.Value != "9000" {
			t.Errorf("expected value '9000', got %q", portFlag.Value)
		}
		if portFlag.DefaultValue != "8080" {
			t.Errorf("expected default '8080', got %q", portFlag.DefaultValue)
		}
		if !portFlag.Changed {
			t.Error("expected Changed to be true")
		}
		if portFlag.Type != "int" {
			t.Errorf("expected type 'int', got %q", portFlag.Type)
		}
	} else {
		t.Error("port flag not found")
	}

	// 检查 debug flag
	if debugFlag, ok := flags["debug"]; ok {
		if !debugFlag.Changed {
			t.Error("expected Changed to be true")
		}
		if debugFlag.Value != "true" {
			t.Errorf("expected value 'true', got %q", debugFlag.Value)
		}
	} else {
		t.Error("debug flag not found")
	}

	// 检查 host flag（未修改）
	if hostFlag, ok := flags["host"]; ok {
		if hostFlag.Changed {
			t.Error("expected Changed to be false")
		}
		if hostFlag.Value != "localhost" {
			t.Errorf("expected value 'localhost', got %q", hostFlag.Value)
		}
	} else {
		t.Error("host flag not found")
	}
}

func TestNewAdapter(t *testing.T) {
	cmd := &cobra.Command{
		Use:     "myapp",
		Version: "1.0.0",
	}

	adapter := NewAdapter(cmd)
	if adapter == nil {
		t.Fatal("expected adapter to be non-nil")
	}

	// 验证实现了接口
	if adapter.GetName() != "myapp" {
		t.Errorf("expected name 'myapp', got %q", adapter.GetName())
	}
	if adapter.GetVersion() != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", adapter.GetVersion())
	}
}
