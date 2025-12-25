package idgen

import "errors"

// Sentinel errors for idgen package.
var (
	// ErrNotInitialized is returned when the default generator is not set.
	ErrNotInitialized = errors.New("gox/idgen: default generator not initialized")

	// ErrAlreadyInitialized is returned when SetDefault is called more than once.
	ErrAlreadyInitialized = errors.New("gox/idgen: default generator already initialized")
)
