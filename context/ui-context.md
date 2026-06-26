# UI Context

## Theme

The visual language is a **modern, clean design system** inspired by shadcn/ui and contemporary desktop tooling (Linear, Raycast, Vercel Dashboard). It is **light and dark** from day one — not dark-only.

Design characteristics:
- **Near-white / near-black base** — clean backgrounds, no muddy grays
- **Layered surfaces** — cards and panels use a slightly elevated background to create depth
- **Muted borders** — subtle, not heavy. Borders define structure without shouting
- **Vibrant accents** — interactive elements (buttons, links, focuses) use a clear primary color
- **Generous whitespace** — breathing room between elements. Cramped UI is the default failure mode in Go GUI apps.
- **System font stack** — uses the OS native font on each platform via Fyne's `theme.DefaultTextFont()`. No bundled fonts.

## Colors

All colors are defined as design tokens in `pkg/theme`. Components consume tokens exclusively — no hardcoded hex values anywhere.

### Light Mode

| Role | Token | Value | Usage |
| --- | --- | --- | --- |
| Page background | `ColorBackground` | `#FFFFFF` | Main window background |
| Surface (cards, panels) | `ColorSurface` | `#F8F9FC` | Elevated containers |
| Surface hover | `ColorSurfaceHover` | `#F1F3F7` | Hover state for interactive surfaces |
| Primary text | `ColorForeground` | `#0F172A` | Headings, body text |
| Muted text | `ColorMuted` | `#64748B` | Placeholders, descriptions, secondary info |
| Disabled text | `ColorDisabled` | `#94A3B8` | Disabled elements |
| Primary accent | `ColorPrimary` | `#2563EB` | Buttons, links, active states, focus rings |
| Primary foreground | `ColorPrimaryForeground` | `#FFFFFF` | Text on primary accent background |
| Secondary accent | `ColorSecondary` | `#F1F5F9` | Secondary buttons, tags, subtle highlights |
| Border | `ColorBorder` | `#E2E8F0` | Card borders, separators, input outlines |
| Border strong | `ColorBorderStrong` | `#CBD5E1` | Focused inputs, active borders |
| Error | `ColorError` | `#EF4444` | Validation errors, destructive actions |
| Error foreground | `ColorErrorForeground` | `#FFFFFF` | Text on error background |
| Success | `ColorSuccess` | `#22C55E` | Success states, confirmations |
| Warning | `ColorWarning` | `#F59E0B` | Warning states, caution |

### Dark Mode

| Role | Token | Value | Usage |
| --- | --- | --- | --- |
| Page background | `ColorBackground` | `#0A0A0F` | Main window background |
| Surface (cards, panels) | `ColorSurface` | `#111118` | Elevated containers |
| Surface hover | `ColorSurfaceHover` | `#1A1A24` | Hover state for interactive surfaces |
| Primary text | `ColorForeground` | `#F8FAFC` | Headings, body text |
| Muted text | `ColorMuted` | `#94A3B8` | Placeholders, descriptions, secondary info |
| Disabled text | `ColorDisabled` | `#475569` | Disabled elements |
| Primary accent | `ColorPrimary` | `#3B82F6` | Buttons, links, active states, focus rings |
| Primary foreground | `ColorPrimaryForeground` | `#FFFFFF` | Text on primary accent background |
| Secondary accent | `ColorSecondary` | `#1E293B` | Secondary buttons, tags, subtle highlights |
| Border | `ColorBorder` | `#1E293B` | Card borders, separators, input outlines |
| Border strong | `ColorBorderStrong` | `#334155` | Focused inputs, active borders |
| Error | `ColorError` | `#EF4444` | Validation errors, destructive actions |
| Error foreground | `ColorErrorForeground` | `#FFFFFF` | Text on error background |
| Success | `ColorSuccess` | `#22C55E` | Success states, confirmations |
| Warning | `ColorWarning` | `#F59E0B` | Warning states, caution |

## Typography

| Role | Token | Fyne Mapping | Size | Weight |
| --- | --- | --- | --- | --- |
| Heading H1 | `FontSizeHeading1` | `theme.TextHeadingSize()` | 30px | Bold |
| Heading H2 | `FontSizeHeading2` | — | 24px | Semi-bold |
| Heading H3 | `FontSizeHeading3` | — | 20px | Semi-bold |
| Body | `FontSizeBody` | `theme.TextSize()` | 14px | Regular |
| Body small | `FontSizeBodySmall` | — | 12px | Regular |
| Caption / label | `FontSizeCaption` | — | 11px | Medium |
| Button | `FontSizeButton` | — | 14px | Medium |
| Mono / code | `FontSizeMono` | — | 13px | Regular |

Font stack:
- **Sans (UI)**: System default via Fyne. No bundled font. Renders as SF Pro on macOS, Segoe UI on Windows, Ubuntu on Linux.
- **Mono**: System monospace via Fyne. Renders as SF Mono on macOS, Consolas on Windows, Ubuntu Mono on Linux.

Line height is 1.5x for body text, 1.2x for headings.

## Spacing & Sizing

A 4px base grid. All spacing and sizing values are multiples of 4.

| Token | Value | Usage |
| --- | --- | --- |
| `Space1` | 4px | Tight inline gaps |
| `Space2` | 8px | Icon-to-text gaps, compact padding |
| `Space3` | 12px | Input internal padding, list item gaps |
| `Space4` | 16px | Card padding, section gaps |
| `Space5` | 24px | Panel padding, large section gaps |
| `Space6` | 32px | Page-level section spacing |
| `Space7` | 48px | Major section breaks |
| `Space8` | 64px | Page margins |

| Token | Value | Usage |
| --- | --- | --- |
| `SizeButtonHeight` | 36px | Default button height |
| `SizeInputHeight` | 36px | Default input height |
| `SizeIcon` | 16px | Inline icon size |
| `SizeIconLarge` | 24px | Button/app icon size |
| `SizeAvatar` | 32px | Default avatar |
| `SizeSidebar` | 240px | Default sidebar width |

## Border Radius

| Token | Value | Usage |
| --- | --- | --- |
| `RadiusSmall` | 4px | Badges, tags, inline elements |
| `RadiusMedium` | 6px | Buttons, inputs, small cards |
| `RadiusLarge` | 8px | Cards, panels, modals |
| `RadiusFull` | 9999px | Avatars, pills, circular buttons |

## Shadows

Fyne has limited shadow support. Shadows are used sparingly and only where Fyne supports them natively (popovers, elevated surfaces via background color contrast).

| Token | Usage |
| --- | --- |
| `ShadowSm` | Subtle card elevation (border-only on Fyne) |
| `ShadowMd` | Popover/dropdown elevation |
| `ShadowLg` | Modal/dialog elevation |

## Component Conventions

### Interactive states
Every interactive component must implement these states:
- **Default** — resting state
- **Hover** — mouse over (desktop)
- **Active/Pressed** — mouse down or touch active
- **Focus** — keyboard focus, visible focus ring using `ColorPrimary` at 50% opacity
- **Disabled** — non-interactive, reduced opacity (0.5)

### Variants
Components that have visual variants use a `Variant` type:
- `VariantPrimary` — main action
- `VariantSecondary` — alternative action
- `VariantOutline` — bordered, transparent background
- `VariantGhost` — no background, hover reveals
- `VariantDestructive` — error/danger action, uses `ColorError`

### Sizes
Components that have size options use a `Size` type:
- `SizeSmall` — compact (28px height)
- `SizeDefault` — standard (36px height)
- `SizeLarge` — prominent (44px height)

## Layout Patterns

- **App shell**: Full viewport with a fixed left sidebar (`SizeSidebar` width), top bar (48px), and scrollable content area.
- **Dashboard**: Grid of stat cards (4 columns on desktop, responsive collapse), followed by content sections.
- **Settings**: Two-column layout — left nav list, right form panel.
- **Form**: Vertical stack of form groups, each containing label + input + help text. Submit actions right-aligned at bottom.
- **Modal**: Centered overlay with backdrop dimming. Close on Escape key and backdrop click.
- **Data table**: Full-width header with column rows. Hover highlights rows. Sticky header on scroll.

## Icons

- **Set**: Lucide icons (stroke-based, 24x24 viewBox, 1.5px stroke weight).
- **Sizing**: `SizeIcon` (16px) for inline, `SizeIconLarge` (24px) for buttons and standalone usage.
- **Color**: Icons inherit `ColorMuted` by default, `ColorForeground` on hover/active.
- **Integration**: Icons are rendered as `fyne.NewIcon()` from embedded SVG resources. No external dependency.
