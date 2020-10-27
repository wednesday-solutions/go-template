<img align="left" src="https://github.com/wednesday-solutions/go-template/blob/master/golang_template_github.svg" width="440" height="480" />

<div>
  <a href="https://www.wednesday.is?utm_source=gthb&utm_medium=repo&utm_campaign=gotemplate" align="left" style="margin-left: 0;">
    <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f5879492fafecdb3e5b0e75_wednesday_logo.svg">
  </a>
  <p>
    <h1 align="left">Go Template
    </h1>
  </p>
  <p>
An enterprise go template application showcasing - Testing strategies, middleware support, and Continuous integration.
  </p>
  
  ___


  <p>
    <h4>
      Expert teams of digital product strategists, developers, and designers.
    </h4>
  </p>

  <div>
    <a href="https://www.wednesday.is/contact-us?utm_source=gthb&utm_medium=repo&utm_campaign=serverless" target="_blank">
      <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f6ae88b9005f9ed382fb2a5_button_get_in_touch.svg" width="121" height="34">
    </a>
    <a href="https://github.com/wednesday-solutions/" target="_blank">
      <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f6ae88bb1958c3253756c39_button_follow_on_github.svg" width="168" height="34">
    </a>
  </div>

  ___

  <span>We’re always looking for people who value their work, so come and join us. <a href="https://www.wednesday.is/hiring">We are hiring!</a></span>


</div>

---

<br/>





The Go Template is a template/starter go project.

## Out of the box support for

- GraphQL
- Dockerization
- Authorization middleware
- Redis Cache
- Graphql Subscription
- Paginated endpoints
- Simplified support for migrations and seeders
- DAO layer for all database interactions

## Getting started

Using go-template requires having Go 1.7 or above. Once you download go-template (either using Git or go get) you need to configure the following:

1. Set the ("ENVIRONMENT_NAME") environment variable, either using terminal or os.Setenv("ENVIRONMENT_NAME","dev").

2. To run all the migrations using the script setup-local.sh as follows `./setup-local.sh`.

3. Generate the graphql models using `./generate-models.sh`

4. Run the app using:

```bash
go run cmd/server/main.go
```

# Setting up database (postgres)

- Requirement [postgresql](https://www.postgresql.org/)
  
Steps to set up database with ```username``` and ```role``` using terminal

  - Enter postgres terminal ```psql postgres```
  - Create new database ```CREATE DATABASE go_template;```
  - Create a new role with password ```CREATE ROLE go_template_role WITH LOGIN PASSWORD 'go_template_role456';```

**NOTE:** Replace these credentials in ```.env``` file of the project

# Using Docker

- Requirement [Docker-Compose](https://docs.docker.com/compose/install/)

You can chose to use Docker-compose for `local`, `prod` and `test` stages. This would take care of all the dependencies, i.e., No need for locally having `Golang`, `PostgreSQL`, `redis` and other `cli` tools used in the template.

There will be 3 containers produced,

- **`go-template_db_1`** this container is based on `postgres:13` image and is referred as `db` when connecting from the `server` source code.
- **`go-template_redis_1`** this container is based on `redis:6-alpine` image and is referred as `redis` when connecting from the `server` source code.
- **`go-template_server_1`** this container is where your source code lives and is based on `golang:1.15.2-buster` image (in `local` and `test` stage) or `scratch` image (in `prod` stage).

When in `local` stage there will be 1 more container,

- **`go-template_pgAdmin_1`** this container is based on `dpage/pgadmin4:4` and requires `PGADMIN_PORT`, `PGADMIN_EMAIL` and `PGADMIN_PASS` env variables to work. To connect to `go-template_db_1` from `pgadmin` put `db` as **Host name/ Address**

To ease the development process a make file is provided with 5 commands

- **`make local`**
  Requires `.env.local` file to be present and set
  This starts the containers in `local` stage, bind the `current directory` to `/go/src/server` inside the `go-template_server_1` container and then starts the `terminal` inside `go-template_server_1`. Once the development is over, `exit` the terminal and call `make tear env=local` to stop all the containers

- **`make test`**
  Requires `.env.test` file to be present and set
  This starts the containers in `test` stage and runs `go test -v ./...` inside `go-template_server_1` container. This can be used for `CI/CD` pipelines too. Remember to call `make tear env=test` to stop the containers.

- **`make prod`**
  Requires `.env.prod` file to be present and set
  This starts the containers in `prod` stage. You cannot connect to `go-template_server_1` because it has no terminal the contain is build from `scratch` image with only the `go binary` running in it. To stop the containers use `make tear env=prod` and to get live logs from the containers use `make logs env=prod`.

- **`make logs env=<stage>`**
  this command gives the logs of all the containers in the given stage, useful because the containers are started in detached mode (`-d`)

- **`make tear env=<stage>`**
  this command tears down (`docker-compose down`) all 3 running containers (`db`, `redis`, `server`).

# Running migrations

Migrations are present in ```internal/migrations``` package. Run below command to run all migrations at once:
```
sql-migrate up -env postgres
```
To drop migration use following
```
sql-migrate down -env postgres -limit=0
```
To check status of migration use following
```
sql-migrate new -env postgres <name>
```
To add new migration use following, it creates a new empty migration template with pattern `<current time>-<name>.sql`
```
sql-migrate new -env postgres <name>
```
append query to above file


For more information on migration package refer [here](https://github.com/rubenv/sql-migrate)

# DB Models

generate your database models
```
sqlboiler psql --no-hooks
```

# graphQL

generate the graphql models from the database schema
```
gqlgen generate
```

## API (for graphQL to operate)

- Graphql endpoint ```POST``` request ```/graphql```

- Playground endpoint for schema ```/playground```

  Take a look at the following file

  - [pkg/api/api.go](pkg/api/api.go)
  
## Schema

- Schema can generated or altered manually

  Take a look at the following file

  - [schema.graphql](schema.graphql)
  
## Resolver

- Queries and mutation can be implemented as per convenience

  Take a look at the following file

  - [resolver.go](resolver.go)

## License

Go Template is licensed under the MIT license. Check the [LICENSE](LICENSE) file for details.
