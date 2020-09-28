package auth

import (
	"context"
	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-boiler/daos"
	"github.com/wednesday-solutions/go-boiler/models"
	resultwrapper "github.com/wednesday-solutions/go-boiler/pkg/utl/result_wrapper"
	"net/http"
	"reflect"
)

type key string

const (
	authorization key = "Authorization"
)

// TokenParser represents JWT token parser
type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}

// Middleware makes JWT implement the Middleware interface.
func Middleware(tokenParser TokenParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if len(c.Request().Header.Get("Authorization")) == 0 && c.Request().RequestURI == "/graphql" && (c.Request().Header.Get("Referer") == "http://localhost:9000/playground" || c.Request().Header.Get("Referer") == "http://127.0.0.1:9000/playground") {
				return next(c)
			}
			token, err := tokenParser.ParseToken(c.Request().Header.Get("Authorization"))
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims := token.Claims.(jwt.MapClaims)

			id := int(claims["id"].(float64))
			companyID := int(claims["c"].(float64))
			locationID := int(claims["l"].(float64))
			username := claims["u"].(string)
			email := claims["e"].(string)
			user, err := models.Users(qm.Where("email=?", email)).One(context.Background(), boil.GetContextDB())
			if err != nil {
				return err
			}

			//ctx := context.WithValue(c.Request().Context(), userCtxKey, user)
			ctx := context.WithValue(c.Request().Context(), userCtxKey, user)
			c.SetRequest(c.Request().WithContext(ctx))

			//cc := &CustomContext{c, ctx}

			c.Set("id", id)
			c.Set("company_id", companyID)
			c.Set("location_id", locationID)
			c.Set("username", username)
			c.Set("email", email)
			c.Set("user", user)

			//c.Request().Context().Value(user)
			cc := &CustomContext{c, ctx}
			return next(cc)
		}
	}
}

// CustomContext ...
type CustomContext struct {
	echo.Context
	ctx context.Context
}

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// FromContextWithCheck ...
func FromContextWithCheck(c echo.Context) (*models.User, bool) {
	user, exists := c.Get("user").(*models.User)
	return user, exists
}

// ExistsInContext ...
func ExistsInContext(ctx context.Context) bool {
	_, exist := ctx.Value(userCtxKey).(*models.User)
	return exist
}

// FromContext finds the user from the context. REQUIRES Middleware to have run.
func FromContext(ctx context.Context) *models.User {
	user, _ := ctx.Value(userCtxKey).(*models.User)
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GraphQLMiddleware ...
func GraphQLMiddleware(ctx context.Context, tokenParser TokenParser, next graphql2.OperationHandler) graphql2.ResponseHandler {
	operationContext := graphql2.GetOperationContext(ctx)

	var needsAuth = false
	for _, selectionSet := range operationContext.Operation.SelectionSet {

		selection := reflect.ValueOf(selectionSet).Elem()
		if !contains(WhiteListedQueries, selection.FieldByName("Name").Interface().(string)) {
			needsAuth = true
		}
	}

	if needsAuth {
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

		email := claims["e"].(string)
		user, err := daos.FindUserByEmail(email)
		if err != nil {
			return resultwrapper.HandleGraphQLError("No user found for this email address")
		}

		ctx = context.WithValue(ctx, userCtxKey, user)
		return next(ctx)
	}

	return next(ctx)
}
