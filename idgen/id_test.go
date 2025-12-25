package idgen

import (
	"testing"

	"github.com/chinayin/gox/idgen/snowflake"
)

func TestID_Int64(t *testing.T) {
	id := NewID(1234567890, "1234567890", nil)

	if got := id.Int64(); got != 1234567890 {
		t.Errorf("ID.Int64() = %v, want %v", got, 1234567890)
	}
}

func TestID_String(t *testing.T) {
	id := NewID(1234567890, "1234567890", nil)

	if got := id.String(); got != "1234567890" {
		t.Errorf("ID.String() = %v, want %v", got, "1234567890")
	}
}

func TestID_IsZero(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want bool
	}{
		{
			name: "zero value",
			id:   ID{},
			want: true,
		},
		{
			name: "non-zero int",
			id:   NewID(1, "", nil),
			want: false,
		},
		{
			name: "non-zero string",
			id:   NewID(0, "abc", nil),
			want: false,
		},
		{
			name: "both non-zero",
			id:   NewID(1, "1", nil),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsZero(); got != tt.want {
				t.Errorf("ID.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID_Unwrap(t *testing.T) {
	// Test with snowflake.ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		t.Fatalf("failed to create snowflake node: %v", err)
	}

	sf := node.Generate()
	id := NewID(sf.Int64(), sf.String(), sf)

	unwrapped := id.Unwrap()
	if unwrapped == nil {
		t.Fatal("ID.Unwrap() returned nil")
	}

	sfID, ok := unwrapped.(snowflake.ID)
	if !ok {
		t.Fatalf("ID.Unwrap() type assertion failed, got %T", unwrapped)
	}

	if sfID.Int64() != sf.Int64() {
		t.Errorf("unwrapped ID = %v, want %v", sfID.Int64(), sf.Int64())
	}

	// Verify we can access snowflake-specific methods
	if sfID.Time() <= 0 { //nolint:staticcheck // testing deprecated function
		t.Error("snowflake Time() should return positive value")
	}
	if sfID.Node() != 1 { //nolint:staticcheck // testing deprecated function
		t.Errorf("snowflake Node() = %v, want 1", sfID.Node()) //nolint:staticcheck // testing deprecated function
	}
}

func TestID_Unwrap_Nil(t *testing.T) {
	id := NewID(0, "", nil)

	if got := id.Unwrap(); got != nil {
		t.Errorf("ID.Unwrap() = %v, want nil", got)
	}
}
