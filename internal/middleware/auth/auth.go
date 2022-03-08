package auth

import (
	"context"
	"reflect"

	"go-template/daos"
	"go-template/models"
	resultwrapper "go-template/pkg/utl/result_wrapper"

	graphql2 "github.com/99designs/gqlgen/graphql"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type key string

const (
	authorization key = "Authorization"
)

// TokenParser represents JWT token parser
type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}

// CustomContext ...
type CustomContext struct {
	echo.Context
	ctx context.Context
}

var UserCtxKey = &ContextKey{"user"}

type ContextKey struct {
	Name string
}

// FromContext finds the user from the context. REQUIRES Middleware to have run.
func FromContext(ctx context.Context) *models.User {
	user, _ := ctx.Value(UserCtxKey).(*models.User)
	return user
}

// UserIDFromContext ...
func UserIDFromContext(ctx context.Context) int {
	user := FromContext(ctx)
	if user != nil {
		return user.ID
	}
	return 0
}

// GqlMiddleware ...
func GqlMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), authorization, c.Request().Header.Get(string(authorization)))
			c.SetRequest(c.Request().WithContext(ctx))
			cc := &CustomContext{c, ctx}
			return next(cc)
		}
	}
}

// WhiteListedQueries ...
var WhiteListedQueries = []string{"__schema", "introspectionquery", "login"}

// AdminQueries ...
var AdminQueries = []string{"users"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GraphQLMiddleware ...
func GraphQLMiddleware(
	ctx context.Context,
	tokenParser TokenParser,
	next graphql2.OperationHandler) graphql2.ResponseHandler {

	operationContext := graphql2.GetOperationContext(ctx)
	var needsAuth = false
	var requiresSuperAdmin = false
	for _, selectionSet := range operationContext.Operation.SelectionSet {

		selection := reflect.ValueOf(selectionSet).Elem()
		if !contains(WhiteListedQueries, selection.FieldByName("Name").Interface().(string)) {
			needsAuth = true
		}
		if contains(AdminQueries, selection.FieldByName("Name").Interface().(string)) {
			requiresSuperAdmin = true
		}
	}

	if needsAuth || requiresSuperAdmin {
		// strip token
		var tokenStr = ctx.Value(authorization).(string)
		if len(tokenStr) == 0 {
			return resultwrapper.HandleGraphQLError("Authorization header is missing")
		}

		token, err := tokenParser.ParseToken(tokenStr)

		if err != nil || !token.Valid {
			return resultwrapper.HandleGraphQLError("Invalid authorization token")
		}
		claims := token.Claims.(jwt.MapClaims)
		if requiresSuperAdmin {

			isSuperAdmin := false

			if claims["role"].(string) == "SUPER_ADMIN" {
				isSuperAdmin = true
			}
			if !isSuperAdmin {
				return resultwrapper.HandleGraphQLError("Unauthorized! \n Only admins are authorized to make this request.")
			}
		}

		email := claims["e"].(string)
		user, err := daos.FindUserByEmail(email)

		if err != nil {
			return resultwrapper.HandleGraphQLError("No user found for this email address")
		}

		ctx = context.WithValue(ctx, UserCtxKey, user)
		return next(ctx)
	}

	return next(ctx)
}
