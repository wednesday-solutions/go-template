<img align="left" src="https://github.com/wednesday-solutions/go-template/blob/master/golang.svg" width="480" height="540" />


<div>
  <a href="https://www.wednesday.is?utm_source=gthb&utm_medium=repo&utm_campaign=serverless" align="left" style="margin-left: 0;">
    <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f5879492fafecdb3e5b0e75_wednesday_logo.svg">
  </a>
  <p>
    <h1 align="left">Go Template
    </h1>
  </p>
  <p>
An enterprise go template application showcasing - Testing strategies, Global state management, middleware support, and Continuous integration.
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


# Go boiler

GO boiler is a starter project with support for both REST as well as graphql   

## Getting started

Using go-template requires having Go 1.7 or above. Once you download go-template (either using Git or go get) you need to configure the following:


1. Set the ("ENVIRONMENT_NAME") environment variable, either using terminal or os.Setenv("ENVIRONMENT_NAME","dev").

2. Set the JWT secret env var ("JWT_SECRET")

3. In cmd/migration/main.go set up psn variable and then run it (go run main.go). It will create all tables, and necessery data, with a new account username/password admin/admin.

4. To run all the migrations using the script setup-local.sh as follows `./setup-local.sh`.

5. Generate the graphql models using ./generate-models.sh

6. Run the app using:

```bash
go run cmd/api/main.go
```

# graphQL

generate a Go ORM tailored to your database schema
```
sqlboiler psql --no-hooks
```

generate the models

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

go-template is licensed under the MIT license. Check the [LICENSE](LICENSE) file for details.
