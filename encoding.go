package bitmask

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"unsafe"
)

// String returns the decimal string representation of the mask.
//
// Example:
//   // m.bits: 0b0101 -> "5"
//   m.String()
func (m Mask[T]) String() string {
	return strconv.FormatUint(uint64(m.bits), 10)
}

// Binary returns a binary string padded with leading zeros to the full bit width of type T.
// For uint8 it returns 8 characters, for uint16 — 16, for uint64 — 64.
//
// Example for uint8:
//   // m.bits: 0b00000101 (5)
//   // result: "00000101"
//   m.Binary()
func (m Mask[T]) Binary() string {
	width := int(unsafe.Sizeof(m.bits) * 8)
	buf := make([]byte, width)
	v := uint64(m.bits)
	for i := width - 1; i >= 0; i-- {
		buf[i] = byte('0' + (v & 1))
		v >>= 1
	}
	return string(buf)
}

// Hex returns a hexadecimal string padded to the byte width of type T.
// For uint8 — 2 characters ("05"), for uint16 — 4 ("0005").
//
// Example:
//   // m.bits: 0b00001111 (15) -> "0f"
//   m.Hex()
func (m Mask[T]) Hex() string {
	width := int(unsafe.Sizeof(m.bits) * 2)
	return fmt.Sprintf("%0*x", width, uint64(m.bits))
}

// Value implements the database/sql/driver.Valuer interface.
// It allows SQL databases (PostgreSQL, MySQL, SQLite, etc.) to store the mask as an integer.
func (m Mask[T]) Value() (driver.Value, error) {
	return int64(m.bits), nil
}

// Scan implements the database/sql.Scanner interface.
// It reads a numeric database value back into the Mask.
func (m *Mask[T]) Scan(value any) error {
	if value == nil {
		m.bits = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		m.bits = T(v)
	case uint64:
		m.bits = T(v)
	case []byte:
		str := string(v)
		num, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return fmt.Errorf("bitmask: cannot scan %s into Mask", str)
		}
		m.bits = T(num)
	case string:
		num, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return fmt.Errorf("bitmask: cannot scan %s into Mask", v)
		}
		m.bits = T(num)
	default:
		return fmt.Errorf("bitmask: cannot scan %T into Mask", value)
	}
	return nil
}

// MarshalJSON encodes the mask as a raw JSON number (without quotes).
//
// Example:
//   // {"perms": 5}
func (m Mask[T]) MarshalJSON() ([]byte, error) {
	return strconv.AppendUint(nil, uint64(m.bits), 10), nil
}

// UnmarshalJSON decodes the mask from JSON.
// It properly handles numbers, null, and empty strings.
func (m *Mask[T]) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" || str == "" {
		m.bits = 0
		return nil
	}
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("bitmask: cannot unmarshal %s into Mask", str)
	}
	m.bits = T(num)
	return nil
}

// MarshalText encodes the mask into text format (for YAML, TOML, query params).
func (m Mask[T]) MarshalText() ([]byte, error) {
	return m.MarshalJSON()
}

// UnmarshalText decodes the mask from text format (for YAML, TOML, query params).
func (m *Mask[T]) UnmarshalText(text []byte) error {
	return m.UnmarshalJSON(text)
}
