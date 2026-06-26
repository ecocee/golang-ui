# Project Overview

## What We Are Building

`golang-ui` is a **shadcn/ui-inspired component library for Go** — a composable, beautiful, design-system-driven UI toolkit built on Fyne. It gives Go desktop applications modern, polished UI out of the box with minimal configuration.

This is **not** a web framework, not a TUI library, not a design tool. It is a **component library** that Go developers `go get` and compose in code.

## Goals

1. **Best-in-class developer experience** — Import, compose, ship. No configuration files, no build steps, no code generation in v1.
2. **Beautiful by default** — Sensible defaults that look professional. No unstyled primitives.
3. **Fully customizable** — Design tokens control the entire visual language. Override anything.
4. **Race-free** — Zero data races. The reactive runtime and layout engine are concurrency-safe by design.
5. **Performant** — 60fps rendering. No blocking on the render thread. Minimal allocations in hot paths.
6. **Composable** — Components are pure functions of props to UI. Compose them freely. No inheritance, no magic.
7. **Fyne-native** — Built on Fyne. Leverages Fyne's cross-platform rendering, input handling, and window management. We wrap Fyne; we don't fight it.

## Non-Goals (Deliberately Out of Scope)

- **Web rendering** — v1 targets desktop only. Web/WASM is a future possibility, not a current target.
- **Terminal TUI** — Different rendering model. Separate product.
- **Visual drag-and-drop builder** — Not in scope. Maybe a future companion tool.
- **Server-side rendering** — This is a client-side library.
- **Authentication, networking, database** — These are application concerns, not library concerns.
- **Native widget rendering** — We use Fyne's canvas rendering, not OS-native widgets. This is a trade-off for cross-platform consistency.
- **Mobile platforms** — Fyne supports mobile but v1 targets desktop (Windows, macOS, Linux).

## Target User

Go developers building desktop applications who are dissatisfied with the current state of Go UI (Fyne's default look, limited component ecosystem, inconsistent design). They want modern UI without leaving Go or reaching for web tech.

## Core User Flow

The user of this library is a **developer**, not an end-user. Their flow:

1. Developer runs `go get github.com/ecocee-internal/golang-ui`
2. Imports `pkg/components` and `pkg/theme` in their Fyne app
3. Creates a theme (or uses the default)
4. Composes UI by calling component functions: `ui.Button(...)`, `ui.Card(...)`, `ui.Flex(...)`
5. Binds signals to interactive components for reactive updates
6. Runs their app — it looks polished out of the box

## Core Features

### Phase 1 — Foundation (v0.1)
- Reactive signals runtime (dependency tracking, glitch prevention, batching, memoization, concurrency safety)
- Design tokens system (color palette, spacing scale, border radius, typography)
- Fyne layout wrappers (temporary ergonomic layer before custom engine)

### Phase 2 — Layout Engine (v0.2)
- Constraint-based Flexbox implementation
- Signal integration (layout recomputes on state change)
- Layout transition hooks (animate size/position changes)

### Phase 3 — Core Components (v0.3)
- **Primitives**: Text, Button, Input, Container, Image, Icon
- **Layout**: Flex, Stack, Card, Separator, Divider
- **Form**: TextField, TextArea, Checkbox, Radio, Select, Toggle, Slider

### Phase 4 — Advanced Components (v0.5)
- **Data**: Table, List, Tree, Avatar, Badge, Tag
- **Overlay**: Modal, Dialog, Tooltip, Popover, Dropdown, Toast, ContextMenu
- **Navigation**: Tabs, Sidebar, Breadcrumb, Stepper
- **Utility**: Skeleton, Spinner, Progress, Accordion, Collapsible

### Phase 5 — Polish (v1.0)
- Comprehensive documentation with live examples
- Interactive playground application
- Theme editor (visual token customization)
- Accessibility audit and keyboard navigation
- Performance benchmarks and optimization pass

## Success Criteria

1. A new user can compose a complete form (text fields, selects, buttons, validation feedback) in under 10 minutes from `go get` to running UI.
2. 60fps sustained with 100+ visible components on screen.
3. Zero data races detected by `go test -race ./...` across the entire codebase.
4. 85%+ test coverage on `internal/signals/` and `internal/layout/`.
5. 100+ GitHub stars within 3 months of v1.0 release.

## Key Decisions & Rationale

| Decision | Rationale |
| --- | --- |
| Fyne as rendering backend | Pure Go, mature, cross-platform, active community. Avoids cgo complexity. |
| Custom layout engine (Flexbox) | Fyne's built-in layouts are too limited for modern UI. Flexbox is the most widely understood layout model. |
| Custom reactive signals | No existing Go reactive library fits the requirements. Building our own gives full control over performance and integration. |
| shadcn-style declarative API | Proven pattern. Familiar to developers coming from React/Vue. Composable. |
| Hybrid theming (tokens + per-component props) | Tokens ensure consistency. Per-component props provide escape hatches. Best of both worlds. |
| No CLI in v1 | Keep it simple. `go get` is enough. Add CLI tooling when the component set justifies it. |
| Desktop-only for v1 | Desktop is where Go's UI story is weakest. Mobile/web are different rendering models — tackle them separately. |

## Risks & Mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| Layout engine complexity overruns timeline | High | Ship Flexbox only. Defer Grid. Use real usage data to guide v2. |
| Reactive signal runtime has subtle concurrency bugs | High | Extensive race-detector testing. Property-based tests. Start with minimal core. |
| Fyne API changes break the bridge layer | Medium | Isolate Fyne imports to `internal/fynebridge/`. Version-pin Fyne dependency. |
| Component library feels incomplete | Medium | Prioritize the 20% of components that cover 80% of use cases. Ship primitives first. |
| Performance doesn't hit 60fps target | Medium | Benchmark early. Profile hot paths. Use object pooling for layout nodes. |
| Low adoption | Medium | Ship fast, document well, engage Go community early. |
