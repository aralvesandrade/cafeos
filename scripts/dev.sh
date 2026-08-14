#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

echo "=== CafeOS Development Environment ==="

case "${1:-help}" in
  up)
    echo "Starting infrastructure (PostgreSQL, Redis, RabbitMQ)..."
    docker compose -f "$ROOT_DIR/docker-compose.yml" up -d
    echo "Waiting for services..."
    sleep 4
    echo "Starting API server (background)..."
    cd "$ROOT_DIR/apps/backend"
    go run ./cmd/api &
    API_PID=$!
    sleep 3
    echo ""
    echo "=== Services ==="
    echo "  API:        http://localhost:5001"
    echo "  Swagger:    http://localhost:5001/swagger/"
    echo "  RabbitMQ:   http://localhost:15672 (cafeos:cafeos)"
    echo ""
    echo "  Admin:      http://localhost:5174  (run: cd apps/admin && npm run dev)"
    echo "  Mobile:     http://localhost:8081  (run: cd apps/mobile && npx expo start --web)"
    echo "  Worker:     cd apps/backend && go run ./cmd/worker/main.go"
    echo ""
    echo "Press Ctrl+C to stop API only."
    wait $API_PID
    ;;
  down)
    echo "Stopping infrastructure..."
    docker compose -f "$ROOT_DIR/docker-compose.yml" down --remove-orphans
    echo "Done."
    ;;
  api)
    echo "Starting API server..."
    cd "$ROOT_DIR/apps/backend"
    go run ./cmd/api
    ;;
  worker)
    echo "Starting RabbitMQ worker..."
    cd "$ROOT_DIR/apps/backend"
    go run ./cmd/worker/main.go
    ;;
  admin)
    echo "Starting admin panel..."
    cd "$ROOT_DIR/apps/admin"
    npm run dev
    ;;
  mobile)
    echo "Starting mobile app (web)..."
    cd "$ROOT_DIR/apps/mobile"
    npx expo start --web
    ;;
  db:migrate)
    echo "Applying GORM AutoMigrate (schema sync only, no seed data)..."
    cd "$ROOT_DIR/apps/backend"
    go run ./cmd/migrate
    ;;
  db:reset)
    echo "Resetting database..."
    docker compose -f "$ROOT_DIR/docker-compose.yml" exec -T postgres psql -U cafeos -d cafeos -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
    cd "$ROOT_DIR/apps/backend"
    go run ./cmd/migrate
    ;;
  db:seed)
    echo "Running seed..."
    cd "$ROOT_DIR/apps/backend"
    go run ./cmd/seed
    ;;
  test)
    echo "Running tests..."
    cd "$ROOT_DIR/apps/backend"
    go test ./... -v
    ;;
  docker:build)
    echo "Building API Docker image..."
    docker buildx build -t aralvesandrade/cafeos-api "$ROOT_DIR/apps/backend"
    ;;
  docker:build:frontend)
    echo "Building frontend Docker image..."
    docker buildx build -f "$ROOT_DIR/apps/frontend/Dockerfile" -t aralvesandrade/cafeos-frontend "$ROOT_DIR"
    ;;
  docker:build:admin)
    echo "Building admin Docker image..."
    docker buildx build -f "$ROOT_DIR/apps/admin/Dockerfile" -t aralvesandrade/cafeos-admin "$ROOT_DIR"
    ;;
  *)
    echo "Usage: $0 {command}"
    echo ""
    echo "  up          Start infra + API (background)"
    echo "  down        Stop all infra containers"
    echo "  api         Start API server (foreground)"
    echo "  worker      Start RabbitMQ sync worker"
    echo "  admin       Start admin panel (React Vite)"
    echo "  mobile      Start mobile app (Expo web)"
    echo "  db:migrate  Run database migrations"
    echo "  db:reset    Reset database and re-run migrations"
    echo "  db:seed     Run seed data"
    echo "  test        Run backend tests"
    echo "  docker:build           Build API Docker image"
    echo "  docker:build:frontend  Build frontend Docker image"
    echo "  docker:build:admin     Build admin Docker image"
    ;;
esac
