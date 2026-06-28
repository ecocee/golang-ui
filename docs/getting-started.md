# Getting Started with Glyra

Glyra is a professional scaffolding and build tool for writing desktop applications using Go (Backend) and web technologies (Frontend).

## Installation

Install the CLI using Go:
```bash
go install github.com/ecocee/golang-ui/cmd/glyra@latest
```

## Scaffolding an Application

Create a new app by running:
```bash
glyra init my-app
```
You'll be prompted to choose a frontend framework:
1. **Vite + React (TypeScript)** - [Recommended]
2. **Vite + React (JavaScript)**
3. **Next.js (Static Export)**
4. **Svelte (Vite)**
5. **Vanilla HTML/CSS/JS**

The CLI generates the project, initializes a Go module, and runs `npm install` automatically.

### Project structure

A scaffolded project keeps the framework wiring and your logic in separate places:

```
my-app/
├── go.mod
├── go.sum
├── main.go                  # thin entrypoint — wires the window, registers services
├── src/
│   └── system.go            # your backend logic (services)
└── frontend/                # your web UI (template-specific)
    └── ...
```

**`main.go`** is intentionally thin: it creates the webview window, serves your
frontend, and registers your backend services with the API bridge. You rarely need to
edit it — it's Glyra's wiring, not your code.

**`src/`** is where your logic lives. Each file defines one or more *services* — plain
Go structs with exported methods. Register a service once in `main.go` and every method
becomes callable from the frontend as `window.<Service>_<Method>()`. To add a new
capability, add a new file under `src/` and register it. Your project scales by adding
files, not by growing a single `main.go`.

See [API Reference](api.md) for how the Go ↔ JavaScript bridge works.

## Development Mode & Hot Reloading

Run the following command inside your project directory to start the backend and frontend concurrently:
```bash
glyra dev
```
In Dev Mode, Glyra acts as a transparent proxy to `localhost:5173` (Vite) or `localhost:3000` (Next.js). This unlocks two extremely powerful development workflows:

1. **Frontend Hot Module Replacement (HMR)**: Any changes you save in your React/Next.js frontend are updated immediately in the app window without refreshing.
2. **Backend Native Hot Reloading**: If you save changes to *any* `.go` file in your project, Glyra detects the change and automatically recompiles and re-launches the backend window in milliseconds!
3. **Manual Reload**: You can manually trigger a backend restart at any time by pressing `r` + `Enter` in the terminal!

## Production Build

To compile a highly optimized, standalone native binary (stripping debugging symbols for a smaller size):
```bash
glyra build
```

During the build process, Glyra will:
1. Build your frontend into static assets (`dist/` or `out/`).
2. Compile a standalone Go binary using `//go:embed` to package your UI directly inside the executable.
3. Automatically package your app natively for the target OS (see [Templates and Icons](templates-and-icons.md)).

Your final optimized binary will be placed in your project root!
