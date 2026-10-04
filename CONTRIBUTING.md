# Contributing to bitmask

Thank you for your interest in contributing to `bitmask`! We welcome bug reports, feature suggestions, and pull requests.

---

## Development Setup

### Prerequisites
- **Go 1.23+** or newer installed.
- Git.

### Clone and Verify
```bash
git clone https://github.com/zefirixhd/bitmask.git
cd bitmask
go test -v -count=1 .
```

---

## Testing and Benchmarks

We maintain a strict **Zero Heap Allocation** policy for all core bit manipulation operations.

Before submitting any code changes, ensure all tests and benchmarks pass:

```bash
# Run all unit tests
go test -v -count=1 .

# Run benchmarks and inspect memory allocations
go test -bench . -benchmem
```

---

## Guidelines for Contributions

1. **Zero-Allocation Core:** Any new methods on `Mask[T]` that perform bitwise operations or queries must not allocate memory on the heap (`0 B/op`, `0 allocs/op`).
2. **Type Safety & Generics:** Methods should work with any type that satisfies `Flag` (`~uint | ~uint64 | ...`). Avoid restricting methods to concrete types unless strictly necessary.
3. **Idiomatic Documentation:** Every exported function, type, and method must have clear GoDoc comments in English, complete with concrete code and binary arithmetic examples (`Example: // m.bits: 0b0001`).
4. **Test Coverage:** Every new function or method must be accompanied by comprehensive unit tests and benchmarks in the corresponding `*_test.go` file.
5. **Receiver Consistency:** 
   - Non-mutating methods (queries, mathematical transformations) must use **value receivers** `(m Mask[T])`.
   - Mutating methods (`Set`, `Unset`, `Toggle`, `Reset`) must use **pointer receivers** `(m *Mask[T])` and return `*Mask[T]` for chaining.

---

## Submitting a Pull Request

1. Fork the repository and create your branch from `main`:
   ```bash
   git checkout -b feature/my-cool-improvement
   ```
2. Make your changes and add tests.
3. Verify formatting:
   ```bash
   go fmt ./...
   go vet ./...
   ```
4. Commit your changes with clear, descriptive commit messages:
   ```bash
   git commit -m "feat: add leading zeros helper with benchmarks"
   ```
5. Push to your fork and submit a Pull Request against the `main` branch.

Thank you for helping make `bitmask` better!
