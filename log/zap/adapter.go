// Package zap provides a zap-based implementation of slog.Handler.
package zap

import (
	"log/slog"

	"github.com/chinayin/gox/log"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

// NewHandler 创建基于 zap 的 slog.Handler
// 这是适配器，将 Options 转换为 zap Handler
func NewHandler(opts log.Options) (slog.Handler, error) {
	// 1. 使用官方配置
	var zapConfig zap.Config
	if opts.Format == log.FormatJSON {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	// 2. 设置级别
	zapConfig.Level = zap.NewAtomicLevelAt(parseLevel(opts.Level))

	// 3. 设置输出
	output := opts.Output
	if output == "" {
		output = log.OutputStdout
	}

	// 如果是文件路径，确保目录存在
	if err := log.EnsureOutputDir(output); err != nil {
		return nil, err
	}

	zapConfig.OutputPaths = []string{output}
	zapConfig.ErrorOutputPaths = []string{output}

	// 4. 创建 zap logger，根据 AddCaller 配置决定是否添加 caller
	zapOpts := []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
	if opts.AddCaller {
		zapOpts = append(zapOpts, zap.AddCaller())
	}

	zapLogger, err := zapConfig.Build(zapOpts...)
	if err != nil {
		return nil, err
	}

	// 5. 使用官方 zapslog 适配器，根据 AddCaller 配置启用 caller 信息
	return zapslog.NewHandler(zapLogger.Core(), zapslog.WithCaller(opts.AddCaller)), nil
}

// New 便捷函数：创建使用 zap 的 Logger
func New(opts log.Options) (*slog.Logger, error) {
	handler, err := NewHandler(opts)
	if err != nil {
		return nil, err
	}
	return slog.New(handler), nil
}

// zapLevelMap 日志级别映射表
var zapLevelMap = map[string]zapcore.Level{
	log.LevelDebug: zapcore.DebugLevel,
	log.LevelInfo:  zapcore.InfoLevel,
	log.LevelWarn:  zapcore.WarnLevel,
	log.LevelError: zapcore.ErrorLevel,
}

// parseLevel 解析日志级别
func parseLevel(level string) zapcore.Level {
	if lvl, ok := zapLevelMap[level]; ok {
		return lvl
	}
	return zapcore.InfoLevel
}
