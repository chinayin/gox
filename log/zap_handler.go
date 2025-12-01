package log

import (
	"context"
	"log/slog"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// slogLevelToZap 日志级别映射表
var slogLevelToZap = map[slog.Level]zapcore.Level{
	slog.LevelDebug: zapcore.DebugLevel,
	slog.LevelInfo:  zapcore.InfoLevel,
	slog.LevelWarn:  zapcore.WarnLevel,
	slog.LevelError: zapcore.ErrorLevel,
}

// ZapHandler 实现 slog.Handler 接口，底层使用 zap
type ZapHandler struct {
	logger *zap.Logger
	level  zapcore.Level
	attrs  []zap.Field // 累积的属性
}

// NewZapHandler 创建新的 ZapHandler
func NewZapHandler(logger *zap.Logger, level zapcore.Level) *ZapHandler {
	return &ZapHandler{
		logger: logger,
		level:  level,
		attrs:  make([]zap.Field, 0),
	}
}

// Enabled 判断是否启用指定级别的日志
func (h *ZapHandler) Enabled(_ context.Context, level slog.Level) bool {
	zapLevel, ok := slogLevelToZap[level]
	if !ok {
		zapLevel = zapcore.InfoLevel
	}
	return zapLevel >= h.level
}

// Handle 处理日志记录（参考 samber/slog-zap 实现）
func (h *ZapHandler) Handle(_ context.Context, r slog.Record) error {
	zapLevel, ok := slogLevelToZap[r.Level]
	if !ok {
		zapLevel = zapcore.InfoLevel
	}

	fields := make([]zap.Field, 0, len(h.attrs)+r.NumAttrs())
	fields = append(fields, h.attrs...)

	r.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, h.slogAttrToZap(attr))
		return true
	})

	// 使用 Check 获取 CheckedEntry
	if ce := h.logger.Check(zapLevel, r.Message); ce != nil {
		// 关键：使用 slog.Record.PC 设置正确的 caller
		if r.PC != 0 {
			frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
			ce.Caller = zapcore.NewEntryCaller(r.PC, frame.File, frame.Line, true)
		}
		ce.Write(fields...)
	}

	return nil
}

// WithAttrs 返回一个新的 Handler，包含指定的属性
func (h *ZapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]zap.Field, len(h.attrs), len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	for _, attr := range attrs {
		newAttrs = append(newAttrs, h.slogAttrToZap(attr))
	}
	return &ZapHandler{
		logger: h.logger,
		level:  h.level,
		attrs:  newAttrs,
	}
}

// WithGroup 返回一个新的 Handler，日志记录在指定的组中
func (h *ZapHandler) WithGroup(name string) slog.Handler {
	// 简化实现：将组名作为命名空间
	return &ZapHandler{
		logger: h.logger.Named(name),
		level:  h.level,
		attrs:  h.attrs,
	}
}

// slogAttrToZap 将 slog.Attr 转换为 zap.Field
func (h *ZapHandler) slogAttrToZap(attr slog.Attr) zap.Field {
	k := attr.Key
	v := attr.Value
	switch v.Kind() {
	case slog.KindString:
		return zap.String(k, v.String())
	case slog.KindInt64:
		return zap.Int64(k, v.Int64())
	case slog.KindUint64:
		return zap.Uint64(k, v.Uint64())
	case slog.KindFloat64:
		return zap.Float64(k, v.Float64())
	case slog.KindBool:
		return zap.Bool(k, v.Bool())
	case slog.KindDuration:
		return zap.Duration(k, v.Duration())
	case slog.KindTime:
		return zap.Time(k, v.Time())
	case slog.KindAny:
		return zap.Any(k, v.Any())
	case slog.KindGroup:
		// 处理分组属性
		attrs := v.Group()
		groupMap := make(map[string]any, len(attrs))
		for _, a := range attrs {
			groupMap[a.Key] = a.Value.Any()
		}
		return zap.Any(k, groupMap)
	default:
		return zap.Any(k, v.Any())
	}
}
