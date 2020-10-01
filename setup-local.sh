#!/bin/bash
set -a
source .env.local
set +a

# drop tables
go run cmd/migration/*.go reset

# run migrations
go run cmd/migration/*.go init
go run cmd/migration/*.go

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

./generate-models.sh