package secure

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Headers adds general security headers for basic security measures
func Headers() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Protects from MimeType Sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			// Prevents browser from prefetching DNS
			c.Response().Header().Set("X-DNS-Prefetch-Control", "off")
			// Denies website content to be served in an iframe
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("Strict-Transport-Security", "max-age=5184000; includeSubDomains")
			// Prevents Internet Explorer from executing downloads in site's context
			c.Response().Header().Set("X-Download-Options", "noopen")
			// Minimal XSS protection
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			return next(c)
		}
	}
}

// CORS adds Cross-Origin Resource Sharing support
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		MaxAge:       86400,
		AllowMethods: []string{
			http.MethodHead,
			http.MethodPost,
			http.MethodGet,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodPut,
		},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
