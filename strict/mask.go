package strict

import (
	"github.com/zefirixhd/bitmask"
)

// Mask is a type-safe, validated wrapper around bitmask.Mask[F] that enforces
// the business rules defined by Policy P.
//
// Mask embeds bitmask.Mask[F], giving it full access to all standard bitwise
// operations (Has, Set, Unset, Count, Bits, etc.). It automatically validates its
// state during Check(), UnmarshalJSON(), UnmarshalText(), and database Scan().
//
// Example:
//
//	type Payload struct {
//		UserID int                           `json:"user_id"`
//		Perms  strict.Mask[uint32, UserRules] `json:"perms"`
//	}
type Mask[F bitmask.Flag, P Policy[F]] struct {
	bitmask.Mask[F]
}

// New creates a new strict Mask initialized with the given flags and immediately
// validates it against Policy P. If any rule in P is violated, an error is returned.
//
// Example:
//
//	mask, err := strict.New[Permission, UserPolicy](PermissionRead, PermissionWrite)
//	if err != nil {
//		log.Fatal(err)
//	}
func New[F bitmask.Flag, P Policy[F]](flags ...F) (Mask[F, P], error) {
	m := Wrap[F, P](bitmask.New(flags...))
	if err := m.Check(); err != nil {
		return Mask[F, P]{}, err
	}
	return m, nil
}

// Wrap wraps an existing bitmask into a strict Mask without immediate validation.
// Call Check() to validate the mask.
//
// Example:
//
//	rawMask := bitmask.New(PermissionRead)
//	guarded := strict.Wrap[Permission, UserPolicy](rawMask)
//	if err := guarded.Check(); err != nil { ... }
func Wrap[F bitmask.Flag, P Policy[F]](m bitmask.Mask[F]) Mask[F, P] {
	return Mask[F, P]{Mask: m}
}

// Check validates the current state of the mask against Policy P.
// It returns nil if the mask complies with all policy rules, or an error if validation fails.
//
// Example:
//
//	if err := userMask.Check(); err != nil {
//		// handle rule violation
//	}
func (s Mask[F, P]) Check() error {
	var p P
	return p.Validate(s.Mask)
}

// UnmarshalJSON deserializes the JSON-encoded data into the mask and automatically
// validates it against Policy P via Check().
//
// If the deserialized bits violate the policy, an error is returned and the unmarshaling fails.
func (s *Mask[F, P]) UnmarshalJSON(data []byte) error {
	if err := s.Mask.UnmarshalJSON(data); err != nil {
		return err
	}
	return s.Check()
}

// UnmarshalText deserializes text data (hex or binary) into the mask and automatically
// validates it against Policy P via Check().
func (s *Mask[F, P]) UnmarshalText(text []byte) error {
	if err := s.Mask.UnmarshalText(text); err != nil {
		return err
	}
	return s.Check()
}

// Scan implements the database/sql.Scanner interface, deserializing from SQL integer/byte
// columns and automatically validating against Policy P via Check().
func (s *Mask[F, P]) Scan(src any) error {
	if err := s.Mask.Scan(src); err != nil {
		return err
	}
	return s.Check()
}
