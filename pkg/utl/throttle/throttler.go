package throttle

import (
	"context"
	"fmt"
	"os"
	"time"

	rediscache "go-template/pkg/utl/rediscache"

	"github.com/99designs/gqlgen/graphql"
	"github.com/labstack/echo/v4"
)

type key string

const (
	userIPAdress key = "userIPAdress"
)

// Check function checks weather the given IP address has already
// tried a given query path 'limit' number of times within past 'dur'
func Check(ctx context.Context, limit int, dur time.Duration) error {
	// disabled throttler in 'local' stage
	if os.Getenv("ENVIRONMENT_NAME") == "local" {
		return nil
	}

	query := graphql.GetPath(ctx).String()
	ip := ctx.Value(userIPAdress).(string)
	key := fmt.Sprintf("rate-limit-%s-%s", query, ip)

	num, err := rediscache.IncVisits(key)
	if err != nil {
		return fmt.Errorf("Internal error")
	}

	if num > limit {
		return fmt.Errorf("You reached the rate limit for this query")
	} else if num == 1 {
		err := rediscache.StartVisits(key, dur)
		if err != nil {
			return fmt.Errorf("Internal error")
		}
	}

	return nil
}

// GqlMiddleware returns a middleware that takes IP address
// from echo context and place it in the context of gqlgen resolvers.
func GqlMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), userIPAdress, c.RealIP())
			c.SetRequest(c.Request().WithContext(ctx))
			cc := &struct {
				echo.Context
				ctx context.Context
			}{c, ctx}
			return next(cc)
		}
	}
}
