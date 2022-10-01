package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	controller "go-template/internal/controller"
	"go-template/internal/middleware/secure"
	"go-template/internal/service/tracer"
	"go-template/pkg/utl/zaplog"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type CustomContext struct {
	echo.Context
	ctx context.Context
}

// New instantates new Echo server
func New() *echo.Echo {
	e := echo.New()
	e.Use(
		otelecho.Middleware(os.Getenv("SERVICE_NAME")),
		middleware.Logger(),
		middleware.Recover(),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				req := c.Request()
				res := c.Response()
				rid := req.Header.Get(echo.HeaderXRequestID)
				if rid == "" {
					rid = random.String(32)
				}
				res.Header().Set(echo.HeaderXRequestID, rid)
				ctx := context.WithValue(c.Request().Context(), zaplog.RequestIdCtxKey, rid)
				c.SetRequest(c.Request().WithContext(ctx))
				cc := &CustomContext{c, ctx}
				return next(cc)
			}
		},
		middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			zaplog.Info(c.Request().Context(), string(reqBody))
			zaplog.Info(c.Request().Context(), string(resBody))
		}),
		secure.Headers(),
		secure.CORS(),
	)
	e.GET("/", controller.HealthCheckHandler)
	e.Validator = &CustomValidator{V: validator.New()}
	custErr := &customErrHandler{e: e}
	e.HTTPErrorHandler = custErr.handler
	e.Binder = &CustomBinder{b: &echo.DefaultBinder{}}
	return e
}

// Config represents server specific config
type Config struct {
	Port                string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	Debug               bool
}

// Start starts echo server
func Start(e *echo.Echo, cfg *Config) {
	tp := tracer.Init()
	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
	e.Debug = cfg.Debug

	// Start server
	go func() {
		zaplog.Logger.Info("Warming up server... ")
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	if err := e.Shutdown(ctx); err != nil {
		e.StdLogger.Fatal(err)
	}
}
