# Progress Tracker

## Current Phase

Phase 2 — CLI, Scaffolding, and Toolchain (Completed).

## Current Goal

Standardizing the Go-to-Web API bridge mechanism for seamless two-way communication.

## Completed

- `CLAUDE.md` — project-wide engineering standards
- `context/architecture.md` — stack, boundaries, storage model, invariants
- `context/project-overview.md` — goals, features, scope, success criteria
- `go mod init github.com/ecocee/golang-ui` — module initialized
- Fyne bridge layer abandoned — Migrated to a Tauri-style WebView architecture (`github.com/webview/webview_go`).
- **Scaffolding Engine (`internal/scaffold`)** — Integrated `text/template` engine using `//go:embed` for serving project templates from memory.
- **CLI Development (`internal/cli`)** — Interactive scaffolding (`charm.land/huh/v2`) and system checks (`go`, `npm`, `node`).
- **Frontend Templates** — Shipped out of the box with React (TS/JS), Next.js (Static Export), and Vanilla HTML/CSS/JS.
- **Custom Hot Reload Engine** — Built a Flutter-style hot reloading engine in `glyra dev`. It watches `.go` files for binary restarts and natively injects Javascript for instant DOM refreshes when editing `.html`, `.css`, or `.js`.
- **Native OS Packaging** — Implemented automated OS application icon and executable packaging in `glyra build` (macOS `.app` bundles, Windows `.syso` and `-H=windowsgui`, Linux `.desktop` entries).
- **Documentation** — Completed `docs/getting-started.md` and `docs/templates-and-icons.md`.

## In Progress

- Standardizing the Go-JS API bridge as an easy-to-use utility.
- Finalizing the `docs/` folder complete content.

## Next Up

1. **Licensing**: Add license key/management support.
2. **Component Library**: Begin scaffolding actual UI components matching shadcn-ui.

## Architecture Decisions

- **WebView as rendering backend** — Tauri-style architecture (`webview_go`).
- **Scaffolding over Library** — Transformed the project into a professional scaffolding tool (`glyra`) rather than just a Go package.
- **Custom Hot Reload** — A hybrid watcher that compiles Go temp binaries and shoots HTTP webhooks to refresh Webview DOM without closing the window.
- **Native OS Tooling** — Relies on `sips`/`iconutil` for macOS and `go-winres` for Windows to natively compile app icons.

## Session Notes

- Project rebranded to **Glyra**.
- The CLI (`glyra`) acts as the central developer tool for `init`, `dev`, and `build`.
- Hot reloading drastically improved DX by separating compiler wraps and using direct JS evaluation hooks.
