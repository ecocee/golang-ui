<div align="center">
  <img src="https://ecocee.in/logo.png" alt="Glyra Logo" width="200" height="auto" />
  <h1>Glyra</h1>
  <p><strong>Build smaller, faster, and more secure desktop applications with Go.</strong></p>
  
  [![Developed by Ecocee](https://img.shields.io/badge/Developed_by-Ecocee_Team-blue?style=for-the-badge)](https://ecocee.in)
  [![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
  [![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
</div>

---

**Glyra** is a next-generation desktop application framework developed by the [Ecocee Team](https://ecocee.in). It allows you to build stunning, high-performance cross-platform desktop applications using web technologies (HTML, CSS, JS, React, Vite) for the frontend, and the raw speed of Go for the backend.

No heavy Node.js environment. No Chromium bloat. Just pure native performance.

## 🚀 Why Glyra?

When comparing Glyra to other popular frameworks, the difference in performance and resource usage is staggering:

| Feature | Glyra | Electron | Tauri |
| :--- | :--- | :--- | :--- |
| **Backend Language** | **Go** (Blazing Fast) | Node.js (Heavy) | Rust (Steep Learning Curve) |
| **Binary Size** | **~5 MB** | ~150 MB | ~10 MB |
| **Memory Usage** | **~20 MB** | ~200+ MB | ~30 MB |
| **Boot Time** | **< 0.1s** | ~2.5s | < 0.5s |
| **Frontend** | Any (React, Vue, Vite, HTML) | Any | Any |

### Key Benefits
- **Zero Bottlenecks**: State and heavy logic run in Go, achieving maximum CPU performance without blocking the UI thread.
- **Vite & React Ready**: Seamlessly integrates with modern JS build tools.
- **Native OS Rendering**: Uses WebKit on macOS, Edge WebView2 on Windows, and WebKitGTK on Linux.
- **Single Command Setup**: Initialize a complete full-stack desktop app in seconds.

## 📦 Getting Started

We provide a zero-configuration CLI to scaffold your new Glyra project instantly.

### 1. Install the Glyra CLI
```bash
go install github.com/ecocee/golang-ui/cmd/glyra@latest
```

### 2. Create a New Project
```bash
glyra init my-app
```
*The CLI instantly scaffolds a full Vite + React environment perfectly wired to a lightweight Go webview backend.*

### 3. Run the App
Start the blazing fast Vite frontend:
```bash
cd my-app/frontend
npm install
npm run dev
```

In a new terminal tab, start the Go backend:
```bash
cd my-app
go run main.go
```

## 🏗️ Architecture

Glyra strictly enforces a clean separation of concerns:

- **The Frontend (UI Layer)**: Lives in `frontend/`. This is completely decoupled from the OS. It only handles rendering and animations.
- **The Backend (Core Logic)**: Lives in `internal/`. Go manages the application state, file system access, and network requests.
- **The Bridge**: Go functions are bound to the `window.ecocee` object, allowing secure, asynchronous communication between JS and Go.

## 📖 Documentation

Check out the `/docs` folder for detailed guides on:
- Passing data between Go and JS
- Using embedded filesystems (`//go:embed`)
- Securing your application
- Customizing the shadcn/ui design tokens

---
<div align="center">
  <sub>Built with ❤️ by the <a href="https://ecocee.in">Ecocee Team</a></sub>
</div>
