package bitmask

import "math/bits"

// HasBit reports whether the bit at position idx (0 to 63) is set.
//
// Example:
//   // m.bits: 0b0100 (4)
//   // idx: 2 -> true (1 << 2 == 4)
//   // idx: 1 -> false
//   m.HasBit(2)
func (m Mask[T]) HasBit(idx int) bool {
	return m.bits&(T(1)<<idx) != 0
}

// SetBit sets the bit at position idx to 1 (m.bits |= 1 << idx).
//
// Example:
//   // m.bits: 0b0001 (1)
//   // idx: 1
//   // result: 0b0011 (3)
//   m.SetBit(1)
func (m *Mask[T]) SetBit(idx int) *Mask[T] {
	m.bits |= T(1) << idx
	return m
}

// UnsetBit clears the bit at position idx to 0 (m.bits &^= 1 << idx).
//
// Example:
//   // m.bits: 0b0011 (3)
//   // idx: 1
//   // result: 0b0001 (1)
//   m.UnsetBit(1)
func (m *Mask[T]) UnsetBit(idx int) *Mask[T] {
	m.bits &^= T(1) << idx
	return m
}

// ToggleBit flips the bit at position idx (m.bits ^= 1 << idx).
//
// Example:
//   // m.bits: 0b0001 (1)
//   // idx: 1
//   // result: 0b0011 (3)
//   m.ToggleBit(1)
func (m *Mask[T]) ToggleBit(idx int) *Mask[T] {
	m.bits ^= T(1) << idx
	return m
}

// Count returns the number of set bits in the mask (Hamming weight / population count).
// It runs in a single CPU hardware instruction (POPCNT).
//
// Example:
//   // m.bits: 0b0101 (bits 0 and 2) -> 2
//   m.Count()
func (m Mask[T]) Count() int {
	return bits.OnesCount64(uint64(m.bits))
}

// TrailingZeros returns the number of consecutive zero bits before the least significant set bit.
// Returns 64 if the mask is empty.
//
// Example:
//   // m.bits: 0b1000 (8) -> 3 trailing zeros
//   m.TrailingZeros()
func (m Mask[T]) TrailingZeros() int {
	return bits.TrailingZeros64(uint64(m.bits))
}

// LeadingZeros returns the number of consecutive zero bits before the most significant set bit in 64-bit representation.
//
// Example:
//   m.LeadingZeros()
func (m Mask[T]) LeadingZeros() int {
	return bits.LeadingZeros64(uint64(m.bits))
}
