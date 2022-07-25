#!/bin/sh
set -a
source .env.local
set +a

export PSQL_HOST=localhost
# drop tables
sql-migrate down -env postgres -limit=0

# run migrations
sql-migrate up -env postgres
sql-migrate status -env postgres

# seed data
./scripts/seed.sh