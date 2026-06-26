# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## Application Building Context

Read the following files in order before implementing
or making any architectural decision:

1. `context/project-overview.md` — product definition,
   goals, features, and scope
2. `context/architecture.md` — system structure,
   boundaries, storage model, and invariants
3. `context/ui-context.md` — theme, colors, typography,
   and component conventions
4. `context/code-standards.md` — implementation rules
   and conventions
5. `context/ai-workflow-rules.md` — development workflow,
   scoping rules, and delivery approach
6. `context/progress-tracker.md` — current phase,
   completed work, open questions, and next steps

Update `context/progress-tracker.md` after each
meaningful implementation change.

If implementation changes the architecture, scope, or
standards documented in the context files, update the
relevant file before continuing.


## Project Overview

`golang-ui` is a Go-based UI application within the `ECOCEE-INTERNAL` monorepo. It is a greenfield project — no `go.mod` or source files exist yet. The goal is a high-performance, race-free, easily maintainable UI application with excellent UX.

This document defines the **non-negotiable engineering standards**, architecture, and workflow for the project. Treat it as the single source of truth.

---

## Build & Development Commands

```bash
# --- Module Initialization (one-time) ---
go mod init github.com/ecocee-internal/golang-ui

# --- Compilation ---
go build ./...                        # build all packages
go build -o bin/ ./cmd/golang-ui     # build binary to bin/

# --- Testing ---
go test ./...                         # run all tests
go test -race ./...                   # run all tests with race detector (MANDATORY before merge)
go test -cover ./...                  # run all tests with coverage
go test ./pkg/component -run TestName  # run a single test
go test -count=1 ./...                # disable test caching (use in CI)

# --- Linting & Static Analysis ---
golangci-lint run ./...               # primary linter (config in .golangci.yml)
go vet ./...                          # standard Go vet

# --- Dependency Management ---
go mod tidy                           # clean up dependencies
go mod verify                         # verify checksums

# --- Running locally ---
go run ./cmd/golang-ui               # run the application
```

**Pre-merge checklist:** `go build ./...` passes, `go test -race ./...` passes, `golangci-lint run ./...` passes with zero issues.

---

## Architecture

### Package Layout (Go Idiomatic)

```
golang-ui/
├── cmd/
│   └── golang-ui/          # Application entrypoint (main.go only — thin)
├── internal/               # Private application code (not importable externally)
│   ├── app/                # Application lifecycle, wiring, dependency injection
│   ├── ui/                 # UI rendering, layout, widgets
│   ├── state/              # State management (single source of truth)
│   ├── service/            # Business logic, use cases
│   └── platform/           # OS-specific abstractions
├── pkg/                    # Public, reusable libraries (importable by other projects)
│   ├── components/         # Reusable UI components
│   └── theme/              # Design tokens, colors, typography
├── api/                    # API contracts (protobuf, OpenAPI, JSON schemas)
├── configs/                # Configuration files
├── scripts/                # Build & CI scripts
├── Makefile                # Build automation
├── go.mod
├── .golangci.yml
└── CLAUDE.md
```

### Core Design Principles

1. **Separation of concerns** — UI rendering, state, and business logic are in separate packages. Never mix them.
2. **Dependency injection** — All dependencies are injected; no global mutable state.
3. **Single source of truth** — Application state lives in `internal/state/`; UI is a pure projection of state.
4. **Immutability** — Prefer value types and immutable data structures. Use pointers only when mutation is intentional and controlled.
5. **Explicit over implicit** — No magic. No hidden side effects. No implicit goroutine spawning.

### UI Rendering Pattern

The UI layer is a **declarative projection of state**. The application holds a state tree; on each state change, the UI re-renders the affected subtree.

```
State Change → State Tree Updated → Diff Calculation → Minimal UI Re-render
```

- Never mutate the UI directly from business logic.
- All UI updates happen on the main/render thread (platform-specific).
- Use channels or an event bus to communicate between service layer and UI layer.

---

## Concurrency & Safety (Zero Race Conditions)

This is a hard requirement. Violations are blocking bugs.

### Rules

1. **No naked goroutines** — Every goroutine must have a clear lifecycle, bounded scope, and a shutdown signal (via `context.Context`).
2. **No shared mutable state** — If data is accessed from multiple goroutines, protect it with:
   - `sync.RWMutex` for read-heavy workloads
   - `sync.Mutex` for write-heavy workloads
   - Channels for goroutine communication (prefer this over shared memory)
3. **Context propagation** — Every function that can block or spawn work must accept `context.Context` as its first parameter. Always select on `ctx.Done()`.
4. **Bounded channels** — Never create unbounded channels. Use buffered channels with a defined capacity, or use a semaphore pattern.
5. **Worker pools** — For parallel work, use bounded worker pools. Never spawn N goroutines for N tasks.
6. **Race detector is law** — `go test -race ./...` must pass on every PR. No exceptions.

### Anti-Patterns (Will Be Rejected in Review)

```go
// ❌ Naked goroutine with no shutdown
go processData(data)

// ❌ Unbounded channel
ch := make(chan Event)

// ❌ Global mutable state
var currentUser User

// ❌ Goroutine spawning per item
for _, item := range items {
    go handle(item)
}

// ✅ Correct patterns:
// Goroutine with context lifecycle
go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case work := <-workCh:
            handle(work)
        }
    }
}(ctx)

// Bounded channel
ch := make(chan Event, 256)

// Worker pool
for i := 0; i < numWorkers; i++ {
    go worker(ctx, workCh)
}
```

---

## UI/UX Standards

### Design Philosophy

- **Clarity over cleverness** — Every UI element must communicate its purpose immediately.
- **Consistency** — Use the design system in `pkg/theme/` exclusively. No hardcoded colors, spacing, or typography.
- **Accessibility** — All interactive elements must be keyboard-navigable. Support screen readers where the framework allows.
- **Responsive** — The UI must adapt gracefully to window resizing and different screen densities.
- **Performance** — UI must render at 60fps. No blocking operations on the render thread. Defer expensive work to background goroutines.

### Component Design

- Components are stateless by default. Accept all data via props/parameters.
- Use composition over inheritance. Go has no inheritance — embrace it.
- Every component must have a clear, single responsibility.
- Components communicate via callbacks or injected interfaces, never via direct references.

### State-Driven UI Contract

```go
// Every UI component follows this contract:
type Component interface {
    // Render returns the visual representation based on current state
    Render(state State) Node
    // Handle processes an event and returns an action (never mutates state directly)
    Handle(event Event) Action
}
```

---

## Testing Strategy

| Layer | Tool | Coverage Target | Notes |
|-------|------|-----------------|-------|
| Unit (pure logic) | `testing` | 90%+ | Table-driven tests, no mocks needed |
| Unit (with deps) | `testify/mock` | 85%+ | Mock all external dependencies |
| Integration | `testing` + testcontainers | 70%+ | Real DB/services in containers |
| UI/Component | framework-specific | Critical paths | Snapshot or visual regression |
| E2E | Playwright/rod | Key user flows | Smoke tests for critical paths |

### Test Rules

- **Table-driven tests** for all functions with multiple input/output cases.
- **No test ordering dependencies** — Tests are independent and can run in any order.
- **No sleeps in tests** — Use synchronization (channels, WaitGroups, `sync/atomic`) to coordinate goroutines.
- **Race detection** — All concurrent code must be tested with `-race`.
- **Golden files** for complex outputs (JSON, rendered UI).

### Example Test Structure

```go
func TestComponent_Render(t *testing.T) {
    tests := []struct {
        name     string
        state    State
        expected Node
    }{
        {"empty state", State{}, EmptyNode},
        {"with data", State{Items: []string{"a"}}, Node{Text: "a"}},
        {"with selection", State{Selected: 1}, Node{Selected: true}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := NewComponent()
            got := c.Render(tt.state)
            if diff := cmp.Diff(tt.expected, got); diff != "" {
                t.Errorf("unexpected render (-want +got):\n%s", diff)
            }
        })
    }
}
```

---

## Code Quality & Review Standards

### What Gets Caught in Review

- Race conditions or unsafe concurrency
- Missing context propagation
- Global mutable state
- Hardcoded values that should be in `pkg/theme/` or config
- Missing error handling (Go convention: always handle errors, never ignore)
- Functions with more than 3 responsibilities
- Packages with circular dependencies

### Error Handling Convention

```go
// Wrap errors with context at every layer
result, err := service.DoWork(ctx, input)
if err != nil {
    return fmt.Errorf("service.DoWork for %q: %w", input.ID, err)
}

// Never use _ = to discard errors
// Never use panic for control flow (only in main() or unrecoverable init failures)
```

### Logging

- Use structured logging (slog or zerolog).
- Never log sensitive data (tokens, passwords, PII).
- Use appropriate levels: `DEBUG` for development, `INFO` for normal flow, `WARN` for recoverable issues, `ERROR` for failures.

---

## Git & Branching

- **Main branch** is always deployable. Direct commits to `main` are blocked.
- **Branch naming**: `feat/<name>`, `fix/<name>`, `refactor/<name>`, `chore/<name>`
- **Commit messages**: Conventional Commits format (`feat:`, `fix:`, `refactor:`, `docs:`, `chore:`)
- **PRs require**: passing CI, race-free tests, linting, and one approval.
- **Squash merge** to `main` for a clean history.

---

## Makefile Targets

```makefile
.PHONY: build test lint race clean

build:
	go build -o bin/ ./cmd/golang-ui

test:
	go test -count=1 ./...

race:
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/
	go clean -testcache
```

---

## Project Governance

- **Architecture decisions** are documented as ADRs (Architecture Decision Records) in `docs/adr/`.
- **Breaking changes** require a migration plan and team consensus.
- **New dependencies** must be justified in the PR — avoid dependency bloat.
- **Performance budgets** are defined and enforced: startup time < X ms, frame render < 16ms.

---

## Notes for Claude

- When generating code, always follow the patterns defined above.
- Prefer simple, readable code over clever abstractions.
- If a design decision is needed, present options with trade-offs rather than guessing.
- When in doubt about concurrency safety, use the most conservative approach (mutex over channel, or channel over shared memory).
- The UI must never block. If an operation takes > 50ms, move it to a background goroutine.
