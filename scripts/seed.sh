#!/bin/sh
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