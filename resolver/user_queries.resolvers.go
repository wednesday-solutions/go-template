package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go-template/daos"
	"go-template/gqlmodels"
	"go-template/internal/middleware/auth"
	"go-template/pkg/utl/convert"
	rediscache "go-template/pkg/utl/redis_cache"
	resultwrapper "go-template/pkg/utl/result_wrapper"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *queryResolver) Me(ctx context.Context) (*gqlmodels.User, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := rediscache.GetUser(userID)
	if err != nil {
		return &gqlmodels.User{}, resultwrapper.ResolverSQLError(err, "data")
	}

	return convert.UserToGraphQlUser(user, 1), err
}

func (r *queryResolver) Users(ctx context.Context, pagination *gqlmodels.UserPagination) (
	*gqlmodels.UsersPayload, error) {
	var queryMods []qm.QueryMod
	if pagination != nil {
		if pagination.Limit != 0 {
			queryMods = append(queryMods, qm.Limit(pagination.Limit), qm.Offset(pagination.Page*pagination.Limit))
		}
	}

	users, count, err := daos.FindAllUsersWithCount(queryMods)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return &gqlmodels.UsersPayload{Total: int(count), Users: convert.UsersToGraphQlUsers(users, 1)}, nil
}

// Query returns gqlmodels.QueryResolver implementation.
func (r *Resolver) Query() gqlmodels.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
