.PHONY: build run migrate-up migrate-down

DB_URL ?= "mysql://root:rootpassword@tcp(localhost:3306)/neocentral_notification"

build:
	go build -o bin/notification-service cmd/main.go

run: build
	./bin/notification-service

migrate-up:
	migrate -path migrations -database $(DB_URL) up

migrate-down:
	migrate -path migrations -database $(DB_URL) down
