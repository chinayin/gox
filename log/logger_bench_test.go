package log

import (
	"testing"
)

// BenchmarkLogger_Info 测试 Info 日志性能
func BenchmarkLogger_Info(b *testing.B) {
	logger, err := New(Options{
		Level:      "info",
		OutputPath: "/dev/null", // 输出到 /dev/null 避免 I/O 影响
		Format:     "json",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test message",
			"key1", "value1",
			"key2", 42,
			"key3", true)
	}
}

// BenchmarkLogger_Debug 测试 Debug 日志性能
func BenchmarkLogger_Debug(b *testing.B) {
	logger, err := New(Options{
		Level:      "debug",
		OutputPath: "/dev/null",
		Format:     "json",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("benchmark test message",
			"key1", "value1",
			"key2", 42)
	}
}

// BenchmarkLogger_Error 测试 Error 日志性能
func BenchmarkLogger_Error(b *testing.B) {
	logger, err := New(Options{
		Level:      "error",
		OutputPath: "/dev/null",
		Format:     "json",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("benchmark test error",
			"error", "test error",
			"code", 500)
	}
}

// BenchmarkLogger_WithFields 测试带多个字段的日志性能
func BenchmarkLogger_WithFields(b *testing.B) {
	logger, err := New(Options{
		Level:      "info",
		OutputPath: "/dev/null",
		Format:     "json",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark with many fields",
			"field1", "value1",
			"field2", "value2",
			"field3", 123,
			"field4", 456,
			"field5", true,
			"field6", 3.14,
			"field7", "value7",
			"field8", "value8")
	}
}

// BenchmarkLogger_ConsoleFormat 测试 Console 格式性能
func BenchmarkLogger_ConsoleFormat(b *testing.B) {
	logger, err := New(Options{
		Level:      "info",
		OutputPath: "/dev/null",
		Format:     "console",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark console format",
			"key1", "value1",
			"key2", 42)
	}
}

// BenchmarkLogger_JSONFormat 测试 JSON 格式性能
func BenchmarkLogger_JSONFormat(b *testing.B) {
	logger, err := New(Options{
		Level:      "info",
		OutputPath: "/dev/null",
		Format:     "json",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark json format",
			"key1", "value1",
			"key2", 42)
	}
}

// BenchmarkLogger_With 测试 With 方法性能
func BenchmarkLogger_With(b *testing.B) {
	logger, err := New(Options{
		Level:      "info",
		OutputPath: "/dev/null",
		Format:     "json",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		contextLogger := logger.With("request_id", "abc123", "user_id", 456)
		contextLogger.Info("processing request")
	}
}
