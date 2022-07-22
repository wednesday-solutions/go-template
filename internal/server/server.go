package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-template/internal/middleware/secure"
	"go-template/internal/service/tracer"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// New instantates new Echo server
func New() *echo.Echo {
	e := echo.New()
	e.Use(
		otelecho.Middleware(os.Getenv("SERVICE_NAME")),
		middleware.Logger(),
		middleware.Recover(),
		secure.Headers(),
		secure.CORS(),
	)
	e.GET("/", healthCheck)
	e.Validator = &CustomValidator{V: validator.New()}
	custErr := &customErrHandler{e: e}
	e.HTTPErrorHandler = custErr.handler
	e.Binder = &CustomBinder{b: &echo.DefaultBinder{}}
	return e
}

type response struct {
	Data string `json:"data"`
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, response{Data: "Go template at your service!🍲"})
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
		fmt.Print("Warming up server... ")
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
