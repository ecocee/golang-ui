# Architecture Context

## Stack

| Layer | Technology | Role |
| --- | --- | --- |
| Language | Go 1.22+ | Primary backend language. Handles state, IO, logic. |
| UI Framework | `webview_go` | Glyra embedded WebView rendering backend. |
| Frontend | HTML / CSS / JS | Standard web technologies for the UI layer. |
| Component Style | shadcn/ui | Pure CSS implementations of shadcn/ui components. |
| Asset Bundling | `//go:embed` | Compiles the frontend directly into the Go binary. |

## System Boundaries

- `cmd/golang-ui/` — Entrypoint. Initializes the WebView, mounts the embedded filesystem, and blocks on `w.Run()`.
- `frontend/` — The pure visual layer. Contains NO business logic. Calls Go functions exposed via the JS window object.
- `internal/bridge/` — The glue layer. Registers Go functions into the WebView (`w.Bind()`) so Javascript can call them.
- `internal/state/` — The Single Source of Truth. Thread-safe Go structs that hold the application state.
- `internal/service/` — Heavy business logic, database access, background workers.

## Storage Model

- **State (in-memory)**: Go manages the state in `internal/state/`. State is completely lost on application exit unless explicitly persisted to disk by `internal/service/`.
- **UI DOM (in-memory)**: The WebView manages the DOM. The frontend fetches state from Go to update the DOM.
- **Assets**: Embedded in the Go binary at compile time.

## Invariants

1. **Frontend is Dumb**: The Javascript layer must not hold canonical state. It should ask Go for data and render it.
2. **Go is Thread-Safe**: Any state accessed by the `internal/bridge/` layer must be protected by a `sync.Mutex` or `sync.RWMutex`, because WebView bindings can be invoked concurrently.
3. **No External Web Server**: The WebView should ideally load assets from a data URI or a local custom scheme, preventing the need to open a local HTTP port (which can trigger firewall warnings).
4. **Component-wise CSS**: Each UI component (Button, Card, Input) must have its own isolated CSS file in `frontend/css/components/` to prevent styling conflicts.
