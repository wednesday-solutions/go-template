// Package goboiler ...
// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
package goboiler

import (
	"context"
	"fmt"
	"github.com/volatiletech/null"
	"github.com/wednesday-solutions/go-boiler/daos"
	fm "github.com/wednesday-solutions/go-boiler/graphql_models"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/config"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/convert"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/middleware/auth"
	resultwrapper "github.com/wednesday-solutions/go-boiler/pkg/utl/result_wrapper"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/service"
)

// Resolver ...
type Resolver struct {
}

func (q queryResolver) Login(ctx context.Context, username string, password string) (*fm.LoginResponse, error) {
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
		return nil, ErrUnauthorized
	}

	token, err := tg.GenerateToken(u)
	if err != nil {
		return nil, ErrUnauthorized
	}

	refreshToken := sec.Token(token)
	u.Token = null.StringFrom(refreshToken)
	_, err = daos.UpdateUserTx(*u, nil)
	if err != nil {
		return nil, err
	}

	return &fm.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (q queryResolver) Me(ctx context.Context) (*fm.User, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := daos.FindUserByID(userID)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
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

func (q mutationResolver) RefreshToken(ctx context.Context, token string) (*fm.RefreshTokenResponse, error) {
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

// Mutation ...
func (r *Resolver) Mutation() fm.MutationResolver { return &mutationResolver{r} }

// Query ...
func (r *Resolver) Query() fm.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
