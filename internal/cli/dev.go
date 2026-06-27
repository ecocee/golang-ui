package cli

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// runDev runs the current Glyra project in development mode. It prints a
// status dashboard showing live frontend / backend URLs with green dots
// that turn red if the process dies.
func runDev(args []string) error {
	if !fileExists("main.go") {
		return fmt.Errorf("no main.go here — run this from inside a Glyra project")
	}

	var frontend *exec.Cmd
	if fileExists("frontend/package.json") {
		frontend = exec.Command("npm", "run", "dev")
		frontend.Dir = "frontend"
		frontend.Stdout = os.Stdout
		frontend.Stderr = os.Stderr
		if err := frontend.Start(); err != nil {
			return fmt.Errorf("starting frontend: %w", err)
		}
		defer func() {
			if frontend.Process != nil {
				frontend.Process.Kill()
			}
		}()
	}

	isVanilla := !fileExists("frontend/package.json")
	return runLiveBackend(isVanilla)
}

func pollFiles(lastMod time.Time, isVanilla bool) (bool, bool, time.Time) {
	var latest time.Time
	changedGo := false
	changedWeb := false

	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules") {
			return filepath.SkipDir
		}
		// In React apps, Vite handles frontend HMR. So we ignore the frontend folder entirely for the Go watcher.
		if d.IsDir() && d.Name() == "frontend" && !isVanilla {
			return filepath.SkipDir
		}
		
		ext := filepath.Ext(path)
		isGoFile := ext == ".go"
		isWebFile := isVanilla && (ext == ".html" || ext == ".css" || ext == ".js")
		
		if !d.IsDir() && (isGoFile || isWebFile) {
			info, _ := d.Info()
			if info != nil {
				mod := info.ModTime()
				if mod.After(lastMod) {
					if mod.After(latest) {
						latest = mod
					}
					if isGoFile {
						changedGo = true
					} else if isWebFile {
						changedWeb = true
					}
				}
				// We also just need to track the absolute latest to initialize
				if mod.After(latest) {
					latest = mod
				}
			}
		}
		return nil
	})
	
	if changedGo || changedWeb {
		return true, changedGo, latest
	}
	return false, false, latest
}

func runLiveBackend(isVanilla bool) error {
	fmt.Println(renderBox(
		theme.Heading.Render("🚀 Glyra dev")+"\n\n"+
			theme.Dim.Render("Backend:  ")+theme.Accent.Render("http://127.0.0.1:8080")+"  "+liveDot(8080)+"\n"+
			theme.Dim.Render("\nPress Ctrl+C to stop.\nEditing *.go files will hot-reload automatically.\nPress 'r' + Enter to manually reload the backend!"),
		"",
	))
	fmt.Println()

	var backend *exec.Cmd
	_, _, lastModTime := pollFiles(time.Time{}, isVanilla)

	startBackend := func() {
		fmt.Println(theme.Dim.Render("Compiling..."))
		buildCmd := exec.Command("go", "build", "-o", ".glyra-dev-bin", ".")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			fmt.Println(theme.Red.Render("\n❌ Build failed! Waiting for next file change..."))
			return
		}

		backend = exec.Command("./.glyra-dev-bin")
		backend.Env = append(os.Environ(), "GLYRA_DEV=true")
		backend.Stdout = os.Stdout
		backend.Stderr = os.Stderr
		backend.Start()
	}

	startBackend()
	defer func() {
		if backend != nil && backend.Process != nil {
			backend.Process.Kill()
		}
		os.Remove(".glyra-dev-bin")
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	// Manual reload trigger via stdin
	reloadSig := make(chan struct{})
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if strings.TrimSpace(text) == "r" || strings.TrimSpace(text) == "R" {
				reloadSig <- struct{}{}
			}
		}
	}()

	for {
		select {
		case <-sig:
			fmt.Println(theme.Dim.Render("\nShutting down…"))
			return nil
		case <-reloadSig:
			_, _, lastModTime = pollFiles(time.Time{}, isVanilla)
			fmt.Println(theme.Accent.Render("\n🔄 Manual reload triggered! Restarting backend..."))
			if backend != nil && backend.Process != nil {
				backend.Process.Kill()
				backend.Wait()
			}
			startBackend()
		case <-ticker.C:
			changed, isGoChange, mod := pollFiles(lastModTime, isVanilla)
			if changed {
				lastModTime = mod
				if isGoChange {
					fmt.Println(theme.Accent.Render("\n🔄 Go code changed! Restarting backend..."))
					if backend != nil && backend.Process != nil {
						backend.Process.Kill()
						backend.Wait()
					}
					startBackend()
				} else {
					fmt.Println(theme.Accent.Render("\n⚡ Web asset changed! Refreshing UI in-place..."))
					// Trigger our custom hot reload endpoint in the Vanilla backend!
					http.Get("http://127.0.0.1:8080/__glyra_reload")
				}
			}
		}
	}
}

// liveDot returns a green ● if the TCP port is reachable, red ● otherwise.
func liveDot(port int) string {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
	if err != nil {
		return theme.Red.Render("●")
	}
	_ = conn.Close()
	return theme.Green.Render("●")
}
