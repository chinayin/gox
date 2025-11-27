// Package validator provides data validation utilities based on go-playground/validator.
//
// This package is a thin wrapper around github.com/go-playground/validator/v10,
// providing a unified validation interface, custom validation rules, and multi-language support.
//
// # Official Validation Rules
//
// The underlying validator library provides extensive built-in validation rules.
// For a complete list of official validation tags, see:
// https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Baked_In_Validators_and_Tags
//
// Common official validation rules (DO NOT reimplement these):
//   - required: Field must be present and non-zero
//   - email: Valid email address
//   - url: Valid URL
//   - uri: Valid URI
//   - min/max: Minimum/maximum value or length
//   - len: Exact length
//   - eq/ne/gt/gte/lt/lte: Comparison operators
//   - oneof: Value must be one of the specified options
//   - alpha/alphanum/numeric: Character type validation
//   - datetime: Date/time format validation (e.g., datetime=2006-01-02)
//   - ip/ipv4/ipv6: IP address validation
//   - uuid: UUID validation
//   - json: Valid JSON string
//   - And many more...
//
// # Custom Validation Rules
//
// This package adds the following custom validation rules specific to gox:
//
//   - snowflake_id: Validates Snowflake ID format
//     Supports int64 and string types
//     Must be a positive integer (> 0)
//     Example usage as a template for adding custom rules
//
// # Adding Custom Rules
//
// The package uses a minimalist design for easy extension.
// To add a new rule, see builtin/README.md for the simple 2-step process.
//
// Architecture:
//   - builtin/yourrule.go: Validation function + translation map
//   - builtin/builtin.go: Register the rule
//   - builtin/translations.go: Automatically handles all languages (one for loop)
//
// # Basic Usage
//
// Use the global validator for simple validation:
//
//	import "github.com/chinayin/gox/validator"
//
//	type User struct {
//	    Email       string `validate:"required,email"`
//	    SnowflakeID int64  `validate:"required,snowflake_id"`
//	}
//
//	user := User{
//	    Email:       "user@example.com",
//	    SnowflakeID: 1234567890123456789,
//	}
//
//	if err := validator.Validate(&user); err != nil {
//	    // Handle validation error
//	}
//
// # Multi-language Support
//
// The validator supports multiple languages for error messages:
//
//	// Create validator with Chinese locale
//	v := validator.New(validator.WithLocale("zh"))
//
//	// Or switch locale at runtime
//	v.SetLocale("zh")
//
//	// Validation errors will be automatically translated
//	err := v.Validate(&user)
//	if err != nil {
//	    fmt.Println(err.Error())  // Output in Chinese
//	}
//
// Supported locales:
//   - "en" - English (default)
//   - "zh" - Simplified Chinese (简体中文)
//
// # Custom Validator Instance
//
// Create a custom validator instance for advanced use cases:
//
//	v := validator.New()
//
//	// Register custom validation rule
//	v.RegisterValidation("custom_rule", func(fl validator.FieldLevel) bool {
//	    // Custom validation logic
//	    return true
//	})
//
//	if err := v.Validate(&user); err != nil {
//	    // Handle validation error
//	}
//
// # Error Handling
//
// Validation errors are automatically translated to the current locale:
//
//	if err := validator.Validate(&user); err != nil {
//	    // Simple error message
//	    fmt.Println(err.Error())
//
//	    // Or get all error messages
//	    if translatedErr, ok := err.(*validator.TranslatedError); ok {
//	        for _, msg := range translatedErr.Errors() {
//	            fmt.Println(msg)
//	        }
//	    }
//	}
//
// # Thread Safety
//
// The global validator instance is lazily initialized and thread-safe.
// Custom validator instances are also safe for concurrent use after initialization.
//
// # Performance
//
// The validator uses reflection and caching internally for optimal performance.
// Validation overhead is typically negligible compared to I/O operations.
package validator
