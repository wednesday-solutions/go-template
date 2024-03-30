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

type loginArgs struct {
	UserName string
	Password string
}

type LoginType struct {
	name     string
	req      loginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
}

func errorFindingUserCase() LoginType {
	return LoginType{
		name: ErrorFindingUser,
		req: loginArgs{
			UserName: TestUsername,
			Password: TestPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFindingUser),
	}
}
func errorPasswordValidationCase() LoginType {
	return LoginType{
		name: ErrorPasswordValidation,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: TestPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgPasswordValidation),
	}
}

func errorActiveStatusCase() LoginType {
	return LoginType{
		name: ErrorActiveStatus,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
	}
}
func errorFromConfigCase() LoginType {
	return LoginType{
		name: ErrorFromConfig,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromConfig),
	}
}

func errorFromJwtCase() LoginType {
	return LoginType{
		name: ErrorFromJwt,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromJwt),
	}
}
func errorFromGenerateTokenCase() LoginType {
	return LoginType{
		name: ErrorFromGenerateToken,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
	}
}
func errorUpdateUserCase() LoginType {
	return LoginType{
		name: ErrorUpdateUser,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgfromUpdateUser),
	}
}
func loginSuccessCase() LoginType {
	return LoginType{
		name: SuccessCase,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantResp: &fm.LoginResponse{
			Token:        "jwttokenstring",
			RefreshToken: TestToken,
		},
	}
}
func loadLoginTestCases() []LoginType {
	return []LoginType{
		errorFindingUserCase(),
		errorPasswordValidationCase(),
		errorActiveStatusCase(),
		errorFromConfigCase(),
		errorFromJwtCase(),
		errorFromGenerateTokenCase(),
		errorUpdateUserCase(),
		loginSuccessCase(),
	}
}

func applyGenerateTokenPatchLogin(name string) *gomonkey.Patches {
	tg := jwt.Service{}
	return gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
		if name == ErrorFromGenerateToken {
			return "", resultwrapper.ErrUnauthorized
		} else {
			return "", nil
		}
	})
}

func setupLoginSQLMocks(name string, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	if name == ErrorFindingUser {
		// get user by username
		return mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (username=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnError(fmt.Errorf(ErrorMsgFindingUser))
	}
	// Handle the case where there is an error validating the password
	if name == ErrorPasswordValidation {
		// get user by username
		rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
			AddRow(testutls.MockID, TestPasswordHash, true, 1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rows)
	}
	// Handle the case where the user is not active
	if name == ErrorActiveStatus {
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
	if name == SuccessCase || name == ErrorUpdateUser {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "ADMIN")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "roles".* FROM "roles" WHERE ("id" = $1) LIMIT 1`)).
			WithArgs([]driver.Value{1}...).
			WillReturnRows(rows)
	}
	return nil
}

func setupLoginExpectedExec(name string, mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
	if name == ErrorUpdateUser {
		fmt.Println(name, " test name ")
		return mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnError(fmt.Errorf(ErrorMsgfromUpdateUser))
	} else {
		fmt.Println(name, " test name ")
		result := driver.Result(driver.RowsAffected(1))
		return mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnResult(result)
	}
}
func TestLogin(
	t *testing.T,
) {
	cases := loadLoginTestCases()
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
					patch := applyPatch(tt.name)
					defer patch.Reset()
				}
				// Handle the case where there is an error creating the JWT service
				if tt.name == ErrorFromJwt {
					patch := applyJWTPatch(tt.name)
					defer patch.Reset()
				}
				patch := applyGenerateTokenPatchLogin(tt.name)
				defer patch.Reset()
				// Load the environment variables from a .env file
				err := config.LoadEnv()
				if err != nil {
					fmt.Print("error loading .env file")
				}
				// Create a mock SQL database connection
				db, mock, _ := sqlmock.New()
				// Inject mock instance into boil.
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)
				// Handle the case where there is an error finding the user
				setupLoginSQLMocks(tt.name, mock)
				// Handle the case where there is an error while updating the user
				setupLoginExpectedExec(tt.name, mock)
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

type changeReq struct {
	OldPassword string
	NewPassword string
}

type changePasswordType struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
}

func loadChangePasswordTestCases() []changePasswordType {
	return []changePasswordType{
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
}
func setupSQLExpectWxecution(name string, mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
	if name == SuccessCase {
		result := driver.Result(driver.RowsAffected(1))
		return mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnResult(result)
	}
	// Handle the case where there is an error while updating the user's password
	if name == ErrorUpdateUser {
		return mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnError(fmt.Errorf("errrorr"))
	}
	return nil
}

func setupSQLMocks(name string, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	if name == ErrorFindingUser {
		return mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
			WithArgs().
			WillReturnError(fmt.Errorf(""))
	}
	// Expect a query to get the user by ID to return a row with mock data entered
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(testutls.MockID, testutls.MockEmail, OldPasswordHash)
	return mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
		WithArgs().
		WillReturnRows(rows)
}

func TestChangePassword(
	t *testing.T,
) {
	// Define a struct to represent the change password request
	cases := loadChangePasswordTestCases()
	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				// Handle the case where there is an error while loading the configuration
				if tt.name == ErrorFromConfig {
					patch := applyPatch(tt.name)
					defer patch.Reset()
				}
				// Load environment variables
				err := config.LoadEnv()
				if err != nil {
					fmt.Print("error loading .env file")
				}
				// Create a mock SQL database connection
				db, mock, _ := sqlmock.New()
				// Inject mock instance into boil.
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)
				// Handle the case where there is an error while finding the user
				setupSQLMocks(tt.name, mock)

				// Handle the case where the password update is successful
				if tt.name == SuccessCase || tt.name == ErrorUpdateUser {
					setupSQLExpectWxecution(tt.name, mock)
				}

				// Handle the case where there is an error while updating the user's passwor

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

type refereshTokenType struct {
	name     string
	req      string
	wantResp *fm.RefreshTokenResponse
	wantErr  bool
	err      error
}

func loadRefereshTokenCases() []refereshTokenType {
	return []refereshTokenType{
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
}

func setupSQLExpectation(name string, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	if name == ErrorInvalidToken {
		return mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnError(fmt.Errorf(ErrorMsginvalidToken))
	}
	rows := sqlmock.NewRows([]string{"id", "email", "token", "role_id"}).
		AddRow(1, testutls.MockEmail, testutls.MockToken, 1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
		WithArgs().
		WillReturnRows(rows)

	if name == SuccessCase {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "ADMIN")
		return mock.ExpectQuery(regexp.QuoteMeta(`SELECT "roles".* FROM "roles" WHERE ("id" = $1) LIMIT 1`)).
			WithArgs([]driver.Value{1}...).
			WillReturnRows(rows)
	}
	return nil
}
func applyPatch(name string) *gomonkey.Patches {
	return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
		if name == ErrorFromConfig {
			return nil, fmt.Errorf("error in loading config")
		} else {
			return &config.Configuration{}, nil
		}
	})
}

func applyJWTPatch(name string) *gomonkey.Patches {
	tg := jwt.Service{}
	return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
		if name == ErrorFromJwt {
			return tg, fmt.Errorf(ErrorMsgFromJwt)
		} else {
			return tg, nil
		}
	})
}

func applyGenerateTokenPatch(name string) *gomonkey.Patches {
	tg := jwt.Service{}
	return gomonkey.ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
		func(jwt.Service, *models.User) (string, error) {
			if name == ErrorFromGenerateToken {
				return "", resultwrapper.ErrUnauthorized
			} else {
				return "token", nil
			}
		})
}
func TestRefreshToken(t *testing.T) {
	cases := loadRefereshTokenCases()
	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				// Create a mock SQL database connection
				db, mock, _ := sqlmock.New()
				// Inject mock instance into boil.
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)
				// Handle the case where authentication token is invalid
				setupSQLExpectation(tt.name, mock)
				// Handle the case where there is an error loading the config
				patch := applyPatch(tt.name)
				defer patch.Reset()
				//initialize a jwt service
				// Handle the case where there is an error creating the JWT service
				patchJWT := applyJWTPatch(tt.name)
				defer patchJWT.Reset()
				// Handle the case where there is an error form token generation service
				patchGenerateToken := applyGenerateTokenPatch(tt.name)
				defer patchGenerateToken.Reset()
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
