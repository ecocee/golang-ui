# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## Application Building Context

Read the following files in order before implementing
or making any architectural decision:

1. `context/project-overview.md` — product definition, goals, features, and scope
2. `context/architecture.md` — system structure (Glyra Architecture), boundaries, and invariants
3. `context/ui-context.md` — theme (CSS/HTML), colors (Zinc), and component conventions
4. `context/code-standards.md` — implementation rules and conventions
5. `context/ai-workflow-rules.md` — development workflow
6. `context/progress-tracker.md` — current phase

Update `context/progress-tracker.md` after each meaningful implementation change.

## Project Overview

`golang-ui` is a Go-based UI application within the `ECOCEE-INTERNAL` monorepo. It uses the **Glyra Architecture** with a high-performance Go backend and a WebView frontend powered by HTML, CSS (shadcn/ui style), and JS.

This document defines the **non-negotiable engineering standards**, architecture, and workflow for the project.

---

## Build & Development Commands

```bash
# --- Compilation ---
go build ./...                        # build all packages
go build -o bin/ ./cmd/golang-ui      # build binary to bin/

# --- Running locally ---
go run ./cmd/golang-ui                # run the application
```

---

## Architecture (Glyra)

### Package Layout

```
golang-ui/
├── cmd/
│   └── golang-ui/          # Application entrypoint (main.go only — thin)
├── frontend/               # Embedded Web Assets (HTML/CSS/JS)
│   ├── index.html          # Main entry UI
│   ├── css/
│   │   ├── theme.css       # Design tokens (CSS variables)
│   │   └── components/     # Component-wise CSS (button.css, card.css)
│   └── js/
│       └── app.js          # Frontend JS logic
├── internal/               # Private Go application code
│   ├── app/                # Application lifecycle, wiring, dependency injection
│   ├── bridge/             # Go <-> JS bindings for the WebView
│   ├── state/              # State management (single source of truth in Go)
│   └── service/            # Business logic, use cases
├── pkg/                    # Public Go libraries
└── docs/                   # ADRs and Documentation
```

### Core Design Principles

1. **Separation of concerns** — The UI rendering is strictly HTML/CSS. Go handles all business logic, state, and heavy computation.
2. **Component-wise UI** — Each UI component has its own CSS file (e.g. `button.css`) following shadcn/ui principles. No bloated monolithic stylesheets.
3. **Go is the Source of Truth** — Application state lives in `internal/state/`. The frontend is just a "dumb" visual layer that calls Go via `internal/bridge/`.
4. **Embedded Assets** — The entire `frontend/` folder is embedded into the Go binary using `//go:embed`.

---

## UI/UX Standards

- **CSS Variables** — All colors and spacing are defined in `theme.css` using CSS variables (`--primary`, `--background`).
- **Flexbox/Grid** — Use standard CSS Flexbox and Grid for all layouts.
- **Component Classes** — Use BEM or shadcn-style utility classes (e.g., `.btn`, `.btn-primary`). No inline styles unless dynamically calculated.

---

## Concurrency & Safety

- UI state updates happen in Go, protected by Mutexes (`sync.Mutex`).
- JS bindings in Go must be thread-safe when modifying state.
- Bounded goroutines with `context.Context`.

## Notes for Claude
- When building components, create the CSS in `frontend/css/components/` and the HTML structure.
- Always use the embedded filesystem (`embed.FS`) to serve the frontend.
- When in doubt, follow the shadcn/ui HTML/CSS structure.
