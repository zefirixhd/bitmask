package bitmask

import (
	"iter"
	"math/bits"
)

// All returns a Go 1.23+ iterator over all active atomic flags (1, 2, 4...).
// It runs with zero heap allocations (0 B/op) and supports early exit via break.
//
// Example:
//   // perms.bits: 0b0101 (PermissionRead | PermissionDelete)
//   for flag := range perms.All() {
//       // 1st iteration: flag = 1 (PermissionRead)
//       // 2nd iteration: flag = 4 (PermissionDelete)
//       fmt.Println("Active flag:", flag)
//   }
func (m Mask[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		v := uint64(m.bits)
		for v != 0 {
			tz := bits.TrailingZeros64(v)
			atomic := T(1) << tz

			if !yield(atomic) {
				return
			}

			v &= v - 1
		}
	}
}

// BitIndices returns a Go 1.23+ iterator over the index positions of all set bits (0, 1, 2... 63).
//
// Example:
//   // perms.bits: 0b0101 (bits at positions 0 and 2)
//   for idx := range perms.BitIndices() {
//       // 1st iteration: idx = 0
//       // 2nd iteration: idx = 2
//       fmt.Println("Set bit index:", idx)
//   }
func (m Mask[T]) BitIndices() iter.Seq[int] {
	return func(yield func(int) bool) {
		v := uint64(m.bits)
		for v != 0 {
			tz := bits.TrailingZeros64(v)

			if !yield(tz) {
				return
			}

			v &= v - 1
		}
	}
}
