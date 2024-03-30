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

	"github.com/agiledragon/gomonkey/v2"
	null "github.com/volatiletech/null/v8"

	fm "go-template/gqlmodels"

	"github.com/stretchr/testify/assert"
)

type createRoleType struct {
	name     string
	req      fm.RoleCreateInput
	wantResp *fm.RolePayload
	wantErr  bool
	init     func() *gomonkey.Patches
}

func errorFromRedisCase() createRoleType {
	return createRoleType{
		name:     ErrorFromRedisCache,
		req:      fm.RoleCreateInput{},
		wantResp: &fm.RolePayload{},
		wantErr:  true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(rediscache.GetUser,
				func(userID int, ctx context.Context) (*models.User, error) {
					return nil, errors.New("redis cache")
				})
		},
	}
}

func errorFromGetRoleCase() createRoleType {
	return createRoleType{
		name:     ErrorFromGetRole,
		req:      fm.RoleCreateInput{},
		wantResp: &fm.RolePayload{},
		wantErr:  true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(rediscache.GetRole,
				func(roleID int, ctx context.Context) (*models.Role, error) {
					return nil, errors.New("data")
				})
		},
	}
}
func errorUnauthorizedUserCase() createRoleType {
	return createRoleType{
		name: ErrorUnauthorizedUser,
		req: fm.RoleCreateInput{
			Name:        UserRoleName,
			AccessLevel: int(constants.UserRole),
		},
		wantResp: &fm.RolePayload{},
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(rediscache.GetRole,
				func(roleID int, ctx context.Context) (*models.Role, error) {
					return &models.Role{
						AccessLevel: int(constants.UserRole),
						Name:        UserRoleName,
					}, nil
				})
		},
	}
}
func successCase() createRoleType {
	return createRoleType{
		name: SuccessCase,
		req: fm.RoleCreateInput{
			Name:        UserRoleName,
			AccessLevel: int(constants.UserRole),
		},
		wantResp: &fm.RolePayload{Role: &fm.Role{
			AccessLevel: int(constants.UserRole),
			Name:        UserRoleName,
		}},
		wantErr: false,
		init:    func() *gomonkey.Patches { return nil },
	}
}

func errorFromCreateRoleCase() createRoleType {
	return createRoleType{
		name: ErrorFromCreateRole,
		req: fm.RoleCreateInput{
			Name:        UserRoleName,
			AccessLevel: int(constants.UserRole),
		},
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.CreateRole,
				func(role models.Role, ctx context.Context) (models.Role, error) {
					return models.Role{}, errors.New("error")
				})
		},
	}
}
func loadTestCases() []createRoleType {
	return []createRoleType{
		errorFromRedisCase(),
		errorFromGetRoleCase(),
		errorUnauthorizedUserCase(),
		successCase(),
		errorFromCreateRoleCase(),
	}
}
func applyUserIdPatch() *gomonkey.Patches {
	return gomonkey.ApplyFunc(auth.UserIDFromContext,
		func(ctx context.Context) int {
			return 1
		})
}
func applyGetUserPatch() *gomonkey.Patches {
	return gomonkey.ApplyFunc(rediscache.GetUser,
		func(userID int, ctx context.Context) (*models.User, error) {
			return &models.User{
				RoleID: null.IntFrom(1),
			}, nil
		})
}
func applyGetRolePatch() *gomonkey.Patches {
	return gomonkey.ApplyFunc(rediscache.GetRole,
		func(roleID int, ctx context.Context) (*models.Role, error) {
			return &models.Role{
				AccessLevel: int(constants.SuperAdminRole),
				Name:        SuperAdminRoleName,
			}, nil
		})
}
func applyCreateRolePatch() *gomonkey.Patches {
	return gomonkey.ApplyFunc(daos.CreateRole,
		func(role models.Role, ctx context.Context) (models.Role, error) {
			return models.Role{
				AccessLevel: int(constants.SuperAdminRole),
				Name:        SuperAdminRoleName,
			}, nil
		})
}

// TestCreateRole tests the CreateRole mutation function.
func TestCreateRole(
	t *testing.T,
) {
	// Define test cases, each case has a name, request input, expected response, and error.
	cases := loadTestCases()
	// Create a new resolver instance.
	resolver1 := resolver.Resolver{}

	// Loop through each test case.
	for _, tt := range cases {
		// Mocking rediscache.GetUserID function
		patchUserID := applyUserIdPatch()
		// Mocking rediscache.GetUser function
		patchGetUser := applyGetUserPatch()
		// Mocking rediscache.GetRole function
		patchGetRole := applyGetRolePatch()
		// Mocking daos.CreateRole function
		patchCreateRole := applyCreateRolePatch()
		// Defer resetting of the monkey patches.
		defer patchUserID.Reset()
		defer patchGetUser.Reset()
		defer patchGetRole.Reset()
		defer patchCreateRole.Reset()
		t.Run(tt.name,
			func(t *testing.T) {
				// Apply additional monkey patches based on test case name.
				patch := tt.init()
				defer patch.Reset()
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
