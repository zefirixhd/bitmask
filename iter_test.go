package bitmask_test

import (
	"reflect"
	"testing"

	"github.com/zefirixhd/bitmask"
)

func TestIterators(t *testing.T) {
	perms := bitmask.New(PermissionLogin, PermissionDelete) // 1 and 8 (bits 0 and 3)

	var collectedFlags []Permission
	for f := range perms.All() {
		collectedFlags = append(collectedFlags, f)
	}
	expectedFlags := []Permission{PermissionLogin, PermissionDelete}
	if !reflect.DeepEqual(collectedFlags, expectedFlags) {
		t.Errorf("All iterator: expected %v, got %v", expectedFlags, collectedFlags)
	}

	// Early break in All
	var earlyFlags []Permission
	for f := range perms.All() {
		earlyFlags = append(earlyFlags, f)
		break
	}
	if len(earlyFlags) != 1 || earlyFlags[0] != PermissionLogin {
		t.Errorf("All iterator early break failed, got: %v", earlyFlags)
	}

	var collectedIndices []int
	for idx := range perms.BitIndices() {
		collectedIndices = append(collectedIndices, idx)
	}
	expectedIndices := []int{0, 3}
	if !reflect.DeepEqual(collectedIndices, expectedIndices) {
		t.Errorf("BitIndices iterator: expected %v, got %v", expectedIndices, collectedIndices)
	}

	// Flags() slice
	flagsSlice := perms.Flags()
	if !reflect.DeepEqual(flagsSlice, expectedFlags) {
		t.Errorf("Flags(): expected %v, got %v", expectedFlags, flagsSlice)
	}

	// Empty mask iteration
	var empty bitmask.Mask[Permission]
	for range empty.All() {
		t.Errorf("expected no iterations on empty mask")
	}
	if empty.Flags() != nil {
		t.Errorf("expected Flags() on empty mask to be nil")
	}
}

func BenchmarkMask_All(b *testing.B) {
	m := bitmask.New(PermissionLogin, PermissionRead, PermissionWrite, PermissionDelete)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for f := range m.All() {
			_ = f
		}
	}
}

func BenchmarkMask_BitIndices(b *testing.B) {
	m := bitmask.New(PermissionLogin, PermissionRead, PermissionWrite, PermissionDelete)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for idx := range m.BitIndices() {
			_ = idx
		}
	}
}
