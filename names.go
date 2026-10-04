package bitmask

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// Names returns a slice of string names for all active flags in the mask using the provided names map.
// If a flag is not found in the names map, its numeric value is used as fallback.
//
// Example:
//
//	names := map[Permission]string{PermissionRead: "read", PermissionWrite: "write"}
//	// perms: 0b0011
//	// result: []string{"read", "write"}
//	perms.Names(names)
func (m Mask[T]) Names(names map[T]string) []string {
	flags := m.Flags()
	resp := make([]string, len(flags))
	for i, f := range flags {
		name, ok := names[f]
		if !ok {
			name = strconv.FormatUint(uint64(f), 10)
		}
		resp[i] = name
	}
	return resp
}

// Format formats the active flags into a single string joined by the separator using the provided map.
//
// Example:
//
//	// result: "read, write"
//	perms.Format(names, ", ")
func (m Mask[T]) Format(names map[T]string, separator string) string {
	return strings.Join(m.Names(names), separator)
}

// Strings returns a slice of string names for all active flags using type T's String() method.
// If type T does not implement fmt.Stringer, decimal string representation is used as fallback.
// No manual map needed!
//
// Example:
//
//	// perms: 0b0011
//	// result: []string{"read", "write"}
//	perms.Strings()
func (m Mask[T]) Strings() []string {
	flags := m.Flags()
	resp := make([]string, len(flags))
	for i, f := range flags {
		if s, ok := any(f).(fmt.Stringer); ok {
			resp[i] = s.String()
		} else {
			resp[i] = strconv.FormatUint(uint64(f), 10)
		}
	}
	return resp
}

// FormatAuto formats the active flags into a single string joined by separator using type T's String() method.
//
// Example:
//
//	// result: "read, write"
//	perms.FormatAuto(", ")
func (m Mask[T]) FormatAuto(separator string) string {
	return strings.Join(m.Strings(), separator)
}

// Parse converts a slice of string names (e.g. from CLI arguments)
// into a Mask[T] using the provided names map.
// It returns an error if any name is not found in the map.
//
// Example:
//
//	// CLI flags: []string{"read", "write"}
//	mask, err := bitmask.Parse([]string{"read", "write"}, names)
func Parse[T Flag](strs []string, names map[T]string) (Mask[T], error) {
	lookup := make(map[string]T, len(names))
	for flag, name := range names {
		lookup[name] = flag
	}
	var m Mask[T]
	for _, s := range strs {
		flag, ok := lookup[s]
		if !ok {
			return m, fmt.Errorf("bitmask: unknown flag %q", s)
		}
		m.Set(flag)
	}
	return m, nil
}

// ParseAuto converts a slice of string names into a Mask[T] automatically
// by probing all possible bit positions (0..width-1) using type T's String() method.
// Type T must implement StringFlag (Flag + fmt.Stringer).
// No manual map needed!
//
// Example:
//
//	mask, err := bitmask.ParseAuto[Permission]("read", "write")
func ParseAuto[T StringFlag](strs ...string) (Mask[T], error) {
	var zero T
	width := int(unsafe.Sizeof(zero) * 8)
	lookup := make(map[string]T, width)
	for i := range width {
		flag := T(1) << i
		lookup[flag.String()] = flag
	}

	var m Mask[T]
	for _, s := range strs {
		flag, ok := lookup[s]
		if !ok {
			return m, fmt.Errorf("bitmask: unknown flag %q for type %T", s, zero)
		}
		m.Set(flag)
	}
	return m, nil
}
