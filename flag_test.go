package bitmask_test

import (
	"reflect"
	"testing"

	"github.com/zefirixhd/bitmask"
)

func TestDecompose(t *testing.T) {
	// 0 returns nil
	if bitmask.Decompose(uint8(0)) != nil {
		t.Errorf("expected Decompose(0) to return nil")
	}

	// Single atomic flag
	single := bitmask.Decompose(uint8(4))
	expectedSingle := []uint8{4}
	if !reflect.DeepEqual(single, expectedSingle) {
		t.Errorf("expected %v, got %v", expectedSingle, single)
	}

	// Composite flag
	composite := bitmask.Decompose(uint8(1 | 2 | 8))
	expectedComposite := []uint8{1, 2, 8}
	if !reflect.DeepEqual(composite, expectedComposite) {
		t.Errorf("expected %v, got %v", expectedComposite, composite)
	}
}

func BenchmarkMask_Decompose(b *testing.B) {
	val := PermissionLogin | PermissionRead | PermissionWrite | PermissionDelete
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bitmask.Decompose(val)
	}
}
