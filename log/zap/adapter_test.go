package zap

import (
	"os"
	"testing"

	"github.com/chinayin/gox/log"
	"go.uber.org/zap/zapcore"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name    string
		opts    log.Options
		wantErr bool
	}{
		{
			name: "default options",
			opts: log.DefaultOptions(),
		},
		{
			name: "json format",
			opts: log.Options{
				Level:  log.LevelInfo,
				Format: log.FormatJSON,
				Output: log.OutputStdout,
			},
		},
		{
			name: "console format",
			opts: log.Options{
				Level:  log.LevelDebug,
				Format: log.FormatConsole,
				Output: log.OutputStdout,
			},
		},
		{
			name: "stderr output",
			opts: log.Options{
				Level:  log.LevelInfo,
				Format: log.FormatJSON,
				Output: log.OutputStderr,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := NewHandler(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if handler == nil {
				t.Error("NewHandler() returned nil")
			}
		})
	}
}

func TestNew(t *testing.T) {
	logger, err := New(log.Options{
		Level:  log.LevelInfo,
		Format: log.FormatConsole,
		Output: log.OutputStdout,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	// Test logging
	logger.Info("test message", "key", "value")
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected zapcore.Level
	}{
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"invalid", zapcore.InfoLevel}, // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEnsureOutputDir(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		shouldErr bool
	}{
		{
			name:      "stdout",
			output:    log.OutputStdout,
			shouldErr: false,
		},
		{
			name:      "stderr",
			output:    log.OutputStderr,
			shouldErr: false,
		},
		{
			name:      "file with nested dirs",
			output:    t.TempDir() + "/logs/app/test.log",
			shouldErr: false,
		},
		{
			name:      "current dir",
			output:    "test.log",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := log.EnsureOutputDir(tt.output)
			if (err != nil) != tt.shouldErr {
				t.Errorf("expected error: %v, got: %v", tt.shouldErr, err)
			}
		})
	}
}

func TestNewHandler_WithFileOutput(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := tmpDir + "/logs/nested/app.log"

	opts := log.Options{
		Level:  log.LevelInfo,
		Format: log.FormatJSON,
		Output: logFile,
	}

	handler, err := NewHandler(opts)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	if handler == nil {
		t.Fatal("handler should not be nil")
	}

	// 验证目录已创建
	dirPath := tmpDir + "/logs/nested"
	if info, err := os.Stat(dirPath); err != nil {
		t.Errorf("directory should be created: %v", err)
	} else if !info.IsDir() {
		t.Errorf("%s should be a directory", dirPath)
	}
}
