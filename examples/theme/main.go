// Command theme demonstrates the design tokens system: light/dark
// palettes, spacing, radius, and typography.
//
// Run: go run ./examples/theme
package main

import (
	"fmt"

	"github.com/ecocee/golang-ui/pkg/theme"
)

func main() {
	fmt.Println("=== Design Tokens ===")

	light := theme.Default()
	dark := theme.DefaultDark()

	fmt.Println("--- Light Theme Colors ---")
	fmt.Printf("Background:      %v\n", light.Colors.Background)
	fmt.Printf("Surface:         %v\n", light.Colors.Surface)
	fmt.Printf("Foreground:      %v\n", light.Colors.Foreground)
	fmt.Printf("Muted:           %v\n", light.Colors.Muted)
	fmt.Printf("Primary:         %v\n", light.Colors.Primary)
	fmt.Printf("Primary FG:      %v\n", light.Colors.PrimaryForeground)
	fmt.Printf("Secondary:       %v\n", light.Colors.Secondary)
	fmt.Printf("Border:          %v\n", light.Colors.Border)
	fmt.Printf("Error:           %v\n", light.Colors.Error)
	fmt.Printf("Success:         %v\n", light.Colors.Success)
	fmt.Printf("Warning:         %v\n", light.Colors.Warning)

	fmt.Println("\n--- Dark Theme Colors ---")
	fmt.Printf("Background:      %v\n", dark.Colors.Background)
	fmt.Printf("Surface:         %v\n", dark.Colors.Surface)
	fmt.Printf("Foreground:      %v\n", dark.Colors.Foreground)
	fmt.Printf("Primary:         %v\n", dark.Colors.Primary)

	fmt.Println("\n--- Spacing Scale (4px grid) ---")
	fmt.Printf("Space1: %.0fpx\n", light.Spacing.Space1)
	fmt.Printf("Space2: %.0fpx\n", light.Spacing.Space2)
	fmt.Printf("Space3: %.0fpx\n", light.Spacing.Space3)
	fmt.Printf("Space4: %.0fpx\n", light.Spacing.Space4)
	fmt.Printf("Space5: %.0fpx\n", light.Spacing.Space5)
	fmt.Printf("Space6: %.0fpx\n", light.Spacing.Space6)
	fmt.Printf("Space7: %.0fpx\n", light.Spacing.Space7)
	fmt.Printf("Space8: %.0fpx\n", light.Spacing.Space8)

	fmt.Println("\n--- Border Radius ---")
	fmt.Printf("Small:  %.0fpx\n", light.Radius.Small)
	fmt.Printf("Medium: %.0fpx\n", light.Radius.Medium)
	fmt.Printf("Large:  %.0fpx\n", light.Radius.Large)
	fmt.Printf("Full:   %.0fpx\n", light.Radius.Full)

	fmt.Println("\n--- Typography ---")
	fmt.Printf("H1:    %.0fpx\n", light.Type.Heading1)
	fmt.Printf("H2:    %.0fpx\n", light.Type.Heading2)
	fmt.Printf("H3:    %.0fpx\n", light.Type.Heading3)
	fmt.Printf("Body:  %.0fpx\n", light.Type.Body)
	fmt.Printf("Small: %.0fpx\n", light.Type.BodySmall)
	fmt.Printf("Cap:   %.0fpx\n", light.Type.Caption)
	fmt.Printf("Btn:   %.0fpx\n", light.Type.Button)
	fmt.Printf("Mono:  %.0fpx\n", light.Type.Mono)

	fmt.Println("\n--- Custom Theme Example ---")
	custom := &theme.Theme{
		Colors: theme.Colors{
			Background:        theme.Hex(0x0D1117),
			Surface:           theme.Hex(0x161B22),
			SurfaceHover:      theme.Hex(0x1C2333),
			Foreground:        theme.Hex(0xE6EDF3),
			Muted:             theme.Hex(0x8B949E),
			Disabled:          theme.Hex(0x6E7681),
			Primary:           theme.Hex(0x58A6FF),
			PrimaryForeground: theme.Hex(0x0D1117),
			Secondary:         theme.Hex(0x21262D),
			Border:            theme.Hex(0x30363D),
			BorderStrong:      theme.Hex(0x484F58),
			Error:             theme.Hex(0xF85149),
			Success:           theme.Hex(0x3FB950),
			Warning:           theme.Hex(0xD29922),
		},
		Spacing: light.Spacing,
		Radius:  light.Radius,
		Type:    light.Type,
	}
	fmt.Printf("Custom BG:    %v\n", custom.Colors.Background)
	fmt.Printf("Custom Primary: %v\n", custom.Colors.Primary)

	fmt.Println("\nTheme example complete.")
}
