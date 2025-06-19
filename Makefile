build_id ?= 7460
checkquality: lint build test

lint:
	go mod tidy
	golangci-lint run

build:
	go build -o bin/client ./cmd/client
	go build -o bin/server ./cmd/server

test:
	go test -v ./... -coverprofile cover.out

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client --build-id=${build_id}