package resultwrapper

import (
	"context"
	"errors"
	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/labstack/echo"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	// ErrGeneric is used for testing purposes and for errors handled later in the callstack
	ErrGeneric = errors.New("generic error")

	// ErrBadRequest (400) is returned for bad request (validation)
	ErrBadRequest = echo.NewHTTPError(400)

	// ErrUnauthorized (401) is returned when user is not authorized
	ErrUnauthorized = echo.ErrUnauthorized
)

// HandleGraphQLError ...
func HandleGraphQLError(errMsg string) graphql2.ResponseHandler {
	return func(ctx context.Context) *graphql2.Response {
		return &graphql2.Response{
			Errors: gqlerror.List{gqlerror.Errorf(errMsg)},
		}
	}
}
