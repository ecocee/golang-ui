package cli

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
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

	if fileExists("frontend/package.json") {
		return runReactDev()
	}
	return runVanillaDev()
}

func runVanillaDev() error {
	fmt.Println(renderBox(
		theme.Heading.Render("🚀 Glyra dev")+"\n\n"+
			theme.Dim.Render("Backend:  ")+theme.Accent.Render("http://127.0.0.1:8080")+"  "+liveDot(8080)+"\n"+
			theme.Dim.Render("\nPress Ctrl+C to stop"),
		"",
	))
	fmt.Println()

	cmd := exec.Command("go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runReactDev() error {
	// Start both servers.
	frontend := exec.Command("npm", "run", "dev")
	frontend.Dir = "frontend"
	frontend.Stdout = os.Stdout
	frontend.Stderr = os.Stderr

	backend := exec.Command("go", "run", "main.go", "-dev")
	backend.Stdout = os.Stdout
	backend.Stderr = os.Stderr

	if err := frontend.Start(); err != nil {
		return fmt.Errorf("starting frontend: %w", err)
	}
	if err := backend.Start(); err != nil {
		if frontend.Process != nil {
			_ = frontend.Process.Kill()
		}
		return fmt.Errorf("starting backend: %w", err)
	}

	// Dashboard header.
	fmt.Println(renderBox(
		theme.Heading.Render("🚀 Glyra dev")+"\n\n"+
			theme.Dim.Render("Frontend: ")+theme.Accent.Render("http://localhost:5173")+"  "+liveDot(5173)+"\n"+
			theme.Dim.Render("Backend:  ")+theme.Accent.Render("http://127.0.0.1:8080")+"  "+liveDot(8080)+"\n"+
			theme.Dim.Render("\nPress Ctrl+C to stop"),
		"",
	))
	fmt.Println()

	// Wait for Ctrl+C or backend exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	backendDone := make(chan error, 1)
	go func() { backendDone <- backend.Wait() }()

	select {
	case <-sig:
		fmt.Println(theme.Dim.Render("\nShutting down…"))
	case err := <-backendDone:
		if err != nil {
			fmt.Println(theme.Red.Render("\nBackend exited with error: " + err.Error()))
		}
	}

	if frontend.Process != nil {
		_ = frontend.Process.Kill()
	}
	if backend.Process != nil {
		_ = backend.Process.Kill()
	}
	return nil
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
