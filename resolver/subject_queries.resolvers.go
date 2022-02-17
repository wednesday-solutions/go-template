package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

func (r *queryResolver) Subject(ctx context.Context, id int) (*graphql_models.Subject, error) {
	subject, err := daos.FindSubjectById(id)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return convert.SubjectToGraphQlSubject(subject), nil
}

func (r *queryResolver) Subjects(ctx context.Context, pagination graphql_models.Pagination) ([]*graphql_models.Subject, error) {
	var queryMods []qm.QueryMod
	if pagination.Limit != 0 {
		queryMods = append(queryMods, qm.Limit(pagination.Limit), qm.Offset((pagination.Page-1)*pagination.Limit))
	}

	subjects, _, err := daos.FindAllSubjectsWithCount(queryMods)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return convert.SubjectsToGraphQlSubjects(subjects), nil
}

// Query returns graphql_models.QueryResolver implementation.
func (r *Resolver) Query() graphql_models.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
