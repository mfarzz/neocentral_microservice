.PHONY: proto test-pkg tidy docker-up docker-down docker-logs clean

# ── Proto generation ─────────────────────────
proto:
	@echo "🔧 Generating protobuf Go code..."
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/auth/auth.proto
	@echo "✅ Proto generation complete"

# ── Tidy ──────────────────────────────────────
tidy:
	go mod tidy

# ── Test Shared Packages ──────────────────────
test-pkg:
	go test ./pkg/... -v -count=1

# ── Docker ────────────────────────────────────
docker-up:
	docker compose -f docker-compose.yml up -d --build

docker-down:
	docker compose -f docker-compose.yml down

docker-logs:
	docker compose -f docker-compose.yml logs -f

# ── Clean ─────────────────────────────────────
clean:
	rm -rf bin/
