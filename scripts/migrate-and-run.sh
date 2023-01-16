#!/bin/sh

echo $ENVIRONMENT_NAME

./migrations

if [[ $ENVIRONMENT_NAME == "docker" ]]; then
    echo "seeding"
    ./seeder
fi

./server