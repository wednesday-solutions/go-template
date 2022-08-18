#!/bin/sh

echo $ENVIRONMENT_NAME

go run ./cmd/migrations/main.go

if [[ $ENVIRONMENT_NAME == "docker" ]]; then
    echo "seeding"
    go run ./cmd/seeder/main.go
fi



./main