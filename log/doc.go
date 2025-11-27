// Package log provides logging functionality based on slog interface, implemented with zap.
//
// Design Philosophy:
//   - Uses standard library log/slog interface for API standardization and portability
//   - Implemented with zap for high performance and rich features
//   - Connects slog and zap through Handler adapter pattern
//
// Basic Usage:
//
//	logger, err := log.New(log.Options{
//	    Level:      "info",
//	    OutputPath: "",        // empty string means stdout
//	    Format:     "console", // or "json"
//	    Color:      "auto",    // auto(default), always, never, no
//	})
//	if err != nil {
//	    panic(err)
//	}
//
//	// Use standard slog API
//	logger.Info("application started")
//	logger.Debug("debug message", "key", "value")
//	logger.Error("error occurred", "error", err)
//
// Structured Logging:
//
//	logger.Info("user login",
//	    "user_id", 123,
//	    "username", "john",
//	    "ip", "192.168.1.1")
//
//	// Using groups
//	logger.Info("request completed",
//	    slog.Group("request",
//	        "method", "POST",
//	        "path", "/api/users",
//	    ),
//	)
//
// Adding Context Fields:
//
//	// Use With to add persistent fields
//	requestLogger := logger.With(
//	    "request_id", "abc123",
//	    "method", "GET")
//	requestLogger.Info("processing request")
//
// Performance Characteristics:
//   - Zero-allocation structured logging (in most cases)
//   - Supports log level filtering
//   - Supports colored console output and JSON format
//   - Automatically adds caller information
//   - Automatically adds stack trace for Error level
package log
