package cli

import "fmt"

// PrintHelp renders the top-level usage screen inside a rounded box with
// the gradient logo — similar in spirit to `claude --help` or `gh`.
func PrintHelp() {
	fmt.Println()
	fmt.Println(gradientLogo())
	fmt.Println()
	fmt.Println(renderBox(
		theme.Heading.Render("Usage")+"\n\n"+
			theme.Accent.Render("  glyra init <name>")+theme.Dim.Render("   Scaffold a new Glyra project")+"\n"+
			theme.Accent.Render("  glyra add <cmp>")+theme.Dim.Render("     Inject a UI component (e.g. button, card)")+"\n"+
			theme.Accent.Render("  glyra dev")+theme.Dim.Render("           Run in development mode (HMR for React)")+"\n"+
			theme.Accent.Render("  glyra build")+theme.Dim.Render("         Compile into a single production binary")+"\n"+
			theme.Accent.Render("  glyra version")+theme.Dim.Render("       Print the Glyra CLI version")+"\n"+
			theme.Accent.Render("  glyra help")+theme.Dim.Render("          Show this message")+"\n\n"+
			theme.Dim.Render("Docs: https://ecocee.in/glyra"),
		"",
	))
	fmt.Println()
}
