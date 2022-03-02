package resolver_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	fm "go-template/graphql_models"
	"go-template/models"
	"go-template/resolver"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestMe(t *testing.T) {
	cases := []struct {
		name     string
		wantResp *fm.User
		wantErr  bool
	}{
		{
			name:     "Success",
			wantResp: &fm.User{},
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load("../.env.local")
			if err != nil {
				fmt.Print("error loading .env file")
			}
			conn := redigomock.NewConn()
			_ = &redis.Pool{
				// Return the same connection mock for each Get() call.
				Dial:    func() (redis.Conn, error) { return conn, nil },
				MaxIdle: 10,
			}
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

			// get user by id
			mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
				WithArgs()

			c := context.Background()
			ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
			response, _ := resolver1.Query().Me(ctx)
			assert.Equal(t, tt.wantResp, response)
		})
	}
}

func TestUsers(t *testing.T) {
	cases := []struct {
		name       string
		pagination *fm.UserPagination
		wantResp   []*models.User
		wantErr    bool
	}{
		{
			name:     "Success",
			wantErr:  false,
			wantResp: testutls.MockUsers(),
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load("../.env.local")
			if err != nil {
				fmt.Print("error loading .env file")
			}
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

			// get user by id
			rows := sqlmock.
				NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "address"}).
				AddRow(
					testutls.MockID,
					testutls.MockEmail,
					"First",
					"Last",
					"+911234567890",
					"username",
					"22 Jump Street",
				)
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\";")).
				WithArgs().
				WillReturnRows(rows)
			rowCount := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM \"users\";")).
				WithArgs().
				WillReturnRows(rowCount)

			c := context.Background()
			ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
			response, err := resolver1.Query().Users(ctx, tt.pagination)

			if tt.wantResp != nil && response != nil {
				assert.Equal(t, len(tt.wantResp), len(response.Users))
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
