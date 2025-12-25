package validator

import "errors"

// errors for validator package.
var (
	// ErrLocaleNotFound is returned when the requested locale is not supported.
	ErrLocaleNotFound = errors.New("gox/validator: locale not found")
)
