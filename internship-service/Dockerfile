FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/internship-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/internship-service/main .
COPY --from=builder /app/internship-service/migrations ./migrations

EXPOSE 8005
CMD ["./main"]
