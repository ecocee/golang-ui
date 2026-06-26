package theme

import (
	"image/color"
)

// Default returns the light theme with default design tokens.
func Default() *Theme {
	return &Theme{
		Colors: LightColors,
		Spacing: Spacing{
			Space1: 4,
			Space2: 8,
			Space3: 12,
			Space4: 16,
			Space5: 24,
			Space6: 32,
			Space7: 48,
			Space8: 64,
		},
		Radius: Radius{
			Small:  4,
			Medium: 6,
			Large:  8,
			Full:   9999,
		},
		Type: Typography{
			Heading1:    30,
			Heading2:    24,
			Heading3:    20,
			Body:        14,
			BodySmall:   12,
			Caption:     11,
			Button:      14,
			Mono:        13,
			LineHeight:  1.5,
			HeadingLine: 1.2,
		},
	}
}

// DefaultDark returns the dark theme with default design tokens.
func DefaultDark() *Theme {
	return &Theme{
		Colors: DarkColors,
		Spacing: Spacing{
			Space1: 4,
			Space2: 8,
			Space3: 12,
			Space4: 16,
			Space5: 24,
			Space6: 32,
			Space7: 48,
			Space8: 64,
		},
		Radius: Radius{
			Small:  4,
			Medium: 6,
			Large:  8,
			Full:   9999,
		},
		Type: Typography{
			Heading1:    30,
			Heading2:    24,
			Heading3:    20,
			Body:        14,
			BodySmall:   12,
			Caption:     11,
			Button:      14,
			Mono:        13,
			LineHeight:  1.5,
			HeadingLine: 1.2,
		},
	}
}

// LightColors is the default light color palette.
var LightColors = Colors{
	Background:        hex(0xFFFFFF),
	Surface:           hex(0xF8F9FC),
	SurfaceHover:      hex(0xF1F3F7),
	Foreground:        hex(0x0F172A),
	Muted:             hex(0x64748B),
	Disabled:          hex(0x94A3B8),
	Primary:           hex(0x2563EB),
	PrimaryForeground: hex(0xFFFFFF),
	Secondary:         hex(0xF1F5F9),
	Border:            hex(0xE2E8F0),
	BorderStrong:      hex(0xCBD5E1),
	Error:             hex(0xEF4444),
	ErrorForeground:   hex(0xFFFFFF),
	Success:           hex(0x22C55E),
	Warning:           hex(0xF59E0B),
}

// DarkColors is the default dark color palette.
var DarkColors = Colors{
	Background:        hex(0x0A0A0F),
	Surface:           hex(0x111118),
	SurfaceHover:      hex(0x1A1A24),
	Foreground:        hex(0xF8FAFC),
	Muted:             hex(0x94A3B8),
	Disabled:          hex(0x475569),
	Primary:           hex(0x3B82F6),
	PrimaryForeground: hex(0xFFFFFF),
	Secondary:         hex(0x1E293B),
	Border:            hex(0x1E293B),
	BorderStrong:      hex(0x334155),
	Error:             hex(0xEF4444),
	ErrorForeground:   hex(0xFFFFFF),
	Success:           hex(0x22C55E),
	Warning:           hex(0xF59E0B),
}

func hex(v uint32) color.Color {
	return color.NRGBA{
		R: byte(v >> 16),
		G: byte(v >> 8),
		B: byte(v),
		A: 0xFF,
	}
}
