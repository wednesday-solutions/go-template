// Package goboiler ...
// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
package goboiler

import (
	"context"
	fm "github.com/wednesday-solutions/go-boiler/graphql_models"
)

// Resolver ...
type Resolver struct {
}

func (r *mutationResolver) CreateComment(ctx context.Context, input fm.CommentCreateInput) (*fm.CommentPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateComments(ctx context.Context, input fm.CommentsCreateInput) (*fm.CommentsPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, input fm.CommentUpdateInput) (*fm.CommentPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateComments(ctx context.Context, filter *fm.CommentFilter, input fm.CommentUpdateInput) (*fm.CommentsUpdatePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*fm.CommentDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteComments(ctx context.Context, filter *fm.CommentFilter) (*fm.CommentsDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateCompany(ctx context.Context, input fm.CompanyCreateInput) (*fm.CompanyPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateCompanies(ctx context.Context, input fm.CompaniesCreateInput) (*fm.CompaniesPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateCompany(ctx context.Context, id string, input fm.CompanyUpdateInput) (*fm.CompanyPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateCompanies(ctx context.Context, filter *fm.CompanyFilter, input fm.CompanyUpdateInput) (*fm.CompaniesUpdatePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteCompany(ctx context.Context, id string) (*fm.CompanyDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteCompanies(ctx context.Context, filter *fm.CompanyFilter) (*fm.CompaniesDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateFollower(ctx context.Context, input fm.FollowerCreateInput) (*fm.FollowerPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateFollowers(ctx context.Context, input fm.FollowersCreateInput) (*fm.FollowersPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateFollower(ctx context.Context, id string, input fm.FollowerUpdateInput) (*fm.FollowerPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateFollowers(ctx context.Context, filter *fm.FollowerFilter, input fm.FollowerUpdateInput) (*fm.FollowersUpdatePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteFollower(ctx context.Context, id string) (*fm.FollowerDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteFollowers(ctx context.Context, filter *fm.FollowerFilter) (*fm.FollowersDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateLocation(ctx context.Context, input fm.LocationCreateInput) (*fm.LocationPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateLocations(ctx context.Context, input fm.LocationsCreateInput) (*fm.LocationsPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateLocation(ctx context.Context, id string, input fm.LocationUpdateInput) (*fm.LocationPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateLocations(ctx context.Context, filter *fm.LocationFilter, input fm.LocationUpdateInput) (*fm.LocationsUpdatePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteLocation(ctx context.Context, id string) (*fm.LocationDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteLocations(ctx context.Context, filter *fm.LocationFilter) (*fm.LocationsDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreatePost(ctx context.Context, input fm.PostCreateInput) (*fm.PostPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreatePosts(ctx context.Context, input fm.PostsCreateInput) (*fm.PostsPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input fm.PostUpdateInput) (*fm.PostPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdatePosts(ctx context.Context, filter *fm.PostFilter, input fm.PostUpdateInput) (*fm.PostsUpdatePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*fm.PostDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeletePosts(ctx context.Context, filter *fm.PostFilter) (*fm.PostsDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateRole(ctx context.Context, input fm.RoleCreateInput) (*fm.RolePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateRoles(ctx context.Context, input fm.RolesCreateInput) (*fm.RolesPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateRole(ctx context.Context, id string, input fm.RoleUpdateInput) (*fm.RolePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateRoles(ctx context.Context, filter *fm.RoleFilter, input fm.RoleUpdateInput) (*fm.RolesUpdatePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteRole(ctx context.Context, id string) (*fm.RoleDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteRoles(ctx context.Context, filter *fm.RoleFilter) (*fm.RolesDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateUser(ctx context.Context, input fm.UserCreateInput) (*fm.UserPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) CreateUsers(ctx context.Context, input fm.UsersCreateInput) (*fm.UsersPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input fm.UserUpdateInput) (*fm.UserPayload, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateUsers(ctx context.Context, filter *fm.UserFilter, input fm.UserUpdateInput) (*fm.UsersUpdatePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*fm.UserDeletePayload, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, filter *fm.UserFilter) (*fm.UsersDeletePayload, error) {
	panic("implement me")
}

func (q queryResolver) Comment(ctx context.Context, id string) (*fm.Comment, error) {
	panic("implement me")
}

func (q queryResolver) Comments(ctx context.Context, filter *fm.CommentFilter, pagination *fm.CommentPagination) ([]*fm.Comment, error) {
	panic("implement me")
}

func (q queryResolver) Company(ctx context.Context, id string) (*fm.Company, error) {
	panic("implement me")
}

func (q queryResolver) Companies(ctx context.Context, filter *fm.CompanyFilter, pagination *fm.CompanyPagination) ([]*fm.Company, error) {
	panic("implement me")
}

func (q queryResolver) Follower(ctx context.Context, id string) (*fm.Follower, error) {
	panic("implement me")
}

func (q queryResolver) Followers(ctx context.Context, filter *fm.FollowerFilter, pagination *fm.FollowerPagination) ([]*fm.Follower, error) {
	panic("implement me")
}

func (q queryResolver) Location(ctx context.Context, id string) (*fm.Location, error) {
	panic("implement me")
}

func (q queryResolver) Locations(ctx context.Context, filter *fm.LocationFilter, pagination *fm.LocationPagination) ([]*fm.Location, error) {
	panic("implement me")
}

func (q queryResolver) Post(ctx context.Context, id string) (*fm.Post, error) {
	panic("implement me")
}

func (q queryResolver) Posts(ctx context.Context, filter *fm.PostFilter, pagination *fm.PostPagination) ([]*fm.Post, error) {
	panic("implement me")
}

func (q queryResolver) Role(ctx context.Context, id string) (*fm.Role, error) {
	panic("implement me")
}

func (q queryResolver) Roles(ctx context.Context, filter *fm.RoleFilter, pagination *fm.RolePagination) ([]*fm.Role, error) {
	panic("implement me")
}

func (q queryResolver) User(ctx context.Context, id string) (*fm.User, error) {
	panic("implement me")
}

func (q queryResolver) Users(ctx context.Context, filter *fm.UserFilter, pagination *fm.UserPagination) ([]*fm.User, error) {
	panic("implement me")
}

// Mutation ...
func (r *Resolver) Mutation() fm.MutationResolver { return &mutationResolver{r} }

// Query ...
func (r *Resolver) Query() fm.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
