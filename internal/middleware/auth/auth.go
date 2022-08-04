package auth

import (
	"context"
	"reflect"

	"go-template/daos"
	"go-template/models"
	resultwrapper "go-template/pkg/utl/resultwrapper"

	graphql2 "github.com/99designs/gqlgen/graphql"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/vektah/gqlparser/v2/ast"
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
			ctx := context.WithValue(
				c.Request().Context(),
				authorization,
				c.Request().Header.Get(string(authorization)),
			)
			c.SetRequest(c.Request().WithContext(ctx))
			cc := &CustomContext{c, ctx}
			return next(cc)
		}
	}
}

// WhiteListedOperations...
var WhiteListedOperations = map[string][]string{
	"query":        {"__schema", "introspectionquery", "userNotification"},
	"mutation":     {"login"},
	"subscription": {"userNotification"},
}

// AdminOperations...
var AdminOperations = map[string][]string{
	"query": {"users"},
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getAccessNeeds(operation *ast.OperationDefinition) (needsAuthAccess bool, needsSuperAdminAccess bool) {
	operationName := string(operation.Operation)
	for _, selectionSet := range operation.SelectionSet {

		selection := reflect.ValueOf(selectionSet).Elem()
		if !contains(WhiteListedOperations[operationName], selection.FieldByName("Name").Interface().(string)) {
			needsAuthAccess = true
		}
		if contains(AdminOperations[operationName], selection.FieldByName("Name").Interface().(string)) {
			needsSuperAdminAccess = true
		}
	}
	return needsAuthAccess, needsSuperAdminAccess
}

// GraphQLMiddleware ...
func GraphQLMiddleware(
	ctx context.Context,
	tokenParser TokenParser,
	next graphql2.OperationHandler) graphql2.ResponseHandler {

	operation := graphql2.GetOperationContext(ctx).Operation
	needsAuthAccess, needsSuperAdminAccess := getAccessNeeds(operation)
	if !needsAuthAccess && !needsSuperAdminAccess {
		return next(ctx)
	}

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

	if needsSuperAdminAccess && claims["role"].(string) != "SUPER_ADMIN" {
		return resultwrapper.HandleGraphQLError(
			"Unauthorized! \n Only admins are authorized to make this request.",
		)
	}

	email := claims["e"].(string)
	user, err := daos.FindUserByEmail(email, ctx)

	if err != nil {
		return resultwrapper.HandleGraphQLError("No user found for this email address")
	}

	ctx = context.WithValue(ctx, UserCtxKey, user)
	return next(ctx)

}
