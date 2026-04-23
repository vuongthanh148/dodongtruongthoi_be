.PHONY: run dev tidy build

run:
	go run ./cmd/server

dev:
	go run ./cmd/server

tidy:
	go mod tidy

build:
	go build -o bin/server ./cmd/server
