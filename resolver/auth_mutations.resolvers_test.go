package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	fm "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/internal/jwt"
	"go-template/internal/service"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/resultwrapper"
	"go-template/resolver"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	UserRoleName               = "UserRole"
	SuperAdminRoleName         = "SuperAdminRole"
	ErrorFromRedisCache        = "RedisCache Error"
	ErrorFromGetRole           = "RedisCache GetRole Error"
	ErrorUnauthorizedUser      = "Unauthorized User"
	ErrorFromCreateRole        = "CreateRole Error"
	ErrorPasswordValidation    = "Fail on PasswordValidation"
	ErrorActiveStatus          = "Fail on ActiveStatus"
	ErrorInsecurePassword      = "Insecure password"
	ErrorInvalidToken          = "Fail on FindByToken"
	ErrorUpdateUser            = "User Update Error"
	ErrorDeleteUser            = "User Delete Error"
	ErrorFromConfig            = "Config Error"
	ErrorFromBool              = "Boolean Error"
	ErrorMsgFromConfig         = "error in loading config"
	ErrorMsginvalidToken       = "error from FindByToken"
	ErrorMsgFindingUser        = "error in finding the user"
	ErrorMsgFromJwt            = "error in creating auth service "
	ErrorMsgfromUpdateUser     = "error while updating user"
	ErrorMsgPasswordValidation = "username or password does not exist "
	TestPasswordHash           = "$2a$10$dS5vK8hHmG5"
	OldPasswordHash            = "$2a$10$dS5vK8hHmG5gzwV8f7TK5.WHviMBqmYQLYp30a3XvqhCW9Wvl2tOS"
	SuccessCase                = "Success"
	ErrorFindingUser           = "Fail on finding user"
	ErrorFromCreateUser        = "Fail on Create User"
	ErrorFromThrottleCheck     = "Throttle error"
	ErrorFromJwt               = "Jwt Error"
	ErrorFromGenerateToken     = "Token Error"
	OldPassword                = "adminuser"
	NewPassword                = "adminuser!A9@"
	TestPassword               = "pass123"
	TestUsername               = "wednesday"
	TestToken                  = "refreshToken"
	ReqToken                   = "refresh_token"
)

func TestLogin(
	t *testing.T,
) {
	type args struct {
		UserName string
		Password string
	}
	cases := []struct {
		name     string
		req      args
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		{
			name: ErrorFindingUser,
			req: args{
				UserName: TestUsername,
				Password: TestPassword,
			},
			wantErr: true,
			err:     fmt.Errorf(ErrorMsgFindingUser),
		},
		{
			name: ErrorPasswordValidation,
			req: args{
				UserName: testutls.MockEmail,
				Password: TestPassword,
			},
			wantErr: true,
			err:     fmt.Errorf(ErrorMsgPasswordValidation),
		},
		{
			name: ErrorActiveStatus,
			req: args{
				UserName: testutls.MockEmail,
				Password: OldPassword,
			},
			wantErr: true,
			err:     resultwrapper.ErrUnauthorized,
		},
		{
			name: ErrorFromConfig,
			req: args{
				UserName: testutls.MockEmail,
				Password: OldPassword,
			},
			wantErr: true,
			err:     fmt.Errorf(ErrorMsgFromConfig),
		},
		{
			name: ErrorFromJwt,
			req: args{
				UserName: testutls.MockEmail,
				Password: OldPassword,
			},
			wantErr: true,
			err:     fmt.Errorf(ErrorMsgFromJwt),
		},
		{
			name: ErrorFromGenerateToken,
			req: args{
				UserName: testutls.MockEmail,
				Password: OldPassword,
			},
			wantErr: true,
			err:     resultwrapper.ErrUnauthorized,
		},
		{
			name: ErrorUpdateUser,
			req: args{
				UserName: testutls.MockEmail,
				Password: OldPassword,
			},
			wantErr: true,
			err:     fmt.Errorf(ErrorMsgfromUpdateUser),
		},
		{
			name: SuccessCase,
			req: args{
				UserName: testutls.MockEmail,
				Password: OldPassword,
			},
			wantResp: &fm.LoginResponse{
				Token:        "jwttokenstring",
				RefreshToken: TestToken,
			},
		},
	}

	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {

				// Apply mock functions using go-monkey for cases where certain errors are expected
				// and defer their resetting

				// Handle the case where there is an error loading the config
				if tt.name == ErrorFromConfig {
					patch := gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
						return nil, fmt.Errorf(ErrorMsgFromConfig)
					})
					defer patch.Reset()
				}

				// Initialize a new JWT service
				var tg jwt.Service

				// Handle the case where there is an error creating the JWT service
				if tt.name == ErrorFromJwt {
					patch := gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {

						return tg, fmt.Errorf(ErrorMsgFromJwt)

					})
					defer patch.Reset()
				}

				patch := gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
					if tt.name == ErrorFromGenerateToken {
						return "", resultwrapper.ErrUnauthorized
					} else {
						return "", nil
					}

				})
				defer patch.Reset()

				// Load the environment variables from a .env file
				err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../"))
				if err != nil {
					fmt.Print("error loading .env file")
				}
				// Create a mock SQL database connection
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

				// Handle the case where there is an error finding the user
				if tt.name == ErrorFindingUser {
					// get user by username
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (username=$1) LIMIT 1;`)).
						WithArgs().
						WillReturnError(fmt.Errorf(ErrorMsgFindingUser))
				}

				// Handle the case where there is an error validating the password
				if tt.name == ErrorPasswordValidation {
					// get user by username
					rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
						AddRow(testutls.MockID, TestPasswordHash, true, 1)
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
						WithArgs().
						WillReturnRows(rows)
				}

				// Handle the case where the user is not active
				if tt.name == ErrorActiveStatus {
					// get user by username
					rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
						AddRow(testutls.MockID, OldPasswordHash, false, 1)
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
						WithArgs().
						WillReturnRows(rows)
				}

				// get user by username
				rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
					AddRow(testutls.MockID, OldPasswordHash, true, 1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
					WithArgs().
					WillReturnRows(rows)

				// Apply mock behavior for a successful query result
				if tt.name == SuccessCase || tt.name == ErrorUpdateUser {
					rows := sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "ADMIN")
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "roles".* FROM "roles" WHERE ("id" = $1) LIMIT 1`)).
						WithArgs([]driver.Value{1}...).
						WillReturnRows(rows)
				}

				// Handle the case where there is an error while updating the user
				if tt.name == ErrorUpdateUser {
					fmt.Println(tt.name, " test name ")
					mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnError(fmt.Errorf(ErrorMsgfromUpdateUser))
				} else {
					fmt.Println(tt.name, " test name ")
					result := driver.Result(driver.RowsAffected(1))
					mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnResult(result)
				}

				c := context.Background()

				// Call the login mutation with the given arguments and check the response and error against the expected values
				response, err := resolver1.Mutation().Login(c, tt.req.UserName, tt.req.Password)
				if tt.wantResp != nil &&
					response != nil {
					tt.wantResp.RefreshToken = response.RefreshToken
					tt.wantResp.Token = response.Token

					// Assert that the expected response matches the actual response
					assert.Equal(t, tt.wantResp, response)
				} else {

					// Assert that the expected error value matches the actual error value
					assert.Equal(t, true, strings.Contains(err.Error(), tt.err.Error()))
					assert.Equal(t, tt.wantErr, err != nil)
				}

			},
		)
	}
}

func TestChangePassword(
	t *testing.T,
) {

	// Define a struct to represent the change password request
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
			name: ErrorFindingUser,
			req: changeReq{
				OldPassword: TestPassword,
				NewPassword: NewPassword,
			},
			wantErr: true,
		},
		{
			name: ErrorPasswordValidation,
			req: changeReq{
				OldPassword: TestPassword,
				NewPassword: NewPassword,
			},
			wantErr: true,
		},
		{
			name: ErrorInsecurePassword,
			req: changeReq{
				OldPassword: OldPassword,
				NewPassword: testutls.MockEmail,
			},
			wantErr: true,
		},
		{
			name: ErrorUpdateUser,
			req: changeReq{
				OldPassword: OldPassword,
				NewPassword: NewPassword,
			},
			wantErr: true,
		},
		{
			name: ErrorFromConfig,
			req: changeReq{
				OldPassword: OldPassword,
				NewPassword: testutls.MockEmail,
			},
			wantErr: true,
		},

		{
			name: SuccessCase,
			req: changeReq{
				OldPassword: OldPassword,
				NewPassword: NewPassword,
			},
			wantResp: &fm.ChangePasswordResponse{
				Ok: true,
			},
			wantErr: false,
		},
	}

	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {

				// Handle the case where there is an error while loading the configuration
				if tt.name == ErrorFromConfig {
					patch := gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
						return nil, fmt.Errorf("error in loading config")
					})
					defer patch.Reset()
				}

				// Load environment variables
				err := config.LoadEnv()
				if err != nil {
					fmt.Print("error loading .env file")
				}

				// Create a mock SQL database connection
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

				// Handle the case where there is an error while finding the user
				if tt.name == ErrorFindingUser {
					mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
						WithArgs().
						WillReturnError(fmt.Errorf(""))
				}
				// Expect a query to get the user by ID to return a row with mock data entered
				rows := sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(testutls.MockID, testutls.MockEmail, OldPasswordHash)
				mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
					WithArgs().
					WillReturnRows(rows)

					// Handle the case where the password update is successful
				if tt.name == SuccessCase {
					result := driver.Result(driver.RowsAffected(1))
					mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnResult(result)
				}

				// Handle the case where there is an error while updating the user's password
				if tt.name == ErrorUpdateUser {

					mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnError(fmt.Errorf("errrorr"))
				}

				// Set up the context with the mock user
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())

				// Call the ChangePassword mutation and check the response and error against the expected values
				response, err := resolver1.Mutation().ChangePassword(ctx, tt.req.OldPassword, tt.req.NewPassword)
				if tt.wantResp != nil {

					// Assert that the expected response matches the actual response
					assert.Equal(t, tt.wantResp, response)
				}
				// Assert that the expected error value matches the actual error value
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}

func TestRefreshToken(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantResp *fm.RefreshTokenResponse
		wantErr  bool
		err      error
	}{
		{
			name:    ErrorInvalidToken,
			req:     TestToken,
			wantErr: true,
			err:     fmt.Errorf(ErrorMsginvalidToken),
		},
		{
			name:    ErrorFromConfig,
			req:     ReqToken,
			wantErr: true,
			err:     fmt.Errorf(ErrorMsgFromConfig),
		},
		{
			name:    ErrorFromJwt,
			req:     ReqToken,
			wantErr: true,
			err:     fmt.Errorf(ErrorMsgFromJwt),
		},
		{
			name:    ErrorFromGenerateToken,
			req:     ReqToken,
			wantErr: true,
			err:     resultwrapper.ErrUnauthorized},
		{
			name: SuccessCase,
			req:  ReqToken,
			wantResp: &fm.RefreshTokenResponse{
				Token: testutls.MockToken,
			},
			wantErr: false,
		},
	}

	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {

		t.Run(
			tt.name,
			func(t *testing.T) {

				// Create a mock SQL database connection
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

				// Handle the case where authentication token is invalid
				if tt.name == ErrorInvalidToken {
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
						WithArgs().
						WillReturnError(fmt.Errorf(ErrorMsginvalidToken))
				}

				// Handle the case where there is an error loading the config

				patch := gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
					if tt.name == ErrorFromConfig {
						return nil, fmt.Errorf("error in loading config")
					} else {
						return &config.Configuration{}, nil
					}
				})
				defer patch.Reset()

				//initialize a jwt service
				tg := jwt.Service{}

				// Handle the case where there is an error creating the JWT service
				patchJWT := gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					if tt.name == ErrorFromJwt {
						return tg, fmt.Errorf(ErrorMsgFromJwt)
					} else {
						return tg, nil
					}
				})
				defer patchJWT.Reset()

				// Handle the case where there is an error form token generation service
				patchGenerateToken := gomonkey.ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						if tt.name == ErrorFromGenerateToken {
							return "", resultwrapper.ErrUnauthorized
						} else {
							return "token", nil
						}
					})
				defer patchGenerateToken.Reset()

				// Expect a query to get the user by ID to return a row with mock data entered
				rows := sqlmock.NewRows([]string{"id", "email", "token", "role_id"}).
					AddRow(1, testutls.MockEmail, testutls.MockToken, 1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
					WithArgs().
					WillReturnRows(rows)

				if tt.name == SuccessCase {
					rows := sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "ADMIN")
					mock.ExpectQuery(regexp.QuoteMeta(`SELECT "roles".* FROM "roles" WHERE ("id" = $1) LIMIT 1`)).
						WithArgs([]driver.Value{1}...).
						WillReturnRows(rows)
				}

				// Set up the context with the mock user
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())

				// Call the refresh token mutation with the given arguments and check the response and error against the expected values
				response, err := resolver1.Mutation().
					RefreshToken(ctx, tt.req)
				if tt.wantResp != nil &&
					response != nil {
					tt.wantResp.Token = response.Token

					// Assert that the expected response matches the actual response
					assert.Equal(t, tt.wantResp, response)
				} else {

					// Assert that the expected error value matches the actual error value
					assert.Equal(t, true, strings.Contains(err.Error(), tt.err.Error()))
				}

			},
		)
	}
}
