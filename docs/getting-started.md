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
You will be prompted to enter a human-readable **App Title** (e.g. "My Awesome App") and choose a frontend framework:
1. **Vite + React (TypeScript)** - [Recommended]
2. **Vite + React (JavaScript)**
3. **Next.js (Static Export)**
4. **Vanilla HTML/CSS/JS**

The CLI will generate all necessary files, initialize a Go module, and run `npm install` automatically.

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
