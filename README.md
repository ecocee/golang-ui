# golang-ui

A shadcn/ui-inspired component library for Go. Beautiful, composable,
design-system-driven UI for desktop applications.

> **Status:** Early development — reactive signals runtime and design
> tokens are complete. Fyne bridge layer and component library coming soon.

## Features

- **Reactive signals** — `Signal[T]`, `Computed[T]`, `Subscribe`, `Batch`
  with dependency tracking, glitch prevention, and concurrency safety
- **Design tokens** — Light/dark palettes, 4px spacing grid, border radius
  scale, typography system
- **Race-free** — All concurrent code passes `go test -race`
- **Fyne-native** — Built on Fyne v2 for cross-platform desktop rendering

## Quick Start

```bash
go get github.com/ecocee/golang-ui
```

```go
package main

import (
    "fmt"
    "github.com/ecocee/golang-ui/internal/signals"
)

func main() {
    count := signals.NewSignal(0)
    double := signals.NewComputed(func() int { return count.Get() * 2 })

    signals.Subscribe(func() {
        fmt.Printf("count=%d, double=%d\n", count.Get(), double.Get())
    })

    count.set(5) // prints: count=5, double=10
}
```

## Examples

```bash
go run ./examples/signals    # Reactive primitives
go run ./examples/counter    # Counter with derived values
go run ./examples/forms      # Reactive form validation
go run ./examples/shopping   # Shopping cart with discounts
go run ./examples/dashboard  # Dashboard with metrics
go run ./examples/theme      # Design tokens
```

## Architecture

See [`context/architecture.md`](context/architecture.md) for the full
system design, package layout, and invariants.

## Development

```bash
go build ./...              # Build all packages
go test -race ./...         # Run tests with race detector
golangci-lint run ./...     # Lint
```

## License

MIT
