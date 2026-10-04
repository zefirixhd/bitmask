package bitmask

// Has reports whether all bits in flag are set: (m.bits & flag) == flag.
// It returns false if flag is 0 or if the mask is empty.
//
// Example:
//   // m.bits: 0b0011 (Read and Write)
//   // flag:   0b0010 (Write)
//   // (0b0011 & 0b0010) == 0b0010 -> true
//   m.Has(PermissionWrite)
func (m Mask[T]) Has(flag T) bool {
	if flag == 0 || m.IsEmpty() {
		return false
	}
	return (m.bits & flag) == flag
}

// Set sets one or more flags (bitwise OR: m.bits |= flag).
//
// Example:
//   // m.bits: 0b0001 (1)
//   // flags:  0b0010 (2)
//   // result: 0b0011 (3)
//   m.Set(PermissionWrite)
func (m *Mask[T]) Set(flags ...T) *Mask[T] {
	for _, f := range flags {
		m.bits |= f
	}
	return m
}

// Unset clears one or more flags (bitwise AND NOT: m.bits &^= flag).
//
// Example:
//   // m.bits: 0b0011 (3)
//   // flags:  0b0010 (2)
//   // result: 0b0001 (1)
//   m.Unset(PermissionWrite)
func (m *Mask[T]) Unset(flags ...T) *Mask[T] {
	for _, f := range flags {
		m.bits &^= f
	}
	return m
}

// Toggle flips the state of the specified flags (bitwise XOR: m.bits ^= flag).
// Active flags become inactive, and inactive flags become active.
//
// Example:
//   // m.bits: 0b0001 (Read is on, Write is off)
//   // flags:  0b0011 (Read | Write)
//   // result: 0b0010 (Read turned off, Write turned on)
//   m.Toggle(PermissionRead, PermissionWrite)
func (m *Mask[T]) Toggle(flags ...T) *Mask[T] {
	for _, f := range flags {
		m.bits ^= f
	}
	return m
}

// HasAny reports whether at least one of the given flags is set.
//
// Example:
//   // m.bits: 0b0001 (only Read)
//   // flags:  PermissionRead, PermissionWrite
//   // result: true (because Read is set)
//   m.HasAny(PermissionRead, PermissionWrite)
func (m Mask[T]) HasAny(flags ...T) bool {
	for _, f := range flags {
		if (m.bits & f) != 0 {
			return true
		}
	}
	return false
}

// HasAll reports whether all of the given flags are set.
//
// Example:
//   // m.bits: 0b0001 (only Read)
//   // flags:  PermissionRead, PermissionWrite
//   // result: false (because Write is missing)
//   m.HasAll(PermissionRead, PermissionWrite)
func (m Mask[T]) HasAll(flags ...T) bool {
	for _, f := range flags {
		if (m.bits & f) == 0 {
			return false
		}
	}
	return true
}

// Invert flips all bits in the mask (bitwise NOT: m.bits = ^m.bits).
//
// Example for uint8:
//   // m.bits: 0b00000001 (1)
//   // result: 0b11111110 (254)
//   m.Invert()
func (m *Mask[T]) Invert() *Mask[T] {
	m.bits = ^m.bits
	return m
}

// With returns a copy of the mask with the given flags added (immutable).
// The original mask m is not modified.
//
// Example:
//   adminPerms := userPerms.With(PermissionDelete)
func (m Mask[T]) With(flags ...T) Mask[T] {
	m.Set(flags...)
	return m
}

// Without returns a copy of the mask with the given flags removed (immutable).
// The original mask m is not modified.
//
// Example:
//   readOnlyPerms := userPerms.Without(PermissionWrite, PermissionDelete)
func (m Mask[T]) Without(flags ...T) Mask[T] {
	m.Unset(flags...)
	return m
}
