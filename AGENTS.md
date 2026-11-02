# AGENTS.md

## Project Overview

This repository contains a Go application.

Follow existing architecture, naming conventions, and established patterns. Prefer minimal, focused changes over broad refactoring.

## Go Conventions

- Target the Go version specified in `go.mod`.
- Follow idiomatic Go conventions and [Effective Go](https://go.dev/doc/effective_go).
- Format code using `gofmt` or `goimports`.
- Use `go vet` and existing project linters.
- Prefer standard library packages over third-party dependencies.
- Prefer simple, readable implementations over clever abstractions.
- Avoid unnecessary allocations and premature optimization.
- Use `any` instead of `interface{}` where appropriate.
- Prefer early returns to deeply nested conditionals.
- Avoid named return values unless they improve clarity.
- Document exported types and functions according to Go conventions.
- Avoid package-level mutable state unless necessary.

## Architecture

- Respect existing package boundaries and dependency direction.
- Keep business logic independent of transport and infrastructure concerns.
- Separate HTTP handlers, application services, domain logic, and persistence where the existing architecture supports it.
- Avoid circular dependencies.
- Define interfaces at the consumer boundary when abstraction is needed.
- Prefer concrete types until multiple implementations or testability justify an interface.
- Use dependency injection through constructors rather than global state.
- Keep interfaces small and focused.
- Do not introduce new architectural layers without a clear benefit.
- Avoid modifying public APIs unless explicitly requested.

## Dependencies

- Prefer the Go standard library whenever practical.
- Do not add dependencies without explaining their purpose.
- Reuse existing dependencies before introducing alternatives.
- Avoid unnecessary transitive dependencies.
- Do not upgrade unrelated modules.
- Keep `go.mod` and `go.sum` consistent.
- Never use `replace` directives as a permanent workaround without approval.
- Avoid deprecated or unmaintained libraries.

## Error Handling

- Never silently ignore errors.
- Return errors to the caller when the caller can handle them.
- Wrap errors with meaningful context using `fmt.Errorf("operation: %w", err)`.
- Use `errors.Is` and `errors.As` for error inspection.
- Prefer sentinel errors only when callers need to distinguish specific conditions.
- Do not use panic for ordinary error handling.
- Do not log and return the same error at every layer.
- Log errors at appropriate application boundaries.
- Never expose secrets or sensitive information in errors or logs.
- Preserve cancellation and deadline errors when propagating them.

## Context and Concurrency

- Pass `context.Context` as the first parameter when appropriate.
- Do not store request-scoped contexts in structs.
- Respect context cancellation and deadlines.
- Avoid starting goroutines without a clear lifecycle.
- Ensure goroutines terminate during shutdown or cancellation.
- Protect shared mutable state using appropriate synchronization.
- Avoid holding mutexes during blocking I/O.
- Never close channels owned by another component.
- Prefer straightforward synchronous code unless concurrency provides a measurable benefit.

## Testing

- Add or update tests when changing behavior.
- Prefer table-driven tests when they improve readability.
- Test error paths and edge cases.
- Use the standard `testing` package by default.
- Avoid unnecessary mocking frameworks.
- Prefer testing observable behavior over implementation details.
- Keep tests deterministic and independent.
- Do not depend on external services in unit tests.
- Run relevant tests after changes.

## Performance

- Prioritize correctness and maintainability.
- Avoid unnecessary allocations in performance-sensitive paths.
- Preallocate slices and maps when capacity is reasonably known.
- Avoid excessive string conversions and copying.
- Use benchmarks before making performance claims.
- Use `pprof` or profiling data to identify bottlenecks.
- Do not sacrifice readability for unmeasured micro-optimizations.

## Security

- Validate untrusted input at system boundaries.
- Use parameterized database queries.
- Never hardcode credentials, API keys, or tokens.
- Avoid logging sensitive information.
- Apply timeouts to outbound network operations.
- Use cryptographically secure randomness for security-sensitive operations.
- Do not disable security checks to make tests pass.

## Code Changes

When modifying code:

1. Inspect relevant packages and existing patterns.
2. Identify the smallest change that satisfies the requirement.
3. Preserve backward compatibility unless explicitly instructed otherwise.
4. Update relevant tests.
5. Format changed Go files.
6. Run appropriate validation.
7. Summarize changes and any remaining risks.

Do not:
- Refactor unrelated code.
- Introduce speculative abstractions.
- Add dependencies without justification.
- Change exported APIs unnecessarily.
- Suppress compiler or linter warnings without explanation.
- Claim tests passed unless they were actually executed.

## Validation

Use the repository's existing Makefile, task runner, or CI commands when available.

Otherwise, run:

```bash
go test -failfast ./...
go vet ./...
```

For concurrency-related changes, additionally run:

```bash
go test -race ./...
```

If validation cannot be completed, explain why.

## Communication

- Keep explanations concise and technical.
- Explain non-obvious design decisions.
- Identify compatibility implications.
- Report test results accurately.
- Distinguish confirmed issues from potential risks.
- When proposing alternatives, explain their tradeoffs.
