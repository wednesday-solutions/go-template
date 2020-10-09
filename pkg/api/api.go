// Package api Go-Template
//
// API Docs for GO-Template v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 2.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Host: localhost:9000
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
package api

import (
	"context"
	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"
	goboiler "github.com/wednesday-solutions/go-template"
	"github.com/wednesday-solutions/go-template/pkg/utl"
	"os"

	_ "github.com/lib/pq" // here
	"github.com/volatiletech/sqlboiler/boil"
	graphql "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/pkg/utl/config"
	"github.com/wednesday-solutions/go-template/pkg/utl/jwt"
	authMw "github.com/wednesday-solutions/go-template/pkg/utl/middleware/auth"
	"github.com/wednesday-solutions/go-template/pkg/utl/postgres"
	"github.com/wednesday-solutions/go-template/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.Connect()
	if err != nil {
		return err
	}
	boil.SetDB(db)

	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}

	e := server.New()

	gqlMiddleware := authMw.GqlMiddleware()

	playgroundHandler := playground.Handler("GraphQL playground", "/graphql")
	playgroundHandler = playground.Handler("Subscription playground", "/subscription")

	events := make(chan *graphql.User, 1)
	chatRoom := goboiler.Chatroom{Observers: map[string]struct {
		Username string
		User     chan *graphql.User
	}{}}
	chatRoom.Observers["name"] = struct {
		Username string
		User     chan *graphql.User
	}{Username: "username", User: events}
	rooms := map[string]*goboiler.Chatroom{}
	rooms["room"] = &chatRoom
	graphqlHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &goboiler.Resolver{Rooms: rooms}}))

	if os.Getenv("ENVIRONMENT_NAME") == "local" {
		boil.DebugMode = true
	}

	// graphql apis
	graphqlHandler.AroundOperations(func(ctx context.Context, next graphql2.OperationHandler) graphql2.ResponseHandler {
		return authMw.GraphQLMiddleware(ctx, jwt, next)
	})
	e.POST("/graphql", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		graphqlHandler.ServeHTTP(res, req)
		return nil
	}, gqlMiddleware)

	// graphql playground
	e.GET("/playground", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		playgroundHandler.ServeHTTP(res, req)
		return nil
	})
	// graphql subscription
	e.GET("/subscription", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		utl.Handler(res, req)
		return nil
	})
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})
	return nil
}
