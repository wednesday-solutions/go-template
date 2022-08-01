#!/bin/sh
set -a && source .env.local && set +a

export PSQL_HOST=localhost
# drop first
go run ./cmd/migrations/main.go down

# run migrations
go run ./cmd/migrations/main.go

# seed data
go run ./cmd/seeder/main.go