// Package log provides logging functionality based on slog interface.
//
// This package provides a simple wrapper around the standard library log/slog
// with convenient configuration options. By default, it uses the standard library
// implementation with no external dependencies.
//
// For high-performance logging, use the zap adapter which wraps the official
// go.uber.org/zap/exp/zapslog implementation.
//
// # Basic Usage (Standard Library)
//
//	import "github.com/chinayin/gox/log"
//
//	logger, err := log.New(log.Options{
//		Level:  log.LevelInfo,
//		Format: log.FormatConsole,
//		Output: log.OutputStdout,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	logger.Info("application started", "port", 8080)
//	logger.Debug("debug message", "key", "value")
//	logger.Error("error occurred", "error", err)
//
// # Using Zap Adapter (High Performance)
//
//	import (
//		"github.com/chinayin/gox/log"
//		zaplog "github.com/chinayin/gox/log/zap"
//	)
//
//	// Method 1: Direct creation
//	logger, err := zaplog.New(log.Options{
//		Level:  log.LevelInfo,
//		Format: log.FormatJSON,
//		Output: log.OutputStdout,
//	})
//
//	// Method 2: Using NewWithHandler (similar to cli.NewStartupWithAdapter)
//	handler, err := zaplog.NewHandler(log.DefaultOptions())
//	logger := log.NewWithHandler(handler)
//
// # Configuration Options
//
// Level constants:
//   - log.LevelDebug - Debug level
//   - log.LevelInfo  - Info level (default)
//   - log.LevelWarn  - Warning level
//   - log.LevelError - Error level
//
// Format constants:
//   - log.FormatJSON    - JSON format (K8s standard)
//   - log.FormatConsole - Console format (human-readable)
//
// Output constants:
//   - log.OutputStdout - Standard output (K8s standard)
//   - log.OutputStderr - Standard error
//   - "/path/to/file"  - File path
//
// # Kubernetes Deployment
//
//	// Production: JSON to stdout for log collection
//	logger, _ := log.New(log.Options{
//		Level:  log.LevelInfo,
//		Format: log.FormatJSON,
//		Output: log.OutputStdout,
//	})
//
// # Local Development
//
//	// Development: Console format for readability
//	logger, _ := log.New(log.Options{
//		Level:  log.LevelDebug,
//		Format: log.FormatConsole,
//		Output: log.OutputStdout,
//	})
//
// # File Output
//
//	// Traditional deployment: Write to file
//	logger, _ := log.New(log.Options{
//		Level:  log.LevelInfo,
//		Format: log.FormatJSON,
//		Output: "/var/log/app.log",
//	})
//
// # Architecture
//
//	Default:  log.New() → slog.Handler (stdlib) → no dependencies
//	With Zap: zaplog.New() → zapslog.Handler → zap (high-perf)
package log
