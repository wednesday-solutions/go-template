// Package gotemplate ...
// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
package gotemplate

import (
	"context"
	"fmt"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-template/daos"
	fm "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/pkg/utl/config"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	"github.com/wednesday-solutions/go-template/pkg/utl/middleware/auth"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
	"github.com/wednesday-solutions/go-template/pkg/utl/service"
)

// Resolver ...
type Resolver struct {
}

func (r queryResolver) Me(ctx context.Context) (*fm.User, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := GetUser(userID)
	if err != nil {
		return &fm.User{}, resultwrapper.ResolverSQLError(err, "data")
	}
	return &fm.User{
		FirstName: convert.NullDotStringToPointerString(user.FirstName),
		LastName:  convert.NullDotStringToPointerString(user.LastName),
		Username:  convert.NullDotStringToPointerString(user.Username),
		Email:     convert.NullDotStringToPointerString(user.Email),
		Mobile:    convert.NullDotStringToPointerString(user.Mobile),
		Phone:     convert.NullDotStringToPointerString(user.Phone),
		Address:   convert.NullDotStringToPointerString(user.Address),
	}, err
}

func (r queryResolver) Users(ctx context.Context, pagination *fm.UserPagination) (*fm.UsersPayload, error) {

	var queryMods []qm.QueryMod
	if pagination != nil {
		if pagination.Limit != 0 {
			queryMods = append(queryMods, qm.Limit(pagination.Limit), qm.Offset(pagination.Page*pagination.Limit))
		}
	}
	users, err := models.Users(queryMods...).All(context.Background(), boil.GetContextDB())
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return &fm.UsersPayload{Users: convert.UsersToGraphQlUsers(users)}, nil
}

func (r mutationResolver) CreateRole(ctx context.Context, input fm.RoleCreateInput) (*fm.RolePayload, error) {
	role := models.Role{
		AccessLevel: input.AccessLevel,
		Name:        input.Name,
	}
	newRole, err := daos.CreateRoleTx(role, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "role")
	}
	return &fm.RolePayload{Role: &fm.Role{
		AccessLevel: newRole.AccessLevel,
		Name:        newRole.Name,
	},
	}, err
}

func (r mutationResolver) Login(ctx context.Context, username string, password string) (*fm.LoginResponse, error) {
	u, err := daos.FindUserByUserName(username)
	if err != nil {
		return nil, err
	}
	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Error in loading config ")
	}
	// creating new secure and token generation service
	sec := service.Secure(cfg)
	tg, err := service.JWT(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error in creating auth service ")
	}

	if !u.Password.Valid || (!sec.HashMatchesPassword(u.Password.String, password)) {
		return nil, fmt.Errorf("Username or password does not exist ")
	}

	if !u.Active.Valid || (!u.Active.Bool) {
		return nil, resultwrapper.ErrUnauthorized
	}

	token, err := tg.GenerateToken(u)
	if err != nil {
		return nil, resultwrapper.ErrUnauthorized
	}

	refreshToken := sec.Token(token)
	u.Token = null.StringFrom(refreshToken)
	_, err = daos.UpdateUserTx(*u, nil)
	if err != nil {
		return nil, err
	}

	return &fm.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (r mutationResolver) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (*fm.ChangePasswordResponse, error) {

	userID := auth.UserIDFromContext(ctx)
	u, err := daos.FindUserByID(userID)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}

	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Error in loading config ")
	}
	// creating new secure service
	sec := service.Secure(cfg)
	if !sec.HashMatchesPassword(convert.NullDotStringToString(u.Password), oldPassword) {
		return nil, fmt.Errorf("incorrect old password")
	}

	if !sec.Password(newPassword, convert.NullDotStringToString(u.FirstName), convert.NullDotStringToString(u.LastName), convert.NullDotStringToString(u.Username), convert.NullDotStringToString(u.Email)) {
		return nil, fmt.Errorf("insecure password")
	}

	u.Password = null.StringFrom(sec.Hash(newPassword))
	_, err = daos.UpdateUserTx(*u, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "new information")
	}
	return &fm.ChangePasswordResponse{Ok: true}, err
}

func (r mutationResolver) RefreshToken(ctx context.Context, token string) (*fm.RefreshTokenResponse, error) {
	user, err := daos.FindUserByToken(token)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "token")
	}
	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Error in loading config ")
	}
	// creating new secure and token generation service
	tg, err := service.JWT(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error in creating auth service ")
	}
	resp, err := tg.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	return &fm.RefreshTokenResponse{Token: resp}, nil
}

func (r mutationResolver) CreateUser(ctx context.Context, input fm.UserCreateInput) (*fm.UserPayload, error) {
	user := models.User{
		Username:  null.StringFromPtr(input.Username),
		Password:  null.StringFromPtr(input.Password),
		Email:     null.StringFromPtr(input.Email),
		FirstName: null.StringFromPtr(input.FirstName),
		LastName:  null.StringFromPtr(input.LastName),
		RoleID:    convert.PointerStringToNullDotInt(input.RoleID),
	}
	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Error in loading config ")
	}
	// creating new secure service
	sec := service.Secure(cfg)
	user.Password = null.StringFrom(sec.Hash(user.Password.String))
	newUser, err := daos.CreateUserTx(user, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "user information")
	}
	return &fm.UserPayload{User: &fm.User{
		FirstName: convert.NullDotStringToPointerString(newUser.FirstName),
		LastName:  convert.NullDotStringToPointerString(newUser.LastName),
		Username:  convert.NullDotStringToPointerString(newUser.Username),
		Email:     convert.NullDotStringToPointerString(newUser.Email),
		Mobile:    convert.NullDotStringToPointerString(newUser.Mobile),
		Phone:     convert.NullDotStringToPointerString(newUser.Phone),
		Address:   convert.NullDotStringToPointerString(newUser.Address),
	},
	}, err
}

func (r mutationResolver) UpdateUser(ctx context.Context, input *fm.UserUpdateInput) (*fm.UserUpdatePayload, error) {
	userID := auth.UserIDFromContext(ctx)
	u := models.User{
		ID:        userID,
		FirstName: null.StringFromPtr(input.FirstName),
		LastName:  null.StringFromPtr(input.LastName),
		Mobile:    null.StringFromPtr(input.Mobile),
		Phone:     null.StringFromPtr(input.Phone),
		Address:   null.StringFromPtr(input.Address),
	}
	_, err := daos.UpdateUserTx(u, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "new information")
	}
	return &fm.UserUpdatePayload{Ok: true}, nil
}

func (r mutationResolver) DeleteUser(ctx context.Context) (*fm.UserDeletePayload, error) {
	userID := auth.UserIDFromContext(ctx)
	u, err := daos.FindUserByID(userID)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	_, err = daos.DeleteUser(*u)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "user")
	}
	return &fm.UserDeletePayload{ID: fmt.Sprint(userID)}, nil
}

// Mutation ...
func (r *Resolver) Mutation() fm.MutationResolver { return &mutationResolver{r} }

// Query ...
func (r *Resolver) Query() fm.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
