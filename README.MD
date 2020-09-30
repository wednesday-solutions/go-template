<img align="left" src="https://github.com/wednesday-solutions/go-template/blob/master/go_template.png" width="480" height="540" />


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

1. To use go-template as a starting point of a real project whose package name is something like `github.com/author/project`, move the directory `$GOPATH/ github.com/wednesday-solutions/go-template` to `$GOPATH/github.com/author/project`. We use [task](https://github.com/go-task/task) as a task runner. Run `task relocate` to rename the packages, e.g. `task relocate TARGET=github.com/author/project`

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

go-template is licensed under the MIT license. Check the [LICENSE](LICENSE) file for details.

# graphQL

generate all the models
```
sqlboiler psql --no-hooks
```

generate the models

```
gqlgen generate
```

