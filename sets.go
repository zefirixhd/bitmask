package bitmask

// And returns a new mask that is the intersection of m and other (bitwise AND: m.bits & other.bits).
// Only bits that are set in BOTH masks are retained.
//
// Example:
//
//	// m:      0b0101 (1, 4)
//	// other:  0b0011 (1, 2)
//	// result: 0b0001 (1)
//	m.And(other)
func (m Mask[T]) And(other Mask[T]) Mask[T] {
	return Wrap(m.bits & other.bits)
}

// Or returns a new mask that is the union of m and other (bitwise OR: m.bits | other.bits).
// Any bit set in EITHER mask is set in the result.
//
// Example:
//
//	// m:      0b0100 (4)
//	// other:  0b0001 (1)
//	// result: 0b0101 (5)
//	m.Or(other)
func (m Mask[T]) Or(other Mask[T]) Mask[T] {
	return Wrap(m.bits | other.bits)
}

// Xor returns a new mask that is the symmetric difference of m and other (bitwise XOR: m.bits ^ other.bits).
// Only bits set in EXACTLY ONE of the masks are set in the result.
//
// Example:
//
//	// m:      0b0110 (2, 4)
//	// other:  0b0011 (1, 2)
//	// result: 0b0101 (1, 4) — common bit 2 is excluded
//	m.Xor(other)
func (m Mask[T]) Xor(other Mask[T]) Mask[T] {
	return Wrap(m.bits ^ other.bits)
}

// AndNot returns a new mask that is the difference of m and other (bitwise AND NOT: m.bits &^ other.bits).
// It clears in m all bits that are set in other.
//
// Example:
//
//	// m:      0b0111 (1, 2, 4)
//	// other:  0b0010 (2)
//	// result: 0b0101 (1, 4)
//	m.AndNot(other)
func (m Mask[T]) AndNot(other Mask[T]) Mask[T] {
	return Wrap(m.bits &^ other.bits)
}

// Contains reports whether mask m contains ALL bits from other (other ⊆ m).
// It returns false if other is empty.
//
// Example:
//
//	// m:      0b0111 (has flags 1, 2, 4)
//	// other:  0b0011 (flags 1, 2)
//	// result: true (all bits of other are in m)
//	m.Contains(other)
func (m Mask[T]) Contains(other Mask[T]) bool {
	return m.Has(other.bits)
}

// Overlaps reports whether m and other share AT LEAST ONE common bit.
//
// Example:
//
//	// m:      0b0100 (4)
//	// other:  0b0101 (1, 4)
//	// result: true (bit 4 is common)
//	m.Overlaps(other)
func (m Mask[T]) Overlaps(other Mask[T]) bool {
	return (m.bits & other.bits) != 0
}
