package strict

import (
	"github.com/zefirixhd/bitmask"
)

// Rule defines a validation predicate for a bitmask.
// It inspects the mask m and returns nil if the mask complies with the rule,
// or an error explaining why validation failed.
type Rule[F bitmask.Flag] func(m bitmask.Mask[F]) error

// Custom creates a validation rule from a boolean predicate function.
// If the predicate returns false, Custom returns err. If err is nil,
// ErrDefaultError is returned by default.
//
// Example:
//
//	rule := strict.Custom(nil, func(m bitmask.Mask[uint32]) bool {
//		return m.Bits()%2 == 0 // only even bit patterns allowed
//	})
func Custom[F bitmask.Flag](err error, rule func(m bitmask.Mask[F]) bool) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if err == nil {
			err = ErrDefaultError
		}
		if !rule(m) {
			return err
		}
		return nil
	}
}

// MutuallyExclusive ensures that at most one flag from the given list v is present in the mask.
// If two or more flags from v are set, ErrConflict is returned.
//
// Example:
//
//	// A user cannot simultaneously have ReadOnly and WriteOnly flags:
//	rule := strict.MutuallyExclusive(ReadOnly, WriteOnly)
func MutuallyExclusive[F bitmask.Flag](v ...F) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if len(v) == 0 {
			return nil
		}
		var seen bool
		for _, f := range v {
			if m.Has(f) {
				if seen {
					return ErrConflict
				}
				seen = true
			}
		}
		return nil
	}
}

// RequireAll ensures that all specified flags v are present in the mask.
// If any flag from v is missing, ErrMissingRequired is returned.
//
// Example:
//
//	// A SuperAdmin must have both ManageUsers and ManageBilling:
//	rule := strict.RequireAll(ManageUsers, ManageBilling)
func RequireAll[F bitmask.Flag](v ...F) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if len(v) == 0 {
			return nil
		}
		if !m.HasAll(v...) {
			return ErrMissingRequired
		}
		return nil
	}
}

// RequireAny ensures that at least one flag from the specified list v is present in the mask.
// If none of the flags from v are set, ErrMissingRequired is returned.
//
// Example:
//
//	// An authenticated user must have at least one role assigned:
//	rule := strict.RequireAny(RoleUser, RoleModerator, RoleAdmin)
func RequireAny[F bitmask.Flag](v ...F) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if len(v) == 0 {
			return nil
		}
		if !m.HasAny(v...) {
			return ErrMissingRequired
		}
		return nil
	}
}

// NoneOf ensures that none of the specified flags v are present in the mask.
// If any flag from v is set, ErrForbidden is returned.
//
// Example:
//
//	// A standard user must not hold Suspended or Banned flags:
//	rule := strict.NoneOf(StatusSuspended, StatusBanned)
func NoneOf[F bitmask.Flag](v ...F) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if len(v) == 0 {
			return nil
		}
		if m.HasAny(v...) {
			return ErrForbidden
		}
		return nil
	}
}

// DependsOn ensures that if flag f is present in the mask, all requirement flags
// must also be present (implication: f => requirements).
// If f is not set, the rule always succeeds.
//
// Example:
//
//	// PermDelete requires both PermRead and PermWrite to be granted:
//	rule := strict.DependsOn(PermDelete, PermRead, PermWrite)
func DependsOn[F bitmask.Flag](f F, requirements ...F) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if m.Has(f) && !m.HasAll(requirements...) {
			return ErrDependencyMissing
		}
		return nil
	}
}

// ExactlyOne ensures that exactly one flag from the given list v is set in the mask.
// If zero flags or more than one flag are set, ErrExactCountMismatch is returned.
//
// Example:
//
//	// An order must be in exactly one state:
//	rule := strict.ExactlyOne(StatePending, StateProcessing, StateCompleted)
func ExactlyOne[F bitmask.Flag](v ...F) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if len(v) == 0 {
			return nil
		}
		var count int
		for _, f := range v {
			if m.Has(f) {
				count++
				if count > 1 {
					return ErrExactCountMismatch
				}
			}
		}
		if count == 0 {
			return ErrExactCountMismatch
		}
		return nil
	}
}

// MaxCount ensures that the total number of set bits (population count) in the mask
// does not exceed threshold. If m.Count() > threshold, ErrCountExceeded is returned.
//
// Example:
//
//	// A free tier user can have at most 3 feature flags enabled:
//	rule := strict.MaxCount[Permission](3)
func MaxCount[F bitmask.Flag](threshold int) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if m.Count() > threshold {
			return ErrCountExceeded
		}
		return nil
	}
}

// MinCount ensures that the total number of set bits (population count) in the mask
// is at least threshold. If m.Count() < threshold, ErrCountInsufficient is returned.
//
// Example:
//
//	// A valid security configuration requires at least 2 active flags:
//	rule := strict.MinCount[Permission](2)
func MinCount[F bitmask.Flag](threshold int) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if m.Count() < threshold {
			return ErrCountInsufficient
		}
		return nil
	}
}

// AllOf combines multiple rules into a single rule using logical AND.
// All rules must pass; validation fails and returns the error of the first rule that fails.
//
// Example:
//
//	policy := strict.AllOf(
//		strict.RequireAny(RoleUser, RoleAdmin),
//		strict.NoneOf(StatusBanned),
//		strict.DependsOn(PermWrite, PermRead),
//	)
func AllOf[F bitmask.Flag](r ...Rule[F]) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if len(r) == 0 {
			return nil
		}
		for _, rule := range r {
			if err := rule(m); err != nil {
				return err
			}
		}
		return nil
	}
}

// AnyOf combines multiple rules using logical OR.
// If at least one rule passes, validation succeeds (returns nil).
// If all rules fail, the error from the last failing rule is returned.
//
// Example:
//
//	// Either the user is an Admin, or they are a Guest with Read-only access:
//	rule := strict.AnyOf(
//		strict.RequireAll(RoleAdmin),
//		strict.AllOf(strict.RequireAll(RoleGuest), strict.NoneOf(PermWrite)),
//	)
func AnyOf[F bitmask.Flag](r ...Rule[F]) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if len(r) == 0 {
			return nil
		}
		var err error
		for _, rule := range r {
			err = rule(m)
			if err == nil {
				return nil
			}
		}
		if err == nil {
			return ErrNoRulePassed
		}
		return err
	}
}

// When applies rule only if conditionFlag is present in the mask.
// If conditionFlag is not set, the rule is skipped and returns nil.
//
// Example:
//
//	// If the user has BetaTester flag, they cannot also have ProductionOnly:
//	rule := strict.When(BetaTester, strict.NoneOf(ProductionOnly))
func When[F bitmask.Flag](conditionFlag F, rule Rule[F]) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if m.Has(conditionFlag) {
			return rule(m)
		}
		return nil
	}
}

// Not inverts the outcome of rule:
// If rule returns nil (succeeded), Not returns ErrInvertedRulePassed.
// If rule returns any error, Not succeeds and returns nil.
//
// Example:
//
//	// Invert a rule requirement:
//	rule := strict.Not(strict.RequireAll(RestrictedFlag))
func Not[F bitmask.Flag](rule Rule[F]) Rule[F] {
	return func(m bitmask.Mask[F]) error {
		if rule(m) == nil {
			return ErrInvertedRulePassed
		}
		return nil
	}
}
