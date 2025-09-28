SHELL := /bin/bash

BIN := transactions-api
PKG := github.com/atalkhandelwal/transactions-api
MAIN := ./cmd/server

.PHONY: all run build test docker up down logs

all: build

build:
	GO111MODULE=on CGO_ENABLED=0 go build -o bin/$(BIN) $(MAIN)

run:
	DB_HOST=${DB_HOST:-localhost} DB_PORT=${DB_PORT:-5432} DB_USER=${DB_USER:-postgres} DB_PASSWORD=${DB_PASSWORD:-postgres} DB_NAME=${DB_NAME:-transactions} go run $(MAIN)

test:
	go test ./... -cover

docker:
	docker build -t $(BIN):local .

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f
