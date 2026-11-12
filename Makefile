.PHONY: run test build vet

run:
	go run ./cmd/collector

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -o bin/collector ./cmd/collector
