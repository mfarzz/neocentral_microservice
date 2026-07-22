# ============================================
# master-data-service/Dockerfile
# Multi-stage build for minimal production image
# ============================================

FROM golang:alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from root
COPY go.mod go.sum ./
RUN go mod download

# Copy shared packages & proto
COPY pkg/ ./pkg/
COPY proto/ ./proto/

# Copy master-data-service source
COPY master-data-service/ ./master-data-service/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /bin/master-data-service ./master-data-service/cmd/main.go

# ── Runtime ───────────────────────────────────
FROM alpine:3.20

RUN apk --no-cache add ca-certificates wget tzdata

# Set timezone
ENV TZ=Asia/Jakarta

COPY --from=builder /bin/master-data-service /usr/local/bin/master-data-service

# Copy migrations (for running via migrate CLI in container)
COPY master-data-service/migrations /migrations

EXPOSE 8002

ENTRYPOINT ["master-data-service"]
