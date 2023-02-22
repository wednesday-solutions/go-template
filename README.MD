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

- GraphQL Relay
- Dockerization
- Authorization middleware
- Redis Cache
- Graphql Subscription
- Paginated endpoints
- Simplified support for migrations and seeders
- DAO layer for all database interactions
- Distributed tracing
- Monitoring
- Alerts

## Getting started

Using go-template requires having Go 1.7 or above. Once you download go-template (either using Git or go get) you need to configure the following:

1. Run the `./scripts/setup-pre-commit.sh` script to install the pre-commit hooks locally.

2. Set the ("ENVIRONMENT_NAME") environment variable, either using terminal or os.Setenv("ENVIRONMENT_NAME","dev").

3. Install the sqlboiler, sql-migrate and gqlgen using
```
  go get -v github.com/rubenv/sql-migrate/... \
  github.com/volatiletech/sqlboiler \
  github.com/99designs/gqlgen
  ```

4. To run all the migrations using the script setup-local.sh as follows `./scripts/setup-local.sh`.

5. Generate the graphql models using `gqlgen generate`

6. Run the app using:

```bash
go run cmd/server/main.go
```

**NOTE:** Please do not delete ```.env.base``` file of the project and rebuild the using docker-compose everytime you make changes to it

# Setting up database (postgres)

- Requirement [postgresql](https://www.postgresql.org/)
  
Steps to set up database with ```username``` and ```role``` using terminal

  - Enter postgres terminal ```psql postgres```
  - Create new database ```CREATE DATABASE go_template;```
  - Create a new role with password ```CREATE ROLE go_template_role WITH LOGIN PASSWORD 'go_template_role456';```

**NOTE:** Replace these credentials in ```.env``` file of the project

# Using Docker

To ease the development process a make file is provided

- **`make docker`**  
  Requires `.env.local` file to be present and set  
  This starts the containers in `local` stage, bind the `current directory` to `/go/src/server` inside the `go-template_server_1` container and then starts the `terminal` inside `go-template_server_1`. Once the development is over, `exit` the terminal and call `make tear env=local` to stop all the containers

# Setting up Signoz

Set up signoz locally by following the steps [here](https://signoz.io/docs/install/docker)
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

# File Structure

```txt
go-template/
└──.github/
│  └──workflow/go-template-ci.yml   # this file contains the config of github action
└──cmd/
│  └──seeder/
│  │  └──v1/1_roles.go              # seed file to load roles into DB
│  │  └──v2/2_users.go              # seed file to load users into DB
│  └──server/main.go                # this is the starting point of the go server
└──daos/                            # this directory will hold info about the DB transactions
└──gqlmodels/                       # this directory contain modules for gqlgen and is mostly auto-generated
└──internal/
│  └──config/                       # this package loads env variables into a config object
│  └──jwt/                          # this package has JWT related middlewares and convertors
│  └──middleware/
│     └──auth/
│     └──secure/
│  └──migrations/                   # these are the migrations to be applied
│  └──postgres/                     # this takes care of connecting to postgre
│  └──server/                       # this package have functionality to start a echo server
│  └──services/                     # this will have services used in the server
└──models/
└──pkg/
│  └──api/api.go                    # the starting point of the api
│  └──utl/
│     └──convert/                   # this package has functionality for type conversions
│     └──mock/                      # this package has mock related to passwords and JWTs
│     └──throttle/                  # this package has functionality for request rate throttling
│     └──rediscache/                # this package has functionality for accessing and using redis
│     └──resultwrapper/             # this package exports the custom errors produced by application
│     └──secure/                    # this package has password related functionalities
│     └──zap/                       # this package has setup for uber-go zap logger
└──resolver/                        # this directory will contain resolvers to populate info for graphQL queries, mutations and subscriptions
└──scripts/                         # this directory will contain different utility scripts
│  └──setup-local.sh                # script to set up database and run the app locally
│  └──setup-ecs.sh                  # script to provision the ECS infrastructure in the cloud
│  └──update-ecs.sh                 # script to deploy new version of the app to the provisioned ECS
│  └──setup-pre-commit.sh           # script to setup the pre-commit hooks and the enviornment
│  └──line-formatter.sh             # auto format to adhere to the lll.line-length criteria
└──schema/                          # this directory will have all the .graphql files which make the graphql api
└──.env.local                       # a sample .env file for reference
└──.env.base                        # a base .env file should be included in all environments 
└──.pre-commit-config.yaml          # config to run pre-commit utility
└──dbconfig.yml
└──docker-compose.*.yml             # docker-compose file corresponding to the state of project (local, prod, test)
└──docker-compose.yml               # docker-compose file which serves as a base to other docker-compose files
└──generate-modules.sh              # script to generate modules
└──gqlgen.yml                       # file to configure gqlgen
└──makefile
└──migrate-run.sh                   # script to run DB migrations
└──setup-local.sh                   # a helper script to setup local env (do migration, etc)
└──sqlboiler.toml                   # sqlboiler config file
└──test.sh                          # a helper script to run test in local env
```

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

  Take a look at the following folder

  - [schema](./schema/)
  
## Resolver

- Queries and mutation are autogenerated using gqlgen and are to be extended. Take a look at the following files

  - [resolver](./resolver)


## Infrastructure

### Create infrastructure

### Precautions
1. Please ensure the maximum limit of number of vpc's in a region has not reached it's limit.
2. The maximum limit of number of buckets has not reached it's limits.
3. Please make the changes to the manifest file of the service.
4. Ensure that the aws cli has been installed and configured with the appropriate credentials and profiles.



Application name should container only lowercase letters. No hyphens and underscores or any other special characters.



```
./scripts/setup-ecs.sh gotemplate dev
```


Please change the ENV_INJECTION variable as true in .env.base file to true, so it will not try to find a local .env file

Also add the environment variables to the task,add this block of yml code in ${service name}/manifest.yaml:

```

variables:                   
 ENVIRONMENT_NAME: develop

 #to inject our .env file from aws s3 inside the container

taskdef_overrides:
  - path: ContainerDefinitions[0].EnvironmentFiles[0]
    value:
      type: 's3'
      value: 'arn:aws:s3:::gotemplate-dev-bucket/develop/.env' 
 
```
Make sure that the manifest.yml has http.path: '/'

```
http:
  # Requests to this path will be forwarded to your service.
  # To match all requests you can use the "/" path.
  path: '/'
  # You can specify a custom health check path. The default is "/".
  # healthcheck: '/'

```
Also make sure the execution role has an appropriate policy attached to it so it can access our .env file inside the s3 bucket, and inject it as environment variables.
### To deploy

```
./scripts/deploy-ecs.sh gotemplate dev
```


### Update infrastructure

```
./scripts/update-ecs.sh gotemplate dev
```

## Testing

Get test coverage using
```
go test -cover
```

## Generate mocks

Sample command to generate mocks

```
  mockgen --build_flags=--mod=mod github.com/go-playground/validator  FieldError
```

## Postman Collection

The postman collection can be found [here](postman/collection.json) and has been auto-generated using [graphql-test](https://www.npmjs.com/package/graphql-testkit)
## License

Go Template is licensed under the MIT license. Check the [LICENSE](LICENSE) file for details.
