rm resolver.go
rm schema.graphql
rm -rf models
sqlboiler psql --no-hooks
go get -t github.com/web-ridge/sqlboiler-graphql-schema
sqlboiler-graphql-schema --pagination=offset
go run cmd/gqlgen/main.go
go run cmd/convert/convert_plugin.go
