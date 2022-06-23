package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"go-template/daos"
	"go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/internal/middleware/auth"
	"go-template/internal/service"
	"go-template/pkg/utl/convert"
	resultwrapper "go-template/pkg/utl/result_wrapper"

	null "github.com/volatiletech/null/v8"
)

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (
	*gqlmodels.LoginResponse, error) {
	u, err := daos.FindUserByUserName(username)
	if err != nil {
		return nil, err
	}
	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading config ")
	}
	// creating new secure and token generation service
	sec := service.Secure(cfg)
	tg, err := service.JWT(cfg)
	if err != nil {
		return nil, fmt.Errorf("error in creating auth service ")
	}

	if !u.Password.Valid || (!sec.HashMatchesPassword(u.Password.String, password)) {
		return nil, fmt.Errorf("username or password does not exist ")
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

	return &gqlmodels.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (
	*gqlmodels.ChangePasswordResponse, error) {
	userID := auth.UserIDFromContext(ctx)
	u, err := daos.FindUserByID(userID)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}

	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading config ")
	}
	// creating new secure service
	sec := service.Secure(cfg)
	if !sec.HashMatchesPassword(convert.NullDotStringToString(u.Password), oldPassword) {
		return nil, fmt.Errorf("incorrect old password")
	}

	if !sec.Password(newPassword,
		convert.NullDotStringToString(u.FirstName),
		convert.NullDotStringToString(u.LastName),
		convert.NullDotStringToString(u.Username),
		convert.NullDotStringToString(u.Email)) {
		return nil, fmt.Errorf("insecure password")
	}

	u.Password = null.StringFrom(sec.Hash(newPassword))
	_, err = daos.UpdateUserTx(*u, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "new information")
	}
	return &gqlmodels.ChangePasswordResponse{Ok: true}, err
}

func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*gqlmodels.RefreshTokenResponse, error) {
	user, err := daos.FindUserByToken(token)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "token")
	}
	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading config ")
	}
	// creating new secure and token generation service
	tg, err := service.JWT(cfg)
	if err != nil {
		return nil, fmt.Errorf("error in creating auth service ")
	}
	resp, err := tg.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	return &gqlmodels.RefreshTokenResponse{Token: resp}, nil
}

// Mutation returns gqlmodels.MutationResolver implementation.
func (r *Resolver) Mutation() gqlmodels.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
