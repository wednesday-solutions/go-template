package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	graphql "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/internal/jwt"
	"github.com/wednesday-solutions/go-template/internal/middleware/auth"
	"github.com/wednesday-solutions/go-template/resolver"
)

func echoHandler(mw ...echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	for _, v := range mw {
		e.Use(v)
	}
	e.GET("/hello", hwHandler)
	return e
}

func hwHandler(c echo.Context) error {
	return c.String(200, "Hello World")
}

func TestGraphQLMiddleware(t *testing.T) {
	cases := map[string]struct {
		wantStatus int
		header     string
		signMethod string
	}{
		"Success": {
			header:     "Bearer 123",
			wantStatus: http.StatusOK,
		},
	}

	err := godotenv.Load(fmt.Sprintf("../../../../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	jWT, err := jwt.New(os.Getenv("JWT_SIGNING_ALGORITHM"), os.Getenv("JWT_SECRET"), 5, 109)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			graphqlHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &resolver.Resolver{}}))
			graphqlHandler.AroundOperations(func(ctx context.Context, next graphql2.OperationHandler) graphql2.ResponseHandler {
				res := auth.GraphQLMiddleware(ctx, jWT, next)
				assert.Equal(t, nil, res)
				return res
			})

			ts := httptest.NewServer(echoHandler())
			defer ts.Close()
			path := ts.URL + "/hello"

			e := echo.New()
			e.POST(path, func(c echo.Context) error {
				req := c.Request()
				res := c.Response()

				graphqlHandler.ServeHTTP(res, req)
				return nil
			})
			req, _ := http.NewRequest("GET", path, nil)
			if tt.wantStatus != 401 {
				req.Header.Set("Authorization", tt.header)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal("Cannot create http request")
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
