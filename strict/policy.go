package strict

import "github.com/zefirixhd/bitmask"

// Policy defines the compile-time validation contract for a strict Mask.
//
// Implementing types are typically zero-sized structs (struct{}) that incur 0 heap allocations
// and zero runtime overhead.
//
// Example:
//
//	type AdminPolicy struct{}
//
//	func (AdminPolicy) Validate(m bitmask.Mask[Permission]) error {
//		return strict.AllOf(
//			strict.MutuallyExclusive(PermissionReadOnly, PermissionWriteOnly),
//			strict.DependsOn(PermissionDelete, PermissionAdmin),
//		)(m)
//	}
type Policy[F bitmask.Flag] interface {
	Validate(m bitmask.Mask[F]) error
}
