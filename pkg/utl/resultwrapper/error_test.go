package resultwrapper_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	resultwrapper "go-template/pkg/utl/resultwrapper"
	"go-template/testutls"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	SuccessCase = "Success"
	ErrorCase   = "error from json"
	ErrMsgJSON  = "Error from JSON"
)

func TestSplitByLabel(t *testing.T) {

	cases := []struct {
		name     string
		req      string
		wantResp string
	}{
		{
			name:     "error string",
			req:      "no rows in sql",
			wantResp: "no rows in sql",
		},
		{
			name:     "having `Error` in string",
			req:      `"Error":{"msg"}`,
			wantResp: "\":{\"msg\"}",
		},
		{
			name:     "having `message` in string",
			req:      `"message":{"msg"}`,
			wantResp: "\":{\"msg\"}",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp := resultwrapper.SplitByLabel(tt.req)
			if len(tt.wantResp) != 0 {
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestErrorFormatter(t *testing.T) {

	cases := []struct {
		name     string
		req      string
		wantResp resultwrapper.ErrorMsg
	}{
		{
			name:     "No string",
			req:      "",
			wantResp: resultwrapper.ErrorMsg{Errors: []string{""}},
		},
		{
			name:     "Having Error string",
			req:      `error message`,
			wantResp: resultwrapper.ErrorMsg{Errors: []string{"error message"}},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp := resultwrapper.ErrorFormatter(tt.req)
			assert.Equal(t, tt.wantResp, resp)
		})
	}
}

func getContext() echo.Context {
	e := echo.New()
	w := httptest.NewRecorder()
	req, _, _, _ := testutls.MakeAndGetRequest(testutls.RequestParameters{
		Pathname:   "/",
		HttpMethod: "GET",
		E:          e,
	})
	return e.NewContext(req, w)
}

func TestResultWrapper(t *testing.T) {
	type args struct {
		errorCode int
		err       error
	}
	e := echo.New()
	req, _, _, _ := testutls.MakeAndGetRequest(testutls.RequestParameters{
		Pathname:   "/",
		HttpMethod: "GET",
		E:          e,
	})

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: SuccessCase,
			args: args{
				errorCode: 400,
				err:       fmt.Errorf("sample error"),
			},
			wantErr: true,
		},
		{
			name: ErrorCase,
			args: args{
				errorCode: 400,
				err:       fmt.Errorf(ErrMsgJSON),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx := e.NewContext(req, w)

			if tt.name == ErrorCase {

				patch := gomonkey.ApplyMethodFunc(ctx, "JSON", func(code int, i interface{}) error {
					return fmt.Errorf(ErrMsgJSON)
				})
				defer patch.Reset()
			}

			err := resultwrapper.ResultWrapper(tt.args.errorCode, ctx, tt.args.err)
			if tt.name == ErrorCase {
				assert.Equal(t, tt.args.err, err)
			} else {
				assert.Equal(t, tt.args.errorCode, ctx.Response().Status)
			}

		})
	}
}

func TestInternalServerError(t *testing.T) {
	type args struct {
		c   echo.Context
		err error
	}
	errorStr := "This is an error"
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			name: SuccessCase,
			err:  errorStr,
			args: args{
				err: fmt.Errorf(errorStr),
				c:   getContext()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.InternalServerError(tt.args.c, tt.args.err)
			assert.Equal(t, http.StatusInternalServerError, tt.args.c.Response().Status)
			assert.Equal(t, err.Error(), tt.err)
		})

	}
}

func TestInternalServerErrorFromMessage(t *testing.T) {
	type args struct {
		c   echo.Context
		err string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: SuccessCase,
			args: args{
				err: "This is an error",
				c:   getContext(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.InternalServerErrorFromMessage(tt.args.c, tt.args.err)
			assert.Equal(t, err, fmt.Errorf(tt.args.err))
			assert.Equal(t, http.StatusInternalServerError, tt.args.c.Response().Status)
		})
	}
}

func TestBadRequest(t *testing.T) {
	type args struct {
		c   echo.Context
		err error
	}
	errorStr := "This is an error"
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			name: SuccessCase,
			err:  errorStr,
			args: args{
				err: fmt.Errorf(errorStr),
				c:   getContext()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.BadRequest(tt.args.c, tt.args.err)
			assert.Equal(t, http.StatusBadRequest, tt.args.c.Response().Status)
			assert.Equal(t, err.Error(), tt.err)
		})

	}
}

func TestBadRequestFromMessage(t *testing.T) {
	type args struct {
		c   echo.Context
		err string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: SuccessCase,
			args: args{
				err: "This is an error",
				c:   getContext(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.BadRequestFromMessage(tt.args.c, tt.args.err)
			assert.Equal(t, err, fmt.Errorf(tt.args.err))
			assert.Equal(t, http.StatusBadRequest, tt.args.c.Response().Status)
		})
	}
}

func TestConflict(t *testing.T) {
	type args struct {
		c   echo.Context
		err error
	}
	errorStr := "This is an error"
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			name: SuccessCase,
			err:  errorStr,
			args: args{
				err: fmt.Errorf(errorStr),
				c:   getContext()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.Conflict(tt.args.c, tt.args.err)
			assert.Equal(t, http.StatusConflict, tt.args.c.Response().Status)
			assert.Equal(t, err.Error(), tt.err)
		})

	}
}

func TestConflictFromMessage(t *testing.T) {
	type args struct {
		c   echo.Context
		err string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: SuccessCase,
			args: args{
				err: "This is an error",
				c:   getContext(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.ConflictFromMessage(tt.args.c, tt.args.err)
			assert.Equal(t, err, fmt.Errorf(tt.args.err))
			assert.Equal(t, http.StatusConflict, tt.args.c.Response().Status)
		})
	}
}

func TestTooManyRequests(t *testing.T) {
	type args struct {
		c   echo.Context
		err error
	}
	errorStr := "This is an error"
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			name: SuccessCase,
			err:  errorStr,
			args: args{
				err: fmt.Errorf(errorStr),
				c:   getContext()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.TooManyRequests(tt.args.c, tt.args.err)
			assert.Equal(t, http.StatusTooManyRequests, tt.args.c.Response().Status)
			assert.Equal(t, err.Error(), tt.err)
		})

	}
}

func TestUnauthorized(t *testing.T) {
	type args struct {
		c   echo.Context
		err error
	}
	errorStr := "This is an error"
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			name: SuccessCase,
			err:  errorStr,
			args: args{
				err: fmt.Errorf(errorStr),
				c:   getContext()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.Unauthorized(tt.args.c, tt.args.err)
			assert.Equal(t, http.StatusUnauthorized, tt.args.c.Response().Status)
			assert.Equal(t, err.Error(), tt.err)
		})

	}
}

func TestUnauthorizedFromMessage(t *testing.T) {
	type args struct {
		c   echo.Context
		err string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: SuccessCase,
			args: args{
				err: "This is an error",
				c:   getContext(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.UnauthorizedFromMessage(tt.args.c, tt.args.err)
			assert.Equal(t, err, fmt.Errorf(tt.args.err))
			assert.Equal(t, http.StatusUnauthorized, tt.args.c.Response().Status)
		})
	}
}

func TestNoDataFound(t *testing.T) {
	type args struct {
		c   echo.Context
		err error
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
		statusCode  int
	}{
		{
			name: "Success_DuplicateData",
			args: args{
				c:   getContext(),
				err: fmt.Errorf("duplicate key value violates unique constraint"),
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("Data already exists. duplicate key value violates unique constraint"),
			statusCode:  http.StatusConflict,
		},
		{
			name: "Success_NoData",
			args: args{
				c:   getContext(),
				err: fmt.Errorf("no rows in result"),
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to fetch data"),
			statusCode:  http.StatusInternalServerError,
		},
		{
			name: "Success_ServerError",
			args: args{
				c:   getContext(),
				err: fmt.Errorf("random error"),
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("random error"),
			statusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.NoDataFound(tt.args.c, tt.args.err)
			assert.Equal(t, err, tt.expectedErr)
			assert.Equal(t, tt.statusCode, tt.args.c.Response().Status)

		})
	}
}

func TestServiceUnavailable(t *testing.T) {
	type args struct {
		c   echo.Context
		err error
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
		errCode     int
	}{
		{
			name: "Success_ServiceUnavailable",
			args: args{
				c:   getContext(),
				err: fmt.Errorf("Service Unavailable"),
			},
			expectedErr: fmt.Errorf("unable to complete process"),
			errCode:     http.StatusServiceUnavailable,
		},
		{
			name: "Success_NoData",
			args: args{
				c:   getContext(),
				err: fmt.Errorf("Random error"),
			},
			expectedErr: fmt.Errorf("Random error"),
			errCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.ServiceUnavailable(tt.args.c, tt.args.err)
			assert.Equal(t, tt.errCode, tt.args.c.Response().Status)
			assert.Equal(t, err, tt.expectedErr)

		})
	}
}

func TestHandleGraphQLError(t *testing.T) {
	type args struct {
		errMsg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: SuccessCase,
			args: args{
				errMsg: "This is an error",
			},
			want: "This is an error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resultwrapper.HandleGraphQLError(tt.args.errMsg)
			c := getContext()
			res := got(c.Request().Context())
			for _, v := range res.Errors {
				assert.Equal(t, tt.want, v.Message)
			}
		})
	}
}

func TestResolverSQLError(t *testing.T) {
	type args struct {
		err    error
		detail string
	}
	tests := []struct {
		name          string
		args          args
		errMsg        string
		dontAddDetail bool
	}{
		{
			name: SuccessCase,
			args: args{
				err:    fmt.Errorf("this is some error"),
				detail: "Some level of detail",
			},
			errMsg:        "this is some error",
			dontAddDetail: true,
		},
		{
			name: "Success_NoData",
			args: args{
				detail: "Some level of detail",
				err:    fmt.Errorf("no rows in result"),
			},
			errMsg: "No data found with provided",
		},
		{
			name: "Success_UnableToUpdate",
			args: args{
				detail: "Some level of detail",
				err:    fmt.Errorf("unable to update"),
			},
			errMsg: "Unable to update",
		},
		{
			name: "Success_UnableToInsert",
			args: args{
				detail: "Some level of detail",
				err:    fmt.Errorf("unable to insert"),
			},
			errMsg: "Unable to save provided",
		},
		{
			name: "Success_UnableToDelete",
			args: args{
				detail: "Some level of detail",
				err:    fmt.Errorf("delete on table. violates foreign key constraint"),
			},
			errMsg:        "Unable to complete the delete operation, it has useful data associated to it",
			dontAddDetail: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := resultwrapper.ResolverSQLError(tt.args.err, tt.args.detail)
			errorMessage := fmt.Sprintf("%s %s", tt.errMsg, tt.args.detail)
			if tt.dontAddDetail {
				errorMessage = tt.errMsg
			}
			assert.Equal(t, fmt.Errorf(errorMessage), err)
		})
	}
}
