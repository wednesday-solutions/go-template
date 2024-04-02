package resolver_test

import (
	"context"
	"fmt"
	"go-template/daos"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/volatiletech/null/v8"

	fm "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/internal/jwt"
	"go-template/internal/service"
	"go-template/models"
	"go-template/pkg/utl/resultwrapper"
	"go-template/pkg/utl/secure"
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
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					return nil, fmt.Errorf(ErrorMsgFindingUser)
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
			// mock FindUserByUserName with the proper password, and active state
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, nil
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
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, nil
				}).
				ApplyFunc(config.Load, func() (*config.Configuration, error) {
					return nil, fmt.Errorf("error in loading config")
				})
		},
	}
}

func errorWhileGeneratingToken() loginType {
	return loginType{
		name: ErrorFromJwt,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(true)
					return user, nil
				}).
				ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, nil
				}).
				ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return "", fmt.Errorf(ErrorMsgFromJwt)
					})
		},
	}
}
func errorUpdateUserCase() loginType {
	err := fmt.Errorf(ErrorMsgfromUpdateUser)
	return loginType{
		name: ErrorUpdateUser,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     err,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			tg := jwt.Service{}
			sec := secure.Service{}
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(true)
					return user, nil
				}).
				ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, nil
				}).
				ApplyFunc(service.Secure, func(cfg *config.Configuration) secure.Service {
					return sec
				}).
				ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return "", nil
					}).
				ApplyMethod(reflect.TypeOf(sec), "Token",
					func(secure.Service, string) string {
						return "refreshToken"
					}).
				ApplyFunc(daos.UpdateUser,
					func(u models.User, ctx context.Context) (models.User, error) {
						return models.User{}, err
					})
		},
	}
}
func loginSuccessCase() loginType {
	jwtToken := "jwttokenstring"
	return loginType{
		name: SuccessCase,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantResp: &fm.LoginResponse{
			Token:        jwtToken,
			RefreshToken: TestToken,
		},
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			tg := jwt.Service{}
			sec := secure.Service{}
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(true)
					return user, nil
				}).
				ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, nil
				}).
				ApplyFunc(service.Secure, func(cfg *config.Configuration) secure.Service {
					return sec
				}).
				ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return jwtToken, nil
					}).
				ApplyMethod(reflect.TypeOf(sec), "Token",
					func(secure.Service, string) string {
						return TestToken
					}).
				ApplyFunc(daos.UpdateUser,
					func(u models.User, ctx context.Context) (models.User, error) {
						return *testutls.MockUser(), nil
					})
		},
	}
}
func errorWhileCreatingJWTService() loginType {
	err := fmt.Errorf("error in creating auth service")
	return loginType{
		name: "Error while creating a JWT Service",
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     err,
		init: func(mock sqlmock.Sqlmock) *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, nil
				}).ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				// mock service.JWT
				return jwt.Service{}, err
			})
		},
	}
}
func loadLoginTestCases() []loginType {
	return []loginType{
		errorFindingUserCase(),
		errorFromConfigCase(),
		errorPasswordValidationCase(),
		errorActiveStatusCase(),
		errorWhileCreatingJWTService(),
		errorWhileGeneratingToken(),
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
				mock, cleanup, _ := testutls.SetupMockDB(t)
				defer cleanup()
				patch := tt.init(mock)
				c := context.Background()
				// Call the login mutation with the given arguments and check the response and error against the expected values
				response, err := resolver1.Mutation().Login(c, tt.req.UserName, tt.req.Password)
				if tt.wantResp != nil &&
					response != nil {
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
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					return nil, fmt.Errorf(ErrorMsgFindingUser)
				})
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
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(tg.GenerateToken, func(u *models.User) (string, error) {
				return "", nil
			})
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
			// mock FindUserByUserName with the proper password, and active state
			return gomonkey.ApplyFunc(daos.FindUserByID,
				func(userID int, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, fmt.Errorf(ErrorInsecurePassword)
				})
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
			return gomonkey.ApplyFunc(daos.FindUserByID,
				func(userID int, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, fmt.Errorf(ErrorInsecurePassword)
				}).ApplyFunc(daos.UpdateUser,
				func(user models.User, ctx context.Context) (*models.User, error) {
					return nil, fmt.Errorf(ErrorUpdateUser)
				})
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
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return nil, fmt.Errorf("error in loading config")
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
			sec := secure.Service{}
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return nil, nil
			}).ApplyFunc(service.Secure, func(cfg *config.Configuration) secure.Service {
				return sec
			}).ApplyMethod(reflect.TypeOf(sec), "Password", func(secure.Service, string, ...string) bool {
				return true
			}).ApplyFunc(daos.FindUserByID,
				func(userID int, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, nil
				}).ApplyFunc(daos.UpdateUser,
				func(user models.User, ctx context.Context) (models.User, error) {
					return *testutls.MockUser(), nil
				})
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
				mock, cleanup, _ := testutls.SetupMockDB(t)
				defer cleanup()
				tt.init(mock)
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
	init     func() *gomonkey.Patches
}

func refreshTokenInvalidCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorInvalidToken,
		req:     TestToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsginvalidToken),
		init: func() *gomonkey.Patches {
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return &config.Configuration{}, nil
			}).ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, nil
			}).ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
				func(jwt.Service, *models.User) (string, error) {
					return "", fmt.Errorf(ErrorMsginvalidToken)
				}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			})
		},
	}
}
func refreshTokenErrorFromConfigCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorFromConfig,
		req:     ReqToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromConfig),
		init: func() *gomonkey.Patches {
			// tg := jwt.Service{}
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return nil, fmt.Errorf(ErrorFromConfig)
			}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			})
		},
	}
}

func refereshTokenerrorWhileGeneratingToken() refereshTokenType {
	return refereshTokenType{
		name:    ErrorFromJwt,
		req:     ReqToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromJwt),
		init: func() *gomonkey.Patches {
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return &config.Configuration{}, nil
			}).ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, fmt.Errorf(ErrorMsgFromJwt)
			}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			})
		},
	}
}

func refereshTokenErrorFromGenerateTokenCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorFromGenerateToken,
		req:     ReqToken,
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
		init: func() *gomonkey.Patches {
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return &config.Configuration{}, nil
			}).ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, nil
			}).ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
				func(jwt.Service, *models.User) (string, error) {
					return "", resultwrapper.ErrUnauthorized
				}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			})
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
		init: func() *gomonkey.Patches {
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return &config.Configuration{}, nil
			}).ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, nil
			}).ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
				func(jwt.Service, *models.User) (string, error) {
					return "", nil
				}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			})
		},
	}
}
func loadRefereshTokenCases() []refereshTokenType {
	return []refereshTokenType{
		refreshTokenInvalidCase(),
		refreshTokenErrorFromConfigCase(),
		refereshTokenerrorWhileGeneratingToken(),
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
				// Handle the case where authentication token is invalid
				patches := tt.init()
				defer patches.Reset()
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
					fmt.Println(err.Error(), tt.err.Error(), strings.Contains(err.Error(), tt.err.Error()))
					// Assert that the expected error value matches the actual error value
					assert.Equal(t, true, strings.Contains(err.Error(), tt.err.Error()))
				}
			},
		)
	}
}
