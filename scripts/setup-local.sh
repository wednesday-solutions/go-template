#!/bin/bash
set -a
source .env.local
set +a

export MYSQL_HOST=localhost
# drop tables
sql-migrate down -env mysql -limit=0

# run migrations
sql-migrate up -env mysql
sql-migrate status -env mysql

# seed data

./scripts/seed.sh