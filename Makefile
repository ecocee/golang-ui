.PHONY: build test lint race clean

build:
	go build -o bin/ ./cmd/golang-ui

test:
	go test -count=1 ./...

race:
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/
	go clean -testcache
