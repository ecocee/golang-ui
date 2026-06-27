package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"charm.land/huh/v2"

	"github.com/ecocee/golang-ui/internal/scaffold"
)

func runInit(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing project name — usage: glyra init <project-name>")
	}
	name := args[0]

	fmt.Println()
	fmt.Println(gradientLogo())
	fmt.Println(theme.BoxHeader.Render("  Build small, fast, native desktop apps"))
	fmt.Println(theme.Dim.Render("  with Go + a web UI"))
	fmt.Println()

	// Interactive template picker. Use huh when we have a TTY, fall back
	// to a plain bufio reader when stdin is piped (CI, `echo 1 | glyra init`).
	var tmpl scaffold.Template
	var choice string

	if isTTY() {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose a frontend template").
					Options(
						huh.NewOption("⚛️  Vite + React (TypeScript) [Recommended]", "react-ts"),
						huh.NewOption("⚛️  Vite + React (JavaScript)", "react"),
						huh.NewOption("▲  Next.js (Static Export)", "nextjs"),
						huh.NewOption("🌐 Vanilla HTML / CSS / JS", "vanilla"),
					).
					Value(&choice),
			),
		).WithTheme(huh.ThemeFunc(huh.ThemeCatppuccin))

		if err := form.Run(); err != nil {
			return err
		}
	} else {
		fmt.Print("Select template  1) React TS  2) React JS  3) Next.js  4) Vanilla HTML/CSS/JS  [1]: ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ = reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
	}

	switch choice {
	case "2", "react":
		tmpl = scaffold.React
	case "3", "nextjs":
		tmpl = scaffold.NextJS
	case "4", "vanilla":
		tmpl = scaffold.Vanilla
	default:
		tmpl = scaffold.ReactTS
	}

	title := name

	// Scaffold.
	if err := scaffold.New(name, title, tmpl); err != nil {
		return err
	}

	// Success summary.
	fmt.Println(renderBox(buildSummary(name, tmpl), "✓ Project ready"))
	fmt.Println()
	fmt.Println(theme.Heading.Render("Next steps:"))
	fmt.Printf("  %s %s\n", theme.Accent.Render("cd"), name)
	if tmpl == scaffold.React {
		fmt.Printf("  %s\n", theme.Dim.Render("cd frontend && npm install && cd .."))
	}
	fmt.Printf("  %s\n", theme.Dim.Render("glyra dev"))
	fmt.Println()

	return nil
}

// buildSummary produces the indented file-tree block shown after a
// successful scaffold.
func buildSummary(name string, tmpl scaffold.Template) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Scaffolded %s with the %s template.\n\n", theme.Accent.Render(name), theme.Accent.Render(string(tmpl)))

	files := scaffoldTree(name, tmpl)
	for i, f := range files {
		fmt.Fprintln(&b, treeLine(f, i == len(files)-1))
	}
	return b.String()
}

// scaffoldTree returns the list of files a scaffold produces, in the
// order they appear in the tree.
func scaffoldTree(name string, tmpl scaffold.Template) []string {
	if tmpl == scaffold.React {
		return []string{
			"main.go",
			"frontend/",
			"  package.json",
			"  vite.config.js",
			"  index.html",
			"  src/",
			"    main.jsx",
			"    App.jsx",
			"    App.css",
			"    index.css",
		}
	}
	return []string{
		"main.go",
		"frontend/",
		"  index.html",
		"  style.css",
		"  app.js",
	}
}
