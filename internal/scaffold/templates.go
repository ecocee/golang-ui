package scaffold

import "embed"

// templatesFS embeds every starter template (vanilla, react, ...) that
// ships with the Glyra CLI. Files ending in .tmpl are rendered through
// text/template when a project is scaffolded; everything else is copied
// verbatim. The "all:" prefix keeps dotfiles like .gitignore.tmpl and
// frontend/dist/.gitkeep in the embedded filesystem.
//
//go:embed all:templates
var templatesFS embed.FS
