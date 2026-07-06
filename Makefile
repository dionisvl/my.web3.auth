.PHONY: up run build test

up:
	docker compose down -t 0 && docker compose up --build

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...
