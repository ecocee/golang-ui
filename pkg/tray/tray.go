package tray

import (
	"github.com/getlantern/systray"
)

// App represents a system tray application.
type App struct {
	onReady func()
	onExit  func()
}

// New creates a new System Tray application.
// Note: On macOS, this takes over the main thread. If you need a WebView and a System Tray
// simultaneously on macOS, you must manage CGO thread dispatching carefully.
func New() *App {
	return &App{}
}

// Run starts the system tray event loop. It blocks until Quit() is called.
func (a *App) Run(onReady func(), onExit func()) {
	a.onReady = onReady
	a.onExit = onExit
	systray.Run(a.onReady, a.onExit)
}

// Quit terminates the system tray app.
func (a *App) Quit() {
	systray.Quit()
}

// SetTitle sets the text shown next to the tray icon.
func SetTitle(title string) {
	systray.SetTitle(title)
}

// SetTooltip sets the tooltip shown when hovering over the tray icon.
func SetTooltip(tooltip string) {
	systray.SetTooltip(tooltip)
}

// AddMenuItem adds a clickable item to the tray menu.
// It returns a channel that emits a signal whenever the item is clicked.
func AddMenuItem(title string, tooltip string) <-chan struct{} {
	item := systray.AddMenuItem(title, tooltip)
	return item.ClickedCh
}

// AddSeparator adds a horizontal line to the tray menu.
func AddSeparator() {
	systray.AddSeparator()
}
