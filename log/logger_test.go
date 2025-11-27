package log

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr bool
	}{
		{
			name:    "default options",
			opts:    DefaultOptions(),
			wantErr: false,
		},
		{
			name: "debug level",
			opts: Options{
				Level:      "debug",
				OutputPath: "",
				Format:     "console",
			},
			wantErr: false,
		},
		{
			name: "json format",
			opts: Options{
				Level:      "info",
				OutputPath: "",
				Format:     "json",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("New() returned nil logger")
			}
		})
	}
}

func TestLogger_BasicLogging(t *testing.T) {
	logger, err := New(Options{
		Level:      "debug",
		OutputPath: "",
		Format:     "console",
	})
	if err != nil {
		t.Fatal(err)
	}

	// 测试基本日志记录
	logger.Debug("debug message", "key", "value")
	logger.Info("info message", "count", 42)
	logger.Warn("warn message", "flag", true)
	logger.Error("error message", "error", "test error")
}

func TestLogger_With(t *testing.T) {
	logger, err := New(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}

	// 测试 With 方法
	contextLogger := logger.With("request_id", "abc123", "user_id", 456)
	contextLogger.Info("processing request")
	contextLogger.Debug("debug info")
}
