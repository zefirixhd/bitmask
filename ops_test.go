package bitmask_test

import (
	"testing"

	"github.com/zefirixhd/bitmask"
)

func TestBasicOps(t *testing.T) {
	m := bitmask.New(PermissionRead, PermissionLogin)

	if !m.Has(PermissionLogin) {
		t.Errorf("expected Has(PermissionLogin) to be true")
	}
	if !m.Has(PermissionRead) {
		t.Errorf("expected Has(PermissionRead) to be true")
	}
	if m.Has(PermissionWrite) {
		t.Errorf("expected Has(PermissionWrite) to be false")
	}
	if m.Has(0) {
		t.Errorf("expected Has(0) to be false")
	}

	m.Set(PermissionWrite)
	if !m.Has(PermissionWrite) {
		t.Errorf("expected Has(PermissionWrite) to be true after Set")
	}

	m.Unset(PermissionRead)
	if m.Has(PermissionRead) {
		t.Errorf("expected Has(PermissionRead) to be false after Unset")
	}

	// Toggle
	m.Toggle(PermissionWrite)
	if m.Has(PermissionWrite) {
		t.Errorf("expected Has(PermissionWrite) to be false after Toggle")
	}
	m.Toggle(PermissionWrite)
	if !m.Has(PermissionWrite) {
		t.Errorf("expected Has(PermissionWrite) to be true after second Toggle")
	}

	// HasAny & HasAll
	if !m.HasAny(PermissionLogin, PermissionDelete) {
		t.Errorf("expected HasAny to be true")
	}
	if m.HasAny(PermissionDelete) {
		t.Errorf("expected HasAny to be false for absent flags")
	}
	if !m.HasAll(PermissionLogin, PermissionWrite) {
		t.Errorf("expected HasAll to be true")
	}
	if m.HasAll(PermissionLogin, PermissionDelete) {
		t.Errorf("expected HasAll to be false")
	}

	// Invert
	var small bitmask.Mask[uint8]
	small.Set(0b00000001)
	small.Invert()
	if small.Bits() != 0b11111110 {
		t.Errorf("expected Invert to flip all bits, got %b", small.Bits())
	}
}

func TestImmutableWithWithout(t *testing.T) {
	base := bitmask.New(PermissionRead)
	derived := base.With(PermissionWrite)

	if base.Has(PermissionWrite) {
		t.Errorf("base mask should not be mutated by With")
	}
	if !derived.Has(PermissionWrite) {
		t.Errorf("derived mask should have PermissionWrite")
	}

	stripped := derived.Without(PermissionRead)
	if !derived.Has(PermissionRead) {
		t.Errorf("derived mask should not be mutated by Without")
	}
	if stripped.Has(PermissionRead) {
		t.Errorf("stripped mask should not have PermissionRead")
	}
}

func BenchmarkMask_Has(b *testing.B) {
	m := bitmask.New(PermissionLogin, PermissionWrite)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Has(PermissionWrite)
	}
}

func BenchmarkMask_Set(b *testing.B) {
	var m bitmask.Mask[Permission]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Set(PermissionLogin)
	}
}

func BenchmarkMask_Toggle(b *testing.B) {
	var m bitmask.Mask[Permission]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Toggle(PermissionLogin)
	}
}
