// Package theme provides design tokens for the golang-ui component library.
//
// Design tokens are the single source of truth for all visual values:
// colors, spacing, typography, and border radius. Components consume
// tokens exclusively — no hardcoded hex values or magic numbers.
//
// A Theme can be customized by modifying the exported fields, or by
// implementing the Theme interface for full control.
package theme

import (
	"image/color"
)

// Theme defines the visual tokens consumed by all components.
type Theme struct {
	Colors  Colors
	Spacing Spacing
	Radius  Radius
	Type    Typography
}

// Colors holds all color tokens for a theme variant.
type Colors struct {
	Background         color.Color
	Surface            color.Color
	SurfaceHover       color.Color
	Foreground         color.Color
	Muted              color.Color
	Disabled           color.Color
	Primary            color.Color
	PrimaryForeground  color.Color
	Secondary          color.Color
	Border             color.Color
	BorderStrong       color.Color
	Error              color.Color
	ErrorForeground    color.Color
	Success            color.Color
	Warning            color.Color
}

// Spacing holds the 4px-based spacing scale.
type Spacing struct {
	Space1 float32 // 4px
	Space2 float32 // 8px
	Space3 float32 // 12px
	Space4 float32 // 16px
	Space5 float32 // 24px
	Space6 float32 // 32px
	Space7 float32 // 48px
	Space8 float32 // 64px
}

// Radius holds the border radius scale.
type Radius struct {
	Small  float32 // 4px
	Medium float32 // 6px
	Large  float32 // 8px
	Full   float32 // 9999px
}

// Typography holds font size tokens.
type Typography struct {
	Heading1     float32
	Heading2     float32
	Heading3     float32
	Body         float32
	BodySmall    float32
	Caption      float32
	Button       float32
	Mono         float32
	LineHeight   float32 // multiplier, e.g. 1.5
	HeadingLine  float32 // multiplier for headings, e.g. 1.2
}
