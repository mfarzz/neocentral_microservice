# ============================================
# auth-service/Dockerfile
# Multi-stage build for minimal production image
# ============================================

FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from root
COPY go.mod go.sum ./
RUN go mod download

# Copy shared packages & proto
COPY pkg/ ./pkg/
COPY proto/ ./proto/

# Copy auth-service source
COPY auth-service/ ./auth-service/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /bin/auth-service ./auth-service/cmd/main.go

# ── Runtime ───────────────────────────────────
FROM alpine:3.20

RUN apk --no-cache add ca-certificates wget tzdata

# Set timezone
ENV TZ=Asia/Jakarta

COPY --from=builder /bin/auth-service /usr/local/bin/auth-service

# Copy migrations (for running via migrate CLI in container)
COPY auth-service/migrations /migrations

EXPOSE 8001 9001

ENTRYPOINT ["auth-service"]
