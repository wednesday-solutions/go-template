#!/bin/bash -x
export DATABASE_URL=postgres://go_boiler_role:go_boiler_role456@localhost:5432/go_boiler?sslmode\=disable
export ENVIRONMENT_NAME=dev
export JWT_SECRET=aa60b52e-7974-11ea-bc55-0242ac130003b08d6654-7974-11ea-bc55-0242ac130003-b08d6654-7974-11ea-bc55-0242ac130003

go run cmd/migration/*.go init
go run cmd/migration/*.go up $1
