package server

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"go-template/testutls"

	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func Test_getVldErrorMsg(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Failure_FailedOnValidation",
			args: args{
				s: "a",
			},
			want: " failed on a validation",
		},
		{
			name: "Failure_FailedOnRequired",
			args: args{
				s: "required",
			},
			want: " is required, but was not received",
		},
		{
			name: "Failure_FailedOnMin",
			args: args{
				s: "required",
			},
			want: " is required, but was not received",
		},
		{
			name: "Failure_FailedOnMax",
			args: args{
				s: "max",
			},
			want: "'s value or length is bigger than allowed",
		},
		{
			name: "Failure_FailedOnMin",
			args: args{
				s: "min",
			},
			want: "'s value or length is less than allowed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getVldErrorMsg(tt.args.s); got != tt.want {
				t.Errorf("getVldErrorMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getValidatorErr(t *testing.T) error {

	fieldError := testutls.NewMockFieldError(gomock.NewController(t))

	fieldError.EXPECT().Field().DoAndReturn(func() string {
		return "FIELD"
	}).AnyTimes()

	fieldError.EXPECT().ActualTag().DoAndReturn(func() string {
		return "ACTUALTAG"
	}).AnyTimes()

	return validator.ValidationErrors{fieldError}
}
func Test_customErrHandler_handler(t *testing.T) {

	type args struct {
		err                error
		expectedStatusCode int
		httpMethod         string
	}
	e := echo.New()
	custErr := &customErrHandler{e: e}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Faliure_Default",
			args: args{
				expectedStatusCode: http.StatusInternalServerError,
				err:                fmt.Errorf("asd"),
				httpMethod:         "GET",
			},
		},
		{
			name: "Faliure_NoContent",
			args: args{
				expectedStatusCode: http.StatusInternalServerError,
				err:                fmt.Errorf("asd"),
				httpMethod:         "HEAD",
			},
		},
		{
			name: "Faliure_HttpError",
			args: args{
				expectedStatusCode: http.StatusBadRequest,
				err:                &echo.HTTPError{Code: http.StatusBadRequest, Message: "asd", Internal: fmt.Errorf("asd")},
				httpMethod:         "GET",
			},
		},
		{
			name: "Faliure_ValidationErrors",
			args: args{
				expectedStatusCode: http.StatusBadRequest,
				err:                getValidatorErr(t),
				httpMethod:         "GET",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// just make a request
			req, _ := http.NewRequest(
				tt.args.httpMethod,
				"/",
				bytes.NewBuffer([]byte("")),
			)

			ctx := testutls.NewMockContext(gomock.NewController(t))

			// mock ctx.Response
			ctx.
				EXPECT().
				Response().
				DoAndReturn(func() *echo.Response {
					return &echo.Response{Status: tt.args.expectedStatusCode}
				}).
				AnyTimes()

			// mock ctx.Request
			ctx.
				EXPECT().
				Request().
				DoAndReturn(func() *http.Request {
					return req
				}).
				AnyTimes()

			if tt.args.httpMethod == "HEAD" {
				// mock ctx.NoContent
				ctx.
					EXPECT().
					NoContent(gomock.Eq(tt.args.expectedStatusCode)).
					DoAndReturn(func(code int) error {
						return nil
					}).
					AnyTimes()
			} else {
				// mock ctx.JSON
				ctx.
					EXPECT().
					JSON(gomock.Eq(tt.args.expectedStatusCode), gomock.Any()).
					DoAndReturn(func(code int, i interface{}) error {
						return fmt.Errorf("error")
					}).
					AnyTimes()
			}

			// call the handler with tt.args.err. We are asserting in the JSON/NoContent call
			custErr.handler(tt.args.err, ctx)
		})
	}
}
