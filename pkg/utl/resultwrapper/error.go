package resultwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"go-template/internal/server"
	"go-template/pkg/utl/zaplog"

	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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

// ErrorMsgLabels ...
type ErrorMsgLabels struct {
	Label string
}

// SplitByLabel ...
func SplitByLabel(errStr string) string {
	labelArray := []ErrorMsgLabels{{"Error"}, {"message"}}

	r := regexp.MustCompile("code=[0-9]*")
	_ = strings.Replace(errStr, r.FindString(errStr), "", -1)

	for _, ele := range labelArray {
		if strings.Contains(errStr, ele.Label) {
			re, _ := regexp.Compile(fmt.Sprintf("(?i)%s[:|=]*", ele.Label))

			return strings.SplitAfter(errStr, re.FindString(errStr))[1]
		}
	}
	return errStr
}

// ResultWrapper ...
func ResultWrapper(errorCode int, c echo.Context, err error) error {
	errMessage := fmt.Sprint(err)
	return WrapperFromMessage(errorCode, c, errMessage)
}

// WrapperFromMessage ...
func WrapperFromMessage(errorCode int, c echo.Context, err string) error {
	errMessage := fmt.Sprint(err)
	e := ErrorFormatter(errMessage)
	e1 := c.JSON(errorCode, e)
	if e1 != nil {
		return e1
	}
	return fmt.Errorf(errMessage)
}

// InternalServerError ...
func InternalServerError(c echo.Context, err error) error {
	return WrapperFromMessage(http.StatusInternalServerError, c, SplitByLabel(fmt.Sprint(err)))
}

// InternalServerErrorFromMessage ...
func InternalServerErrorFromMessage(c echo.Context, err string) error {
	return WrapperFromMessage(http.StatusInternalServerError, c, err)
}

// BadRequest ...
func BadRequest(c echo.Context, err error) error {
	return WrapperFromMessage(http.StatusBadRequest, c, SplitByLabel(fmt.Sprint(err)))
}

// BadRequestFromMessage ...
func BadRequestFromMessage(c echo.Context, err string) error {
	return WrapperFromMessage(http.StatusBadRequest, c, err)
}

// Conflict ...
func Conflict(c echo.Context, err error) error {
	return WrapperFromMessage(http.StatusConflict, c, SplitByLabel(fmt.Sprint(err)))
}

// ConflictFromMessage ...
func ConflictFromMessage(c echo.Context, err string) error {
	return WrapperFromMessage(http.StatusConflict, c, err)
}

// TooManyRequests ...
func TooManyRequests(c echo.Context, err error) error {
	return WrapperFromMessage(http.StatusTooManyRequests, c, fmt.Sprint(err))
}

// Unauthorized ...
func Unauthorized(c echo.Context, err error) error {
	return WrapperFromMessage(http.StatusUnauthorized, c, SplitByLabel(fmt.Sprint(err)))
}

// UnauthorizedFromMessage ...
func UnauthorizedFromMessage(c echo.Context, err string) error {
	return WrapperFromMessage(http.StatusUnauthorized, c, err)
}

// NoDataFound ...
func NoDataFound(c echo.Context, err error) error {
	substr := "duplicate key value violates unique constraint"
	if strings.Contains(fmt.Sprintf("%s", err), substr) {
		return ConflictFromMessage(c, "Data already exists. "+err.Error())
	}
	substr = "no rows in result"
	if strings.Contains(fmt.Sprintf("%s", err), substr) {
		return InternalServerErrorFromMessage(c, "failed to fetch data")
	}
	return BadRequest(c, err)
}

// ServiceUnavailable ...
func ServiceUnavailable(c echo.Context, err error) error {
	substr := "Service Unavailable"
	if strings.Contains(fmt.Sprintf("%s", err), substr) {
		return WrapperFromMessage(http.StatusServiceUnavailable, c, "unable to complete process")
	}
	return NoDataFound(c, err)
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
	zaplog.Logger.Info(err.Error(), detail)
	if strings.Contains(fmt.Sprint(err), "no rows in result") {
		return ResolverWrapperFromMessage(http.StatusBadRequest, fmt.Sprint("No data found with provided ", detail))
	}
	if strings.Contains(fmt.Sprint(err), "unable to update") {
		return ResolverWrapperFromMessage(http.StatusInternalServerError, fmt.Sprint("Unable to update ", detail))
	}
	if strings.Contains(fmt.Sprint(err), "unable to insert") {
		return ResolverWrapperFromMessage(http.StatusInternalServerError, fmt.Sprint("Unable to save provided ", detail))
	}
	if strings.Contains(fmt.Sprint(err), "delete on table") &&
		strings.Contains(fmt.Sprint(err), "violates foreign key constraint") {

		return ResolverWrapperFromMessage(http.StatusInternalServerError,
			"Unable to complete the delete operation, it has useful data associated to it")
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
