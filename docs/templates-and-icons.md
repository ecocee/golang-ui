# Icons and OS Packaging

Glyra natively packages your application based on your Operating System to provide a true "native app" feel when you run `glyra build`.

## Adding an App Icon
To give your application a custom logo:
1. Create an `assets/` directory at the root of your project.
2. Place a high-resolution PNG image inside named exactly `appicon.png` (e.g., `assets/appicon.png`).

When you run `glyra build`, the CLI automatically handles the complex OS-level icon injection!

### macOS (.app Bundles)
If you build on a Mac, Glyra intercepts the binary and creates a beautiful native `YourApp.app` bundle.
- It uses Apple's built-in `sips` and `iconutil` CLI tools to dynamically compile your `appicon.png` into a multi-layered `.icns` Apple Icon Image.
- It automatically generates an `Info.plist` and wires your icon so it appears seamlessly on your Dock, Finder, and Launchpad.

### Windows (.exe Resources)
If you build on Windows, Glyra automatically generates a `winres.json` file.
- It downloads and executes `go-winres` on the fly to convert your `appicon.png` into a `.syso` Windows resource file.
- It injects the `-H=windowsgui` flag so your `.exe` runs silently in the background without launching a black command prompt, and displays your icon cleanly on the Taskbar.

### Linux (.desktop Entries)
On Linux, Glyra will automatically output a `.desktop` shortcut file mapped to your binary and wire the `Icon=` property directly to your PNG asset.

---

# Customizing Templates

Glyra's scaffolding engine uses Go's `embed.FS` to serve project templates directly from memory.

If you wish to fork Glyra or add your own template architectures to the `glyra init` command:
1. Navigate to the `internal/scaffold/templates/` folder in the Glyra CLI source code.
2. Create a new directory for your framework (e.g., `svelte/`).
3. Add your base files (`main.go.tmpl`, `package.json`, etc.).
   - Files ending in `.tmpl` will be parsed by the Go `text/template` engine. This is useful for injecting the `{{.ProjectName}}`, `{{.PackageName}}`, and `{{.Title}}` variables.
   - Files *without* `.tmpl` (like React's `.tsx` files) bypass the parser and are copied byte-for-byte. (This prevents Go from crashing when it sees standard React `{{ style }}` double-braces!)
4. Update the `huh.NewOption` dropdown in `internal/cli/init.go`.
5. Update the switch case routing in `internal/cli/init.go` and add the template string to `internal/scaffold/scaffold.go`.
6. Run `go install ./cmd/glyra` to embed the new template into your local binary!
