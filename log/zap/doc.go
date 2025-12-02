// Package zap provides a zap-based implementation of slog.Handler.
//
// This package wraps the official go.uber.org/zap/exp/zapslog adapter
// and provides convenient configuration options.
//
// # Basic Usage
//
//	import zaplog "github.com/chinayin/gox/log/zap"
//
//	logger, err := zaplog.New(log.Options{
//		Level:  log.LevelInfo,
//		Format: log.FormatJSON,
//		Output: log.OutputStdout,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	logger.Info("application started", "port", 8080)
//
// # Using with log.NewWithHandler
//
//	import (
//		"github.com/chinayin/gox/log"
//		zaplog "github.com/chinayin/gox/log/zap"
//	)
//
//	handler, err := zaplog.NewHandler(log.DefaultOptions())
//	if err != nil {
//		panic(err)
//	}
//	logger := log.NewWithHandler(handler)
//
// # Kubernetes Deployment
//
//	// Production: JSON to stdout
//	logger, _ := zaplog.New(log.Options{
//		Level:  log.LevelInfo,
//		Format: log.FormatJSON,
//		Output: log.OutputStdout,
//	})
//
// # Performance
//
// This implementation uses zap for high-performance structured logging
// with zero-allocation in most cases.
package zap
