package cli

import (
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
)

// theme holds every lipgloss style used by the CLI. Centralising them here
// keeps the visual identity consistent across commands and makes it trivial
// to swap the palette later.
var theme = struct {
	// Logo is the gradient "Glyra" wordmark (purple → pink → cyan).
	Logo lipgloss.Style
	// Box is the rounded border used for welcome banners and summaries.
	Box lipgloss.Style
	// BoxHeader is the title line inside a Box.
	BoxHeader lipgloss.Style
	// StepActive is the "→ Step N/M: ..." line while it's running.
	StepActive lipgloss.Style
	// StepDone is the "✓ ..." line after a step completes.
	StepDone lipgloss.Style
	// StepErr is the "✗ ..." line after a step fails.
	StepErr lipgloss.Style
	// StepDim is the label part of a done step (faded so the eye is drawn
	// to the next active step).
	StepDim lipgloss.Style
	// Heading is used for section headers inside boxes.
	Heading lipgloss.Style
	// Dim is for secondary / helper text.
	Dim lipgloss.Style
	// Accent is for highlighted values (URLs, file names).
	Accent lipgloss.Style
	// TreeBranch / TreeLast render the ├── / └── file tree.
	TreeBranch lipgloss.Style
	TreeLast  lipgloss.Style
	// Spinner is the animated spinner character style.
	Spinner lipgloss.Style
	// Green / Red are the status colors.
	Green lipgloss.Style
	Red   lipgloss.Style
}{
	Logo:       lipgloss.NewStyle().Bold(true),
	Box:        lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(1, 2).Width(52),
	BoxHeader:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252")),
	StepActive: lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
	StepDone:   lipgloss.NewStyle().Foreground(lipgloss.Color("42")),
	StepErr:    lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
	StepDim:    lipgloss.NewStyle().Foreground(lipgloss.Color("243")),
	Heading:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252")),
	Dim:        lipgloss.NewStyle().Foreground(lipgloss.Color("243")),
	Accent:     lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true),
	TreeBranch: lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
	TreeLast:  lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
	Spinner:   lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
	Green:     lipgloss.NewStyle().Foreground(lipgloss.Color("42")),
	Red:       lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
}

// isTTY reports whether stdout is a terminal. When it isn't (piped to a
// file, running in CI without a TTY, etc.) we fall back to plain text so
// the output remains readable.
func isTTY() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// gradientLogo renders "Glyra" with a purple → pink → cyan gradient using
// per-character TrueColor styles. It degrades gracefully on terminals
// without 24-bit color support thanks to lipgloss's automatic downsample.
func gradientLogo() string {
	text := "Glyra"
	var ranges []lipgloss.Range
	n := len(text)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		r := uint8(140 + t*115)
		g := uint8(50 + t*30)
		b := uint8(220 - t*90)
		hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		ranges = append(ranges, lipgloss.NewRange(i, i+1,
			lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Bold(true)))
	}
	return lipgloss.StyleRanges(text, ranges...)
}

// spinnerFrames are the braille-dot frames cycled by the step counter.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// frame returns the spinner character for the given tick index.
func frame(i int) string {
	return spinnerFrames[i%len(spinnerFrames)]
}

// treeLine renders a single file-tree line. `last` draws a └ instead of ├.
func treeLine(name string, last bool) string {
	prefix := "├── "
	if last {
		prefix = "└── "
	}
	return theme.TreeBranch.Render(prefix) + theme.Accent.Render(name)
}

// renderBox renders content inside a rounded border with an optional header.
func renderBox(content, header string) string {
	b := theme.Box
	if header != "" {
		fmt.Fprintf(os.Stdout, "%s\n", theme.BoxHeader.Render(header))
	}
	return b.Render(content)
}

// indent prefixes every line in s with two spaces.
func indent(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = "  " + l
	}
	return strings.Join(lines, "\n")
}
