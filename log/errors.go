package log

import "errors"

// errors for log package.
var (
	// ErrOpenFile is returned when opening a log file fails.
	ErrOpenFile = errors.New("gox/log: failed to open log file")
)
