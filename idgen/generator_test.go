package idgen

import (
	"errors"
	"sync"
	"testing"

	"github.com/chinayin/gox/idgen/snowflake"
)

func TestNewSnowflake(t *testing.T) {
	tests := []struct {
		name    string
		nodeID  int64
		wantErr bool
	}{
		{
			name:    "valid node 0",
			nodeID:  0,
			wantErr: false,
		},
		{
			name:    "valid node 1",
			nodeID:  1,
			wantErr: false,
		},
		{
			name:    "valid node max",
			nodeID:  1023,
			wantErr: false,
		},
		{
			name:    "invalid negative",
			nodeID:  -1,
			wantErr: true,
		},
		{
			name:    "invalid too large",
			nodeID:  1024,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewSnowflake(tt.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSnowflake() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gen == nil {
				t.Error("NewSnowflake() returned nil generator")
			}
		})
	}
}

func TestNewSnowflake_Default(t *testing.T) {
	// Test NewSnowflake with no arguments (uses default nodeID=1)
	gen, err := NewSnowflake()
	if err != nil {
		t.Fatalf("NewSnowflake() error = %v", err)
	}
	if gen == nil {
		t.Fatal("NewSnowflake() returned nil generator")
	}

	id := gen.Generate()
	sf, ok := id.Unwrap().(snowflake.ID)
	if !ok {
		t.Fatal("failed to unwrap snowflake.ID")
	}
	if sf.Node() != DefaultNodeID { //nolint:staticcheck // testing deprecated function
		t.Errorf("default nodeID = %v, want %v", sf.Node(), DefaultNodeID) //nolint:staticcheck // testing deprecated function
	}
}

func TestSnowflake_Generate(t *testing.T) {
	gen, err := NewSnowflake(1)
	if err != nil {
		t.Fatalf("NewSnowflake() error = %v", err)
	}

	id := gen.Generate()

	// Should have valid int64
	if id.Int64() <= 0 {
		t.Errorf("ID.Int64() = %v, want positive", id.Int64())
	}

	// Should have valid string
	if id.String() == "" {
		t.Error("ID.String() should not be empty")
	}

	// Should be able to unwrap to snowflake.ID
	sf, ok := id.Unwrap().(snowflake.ID)
	if !ok {
		t.Fatalf("ID.Unwrap() type assertion failed")
	}

	if sf.Node() != 1 { //nolint:staticcheck // testing deprecated function
		t.Errorf("snowflake.Node() = %v, want 1", sf.Node()) //nolint:staticcheck // testing deprecated function
	}
}

func TestSnowflake_Generate_Uniqueness(t *testing.T) {
	gen, err := NewSnowflake(1)
	if err != nil {
		t.Fatalf("NewSnowflake() error = %v", err)
	}

	const count = 1000
	ids := make(map[int64]bool, count)

	for i := 0; i < count; i++ {
		id := gen.Generate()
		if ids[id.Int64()] {
			t.Fatalf("duplicate ID generated: %v", id.Int64())
		}
		ids[id.Int64()] = true
	}
}

func TestSnowflake_Generate_Concurrent(t *testing.T) {
	gen, err := NewSnowflake(1)
	if err != nil {
		t.Fatalf("NewSnowflake() error = %v", err)
	}

	const goroutines = 10
	const idsPerGoroutine = 100

	var wg sync.WaitGroup
	idsChan := make(chan int64, goroutines*idsPerGoroutine)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id := gen.Generate()
				idsChan <- id.Int64()
			}
		}()
	}

	wg.Wait()
	close(idsChan)

	ids := make(map[int64]bool)
	for id := range idsChan {
		if ids[id] {
			t.Fatalf("duplicate ID in concurrent generation: %v", id)
		}
		ids[id] = true
	}

	if len(ids) != goroutines*idsPerGoroutine {
		t.Errorf("expected %d unique IDs, got %d", goroutines*idsPerGoroutine, len(ids))
	}
}

func TestSetDefault(t *testing.T) {
	// Reset before test
	ResetDefault()

	gen, err := NewSnowflake(1)
	if err != nil {
		t.Fatalf("NewSnowflake() error = %v", err)
	}

	// First call should succeed
	if err := SetDefault(gen); err != nil {
		t.Errorf("SetDefault() error = %v", err)
	}

	// Second call should fail
	gen2, _ := NewSnowflake(2)
	if err := SetDefault(gen2); !errors.Is(err, ErrAlreadyInitialized) {
		t.Errorf("SetDefault() error = %v, want ErrAlreadyInitialized", err)
	}

	// Cleanup
	ResetDefault()
}

func TestDefault(t *testing.T) {
	// Reset before test
	ResetDefault()

	// Should return nil before SetDefault
	if got := Default(); got != nil {
		t.Errorf("Default() = %v, want nil", got)
	}

	gen, _ := NewSnowflake(1)
	_ = SetDefault(gen)

	// Should return generator after SetDefault
	if got := Default(); got == nil {
		t.Error("Default() = nil after SetDefault")
	}

	// Cleanup
	ResetDefault()
}

func TestMustDefault_Panic(t *testing.T) {
	// Reset before test
	ResetDefault()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDefault() should panic when not initialized")
		}
	}()

	_ = MustDefault()
}

func TestGenerate_Global(t *testing.T) {
	// Reset before test
	ResetDefault()

	gen, _ := NewSnowflake(1)
	_ = SetDefault(gen)

	id := Generate()
	if id.IsZero() {
		t.Error("Generate() returned zero ID")
	}

	// Cleanup
	ResetDefault()
}

func TestGenerate_Panic(t *testing.T) {
	// Reset before test
	ResetDefault()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Generate() should panic when not initialized")
		}
	}()

	_ = Generate()
}
