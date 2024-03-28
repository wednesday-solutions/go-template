package resolver_test

import (
	"context"
	"errors"
	"fmt"
	"go-template/pkg/utl/cnvrttogql"
	"regexp"
	"testing"

	"go-template/gqlmodels"
	fm "go-template/gqlmodels"
	"go-template/models"
	//	"go-template/pkg/utl/cnvrttogql"
	"go-template/pkg/utl/rediscache"
	"go-template/resolver"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/assert"
)

type args struct {
	user *models.User
}

type queryMeTestArgs struct {
	name     string
	wantResp *fm.User
	wantErr  bool
	args     args
	init     func(args) *Patches
}

func initializeCases() []queryMeTestArgs {
	return []queryMeTestArgs{
		{
			name:     SuccessCase,
			args:     args{user: testutls.MockUser()},
			wantResp: cnvrttogql.UserToGraphQlUser(testutls.MockUser(), 4),
			init: func(args args) *Patches {
				return gomonkey.ApplyFunc(rediscache.GetUser,
					func(userID int, ctx context.Context) (*models.User, error) {
						return args.user, nil
					})
			},
		},
		{
			name:     ErrorFromRedisCache,
			args:     args{user: testutls.MockUser()},
			wantErr:  true,
			wantResp: &gqlmodels.User{},
			init: func(args args) *Patches {
				return gomonkey.ApplyFunc(rediscache.GetUser,
					func(userID int, ctx context.Context) (*models.User, error) {
						return nil, errors.New("redis cache")
					})
			},
		},
	}
}

func TestMe(t *testing.T) {
	cases := initializeCases()
	_, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()
	resolver1 := &resolver.Resolver{}
	for _, tt := range cases {
		patches := tt.init(tt.args)
		response, err := resolver1.Query().Me(context.TODO())
		if tt.wantResp != nil && response != nil {
			assert.Equal(t, tt.wantResp, response)
		}
		assert.Equal(t, tt.wantErr, err != nil)
		if patches != nil {
			patches.Reset()
		}
	}
}

type queryUsersArgs struct {
	name       string
	pagination *fm.UserPagination
	wantResp   []*models.User
	wantErr    bool
	init       func(sqlmock.Sqlmock)
}

// TestUsers is a unit test function for testing user queries.
func initializeTestCases() []queryUsersArgs {
	return []queryUsersArgs{
		{
			name:    ErrorFindingUser,
			wantErr: true,
			init: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users";`)).
					WithArgs().
					WillReturnError(fmt.Errorf(""))
			},
		},
		{
			name:    "Paginated Users are returned when the request is paginated",
			wantErr: false,
			pagination: &fm.UserPagination{
				Limit: 1,
				Page:  1,
			},
			wantResp: testutls.MockUsers(),
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "address"}).
					AddRow(testutls.MockID, testutls.MockEmail, "First", "Last", "+911234567890", "username", "22 Jump Street")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users" LIMIT 1 OFFSET 1;`)).WithArgs().WillReturnRows(rows)

				rowCount := sqlmock.NewRows([]string{"count"}).
					AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM "users"`)).
					WithArgs().
					WillReturnRows(rowCount)
			},
		},
		{
			name:     "Successfully fetches paginated users even without pagination request payload",
			wantErr:  false,
			wantResp: testutls.MockUsers(),
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "email", "first_name", "last_name", "mobile", "username", "address"}).
					AddRow(testutls.MockID, testutls.MockEmail, "First", "Last", "+911234567890", "username", "22 Jump Street")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users".* FROM "users";`)).WithArgs().WillReturnRows(rows)

				rowCount := sqlmock.NewRows([]string{"count"}).
					AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM "users";`)).
					WithArgs().
					WillReturnRows(rowCount)
			},
		},
	}
}

func executeQuery(resolver1 *resolver.Resolver,
	ctx context.Context, pagination *fm.UserPagination) (*gqlmodels.UsersPayload, error) {
	return resolver1.Query().Users(ctx, pagination)
}

func TestUsers(
	t *testing.T,
) {
	cases := initializeTestCases()
	resolver1 := resolver.Resolver{}

	for _, tt := range cases {
		mock, cleanup, _ := testutls.SetupMockDB(t)
		tt.init(mock)
		response, err := executeQuery(&resolver1, context.Background(), tt.pagination)
		if tt.wantResp != nil && response != nil {
			assert.Equal(t, len(tt.wantResp), len(response.Users))
		}
		assert.Equal(t, tt.wantErr, err != nil)
		cleanup()
	}
}
