---
name: go-review
description: Review Go code for correctness, concurrency, performance, and idiomatic design.
---

# Go Code Review

Review code using Go 1.27 conventions.

## Correctness
- Check error handling and wrapping.
- Check context cancellation and propagation.
- Identify resource leaks.
- Check nil pointer risks.

## Concurrency
- Detect goroutine leaks.
- Check channel ownership and closure.
- Identify data races.
- Verify synchronization correctness.

## Performance
- Identify unnecessary allocations.
- Review slice and map capacity.
- Avoid premature optimization.
- Recommend benchmarks when appropriate.

## Design
- Prefer small, focused interfaces.
- Follow existing package boundaries.
- Avoid unnecessary abstractions.

## Validation
- Run `go test -failfast ./...`
- Run `go vet ./...`
- Run `go test -race ./...` when relevant.

## Output
Report findings ordered by severity.
Include file paths, line numbers, and suggested fixes.