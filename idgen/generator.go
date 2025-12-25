package idgen

import "sync/atomic"

// Generator is the interface for ID generators.
type Generator interface {
	Generate() ID
}

// ========== Global Generator ==========

var defaultGen atomic.Pointer[Generator]

// SetDefault sets the default global generator.
// This should be called once at application startup.
// Returns ErrAlreadyInitialized if called more than once.
func SetDefault(g Generator) error {
	if !defaultGen.CompareAndSwap(nil, &g) {
		return ErrAlreadyInitialized
	}
	return nil
}

// Default returns the default global generator.
// Returns nil if SetDefault has not been called.
func Default() Generator {
	if g := defaultGen.Load(); g != nil {
		return *g
	}
	return nil
}

// MustDefault returns the default global generator.
// Panics if SetDefault has not been called.
func MustDefault() Generator {
	g := Default()
	if g == nil {
		panic(ErrNotInitialized)
	}
	return g
}

// Generate creates a new ID using the default global generator.
// Panics if SetDefault has not been called.
func Generate() ID {
	return MustDefault().Generate()
}

// ResetDefault resets the default generator to nil.
// This is primarily for testing purposes.
func ResetDefault() {
	defaultGen.Store(nil)
}
