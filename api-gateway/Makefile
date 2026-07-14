.PHONY: build run clean

# ── Build ─────────────────────────────────────
build:
	@echo "🏗️  Building api-gateway..."
	go build -o ../bin/api-gateway ./cmd/main.go
	@echo "✅ Build complete → bin/api-gateway"

# ── Run ───────────────────────────────────────
run:
	go run ./cmd/main.go

# ── Clean ─────────────────────────────────────
clean:
	rm -rf ../bin/api-gateway
