#!/bin/bash
# scripts/migrate.sh

# Load environment variables if not set
DB_URL=${DATABASE_URL:-"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"}
MIGRATIONS_PATH="/workspaces/waypoint/migrations"

command=$1
shift # remove the first argument (command)

case $command in
  "up")
    migrate -path $MIGRATIONS_PATH -database "$DB_URL" up "$@"
    ;;
  "down")
    migrate -path $MIGRATIONS_PATH -database "$DB_URL" down "$@"
    ;;
  "create")
    migrate create -ext sql -dir $MIGRATIONS_PATH -seq "$@"
    ;;
  *)
    echo "Usage: ./scripts/migrate.sh {up|down|create} [args]"
    exit 1
    ;;
esac