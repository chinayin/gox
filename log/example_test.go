package log_test

import (
	"errors"
	"log/slog"
	"time"

	"github.com/chinayin/gox/log"
)

func ExampleNew() {
	// 创建日志实例
	logger, err := log.New(log.Options{
		Level:      "info",
		OutputPath: "",        // 空字符串表示输出到 stdout
		Format:     "console", // 或 "json"
	})
	if err != nil {
		panic(err)
	}

	// 基本日志记录
	logger.Info("application started")
	logger.Debug("this will not be printed because level is info")
	logger.Warn("warning message")
	logger.Error("error occurred")
}

func ExampleNew_structuredLogging() {
	logger, _ := log.New(log.DefaultOptions())

	// 结构化日志
	logger.Info("user login",
		"user_id", 123,
		"username", "john",
		"ip", "192.168.1.1",
		"timestamp", time.Now())

	// 记录错误
	err := errors.New("database connection failed")
	logger.Error("operation failed",
		"error", err,
		"retry_count", 3,
		"timeout", 30*time.Second)
}

func ExampleNew_withContext() {
	logger, _ := log.New(log.DefaultOptions())

	// 使用 With 添加上下文字段
	requestLogger := logger.With(
		"request_id", "abc-123",
		"method", "GET",
		"path", "/api/users")

	// 所有日志都会包含上下文字段
	requestLogger.Info("processing request")
	requestLogger.Debug("validating input")
	requestLogger.Info("request completed", "duration", 150*time.Millisecond)
}

func ExampleNew_jsonFormat() {
	// JSON 格式适合生产环境
	logger, _ := log.New(log.Options{
		Level:      "info",
		OutputPath: "/var/log/app.log",
		Format:     "json",
	})

	logger.Info("server started",
		"port", 8080,
		"env", "production")
}

func ExampleNew_groupedAttributes() {
	logger, _ := log.New(log.DefaultOptions())

	// 使用 slog.Group 组织相关字段
	logger.Info("请求完成",
		slog.Group("request",
			"method", "POST",
			"path", "/api/users",
			"duration_ms", 150,
		),
		slog.Group("response",
			"status", 200,
			"size_bytes", 1024,
		),
	)

	// 嵌套分组
	logger.Info("数据库操作",
		slog.Group("db",
			"host", "localhost",
			"port", 5432,
			slog.Group("stats",
				"connections", 10,
				"queries", 1500,
			),
		),
	)
}

func ExampleNew_colorControl() {
	// 自动检测（默认）
	logger1, _ := log.New(log.Options{
		Level:  "info",
		Format: "console",
		Color:  "auto", // 终端输出有颜色，管道/文件无颜色
	})
	logger1.Info("auto color mode")

	// 强制启用颜色
	logger2, _ := log.New(log.Options{
		Level:  "info",
		Format: "console",
		Color:  "always",
	})
	logger2.Info("always color mode")

	// 禁用颜色（适合 k8s 容器日志收集）
	logger3, _ := log.New(log.Options{
		Level:  "info",
		Format: "console",
		Color:  "never", // 或 "no"
	})
	logger3.Info("no color mode")
}
