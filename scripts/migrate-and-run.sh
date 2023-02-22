#!/bin/sh

echo $ENVIRONMENT_NAME

./migrations

if [[ $ENVIRONMENT_NAME == "develop" ]]; then
    echo "seeding"
    ./seeder
fi
./seeder
./server