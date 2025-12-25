// Package snowflake implements Twitter Snowflake ID generation algorithm.
//
// This package is based on github.com/bwmarrin/snowflake (BSD 2-Clause License)
// and has been vendored locally for long-term stability and customization.
//
// Original project: https://github.com/bwmarrin/snowflake
// Original author: Bruce (bwmarrin)
//
// # ID Structure (64 bit)
//
//	+--------------------------------------------------------------------------+
//	| 1 Bit Unused | 41 Bit Timestamp | 10 Bit NodeID | 12 Bit Sequence ID     |
//	+--------------------------------------------------------------------------+
//
// # Features
//
//   - Generates up to 4096 unique IDs per millisecond per node
//   - Supports up to 1024 nodes
//   - Timestamp epoch starts from 2010-11-04, usable for ~69 years
//   - Thread-safe ID generation
//
// # Basic Usage
//
//	node, err := snowflake.NewNode(1)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	id := node.Generate()
//	fmt.Println(id.Int64())
//	fmt.Println(id.String())
//	fmt.Println(id.Time())
//	fmt.Println(id.Node())
//	fmt.Println(id.Step())
package snowflake
