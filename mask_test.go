package bitmask_test

import (
	"strconv"
	"testing"

	"github.com/zefirixhd/bitmask"
)

type Permission uint8

const (
	PermissionLogin Permission = 1 << iota
	PermissionRead
	PermissionWrite
	PermissionDelete
)

func (p Permission) String() string {
	switch p {
	case PermissionLogin:
		return "login"
	case PermissionRead:
		return "read"
	case PermissionWrite:
		return "write"
	case PermissionDelete:
		return "delete"
	default:
		return strconv.Itoa(int(p))
	}
}

var permNames = map[Permission]string{
	PermissionLogin:  "login",
	PermissionRead:   "read",
	PermissionWrite:  "write",
	PermissionDelete: "delete",
}

func TestCoreLifecycle(t *testing.T) {
	m := bitmask.New(PermissionRead, PermissionLogin)

	if m.IsEmpty() {
		t.Errorf("expected mask not to be empty")
	}
	if m.Bits() != (PermissionRead | PermissionLogin) {
		t.Errorf("expected bits %d, got %d", PermissionRead|PermissionLogin, m.Bits())
	}

	wrapped := bitmask.Wrap(PermissionWrite)
	if wrapped.Bits() != PermissionWrite {
		t.Errorf("expected wrapped bits %d, got %d", PermissionWrite, wrapped.Bits())
	}

	m.Reset()
	if !m.IsEmpty() {
		t.Errorf("expected mask to be empty after Reset")
	}
	if m.Bits() != 0 {
		t.Errorf("expected bits to be 0 after Reset, got %d", m.Bits())
	}
}
