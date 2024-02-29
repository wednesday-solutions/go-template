package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"fmt"
	"go-template/daos"
	"go-template/gqlmodels"
	"go-template/models"
	"go-template/pkg/utl/cnvrttogql"
	"go-template/pkg/utl/resultwrapper"
	"go-template/pkg/utl/throttle"
	"strconv"
	"time"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input gqlmodels.PostCreateInput) (*gqlmodels.Post, error) {
	err := throttle.Check(ctx, 5, 10*time.Second)
	if err != nil {
		return nil, err
	}
	authorId, err := strconv.Atoi(input.AuthorID)
	if err != nil {
		return nil, err
	}
	post := models.Post{
		AuthorID: authorId,
		Post:     input.Post,
	}
	newPost, err := daos.CreatePost(post, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "Post information")
	}
	graphPost := cnvrttogql.PostToGraphQlPost(&newPost, 1)

	return graphPost, err
}

// UpdatePost is the resolver for the updatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, input gqlmodels.PostUpdateInput) (*gqlmodels.Post, error) {
	var u *models.Post
	var err error
	u, err = daos.FindPostWithId(input.ID, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	if input.Post != nil {
		u.Post = *input.Post
	}
	// update the post in the database
	_, err = daos.UpdatePost(*u, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "new information")
	}
	// convert the updated post to a graphql post
	graphPost := cnvrttogql.PostToGraphQlPost(&*u, 1)
	return graphPost, nil
}

// DeletePost is the resolver for the deletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, input gqlmodels.PostDeleteInput) (*gqlmodels.PostDeletePayload, error) {
	u, err := daos.FindPostWithId(input.ID, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	_, err = daos.DeletePost(*u, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "post")
	}
	return &gqlmodels.PostDeletePayload{ID: fmt.Sprint(input.ID)}, nil
}
