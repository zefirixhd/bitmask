# bitmask

[![Go Reference](https://pkg.go.dev/badge/github.com/zefirixhd/bitmask.svg)](https://pkg.go.dev/github.com/zefirixhd/bitmask)
[![Go Report Card](https://goreportcard.com/badge/github.com/zefirixhd/bitmask)](https://goreportcard.com/report/github.com/zefirixhd/bitmask)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Zero Allocations](https://img.shields.io/badge/allocs-0%20B%2Fop-brightgreen)](mask_test.go)

`bitmask` is a high-performance, type-safe, generic bitmask library for **Go 1.23+**. 

Designed for high-throughput backends, Role-Based Access Control (RBAC), and remote CLI management tools where performance, type safety, and zero memory allocations matter.

---

## Features

- 🐘 **Sub-nanosecond Performance:** Core operations (`Has`, `Overlaps`, `Contains`) execute in ~0.30–0.90 ns using hardware CPU instructions (`POPCNT`, `TZCNT`).
- 🐘 **Zero Heap Allocations:** All core mask operations, bit twiddling, iterators, and strict policy validations run at `0 B/op` and `0 allocs/op`.
- 🐘 **Strict Masks & Validation (`strict` subpackage):** Declarative validation rules (mutual exclusion, requirements, dependencies) that auto-validate during JSON unmarshaling and database scans.
- 🐘 **Type-Safe Generics:** Works seamlessly with custom enum types (`type Permission uint64`) via Go generics (`~uint...`).
- 🐘 **Database & JSON Ready:** First-class support for `database/sql` (`sql.Scanner` / `driver.Valuer`), `json.Marshaler`, and `encoding.TextMarshaler` (YAML / TOML / query params).
- 🐘 **Go 1.23+ Native Iterators:** Efficient iteration over set flags (`for flag := range mask.All()`) without allocating intermediate slices.
- 🐘 **CLI-Ready String Mapping:** Automatic parsing and formatting (`ParseAuto`, `Strings`) using your enum's `String()` method without manual map declarations.

---

## Installation

```bash
go get github.com/zefirixhd/bitmask
```

*Requires Go 1.23 or newer.*

---

## Terminology & Concepts

| Term | Definition | Representation / Example |
| :--- | :--- | :--- |
| **Bit** | The smallest unit of binary data: `0` (cleared / off) or `1` (set / on). | `0` or `1` |
| **Bit Index** | The zero-based position of a bit from least to most significant ($0..N-1$). | Index `3` $\to$ $2^3 = 8$ |
| **Flag** | A named power-of-two constant ($2^k$, a single set bit) representing an atomic permission or feature. | `const Read = 1 << 0`<br>`const Write = 1 << 1` |
| **Bitmask** | An integer representing a combination of flags via bitwise OR (`\|`). | `Read \| Write` (`0b0011`) |
| **Cardinality** | The total count of set bits (1s) in a mask, also known as **Hamming Weight** or **Population Count (popcount)**. | `mask.Count() == 2` |
| **Set Operations** | Mathematical set theory equivalents for bitwise operations: | |
| ↳ *Union* | Flags set in either mask. | `m.Or(other)` (bitwise `\|`) |
| ↳ *Intersection* | Flags set in both masks. | `m.And(other)` (bitwise `&`) |
| ↳ *Difference* | Flags in `m` that are not in `other`. | `m.AndNot(other)` (bitwise `&^`) |
| ↳ *Symmetric Diff* | Flags in exactly one of the masks. | `m.Xor(other)` (bitwise `^`) |
| **Rule** | A validation predicate (`func(Mask[F]) error`) enforcing an invariant. | `strict.MutuallyExclusive(...)` |
| **Policy** | A zero-allocation compile-time struct binding rules to a strict mask. | `strict.Policy[F]` |
| **Strict Mask** | A guarded mask wrapper that enforces policy rules on creation, updates, and deserialization. | `strict.Mask[F, P]` |

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

	// Set algebra
	adminRequired := bitmask.New(PermissionLogin, PermissionDelete)
	if !perms.Contains(adminRequired) {
		fmt.Println("Access denied: missing admin permissions")
	}
}
```

---

## Strict Masks & Rule Validation (`strict` subpackage)

The `strict` package provides compile-time policies and composable validation rules for masks. It enforces the **"Parse, don't validate"** principle: invalid bitmasks are rejected immediately during API deserialization (`json.Unmarshal`, `UnmarshalText`, SQL `Scan`).

### 1. Define a Policy

Policies are zero-sized structs (`struct{}`) that cost **0 bytes** of memory:

```go
package main

import (
	"github.com/zefirixhd/bitmask"
	"github.com/zefirixhd/bitmask/strict"
)

type UserPolicy struct{}

func (UserPolicy) Validate(m bitmask.Mask[Permission]) error {
	return strict.AllOf(
		// Read and Write cannot both be assigned in restricted mode:
		strict.MutuallyExclusive(PermissionRead, PermissionWrite),
		// Deleting requires login permission:
		strict.DependsOn(PermissionDelete, PermissionLogin),
		// Banned or guest users cannot be granted permissions here:
		strict.NoneOf(PermissionGuest),
	)(m)
}
```

### 2. Embed into Struct Payloads (Auto-validating JSON)

```go
type UpdateUserRequest struct {
	UserID int                                      `json:"user_id"`
	Perms  strict.Mask[Permission, UserPolicy]      `json:"perms"`
}

// In your HTTP handler:
var req UpdateUserRequest
err := json.Unmarshal(body, &req)
if err != nil {
	// If the JSON payload contains conflicting flags (e.g. Read + Write),
	// Unmarshal fails immediately with strict.ErrConflict!
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}
```

### Available Rules in `strict`

| Rule | Description |
| :--- | :--- |
| `MutuallyExclusive(v...)` | At most one flag from `v` can be present simultaneously. |
| `RequireAll(v...)` | All specified flags must be present. |
| `RequireAny(v...)` | At least one flag from `v` must be present. |
| `NoneOf(v...)` | None of the specified flags are allowed. |
| `DependsOn(flag, reqs...)` | If `flag` is set, all `reqs` must also be set ($A \implies B$). |
| `ExactlyOne(v...)` | Exactly one flag from `v` must be set (neither 0 nor >1). |
| `MaxCount(threshold)` | Total set bits (`Count()`) must not exceed `threshold`. |
| `MinCount(threshold)` | Total set bits (`Count()`) must be at least `threshold`. |
| `When(flag, rule)` | Applies `rule` only if `flag` is set in the mask. |
| `Not(rule)` | Inverts the outcome of a rule. |
| `AllOf(rules...)` | Logical AND: all rules must succeed. |
| `AnyOf(rules...)` | Logical OR: succeeds if at least one rule passes. |
| `Custom(err, predicate)` | Custom boolean predicate (`func(Mask[F]) bool`). |

---

## String Mapping & CLI Integration

```go
func (p Permission) String() string {
	switch p {
	case PermissionRead:
		return "read"
	case PermissionWrite:
		return "write"
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
| `BenchmarkStrict_Check` | 50,812,362 | **24.99 ns/op** | `0 B/op` | `0 allocs/op` |
| `BenchmarkStrict_MutuallyExclusive` | 231,548,985 | **5.44 ns/op** | `0 B/op` | `0 allocs/op` |
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
