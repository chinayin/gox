package config

import "errors"

// errors for config package.
var (
	// ErrReadFailed is returned when reading a config file fails.
	ErrReadFailed = errors.New("gox/config: failed to read config")

	// ErrUnmarshalFailed is returned when unmarshalling config data fails.
	ErrUnmarshalFailed = errors.New("gox/config: failed to unmarshal config")

	// ErrValidationFailed is returned when config validation fails.
	ErrValidationFailed = errors.New("gox/config: validation failed")

	// ErrMergeFailed is returned when merging configurations fails.
	ErrMergeFailed = errors.New("gox/config: failed to merge config")

	// ErrNotFound is returned when no config files are found or a specific file is missing.
	ErrNotFound = errors.New("gox/config: config file not found")
)
