package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	null "github.com/volatiletech/null/v8"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/internal/config"
	"github.com/wednesday-solutions/go-template/internal/middleware/auth"
	"github.com/wednesday-solutions/go-template/internal/service"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	"github.com/wednesday-solutions/go-template/pkg/utl/rate_throttle"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

func (r *mutationResolver) CreateUser(ctx context.Context, createUserInput graphql_models.UserCreateInput) (*graphql_models.User, error) {
	err := throttle.Check(ctx, 5, 10*time.Second)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:  null.StringFromPtr(createUserInput.Username),
		Password:  null.StringFromPtr(createUserInput.Password),
		Email:     null.StringFromPtr(createUserInput.Email),
		FirstName: null.StringFromPtr(createUserInput.FirstName),
		LastName:  null.StringFromPtr(createUserInput.LastName),
		RoleID:    convert.PointerStringToNullDotInt(createUserInput.RoleID),
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
	graphUser := convert.UserToGraphQlUser(&newUser)

	r.Lock()
	for _, observer := range r.Observers {
		observer <- graphUser
	}
	r.Unlock()

	return graphUser, err
}

func (r *mutationResolver) UpdateUser(ctx context.Context, updateUserInput *graphql_models.UserUpdateInput) (*graphql_models.User, error) {
	userID := auth.UserIDFromContext(ctx)
	u := models.User{
		ID:        userID,
		FirstName: null.StringFromPtr(updateUserInput.FirstName),
		LastName:  null.StringFromPtr(updateUserInput.LastName),
		Mobile:    null.StringFromPtr(updateUserInput.Mobile),
		Phone:     null.StringFromPtr(updateUserInput.Phone),
		Address:   null.StringFromPtr(updateUserInput.Address),
	}
	_, err := daos.UpdateUserTx(u, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "new information")
	}

	graphUser := convert.UserToGraphQlUser(&u)
	r.Lock()
	for _, observer := range r.Observers {
		observer <- graphUser
	}
	r.Unlock()

	return graphUser, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (string, error) {
	userID := auth.UserIDFromContext(ctx)
	u, err := daos.FindUserByID(userID)
	if err != nil {
		return "", resultwrapper.ResolverSQLError(err, "data")
	}
	_, err = daos.DeleteUser(*u)
	if err != nil {
		return "", resultwrapper.ResolverSQLError(err, "user")
	}
	return fmt.Sprint(userID), nil
}
