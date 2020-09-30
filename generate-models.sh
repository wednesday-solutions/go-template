rm resolver.go
rm schema.graphql
rm -rf models
sqlboiler psql --no-hooks
gqlgen generate