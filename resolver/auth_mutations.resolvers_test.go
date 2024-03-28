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
	ErrorMessage               = "an error '%s' was not expected when opening a stub database connection"
)

func TestLogin(t *testing.T) {
	cases := prepareTestCases()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			handleTestCases(t, tt)
		})
	}
}

type LoginArgs struct {
	UserName string
	Password string
}

var req = LoginArgs{
	UserName: testutls.MockEmail,
	Password: OldPassword,
}

func prepareTestCases() []struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	cases := []struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		prepareErrorFindingUser(),
		prepareErrorPasswordValidation(),
		prepareErrorActiveStatus(),
		prepareErrorFromConfig(),
		prepareErrorFromJwt(),
		prepareErrorFromGenerateToken(),
		prepareErrorUpdateUser(),
		prepareSuccessCase(),
	}
	return cases
}

func prepareErrorFindingUser() struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	return struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		name: ErrorFindingUser,
		req: LoginArgs{
			UserName: TestUsername,
			Password: TestPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFindingUser),
	}
}

func prepareErrorPasswordValidation() struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	return struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		name: ErrorPasswordValidation,
		req: LoginArgs{
			UserName: testutls.MockEmail,
			Password: TestPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgPasswordValidation),
	}
}
func prepareErrorActiveStatus() struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	return struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		name:    ErrorActiveStatus,
		req:     req, // Assuming req is defined somewhere in your code
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
	}
}

func prepareErrorFromConfig() struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	return struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		name: ErrorFromConfig,
		req: LoginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromConfig),
	}
}

func prepareErrorFromJwt() struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	return struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		name:    ErrorFromJwt,
		req:     req,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromJwt),
	}
}

func prepareErrorFromGenerateToken() struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	return struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		name:    ErrorFromGenerateToken,
		req:     req,
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
	}
}

func prepareErrorUpdateUser() struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	return struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		name:    ErrorUpdateUser,
		req:     req,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgfromUpdateUser),
	}
}

func prepareSuccessCase() struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
} {
	return struct {
		name     string
		req      LoginArgs
		wantResp *fm.LoginResponse
		wantErr  bool
		err      error
	}{
		name: SuccessCase,
		req:  req,
		wantResp: &fm.LoginResponse{
			Token:        "jwttokenstring",
			RefreshToken: TestToken,
		},
	}
}

func handleTestCases(t *testing.T, tt struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
}) {
	// Prepare necessary mocks and patches
	prepareMocksAndPatches(tt)

	// Load environment variables
	loadEnvironmentVariables()
	// Handle specific test cases
	handleSpecificTestCase(t, tt)
}

func prepareMocksAndPatches(tt struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
}) {
	if tt.name == ErrorFromConfig {
		patch := gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
			return nil, fmt.Errorf(ErrorMsgFromConfig)
		})
		defer patch.Reset()
	}

	var tg jwt.Service
	if tt.name == ErrorFromJwt {
		patch := gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
			return tg, fmt.Errorf(ErrorMsgFromJwt)
		})
		defer patch.Reset()
	}
	patch := gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
		if tt.name == ErrorFromGenerateToken {
			return "", resultwrapper.ErrUnauthorized
		}
		return "", nil
	})
	defer patch.Reset()
}

func loadEnvironmentVariables() {
	err := config.LoadEnv()
	if err != nil {
		fmt.Print("error loading .env file")
	}
}

func handleSpecificTestCase(t *testing.T, tt struct {
	name     string
	req      LoginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
}) {
	resolver1 := resolver.Resolver{}
	mock, cleanup, _ := testutls.SetupMockDB(t)
	// Mock database queries based on test case
	switch tt.name {
	case ErrorFindingUser:
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (username=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnError(fmt.Errorf(ErrorMsgFindingUser))
	case ErrorPasswordValidation:
		rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
			AddRow(testutls.MockID, TestPasswordHash, true, 1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rows)
	case ErrorActiveStatus:
		rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
			AddRow(testutls.MockID, OldPasswordHash, false, 1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rows)
	default:
		rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
			AddRow(testutls.MockID, OldPasswordHash, true, 1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rows)
	}

	if tt.name == SuccessCase || tt.name == ErrorUpdateUser {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "ADMIN")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "roles".* FROM "roles" WHERE ("id" = $1) LIMIT 1`)).
			WithArgs([]driver.Value{1}...).
			WillReturnRows(rows)
	}

	if tt.name == ErrorUpdateUser {
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnError(fmt.Errorf(ErrorMsgfromUpdateUser))
	} else {
		result := driver.Result(driver.RowsAffected(1))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnResult(result)
	}

	// Execute resolver function
	c := context.Background()
	response, err := resolver1.Mutation().Login(c, tt.req.UserName, tt.req.Password)

	// Assert results
	if tt.wantResp != nil && response != nil {
		tt.wantResp.RefreshToken = response.RefreshToken
		tt.wantResp.Token = response.Token
		assert.Equal(t, tt.wantResp, response)
	} else {
		assert.Equal(t, true, strings.Contains(err.Error(), tt.err.Error()))
		assert.Equal(t, tt.wantErr, err != nil)
	}
	cleanup()
}

type changeReq struct {
	OldPassword string
	NewPassword string
}

func GetChangePasswordTestCases() []struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
} {
	cases := []struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		ErrorFindingUserCase(),
		ErrorPasswordValidationCase(),
		ErrorInsecurePasswordCase(),
		ErrorUpdateUserCase(),
		ErrorFromConfigCase(),
		GetSuccessCase(),
	}
	return cases
}
func ErrorFindingUserCase() struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
} {
	return struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		name: ErrorFindingUser,
		req: changeReq{
			OldPassword: TestPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
	}
}

func ErrorPasswordValidationCase() struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
} {
	return struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		name: ErrorPasswordValidation,
		req: changeReq{
			OldPassword: TestPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
	}
}

func ErrorInsecurePasswordCase() struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
} {
	return struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		name: ErrorInsecurePassword,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: testutls.MockEmail,
		},
		wantErr: true,
	}
}
func ErrorUpdateUserCase() struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
} {
	return struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		name: ErrorUpdateUser,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
	}
}

func ErrorFromConfigCase() struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
} {
	return struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		name: ErrorFromConfig,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: testutls.MockEmail,
		},
		wantErr: true,
	}
}

func GetSuccessCase() struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
} {
	return struct {
		name     string
		req      changeReq
		wantResp *fm.ChangePasswordResponse
		wantErr  bool
	}{
		name: SuccessCase,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: NewPassword,
		},
		wantResp: &fm.ChangePasswordResponse{
			Ok: true,
		},
		wantErr: false,
	}
}
func TestChangePassword(
	t *testing.T,
) {
	// Define a struct to represent the change password request
	cases := GetChangePasswordTestCases()
	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				// Handle the case where there is an error while loading the configuration
				if tt.name == ErrorFromConfig {
					patch := gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
						return nil, fmt.Errorf(ErrorMsgFromConfig)
					})
					defer patch.Reset()
				}
				// Load environment variables
				err := config.LoadEnv()
				if err != nil {
					fmt.Print("error loading .env file")
				}
				// Create a mock SQL database connection
				mock, cleanup, _ := testutls.SetupMockDB(t)
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
				cleanup()
			},
		)
	}
}

func TestRefreshToken(t *testing.T) {
	cases := prepareRefreshTokenCases()

	resolver := resolver.Resolver{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mock, cleanup, _ := testutls.SetupMockDB(t)

			setupMockDBExpectations(tt.name, mock)

			mockConfigLoad := prepareMockConfigLoad(tt.name)
			defer mockConfigLoad.Reset()

			mockJWTService := prepareMockJWTService(tt.name)
			defer mockJWTService.Reset()

			mockTokenGeneration := prepareMockTokenGeneration(tt.name)
			defer mockTokenGeneration.Reset()

			ctx := prepareContextRefereshToken()

			response, err := resolver.Mutation().RefreshToken(ctx, tt.req)

			if tt.wantResp != nil && response != nil {
				tt.wantResp.Token = response.Token
				assert.Equal(t, tt.wantResp, response)
			} else {
				assert.Equal(t, true, strings.Contains(err.Error(), tt.err.Error()))
			}
			cleanup()
		})
	}
}

func prepareRefreshTokenCases() []struct {
	name     string
	req      string
	wantResp *fm.RefreshTokenResponse
	wantErr  bool
	err      error
} {
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
	return cases
}

func setupMockDBExpectations(name string, mock sqlmock.Sqlmock) {
	switch name {
	case ErrorInvalidToken:
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnError(fmt.Errorf(ErrorMsginvalidToken))
	case SuccessCase:
		rows := sqlmock.NewRows([]string{"id", "email", "token", "role_id"}).
			AddRow(1, testutls.MockEmail, testutls.MockToken, 1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(rows)
		rows = sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "ADMIN")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "roles".* FROM "roles" WHERE ("id" = $1) LIMIT 1`)).
			WithArgs([]driver.Value{1}...).
			WillReturnRows(rows)
	case ErrorFromConfig:
		// Expectation for error loading the config
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnError(fmt.Errorf(ErrorMsgFromConfig))
	case ErrorFromJwt:
		// Expectation for error in creating the JWT service
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "token", "role_id"}))
	}
}

func prepareMockConfigLoad(name string) *gomonkey.Patches {
	return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
		if name == ErrorFromConfig {
			return nil, fmt.Errorf(ErrorMsgFromConfig)
		}
		return &config.Configuration{}, nil
	})
}

func prepareMockJWTService(name string) *gomonkey.Patches {
	return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
		if name == ErrorFromJwt {
			return jwt.Service{}, fmt.Errorf(ErrorMsgFromJwt)
		}
		return jwt.Service{}, nil
	})
}

func prepareMockTokenGeneration(name string) *gomonkey.Patches {
	return gomonkey.ApplyMethod(reflect.TypeOf(jwt.Service{}), "GenerateToken",
		func(jwt.Service, *models.User) (string, error) {
			if name == ErrorFromGenerateToken {
				return "", resultwrapper.ErrUnauthorized
			}
			return "token", nil
		})
}

func prepareContextRefereshToken() context.Context {
	c := context.Background()
	return context.WithValue(c, testutls.UserKey, testutls.MockUser())
}
