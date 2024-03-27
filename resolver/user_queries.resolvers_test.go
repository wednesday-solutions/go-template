package resolver_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"testing"

	"go-template/gqlmodels"
	fm "go-template/gqlmodels"
	"go-template/models"
	"go-template/pkg/utl/cnvrttogql"
	"go-template/pkg/utl/rediscache"
	"go-template/resolver"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

type args struct {
	user *models.User
}

func TestMe(t *testing.T) {
	cases := initializeCases()

	_, cleanup, err := testutls.SetupMockDB(t)
	defer cleanup()
	if err != nil {
		log.Fatal(err)
	}

	resolver1 := &resolver.Resolver{}
	conn := redigomock.NewConn()
	patches := ApplyFunc(
		redis.Dial,
		func(network string, address string, options ...redis.DialOption) (redis.Conn, error) {
			return conn, nil
		},
	)
	defer patches.Reset()

	for _, tt := range cases {
		setupTestEnvironment(tt.name, t)
		runTestCase(t, tt, resolver1)
	}
}

func initializeCases() []struct {
	name     string
	wantResp *fm.User
	wantErr  bool
	args     args
} {
	return []struct {
		name     string
		wantResp *fm.User
		wantErr  bool
		args     args
	}{
		{
			name:     SuccessCase,
			args:     args{user: testutls.MockUser()},
			wantResp: cnvrttogql.UserToGraphQlUser(testutls.MockUser(), 4),
		},
		{
			name:     ErrorFromRedisCache,
			args:     args{user: testutls.MockUser()},
			wantResp: nil,
		},
	}
}

func setupTestEnvironment(name string, t *testing.T) {
	if name == ErrorFromRedisCache {
		patchGetUser := patchRedisCache()
		defer patchGetUser.Reset()
	}
}

func runTestCase(t *testing.T, tt struct {
	name     string
	wantResp *fm.User
	wantErr  bool
	args     args
}, resolver1 *resolver.Resolver) {
	t.Run(
		tt.name,
		func(t *testing.T) {
			testRedisConnection(tt.args.user)

			ctx := setupTestMeContext()

			response, err := resolver1.Query().Me(ctx)
			assertTestMeResponse(t, tt, response, err)
		},
	)
}

func patchRedisCache() *gomonkey.Patches {
	return gomonkey.ApplyFunc(rediscache.GetUser,
		func(userID int, ctx context.Context) (*models.User, error) {
			return nil, errors.New("redis cache")
		})
}

func testRedisConnection(user *models.User) {
	b, _ := json.Marshal(user)
	conn := redigomock.NewConn()
	conn.Command("GET", "user0").Expect(b)
}

func setupTestMeContext() context.Context {
	c := context.Background()
	ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
	return ctx
}

func assertTestMeResponse(t *testing.T, tt struct {
	name     string
	wantResp *fm.User
	wantErr  bool
	args     args
}, response *fm.User, err error) {
	if tt.wantResp != nil && response != nil {
		assert.Equal(t, tt.wantResp, response)
	}
	assert.Equal(t, tt.wantErr, err != nil)
}

// TestUsers is a unit test function for testing user queries.
func initializeTestCases() []struct {
	name       string
	pagination *fm.UserPagination
	wantResp   []*models.User
	wantErr    bool
} {
	return []struct {
		name       string
		pagination *fm.UserPagination
		wantResp   []*models.User
		wantErr    bool
	}{
		{
			name:    ErrorFindingUser,
			wantErr: true,
		},
		{
			name:    "pagination",
			wantErr: false,
			pagination: &fm.UserPagination{
				Limit: 1,
				Page:  1,
			},
			wantResp: testutls.MockUsers(),
		},
		{
			name:     SuccessCase,
			wantErr:  false,
			wantResp: testutls.MockUsers(),
		},
	}
}

func loadEnvVars() error {
	err := godotenv.Load("../.env.local")
	if err != nil {
		return err
	}
	return nil
}

func setExpectations(mock sqlmock.Sqlmock, tt struct {
	name       string
	pagination *fm.UserPagination
	wantResp   []*models.User
	wantErr    bool
}) {
	if tt.name == ErrorFindingUser {
		mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
			WithArgs().
			WillReturnError(fmt.Errorf(""))
	}
	if tt.name == "pagination" {
		rows := sqlmock.
			NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "address"}).
			AddRow(testutls.MockID, testutls.MockEmail, "First", "Last", "+911234567890", "username", "22 Jump Street")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" LIMIT 1 OFFSET 1;`)).WithArgs().WillReturnRows(rows)

		rowCount := sqlmock.NewRows([]string{"count"}).
			AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM "users" LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rowCount)
	} else {
		rows := sqlmock.
			NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "address"}).
			AddRow(testutls.MockID, testutls.MockEmail, "First", "Last", "+911234567890", "username", "22 Jump Street")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users";`)).WithArgs().WillReturnRows(rows)

		rowCount := sqlmock.NewRows([]string{"count"}).
			AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM "users";`)).
			WithArgs().
			WillReturnRows(rowCount)
	}
}

func setupContext() context.Context {
	c := context.Background()
	return context.WithValue(c, testutls.UserKey, testutls.MockUser())
}

func executeQuery(resolver1 *resolver.Resolver,
	ctx context.Context, pagination *fm.UserPagination) (*gqlmodels.UsersPayload, error) {
	return resolver1.Query().Users(ctx, pagination)
}

func assertResponse(t *testing.T, tt struct {
	name       string
	pagination *fm.UserPagination
	wantResp   []*models.User
	wantErr    bool
}, response *gqlmodels.UsersPayload, err error) {
	if tt.wantResp != nil && response != nil {
		assert.Equal(t, len(tt.wantResp), len(response.Users))
	}
	assert.Equal(t, tt.wantErr, err != nil)
}

func TestUsers(
	t *testing.T,
) {
	cases := initializeTestCases()
	resolver1 := resolver.Resolver{}

	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				err := loadEnvVars()
				if err != nil {
					fmt.Print("error loading .env file")
				}

				mock, cleanup, _ := testutls.SetupMockDB(t)

				setExpectations(mock, tt)

				ctx := setupContext()

				response, err := executeQuery(&resolver1, ctx, tt.pagination)

				assertResponse(t, tt, response, err)
				cleanup()
			},
		)
	}
}
