<div align="center">
  <img src="https://ecocee.in/logo.png" alt="Glyra Logo" width="200" height="auto" />
  <h1>Glyra</h1>
  <p><strong>Build smaller, faster, and more secure desktop applications with Go.</strong></p>

  [![Developed by Ecocee](https://img.shields.io/badge/Developed_by-Ecocee_Team-blue?style=for-the-badge)](https://ecocee.in)
  [![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
  [![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
</div>

---

**Glyra** is a desktop application framework developed by the [Ecocee Team](https://ecocee.in). It lets you build cross-platform desktop apps using web technologies (HTML, CSS, JS, React, Vite) for the frontend, and the speed of Go for the backend.

No Node.js runtime to ship. No Chromium bundle. The OS's own webview, a Go binary, and your frontend embedded straight inside it.

## 🚀 Why Glyra?

| Feature | Glyra | Electron | Tauri |
| :--- | :--- | :--- | :--- |
| **Backend Language** | **Go** | Node.js | Rust |
| **Binary Size** | **~5 MB** | ~150 MB | ~10 MB |
| **Memory Usage** | **~20 MB** | ~200+ MB | ~30 MB |
| **Boot Time** | **< 0.1s** | ~2.5s | < 0.5s |
| **Frontend** | Any (React, Vue, Vite, plain HTML) | Any | Any |

## ✨ Features

- **Blazing Fast Native Desktop Apps**: Combines the performance of Go with modern web interfaces.
- **Multiple Frontends**: Supports React (TypeScript/JavaScript), Next.js, and Vanilla HTML/CSS out of the box.
- **Flutter-like Hot Reloading**: Native automatic code reloading! Editing frontend code triggers HMR, and editing `.go` files automatically restarts the native backend! (You can also press `r` + `Enter` to manually reload).
- **Native OS Packaging**: Automatically creates macOS `.app` bundles, Windows `.syso` resources (via `go-winres`), and Linux `.desktop` entries on build.

## 🚀 Getting Started

Check out the full guide in the [docs/getting-started.md](docs/getting-started.md) or dive right in:

```bash
# 1. Install the CLI
go install github.com/ecocee/golang-ui/cmd/glyra@latest

# 2. Scaffold a new project
glyra init my-app
You'll be asked to pick a frontend: **Vite + React** (recommended) or **plain HTML/CSS/JS**. Glyra scaffolds the project and wires up its Go module for you.

### 3. Develop
```bash
cd my-app
glyra dev
```
For the React template this starts the Vite dev server *and* the Go backend together, with hot reload, and stops both when you hit Ctrl+C. The vanilla template just runs the Go backend directly — there's no bundler to babysit.

### 4. Ship a production binary
```bash
glyra build
```
This runs `npm run build` (if there's a frontend build step) and then compiles a single Go binary with the built frontend embedded inside it via `go:embed`. The result is one file — copy it anywhere, no `frontend/` directory required.

## 🏗️ Project Structure

A scaffolded project looks like this:

```
my-app/
├── go.mod
├── main.go              # Go backend + webview bootstrap
└── frontend/
    ├── package.json      # React template only
    ├── vite.config.js
    ├── index.html
    └── src/
        ├── main.jsx
        ├── App.jsx
        ├── App.css
        └── index.css
```

The Go ↔ JS bridge is a single call: `w.Bind("GetSystemStatus", func() string { ... })` on the Go side becomes `window.GetSystemStatus()` in the frontend — async, and ready to extend with your own bound functions.

## 🛠️ Glyra CLI internals

The CLI itself (this repository) is organized for easy extension:

```
cmd/glyra/             entrypoint
internal/cli/          command implementations (init, dev, build, help, ...)
internal/scaffold/     the templating engine
internal/scaffold/templates/
    ├── vanilla/        plain HTML/CSS/JS starter
    └── react/          Vite + React starter
```

Starter templates live as real files under `internal/scaffold/templates/`, embedded into the CLI binary at compile time with `go:embed`. Any file ending in `.tmpl` is rendered through `text/template` (so it can use `{{.ProjectName}}` / `{{.PackageName}}`) when a project is generated; everything else is copied byte-for-byte. Adding a new starter template is just adding a new folder there — no Go string literals to edit.

## 📖 Documentation

Check out the `/docs` folder for guides on:
- Adding your own Go ↔ JS bound functions
- Customizing the webview window (size, title, menus)
- Using the design tokens in the included templates

---
<div align="center">
  <sub>Built with ❤️ by the <a href="https://ecocee.in">Ecocee Team</a></sub>
</div>
