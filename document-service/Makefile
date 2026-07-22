.PHONY: migrate-up migrate-down run

DB_URL=mysql://root:rootpassword@tcp(localhost:3306)/neocentral_document

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down -all

run:
	go run cmd/main.go
