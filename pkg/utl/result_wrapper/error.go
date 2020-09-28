package resultwrapper

import (
	"context"
	"errors"
	"fmt"
	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/server"
	"net/http"
	"net/http/httptest"
	"strings"
)

var (
	// ErrGeneric is used for testing purposes and for errors handled later in the callstack
	ErrGeneric = errors.New("generic error")

	// ErrBadRequest (400) is returned for bad request (validation)
	ErrBadRequest = echo.NewHTTPError(400)

	// ErrUnauthorized (401) is returned when user is not authorized
	ErrUnauthorized = echo.ErrUnauthorized
)

// ErrorMsg ...
type ErrorMsg struct {
	Errors []string `json:"errors"`
}

// ErrorFormatter ...
func ErrorFormatter(err string) ErrorMsg {
	e := ErrorMsg{
		Errors: []string{err},
	}
	return e
}

// HandleGraphQLError ...
func HandleGraphQLError(errMsg string) graphql2.ResponseHandler {
	return func(ctx context.Context) *graphql2.Response {
		return &graphql2.Response{
			Errors: gqlerror.List{gqlerror.Errorf(errMsg)},
		}
	}
}

// ResolverSQLError ...
func ResolverSQLError(err error, detail string) error {
	if strings.Contains(fmt.Sprint(err), "no rows in result") {
		return ResolverWrapperFromMessage(http.StatusBadRequest, fmt.Sprint("No data found with provided ", detail))
	}
	if strings.Contains(fmt.Sprint(err), "unable to update") {
		return ResolverWrapperFromMessage(http.StatusInternalServerError, fmt.Sprint("Unable to update ", detail))
	}
	if strings.Contains(fmt.Sprint(err), "unable to insert") {
		return ResolverWrapperFromMessage(http.StatusInternalServerError, fmt.Sprint("Unable to save provided ", detail))
	}
	if strings.Contains(fmt.Sprint(err), "delete on table") && strings.Contains(fmt.Sprint(err), "violates foreign key constraint") {
		return ResolverWrapperFromMessage(http.StatusInternalServerError, "Unable to complete the delete operation, it has useful data associated to it")
	}
	return ResolverWrapperFromMessage(http.StatusBadRequest, fmt.Sprint(err))
}

// ResolverWrapperFromMessage ...
func ResolverWrapperFromMessage(errorCode int, err string) error {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "", nil)
	req.Header.Set("Content-Type", "application/json")
	e := echo.New()
	e.Validator = &server.CustomValidator{V: validator.New()}
	e.Binder = server.NewBinder()
	c := e.NewContext(req, w)

	errMessage := fmt.Sprint(err)
	er := ErrorFormatter(errMessage)

	c.Echo().Debug = true

	e1 := c.JSON(errorCode, er)
	if e1 != nil {
		return e1
	}
	return fmt.Errorf(errMessage)
}
