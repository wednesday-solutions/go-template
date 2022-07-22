#!/usr/bin/env bash

source ./scripts/source-local.sh

go test -gcflags=all=-l $(go list ./... | grep -v models | grep -v testutls | grep -v gqlmodels | grep -v cmd/seeder)  -coverprofile=coverage.out
