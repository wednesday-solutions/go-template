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

	graphql "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/internal/middleware/auth"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/resolver"
	testutls "go-template/testutls"

	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/DATA-DOG/go-sqlmock"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const SuccessCase = "Success"

var parseTokenMock func(token string) (*jwt.Token, error)

type tokenParserMock struct {
}

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
		whiteListedQuery bool
	}{
		SuccessCase: {
			whiteListedQuery: false,
			header:           "Bearer 123",
			wantStatus:       http.StatusOK,
			err:              "",
			tokenParser: func(token string) (*jwt.Token, error) {
				return testutls.MockJwt("SUPER_ADMIN"), nil
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
					Query:   `SELECT "users".* FROM "users" WHERE (email=$1) LIMIT 1`,
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
		"Success__WhitelistedQuery": {
			whiteListedQuery: true,
			header:           "bearer 123",
			wantStatus:       http.StatusOK,
			err:              "",
			tokenParser: func(token string) (*jwt.Token, error) {
				// even without mocking the database or the token parser the middleware
				// doesn't throw an error since it skips all the checks and directly calls next
				return nil, nil
			},
			operationHandler: func(ctx context.Context) graphql2.ResponseHandler {
				var handler = func(ctx context.Context) *graphql2.Response {
					return &graphql2.Response{
						Data: json.RawMessage([]byte(`{ "data": { "user": { "id": 1 } } } `)),
					}
				}
				return handler
			},
			dbQueries: []testutls.QueryData{},
		},
		"Failure__NoAuthorizationToken": {
			whiteListedQuery: false,
			header:           "",
			wantStatus:       http.StatusOK,
			err:              "Authorization header is missing",
			tokenParser: func(token string) (*jwt.Token, error) {
				return nil, fmt.Errorf("token is invalid")
			},
			operationHandler: func(ctx context.Context) graphql2.ResponseHandler {
				return nil
			},
			dbQueries: []testutls.QueryData{},
		},
		"Failure__InvalidAuthorizationToken": {
			whiteListedQuery: false,
			header:           "bearer 123",
			wantStatus:       http.StatusOK,
			err:              "Invalid authorization token",
			tokenParser: func(token string) (*jwt.Token, error) {
				return nil, fmt.Errorf("token is invalid")
			},
			operationHandler: func(ctx context.Context) graphql2.ResponseHandler {
				return nil
			},
			dbQueries: []testutls.QueryData{},
		},
		"Failure__NotAnAdmin": {
			whiteListedQuery: false,
			header:           "bearer 123",
			wantStatus:       http.StatusOK,
			err:              "Unauthorized! \n Only admins are authorized to make this request.",
			tokenParser: func(token string) (*jwt.Token, error) {
				return testutls.MockJwt("USER"), nil
			},
			operationHandler: func(ctx context.Context) graphql2.ResponseHandler {
				return nil
			},
			dbQueries: []testutls.QueryData{},
		},
		"Failure__NoUserWithThatEmail": {
			whiteListedQuery: false,
			header:           "bearer 123",
			wantStatus:       http.StatusOK,
			err:              "No user found for this email address",
			tokenParser: func(token string) (*jwt.Token, error) {
				return testutls.MockJwt("SUPER_ADMIN"), nil
			},
			operationHandler: func(ctx context.Context) graphql2.ResponseHandler {
				return nil
			},
			dbQueries: []testutls.QueryData{},
		},
	}

	oldDB := boil.GetDB()
	err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../../../"))
	if err != nil {
		log.Fatal(err)
	}
	mock, db, _ := testutls.SetupMockDB(t)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			for _, dbQuery := range tt.dbQueries {
				mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
					WithArgs(*dbQuery.Actions...).
					WillReturnRows(dbQuery.DbResponse)
			}

			requestQuery := testutls.MockQuery
			if tt.whiteListedQuery {
				requestQuery = testutls.MockWhitelistedQuery
			}
			makeRequest(t, requestQuery, tt)
		})
	}
	boil.SetDB(oldDB)
	db.Close()
}

func makeRequest(t *testing.T, requestQuery string, tt struct {
	wantStatus       int
	header           string
	signMethod       string
	err              string
	dbQueries        []testutls.QueryData
	operationHandler func(ctx context.Context) graphql2.ResponseHandler
	tokenParser      func(token string) (*jwt.Token, error)
	whiteListedQuery bool
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
		bytes.NewBuffer([]byte(requestQuery)),
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

func TestUserIDFromContext(t *testing.T) {
	cases := map[string]struct {
		user   *models.User
		userID int
	}{
		SuccessCase: {
			user:   &models.User{ID: testutls.MockID},
			userID: testutls.MockID,
		},
		"Failure": {
			user:   nil,
			userID: 0,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			userID := auth.UserIDFromContext(context.WithValue(testutls.MockCtx{}, auth.UserCtxKey, tt.user))
			assert.Equal(t, tt.userID, userID)
		})
	}

}

func TestFromContext(t *testing.T) {
	user := &models.User{ID: testutls.MockID}
	u := auth.FromContext(context.WithValue(testutls.MockCtx{}, auth.UserCtxKey, user))
	assert.Equal(t, user, u)
	assert.Equal(t, user.ID, testutls.MockID)
}
