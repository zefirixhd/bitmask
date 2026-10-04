package bitmask

// Mask represents a type-safe, lightweight wrapper around an integer bitmask.
// T must satisfy the Flag constraint (any unsigned integer type).
type Mask[T Flag] struct {
	bits T
}

// New creates a new Mask with the specified flags set.
//
// Example:
//   perms := bitmask.New(PermissionRead, PermissionWrite)
func New[T Flag](flags ...T) Mask[T] {
	var m Mask[T]
	m.Set(flags...)
	return m
}

// Wrap creates a Mask directly from an integer bit value.
//
// Example:
//   perms := bitmask.Wrap(PermissionRead | PermissionWrite)
func Wrap[T Flag](i T) Mask[T] {
	return Mask[T]{
		bits: i,
	}
}

// IsEmpty reports whether no bits are set in the mask (value is 0).
func (m Mask[T]) IsEmpty() bool {
	return m.bits == 0
}

// Bits returns the raw underlying integer value of the mask.
func (m Mask[T]) Bits() T {
	return m.bits
}

// Reset clears all bits in the mask, setting its value to 0.
//
// Example:
//   m.Reset() // m.IsEmpty() == true
func (m *Mask[T]) Reset() *Mask[T] {
	m.bits = T(0)
	return m
}

// Flags decomposes the mask into a slice of individual atomic flags (powers of 2).
//
// Example:
//   // m.bits: 0b0101 (1, 4)
//   // result: []Permission{PermissionRead, PermissionDelete}
//   flags := m.Flags()
func (m Mask[T]) Flags() []T {
	return Decompose(m.bits)
}
