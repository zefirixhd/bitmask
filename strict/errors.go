package strict

import "errors"

var (
	// ErrDefaultError is the fallback error returned by Custom when no specific error is provided.
	ErrDefaultError = errors.New("strict: failed to validate bitmask value")

	// ErrLogicError is a generic sentinel error for validation failures (kept for backwards compatibility).
	ErrLogicError = ErrDefaultError

	// ErrConflict is returned when mutually exclusive flags are set simultaneously.
	ErrConflict = errors.New("strict: mutually exclusive flags present")

	// ErrMissingRequired is returned when one or more required flags are missing.
	ErrMissingRequired = errors.New("strict: required flags are missing")

	// ErrForbidden is returned when a forbidden flag is present in the mask.
	ErrForbidden = errors.New("strict: forbidden flags present")

	// ErrDependencyMissing is returned when a dependent flag is set without its prerequisite flags.
	ErrDependencyMissing = errors.New("strict: dependent flag requires missing prerequisite flag")

	// ErrExactCountMismatch is returned when the number of set flags does not match the exact requirement.
	ErrExactCountMismatch = errors.New("strict: exactly one flag must be set")

	// ErrCountExceeded is returned when the number of set flags exceeds the allowed maximum.
	ErrCountExceeded = errors.New("strict: active flag count exceeds maximum threshold")

	// ErrCountInsufficient is returned when the number of set flags is below the allowed minimum.
	ErrCountInsufficient = errors.New("strict: active flag count below minimum threshold")

	// ErrInvertedRulePassed is returned by Not when the negated rule unexpectedly succeeded.
	ErrInvertedRulePassed = errors.New("strict: negated rule unexpectedly succeeded")

	// ErrNoRulePassed is returned by AnyOf when none of the candidate rules succeeded.
	ErrNoRulePassed = errors.New("strict: none of the alternative rules passed")
)
