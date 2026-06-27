package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a project name. Example: glyra init my-app")
			os.Exit(1)
		}
		initProject(os.Args[2])
	case "dev":
		fmt.Println("🚀 Starting Glyra development server...")
		fmt.Println("Please run 'npm run dev' in the frontend folder, and 'go run main.go' in the root.")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Glyra CLI - Build high-performance desktop apps")
	fmt.Println("\nUsage:")
	fmt.Println("  glyra init <project-name>  Initialize a new Glyra project")
	fmt.Println("  glyra build                Build the binary for production")
	fmt.Println("  glyra dev                  Run the application in development mode")
}

func initProject(name string) {
	fmt.Printf("\n✨ Welcome to Glyra! Let's build '%s'.\n", name)
	
	fmt.Println("\nSelect a frontend template:")
	fmt.Println("  1) ⚛️  Vite + React (Recommended)")
	fmt.Println("  2) 🌐 Vanilla HTML/CSS/JS")
	fmt.Print("\nEnter choice [1]: ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Scaffold base Go structure
	createBaseStructure(name)
	
	if choice == "2" {
		scaffoldVanilla(name)
	} else {
		scaffoldReact(name)
	}

	// Initialize Go module
	fmt.Println("\n📦 Initializing Go module...")
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = name
	cmd.Run()

	cmd = exec.Command("go", "get", "github.com/webview/webview_go")
	cmd.Dir = name
	cmd.Run()

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = name
	cmd.Run()

	fmt.Println("\n✅ Success! Your Glyra project is ready.")
	
	if choice == "2" {
		fmt.Printf("\nNext steps:\n  cd %s\n  go run main.go\n\n", name)
	} else {
		fmt.Printf("\nNext steps:\n  cd %s/frontend\n  npm install\n  npm run dev\n\nIn a new terminal tab:\n  cd %s\n  go run main.go\n\n", name, name)
	}
}

func createBaseStructure(name string) {
	dirs := []string{
		name,
		filepath.Join(name, "frontend"),
	}
	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}
}

func writeString(path, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}
}

// ---------------------------------------------------------
// VANILLA HTML/CSS/JS TEMPLATE
// ---------------------------------------------------------

func scaffoldVanilla(name string) {
	fmt.Println("⚡ Scaffolding Vanilla HTML/CSS/JS template...")
	
	writeString(filepath.Join(name, "main.go"), goVanillaMainTemplate)
	
	writeString(filepath.Join(name, "frontend", "index.html"), htmlVanillaTemplate)
}

var goVanillaMainTemplate = `package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/webview/webview_go"
)

func main() {
	// 1. Serve the frontend folder locally
	go func() {
		http.Handle("/", http.FileServer(http.Dir("./frontend")))
		http.ListenAndServe("127.0.0.1:8080", nil)
	}()
	time.Sleep(100 * time.Millisecond) // Give server a moment to start

	// 2. Start WebView
	debug := true
	w := webview.New(debug)
	defer w.Destroy()

	w.SetTitle("Glyra — Vanilla UI")
	w.SetSize(900, 650, webview.HintNone)

	// Bridge Go function to JS
	w.Bind("GetSystemStatus", func() string {
		return "All systems optimal. Go backend running perfectly!"
	})

	w.Navigate("http://127.0.0.1:8080")
	fmt.Println("Glyra Backend started.")
	w.Run()
}
`

var htmlVanillaTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Glyra Vanilla UI</title>
  <style>
    :root {
      --bg: #09090b;
      --card-bg: rgba(255, 255, 255, 0.03);
      --text: #fafafa;
      --text-muted: #a1a1aa;
      --primary: #3b82f6;
      --primary-hover: #2563eb;
    }
    
    body {
      margin: 0;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
      background: radial-gradient(circle at top right, #1e1e2e, var(--bg));
      color: var(--text);
      display: flex;
      justify-content: center;
      align-items: center;
      height: 100vh;
      overflow: hidden;
    }

    .glass-card {
      background: var(--card-bg);
      backdrop-filter: blur(16px);
      -webkit-backdrop-filter: blur(16px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 16px;
      padding: 3rem;
      text-align: center;
      box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
      max-width: 500px;
      width: 100%;
    }

    h1 { margin: 0 0 0.5rem 0; font-size: 2.5rem; font-weight: 700; letter-spacing: -0.05em; }
    p { color: var(--text-muted); margin-bottom: 2rem; line-height: 1.6; }

    .status-box {
      background: rgba(0,0,0,0.4);
      padding: 1rem;
      border-radius: 8px;
      margin-bottom: 2rem;
      font-family: monospace;
      color: #10b981;
    }

    button {
      background: linear-gradient(135deg, var(--primary), var(--primary-hover));
      color: white;
      border: none;
      padding: 0.875rem 2rem;
      font-size: 1rem;
      font-weight: 600;
      border-radius: 999px;
      cursor: pointer;
      transition: transform 0.2s, box-shadow 0.2s;
    }
    button:hover {
      transform: translateY(-2px);
      box-shadow: 0 10px 15px -3px rgba(59, 130, 246, 0.4);
    }
  </style>
</head>
<body>

  <div class="glass-card">
    <h1>Glyra</h1>
    <p>Pure HTML/CSS/JS frontend powered by a blazing fast Go backend.</p>
    
    <div class="status-box" id="statusBox">Waiting for Go bridge...</div>
    
    <button onclick="pingBackend()">Ping Go Backend</button>
  </div>

  <script>
    async function pingBackend() {
      const box = document.getElementById('statusBox');
      box.innerText = "Pinging...";
      try {
        const response = await window.GetSystemStatus();
        setTimeout(() => { box.innerText = response; }, 400); // slight delay for effect
      } catch (err) {
        box.innerText = "Error: Cannot reach Go backend.";
        box.style.color = "#ef4444";
      }
    }
  </script>
</body>
</html>
`

// ---------------------------------------------------------
// REACT + VITE TEMPLATE
// ---------------------------------------------------------

func scaffoldReact(name string) {
	fmt.Println("⚛️  Scaffolding Vite + React template...")
	
	os.MkdirAll(filepath.Join(name, "frontend", "src"), 0755)

	writeString(filepath.Join(name, "main.go"), goReactMainTemplate)
	writeString(filepath.Join(name, "frontend", "package.json"), packageJsonTemplate)
	writeString(filepath.Join(name, "frontend", "vite.config.js"), viteConfigTemplate)
	writeString(filepath.Join(name, "frontend", "index.html"), htmlReactTemplate)
	writeString(filepath.Join(name, "frontend", "src", "main.jsx"), reactMainTemplate)
	writeString(filepath.Join(name, "frontend", "src", "App.jsx"), reactAppTemplate)
}

var goReactMainTemplate = `package main

import (
	"fmt"
	"github.com/webview/webview_go"
)

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()

	w.SetTitle("Glyra — React + Vite")
	w.SetSize(900, 650, webview.HintNone)

	// Bridge Go function to React
	w.Bind("GetSystemStatus", func() string {
		return "All systems optimal. Go backend running perfectly!"
	})

	// Connect to Vite HMR server!
	w.Navigate("http://localhost:5173")
	
	fmt.Println("Glyra Backend started. Waiting for Vite on port 5173...")
	w.Run()
}
`

var packageJsonTemplate = `{
  "name": "frontend",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@vitejs/plugin-react": "^4.2.1",
    "vite": "^5.2.0"
  }
}
`

var viteConfigTemplate = `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    strictPort: true,
  }
})
`

var htmlReactTemplate = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Glyra App</title>
    <style>
      body { margin: 0; padding: 0; background-color: #09090b; }
    </style>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.jsx"></script>
  </body>
</html>
`

var reactMainTemplate = `import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
`

var reactAppTemplate = `import { useState } from 'react'

function App() {
  const [status, setStatus] = useState("Waiting for Go backend...")

  const checkStatus = async () => {
    try {
      setStatus("Pinging...");
      const res = await window.GetSystemStatus();
      setTimeout(() => setStatus(res), 400);
    } catch (e) {
      setStatus("Error communicating with Go backend.");
    }
  }

  return (
    <div style={{
      display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh',
      fontFamily: 'system-ui, sans-serif', color: '#fafafa',
      background: 'radial-gradient(circle at top right, #1e1e2e, #09090b)'
    }}>
      
      <div style={{
        background: 'rgba(255,255,255,0.03)', backdropFilter: 'blur(16px)',
        border: '1px solid rgba(255,255,255,0.1)', borderRadius: '16px',
        padding: '3rem', textAlign: 'center', boxShadow: '0 25px 50px -12px rgba(0,0,0,0.5)',
        maxWidth: '500px', width: '100%'
      }}>
        <h1 style={{ margin: '0 0 0.5rem 0', fontSize: '2.5rem', fontWeight: 700, letterSpacing: '-0.05em' }}>
          Glyra React
        </h1>
        <p style={{ color: '#a1a1aa', marginBottom: '2rem', lineHeight: 1.6 }}>
          Blazing fast React frontend powered by Vite, running inside a native webview driven by Go.
        </p>
        
        <div style={{
          background: 'rgba(0,0,0,0.4)', padding: '1rem', borderRadius: '8px', 
          marginBottom: '2rem', fontFamily: 'monospace', color: '#10b981'
        }}>
          {status}
        </div>
        
        <button 
          onClick={checkStatus}
          style={{
            background: 'linear-gradient(135deg, #3b82f6, #2563eb)', color: 'white',
            border: 'none', padding: '0.875rem 2rem', fontSize: '1rem', fontWeight: 600,
            borderRadius: '999px', cursor: 'pointer', transition: 'transform 0.2s',
            boxShadow: '0 4px 6px -1px rgba(59, 130, 246, 0.4)'
          }}
          onMouseOver={e => e.currentTarget.style.transform = 'translateY(-2px)'}
          onMouseOut={e => e.currentTarget.style.transform = 'translateY(0)'}
        >
          Ping Go Backend
        </button>
      </div>

    </div>
  )
}

export default App
`
