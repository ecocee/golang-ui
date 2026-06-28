package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed all:templates
var templatesFS embed.FS

// Template identifies which starter frontend a project is scaffolded with.
type Template string

const (
	React   Template = "react"
	ReactTS Template = "react-ts"
	NextJS  Template = "nextjs"
	Svelte  Template = "svelte"
	Vanilla Template = "vanilla"
)

// Data is the set of variables available inside every *.tmpl file.
type Data struct {
	ProjectName string // used for package names or directory names
	PackageName string // lowercased, npm-safe form of ProjectName
	Title       string // Human readable app title for window and HTML
	ModulePath  string // Go module path; lets cmd/<name>/main.go import <ModulePath>/src
}

// New scaffolds a brand new Glyra project named `name` on disk using
// `tmpl`, then initializes its Go module. modulePath is the Go module path
// written into go.mod (typically the project name); it lets the generated
// cmd/<name>/main.go import the project's own src/ package.
func New(name, title, modulePath string, tmpl Template) error {
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory %q already exists", name)
	}

	data := Data{
		ProjectName: name,
		PackageName: strings.ToLower(name),
		Title:       title,
		ModulePath:  modulePath,
	}

	fmt.Printf("⚡ Scaffolding %s template...\n", tmpl)
	if err := copyTemplate(string(tmpl), name, data); err != nil {
		return err
	}

	return initGoModule(name)
}

// copyTemplate walks the embedded template directory `tmpl`, rendering
// every *.tmpl file through text/template and copying everything else
// into the destination directory byte-for-byte. Template syntax in
// directory and file names (e.g. cmd/{{.ProjectName}}) is also rendered,
// so dynamic paths resolve alongside file contents.
func copyTemplate(tmpl, dest string, data Data) error {
	root := filepath.Join("templates", tmpl)

	return fs.WalkDir(templatesFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		// Render template syntax in the path itself (both dirs and files),
		// then strip the .tmpl suffix to get the on-disk target.
		renderedRel, err := renderPath(rel, data)
		if err != nil {
			return fmt.Errorf("rendering path %s: %w", rel, err)
		}
		target := filepath.Join(dest, strings.TrimSuffix(renderedRel, ".tmpl"))

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		raw, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}

		if !strings.HasSuffix(renderedRel, ".tmpl") {
			return os.WriteFile(target, raw, 0o644)
		}

		t, err := template.New(d.Name()).Parse(string(raw))
		if err != nil {
			return fmt.Errorf("parsing template %s: %w", path, err)
		}

		f, err := os.Create(target)
		if err != nil {
			return err
		}
		defer f.Close()

		return t.Execute(f, data)
	})
}

// renderPath renders template syntax (e.g. {{.ProjectName}}) in every
// segment of rel, so dynamic directory names like cmd/{{.ProjectName}} resolve
// correctly. Each segment is rendered independently so template errors
// reference a short, meaningful name.
func renderPath(rel string, data Data) (string, error) {
	rel = filepath.ToSlash(rel)
	parts := strings.Split(rel, "/")
	for i, p := range parts {
		t, err := template.New("path").Parse(p)
		if err != nil {
			return "", err
		}
		var buf strings.Builder
		if err := t.Execute(&buf, data); err != nil {
			return "", err
		}
		parts[i] = buf.String()
	}
	return filepath.FromSlash(strings.Join(parts, "/")), nil
}
