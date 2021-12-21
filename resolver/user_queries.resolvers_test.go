package resolver_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	fm "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/resolver"
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
			err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
			if err != nil {
				fmt.Print("Error loading .env file")
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
			//rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "phone", "address"}).AddRow(1, "mac@wednesday.is", "First", "Last", "+911234567890", "username", "05943-1123", "22 Jump Street")
			mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
				WithArgs()

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
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
			name:    "Success",
			wantErr: false,
			wantResp: []*models.User{
				{
					FirstName: null.StringFrom("First"),
					LastName:  null.StringFrom("Last"),
					Username:  null.StringFrom("username"),
					Email:     null.StringFrom("mac@wednesday.is"),
					Mobile:    null.StringFrom("+911234567890"),
					Phone:     null.StringFrom("05943-1123"),
					Address:   null.StringFrom("22 Jump Street"),
				},
			},
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
			if err != nil {
				fmt.Print("Error loading .env file")
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
			rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "phone", "address"}).AddRow(1, "mac@wednesday.is", "First", "Last", "+911234567890", "username", "05943-1123", "22 Jump Street")
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\";")).
				WithArgs().
				WillReturnRows(rows)

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver1.Query().Users(ctx, tt.pagination)
			if tt.wantResp != nil && response != nil {
				assert.Equal(t, len(tt.wantResp), len(response.Users))
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
