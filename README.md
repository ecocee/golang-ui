# golang-ui

A shadcn/ui-inspired component library for Go. Beautiful, composable,
design-system-driven UI built on Fyne.

## Status

Early development — see `context/progress-tracker.md` for current phase.

## Quick Start

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "github.com/ecocee/golang-ui/pkg/theme"
    "github.com/ecocee/golang-ui/pkg/ui"
)

func main() {
    a := app.New()
    w := a.NewWindow("My App")
    w.Resize(fyne.NewSize(800, 600))

    theme := theme.Default()
    _ = theme // apply to your app

    greeting := ui.Text("Hello, world!", ui.TextProps{
        Size: ui.TextSizeHeading2,
    })
    _ = greeting

    w.SetContent(container.NewVBox(greeting))
    w.ShowAndRun()
}
```

## Architecture

See `context/architecture.md` for the full system design.

## License

MIT
