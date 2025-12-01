package log

import (
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options 日志配置选项
type Options struct {
	Level      string // 日志级别: debug, info, warn, error
	OutputPath string // 输出路径，空字符串表示输出到 stdout
	Format     string // 格式: json, console
	Color      string // 颜色: auto(默认), always, never, no
}

// DefaultOptions 返回默认配置
func DefaultOptions() Options {
	return Options{
		Level:      "info",
		OutputPath: "",
		Format:     "console",
		Color:      "auto",
	}
}

// New 创建新的 slog.Logger，底层使用 zap 实现
func New(opts Options) (*slog.Logger, error) {
	// 解析日志级别
	lvl, err := parseLevel(opts.Level)
	if err != nil {
		return nil, err
	}

	// 配置编码器
	cfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出
	ws, err := getWriteSyncer(opts)
	if err != nil {
		return nil, err
	}

	// 配置编码器
	enc := getEncoder(opts, cfg)

	// 创建 zap core
	core := zapcore.NewCore(enc, ws, lvl)

	// 创建 zap logger
	zl := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))

	// 创建 slog Handler
	h := NewZapHandler(zl, lvl)

	// 创建 slog.Logger
	return slog.New(h), nil
}

func getWriteSyncer(opts Options) (zapcore.WriteSyncer, error) {
	if opts.OutputPath == "" {
		return zapcore.AddSync(os.Stdout), nil
	}

	dir := filepath.Dir(opts.OutputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil { //nolint:gosec
		return nil, err
	}

	f, err := os.OpenFile(opts.OutputPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644) //nolint:gosec
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(f), nil
}

func getEncoder(opts Options, cfg zapcore.EncoderConfig) zapcore.Encoder {
	if opts.Format == "json" {
		return zapcore.NewJSONEncoder(cfg)
	}

	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	if shouldUseColor(opts) {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return zapcore.NewConsoleEncoder(cfg)
}

// parseLevel 解析日志级别
func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, nil
	}
}

// shouldUseColor 判断是否应该使用颜色
func shouldUseColor(opts Options) bool {
	switch opts.Color {
	case "always":
		return true
	case "never", "no":
		return false
	case "auto", "":
		// auto: 输出到终端时启用颜色
		if opts.OutputPath != "" {
			return false // 文件输出不使用颜色
		}
		// 检查 stdout 是否是终端
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			return true
		}
		return false
	default:
		return false
	}
}
