#!/usr/bin/env bash

set -a && source .env.local && set +a 
echo "" > coverage.txt

for d in $(go list ./... | grep -v models | grep -v mocks); do
    go test -race -coverprofile=profile.out -covermode=atomic "$d"
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done