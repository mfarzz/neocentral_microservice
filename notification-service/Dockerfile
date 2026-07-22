FROM golang:alpine AS builder

WORKDIR /app

# Copy the monorepo go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire workspace
COPY . .

# Build the service
WORKDIR /app/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/notification-service/notification-service .
COPY --from=builder /app/notification-service/migrations ./migrations

EXPOSE 8004
CMD ["./notification-service"]
