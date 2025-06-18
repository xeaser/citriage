.PHONY: lint build test

check: lint build test

lint:
	go mod tidy
	golangci-lint run

build:
	go build -o bin/client ./cmd/client
	go build -o bin/server ./cmd/server

test:
	go test -v ./...

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client