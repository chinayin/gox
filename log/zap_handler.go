package log

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
	return h.slogLevelToZap(level) >= h.level
}

// Handle 处理日志记录
func (h *ZapHandler) Handle(_ context.Context, r slog.Record) error {
	// 转换日志级别
	lvl := h.slogLevelToZap(r.Level)

	// 收集所有字段
	fields := make([]zap.Field, 0, len(h.attrs)+r.NumAttrs())

	// 添加累积的属性
	fields = append(fields, h.attrs...)

	// 添加当前记录的属性
	r.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, h.slogAttrToZap(attr))
		return true
	})

	// 记录日志
	if ce := h.logger.Check(lvl, r.Message); ce != nil {
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

// slogLevelToZap 将 slog.Level 转换为 zapcore.Level
func (h *ZapHandler) slogLevelToZap(level slog.Level) zapcore.Level {
	switch {
	case level >= slog.LevelError:
		return zapcore.ErrorLevel
	case level >= slog.LevelWarn:
		return zapcore.WarnLevel
	case level >= slog.LevelInfo:
		return zapcore.InfoLevel
	default:
		return zapcore.DebugLevel
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
		// 将 group 转换为 map，因为 zap.Namespace 不能用于此处的值上下文
		// 且 zap.Any 需要具体的值
		groupMap := make(map[string]any, len(attrs))
		for _, a := range attrs {
			groupMap[a.Key] = a.Value.Any()
		}
		return zap.Any(k, groupMap)
	default:
		return zap.Any(k, v.Any())
	}
}
