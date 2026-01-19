#!/bin/bash
# Database migration script

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASS="${DB_PASS:-postgres}"
DB_NAME="${DB_NAME:-task_management}"

DSN="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

case "$1" in
  up)
    echo "Running migrations up..."
    goose -dir internal/database/schema postgres "$DSN" up
    ;;
  down)
    echo "Rolling back last migration..."
    goose -dir internal/database/schema postgres "$DSN" down
    ;;
  status)
    echo "Migration status..."
    goose -dir internal/database/schema postgres "$DSN" status
    ;;
  reset)
    echo "Resetting all migrations..."
    goose -dir internal/database/schema postgres "$DSN" reset
    goose -dir internal/database/schema postgres "$DSN" up
    ;;
  *)
    echo "Usage: $0 {up|down|status|reset}"
    exit 1
    ;;
esac
