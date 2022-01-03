package auth_test

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/boil"
	graphql "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/internal/middleware/auth"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/resolver"
	testutls "github.com/wednesday-solutions/go-template/testutls"
)

type tokenParserMock struct {
}

var parseTokenMock func(token string) (*jwt.Token, error)

func (s tokenParserMock) ParseToken(token string) (*jwt.Token, error) {
	return parseTokenMock(token)
}

var operationHandlerMock func(ctx context.Context) graphql2.ResponseHandler

func TestGraphQLMiddleware(t *testing.T) {
	cases := map[string]struct {
		wantStatus       int
		header           string
		signMethod       string
		err              string
		dbQueries        []testutls.QueryData
		operationHandler func(ctx context.Context) graphql2.ResponseHandler
		tokenParser      func(token string) (*jwt.Token, error)
	}{
		"Success": {
			header:     "Bearer 123",
			wantStatus: http.StatusOK,
			err:        "",
			tokenParser: func(token string) (*jwt.Token, error) {
				return testutls.MockJwt(), nil
			},
			operationHandler: func(ctx context.Context) graphql2.ResponseHandler {
				user := ctx.Value(auth.UserCtxKey).(*models.User)

				// add your assertions here
				assert.Equal(t, testutls.MockEmail, user.Email.String)
				assert.Equal(t, testutls.MockID, user.ID)
				assert.Equal(t, testutls.MockToken, user.Token.String)

				// if you want a custom response you can add it here
				var handler = func(ctx context.Context) *graphql2.Response {
					return &graphql2.Response{
						Data: json.RawMessage([]byte("{}")),
					}
				}
				return handler
			},
			dbQueries: []testutls.QueryData{
				{
					Actions: &[]driver.Value{testutls.MockEmail},
					Query:   "SELECT * FROM \"users\" WHERE (email=$1) LIMIT 1",
					DbResponse: sqlmock.NewRows([]string{
						"id", "email", "token",
					}).AddRow(
						testutls.MockID,
						testutls.MockEmail,
						testutls.MockToken,
					),
				},
			},
		},
		"Failure__NoAuthorizationToken": {
			header:     "",
			wantStatus: http.StatusOK,
			err:        "Authorization header is missing",
			tokenParser: func(token string) (*jwt.Token, error) {
				return nil, fmt.Errorf("token is invalid")
			},
			operationHandler: func(ctx context.Context) graphql2.ResponseHandler {
				user := ctx.Value(auth.UserCtxKey).(*models.User)

				// add your assertions here
				assert.Equal(t, testutls.MockEmail, user.Email.String)
				assert.Equal(t, testutls.MockID, user.ID)
				assert.Equal(t, testutls.MockToken, user.Token.String)

				// if you want a custom response you can add it here
				var handler = func(ctx context.Context) *graphql2.Response {
					return &graphql2.Response{
						Data: json.RawMessage([]byte("{}")),
					}
				}
				return handler
			},
			dbQueries: []testutls.QueryData{},
		},
		"Failure__InvalidAuthorizationToken": {
			header:     "bearer 123",
			wantStatus: http.StatusOK,
			err:        "Invalid authorization token",
			tokenParser: func(token string) (*jwt.Token, error) {
				return nil, fmt.Errorf("token is invalid")
			},
			operationHandler: func(ctx context.Context) graphql2.ResponseHandler {
				user := ctx.Value(auth.UserCtxKey).(*models.User)

				// add your assertions here
				assert.Equal(t, testutls.MockEmail, user.Email.String)
				assert.Equal(t, testutls.MockID, user.ID)
				assert.Equal(t, testutls.MockToken, user.Token.String)

				// if you want a custom response you can add it here
				var handler = func(ctx context.Context) *graphql2.Response {
					return &graphql2.Response{
						Data: json.RawMessage([]byte("{}")),
					}
				}
				return handler
			},
			dbQueries: []testutls.QueryData{},
		},
	}

	oldDB := boil.GetDB()
	mock, db, _ := testutls.SetupEnvAndDB(t, testutls.Parameters{EnvFileLocation: "../../../.env.local"})

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			for _, dbQuery := range tt.dbQueries {
				mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
					WithArgs(*dbQuery.Actions...).
					WillReturnRows(dbQuery.DbResponse)
			}
			makeRequest(t, tt)
		})
	}
	boil.SetDB(oldDB)
	db.Close()
}

func makeRequest(t *testing.T, tt struct {
	wantStatus       int
	header           string
	signMethod       string
	err              string
	dbQueries        []testutls.QueryData
	operationHandler func(ctx context.Context) graphql2.ResponseHandler
	tokenParser      func(token string) (*jwt.Token, error)
}) {
	// mock token parser to handle the different cases for when the token us valid, invalid, empty
	parseTokenMock = tt.tokenParser

	// mock operation handler, and assert different conditions
	operationHandlerMock = tt.operationHandler

	tokenParser := tokenParserMock{}
	client := &http.Client{}
	observers := map[string]chan *graphql.User{}
	graphqlHandler := handler.New(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &resolver.Resolver{Observers: observers},
	}))

	graphqlHandler.
		AroundOperations(func(ctx context.Context, next graphql2.OperationHandler) graphql2.ResponseHandler {
			res := auth.GraphQLMiddleware(ctx, tokenParser, operationHandlerMock)
			return res
		})

	graphqlHandler.AddTransport(transport.POST{})
	pathName := "/graphql"

	e := echo.New()
	e.POST(pathName, func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		graphqlHandler.ServeHTTP(res, req)
		return nil
	}, auth.GqlMiddleware())

	ts := httptest.NewServer(e)
	path := ts.URL + pathName
	defer ts.Close()

	req, _ := http.NewRequest(
		"POST",
		path,
		bytes.NewBuffer([]byte(`{"query":"query Me {me{id}}","variables":{}}`)),
	)

	if tt.wantStatus != 401 {
		req.Header.Set("authorization", tt.header)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Fatal("Cannot create http request")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var jsonRes graphql2.Response
	err = json.Unmarshal(bodyBytes, &jsonRes)

	if err != nil {
		log.Fatal(err)
	}
	for _, errorString := range jsonRes.Errors {
		assert.Equal(t, tt.err, errorString.Message)
	}
	assert.Equal(t, tt.wantStatus, res.StatusCode)
}
