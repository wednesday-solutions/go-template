package resolver_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"testing"

	fm "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/models"
	"go-template/pkg/utl/cnvrttogql"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/rediscache"
	"go-template/resolver"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

func TestMe(
	t *testing.T,
) {
	type args struct {
		user *models.User
	}
	cases := []struct {
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

	err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../"))
	if err != nil {
		log.Fatal(err)
	}
	_, db, _ := testutls.SetupMockDB(t)
	oldDb := boil.GetDB()
	boil.SetDB(db)
	defer func() {
		db.Close()
		boil.SetDB(oldDb)
	}()
	conn := redigomock.NewConn()
	ApplyFunc(
		redis.Dial,
		func(network string, address string, options ...redis.DialOption) (redis.Conn, error) {
			return conn, nil
		},
	)
	//
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {

		if tt.name == ErrorFromRedisCache {
			patchGetUser := gomonkey.ApplyFunc(rediscache.GetUser,
				func(userID int, ctx context.Context) (*models.User, error) {
					return nil, errors.New("redis cache")
				})
			defer patchGetUser.Reset()
		}
		_, db, err := testutls.SetupEnvAndDB(t, testutls.Parameters{EnvFileLocation: `../.env.local`})
		if err != nil {
			panic("failed to setup env and db")
		}
		oldDb := boil.GetDB()
		boil.SetDB(db)
		defer func() {
			db.Close()
			boil.SetDB(oldDb)
		}()
		conn := redigomock.NewConn()
		ApplyFunc(
			redis.Dial,
			func(network string, address string, options ...redis.DialOption) (redis.Conn, error) {
				return conn, nil
			},
		)
		t.Run(
			tt.name,
			func(t *testing.T) {

				b, _ := json.Marshal(tt.args.user)
				conn.Command("GET", "user0").Expect(b)
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, _ := resolver1.Query().Me(ctx)
				if tt.wantResp != nil &&
					response != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}

// TestUsers is a unit test function for testing user queries.
func TestUsers(
	t *testing.T,
) {
	cases := []struct {
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

	// Create a new instance of the resolver.
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {

				// Load environment variables from the .env.local file.
				err := godotenv.Load(
					"../.env.local",
				)
				if err != nil {
					fmt.Print("error loading .env file")
				}

				// Create a new mock database connection.
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				// Inject mock instance into boil.
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)

				//fail on finding user case
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
				// Define a mock result set for user queries.

				// Define a mock result set.

				// Create a new context with a mock user.
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())

				// Query for users using the resolver and get the response and error.
				response, err := resolver1.Query().
					Users(ctx, tt.pagination)

				// Check if the response matches the expected response length.
				if tt.wantResp != nil &&
					response != nil {
					assert.Equal(t, len(tt.wantResp), len(response.Users))

				}
				// Check if the error matches the expected error value.
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}
