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
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/rate_throttle"
	resultwrapper "go-template/pkg/utl/result_wrapper"
	"strconv"
	"time"

	null "github.com/volatiletech/null/v8"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input gqlmodels.UserCreateInput) (*gqlmodels.User, error) {
	err := throttle.Check(ctx, 5, 10*time.Second)
	if err != nil {
		return nil, err
	}

	roleId, _ := strconv.Atoi(input.RoleID)
	user := models.User{
		Username:  null.StringFrom(input.Username),
		Password:  null.StringFrom(input.Password),
		Email:     null.StringFrom(input.Email),
		FirstName: null.StringFrom(input.FirstName),
		LastName:  null.StringFrom(input.LastName),
		RoleID:    null.IntFrom(roleId),
	}
	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading config ")
	}
	// creating new secure service
	sec := service.Secure(cfg)
	user.Password = null.StringFrom(sec.Hash(user.Password.String))
	newUser, err := daos.CreateUserTx(user, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "user information")
	}
	graphUser := convert.UserToGraphQlUser(&newUser, 1)

	r.Lock()
	for _, observer := range r.Observers {
		observer <- graphUser
	}
	r.Unlock()

	return graphUser, err
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input *gqlmodels.UserUpdateInput) (*gqlmodels.User, error) {
	userID := auth.UserIDFromContext(ctx)
	user, _ := daos.FindUserByID(userID)
	var u models.User
	if user != nil {
		u = *user
	} else {
		return nil, resultwrapper.ResolverWrapperFromMessage(404, "user not found")
	}
	if input.FirstName != nil {
		u.FirstName = null.StringFromPtr(input.FirstName)
	}
	if input.LastName != nil {
		u.LastName = null.StringFromPtr(input.LastName)
	}
	if input.Mobile != nil {
		u.Mobile = null.StringFromPtr(input.Mobile)
	}
	if input.Address != nil {
		u.Address = null.StringFromPtr(input.Address)
	}
	_, err := daos.UpdateUserTx(u, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "new information")
	}

	graphUser := convert.UserToGraphQlUser(&u, 1)
	r.Lock()
	for _, observer := range r.Observers {
		observer <- graphUser
	}
	r.Unlock()

	return graphUser, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (*gqlmodels.UserDeletePayload, error) {
	userID := auth.UserIDFromContext(ctx)
	u, err := daos.FindUserByID(userID)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	_, err = daos.DeleteUser(*u)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "user")
	}
	return &gqlmodels.UserDeletePayload{ID: fmt.Sprint(userID)}, nil
}
