package strict_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/zefirixhd/bitmask"
	"github.com/zefirixhd/bitmask/strict"
)

type Permission uint32

const (
	PermRead Permission = 1 << iota
	PermWrite
	PermDelete
	PermAdmin
	PermGuest
)

func TestMutuallyExclusive(t *testing.T) {
	rule := strict.MutuallyExclusive(PermRead, PermWrite)

	// Empty rule call
	emptyRule := strict.MutuallyExclusive[Permission]()
	if err := emptyRule(bitmask.New(PermRead)); err != nil {
		t.Fatalf("expected nil for empty rule, got %v", err)
	}

	// 0 flags set
	if err := rule(bitmask.New[Permission]()); err != nil {
		t.Fatalf("expected nil for 0 flags, got %v", err)
	}

	// 1 flag set (Read)
	if err := rule(bitmask.New(PermRead)); err != nil {
		t.Fatalf("expected nil for 1 flag, got %v", err)
	}

	// 1 flag set (Write)
	if err := rule(bitmask.New(PermWrite)); err != nil {
		t.Fatalf("expected nil for 1 flag, got %v", err)
	}

	// Both flags set -> conflict
	err := rule(bitmask.New(PermRead, PermWrite))
	if !errors.Is(err, strict.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestRequireAll(t *testing.T) {
	rule := strict.RequireAll(PermRead, PermWrite)

	// Empty rule call
	if err := strict.RequireAll[Permission]()(bitmask.New[Permission]()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Missing both
	if err := rule(bitmask.New[Permission]()); !errors.Is(err, strict.ErrMissingRequired) {
		t.Fatalf("expected ErrMissingRequired, got %v", err)
	}

	// Missing Write
	if err := rule(bitmask.New(PermRead)); !errors.Is(err, strict.ErrMissingRequired) {
		t.Fatalf("expected ErrMissingRequired, got %v", err)
	}

	// Has both
	if err := rule(bitmask.New(PermRead, PermWrite)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRequireAny(t *testing.T) {
	rule := strict.RequireAny(PermRead, PermWrite)

	// Empty rule call
	if err := strict.RequireAny[Permission]()(bitmask.New[Permission]()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Has neither
	if err := rule(bitmask.New[Permission]()); !errors.Is(err, strict.ErrMissingRequired) {
		t.Fatalf("expected ErrMissingRequired, got %v", err)
	}

	// Has Read
	if err := rule(bitmask.New(PermRead)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Has Write
	if err := rule(bitmask.New(PermWrite)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestNoneOf(t *testing.T) {
	rule := strict.NoneOf(PermGuest, PermDelete)

	// Empty rule call
	if err := strict.NoneOf[Permission]()(bitmask.New(PermGuest)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Has neither forbidden flag
	if err := rule(bitmask.New(PermRead, PermWrite)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Has one forbidden flag (Guest)
	if err := rule(bitmask.New(PermRead, PermGuest)); !errors.Is(err, strict.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	// Has other forbidden flag (Delete)
	if err := rule(bitmask.New(PermDelete)); !errors.Is(err, strict.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDependsOn(t *testing.T) {
	// PermDelete depends on PermRead and PermWrite
	rule := strict.DependsOn(PermDelete, PermRead, PermWrite)

	// Flag not set -> always succeeds
	if err := rule(bitmask.New(PermAdmin)); err != nil {
		t.Fatalf("expected nil when flag not set, got %v", err)
	}

	// Delete set, but missing Read and Write
	if err := rule(bitmask.New(PermDelete)); !errors.Is(err, strict.ErrDependencyMissing) {
		t.Fatalf("expected ErrDependencyMissing, got %v", err)
	}

	// Delete set, has Read but missing Write
	if err := rule(bitmask.New(PermDelete, PermRead)); !errors.Is(err, strict.ErrDependencyMissing) {
		t.Fatalf("expected ErrDependencyMissing, got %v", err)
	}

	// Delete set, has both Read and Write
	if err := rule(bitmask.New(PermDelete, PermRead, PermWrite)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestExactlyOne(t *testing.T) {
	rule := strict.ExactlyOne(PermRead, PermWrite, PermAdmin)

	// Empty rule call
	if err := strict.ExactlyOne[Permission]()(bitmask.New(PermRead)); err != nil {
		t.Fatalf("expected nil for empty rule, got %v", err)
	}

	// 0 flags set
	if err := rule(bitmask.New[Permission]()); !errors.Is(err, strict.ErrExactCountMismatch) {
		t.Fatalf("expected ErrExactCountMismatch for 0 flags, got %v", err)
	}

	// 1 flag set (Read)
	if err := rule(bitmask.New(PermRead)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// 2 flags set (Read, Write)
	if err := rule(bitmask.New(PermRead, PermWrite)); !errors.Is(err, strict.ErrExactCountMismatch) {
		t.Fatalf("expected ErrExactCountMismatch for 2 flags, got %v", err)
	}
}

func TestCountRules(t *testing.T) {
	maxRule := strict.MaxCount[Permission](2)
	minRule := strict.MinCount[Permission](2)

	m1 := bitmask.New(PermRead)
	m2 := bitmask.New(PermRead, PermWrite)
	m3 := bitmask.New(PermRead, PermWrite, PermAdmin)

	if err := maxRule(m2); err != nil {
		t.Fatalf("expected nil for max 2 with 2 flags, got %v", err)
	}
	if err := maxRule(m3); !errors.Is(err, strict.ErrCountExceeded) {
		t.Fatalf("expected ErrCountExceeded, got %v", err)
	}

	if err := minRule(m2); err != nil {
		t.Fatalf("expected nil for min 2 with 2 flags, got %v", err)
	}
	if err := minRule(m1); !errors.Is(err, strict.ErrCountInsufficient) {
		t.Fatalf("expected ErrCountInsufficient, got %v", err)
	}
}

func TestAllOfAndAnyOf(t *testing.T) {
	allOf := strict.AllOf(
		strict.RequireAll(PermRead),
		strict.NoneOf(PermGuest),
	)

	// Empty AllOf
	if err := strict.AllOf[Permission]()(bitmask.New[Permission]()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Passes both
	if err := allOf(bitmask.New(PermRead)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Fails first
	if err := allOf(bitmask.New[Permission]()); !errors.Is(err, strict.ErrMissingRequired) {
		t.Fatalf("expected ErrMissingRequired, got %v", err)
	}

	// Fails second
	if err := allOf(bitmask.New(PermRead, PermGuest)); !errors.Is(err, strict.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	// AnyOf
	anyOf := strict.AnyOf(
		strict.RequireAll(PermAdmin),
		strict.RequireAll(PermGuest),
	)

	// Empty AnyOf
	if err := strict.AnyOf[Permission]()(bitmask.New[Permission]()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Passes first
	if err := anyOf(bitmask.New(PermAdmin)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Passes second
	if err := anyOf(bitmask.New(PermGuest)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Fails both
	if err := anyOf(bitmask.New(PermRead)); err == nil {
		t.Fatalf("expected error when both fail, got nil")
	}
}

func TestWhenAndNot(t *testing.T) {
	whenRule := strict.When(PermAdmin, strict.RequireAll(PermWrite))

	// Condition not met (no Admin) -> succeeds even without Write
	if err := whenRule(bitmask.New(PermRead)); err != nil {
		t.Fatalf("expected nil when condition not met, got %v", err)
	}

	// Condition met (Admin set), but missing Write
	if err := whenRule(bitmask.New(PermAdmin)); !errors.Is(err, strict.ErrMissingRequired) {
		t.Fatalf("expected ErrMissingRequired, got %v", err)
	}

	// Condition met and requirement satisfied
	if err := whenRule(bitmask.New(PermAdmin, PermWrite)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Not
	notRule := strict.Not(strict.RequireAll(PermGuest))
	// Has guest -> inner rule succeeds -> Not fails
	if err := notRule(bitmask.New(PermGuest)); !errors.Is(err, strict.ErrInvertedRulePassed) {
		t.Fatalf("expected ErrInvertedRulePassed, got %v", err)
	}
	// Does not have guest -> inner rule fails -> Not succeeds
	if err := notRule(bitmask.New(PermRead)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCustom(t *testing.T) {
	customErr := errors.New("custom error")
	rule := strict.Custom(customErr, func(m bitmask.Mask[Permission]) bool {
		return m.Count()%2 == 0
	})

	if err := rule(bitmask.New(PermRead, PermWrite)); err != nil {
		t.Fatalf("expected nil for even count, got %v", err)
	}

	if err := rule(bitmask.New(PermRead)); !errors.Is(err, customErr) {
		t.Fatalf("expected customErr, got %v", err)
	}

	// Default error when err is nil
	defaultRule := strict.Custom(nil, func(m bitmask.Mask[Permission]) bool {
		return false
	})
	if err := defaultRule(bitmask.New[Permission]()); !errors.Is(err, strict.ErrDefaultError) {
		t.Fatalf("expected ErrDefaultError, got %v", err)
	}
}

// Test policy for Mask tests
type TestUserPolicy struct{}

func (TestUserPolicy) Validate(m bitmask.Mask[Permission]) error {
	return strict.AllOf(
		strict.MutuallyExclusive(PermRead, PermWrite),
		strict.NoneOf(PermGuest),
	)(m)
}

func TestStrictMask(t *testing.T) {
	// New: valid
	m, err := strict.New[Permission, TestUserPolicy](PermRead)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Has(PermRead) {
		t.Fatalf("expected mask to have PermRead")
	}

	// New: invalid (conflict)
	_, err = strict.New[Permission, TestUserPolicy](PermRead, PermWrite)
	if !errors.Is(err, strict.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}

	// Wrap & Check
	raw := bitmask.New(PermGuest)
	wrapped := strict.Wrap[Permission, TestUserPolicy](raw)
	if err := wrapped.Check(); !errors.Is(err, strict.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	// UnmarshalJSON valid
	type Payload struct {
		Perms strict.Mask[Permission, TestUserPolicy] `json:"perms"`
	}

	var p Payload
	validJSON := `{"perms": 1}` // PermRead (1)
	if err := json.Unmarshal([]byte(validJSON), &p); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if !p.Perms.Has(PermRead) {
		t.Fatalf("expected unmarshaled mask to have PermRead")
	}

	// UnmarshalJSON invalid
	invalidJSON := `{"perms": 3}` // PermRead (1) | PermWrite (2) = conflict!
	var pBad Payload
	err = json.Unmarshal([]byte(invalidJSON), &pBad)
	if !errors.Is(err, strict.ErrConflict) {
		t.Fatalf("expected ErrConflict on invalid JSON unmarshal, got %v", err)
	}

	// UnmarshalText valid & invalid
	var textMask strict.Mask[Permission, TestUserPolicy]
	if err := textMask.UnmarshalText([]byte("1")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := textMask.UnmarshalText([]byte("3")); !errors.Is(err, strict.ErrConflict) {
		t.Fatalf("expected ErrConflict on text unmarshal, got %v", err)
	}

	// Scan valid & invalid
	var sqlMask strict.Mask[Permission, TestUserPolicy]
	if err := sqlMask.Scan(int64(1)); err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	if err := sqlMask.Scan(int64(3)); !errors.Is(err, strict.ErrConflict) {
		t.Fatalf("expected ErrConflict on sql scan, got %v", err)
	}
}

func BenchmarkStrictCheck(b *testing.B) {
	m, _ := strict.New[Permission, TestUserPolicy](PermRead)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := m.Check(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMutuallyExclusive(b *testing.B) {
	rule := strict.MutuallyExclusive(PermRead, PermWrite)
	m := bitmask.New(PermRead)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := rule(m); err != nil {
			b.Fatal(err)
		}
	}
}
