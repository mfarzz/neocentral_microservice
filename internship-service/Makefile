.PHONY: run build migrate-up migrate-down

run:
	go run cmd/main.go

build:
	go build -o main cmd/main.go

migrate-up:
	migrate -path migrations -database "mysql://root:rootpassword@tcp(localhost:3306)/neocentral_internship?charset=utf8mb4&parseTime=True&loc=Local" -verbose up

migrate-down:
	migrate -path migrations -database "mysql://root:rootpassword@tcp(localhost:3306)/neocentral_internship?charset=utf8mb4&parseTime=True&loc=Local" -verbose down
