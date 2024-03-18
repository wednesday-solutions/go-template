package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"go-template/daos"
	"go-template/gqlmodels"
	"go-template/models"
	"go-template/pkg/utl/cnvrttogql"
	"go-template/pkg/utl/resultwrapper"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Author is the resolver for the author field.
func (r *queryResolver) Author(ctx context.Context, input gqlmodels.AuthorQueryInput) (*gqlmodels.Author, error) {
	var author *models.Author
	var err error
	author, err = daos.FindAuthorWithId(input.ID, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return cnvrttogql.AuthorToGraphQlAuthor(author), nil
}

// Authors is the resolver for the authors field.
func (r *queryResolver) Authors(ctx context.Context, pagination *gqlmodels.AuthorPagination) (*gqlmodels.AuthorsPayload, error) {
	var queryMods []qm.QueryMod
	if pagination != nil {
		if pagination.Limit != 0 {
			queryMods = append(queryMods, qm.Limit(pagination.Limit), qm.Offset(pagination.Page*pagination.Limit))
		}
	}
	authors, count, err := daos.FindAllAuthorsWithCount(queryMods, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return &gqlmodels.AuthorsPayload{Total: int(count), Authors: cnvrttogql.AuthorsToGraphQlAuthors(authors)}, nil
}

// Query returns gqlmodels.QueryResolver implementation.
func (r *Resolver) Query() gqlmodels.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }