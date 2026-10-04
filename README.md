# bitmask

[![Go Reference](https://pkg.go.dev/badge/github.com/zefirixhd/bitmask.svg)](https://pkg.go.dev/github.com/zefirixhd/bitmask)
[![Go Report Card](https://goreportcard.com/badge/github.com/zefirixhd/bitmask)](https://goreportcard.com/report/github.com/zefirixhd/bitmask)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Zero Allocations](https://img.shields.io/badge/allocs-0%20B%2Fop-brightgreen)](mask_test.go)

`bitmask` is a high-performance, type-safe, generic bitmask library for **Go 1.23+**. 

Designed for high-throughput backends, Role-Based Access Control (RBAC), and remote CLI management tools where performance, type safety, and zero memory allocations matter.

---

## Features

- рџђ **Sub-nanosecond Performance:** Core operations (`Has`, `Overlaps`, `Contains`) execute in ~0.30вЂ“0.90 ns using hardware CPU instructions (`POPCNT`, `TZCNT`).
- рџђ **Zero Heap Allocations:** All core mask operations, bit twiddling, and iterators run at `0 B/op` and `0 allocs/op`.
- рџђ **Type-Safe Generics:** Works seamlessly with custom enum types (`type Permission uint64`) via Go generics (`~uint...`).
- рџђ **Database & JSON Ready:** First-class support for `database/sql` (`sql.Scanner` / `driver.Valuer`), `json.Marshaler`, and `encoding.TextMarshaler` (YAML / TOML / query params).
- рџђ **Go 1.23+ Native Iterators:** Efficient iteration over set flags (`for flag := range mask.All()`) without allocating intermediate slices.
- рџђ **CLI-Ready String Mapping:** Automatic parsing and formatting (`ParseAuto`, `Strings`) using your enum's `String()` method without manual map declarations.

---

## Installation

```bash
go get github.com/zefirixhd/bitmask
```

*Requires Go 1.23 or newer.*

---

## Quick Start

Define your permissions using standard Go `iota` constants:

```go
package main

import (
	"fmt"
	"github.com/zefirixhd/bitmask"
)

type Permission uint8

const (
	PermissionLogin Permission = 1 << iota
	PermissionRead
	PermissionWrite
	PermissionDelete
)

func main() {
	// Create a mask with initial permissions
	perms := bitmask.New(PermissionRead, PermissionLogin)

	// Check permissions
	if perms.Has(PermissionLogin) {
		fmt.Println("User can login")
	}

	// Grant permissions (chainable)
	perms.Set(PermissionWrite)

	// Revoke permissions
	perms.Unset(PermissionRead)

	// Toggle state
	perms.Toggle(PermissionDelete)

	fmt.Printf("Active flags count: %d\n", perms.Count())
	fmt.Printf("Binary representation: %s\n", perms.Binary()) // "00001010"
}
```

---

## Real-World RBAC & Database Integration

`bitmask.Mask` implements standard Go SQL and JSON interfaces, allowing zero-boilerplate persistence:

```go
type User struct {
	ID    int64                       `json:"id" db:"id"`
	Email string                      `json:"email" db:"email"`
	Perms bitmask.Mask[Permission]    `json:"perms" db:"perms"`
}

// 1. Saved directly as a BIGINT in PostgreSQL/MySQL/SQLite:
// db.Exec("INSERT INTO users (email, perms) VALUES (?, ?)", user.Email, user.Perms)

// 2. Scanned directly from database rows:
// db.QueryRow("SELECT id, email, perms FROM users WHERE id = ?", id).Scan(&user.ID, &user.Email, &user.Perms)

// 3. Serialized directly to JSON as an integer:
// {"id": 42, "email": "admin@example.com", "perms": 7}
```

### Immutable Permissions for Read-Only Tokens

Derive scoped permissions for tokens without mutating the original user:

```go
// Create a read-only token by stripping write and delete privileges
tokenPerms := user.Perms.Without(PermissionWrite, PermissionDelete)
```

---

## Set Operations

Compare and combine permission sets using boolean set algebra:

```go
adminRole := bitmask.New(PermissionRead, PermissionWrite, PermissionDelete)
userPerms := bitmask.New(PermissionRead, PermissionWrite)

// Subset check: Does adminRole contain all userPerms?
isAdmin := userPerms.Contains(adminRole) // false

// Overlap check: Does userPerms have ANY permission from adminRole?
hasAny := userPerms.Overlaps(adminRole) // true

// Set intersection, union, and difference
commonPerms  := userPerms.And(adminRole)    // Intersection (AND)
allPerms     := userPerms.Or(adminRole)     // Union (OR)
diffPerms    := adminRole.AndNot(userPerms) // Difference (AND NOT)
```

---

## Enum Mapping

If your custom enum type implements `fmt.Stringer` (e.g. via `go generate` or custom `.String()`), you can parse and format flags **without defining manual maps**:

```go
func (p Permission) String() string {
	switch p {
	case PermissionLogin:
		return "login"
	case PermissionRead:
		return "read"
	case PermissionWrite:
		return "write"
	case PermissionDelete:
		return "delete"
	default:
		return strconv.Itoa(int(p))
	}
}

// 1. Parse CLI arguments directly:
// $ cli user add-role --perms read,write
mask, err := bitmask.ParseAuto[Permission]("read", "write")

// 2. Format to human-readable strings:
fmt.Println(mask.Strings())          // ["read", "write"]
fmt.Println(mask.FormatAuto(", "))   // "read, write"
```

---

## Go 1.23+ Iterators

Iterate through active flags with zero memory allocations and early exit:

```go
// Iterate over set atomic flags (powers of 2)
for flag := range perms.All() {
	fmt.Println("Active flag:", flag)
}

// Iterate over bit position indices (0..63)
for idx := range perms.BitIndices() {
	fmt.Println("Bit index:", idx)
}
```

---

## Benchmarks

Run on Intel(R) Core(TM) i5-3470 CPU @ 3.20GHz (Windows amd64):

| Benchmark | Operations | Speed | Memory | Allocs |
| :--- | :--- | :--- | :--- | :--- |
| `BenchmarkMask_Has` | 1,000,000,000 | **0.30 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_Overlaps` | 1,000,000,000 | **0.30 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_Count` | 1,000,000,000 | **0.59 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_Contains` | 1,000,000,000 | **0.90 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_Set` | 485,500,918 | **2.64 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_Toggle` | 562,390,948 | **2.42 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_BitIndices` | 424,501,369 | **2.39 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_All` | 435,986,151 | **2.56 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_Decompose` | 149,277,348 | **8.01 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_Binary` | 72,363,501 | **14.48 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkMask_Strings` | 13,390,375 | **90.57 ns/op** | `34 B/op` | `2 allocs/op` |
| `BenchmarkMask_ParseAuto` | 4,676,766 | **252.70 ns/op** | `3 B/op` | `1 allocs/op` |

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on code style, running tests, and benchmark requirements.

---

## Security

For security vulnerabilities and responsible disclosure guidelines, please see [SECURITY.md](SECURITY.md).

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
