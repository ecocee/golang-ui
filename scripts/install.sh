#!/usr/bin/env bash
# Install the Glyra CLI from source.
#
# Useful in CI or on a dev machine when you want the binary built from the
# current tree (with a real version stamp) rather than `go run`.
#
# Usage:
#   scripts/install.sh                 # install dev build to $(go env GOPATH)/bin
#   scripts/install.sh v0.2.0-beta.1   # install a stamped build
#   scripts/install.sh --path ./bin    # install into a custom directory

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

VERSION="dev"
INSTALL_DIR="$(go env GOPATH)/bin"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --path)
      INSTALL_DIR="$2"
      shift 2
      ;;
    *)
      VERSION="$1"
      shift
      ;;
  esac
done

LDFLAGS="-s -w -X github.com/ecocee/golang-ui/internal/cli.Version=${VERSION}"

echo "==> building glyra (version=${VERSION})"
mkdir -p "$INSTALL_DIR"
go build -ldflags="${LDFLAGS}" -o "${INSTALL_DIR}/glyra" ./cmd/glyra

echo "==> installed to ${INSTALL_DIR}/glyra"
"${INSTALL_DIR}/glyra" --version
