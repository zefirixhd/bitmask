package bitmask_test

import (
	"encoding/json"
	"testing"

	"github.com/zefirixhd/bitmask"
)

func TestEncoding(t *testing.T) {
	m := bitmask.New(Permission(1), Permission(4)) // 0b00000101 = 5

	if m.String() != "5" {
		t.Errorf("expected String() '5', got %q", m.String())
	}
	if m.Binary() != "00000101" {
		t.Errorf("expected Binary() '00000101', got %q", m.Binary())
	}
	if m.Hex() != "05" {
		t.Errorf("expected Hex() '05', got %q", m.Hex())
	}

	// JSON
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	if string(data) != "5" {
		t.Errorf("expected JSON '5', got %q", string(data))
	}

	var unmarshaled bitmask.Mask[Permission]
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if unmarshaled.Bits() != m.Bits() {
		t.Errorf("expected unmarshaled bits %d, got %d", m.Bits(), unmarshaled.Bits())
	}

	// JSON null
	var nullMask bitmask.Mask[Permission]
	nullMask.Set(PermissionLogin)
	if err := json.Unmarshal([]byte("null"), &nullMask); err != nil {
		t.Fatalf("json.Unmarshal 'null' failed: %v", err)
	}
	if !nullMask.IsEmpty() {
		t.Errorf("expected null to clear mask")
	}

	// Text marshaler & unmarshaler
	textData, err := m.MarshalText()
	if err != nil || string(textData) != "5" {
		t.Errorf("MarshalText failed: %v, text: %s", err, textData)
	}
	var textMask bitmask.Mask[Permission]
	if err := textMask.UnmarshalText([]byte("5")); err != nil || textMask.Bits() != 5 {
		t.Errorf("UnmarshalText failed: %v, bits: %d", err, textMask.Bits())
	}

	// SQL Valuer & Scanner
	val, err := m.Value()
	if err != nil || val != int64(5) {
		t.Errorf("Value() failed: %v, val: %v", err, val)
	}

	var scanned bitmask.Mask[Permission]
	if err := scanned.Scan(int64(5)); err != nil || scanned.Bits() != 5 {
		t.Errorf("Scan() int64 failed: %v, bits: %d", err, scanned.Bits())
	}
	if err := scanned.Scan([]byte("5")); err != nil || scanned.Bits() != 5 {
		t.Errorf("Scan() []byte failed: %v, bits: %d", err, scanned.Bits())
	}
	if err := scanned.Scan("5"); err != nil || scanned.Bits() != 5 {
		t.Errorf("Scan() string failed: %v, bits: %d", err, scanned.Bits())
	}
	if err := scanned.Scan(nil); err != nil || !scanned.IsEmpty() {
		t.Errorf("Scan() nil failed: %v", err)
	}
}

func BenchmarkMask_Binary(b *testing.B) {
	m := bitmask.New(PermissionLogin, PermissionWrite)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Binary()
	}
}

func BenchmarkMask_Hex(b *testing.B) {
	m := bitmask.New(PermissionLogin, PermissionWrite)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Hex()
	}
}
