# Code Standards

## General

- **Single responsibility** — Every function, type, and package does one thing. If you can name it with "And", it should be split.
- **Fix root causes** — Never paper over a bug with a workaround. Find the actual source and fix it there.
- **No speculative generality** — Don't add an interface, config option, or abstraction until you have at least two concrete use cases that need it.
- **Readability over cleverness** — Code is read 10x more than it is written. Optimize for the reader.
- **Error handling is not optional** — Every error must be handled, wrapped with context, or explicitly documented as impossible.

## Go

- **Module path**: `github.com/ecocee-internal/golang-ui`
- **Go version**: 1.22+ (use range over func, generics where they add clarity)
- **No cgo** in the core library. CGO is allowed only in platform-specific internal packages if absolutely necessary.
- **Generics**: Use generics for truly generic data structures (e.g., `Signal[T]`). Do not use generics to abstract over domain concepts — use interfaces for that.
- **Interfaces are small**: 1-3 methods. If an interface has more than 3 methods, question whether it should be split.
- **Return structs, accept interfaces**: Functions return concrete types. Functions accept interfaces (or concrete types when the set is closed).
- **Named return values**: Only when they improve documentation in `godoc`. Never for control flow.
- **Package names**: Short, lowercase, singular. `components`, not `uiComponents`. `theme`, not `theming`.
- **No `init()` functions** in the library. They make testing hard and create hidden ordering dependencies.

## Fyne Integration

- **Fyne is isolated to `internal/fynebridge/`**: No other package in the core library imports `fyne/v2`. This is an architectural invariant.
- **Wrap, don't replace**: Fyne widgets are wrapped with our component API. We don't reimplement Fyne's rendering.
- **Theme integration**: Our design tokens are mapped to Fyne's `fyne.Theme` interface via `internal/fynebridge/theme.go`. This is the bridge layer's job.
- **Window management**: The consuming app manages windows. Our library never calls `app.New()` or `window.Show()`.

## Reactive Signals

- **Signal creation**: `ui.NewSignal[T](initial)` for local state. `ui.Computed[T](compute)` for derived state.
- **Signal reads are tracked**: Reading a signal inside a `Computed` automatically creates a dependency. Reading inside a `Subscribe` registers a side-effect.
- **Signal writes are batched**: Multiple `Set()` calls within the same synchronous block are batched into a single re-computation pass.
- **No circular dependencies**: The signal graph must be a DAG. If a cycle is detected, panic immediately with a clear message — don't silently loop.
- **Subscription cleanup**: Every `Subscribe()` returns a `Unsubscribe` function. Call it when the subscriber is done. Leaking subscriptions causes memory leaks.
- **Concurrency**: Signals are safe for concurrent read/write from multiple goroutines. Use `sync.RWMutex` internally — the consumer never sees a lock.

## Layout Engine

- **Node creation**: `layout.NewNode()` with `layout.NodeConfig{MinWidth, MaxWidth, MinHeight, MaxHeight, FlexGrow, FlexShrink}`.
- **Constraints are always respected**: A node never renders outside its constraints. If content overflows, it is clipped or scrollable — never expanding the parent.
- **Layout is pure**: Given the same node tree and same constraints, `Compute()` always produces the same result. No time-based or random behavior.
- **No layout in render**: Layout is computed in a separate pass before rendering. The render pass only reads computed positions.
- **Measure text correctly**: Use Fyne's `canvas.Text` or `widget.RichText` for text measurement. Never assume a fixed character width.

## Components

- **Props pattern**: Components accept a props struct and return a `fyne.Widget`. Example: `Button(props ButtonProps) fyne.Widget`.
- **No global state**: Components do not access global variables, singletons, or package-level state. Everything is passed in.
- **Variant enum**: Components with visual variants use a `Variant` type (`Primary`, `Secondary`, `Outline`, `Ghost`, `Destructive`).
- **Size enum**: Components with size options use a `Size` type (`Small`, `Default`, `Large`).
- **Callback naming**: Event handlers are named `onClick`, `onSubmit`, `onChange` — not `ClickHandler`, `SubmitFunc`.
- **Children pattern**: Components that accept children use `...fyne.CanvasObject` variadic, matching Fyne conventions.
- **Zero-value defaults**: A zero-value props struct should produce a valid (if plain) component. No required fields unless truly necessary.

## Testing

- **Table-driven tests**: Use `[]struct{name string, input T, expected U}` for all functions with multiple cases.
- **Race detector**: All concurrent code must pass `go test -race`. This is non-negotiable.
- **No sleeps in tests**: Use channels, `sync.WaitGroup`, or `atomic` to coordinate goroutines. `time.Sleep` in tests is a bug.
- **Test names**: `TestType_Method_Scenario` — e.g., `TestSignal_Set_NotifiesSubscribers`.
- **Golden files**: Use `testdata/` for complex expected outputs (rendered widget trees, JSON responses).
- **Fakes over mocks**: Prefer hand-written fakes over mock frameworks. Fakes are simpler and catch more bugs.
- **Coverage target**: 85%+ on `internal/signals/` and `internal/layout/`. 70%+ on `pkg/components/`.

## Error Handling

- **Wrap with context**: `return fmt.Errorf("operation: %w", err)`. The wrapping message should answer "what was I trying to do when this failed?"
- **Never discard errors**: `_ := foo()` is forbidden. If you intentionally ignore an error, leave a comment explaining why it's safe.
- **Sentinel errors for known cases**: `var ErrSignalClosed = errors.New("signal closed")` for recoverable, expected failures.
- **Panic is for unrecoverable only**: `panic()` is allowed only in `main()` and in programming errors (violated invariants). Never in library code for runtime errors.

## File Organization

- `cmd/golang-ui/` — Demo/playground app entrypoint
- `cmd/ecocee-ui/` — CLI tool (future)
- `internal/signals/` — Reactive signals runtime
- `internal/layout/` — Flexbox layout engine
- `internal/fynebridge/` — Fyne adapter (only package importing fyne/v2)
- `pkg/components/` — Public component library
- `pkg/theme/` — Public design tokens and theme system
- `pkg/layout/` — Public layout primitives (Flex, Stack, Grid)
- `api/` — API contracts (future)
- `docs/` — ADRs, design docs
- `examples/` — Runnable example apps
- `scripts/` — Build/CI helpers

## Naming Conventions

- **Files**: `snake_case.go` for Go files (`signal.go`, `button.go`, `flex_layout.go`).
- **Types**: `PascalCase` for exported, `camelCase` for unexported.
- **Constants**: `PascalCase` for exported (`ColorPrimary`), `camelCase` for unexported.
- **Acronyms**: Treat as a word, not all-caps (`Url` not `URL`, `Http` not `HTTP` — following Go's convention). Exception: `ID` is acceptable as it's universally understood.
- **Packages**: Short, lowercase, no underscores (`components`, `theme`, `layout`).

## Documentation

- **Every exported symbol has a doc comment**: `// Button renders an interactive button component.`
- **README in every package**: `pkg/components/README.md`, `internal/signals/README.md` — explains purpose and when to use it.
- **Examples in doc comments**: `// Example:` blocks in godoc for non-trivial types.
- **ADRs for decisions**: Any significant architectural decision gets a `docs/adr/NNNN-title.md` entry.

## Git

- **Branch naming**: `feat/<name>`, `fix/<name>`, `refactor/<name>`, `chore/<name>`.
- **Commit messages**: Conventional Commits — `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`, `test:`.
- **Squash merge** to main. Clean, atomic commits in the main history.
- **No force push** to main. Ever.
