package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// New 创建 slog.Logger，使用标准库实现
// 这是默认实现，无外部依赖
func New(opts Options) (*slog.Logger, error) {
	// 解析级别
	level := parseLevel(opts.Level)

	// 配置 HandlerOptions
	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	// 选择输出
	writer, err := getWriter(opts.Output)
	if err != nil {
		return nil, err
	}

	// 选择格式
	handler := slog.Handler(nil)
	if opts.Format == FormatJSON {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(writer, handlerOpts)
	}

	return slog.New(handler), nil
}

// NewWithHandler 使用自定义 Handler 创建 Logger
// 类似 cli.NewStartupWithAdapter
func NewWithHandler(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}

// slogLevelMap 日志级别映射表
var slogLevelMap = map[string]slog.Level{
	LevelDebug: slog.LevelDebug,
	LevelInfo:  slog.LevelInfo,
	LevelWarn:  slog.LevelWarn,
	LevelError: slog.LevelError,
}

// getWriter 根据 Output 获取输出目标
func getWriter(output string) (io.Writer, error) {
	switch output {
	case OutputStdout, "":
		return os.Stdout, nil
	case OutputStderr:
		return os.Stderr, nil
	default:
		// 文件路径
		f, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		return f, nil
	}
}

// parseLevel 解析日志级别
func parseLevel(level string) slog.Level {
	if lvl, ok := slogLevelMap[level]; ok {
		return lvl
	}
	return slog.LevelInfo
}
