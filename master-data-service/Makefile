.PHONY: build run test migrate-up migrate-down migrate-create swagger

# ── Variables ─────────────────────────────────
DB_USER ?= root
DB_PASSWORD ?= rootpassword
DB_HOST ?= localhost
DB_PORT ?= 3306
DB_NAME ?= neocentral_master

# ── Build ─────────────────────────────────────
build:
	@echo "🏗️  Building master-data-service..."
	go build -o ../bin/master-data-service ./cmd/main.go
	@echo "✅ Build complete → bin/master-data-service"

# ── Run ───────────────────────────────────────
run:
	go run ./cmd/main.go

# ── Test ──────────────────────────────────────
test:
	go test ./... -v -count=1

# ── Database Migrations ──────────────────────
migrate-up:
	migrate -path migrations \
	        -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" \
	        up

migrate-down:
	migrate -path migrations \
	        -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" \
	        down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# ── Swagger ──────────────────────────────────
swagger:
	swag init -g cmd/main.go -o docs --parseDependency --parseInternal
