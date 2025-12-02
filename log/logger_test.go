package log

import (
	"bytes"
	"log/slog"
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
