package cli

import (
	"fmt"
)

// Run parses CLI arguments and dispatches to the matching command.
func Run(args []string) error {
	if len(args) == 0 {
		PrintHelp()
		return nil
	}

	switch args[0] {
	case "init":
		return runInit(args[1:])
	case "add":
		return runAdd(args[1:])
	case "dev":
		return runDev(args[1:])
	case "build":
		return runBuild(args[1:])
	case "version", "-v", "--version":
		fmt.Println(gradientLogo() + " " + theme.Dim.Render(Version))
		return nil
	case "help", "-h", "--help":
		PrintHelp()
		return nil
	default:
		PrintHelp()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

