# Glyra Roadmap & Professional Plan

## Vision
Glyra aims to be the undisputed best framework for building Native Desktop Applications using Go and Web Technologies (React, Next.js, Vanilla HTML/CSS). By combining a blazing-fast Go backend with the limitless styling of modern web frontends, Glyra provides a Tauri-like experience but purely localized to Go developers.

## Current State (v0.1.0)
- **Scaffolding CLI Engine**: Interactive terminal UI to generate projects with React, Next.js, and Vanilla HTML/CSS.
- **Custom Hot Reload Engine**: Flutter-style real-time DX. Watchers automatically kill/restart Go binaries for backend changes, and inject JavaScript `window.location.reload()` hooks for instant, state-preserving DOM refreshes when editing HTML/CSS.
- **Native OS Packaging**: Automated cross-platform packaging during `glyra build` (`.app` macOS bundles with `sips` icon generation, Windows `.syso` hidden-prompt execution via `go-winres`, Linux `.desktop` entries).

---

## Upcoming Version Needs & Roadmap

### v0.2.0: The API Bridge & Developer Ergonomics (Next Immediate Goal)
Currently, Go and the WebView frontend communicate via raw `w.Bind()` and `window.FunctionName()` calls. This needs to be formalized.
- **Automated Type-Safe Bridge**: A system (potentially using code-generation or reflections) to automatically map Go structs and methods to a TypeScript client interface.
- **Standardized IPC Protocol**: Ensure that passing large JSON payloads between Go and JS is highly optimized.
- **Frontend State Sync**: Introduce a lightweight signals/store system so frontend React components can natively subscribe to Go backend state changes.

### v0.3.0: Commercial & Monetization Features (License Engine)
The user explicitly requested license key capabilities for professional desktop software distribution.
- **Cryptographic License Verification**: Implement a module `glyra/license` using ECDSA or RSA signatures to validate offline license keys.
- **Hardware Binding**: Utility to lock a license key to a specific machine's hardware ID.
- **Trial Mechanisms**: Built-in logic to enforce 14-day or 30-day trial periods with secure local storage.

### v0.4.0: The Component Library & `glyra add`
Developers love shadcn/ui. Glyra should provide pre-built, beautifully styled components that are inserted directly into the user's codebase.
- **CLI Component Generator**: Implement `glyra add button`, `glyra add card`, etc., which downloads and injects raw React/Tailwind/CSS components tailored for desktop.
- **Desktop-Optimized Behaviors**: Components should feel like native desktop apps (context menus, keyboard shortcuts, focus trapping).

### v0.5.0: Advanced Window & OS Management
`webview_go` is lightweight but basic. Advanced desktop apps need more control.
- **Multi-Window Support**: Ability to spawn child windows from Go.
- **System Tray & Menus**: Native OS menu bars and system tray icons.
- **Native Dialogs**: Open File, Save File, and Alert dialogs using native OS UI instead of HTML prompts.

---

## Strategic Action Plan

1. **Step 1 (Today)**: Design the `glyra/api` and `glyra/license` packages. Determine if we will use code-gen or generic reflection for the API Bridge.
2. **Step 2**: Build out a prototype React hook `useGlyra()` that wraps the API calls elegantly.
3. **Step 3**: Draft the cryptographic requirements for the offline license key validator.
