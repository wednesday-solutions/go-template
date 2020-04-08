#!/bin/bash -x
export DATABASE_URL=postgres://go_boiler_role:go_boiler_role456@localhost:5432/go_boiler?sslmode\=disable
export ENVIRONMENT_NAME=dev
go run cmd/migration/*.go init
go run cmd/migration/*.go up $1
