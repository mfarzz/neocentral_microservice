FROM golang:alpine AS builder

WORKDIR /app

# Copy the monorepo go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire workspace
COPY . .

# Build the service
WORKDIR /app/document-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o document-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/document-service/document-service .

EXPOSE 8003
CMD ["./document-service"]
