all: lint build run
lint:
	go mod tidy
	golangci-lint run

build:
 	go build -o bin/cli ./cmd/cli
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server
	go run ./cmd/cli