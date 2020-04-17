# Go boiler

GO boiler is a starter project with support for both REST as well as graphql   

## Getting started

Using go-boiler requires having Go 1.7 or above. Once you download go-boiler (either using Git or go get) you need to configure the following:

1. To use go-boiler as a starting point of a real project whose package name is something like `github.com/author/project`, move the directory `$GOPATH/ github.com/wednesday-solutions/go-boiler` to `$GOPATH/github.com/author/project`. We use [task](https://github.com/go-task/task) as a task runner. Run `task relocate` to rename the packages, e.g. `task relocate TARGET=github.com/author/project`

2. Change the configuration file according to your needs, or create a new one.

3. Set the ("ENVIRONMENT_NAME") environment variable, either using terminal or os.Setenv("ENVIRONMENT_NAME","dev").

4. Set the JWT secret env var ("JWT_SECRET")

5. In cmd/migration/main.go set up psn variable and then run it (go run main.go). It will create all tables, and necessery data, with a new account username/password admin/admin.

6. To run all the migrations using the script setup-local.sh as follows `./setup-local.sh 4`, the 4 represent the version that you need to migrate to.

7. Generate the graphql models using ./generate-models.sh

8. Run the app using:

```bash
go run cmd/api/main.go
```


## License

go-boiler is licensed under the MIT license. Check the [LICENSE](LICENSE) file for details.

## Author

[Mohammed Ali Chherawalla](https://github.com/aliChherawalla)

# graphQL

generate all the models
```
sqlboiler psql --no-hooks
```
create the schema from the sqlboiler models

```
go run github.com/web-ridge/sqlboiler-graphql-schema
```

create the models and resolvers

```
go run cmd/gqlgen/main.go
```

create resolvers with mapping between sqlboiler models to graphql models

```
 go run cmd/convert/convert_plugin.go
```
