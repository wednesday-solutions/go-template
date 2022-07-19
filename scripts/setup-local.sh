#!/bin/bash
set -a
source .env.local
set +a

export DB_HOST=localhost
# drop tables
sql-migrate down -env postgres -limit=0

# run migrations
sql-migrate up -env postgres
sql-migrate status -env postgres

# seed data

# shellcheck disable=SC2164
cd ./cmd/seeder/

# shellcheck disable=SC2207
seeders=($(ls -d ./*))
for i in "${seeders[@]}"
do
  # shellcheck disable=SC2207
  file=($(ls -d "$i"/*))
   :
   # shellcheck disable=SC2128
   go run "$file"
done
cd  ../../