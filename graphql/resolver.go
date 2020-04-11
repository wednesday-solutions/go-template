package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	gqlmodels "github.com/wednesday-solution/go-boiler/graphql/models"
)

type Resolver struct{}

func (r *mutationResolver) CreateComment(ctx context.Context, input gqlmodels.CommentCreateInput) (*gqlmodels.CommentPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateComments(ctx context.Context, input gqlmodels.CommentsCreateInput) (*gqlmodels.CommentsPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, input gqlmodels.CommentUpdateInput) (*gqlmodels.CommentPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateComments(ctx context.Context, filter *gqlmodels.CommentFilter, input gqlmodels.CommentsUpdateInput) (*gqlmodels.CommentsUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*gqlmodels.CommentDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteComments(ctx context.Context, filter *gqlmodels.CommentFilter) (*gqlmodels.CommentsDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateCompany(ctx context.Context, input gqlmodels.CompanyCreateInput) (*gqlmodels.CompanyPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateCompanies(ctx context.Context, input gqlmodels.CompaniesCreateInput) (*gqlmodels.CompaniesPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateCompany(ctx context.Context, id string, input gqlmodels.CompanyUpdateInput) (*gqlmodels.CompanyPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateCompanies(ctx context.Context, filter *gqlmodels.CompanyFilter, input gqlmodels.CompaniesUpdateInput) (*gqlmodels.CompaniesUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteCompany(ctx context.Context, id string) (*gqlmodels.CompanyDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteCompanies(ctx context.Context, filter *gqlmodels.CompanyFilter) (*gqlmodels.CompaniesDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateFollower(ctx context.Context, input gqlmodels.FollowerCreateInput) (*gqlmodels.FollowerPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateFollowers(ctx context.Context, input gqlmodels.FollowersCreateInput) (*gqlmodels.FollowersPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateFollower(ctx context.Context, id string, input gqlmodels.FollowerUpdateInput) (*gqlmodels.FollowerPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateFollowers(ctx context.Context, filter *gqlmodels.FollowerFilter, input gqlmodels.FollowersUpdateInput) (*gqlmodels.FollowersUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteFollower(ctx context.Context, id string) (*gqlmodels.FollowerDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteFollowers(ctx context.Context, filter *gqlmodels.FollowerFilter) (*gqlmodels.FollowersDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateLocation(ctx context.Context, input gqlmodels.LocationCreateInput) (*gqlmodels.LocationPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateLocations(ctx context.Context, input gqlmodels.LocationsCreateInput) (*gqlmodels.LocationsPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateLocation(ctx context.Context, id string, input gqlmodels.LocationUpdateInput) (*gqlmodels.LocationPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateLocations(ctx context.Context, filter *gqlmodels.LocationFilter, input gqlmodels.LocationsUpdateInput) (*gqlmodels.LocationsUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteLocation(ctx context.Context, id string) (*gqlmodels.LocationDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteLocations(ctx context.Context, filter *gqlmodels.LocationFilter) (*gqlmodels.LocationsDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreatePost(ctx context.Context, input gqlmodels.PostCreateInput) (*gqlmodels.PostPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreatePosts(ctx context.Context, input gqlmodels.PostsCreateInput) (*gqlmodels.PostsPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input gqlmodels.PostUpdateInput) (*gqlmodels.PostPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePosts(ctx context.Context, filter *gqlmodels.PostFilter, input gqlmodels.PostsUpdateInput) (*gqlmodels.PostsUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*gqlmodels.PostDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeletePosts(ctx context.Context, filter *gqlmodels.PostFilter) (*gqlmodels.PostsDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateRole(ctx context.Context, input gqlmodels.RoleCreateInput) (*gqlmodels.RolePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateRoles(ctx context.Context, input gqlmodels.RolesCreateInput) (*gqlmodels.RolesPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateRole(ctx context.Context, id string, input gqlmodels.RoleUpdateInput) (*gqlmodels.RolePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateRoles(ctx context.Context, filter *gqlmodels.RoleFilter, input gqlmodels.RolesUpdateInput) (*gqlmodels.RolesUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteRole(ctx context.Context, id string) (*gqlmodels.RoleDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteRoles(ctx context.Context, filter *gqlmodels.RoleFilter) (*gqlmodels.RolesDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateUser(ctx context.Context, input gqlmodels.UserCreateInput) (*gqlmodels.UserPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateUsers(ctx context.Context, input gqlmodels.UsersCreateInput) (*gqlmodels.UsersPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input gqlmodels.UserUpdateInput) (*gqlmodels.UserPayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateUsers(ctx context.Context, filter *gqlmodels.UserFilter, input gqlmodels.UsersUpdateInput) (*gqlmodels.UsersUpdatePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*gqlmodels.UserDeletePayload, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, filter *gqlmodels.UserFilter) (*gqlmodels.UsersDeletePayload, error) {
	panic("not implemented")
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*gqlmodels.Comment, error) {
	panic("not implemented")
}

func (r *queryResolver) Comments(ctx context.Context, filter *gqlmodels.CommentFilter) ([]*gqlmodels.Comment, error) {
	panic("not implemented")
}

func (r *queryResolver) Company(ctx context.Context, id string) (*gqlmodels.Company, error) {
	panic("not implemented")
}

func (r *queryResolver) Companies(ctx context.Context, filter *gqlmodels.CompanyFilter) ([]*gqlmodels.Company, error) {
	panic("not implemented")
}

func (r *queryResolver) Follower(ctx context.Context, id string) (*gqlmodels.Follower, error) {
	panic("not implemented")
}

func (r *queryResolver) Followers(ctx context.Context, filter *gqlmodels.FollowerFilter) ([]*gqlmodels.Follower, error) {
	panic("not implemented")
}

func (r *queryResolver) Location(ctx context.Context, id string) (*gqlmodels.Location, error) {
	panic("not implemented")
}

func (r *queryResolver) Locations(ctx context.Context, filter *gqlmodels.LocationFilter) ([]*gqlmodels.Location, error) {
	panic("not implemented")
}

func (r *queryResolver) Post(ctx context.Context, id string) (*gqlmodels.Post, error) {
	panic("not implemented")
}

func (r *queryResolver) Posts(ctx context.Context, filter *gqlmodels.PostFilter) ([]*gqlmodels.Post, error) {
	panic("not implemented")
}

func (r *queryResolver) Role(ctx context.Context, id string) (*gqlmodels.Role, error) {
	panic("not implemented")
}

func (r *queryResolver) Roles(ctx context.Context, filter *gqlmodels.RoleFilter) ([]*gqlmodels.Role, error) {
	panic("not implemented")
}

func (r *queryResolver) User(ctx context.Context, id string) (*gqlmodels.User, error) {
	panic("not implemented")
}

func (r *queryResolver) Users(ctx context.Context, filter *gqlmodels.UserFilter) ([]*gqlmodels.User, error) {
	panic("not implemented")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
