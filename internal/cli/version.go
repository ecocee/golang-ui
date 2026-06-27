package cli

// Version is the current Glyra CLI release. It defaults to the value
// baked in at compile time but can be overridden via -ldflags for CI
// builds, e.g.:
//
//	go build -ldflags="-X github.com/ecocee/golang-ui/internal/cli.Version=v0.2.0-beta.1"
var Version = "0.2.0"
