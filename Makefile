.PHONY: build test lint race clean

build:
	go build -o glyra ./cmd/glyra

test:
	go test -count=1 ./...

race:
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/
	go clean -testcache
