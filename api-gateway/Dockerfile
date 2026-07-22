# ============================================
# api-gateway/Dockerfile
# Multi-stage build for minimal production image
# ============================================

FROM golang:alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from root
COPY go.mod go.sum ./
RUN go mod download

# Copy api-gateway source
COPY api-gateway/ ./api-gateway/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /bin/api-gateway ./api-gateway/cmd/main.go

# ── Runtime ───────────────────────────────────
FROM alpine:3.20

RUN apk --no-cache add ca-certificates wget tzdata

# Set timezone
ENV TZ=Asia/Jakarta

COPY --from=builder /bin/api-gateway /usr/local/bin/api-gateway

EXPOSE 8080

ENTRYPOINT ["api-gateway"]
