# AI Workflow Rules

## Approach

This project is built **incrementally and spec-driven**. The context files define what to build, how to build it, and the current state of progress. Always implement against these specs — do not infer or invent behavior from scratch.

Each implementation step is a small, verifiable increment. A step is complete when:
1. It compiles (`go build ./...`)
2. Tests pass (`go test -race ./...`)
3. Lint passes (`golangci-lint run ./...`)
4. The relevant context file is updated if the implementation changed scope, architecture, or conventions

## Scoping Rules

- **One feature unit at a time.** Never implement two components, two systems, or a system plus its tests in the same step unless they are trivially small.
- **Prefer small, verifiable increments.** A working signal `NewSignal` with one test is better than a half-finished signal system with ten types.
- **Do not mix system boundaries in a single step.** Signals code in the same step as layout code? No. Component code in the same step as the bridge layer? No.
- **If a change cannot be verified end-to-end in under 5 minutes, the scope is too broad.** Split it.

## When to Split Work

Split an implementation step if it combines:

- **Multiple systems** — e.g., building the signal runtime *and* the layout engine in the same step
- **Multiple components** — e.g., building Button *and* Input *and* Select in the same step
- **Implementation and tests** — Implementation first (verify it compiles), then tests (verify correctness). Exception: trivial tests that take < 30 seconds.
- **Behavior not defined in context files** — If you're inventing behavior, stop. Update the context file first.

## Handling Missing Requirements

- **Do not invent product behavior** not defined in the context files. If a component's behavior is unspecified, it doesn't exist yet.
- **If a requirement is ambiguous**, resolve it in the relevant context file before implementing. Update `ui-context.md` for visual questions, `architecture.md` for technical questions.
- **If a requirement is missing**, add it as an open question in `progress-tracker.md` before continuing. Don't guess.

## Protected Files

Do not modify the following unless explicitly instructed:

- `go.mod` / `go.sum` — dependency changes require explicit approval
- `context/*.md` — only update when implementation genuinely changes architecture, scope, or conventions
- Any file in a package you are not currently working in — no "while I'm here" refactors

## Keeping Docs in Sync

Update the relevant context file whenever implementation changes:

- `architecture.md` — system boundaries, stack, invariants, storage model
- `project-overview.md` — feature scope, success criteria, decisions
- `ui-context.md` — design tokens, component conventions, layout patterns
- `code-standards.md` — conventions, naming, testing rules
- `progress-tracker.md` — current phase, completed work, next up

## Before Moving to the Next Unit

1. The current unit **works end-to-end** within its defined scope (compiles, tests pass, lint passes)
2. **No invariant** defined in `architecture.md` was violated
3. `progress-tracker.md` reflects the completed work
4. Any open questions raised during implementation are documented

## Implementation Order

Follow the phases defined in `project-overview.md` strictly:

1. Foundation (signals → tokens → bridge)
2. Layout engine
3. Core components
4. Advanced components
5. Polish

Do not skip ahead. The layout engine depends on signals. Components depend on the layout engine. Each phase is a prerequisite for the next.

## Testing Discipline

- Write tests **immediately after** the implementation, not at the end of the phase.
- Every `internal/signals/` and `internal/layout/` function must have a corresponding test.
- Every `pkg/components/` component must have at least a render test (does it produce a valid widget?) and an interaction test (does it respond to events?).
- If a test requires > 10 lines of setup, extract it into a helper. Test code is still code.

## Code Review Checklist (Self-Before submitting any work)

- [ ] No Fyne imports outside `internal/fynebridge/`
- [ ] No global mutable state
- [ ] All errors handled or wrapped
- [ ] All goroutines have bounded lifecycle
- [ ] No hardcoded visual values (all from `pkg/theme/`)
- [ ] Race detector passes on any concurrent code
- [ ] Tests are table-driven where multiple cases exist
- [ ] Exported symbols have doc comments
