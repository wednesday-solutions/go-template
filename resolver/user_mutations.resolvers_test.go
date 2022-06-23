package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	fm "go-template/gqlmodels"
	"go-template/pkg/utl/convert"
	"go-template/resolver"
	"go-template/testutls"
	"regexp"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
func TestCreateUser(t *testing.T) {
	cases := []struct {
		name     string
		req      fm.UserCreateInput
		wantResp *fm.User
		wantErr  bool
	}{
		{
			name:    "Fail on Create User",
			req:     fm.UserCreateInput{},
			wantErr: true,
		},
		{
			name: "Success",
			req: fm.UserCreateInput{
				FirstName: testutls.MockUser().FirstName.String,
				LastName:  testutls.MockUser().LastName.String,
				Username:  testutls.MockUser().Username.String,
				Email:     testutls.MockUser().Email.String,
				RoleID:    fmt.Sprint(testutls.MockUser().RoleID.Int),
			},
			wantResp: &fm.User{
				ID:                 fmt.Sprint(testutls.MockUser().ID),
				Email:              convert.NullDotStringToPointerString(testutls.MockUser().Email),
				FirstName:          convert.NullDotStringToPointerString(testutls.MockUser().FirstName),
				LastName:           convert.NullDotStringToPointerString(testutls.MockUser().LastName),
				Username:           convert.NullDotStringToPointerString(testutls.MockUser().Username),
				Mobile:             convert.NullDotStringToPointerString(testutls.MockUser().Mobile),
				Address:            convert.NullDotStringToPointerString(testutls.MockUser().Address),
				Active:             convert.NullDotBoolToPointerBool(testutls.MockUser().Active),
				LastLogin:          convert.NullDotTimeToPointerInt(testutls.MockUser().LastLogin),
				LastPasswordChange: convert.NullDotTimeToPointerInt(testutls.MockUser().LastPasswordChange),
				DeletedAt:          convert.NullDotTimeToPointerInt(testutls.MockUser().DeletedAt),
				UpdatedAt:          convert.NullDotTimeToPointerInt(testutls.MockUser().UpdatedAt),
			},
			wantErr: false,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mock, db, _ := testutls.SetupEnvAndDB(t, testutls.Parameters{
				EnvFileLocation: "../.env.local",
			})
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on Create User" {
				// insert new user
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// insert new user
			rows := sqlmock.NewRows([]string{
				"id",
				"mobile",
				"address",
				"active",
				"last_login",
				"last_password_change",
				"token",
				"deleted_at",
			}).AddRow(
				testutls.MockUser().ID,
				testutls.MockUser().Mobile,
				testutls.MockUser().Address,
				testutls.MockUser().Active,
				testutls.MockUser().LastLogin,
				testutls.MockUser().LastPasswordChange,
				testutls.MockUser().Token,
				testutls.MockUser().DeletedAt,
			)
			ApplyFunc(bcrypt.GenerateFromPassword, func([]uint8, int) ([]uint8, error) {
				var a []uint8
				return a, nil
			})
			mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
				WithArgs(
					testutls.MockUser().FirstName,
					testutls.MockUser().LastName,
					testutls.MockUser().Username,
					"",
					testutls.MockUser().Email,
					testutls.MockUser().RoleID,
					AnyTime{},
					AnyTime{},
				).WillReturnRows(rows)

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
		wantResp *fm.User
		wantErr  bool
	}{
		{
			name:    "Fail on finding User",
			req:     &fm.UserUpdateInput{},
			wantErr: true,
		},
		{
			name: "Success",
			req: &fm.UserUpdateInput{
				FirstName: &testutls.MockUser().FirstName.String,
				LastName:  &testutls.MockUser().LastName.String,
				Address:   &testutls.MockUser().Address.String,
			},
			wantResp: &fm.User{
				ID:        "0",
				FirstName: &testutls.MockUser().FirstName.String,
				LastName:  &testutls.MockUser().LastName.String,
				Address:   &testutls.MockUser().Address.String,
			},
			wantErr: false,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mock, db, _ := testutls.SetupEnvAndDB(t, testutls.Parameters{
				EnvFileLocation: "../.env.local",
			})
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on finding User" {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "users"`)).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}

			rows := sqlmock.NewRows([]string{
				"first_name",
			}).AddRow(
				testutls.MockUser().FirstName,
			)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users"`)).
				WithArgs(0).
				WillReturnRows(rows)

			// update users with new information
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
				WillReturnResult(result)

			c := context.Background()
			ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
			response, err := resolver1.Mutation().UpdateUser(ctx, tt.req)
			if tt.wantResp != nil && response != nil {
				assert.Equal(t, tt.wantResp, response)
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
			err := godotenv.Load("../.env.local")
			if err != nil {
				fmt.Print("error loading .env file")
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
			ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
			response, err := resolver1.Mutation().DeleteUser(ctx)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
