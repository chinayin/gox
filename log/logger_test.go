package log

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr bool
	}{
		{
			name: "default options",
			opts: DefaultOptions(),
		},
		{
			name: "json format",
			opts: Options{
				Level:  LevelInfo,
				Format: FormatJSON,
				Output: OutputStdout,
			},
		},
		{
			name: "debug level",
			opts: Options{
				Level:  LevelDebug,
				Format: FormatConsole,
				Output: OutputStdout,
			},
		},
		{
			name: "stderr output",
			opts: Options{
				Level:  LevelInfo,
				Format: FormatJSON,
				Output: OutputStderr,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if logger == nil {
				t.Error("New() returned nil logger")
			}
		})
	}
}

func TestNewWithHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)

	logger := NewWithHandler(handler)
	if logger == nil {
		t.Fatal("NewWithHandler() returned nil")
	}

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("expected log output to contain 'test message', got: %s", output)
	}
}

func TestNewNop(t *testing.T) {
	logger := NewNop()
	if logger == nil {
		t.Fatal("NewNop() returned nil")
	}
	defer logger.Close()

	// 捕获 stdout/stderr 比较困难，因为 slog 直接写入 io.Discard
	// 这里主要验证调用不 panic，且资源管理正常
	logger.Info("this should be discarded")
	logger.Error("this should also be discarded")
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"invalid", slog.LevelInfo}, // default
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

func TestNew_WithFileOutput(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := tmpDir + "/logs/nested/app.log"

	opts := Options{
		Level:  LevelInfo,
		Format: FormatJSON,
		Output: logFile,
	}

	logger, err := New(opts)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer func() { _ = logger.Close() }()

	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	// 验证目录已创建
	dirPath := tmpDir + "/logs/nested"
	if info, err := os.Stat(dirPath); err != nil {
		t.Errorf("directory should be created: %v", err)
	} else if !info.IsDir() {
		t.Errorf("%s should be a directory", dirPath)
	}

	// 写入日志
	logger.Info("test message", "key", "value")

	// 验证文件已创建
	if _, err := os.Stat(logFile); err != nil {
		t.Errorf("log file should be created: %v", err)
	}
}

func TestLogger_Close(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := tmpDir + "/test.log"

	opts := Options{
		Level:  LevelInfo,
		Format: FormatJSON,
		Output: logFile,
	}

	logger, err := New(opts)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// 写入日志
	logger.Info("test message")

	// 关闭 logger
	if err := logger.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(logFile); err != nil {
		t.Errorf("log file should exist after Close(): %v", err)
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
			output:    OutputStdout,
			shouldErr: false,
		},
		{
			name:      "stderr",
			output:    OutputStderr,
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
			err := EnsureOutputDir(tt.output)
			if (err != nil) != tt.shouldErr {
				t.Errorf("expected error: %v, got: %v", tt.shouldErr, err)
			}
		})
	}
}
