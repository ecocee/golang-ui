package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// step is a single named operation shown in the step counter.
type step struct {
	name string
	run  func() error
}

// runBuild produces a production-ready binary for the current Glyra project.
// It renders an animated step counter:
//
//	→ Step 1/3: Building frontend...            ⠋
//	✓ Step 2/3: Compiling Go binary...          ⠙
//	✓ Step 3/3: Build complete!                 ✅
func runBuild(args []string) error {
	if !fileExists("main.go") {
		return fmt.Errorf("no main.go here — run this from inside a Glyra project")
	}

	// Build the ordered list of steps. Each step has a name and a func.
	// The step counter drives the UI; the func does the real work.
	var steps []step

	hasFrontend := fileExists("frontend/package.json")
	if hasFrontend {
		steps = append(steps, step{
			name: "Building frontend (Vite)",
			run: func() error {
				cmd := exec.Command("npm", "run", "build")
				cmd.Dir = "frontend"
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				return cmd.Run()
			},
		})
	}

	steps = append(steps, step{
		name: "Generating App Icons",
		run: func() error {
			if !fileExists("assets/appicon.png") {
				return nil
			}
			if runtime.GOOS == "windows" {
				os.WriteFile("winres.json", []byte(`{"RT_ICON":{"MAIN":{"0000":["assets/appicon.png"]}}}`), 0644)
				cmd := exec.Command("go", "run", "github.com/tc-hib/go-winres@latest", "make")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				os.Remove("winres.json")
			}
			return nil
		},
	})

	steps = append(steps,
		step{
			name: "Compiling Go binary",
			run: func() error {
				output := binaryName()
				ldflags := "-s -w"
				if runtime.GOOS == "windows" {
					ldflags += " -H=windowsgui"
				}
				cmd := exec.Command("go", "build", "-ldflags="+ldflags, "-o", output, ".")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				return cmd.Run()
			},
		},
		step{
			name: "Verifying binary",
			run: func() error {
				output := binaryName()
				if !fileExists(output) {
					return fmt.Errorf("binary %s not found after build", output)
				}
				return nil
			},
		},
		step{
			name: "Packaging for OS",
			run: func() error {
				output := binaryName()
				if runtime.GOOS == "darwin" {
					appName := output + ".app"
					os.MkdirAll(filepath.Join(appName, "Contents", "MacOS"), 0755)
					os.MkdirAll(filepath.Join(appName, "Contents", "Resources"), 0755)
					
					os.Rename(output, filepath.Join(appName, "Contents", "MacOS", output))
					
					plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>%s</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleIdentifier</key>
	<string>com.glyra.%s</string>
</dict>
</plist>`, output, output)
					os.WriteFile(filepath.Join(appName, "Contents", "Info.plist"), []byte(plist), 0644)
					
					if fileExists("assets/appicon.png") {
						os.MkdirAll("assets/icon.iconset", 0755)
						
						sizes := []string{"128", "256", "512"}
						for _, size := range sizes {
							cmd := exec.Command("sips", "-z", size, size, "assets/appicon.png", "--out", fmt.Sprintf("assets/icon.iconset/icon_%sx%s.png", size, size))
							cmd.Run()
							
							// Also add @2x retina variants
							cmd2x := exec.Command("sips", "-z", size, size, "assets/appicon.png", "--out", fmt.Sprintf("assets/icon.iconset/icon_%sx%s@2x.png", size, size))
							cmd2x.Run()
						}
						
						cmd := exec.Command("iconutil", "-c", "icns", "assets/icon.iconset", "-o", filepath.Join(appName, "Contents", "Resources", "icon.icns"))
						if err := cmd.Run(); err != nil {
							fmt.Println("\nWarning: Failed to generate Mac icon. Ensure assets/appicon.png is a valid square PNG.")
						}
						os.RemoveAll("assets/icon.iconset")
					}
					
					// Force macOS Finder to refresh the app icon cache
					exec.Command("touch", appName).Run()
				} else if runtime.GOOS == "linux" {
					wd, _ := os.Getwd()
					desktop := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s/%s
Icon=%s/assets/appicon.png
Terminal=false`, output, wd, output, wd)
					os.WriteFile(output+".desktop", []byte(desktop), 0644)
				}
				return nil
			},
		},
	)

	// Render the step counter.
	if err := renderSteps(steps); err != nil {
		return err
	}

	// Final success banner.
	fmt.Println()
	name := binaryName()
	fmt.Println(renderBox(
		theme.Heading.Render("✓ Build complete!")+"\n\n"+
			theme.Dim.Render("Binary: ")+theme.Accent.Render(name)+"\n"+
			theme.Dim.Render("Run it:  ")+theme.Accent.Render("./"+name),
		"",
	))
	return nil
}

// renderSteps runs each step with an animated spinner and a live N/M counter.
// Each step occupies its own line; the spinner animates in place via the
// returned carriage-return trick (no full TUI event loop needed).
func renderSteps(steps []step) error {
	total := len(steps)
	for i, s := range steps {
		// Print the "→ Step N/M: name... ⠋" line.
		prefix := fmt.Sprintf("→ Step %d/%d: ", i+1, total)
		line := theme.StepActive.Render(prefix) + theme.StepActive.Render(s.name)
		fmt.Fprintf(os.Stdout, "%s  %s", line, theme.Spinner.Render(frame(0)))

		// Animate the spinner while the step runs.
		done := make(chan error, 1)
		go func() { done <- s.run() }()

		tick := 0
		ticker := time.NewTicker(80 * time.Millisecond)
		running := true
		for running {
			select {
			case err := <-done:
				<-ticker.C // drain
				running = false
				if err != nil {
					// Replace spinner with red ✗ and print the error.
					fmt.Fprintf(os.Stdout, "\r%s  %s\n",
						theme.StepErr.Render(prefix),
						theme.StepErr.Render(s.name+" ✗"))
					fmt.Println(theme.Red.Render(indent("  Error: " + err.Error())))
					return err
				}
				// Replace spinner with green ✓.
				fmt.Fprintf(os.Stdout, "\r%s  %s\n",
					theme.StepDone.Render(prefix),
					theme.StepDone.Render(s.name+" ✓"))
			case <-ticker.C:
				tick++
				fmt.Fprintf(os.Stdout, "\r%s  %s  %s",
					theme.StepActive.Render(prefix),
					theme.StepActive.Render(s.name),
					theme.Spinner.Render(frame(tick)))
			}
		}
		ticker.Stop()
	}
	return nil
}

func binaryName() string {
	name := "app"
	if wd, err := os.Getwd(); err == nil {
		name = filepath.Base(wd)
	}
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}
