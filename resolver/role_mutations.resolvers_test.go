package resolver_test

import (
	"context"
	"errors"
	"go-template/daos"
	"go-template/internal/constants"
	"go-template/models"
	"go-template/pkg/utl/rediscache"
	"go-template/resolver"
	"go-template/testutls"
	"testing"

	"github.com/agiledragon/gomonkey/v2"

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
	role := models.Role{
		AccessLevel: int(constants.SuperAdminRole),
		Name:        SuperAdminRoleName,
	}
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
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.CreateRole,
				func(role models.Role, ctx context.Context) (models.Role, error) {
					return role, nil
				}).
				ApplyFunc(rediscache.GetUser,
					func(userID int, ctx context.Context) (*models.User, error) {
						return testutls.MockUser(), nil
					}).
				ApplyFunc(rediscache.GetRole,
					func(userID int, ctx context.Context) (*models.Role, error) {
						return &role, nil
					})
		},
	}
}

// Suggested refactoring to use table-driven tests and helper functions for common setup.

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
		errorFromCreateRoleCase(),
		successCase(),
	}
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
		t.Run(tt.name,
			func(t *testing.T) {
				// Apply additional monkey patches based on test case name.
				patch := tt.init()
				// Create a new context
				c := context.Background()
				// Call the resolver function
				response, err := resolver1.Mutation().CreateRole(c, tt.req)
				// Check if the error matches the expected error
				if tt.wantErr {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, tt.wantResp, response)
				}
				if patch != nil {
					patch.Reset()
				}
			})
	}
}
