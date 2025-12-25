// Package idgen provides unified ID generation capabilities.
//
// This package offers a consistent interface for generating unique IDs
// using various algorithms. Currently supported:
//   - Snowflake: Distributed 64-bit ID, trend-increasing
//
// # Basic Usage (Recommended)
//
// Create a generator instance and inject it via dependency injection:
//
//	// Single-node deployment (default nodeID=1)
//	gen, err := idgen.NewSnowflake()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Distributed deployment (each node must have unique ID 0-1023)
//	gen, err := idgen.NewSnowflake(cfg.NodeID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	id := gen.Generate()
//	fmt.Println(id.Int64())   // int64 representation
//	fmt.Println(id.String())  // string representation
//
// # Accessing Raw Type
//
// Use Unwrap() to access the underlying implementation-specific type:
//
//	if sf, ok := id.Unwrap().(snowflake.ID); ok {
//	    fmt.Println(sf.Time())   // timestamp
//	    fmt.Println(sf.Node())   // node ID
//	    fmt.Println(sf.Step())   // sequence number
//	    fmt.Println(sf.Base64()) // base64 encoding
//	}
//
// # Global Generator (Use with Caution)
//
// A global generator is provided for simple use cases, but has limitations:
//
//   - Makes testing harder (global state)
//   - Hides dependencies
//   - Can cause issues in distributed systems if not properly configured
//
// If you choose to use it:
//
//	// At application startup (ONCE only)
//	gen, _ := idgen.NewSnowflake(cfg.NodeID)
//	idgen.SetDefault(gen)
//
//	// Anywhere in your code
//	id := idgen.Generate()
//
// Prefer dependency injection over global state when possible.
//
// # Thread Safety
//
// All generators are safe for concurrent use.
package idgen
