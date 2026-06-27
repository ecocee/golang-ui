#!/usr/bin/env bash
# Cross-compile the Glyra CLI for every supported platform.
#
# Mirrors the GitHub Actions release matrix so a local build reproduces CI
# artifacts. The resulting binaries (and tarballs/zips) land in ./dist/.
#
# Usage:
#   scripts/build.sh                 # build all platforms, untagged dev build
#   scripts/build.sh v0.2.0-beta.1   # build all platforms, stamped with version
#   scripts/build.sh --single        # build only the current host platform

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

VERSION="${1:-dev}"
SINGLE=0
if [[ "${1:-}" == "--single" ]]; then
  SINGLE=1
  VERSION="${2:-dev}"
fi

LDFLAGS="-s -w -X github.com/ecocee/golang-ui/internal/cli.Version=${VERSION}"

# (target, goos, goarch)
platforms=(
  "darwin-arm64 darwin arm64"
  "darwin-amd64 darwin amd64"
  "linux-amd64 linux amd64"
  "linux-arm64 linux arm64"
  "windows-amd64 windows amd64"
)

mkdir -p dist release

build_one() {
  local target="$1" goos="$2" goarch="$3"
  local output="glyra-${target}"
  if [[ "$goos" == "windows" ]]; then
    output="${output}.exe"
  fi

  echo "==> building ${target} (${goos}/${goarch})"
  GOOS="$goos" GOARCH="$goarch" go build \
    -ldflags="${LDFLAGS}" \
    -o "dist/${output}" \
    ./cmd/glyra

  if [[ "$goos" == "windows" ]]; then
    (cd dist && zip -q "../release/${target}.zip" "$output")
  else
    (cd dist && tar czf "../release/${target}.tar.gz" "$output")
  fi
}

if [[ "$SINGLE" -eq 1 ]]; then
  build_one "$(go env GOOS)-$(go env GOARCH)" "$(go env GOOS)" "$(go env GOARCH)"
else
  for p in "${platforms[@]}"; do
    build_one $p
  done
fi

echo "==> artifacts in ./release/"
ls -1 release/
