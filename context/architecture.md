# Architecture Context

## Stack

| Layer | Technology | Role |
| --- | --- | --- |
| Language | Go 1.22+ | Primary language. No cgo in core library. |
| UI Framework | Fyne v2 | Rendering backend. Provides canvas, input events, window management. All Fyne widgets are wrapped, not replaced. |
| Layout Engine | Custom constraint-based Flexbox | Replaces Fyne's built-in layouts. Computes child positions/sizes from parent constraints. Integrated with the reactive signal system. |
| State Management | Custom reactive signals runtime | Fine-grained reactivity. Signals drive layout re-computation and component re-render. |
| Component Library | `pkg/components` | shadcn/ui-inspired declarative components. Stateless, composable, theme-driven. |
| Design Tokens | `pkg/theme` | Type-safe design tokens (color, spacing, radius, typography). Single source of truth for all visual values. |
| CLI Tooling (future) | `cmd/ecocee-ui` | `ecocee-ui add <component>` — scaffold components into user projects. Not in v1. |

## System Boundaries

- `cmd/` — Entrypoints only. `cmd/golang-ui` is the demo/playground app. `cmd/ecocee-ui` (future) is the CLI tool. Thin: no business logic.
- `internal/` — Library-private code. Not importable by consumers. Contains the reactive runtime internals, layout engine internals, and Fyne adapter layers.
- `internal/signals/` — Reactive signals runtime. Dependency tracking, glitch prevention, batching, subscription lifecycle, concurrency safety.
- `internal/layout/` — Constraint-based Flexbox engine. Node tree, constraint solving, size computation, position allocation.
- `internal/fynebridge/` — Adapts the layout + signal system to Fyne's `fyne.Widget` interface. This is the only package that imports `fyne/v2` directly in the core library.
- `pkg/` — Public, reusable packages. Safe for external import.
- `pkg/components/` — UI components. Each component is a pure function of props + state → `fyne.Widget`. No global state, no side effects.
- `pkg/theme/` — Design tokens and theme system. Consumers override tokens to customize appearance.
- `pkg/layout/` — Public layout primitives (`Flex`, `Stack`, `Grid` when added). Thin wrappers over `internal/layout` with signal integration.
- `api/` — API contracts for future services (theme registry, component playground backend). Empty for v1.
- `docs/` — Architecture Decision Records (ADRs), design docs, contribution guides.
- `examples/` — Runnable example applications demonstrating patterns.
- `scripts/` — Build scripts, CI helpers, code generators.

## Storage Model

This is a **client-side UI library** — there is no server, no database, and no file storage owned by the library. Clarifying what lives where:

- **Signal graph (in-memory)**: The reactive runtime maintains a directed acyclic graph of signal dependencies in memory. Lives for the lifetime of the application. Never persisted.
- **Layout tree (in-memory)**: The Flexbox engine maintains a tree of layout nodes with computed sizes/positions. Recomputed on signal changes. Never persisted.
- **Theme configuration (in-memory)**: Design tokens loaded at startup. Consumers may load from a file or embed, but the library itself does not manage persistence.
- **Fyne canvas (render target)**: Fyne manages the final pixel output. The library writes to it; does not manage its lifecycle.

If a future companion service is built (theme registry, component playground), that service will have its own storage — it is a separate product with its own architecture.

## Auth and Access Model

Not applicable. This is a library, not a service. Authentication, ownership, and access control are the responsibility of the application that consumes this library.

If a companion service is built in the future (e.g., a theme registry where users publish/share themes), that service will define its own auth model in its own architecture document.

## Invariants

1. **The core library (`pkg/`, `internal/`) never imports `fyne/v2` outside of `internal/fynebridge/`**. The layout engine, signal runtime, and component logic must remain decoupled from Fyne. This allows swapping the rendering backend in the future (e.g., to web, to a different desktop framework) without rewriting the core.

2. **Components are pure functions of props and state**. A component called with the same props and same signal values must always produce the same widget tree. No global mutable state, no time-dependent logic, no network calls inside components.

3. **Signal reads are tracked, signal writes are batched**. The reactive runtime must guarantee: (a) no glitches — observers see a consistent state even when multiple signals change; (b) no missed updates — every dependency is notified; (c) no redundant computation — derived signals cache until a dependency changes.

4. **Layout is single-pass and deterministic**. Given the same constraints, parent size, and children, the layout engine must produce the same positions and sizes. No iterative relaxation, no order-dependent results, no floating-point drift across frames.

5. **All goroutines spawned by the library have bounded lifecycle**. Every goroutine created in `internal/signals/` or `internal/layout/` must be tied to a `context.Context` and exit on cancellation. The library never leaks goroutines.

6. **No blocking on the render thread**. Signal propagation and layout computation must complete in under 1ms for trees under 1,000 nodes. Expensive work (image decoding, data transformation) is the consumer's responsibility to offload.

7. **Theme tokens are the single source of truth for visual values**. No component may hardcode a color, spacing value, font size, or radius. All visual values flow through `pkg/theme/`. Consumers override tokens; components consume them.
