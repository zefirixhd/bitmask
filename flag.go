package bitmask

import (
	"fmt"
	"math/bits"
)

// Flag is a constraint for unsigned integer types suitable for bitmasks.
// It supports all standard unsigned integers and custom types defined over them (~).
type Flag interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uintptr
}

// StringFlag is a constraint for unsigned integer flags that also implement fmt.Stringer.
// Used for automatic string conversion and parsing without manual maps.
type StringFlag interface {
	Flag
	fmt.Stringer
}

// Decompose breaks down a composite flag or mask value into a slice of atomic flags (powers of 2).
// It uses the hardware TZCNT instruction to find set bits in O(k) time, where k is the number of set bits.
//
// Example:
//   // val: 0b0101 (number 5)
//   // result: []T{1, 4} (0b0001, 0b0100)
//   atomicFlags := bitmask.Decompose(5)
func Decompose[T Flag](val T) []T {
	if val == 0 {
		return nil
	}
	if isPowerOfTwo(val) {
		return []T{val}
	}
	result := make([]T, 0, bits.OnesCount64(uint64(val)))
	for val != 0 {
		tz := bits.TrailingZeros64(uint64(val))
		atomic := T(1) << tz

		result = append(result, atomic)
		val &= val - 1 // Brian Kernighan's algorithm: clear the lowest set bit
	}
	return result
}

// isPowerOfTwo reports whether i is strictly a power of two (has exactly one bit set).
//
// Example:
//   // 1, 2, 4, 8, 16 -> true
//   // 0, 3, 5, 6, 7  -> false
func isPowerOfTwo[T Flag](i T) bool {
	if i == 0 {
		return false
	}
	return (i & (i - 1)) == 0
}
