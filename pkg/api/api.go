package api

import (
	"context"
	"net/http"
	"os"
	"time"

	graphql "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/internal/jwt"
	authMw "go-template/internal/middleware/auth"
	"go-template/internal/postgres"
	"go-template/internal/server"
	throttle "go-template/pkg/utl/throttle"
	"go-template/resolver"

	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // here

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Start starts the API service
func Start(cfg *config.Configuration) (*echo.Echo, error) {
	// Initialize Echo instance
	e := server.New()

	// Set up database connection
	if err := setupDatabase(); err != nil {
		return nil, err
	}

	// Set up JWT
	jwt, err := jwt.New(
		cfg.JWT.SigningAlgorithm,
		os.Getenv("JWT_SECRET"),
		cfg.JWT.DurationMinutes,
		cfg.JWT.MinSecretLength)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// Set up GraphQL
	observers := map[string]chan *graphql.User{}
	graphqlHandler := handler.New(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &resolver.Resolver{Observers: observers},
	}))

	graphqlHandler.AroundOperations(func(ctx context.Context, next graphql2.OperationHandler) graphql2.ResponseHandler {
		return authMw.GraphQLMiddleware(ctx, jwt, next)
	})

	graphqlHandler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	graphqlHandler.AddTransport(transport.Options{})
	graphqlHandler.AddTransport(transport.GET{})
	graphqlHandler.AddTransport(transport.POST{})
	graphqlHandler.AddTransport(transport.MultipartForm{})
	graphqlHandler.SetQueryCache(lru.New(1000))
	graphqlHandler.Use(extension.Introspection{})
	graphqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	// Set up GraphQL endpoints
	setupGraphQLEndpoints(e, graphqlHandler)

	// Set up GraphQL playground
	setupGraphQLPlayground(e)

	// Start the server
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return e, nil
}

func setupDatabase() error {
	db, err := postgres.Connect()
	if err != nil {
		return err
	}
	boil.SetDB(db)
	if os.Getenv("ENVIRONMENT_NAME") == "local" {
		boil.DebugMode = true
	}
	return nil
}

func setupGraphQLEndpoints(e *echo.Echo, graphqlHandler *handler.Server) {
	graphQLPathname := "/graphql"
	gqlMiddleware := authMw.GqlMiddleware()
	throttlerMiddleware := throttle.GqlMiddleware()

	e.POST(graphQLPathname, func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		graphqlHandler.ServeHTTP(res, req)
		return nil
	}, gqlMiddleware, throttlerMiddleware)

	e.GET(graphQLPathname, func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		graphqlHandler.ServeHTTP(res, req)
		return nil
	}, gqlMiddleware, throttlerMiddleware)
}

func setupGraphQLPlayground(e *echo.Echo) {
	graphQLPathname := "/graphql"
	playgroundHandler := playground.Handler("GraphQL playground", graphQLPathname)

	e.GET("/playground", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		playgroundHandler.ServeHTTP(res, req)
		return nil
	})
}
