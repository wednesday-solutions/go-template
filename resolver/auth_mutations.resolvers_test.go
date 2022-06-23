package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"

	fm "go-template/gqlmodels"
	"go-template/resolver"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestLogin(t *testing.T) {
	type args struct {
		UserName string
		Password string
	}
	cases := []struct {
		name     string
		req      args
		wantResp *fm.LoginResponse
		wantErr  bool
	}{
		{
			name:    "Fail on FindByUser",
			req:     args{UserName: "wednesday", Password: "pass123"},
			wantErr: true,
		},
		{
			name:     "Success",
			req:      args{UserName: testutls.MockEmail, Password: "adminuser"},
			wantResp: &fm.LoginResponse{Token: "jwttokenstring", RefreshToken: "refreshtoken"},
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
			if tt.name == "Fail on FindByUser" {
				// get user by username
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (username=$1) LIMIT 1;")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// get user by username
			rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
				AddRow(testutls.MockID, "$2a$10$dS5vK8hHmG5gzwV8f7TK5.WHviMBqmYQLYp30a3XvqhCW9Wvl2tOS", true, 1)
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (username=$1) LIMIT 1;")).
				WithArgs().
				WillReturnRows(rows)

			if tt.name == "Success" {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "ADMIN")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"roles\" WHERE (\"id\" = $1) LIMIT 1")).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(rows)
			}

			// update users with token
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" ")).
				WillReturnResult(result)

			c := context.Background()
			response, err := resolver1.Mutation().Login(c, tt.req.UserName, tt.req.Password)
			if tt.wantResp != nil && response != nil {
				tt.wantResp.RefreshToken = response.RefreshToken
				tt.wantResp.Token = response.Token
				assert.Equal(t, tt.wantResp, response)
			}

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestChangePassword(t *testing.T) {

	type changeReq struct {
		OldPassword string
		NewPassword string
	}
	cases := []struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		{
			name:    "Fail on FindByUser",
			req:     changeReq{OldPassword: "adminuser!A9@@@@", NewPassword: "adminuser!A9@"},
			wantErr: true,
		},
		{
			name:    "Incorrect Old Password",
			req:     changeReq{OldPassword: "admin", NewPassword: "adminuser!A9@"},
			wantErr: true,
		},
		{
			name:     "Success",
			req:      changeReq{OldPassword: "adminuser", NewPassword: "adminuser!A9@"},
			wantResp: &fm.ChangePasswordResponse{Ok: true},
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
			// Inject mock instance into boil.
			oldDB := boil.GetDB()
			defer func() {
				db.Close()
				boil.SetDB(oldDB)
			}()
			boil.SetDB(db)

			if tt.name == "Fail on FindByUser" {
				// get user by id
				mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// get user by id
			rows := sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow(testutls.MockID, testutls.MockEmail, "$2a$10$dS5vK8hHmG5gzwV8f7TK5.WHviMBqmYQLYp30a3XvqhCW9Wvl2tOS")
			mock.ExpectQuery(regexp.QuoteMeta("select * from \"users\" where \"id\"=$1")).
				WithArgs().
				WillReturnRows(rows)
			if tt.name == "Success" {
				// update password
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" ")).
					WillReturnResult(result)
			}

			c := context.Background()
			ctx := context.
				WithValue(c, testutls.UserKey, testutls.MockUser())
			response, err := resolver1.Mutation().ChangePassword(ctx, tt.req.OldPassword, tt.req.NewPassword)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantResp *fm.RefreshTokenResponse
		wantErr  bool
	}{
		{
			name:    "Fail on FindByToken",
			req:     "refreshToken",
			wantErr: true,
		},
		{
			name:     "Success",
			req:      "refresh_token",
			wantResp: &fm.RefreshTokenResponse{Token: testutls.MockToken},
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

			if tt.name == "Fail on FindByToken" {
				// get user by token
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (token=$1) LIMIT 1;")).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			}
			// get user by token
			rows := sqlmock.NewRows([]string{"id", "email", "token", "role_id"}).
				AddRow(1, testutls.MockEmail, testutls.MockToken, 1)
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"users\" WHERE (token=$1) LIMIT 1;")).
				WithArgs().
				WillReturnRows(rows)

			if tt.name == "Success" {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "ADMIN")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"roles\" WHERE (\"id\" = $1) LIMIT 1")).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(rows)
			}

			c := context.Background()
			ctx := context.
				WithValue(c,
					testutls.UserKey,
					testutls.MockUser(),
				)
			response, err := resolver1.Mutation().RefreshToken(ctx, tt.req)
			if tt.wantResp != nil && response != nil {
				tt.wantResp.Token = response.Token
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
