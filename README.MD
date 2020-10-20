<img align="left" src="https://github.com/wednesday-solutions/go-template/blob/master/golang_template_github.svg" width="440" height="480" />

<div>
  <a href="https://www.wednesday.is?utm_source=gthb&utm_medium=repo&utm_campaign=serverless" align="left" style="margin-left: 0;">
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
