#!/usr/bin/env bash

set -a && source .env.local && set +a 
go test $(go list ./... | grep -v models | grep -v testutls | grep -v graphql_models | grep -v cmd/seeder)  -coverprofile=coverage.out