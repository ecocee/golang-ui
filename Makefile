.PHONY: build build-single install test lint race clean release

VERSION ?= dev

# Delegate to scripts/ so CI, local builds, and the release flow share one
# source of truth. See scripts/build.sh and scripts/install.sh for the
# full cross-compilation and install logic.

## build all platforms into ./release (default: dev build)
build:
	bash scripts/build.sh "$(VERSION)"

## build for the current host platform only
build-single:
	bash scripts/build.sh --single "$(VERSION)"

## install the CLI for the current host into $(GOPATH)/bin
install:
	bash scripts/install.sh "$(VERSION)"

## run tests
test:
	go test -count=1 ./...

## run tests with the race detector
race:
	go test -race -count=1 ./...

## lint (requires golangci-lint)
lint:
	golangci-lint run ./...

## cut a new release: bump version, build, tag, push (use with care)
release:
	@test -n "$(TAG)" && bash scripts/release.sh "$(TAG)" || \
		(echo "usage: make release TAG=v0.2.0-beta.1" >&2; exit 1)

## remove build artifacts
clean:
	rm -rf bin/ dist/ release/ RELEASE_NOTES.md
	go clean -testcache
