package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/internal/middleware/auth"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	rediscache "github.com/wednesday-solutions/go-template/pkg/utl/redis_cache"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

func (r *queryResolver) Me(ctx context.Context) (*graphql_models.User, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := rediscache.GetUser(userID)
	if err != nil {
		return &graphql_models.User{}, resultwrapper.ResolverSQLError(err, "data")
	}
	return convert.UserToGraphQlUser(user), err
}

func (r *queryResolver) Users(ctx context.Context, pagination *graphql_models.UserPagination) ([]*graphql_models.User, error) {
	var queryMods []qm.QueryMod
	fmt.Println("HIIII")
	if pagination != nil {
		if pagination.Limit != 0 {
			queryMods = append(queryMods, qm.Limit(pagination.Limit), qm.Offset(pagination.Page*pagination.Limit))
		}
	}
	users, _, err := daos.FindAllUsersWithCount(queryMods)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return convert.UsersToGraphQlUsers(users), nil
}
