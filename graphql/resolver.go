package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	graphql1 "github.com/wednesday-solution/go-boiler/graphql/models"
)

type Resolver struct{}

func (r *mutationResolver) CreateComment(ctx context.Context, input graphql1.CommentCreateInput) (*graphql1.CommentPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateComments(ctx context.Context, input graphql1.CommentsCreateInput) (*graphql1.CommentsPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, input graphql1.CommentUpdateInput) (*graphql1.CommentPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateComments(ctx context.Context, filter *graphql1.CommentFilter, input graphql1.CommentsUpdateInput) (*graphql1.CommentsUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*graphql1.CommentDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteComments(ctx context.Context, filter *graphql1.CommentFilter) (*graphql1.CommentsDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateCompany(ctx context.Context, input graphql1.CompanyCreateInput) (*graphql1.CompanyPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateCompanies(ctx context.Context, input graphql1.CompaniesCreateInput) (*graphql1.CompaniesPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateCompany(ctx context.Context, id string, input graphql1.CompanyUpdateInput) (*graphql1.CompanyPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateCompanies(ctx context.Context, filter *graphql1.CompanyFilter, input graphql1.CompaniesUpdateInput) (*graphql1.CompaniesUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteCompany(ctx context.Context, id string) (*graphql1.CompanyDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteCompanies(ctx context.Context, filter *graphql1.CompanyFilter) (*graphql1.CompaniesDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateFollower(ctx context.Context, input graphql1.FollowerCreateInput) (*graphql1.FollowerPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateFollowers(ctx context.Context, input graphql1.FollowersCreateInput) (*graphql1.FollowersPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateFollower(ctx context.Context, id string, input graphql1.FollowerUpdateInput) (*graphql1.FollowerPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateFollowers(ctx context.Context, filter *graphql1.FollowerFilter, input graphql1.FollowersUpdateInput) (*graphql1.FollowersUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteFollower(ctx context.Context, id string) (*graphql1.FollowerDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteFollowers(ctx context.Context, filter *graphql1.FollowerFilter) (*graphql1.FollowersDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateLocation(ctx context.Context, input graphql1.LocationCreateInput) (*graphql1.LocationPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateLocations(ctx context.Context, input graphql1.LocationsCreateInput) (*graphql1.LocationsPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateLocation(ctx context.Context, id string, input graphql1.LocationUpdateInput) (*graphql1.LocationPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateLocations(ctx context.Context, filter *graphql1.LocationFilter, input graphql1.LocationsUpdateInput) (*graphql1.LocationsUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteLocation(ctx context.Context, id string) (*graphql1.LocationDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteLocations(ctx context.Context, filter *graphql1.LocationFilter) (*graphql1.LocationsDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreatePost(ctx context.Context, input graphql1.PostCreateInput) (*graphql1.PostPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreatePosts(ctx context.Context, input graphql1.PostsCreateInput) (*graphql1.PostsPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input graphql1.PostUpdateInput) (*graphql1.PostPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePosts(ctx context.Context, filter *graphql1.PostFilter, input graphql1.PostsUpdateInput) (*graphql1.PostsUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*graphql1.PostDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeletePosts(ctx context.Context, filter *graphql1.PostFilter) (*graphql1.PostsDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateRole(ctx context.Context, input graphql1.RoleCreateInput) (*graphql1.RolePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateRoles(ctx context.Context, input graphql1.RolesCreateInput) (*graphql1.RolesPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateRole(ctx context.Context, id string, input graphql1.RoleUpdateInput) (*graphql1.RolePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateRoles(ctx context.Context, filter *graphql1.RoleFilter, input graphql1.RolesUpdateInput) (*graphql1.RolesUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteRole(ctx context.Context, id string) (*graphql1.RoleDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteRoles(ctx context.Context, filter *graphql1.RoleFilter) (*graphql1.RolesDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateUser(ctx context.Context, input graphql1.UserCreateInput) (*graphql1.UserPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateUsers(ctx context.Context, input graphql1.UsersCreateInput) (*graphql1.UsersPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input graphql1.UserUpdateInput) (*graphql1.UserPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateUsers(ctx context.Context, filter *graphql1.UserFilter, input graphql1.UsersUpdateInput) (*graphql1.UsersUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*graphql1.UserDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, filter *graphql1.UserFilter) (*graphql1.UsersDeletePayload, error) {
	panic("not implemented")
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*graphql1.Comment, error) {
	panic("not implemented")
}

func (r *queryResolver) Comments(ctx context.Context, filter *graphql1.CommentFilter) ([]*graphql1.Comment, error) {
	panic("not implemented")
}

func (r *queryResolver) Company(ctx context.Context, id string) (*graphql1.Company, error) {
	panic("not implemented")
}

func (r *queryResolver) Companies(ctx context.Context, filter *graphql1.CompanyFilter) ([]*graphql1.Company, error) {
	panic("not implemented")
}

func (r *queryResolver) Follower(ctx context.Context, id string) (*graphql1.Follower, error) {
	panic("not implemented")
}

func (r *queryResolver) Followers(ctx context.Context, filter *graphql1.FollowerFilter) ([]*graphql1.Follower, error) {
	panic("not implemented")
}

func (r *queryResolver) Location(ctx context.Context, id string) (*graphql1.Location, error) {
	panic("not implemented")
}

func (r *queryResolver) Locations(ctx context.Context, filter *graphql1.LocationFilter) ([]*graphql1.Location, error) {
	panic("not implemented")
}

func (r *queryResolver) Post(ctx context.Context, id string) (*graphql1.Post, error) {
	panic("not implemented")
}

func (r *queryResolver) Posts(ctx context.Context, filter *graphql1.PostFilter) ([]*graphql1.Post, error) {
	panic("not implemented")
}

func (r *queryResolver) Role(ctx context.Context, id string) (*graphql1.Role, error) {
	panic("not implemented")
}

func (r *queryResolver) Roles(ctx context.Context, filter *graphql1.RoleFilter) ([]*graphql1.Role, error) {
	panic("not implemented")
}

func (r *queryResolver) User(ctx context.Context, id string) (*graphql1.User, error) {
	panic("not implemented")
}

func (r *queryResolver) Users(ctx context.Context, filter *graphql1.UserFilter) ([]*graphql1.User, error) {
	panic("not implemented")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
