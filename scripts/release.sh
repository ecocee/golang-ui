#!/usr/bin/env bash
# Cut a new Glyra release.
#
# Does the hand-heavy parts of a release in one place so the tag, the
# release notes, and the version baked into the CLI never drift:
#
#   1. sanity-checks the working tree and the version string
#   2. stamps cli.Version, builds every platform, and smoke-tests the binary
#   3. generates RELEASE_NOTES.md (consumed by .github/workflows/release.yml)
#   4. commits the version bump, tags, and pushes
#
# Usage:
#   scripts/release.sh v0.2.0-beta.1

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ $# -ne 1 ]]; then
  echo "usage: scripts/release.sh <tag>" >&2
  echo "example: scripts/release.sh v0.2.0-beta.1" >&2
  exit 1
fi

TAG="$1"
# Require a leading 'v' — matches the 'v*' tag pattern in release.yml.
if [[ ! "$TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
  echo "error: tag must match vMAJOR.MINOR.PATCH[-prerelease], got '$TAG'" >&2
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "error: working tree is dirty — commit or stash before releasing" >&2
  git status --short >&2
  exit 1
fi

if git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "error: tag $TAG already exists" >&2
  exit 1
fi

echo "==> bumping cli.Version -> ${TAG}"
if grep -q 'var Version = ' internal/cli/version.go; then
  sed -i.bak "s/var Version = \".*\"/var Version = \"${TAG}\"/" internal/cli/version.go
  rm -f internal/cli/version.go.bak
else
  echo "error: could not find Version declaration in internal/cli/version.go" >&2
  exit 1
fi

echo "==> building all platforms"
bash scripts/build.sh "$TAG"

echo "==> smoke-testing host binary"
"dist/glyra-$(go env GOOS)-$(go env GOARCH)" --version

echo "==> generating RELEASE_NOTES.md"
cat > RELEASE_NOTES.md <<EOF
## Glyra ${TAG}

$(sed -n '/^## /,/^## /p' CHANGELOG.md 2>/dev/null | sed '$d' | grep -v "^## Glyra ${TAG}" || echo "See the [CHANGELOG](CHANGELOG.md) for details.")

### Platforms
| Asset | OS | Arch |
|-------|----|------|
| \`glyra-darwin-arm64.tar.gz\` | macOS | Apple Silicon (M1/M2/M3) |
| \`glyra-darwin-amd64.tar.gz\` | macOS | Intel (x86_64) |
| \`glyra-linux-amd64.tar.gz\` | Linux | x86_64 |
| \`glyra-linux-arm64.tar.gz\` | Linux | ARM64 |
| \`glyra-windows-amd64.zip\` | Windows | x86_64 |

### Install
\`\`\`bash
# macOS / Linux
tar xzf glyra-<platform>.tar.gz
./glyra --version

# Windows
Expand-Archive glyra-windows-amd64.zip .
.\\glyra.exe --version
\`\`\`

### Quick start
\`\`\`bash
glyra init my-app
cd my-app
glyra dev
\`\`\`
EOF

echo "==> committing, tagging, pushing"
git add -A
git commit -m "release: ${TAG}" || true
git tag "$TAG"
git push origin main
git push origin "$TAG"

echo "==> released ${TAG}"
echo "    GitHub Actions will pick up the tag and publish the assets."
