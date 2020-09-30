#!/bin/bash
set -a
source .env.local
set +a

# drop tables
go run cmd/migrations/*.go reset

# run migrations
go run cmd/migrations/*.go init
go run cmd/migrations/*.go

# seed data

# shellcheck disable=SC2164
cd ./cmd/seeders/

# shellcheck disable=SC2207
seeders=($(ls -d ./*))
for i in "${seeders[@]}"
do
  file=($(ls -d $i/*))
   :
   go run "$file"
done
cd  ../../

./generate-models.sh