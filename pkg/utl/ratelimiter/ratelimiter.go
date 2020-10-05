package ratelimiter

import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/labstack/echo"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
	"time"
)

// LimitMiddleware ...
func LimitMiddleware(lmt *limiter.Limiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			httpError := tollbooth.LimitByRequest(lmt, c.Response(), c.Request())
			if httpError != nil {
				return resultwrapper.TooManyRequests(c, fmt.Errorf(httpError.Message))
			}
			return next(c)
		})
	}
}

// LimitHandler ...
func LimitHandler(lmt *limiter.Limiter) echo.MiddlewareFunc {
	return LimitMiddleware(lmt)
}

// RateHandler ...
func RateHandler(burstLimit int) echo.MiddlewareFunc {

	lmt := tollbooth.NewLimiter(0.0042, &limiter.ExpirableOptions{ExpireJobInterval: time.Hour}).SetBurst(burstLimit)
	lmt.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})

	return LimitHandler(lmt)
}
