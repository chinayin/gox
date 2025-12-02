package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// Logger 包装 slog.Logger 并管理资源
// 实现 io.Closer 接口
type Logger struct {
	*slog.Logger
	cleanup func()
}

// 确保 Logger 实现 io.Closer 接口
var _ io.Closer = (*Logger)(nil)

// Close 释放日志资源（如文件句柄）
func (l *Logger) Close() error {
	if l.cleanup != nil {
		l.cleanup()
	}
	return nil
}

// New 创建 Logger，使用标准库实现
// 返回的 Logger 需要在应用退出时调用 Close() 释放资源
func New(opts Options) (*Logger, error) {
	// 解析级别
	level := parseLevel(opts.Level)

	// 配置 HandlerOptions
	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	// 选择输出
	writer, cleanup, err := getWriter(opts.Output)
	if err != nil {
		return nil, err
	}

	// 选择格式
	var handler slog.Handler
	if opts.Format == FormatJSON {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(writer, handlerOpts)
	}

	return &Logger{
		Logger:  slog.New(handler),
		cleanup: cleanup,
	}, nil
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
// 返回 writer, cleanup 函数, error
func getWriter(output string) (io.Writer, func(), error) {
	switch output {
	case OutputStdout, "":
		return os.Stdout, func() {}, nil
	case OutputStderr:
		return os.Stderr, func() {}, nil
	default:
		// 确保目录存在
		if err := EnsureOutputDir(output); err != nil {
			return nil, nil, err
		}

		// 打开文件
		// #nosec G304 -- output 来自配置文件，由用户控制
		f, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open log file: %w", err)
		}

		// 返回文件和 cleanup 函数
		cleanup := func() {
			_ = f.Close()
		}
		return f, cleanup, nil
	}
}

// parseLevel 解析日志级别
func parseLevel(level string) slog.Level {
	if lvl, ok := slogLevelMap[level]; ok {
		return lvl
	}
	return slog.LevelInfo
}

// isFileOutput 判断是否为文件输出
func isFileOutput(output string) bool {
	return output != OutputStdout && output != OutputStderr && output != ""
}

// EnsureOutputDir 确保日志文件的目录存在
func EnsureOutputDir(output string) error {
	if !isFileOutput(output) {
		return nil
	}
	return os.MkdirAll(filepath.Dir(output), 0o750)
}
