package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"go-template/daos"
	"go-template/gqlmodels"
	"go-template/internal/constants"
	"go-template/internal/middleware/auth"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/rediscache"
	"go-template/pkg/utl/resultwrapper"
)

// CreateRole is the resolver for the createRole field.
func (r *mutationResolver) CreateRole(ctx context.Context, input gqlmodels.RoleCreateInput) (*gqlmodels.RolePayload, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := rediscache.GetUser(userID, ctx)
	if err != nil {
		return &gqlmodels.RolePayload{}, resultwrapper.ResolverSQLError(err, "data")
	}
	userRole, err := rediscache.GetRole(convert.NullDotIntToInt(user.RoleID), ctx)
	if err != nil {
		return &gqlmodels.RolePayload{}, resultwrapper.ResolverSQLError(err, "data")
	}
	role := models.Role{
		AccessLevel: input.AccessLevel,
		Name:        input.Name,
	}
	if userRole.AccessLevel != int(constants.SuperAdminRole) {
		return &gqlmodels.RolePayload{}, fmt.Errorf("You don't appear to have enough access level for this request ")
	}

	newRole, err := daos.CreateRole(role, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "role")
	}
	return &gqlmodels.RolePayload{Role: &gqlmodels.Role{
		AccessLevel: newRole.AccessLevel,
		Name:        newRole.Name,
	},
	}, err
}
