package throttle

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	rediscache "go-template/pkg/utl/rediscache"
	"go-template/testutls"

	"github.com/99designs/gqlgen/graphql"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/ast"
)

type args struct {
	ctx            context.Context
	limit          int
	dur            time.Duration
	isLocal        bool
	visits         int
	visitsErr      error
	startVisitsErr error
	ip             string
}

type testCase struct {
	name    string
	args    args
	wantErr bool
}

// CreateTestCases creates and returns test cases.
func CreateTestCases(ctx context.Context) []testCase {
	return []testCase{
		createSuccessLocalTestCase(ctx),
		createSuccessNotLocalFirstVisitTestCase(ctx),
		createSuccessNotLocalSecondVisitTestCase(ctx),
		createFailureNotLocalFirstVisitTestCase(ctx),
		createFailureNotLocalFirstVisitStartVisitErrTestCase(ctx),
		createFailureNotLocalRateLimitExceededTestCase(ctx),
	}
}

func createSuccessLocalTestCase(ctx context.Context) testCase {
	return testCase{
		name: "Success_Local",
		args: args{
			ctx:     ctx,
			limit:   10,
			dur:     time.Second,
			isLocal: true,
		},
	}
}

func createSuccessNotLocalFirstVisitTestCase(ctx context.Context) testCase {
	return testCase{
		name: "Success_NotLocal_FirstVisit",
		args: args{
			ctx:    ctx,
			limit:  10,
			visits: 1,
			dur:    time.Second,
			ip:     "mock IP address",
		},
	}
}

func createSuccessNotLocalSecondVisitTestCase(ctx context.Context) testCase {
	return testCase{
		name: "Success_NotLocal_SecondVisit",
		args: args{
			ctx:    ctx,
			limit:  10,
			visits: 2,
			dur:    time.Second,
			ip:     "mock IP address",
		},
	}
}

func createFailureNotLocalFirstVisitTestCase(ctx context.Context) testCase {
	return testCase{
		name: "Failure_NotLocal_FirstVisit",
		args: args{
			ctx:       ctx,
			limit:     10,
			visits:    1,
			visitsErr: fmt.Errorf("Internal error"),
			dur:       time.Second,
			ip:        "mock IP address",
		},
		wantErr: true,
	}
}

func createFailureNotLocalFirstVisitStartVisitErrTestCase(ctx context.Context) testCase {
	return testCase{
		name: "Failure_NotLocal_FirstVisit_StartVisitErr",
		args: args{
			ctx:            ctx,
			limit:          10,
			visits:         1,
			startVisitsErr: fmt.Errorf("Internal error"),
			dur:            time.Second,
			ip:             "mock IP address",
		},
		wantErr: true,
	}
}

func createFailureNotLocalRateLimitExceededTestCase(ctx context.Context) testCase {
	return testCase{
		name: "Failure_NotLocal_RateLimitExceeded",
		args: args{
			ctx:    ctx,
			limit:  10,
			visits: 11,
			dur:    time.Second,
			ip:     "mock IP address",
		},
		wantErr: true,
	}
}

// Suggested refactoring to use table-driven tests
func TestCheck(t *testing.T) {
	var ctx context.Context = testutls.MockCtx{}
	tests := CreateTestCases(ctx)
	ApplyFunc(graphql.GetPath, func(ctx context.Context) ast.Path {
		return ast.Path{ast.PathName("users")}
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.ctx = context.WithValue(tt.args.ctx, userIPAdress, tt.args.ip)
			patches := ApplyFunc(os.Getenv, func(key string) string {
				if key == "ENVIRONMENT_NAME" {
					if tt.args.isLocal {
						return "local"
					}
				}
				return ""
			})
			defer patches.Reset()
			ApplyFunc(rediscache.IncVisits, func(path string) (int, error) {
				if tt.args.visitsErr != nil {
					return 0, tt.args.visitsErr
				}
				return tt.args.visits, nil
			})
			ApplyFunc(rediscache.StartVisits, func(path string, exp time.Duration) error {
				return tt.args.startVisitsErr
			})
			if err := Check(tt.args.ctx, tt.args.limit, tt.args.dur); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestGqlMiddleware(t *testing.T) {
	type args struct {
		handler echo.HandlerFunc
	}
	tests := []struct {
		name string
		want echo.MiddlewareFunc
		args args
	}{
		{
			name: "Success",
			want: func(h echo.HandlerFunc) echo.HandlerFunc {
				return nil
			},
			args: args{
				handler: func(c echo.Context) error {
					ipAddress := c.Request().Context().Value(userIPAdress)
					assert.Equal(t, ipAddress, testutls.MockIpAddress)
					return nil
				},
			},
		},
	}
	req, _ := http.NewRequest(
		"POST",
		"/",
		bytes.NewBuffer([]byte("")),
	)
	req.Header.Set("X-Real-IP", testutls.MockIpAddress)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			ctx := e.NewContext(req, nil)
			got := GqlMiddleware()
			next := got(tt.args.handler)

			// assertion is in tt.args.handler
			_ = next(ctx)
		})
	}
}
