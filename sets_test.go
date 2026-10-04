package bitmask_test

import (
	"testing"

	"github.com/zefirixhd/bitmask"
)

func TestSetOperations(t *testing.T) {
	a := bitmask.New(PermissionRead, PermissionLogin)
	b := bitmask.New(PermissionRead, PermissionWrite)

	and := a.And(b)
	if and.Bits() != PermissionRead {
		t.Errorf("expected And to have only PermissionRead, got %d", and.Bits())
	}

	or := a.Or(b)
	if or.Bits() != (PermissionRead | PermissionLogin | PermissionWrite) {
		t.Errorf("expected Or to combine flags, got %d", or.Bits())
	}

	xor := a.Xor(b)
	if xor.Bits() != (PermissionLogin | PermissionWrite) {
		t.Errorf("expected Xor to exclude common flags, got %d", xor.Bits())
	}

	andNot := a.AndNot(b)
	if andNot.Bits() != PermissionLogin {
		t.Errorf("expected AndNot to leave only PermissionLogin, got %d", andNot.Bits())
	}

	sub := bitmask.New(PermissionRead)
	if !a.Contains(sub) {
		t.Errorf("expected a to contain sub")
	}
	if sub.Contains(a) {
		t.Errorf("expected sub not to contain a")
	}
	if a.Contains(bitmask.Mask[Permission]{}) {
		t.Errorf("expected Contains(empty) to be false")
	}

	if !a.Overlaps(b) {
		t.Errorf("expected a to overlap with b (PermissionRead is common)")
	}

	c := bitmask.New(PermissionDelete)
	if a.Overlaps(c) {
		t.Errorf("expected a not to overlap with c")
	}
}

func BenchmarkMask_Contains(b *testing.B) {
	a := bitmask.New(PermissionLogin, PermissionRead, PermissionWrite)
	sub := bitmask.New(PermissionLogin, PermissionWrite)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.Contains(sub)
	}
}

func BenchmarkMask_Overlaps(b *testing.B) {
	a := bitmask.New(PermissionLogin, PermissionRead)
	sub := bitmask.New(PermissionRead, PermissionWrite)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.Overlaps(sub)
	}
}
