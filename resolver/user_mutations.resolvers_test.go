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
func getCreateUserTestCase() []struct {
	name     string
	req      fm.UserCreateInput
	wantResp *fm.User
	wantErr  bool
} {
	cases := []struct {
		name     string
		req      fm.UserCreateInput
		wantResp *fm.User
		wantErr  bool
	}{
		{
			name:    ErrorFromCreateUser,
			req:     fm.UserCreateInput{},
			wantErr: true,
		},
		{
			name:    ErrorFromThrottleCheck,
			req:     fm.UserCreateInput{},
			wantErr: true,
		},
		{
			name:    ErrorFromConfig,
			req:     fm.UserCreateInput{},
			wantErr: true,
		},
		{
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
		},
	}
	return cases
}
func TestCreateUser(t *testing.T) {
	cases := getCreateUserTestCase()
	resolver := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mock, cleanup, _ := testutls.SetupMockDB(t)

			if tt.name == ErrorFromCreateUser {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			} else {
				expectInsertUser(mock, *testutls.MockUser())
			}

			if tt.name == ErrorFromThrottleCheck {
				patch := gomonkey.ApplyFunc(throttle.Check, func(ctx context.Context, limit int, dur time.Duration) error {
					return fmt.Errorf("Internal error")
				})
				defer patch.Reset()
			}

			if tt.name == ErrorFromConfig {
				patch := gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
					return nil, fmt.Errorf("error in loading config")
				})
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

func GetUpdateUserTestCase() []struct {
	name     string
	req      *fm.UserUpdateInput
	wantResp *fm.User
	wantErr  bool
} {
	cases := []struct {
		name     string
		req      *fm.UserUpdateInput
		wantResp *fm.User
		wantErr  bool
	}{
		{
			name:    ErrorFindingUser,
			req:     &fm.UserUpdateInput{},
			wantErr: true,
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
		},
	}
	return cases
}
func TestUpdateUser(
	t *testing.T,
) {
	cases := GetUpdateUserTestCase()
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				if tt.name == ErrorUpdateUser {
					patch := gomonkey.ApplyFunc(daos.UpdateUser,
						func(user models.User, ctx context.Context) (models.User, error) {
							return user, fmt.Errorf("error for update user")
						})
					defer patch.Reset()
				}
				err := config.LoadEnv()
				if err != nil {
					log.Fatal(err)
				}
				mock, cleanup, _ := testutls.SetupMockDB(t)

				if tt.name == ErrorFindingUser {
					mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "users"`)).WithArgs().WillReturnError(fmt.Errorf(""))
				}
				rows := sqlmock.NewRows([]string{"first_name"}).AddRow(testutls.MockUser().FirstName)
				mock.ExpectQuery(regexp.QuoteMeta(`select * from "users"`)).WithArgs(0).WillReturnRows(rows)
				// update users with new information
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).WillReturnResult(result)
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Mutation().UpdateUser(ctx, tt.req)
				if tt.wantResp != nil &&
					response != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
				cleanup()
			},
		)
	}
}

func GetDeleteTestCases() []struct {
	name     string
	wantResp *fm.UserDeletePayload
	wantErr  bool
} {
	cases := []struct {
		name     string
		wantResp *fm.UserDeletePayload
		wantErr  bool
	}{
		{
			name:    ErrorFindingUser,
			wantErr: true,
		},
		{
			name:    ErrorDeleteUser,
			wantErr: true,
		},
		{
			name: SuccessCase,
			wantResp: &fm.UserDeletePayload{
				ID: "0",
			},
			wantErr: false,
		},
	}
	return cases
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
				if tt.name == ErrorDeleteUser {
					patch := gomonkey.ApplyFunc(daos.DeleteUser,
						func(user models.User, ctx context.Context) (int64, error) {
							return 0, fmt.Errorf("error for delete user")
						})
					defer patch.Reset()
				}
				mock, cleanup, _ := testutls.SetupMockDB(t)
				if tt.name == ErrorFindingUser {
					mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
						WithArgs().
						WillReturnError(fmt.Errorf(""))
				}
				// get user by id
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where "id"=$1`)).
					WithArgs().
					WillReturnRows(rows)
				// delete user
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "id"=$1`)).
					WillReturnResult(result)
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Mutation().
					DeleteUser(ctx)
				if tt.wantResp != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
				cleanup()
			},
		)
	}
}
