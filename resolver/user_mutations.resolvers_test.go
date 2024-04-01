package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"go-template/daos"
	fm "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/throttle"
	"go-template/resolver"
	"go-template/testutls"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func expectInsertUser(mock sqlmock.Sqlmock, mockUser models.User) {
	rows := sqlmock.NewRows([]string{
		"id", "mobile", "address", "active", "last_login", "last_password_change", "token", "deleted_at",
	}).AddRow(
		mockUser.ID, mockUser.Mobile, mockUser.Address, mockUser.Active,
		mockUser.LastLogin, mockUser.LastPasswordChange, mockUser.Token, mockUser.DeletedAt,
	)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(
			mockUser.FirstName, mockUser.LastName, mockUser.Username, AnyString{}, mockUser.Email,
			mockUser.RoleID, AnyTime{}, AnyTime{},
		).
		WillReturnRows(rows)
}

type createUserType struct {
	name      string
	req       fm.UserCreateInput
	wantResp  *fm.User
	wantErr   bool
	init      func() *gomonkey.Patches
	initmocks func(mock sqlmock.Sqlmock, mockUser models.User)
}

func errorFromCreateUserCase() createUserType {
	return createUserType{
		name:    ErrorFromCreateUser,
		req:     fm.UserCreateInput{},
		wantErr: true,
		init:    func() *gomonkey.Patches { return nil },
		initmocks: func(mock sqlmock.Sqlmock, mockUser models.User) {
			mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
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
		initmocks: expectInsertUser,
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
		initmocks: expectInsertUser,
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
		wantErr:   false,
		init:      func() *gomonkey.Patches { return nil },
		initmocks: expectInsertUser,
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
			mock, cleanup, _ := testutls.SetupMockDB(t)
			tt.initmocks(mock, *testutls.MockUser())
			patch := tt.init()
			if patch != nil {
				defer patch.Reset()
			}
			response, err := resolver.Mutation().CreateUser(context.Background(), tt.req)
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
			cleanup()
		})
	}
}

type updateUserType struct {
	name      string
	req       *fm.UserUpdateInput
	wantResp  *fm.User
	wantErr   bool
	init      func()
	initMocks func(mock sqlmock.Sqlmock)
}

func loadUpdateUserTestCases() []updateUserType {
	return []updateUserType{
		{
			name:    ErrorFindingUser,
			req:     &fm.UserUpdateInput{},
			wantErr: true,
			init:    func() {},
			initMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "users"`)).WithArgs().WillReturnError(fmt.Errorf(""))
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
			init: func() {
				gomonkey.ApplyFunc(daos.UpdateUser,
					func(user models.User, ctx context.Context) (models.User, error) {
						return user, fmt.Errorf("error for update user")
					})
			},
			initMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"first_name"}).AddRow(testutls.MockUser().FirstName)
				mock.ExpectQuery(regexp.QuoteMeta(`select * from "users"`)).WithArgs(0).WillReturnRows(rows)
				// update users with new information
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).WillReturnResult(result)
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
				ID:        "0",
				FirstName: &testutls.MockUser().FirstName.String,
				LastName:  &testutls.MockUser().LastName.String,
				Mobile:    &testutls.MockUser().Mobile.String,
				Address:   &testutls.MockUser().Address.String,
			},
			wantErr: false,
			init:    func() {},
			initMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"first_name"}).AddRow(testutls.MockUser().FirstName)
				mock.ExpectQuery(regexp.QuoteMeta(`select * from "users"`)).WithArgs(0).WillReturnRows(rows)
				// update users with new information
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).WillReturnResult(result)
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
				tt.init()
				err := config.LoadEnv()
				if err != nil {
					log.Fatal(err)
				}
				mock, cleanup, _ := testutls.SetupMockDB(t)
				oldDB := boil.GetDB()
				defer cleanup()
				boil.SetDB(oldDB)
				tt.initMocks(mock)
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Mutation().UpdateUser(ctx, tt.req)
				if tt.wantResp != nil &&
					response != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}

type deleteUserType struct {
	name      string
	wantResp  *fm.UserDeletePayload
	wantErr   bool
	init      func() *gomonkey.Patches
	initMocks func(mock sqlmock.Sqlmock)
}

func errorFindinguserCaseDelete() deleteUserType {
	return deleteUserType{
		name:    ErrorFindingUser,
		wantErr: true,
		init:    func() *gomonkey.Patches { return nil },
		initMocks: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnError(fmt.Errorf(""))
			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnRows(rows)
			// delete user
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "id"=$1`)).
				WillReturnResult(result)
		},
	}
}

func errorDeleteUserCase() deleteUserType {
	return deleteUserType{
		name:    ErrorDeleteUser,
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.DeleteUser,
				func(user models.User, ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("error for delete user")
				})
		},
		initMocks: func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnRows(rows)
			// delete user
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "id"=$1`)).
				WillReturnResult(result)
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
		init:    func() *gomonkey.Patches { return nil },
		initMocks: func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(1)
			mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
				WithArgs().
				WillReturnRows(rows)
			// delete user
			result := driver.Result(driver.RowsAffected(1))
			mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "id"=$1`)).
				WillReturnResult(result)
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
				if patch != nil {
					defer patch.Reset()
				}
				mock, cleanup, _ := testutls.SetupMockDB(t)
				// get user by id
				tt.initMocks(mock)
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Mutation().DeleteUser(ctx)
				if tt.wantResp != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
				cleanup()
			},
		)
	}
}
