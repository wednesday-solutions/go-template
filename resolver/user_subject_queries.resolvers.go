package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

func (r *queryResolver) UserSubject(ctx context.Context, userID *int, subjectID *int) (*graphql_models.UserSubject, error) {
	userSubject, err := daos.FindUserSubjectById(userID, subjectID)
	fmt.Print("userSubject", userSubject)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return convert.UserSubjectToGraphQlUserSubject(userSubject), nil
}
