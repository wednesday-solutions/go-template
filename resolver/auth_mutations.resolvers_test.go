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

type loginType struct {
	name     string
	req      loginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
	init     func(mock sqlmock.Sqlmock) *gomonkey.Patches
}

func errorFindingUserCase() loginType {
	return loginType{
		name: ErrorFindingUser,
		req: loginArgs{
			UserName: TestUsername,
			Password: TestPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFindingUser),
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (username=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnError(fmt.Errorf(ErrorMsgFindingUser))
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
				return "", nil
			})
		},
	}
}
func errorPasswordValidationCase() loginType {
	return loginType{
		name: ErrorPasswordValidation,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: TestPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgPasswordValidation),
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
				AddRow(testutls.MockID, TestPasswordHash, true, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
				return "", nil
			})
		},
	}
}

func errorActiveStatusCase() loginType {
	return loginType{
		name: ErrorActiveStatus,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
				AddRow(testutls.MockID, OldPasswordHash, true, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
				return "", nil
			})
		},
	}
}
func errorFromConfigCase() loginType {
	return loginType{
		name: ErrorFromConfig,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromConfig),
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
				AddRow(testutls.MockID, OldPasswordHash, true, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return nil, fmt.Errorf("error in loading config")
			})
		},
	}
}

func errorFromJwtCase() loginType {
	return loginType{
		name: ErrorFromJwt,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromJwt),
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
				AddRow(testutls.MockID, OldPasswordHash, true, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, fmt.Errorf(ErrorMsgFromJwt)
			})
		},
	}
}
func errorFromGenerateTokenCase() loginType {
	return loginType{
		name: ErrorFromGenerateToken,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "password", "active", "role_id"}).
				AddRow(testutls.MockID, OldPasswordHash, true, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users"  WHERE (username=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
				return "", resultwrapper.ErrUnauthorized
			})
		},
	}
}
func errorUpdateUserCase() loginType {
	return loginType{
		name: ErrorUpdateUser,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgfromUpdateUser),
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "ADMIN")
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "roles".* FROM "roles" WHERE ("id" = $1) LIMIT 1`)).
				WithArgs([]driver.Value{1}...).
				WillReturnRows(rows)
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnError(fmt.Errorf(ErrorMsgfromUpdateUser))
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
				return "", nil
			})
		},
	}
}
func loginSuccessCase() loginType {
	return loginType{
		name: SuccessCase,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantResp: &fm.LoginResponse{
			Token:        "jwttokenstring",
			RefreshToken: TestToken,
		},
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "ADMIN")
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "roles".* FROM "roles" WHERE ("id" = $1) LIMIT 1`)).
				WithArgs([]driver.Value{1}...).
				WillReturnRows(rows)
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnResult(result)
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
				return "", nil
			})
		},
	}
}
func loadLoginTestCases() []loginType {
	return []loginType{
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
				db, mock, _ := sqlmock.New()
				// Inject mock instance into boil.
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)
				patch := tt.init(mock)
				// Load the environment variables from a .env file
				err := config.LoadEnv()
				if err != nil {
					fmt.Print("error loading .env file")
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
				patch.Reset()
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
	init     func(mock sqlmock.Sqlmock) *gomonkey.Patches
}

func changePasswordErrorFindingUserCase() changePasswordType {
	return changePasswordType{
		name: ErrorFindingUser,
		req: changeReq{
			OldPassword: TestPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
			return nil
		},
	}
}
func changePasswordErrorPasswordValidationcase() changePasswordType {
	return changePasswordType{
		name: ErrorPasswordValidation,
		req: changeReq{
			OldPassword: TestPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow(testutls.MockID, testutls.MockEmail, OldPasswordHash)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnRows(rows)
			return nil
		},
	}
}

func changePasswordErrorInsecurePasswordCase() changePasswordType {
	return changePasswordType{
		name: ErrorInsecurePassword,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: testutls.MockEmail,
		},
		wantErr: true,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow(testutls.MockID, testutls.MockEmail, OldPasswordHash)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnRows(rows)
			return nil
		},
	}
}

func changePasswordErrorUpdateUserCase() changePasswordType {
	return changePasswordType{
		name: ErrorUpdateUser,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow(testutls.MockID, testutls.MockEmail, OldPasswordHash)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnRows(rows)
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnError(fmt.Errorf("errrorr"))
			return nil
		},
	}
}

func changePasswordErrorFromConfigCase() changePasswordType {
	return changePasswordType{
		name: ErrorFromConfig,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: testutls.MockEmail,
		},
		wantErr: true,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow(testutls.MockID, testutls.MockEmail, OldPasswordHash)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnRows(rows)
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, fmt.Errorf(ErrorMsgFromJwt)
			})
		},
	}
}

func changePasswordSuccessCase() changePasswordType {
	return changePasswordType{
		name: SuccessCase,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: NewPassword,
		},
		wantResp: &fm.ChangePasswordResponse{
			Ok: true,
		},
		wantErr: false,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			rows := sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow(testutls.MockID, testutls.MockEmail, OldPasswordHash)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnRows(rows)
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" `)).WillReturnResult(result)
			return nil
		},
	}
}
func loadChangePasswordTestCases() []changePasswordType {
	return []changePasswordType{
		changePasswordErrorFindingUserCase(),
		changePasswordErrorPasswordValidationcase(),
		changePasswordErrorInsecurePasswordCase(),
		changePasswordErrorUpdateUserCase(),
		changePasswordErrorFromConfigCase(),
		changePasswordSuccessCase(),
	}
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
				db, mock, _ := sqlmock.New()
				// Inject mock instance into boil.
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)
				tt.init(mock)
				// Load environment variables
				err := config.LoadEnv()
				if err != nil {
					fmt.Print("error loading .env file")
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

type refereshTokenType struct {
	name      string
	req       string
	wantResp  *fm.RefreshTokenResponse
	wantErr   bool
	err       error
	init      refereshTokenPatches
	initMocks func(mock sqlmock.Sqlmock)
}
type refereshTokenPatches struct {
	configPatch func() *gomonkey.Patches
	jwtPatch    func() *gomonkey.Patches
	tokenPatch  func() *gomonkey.Patches
}

func refreshTokenInvalidCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorInvalidToken,
		req:     TestToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsginvalidToken),
		init: refereshTokenPatches{
			configPatch: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
					return &config.Configuration{}, nil
				})
			},
			jwtPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, nil
				})
			},
			tokenPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return "token", nil
					})
			},
		},
		initMocks: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnError(fmt.Errorf(ErrorMsginvalidToken))
		},
	}
}
func refreshTokenErrorFromConfigCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorFromConfig,
		req:     ReqToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromConfig),
		init: refereshTokenPatches{
			configPatch: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
					return nil, fmt.Errorf("error in loading config")
				})
			},
			jwtPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, nil
				})
			},
			tokenPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return "token", nil
					})
			},
		},
		initMocks: func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id", "email", "token", "role_id"}).
				AddRow(1, testutls.MockEmail, testutls.MockToken, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
		},
	}
}

func refereshTokenErrorFromJWTCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorFromJwt,
		req:     ReqToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromJwt),
		init: refereshTokenPatches{
			configPatch: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
					return &config.Configuration{}, nil
				})
			},
			jwtPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, fmt.Errorf(ErrorMsgFromJwt)
				})
			},
			tokenPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return "token", nil
					})
			},
		},
		initMocks: func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id", "email", "token", "role_id"}).
				AddRow(1, testutls.MockEmail, testutls.MockToken, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
		},
	}
}

func refereshTokenErrorFromGenerateTokenCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorFromGenerateToken,
		req:     ReqToken,
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
		init: refereshTokenPatches{
			configPatch: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
					return &config.Configuration{}, nil
				})
			},
			jwtPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, nil
				})
			},
			tokenPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return "", resultwrapper.ErrUnauthorized
					})
			},
		},
		initMocks: func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id", "email", "token", "role_id"}).
				AddRow(1, testutls.MockEmail, testutls.MockToken, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
		},
	}
}

func refreshTokenSuccessCase() refereshTokenType {
	return refereshTokenType{
		name: SuccessCase,
		req:  ReqToken,
		wantResp: &fm.RefreshTokenResponse{
			Token: testutls.MockToken,
		},
		wantErr: false,
		init: refereshTokenPatches{
			configPatch: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
					return &config.Configuration{}, nil
				})
			},
			jwtPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, nil
				})
			},
			tokenPatch: func() *gomonkey.Patches {
				tg := jwt.Service{}
				return gomonkey.ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return "token", nil
					})
			},
		},
		initMocks: func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id", "email", "token", "role_id"}).
				AddRow(1, testutls.MockEmail, testutls.MockToken, 1)
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" WHERE (token=$1) LIMIT 1;`)).
				WithArgs().
				WillReturnRows(rows)
		},
	}
}
func loadRefereshTokenCases() []refereshTokenType {
	return []refereshTokenType{
		refreshTokenInvalidCase(),
		refreshTokenErrorFromConfigCase(),
		refereshTokenErrorFromJWTCase(),
		refereshTokenErrorFromGenerateTokenCase(),
		refreshTokenSuccessCase(),
	}
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
				tt.initMocks(mock)
				// Handle the case where there is an error loading the config
				configpatch := tt.init.configPatch()
				//initialize a jwt service
				// Handle the case where there is an error creating the JWT service
				patchJWT := tt.init.jwtPatch()
				// Handle the case where there is an error form token generation service
				patchGenerateToken := tt.init.tokenPatch()
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
				configpatch.Reset()
				patchJWT.Reset()
				patchGenerateToken.Reset()
			},
		)
	}
}
