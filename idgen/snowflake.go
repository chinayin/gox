package idgen

import (
	"github.com/chinayin/gox/idgen/snowflake"
)

// DefaultNodeID is the default node ID used when NewSnowflake is called without arguments.
const DefaultNodeID int64 = 1

// Snowflake wraps the snowflake.Node to implement the Generator interface.
type Snowflake struct {
	node *snowflake.Node
}

// NewSnowflake creates a new Snowflake ID generator.
//
// If nodeID is not provided, defaults to 1 (suitable for single-node deployment).
// For distributed deployments, each node MUST have a unique nodeID (0-1023).
//
// Example:
//
//	// Single-node deployment
//	gen, _ := idgen.NewSnowflake()
//
//	// Distributed deployment (each node must have unique ID)
//	gen, _ := idgen.NewSnowflake(cfg.NodeID)
func NewSnowflake(nodeID ...int64) (*Snowflake, error) {
	id := DefaultNodeID
	if len(nodeID) > 0 {
		id = nodeID[0]
	}
	node, err := snowflake.NewNode(id)
	if err != nil {
		return nil, err
	}
	return &Snowflake{node: node}, nil
}

// Generate creates a new unique ID.
func (s *Snowflake) Generate() ID {
	sf := s.node.Generate()
	return NewID(sf.Int64(), sf.String(), sf)
}
