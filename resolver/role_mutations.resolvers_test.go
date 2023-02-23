package resolver_test

import (
	"context"
	"errors"
	"go-template/daos"
	"go-template/internal/constants"
	"go-template/internal/middleware/auth"
	"go-template/models"
	"go-template/pkg/utl/rediscache"
	"go-template/resolver"
	"testing"

	"github.com/volatiletech/null/v8"

	"github.com/agiledragon/gomonkey/v2"

	fm "go-template/gqlmodels"

	"github.com/stretchr/testify/assert"
)

// TestCreateRole tests the CreateRole mutation function.
func TestCreateRole(
	t *testing.T,
) {

	// Define test cases, each case has a name, request input, expected response, and error.
	cases := []struct {
		name     string
		req      fm.RoleCreateInput
		wantResp *fm.RolePayload
		wantErr  bool
	}{
		{
			name:     constants.ErrorFromRedisCache,
			req:      fm.RoleCreateInput{},
			wantResp: &fm.RolePayload{},
			wantErr:  true,
		},
		{
			name:     constants.ErrorFromGetRole,
			req:      fm.RoleCreateInput{},
			wantResp: &fm.RolePayload{},
			wantErr:  true,
		},
		{
			name: constants.ErrorUnauthorizedUser,
			req: fm.RoleCreateInput{
				Name:        "Role",
				AccessLevel: int(constants.UserRole),
			},
			wantResp: &fm.RolePayload{},
		},

		{
			name: constants.SuccessCase,
			req: fm.RoleCreateInput{
				Name:        "Role",
				AccessLevel: int(constants.UserRole),
			},
			wantResp: &fm.RolePayload{Role: &fm.Role{

				AccessLevel: int(constants.UserRole),
				Name:        "Role",
			}},
		},
		{
			name: constants.ErrorFromCreateRole,
			req: fm.RoleCreateInput{
				Name:        "UserRole",
				AccessLevel: int(constants.UserRole),
			},
			wantErr: true,
		},
	}
	// Create a new resolver instance.
	resolver1 := resolver.Resolver{}

	// Loop through each test case.
	for _, tt := range cases {

		// Mocking rediscache.GetUserID function
		patchUserID := gomonkey.ApplyFunc(auth.UserIDFromContext,
			func(ctx context.Context) int {
				return 1
			})

		// Mocking rediscache.GetUser function
		patchGetUser := gomonkey.ApplyFunc(rediscache.GetUser,
			func(userID int, ctx context.Context) (*models.User, error) {
				return &models.User{
					RoleID: null.IntFrom(1),
				}, nil
			})

		// Mocking rediscache.GetRole function
		patchGetRole := gomonkey.ApplyFunc(rediscache.GetRole,
			func(roleID int, ctx context.Context) (*models.Role, error) {
				return &models.Role{
					AccessLevel: 100,
					Name:        "SuperAdminRole",
				}, nil
			})

		// Mocking daos.CreateRole function
		patchCreateRole := gomonkey.ApplyFunc(daos.CreateRole,
			func(role models.Role, ctx context.Context) (models.Role, error) {
				return models.Role{
					AccessLevel: int(constants.UserRole),
					Name:        "Role",
				}, nil
			})

		// Defer resetting of the monkey patches.
		defer patchUserID.Reset()
		defer patchGetUser.Reset()
		defer patchGetRole.Reset()
		defer patchCreateRole.Reset()
		t.Run(tt.name,
			func(t *testing.T) {

				// Apply additional monkey patches based on test case name.
				if tt.name == constants.ErrorFromRedisCache {
					patchGetUser := gomonkey.ApplyFunc(rediscache.GetUser,
						func(userID int, ctx context.Context) (*models.User, error) {
							return nil, errors.New("redis cache")
						})
					defer patchGetUser.Reset()
				}

				if tt.name == constants.ErrorFromGetRole {
					patchGetRole := gomonkey.ApplyFunc(rediscache.GetRole,
						func(roleID int, ctx context.Context) (*models.Role, error) {
							return nil, errors.New("data")
						})
					defer patchGetRole.Reset()
				}
				if tt.name == constants.ErrorUnauthorizedUser {
					patchGetRole := gomonkey.ApplyFunc(rediscache.GetRole,
						func(roleID int, ctx context.Context) (*models.Role, error) {
							return &models.Role{
								AccessLevel: int(constants.UserRole),
								Name:        "Role",
							}, nil
						})
					defer patchGetRole.Reset()
				}

				if tt.name == constants.ErrorFromCreateRole {
					patchCreateRole := gomonkey.ApplyFunc(daos.CreateRole,
						func(role models.Role, ctx context.Context) (models.Role, error) {
							return models.Role{}, errors.New("error")
						})

					defer patchCreateRole.Reset()
				}

				// Create a new context
				c := context.Background()

				// Call the resolver function
				response, err := resolver1.Mutation().CreateRole(c, tt.req)

				// Check if the error matches the expected error
				if tt.wantErr {
					assert.NotNil(t, err)
				}

				// Check if the response matches the expected response
				assert.Equal(t, tt.wantResp, response)

			})
	}
}
