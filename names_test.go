package bitmask_test

import (
	"reflect"
	"testing"

	"github.com/zefirixhd/bitmask"
)

func TestNamesAndParse(t *testing.T) {
	perms := bitmask.New(PermissionLogin, PermissionWrite)

	// Explicit map
	names := perms.Names(permNames)
	expectedNames := []string{"login", "write"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf("expected names %v, got %v", expectedNames, names)
	}

	formatted := perms.Format(permNames, ", ")
	if formatted != "login, write" {
		t.Errorf("expected format 'login, write', got %q", formatted)
	}

	// Fallback for unmapped flag
	permsWithUnknown := perms.With(Permission(64))
	namesWithFallback := permsWithUnknown.Names(permNames)
	expectedFallback := []string{"login", "write", "64"}
	if !reflect.DeepEqual(namesWithFallback, expectedFallback) {
		t.Errorf("expected fallback %v, got %v", expectedFallback, namesWithFallback)
	}

	// Parse with map
	parsed, err := bitmask.Parse([]string{"login", "write"}, permNames)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if parsed.Bits() != perms.Bits() {
		t.Errorf("expected parsed bits %d, got %d", perms.Bits(), parsed.Bits())
	}

	_, err = bitmask.Parse([]string{"unknown_permission"}, permNames)
	if err == nil {
		t.Errorf("expected error for unknown flag in Parse")
	}
}

func TestAutoNamesAndParse(t *testing.T) {
	perms := bitmask.New(PermissionLogin, PermissionWrite)

	// Automatic Strings() without map (using Permission.String())
	autoStrings := perms.Strings()
	expected := []string{"login", "write"}
	if !reflect.DeepEqual(autoStrings, expected) {
		t.Errorf("expected auto strings %v, got %v", expected, autoStrings)
	}

	// Automatic FormatAuto()
	autoFormatted := perms.FormatAuto(" | ")
	if autoFormatted != "login | write" {
		t.Errorf("expected 'login | write', got %q", autoFormatted)
	}

	// Automatic ParseAuto() without map
	autoParsed, err := bitmask.ParseAuto[Permission]("login", "write")
	if err != nil {
		t.Fatalf("ParseAuto failed: %v", err)
	}
	if autoParsed.Bits() != perms.Bits() {
		t.Errorf("expected autoParsed bits %d, got %d", perms.Bits(), autoParsed.Bits())
	}

	// ParseAuto with unknown flag
	_, err = bitmask.ParseAuto[Permission]("nonexistent")
	if err == nil {
		t.Errorf("expected error for unknown flag in ParseAuto")
	}

	// Raw uint8 mask without fmt.Stringer fallback in Strings()
	var raw bitmask.Mask[uint8]
	raw.Set(1, 4)
	rawStrings := raw.Strings()
	expectedRaw := []string{"1", "4"}
	if !reflect.DeepEqual(rawStrings, expectedRaw) {
		t.Errorf("expected raw strings %v, got %v", expectedRaw, rawStrings)
	}
}

func BenchmarkMask_Parse(b *testing.B) {
	strs := []string{"login", "write"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bitmask.Parse(strs, permNames)
	}
}

func BenchmarkMask_ParseAuto(b *testing.B) {
	strs := []string{"login", "write"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bitmask.ParseAuto[Permission](strs...)
	}
}

func BenchmarkMask_Names(b *testing.B) {
	m := bitmask.New(PermissionLogin, PermissionWrite)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Names(permNames)
	}
}

func BenchmarkMask_Strings(b *testing.B) {
	m := bitmask.New(PermissionLogin, PermissionWrite)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Strings()
	}
}
