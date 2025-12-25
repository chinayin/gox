package idgen

// ID represents a unified ID type that wraps various ID implementations.
// It provides common accessors for int64 and string representations,
// while allowing access to the underlying implementation via Unwrap().
type ID struct {
	intVal int64
	strVal string
	raw    any
}

// Int64 returns the int64 representation of the ID.
// For some implementations (e.g., UUID), this may return 0.
func (id ID) Int64() int64 {
	return id.intVal
}

// String returns the string representation of the ID.
func (id ID) String() string {
	return id.strVal
}

// IsZero reports whether the ID is a zero value.
func (id ID) IsZero() bool {
	return id.intVal == 0 && id.strVal == ""
}

// Unwrap returns the underlying implementation-specific value.
// Use type assertion to access implementation-specific methods.
//
// Example (Snowflake):
//
//	if sf, ok := id.Unwrap().(snowflake.ID); ok {
//	    fmt.Println(sf.Time())   // timestamp in milliseconds
//	    fmt.Println(sf.Node())   // node ID
//	    fmt.Println(sf.Step())   // sequence number
//	    fmt.Println(sf.Base64()) // base64 encoding
//	}
func (id ID) Unwrap() any {
	return id.raw
}

// NewID creates a new ID with the given values.
// This is typically used by generator implementations.
func NewID(intVal int64, strVal string, raw any) ID {
	return ID{
		intVal: intVal,
		strVal: strVal,
		raw:    raw,
	}
}
