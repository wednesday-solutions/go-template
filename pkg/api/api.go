// Package api Go-Template
package api

import (
	"context"
	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	goboiler "github.com/wednesday-solutions/go-template"
	"net/http"
	"os"
	"time"

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

	observers := map[string]chan *graphql.User{}
	graphqlHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &goboiler.Resolver{Observers: observers},
	}))

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

	subscriptionHandler := handler.New(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &goboiler.Resolver{Observers: observers},
	}))
	subscriptionHandler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 1 * time.Second,
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			return ctx, nil
		},
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	subscriptionHandler.Use(extension.Introspection{})
	e.GET("/graphql", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		subscriptionHandler.ServeHTTP(res, req)
		return nil
	}, gqlMiddleware)

	// graphql playground
	e.GET("/playground", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		playgroundHandler.ServeHTTP(res, req)
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
