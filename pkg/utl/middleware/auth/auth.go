package auth

import (
	"context"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-boiler/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
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

func FromContextWithCheck(c echo.Context) (*models.User, bool) {
	user, exists := c.Get("user").(*models.User)
	return user, exists
}

func ExistsInContext(ctx context.Context) bool {
	_, exist := ctx.Value(userCtxKey).(*models.User)
	return exist
}

// FromContext finds the user from the context. REQUIRES Middleware to have run.
func FromContext(ctx context.Context) *models.User {
	user, _ := ctx.Value(userCtxKey).(*models.User)
	return user
}

func UserIDFromContext(ctx context.Context) int {
	user := FromContext(ctx)
	if user != nil {
		return user.ID
	}
	return 0
}
