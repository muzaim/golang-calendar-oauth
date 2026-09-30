include .env
export

DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

.PHONY: dev build run migrate-create migrate-up migrate-down migrate-force docker-up docker-down test

# Development with Air hot reload
dev:
	air

# Build app binary
build:
	go build -o bin/server cmd/server/main.go

# Run app binary
run: build
	./bin/server

# Migrations using golang-migrate
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir database/migrations -seq $$name

migrate-up:
	migrate -path database/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path database/migrations -database "$(DB_URL)" down 1

migrate-force:
	@read -p "Enter version: " version; \
	migrate -path database/migrations -database "$(DB_URL)" force $$version

# Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Run unit tests
test:
	go test -v ./...
