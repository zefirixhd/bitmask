// Package strict provides constraint-based, validated bitmasks with compile-time policies
// and composable validation rules for the bitmask library.
//
// It enables declarative validation rules (mutual exclusion, requirements, dependencies,
// limits, and conditions) that can be checked on demand or automatically enforced during
// JSON/Text unmarshaling and database scanning.
//
// By using zero-sized struct policies, strict.Mask incurs zero heap allocations during
// validation while remaining fully type-safe.
package strict
