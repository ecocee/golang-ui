# Progress Tracker

## Current Phase

Phase 1 — Foundation (in progress).

## Current Goal

Build the reactive signals runtime and design tokens system.

## Completed

- `CLAUDE.md` — project-wide engineering standards
- `context/architecture.md` — stack, boundaries, storage model, invariants
- `context/project-overview.md` — goals, features, scope, success criteria
- `context/ui-context.md` — design tokens, typography, spacing, component conventions
- `context/code-standards.md` — Go conventions, testing rules, naming, git
- `context/ai-workflow-rules.md` — development workflow, scoping rules, delivery approach
- `go mod init github.com/ecocee/golang-ui` — module initialized
- `pkg/theme/` — design tokens (light + dark palettes, spacing, radius, typography)
- `internal/signals/` — reactive signals runtime (Signal, Computed, Subscribe, Batch) with full race-detector safety

## In Progress

- None yet.

## Next Up

1. Initialize Go module (`go mod init github.com/ecocee-internal/golang-ui`)
2. Create package directory structure
3. Implement reactive signals runtime (`internal/signals/`)
   - `Signal[T]` type with `Get()` / `Set()`
   - `Computed[T]` type with dependency tracking
   - `Subscribe()` with cleanup
   - Batching, glitch prevention, concurrency safety
4. Implement design tokens (`pkg/theme/`)
5. Implement Fyne bridge (`internal/fynebridge/`)

## Open Questions

- None currently. All decisions from the planning conversation are documented in the relevant context files.

## Architecture Decisions

- **Fyne as rendering backend** — pure Go, mature, cross-platform. Isolated to `internal/fynebridge/`.
- **Custom Flexbox layout engine** — Fyne's built-in layouts are too limited. Constraint-based model. Flexbox only for v1; Grid deferred.
- **Custom reactive signals runtime** — no existing Go library fits. Full system: dependency tracking, glitch prevention, batching, memoization, concurrency safety.
- **shadcn-style declarative API** — components as pure functions of props to `fyne.Widget`. Familiar, composable.
- **Hybrid theming** — global design tokens for consistency, per-component variant/size props for escape hatches.
- **Desktop-only v1** — Windows, macOS, Linux. Mobile and web are future possibilities, not current targets.
- **No CLI in v1** — `go get` is sufficient. CLI tooling deferred until component set justifies it.

## Session Notes

- This is a greenfield project. No existing code to work around.
- Full-time development. Estimated 6-8 months to v1.0.
- The layout engine and signals runtime are the two highest-risk components. Start with signals (it's a dependency for everything else), then layout engine, then components.
- Follow the phases in order. Do not skip ahead.
