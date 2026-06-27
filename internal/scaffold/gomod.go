package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// initGoModule sets up the Go module for a freshly scaffolded project:
// it creates go.mod and fetches the webview dependency every Glyra app
// needs to render its native window.
func initGoModule(dir string) error {
	fmt.Println("📦 Initializing Go module...")

	steps := [][]string{
		{"go", "mod", "init", filepath.Base(dir)},
		{"go", "get", "github.com/webview/webview_go"},
		{"go", "mod", "tidy"},
	}

	for _, args := range steps {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("running %q: %w", strings.Join(args, " "), err)
		}
	}
	return nil
}
