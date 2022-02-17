package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	throttle "github.com/wednesday-solutions/go-template/pkg/utl/rate_throttle"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

func (r *mutationResolver) CreateSubject(ctx context.Context, createSubjectInput *graphql_models.CreateSubjectInput) (*graphql_models.Subject, error) {
	err := throttle.Check(ctx, 5, 10*time.Second)
	if err != nil {
		return nil, err
	}
	subject := models.Subject{
		Name: createSubjectInput.Name,
	}

	insertedSubject, err := daos.CreateSubject(subject)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "Could not Create a new Subject")
	}
	graphqlSubjectObject := convert.SubjectToGraphQlSubject(&insertedSubject)
	return graphqlSubjectObject, nil
}

func (r *mutationResolver) UpdateSubject(ctx context.Context, updateSubjectInput *graphql_models.UpdateSubjectInput) (*graphql_models.Subject, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteSubject(ctx context.Context, deleteSubjectInput *graphql_models.DeleteSubjectInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
