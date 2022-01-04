package server

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/wednesday-solutions/go-template/testutls"
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

func Test_customErrHandler_handler(t *testing.T) {

	type args struct {
		errorFunc          func(c echo.Context) error
		expectedStatusCode int
	}
	e := echo.New()
	custErr := &customErrHandler{e: e}
	e.HTTPErrorHandler = custErr.handler

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Faliure_Default",
			args: args{
				expectedStatusCode: http.StatusInternalServerError,
				errorFunc: func(c echo.Context) error {
					return fmt.Errorf("asd")
				},
			},
		},
		{
			name: "Faliure_HttpError",
			args: args{
				expectedStatusCode: http.StatusBadRequest,
				errorFunc: func(c echo.Context) error {
					return &echo.HTTPError{Code: http.StatusBadRequest, Message: "asd"}
				},
			},
		},
		{
			name: "Faliure_ValidationErrors",
			args: args{
				expectedStatusCode: http.StatusBadRequest,
				errorFunc: func(c echo.Context) error {

					fieldError := testutls.NewMockFieldError(gomock.NewController(t))

					fieldError.EXPECT().Field().DoAndReturn(func() string {
						return "FIELD"
					}).AnyTimes()

					fieldError.EXPECT().ActualTag().DoAndReturn(func() string {
						return "ACTUALTAG"
					}).AnyTimes()

					return validator.ValidationErrors{fieldError}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e.GET("/", tt.args.errorFunc)
			_, res, _, err := testutls.MakeAndGetRequest(testutls.RequestParameters{
				Pathname:   "/",
				HttpMethod: "GET",
				E:          e,
			})

			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, tt.args.expectedStatusCode, res.StatusCode)

		})
	}
}
