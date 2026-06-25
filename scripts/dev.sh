#!/bin/bash
set -euo pipefail

echo "=== CafeOS Development Environment ==="

case "${1:-help}" in
  up)
    echo "Starting infrastructure..."
    docker compose -f infra/dev/docker-compose.yml up -d
    echo "Waiting for PostgreSQL..."
    sleep 3
    echo "Running migrations..."
    docker compose -f infra/dev/docker-compose.yml exec -T postgres psql -U cafeos -d cafeos -f /docker-entrypoint-initdb.d/001_initial_schema.sql 2>/dev/null || true
    echo "Starting API server..."
    cd apps/backend
    go run ./cmd/api
    ;;
  down)
    echo "Stopping infrastructure..."
    docker compose -f infra/dev/docker-compose.yml down
    ;;
  db:migrate)
    echo "Running migrations..."
    docker compose -f infra/dev/docker-compose.yml exec -T postgres psql -U cafeos -d cafeos -f /docker-entrypoint-initdb.d/001_initial_schema.sql
    ;;
  db:reset)
    echo "Resetting database..."
    docker compose -f infra/dev/docker-compose.yml exec -T postgres psql -U cafeos -d cafeos -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
    docker compose -f infra/dev/docker-compose.yml exec -T postgres psql -U cafeos -d cafeos -f /docker-entrypoint-initdb.d/001_initial_schema.sql
    ;;
  test)
    echo "Running tests..."
    cd apps/backend
    go test ./... -v
    ;;
  *)
    echo "Usage: $0 {up|down|db:migrate|db:reset|test}"
    echo ""
    echo "  up          Start infrastructure and API server"
    echo "  down        Stop infrastructure"
    echo "  db:migrate  Run database migrations"
    echo "  db:reset    Reset database and re-run migrations"
    echo "  test        Run all tests"
    ;;
esac
