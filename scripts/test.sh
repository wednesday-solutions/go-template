#!/usr/bin/env bash

set -a && source .env.local && set +a 
go test -gcflags=all=-l $(go list ./... | grep -v models | grep -v cmd | grep -v testutls | grep -v gqlmodels | grep -v cmd/seeder)  -coverprofile=coverage.out
