package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	gotemplate "github.com/wednesday-solutions/go-template"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/internal/config"
	"github.com/wednesday-solutions/go-template/internal/middleware/auth"
	"github.com/wednesday-solutions/go-template/internal/service"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/pkg/utl"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	"github.com/wednesday-solutions/go-template/pkg/utl/rate_throttle"
	rediscache "github.com/wednesday-solutions/go-template/pkg/utl/redis_cache"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

func (r *mutationResolver) CreateRole(ctx context.Context, input graphql_models.RoleCreateInput) (*graphql_models.RolePayload, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := rediscache.GetUser(userID)
	if err != nil {
		return &graphql_models.RolePayload{}, resultwrapper.ResolverSQLError(err, "data")
	}
	userRole, err := rediscache.GetRole(convert.NullDotIntToInt(user.RoleID))
	if err != nil {
		return &graphql_models.RolePayload{}, resultwrapper.ResolverSQLError(err, "data")
	}
	role := models.Role{
		AccessLevel: input.AccessLevel,
		Name:        input.Name,
	}
	if userRole.AccessLevel != int(gotemplate.SuperAdminRole) {
		return &graphql_models.RolePayload{}, fmt.Errorf("You don't appear to have enough access level for this request ")
	}

	newRole, err := daos.CreateRoleTx(role, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "role")
	}
	return &graphql_models.RolePayload{Role: &graphql_models.Role{
		AccessLevel: newRole.AccessLevel,
		Name:        newRole.Name,
	},
	}, err
}

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*graphql_models.LoginResponse, error) {
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

	return &graphql_models.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (*graphql_models.ChangePasswordResponse, error) {
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
	return &graphql_models.ChangePasswordResponse{Ok: true}, err
}

func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*graphql_models.RefreshTokenResponse, error) {
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
	return &graphql_models.RefreshTokenResponse{Token: resp}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input graphql_models.UserCreateInput) (*graphql_models.UserPayload, error) {
	err := throttle.Check(ctx, 5, 10*time.Second)
	if err != nil {
		return nil, err
	}

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
	graphUser := &graphql_models.User{
		FirstName: convert.NullDotStringToPointerString(newUser.FirstName),
		LastName:  convert.NullDotStringToPointerString(newUser.LastName),
		Username:  convert.NullDotStringToPointerString(newUser.Username),
		Email:     convert.NullDotStringToPointerString(newUser.Email),
		Mobile:    convert.NullDotStringToPointerString(newUser.Mobile),
		Phone:     convert.NullDotStringToPointerString(newUser.Phone),
		Address:   convert.NullDotStringToPointerString(newUser.Address),
	}

	r.Lock()
	for _, observer := range r.Observers {
		observer <- graphUser
	}
	r.Unlock()

	return &graphql_models.UserPayload{User: graphUser}, err
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input *graphql_models.UserUpdateInput) (*graphql_models.UserUpdatePayload, error) {
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

	graphUser := &graphql_models.User{
		FirstName: convert.NullDotStringToPointerString(u.FirstName),
		LastName:  convert.NullDotStringToPointerString(u.LastName),
		Mobile:    convert.NullDotStringToPointerString(u.Mobile),
		Phone:     convert.NullDotStringToPointerString(u.Phone),
		Address:   convert.NullDotStringToPointerString(u.Address),
	}
	r.Lock()
	for _, observer := range r.Observers {
		observer <- graphUser
	}
	r.Unlock()

	return &graphql_models.UserUpdatePayload{Ok: true}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (*graphql_models.UserDeletePayload, error) {
	userID := auth.UserIDFromContext(ctx)
	u, err := daos.FindUserByID(userID)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	_, err = daos.DeleteUser(*u)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "user")
	}
	return &graphql_models.UserDeletePayload{ID: fmt.Sprint(userID)}, nil
}

func (r *queryResolver) Me(ctx context.Context) (*graphql_models.User, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := rediscache.GetUser(userID)
	if err != nil {
		return &graphql_models.User{}, resultwrapper.ResolverSQLError(err, "data")
	}
	return &graphql_models.User{
		FirstName: convert.NullDotStringToPointerString(user.FirstName),
		LastName:  convert.NullDotStringToPointerString(user.LastName),
		Username:  convert.NullDotStringToPointerString(user.Username),
		Email:     convert.NullDotStringToPointerString(user.Email),
		Mobile:    convert.NullDotStringToPointerString(user.Mobile),
		Phone:     convert.NullDotStringToPointerString(user.Phone),
		Address:   convert.NullDotStringToPointerString(user.Address),
	}, err
}

func (r *queryResolver) Users(ctx context.Context, pagination *graphql_models.UserPagination) (*graphql_models.UsersPayload, error) {
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
	return &graphql_models.UsersPayload{Users: convert.UsersToGraphQlUsers(users)}, nil
}

func (r *subscriptionResolver) UserNotification(ctx context.Context) (<-chan *graphql_models.User, error) {
	id := utl.RandomSequence(5)
	event := make(chan *graphql_models.User, 1)

	go func() {
		<-ctx.Done()
		r.Lock()
		delete(r.Observers, id)
		r.Unlock()
	}()

	r.Lock()
	r.Observers[id] = event
	r.Unlock()

	return event, nil
}

// Mutation returns graphql_models.MutationResolver implementation.
func (r *Resolver) Mutation() graphql_models.MutationResolver { return &mutationResolver{r} }

// Query returns graphql_models.QueryResolver implementation.
func (r *Resolver) Query() graphql_models.QueryResolver { return &queryResolver{r} }

// Subscription returns graphql_models.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graphql_models.SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
