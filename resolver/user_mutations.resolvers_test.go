package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	fm "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	"github.com/wednesday-solutions/go-template/resolver"
)

func TestCreateUser(t *testing.T) {
	cases := []struct {
		name     string
		req      fm.UserCreateInput
		wantResp *fm.UserPayload
		wantErr  bool
	}{
		{
			name:    "Fail on Create User",
			req:     fm.UserCreateInput{},
			wantErr: true,
		},
		{
			name:     "Success",
			req:      fm.UserCreateInput{FirstName: convert.StringToPointerString("First"), LastName: convert.StringToPointerString("Last"), Username: convert.StringToPointerString("username"), Email: convert.StringToPointerString("mac@wednesday.is")},
			wantResp: &fm.UserPayload{User: &fm.User{ID: "1", FirstName: convert.StringToPointerString("First"), LastName: convert.StringToPointerString("Last"), Username: convert.StringToPointerString("username"), Email: convert.StringToPointerString("mac@wednesday.is")}},
			wantErr:  false,
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
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on Create User" {
				// insert new user
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"users\" (\"first_name\",\"last_name\",\"username\",\"password\",\"email\",\"mobile\",\"phone\",\"address\",\"active\",\"last_login\",\"last_password_change\",\"token\",\"role_id\",\"created_at\",\"updated_at\",\"deleted_at\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// insert new user
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"users\" (\"first_name\",\"last_name\",\"username\",\"password\",\"email\",\"mobile\",\"phone\",\"address\",\"active\",\"last_login\",\"last_password_change\",\"token\",\"role_id\",\"created_at\",\"updated_at\",\"deleted_at\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)")).
				WithArgs().
				WillReturnRows(rows)

			c := context.Background()
			response, err := resolver1.Mutation().CreateUser(c, tt.req)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	cases := []struct {
		name     string
		req      *fm.UserUpdateInput
		wantResp *fm.UserUpdatePayload
		wantErr  bool
	}{
		{
			name:    "Fail on finding User",
			req:     &fm.UserUpdateInput{},
			wantErr: true,
		},
		{
			name:     "Success",
			req:      &fm.UserUpdateInput{FirstName: convert.StringToPointerString("First"), LastName: convert.StringToPointerString("Last"), Address: convert.StringToPointerString("address")},
			wantResp: &fm.UserUpdatePayload{Ok: true},
			wantErr:  false,
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
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on finding User" {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE \"users\" SET \"first_name\"=$1,\"last_name\"=$2,\"username\"=$3,\"password\"=$4,\"email\"=$5,\"mobile\"=$6,\"phone\"=$7,\"address\"=$8,\"active\"=$9,\"last_login\"=$10,\"last_password_change\"=$11,\"token\"=$12,\"role_id\"=$13,\"updated_at\"=$14,\"deleted_at\"=$15 WHERE \"id\"=$16")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// update users with new information
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" SET \"first_name\"=$1,\"last_name\"=$2,\"username\"=$3,\"password\"=$4,\"email\"=$5,\"mobile\"=$6,\"phone\"=$7,\"address\"=$8,\"active\"=$9,\"last_login\"=$10,\"last_password_change\"=$11,\"token\"=$12,\"role_id\"=$13,\"updated_at\"=$14,\"deleted_at\"=$15 WHERE \"id\"=$16")).
				WillReturnResult(result)

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver1.Mutation().UpdateUser(ctx, tt.req)
			if tt.wantResp != nil && response != nil {
				assert.Equal(t, tt.wantResp.Ok, response.Ok)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	cases := []struct {
		name     string
		wantResp *fm.UserDeletePayload
		wantErr  bool
	}{
		{
			name:    "Fail on finding user",
			wantErr: true,
		},
		{
			name:     "Success",
			wantResp: &fm.UserDeletePayload{ID: "0"},
			wantErr:  false,
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
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on finding user" {
				mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// get user by id
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
				WithArgs().
				WillReturnRows(rows)
			// delete user
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM \"users\" WHERE \"id\"=$1")).
				WillReturnResult(result)

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver1.Mutation().DeleteUser(ctx)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
