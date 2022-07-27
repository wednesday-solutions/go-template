#!/bin/sh

go run /app/cmd/migrations/main.go

if [[ $ENVIRONMENT_NAME == "docker" ]]; then
    echo "seeding"
    /app/scripts/seed.sh
fi

/app/main