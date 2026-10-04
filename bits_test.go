package bitmask_test

import (
	"testing"

	"github.com/zefirixhd/bitmask"
)

func TestBitsOperations(t *testing.T) {
	var m bitmask.Mask[Permission]
	m.SetBit(0).SetBit(2)

	if !m.HasBit(0) || !m.HasBit(2) {
		t.Errorf("expected bits 0 and 2 to be set")
	}
	if m.HasBit(1) {
		t.Errorf("expected bit 1 to not be set")
	}
	if m.Count() != 2 {
		t.Errorf("expected Count to be 2, got %d", m.Count())
	}
	if m.TrailingZeros() != 0 {
		t.Errorf("expected TrailingZeros to be 0, got %d", m.TrailingZeros())
	}

	m.UnsetBit(0)
	if m.HasBit(0) {
		t.Errorf("expected bit 0 to be cleared")
	}
	if m.TrailingZeros() != 2 {
		t.Errorf("expected TrailingZeros to be 2, got %d", m.TrailingZeros())
	}

	m.ToggleBit(1)
	if !m.HasBit(1) {
		t.Errorf("expected bit 1 to be toggled on")
	}
}

func BenchmarkMask_Count(b *testing.B) {
	m := bitmask.New(PermissionLogin, PermissionRead, PermissionWrite, PermissionDelete)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Count()
	}
}
