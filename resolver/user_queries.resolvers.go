package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go-template/daos"
	"go-template/gqlmodels"
	"go-template/internal/middleware/auth"
	"go-template/pkg/utl/cnvrttogql"
	"go-template/pkg/utl/rediscache"
	"go-template/pkg/utl/resultwrapper"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*gqlmodels.User, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := rediscache.GetUser(userID, ctx)
	if err != nil {
		return &gqlmodels.User{}, resultwrapper.ResolverSQLError(err, "data")
	}

	return cnvrttogql.UserToGraphQlUser(user, 1), err
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, pagination *gqlmodels.UserPagination) (*gqlmodels.UsersPayload, error) {
	var queryMods []qm.QueryMod
	if pagination != nil {
		if pagination.Limit != 0 {
			queryMods = append(queryMods, qm.Limit(pagination.Limit), qm.Offset(pagination.Page*pagination.Limit))
		}
	}

	users, count, err := daos.FindAllUsersWithCount(queryMods, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return &gqlmodels.UsersPayload{Total: int(count), Users: cnvrttogql.UsersToGraphQlUsers(users, 1)}, nil
}

// Query returns gqlmodels.QueryResolver implementation.
func (r *Resolver) Query() gqlmodels.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
