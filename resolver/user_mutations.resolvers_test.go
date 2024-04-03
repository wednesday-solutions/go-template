package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"go-template/daos"
	fm "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/internal/service"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/secure"
	"go-template/pkg/utl/throttle"
	"go-template/resolver"
	"go-template/testutls"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(
	v driver.Value,
) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyString struct{}

func (a AnyString) Match(
	v driver.Value,
) bool {
	_, ok := v.(string)
	return ok
}

type createUserType struct {
	name     string
	req      fm.UserCreateInput
	wantResp *fm.User
	wantErr  bool
	init     func() *gomonkey.Patches
}

func errorFromCreateUserCase() createUserType {
	return createUserType{
		name:    ErrorFromCreateUser,
		req:     fm.UserCreateInput{},
		wantErr: true,
		init: func() *gomonkey.Patches {
			sec := secure.Service{}
			return gomonkey.ApplyFunc(daos.CreateUser, func(user models.User, ctx context.Context) (models.User, error) {
				return *testutls.MockUser(), fmt.Errorf("error")
			}).ApplyFunc(service.Secure, func(cfg *config.Configuration) secure.Service {
				return sec
			}).ApplyFunc(throttle.Check, func(ctx context.Context, limit int, dur time.Duration) error {
				return nil
			})
		},
	}
}

func errorFromThrottleCheck() createUserType {
	return createUserType{
		name:    ErrorFromThrottleCheck,
		req:     fm.UserCreateInput{},
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(throttle.Check, func(ctx context.Context, limit int, dur time.Duration) error {
				return fmt.Errorf("Internal error")
			})
		},
	}
}
func errorFromCreateUserConfigCase() createUserType {
	return createUserType{
		name:    ErrorFromConfig,
		req:     fm.UserCreateInput{},
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return nil, fmt.Errorf("error in loading config")
			})
		},
	}
}

func createUserSuccessCase() createUserType {
	return createUserType{
		name: SuccessCase,
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
		init: func() *gomonkey.Patches {
			sec := secure.Service{}
			return gomonkey.ApplyFunc(daos.CreateUser, func(user models.User, ctx context.Context) (models.User, error) {
				return models.User{
					ID:                 testutls.MockUser().ID,
					Email:              testutls.MockUser().Email,
					FirstName:          testutls.MockUser().FirstName,
					LastName:           testutls.MockUser().LastName,
					Username:           testutls.MockUser().Username,
					Mobile:             testutls.MockUser().Mobile,
					Address:            testutls.MockUser().Address,
					Active:             testutls.MockUser().Active,
					LastLogin:          testutls.MockUser().LastLogin,
					LastPasswordChange: testutls.MockUser().LastPasswordChange,
					DeletedAt:          testutls.MockUser().DeletedAt,
					UpdatedAt:          testutls.MockUser().UpdatedAt,
				}, nil
			}).ApplyFunc(throttle.Check, func(ctx context.Context, limit int, dur time.Duration) error {
				return nil
			}).ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return nil, nil
			}).ApplyFunc(service.Secure, func(cfg *config.Configuration) secure.Service {
				return sec
			})
		},
	}
}
func getCreateUserTestCase() []createUserType {
	cases := []createUserType{
		errorFromCreateUserCase(),
		errorFromThrottleCheck(),
		errorFromCreateUserConfigCase(),
		createUserSuccessCase(),
	}
	return cases
}
func TestCreateUser(t *testing.T) {
	cases := getCreateUserTestCase()
	resolver := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			patch := tt.init()
			time.Sleep(time.Duration(100000))
			response, err := resolver.Mutation().CreateUser(context.Background(), tt.req)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
			if patch != nil {
				patch.Reset()
			}
		})
	}
}

type updateUserType struct {
	name     string
	req      *fm.UserUpdateInput
	wantResp *fm.User
	wantErr  bool
	init     func() *gomonkey.Patches
}

func loadUpdateUserTestCases() []updateUserType {
	return []updateUserType{
		{
			name:    ErrorFindingUser,
			req:     &fm.UserUpdateInput{},
			wantErr: true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindUserByID, func(userID int, ctx context.Context) (*models.User, error) {
					return nil, fmt.Errorf("")
				})
			},
		},
		{
			name: ErrorUpdateUser,
			req: &fm.UserUpdateInput{
				FirstName: &testutls.MockUser().FirstName.String,
				LastName:  &testutls.MockUser().LastName.String,
				Mobile:    &testutls.MockUser().Mobile.String,
				Address:   &testutls.MockUser().Address.String,
			},
			wantErr: true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.UpdateUser,
					func(models.User, context.Context) (models.User, error) {
						return *testutls.MockUser(), fmt.Errorf("error for update user")
					})
			},
		},
		{
			name: SuccessCase,
			req: &fm.UserUpdateInput{
				FirstName: &testutls.MockUser().FirstName.String,
				LastName:  &testutls.MockUser().LastName.String,
				Mobile:    &testutls.MockUser().Mobile.String,
				Address:   &testutls.MockUser().Address.String,
			},
			wantResp: &fm.User{
				ID:        "1",
				FirstName: convert.NullDotStringToPointerString(testutls.MockUser().FirstName),
				LastName:  convert.NullDotStringToPointerString(testutls.MockUser().LastName),
				Username:  convert.NullDotStringToPointerString(testutls.MockUser().Username),
				Mobile:    convert.NullDotStringToPointerString(testutls.MockUser().Mobile),
				Address:   convert.NullDotStringToPointerString(testutls.MockUser().Address),
				Email:     convert.NullDotStringToPointerString(testutls.MockUser().Email),
				Active:    &testutls.MockUser().Active.Bool,
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindUserByID, func(userID int, ctx context.Context) (*models.User, error) {
					return testutls.MockUser(), nil
				}).ApplyFunc(daos.UpdateUser, func(user models.User, ctx context.Context) (models.User, error) {
					return *testutls.MockUser(), nil
				})
			},
		},
	}
}

func TestUpdateUser(
	t *testing.T,
) {
	cases := loadUpdateUserTestCases()
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				patches := tt.init()
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Mutation().UpdateUser(ctx, tt.req)
				if tt.wantResp != nil &&
					response != nil {
					assert.Equal(t, tt.wantResp, response)
				} else {
					assert.Equal(t, tt.wantErr, err != nil)
				}
				if patches != nil {
					patches.Reset()
				}
			},
		)
	}
}

type deleteUserType struct {
	name     string
	wantResp *fm.UserDeletePayload
	wantErr  bool
	init     func() *gomonkey.Patches
}

func errorFindinguserCaseDelete() deleteUserType {
	return deleteUserType{
		name:    ErrorFindingUser,
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByID, func(userID int, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), fmt.Errorf("")
			})
		},
	}
}

func errorDeleteUserCase() deleteUserType {
	return deleteUserType{
		name:    ErrorDeleteUser,
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByID, func(userID int, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), fmt.Errorf("")
			}).ApplyFunc(daos.DeleteUser,
				func(user models.User, ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("error for delete user")
				})
		},
	}
}

func deleteUserSuccessCase() deleteUserType {
	return deleteUserType{
		name: SuccessCase,
		wantResp: &fm.UserDeletePayload{
			ID: "0",
		},
		wantErr: false,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByID, func(userID int, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			}).ApplyFunc(daos.DeleteUser,
				func(user models.User, ctx context.Context) (int64, error) {
					return 0, nil
				})
		},
	}
}

func GetDeleteTestCases() []deleteUserType {
	return []deleteUserType{
		errorFindinguserCaseDelete(),
		errorDeleteUserCase(),
		deleteUserSuccessCase(),
	}
}

func TestDeleteUser(
	t *testing.T,
) {
	cases := GetDeleteTestCases()
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				patch := tt.init()
				// get user by id
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Mutation().DeleteUser(ctx)
				if tt.wantResp != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
				if patch != nil {
					patch.Reset()
				}
			},
		)
	}
}
