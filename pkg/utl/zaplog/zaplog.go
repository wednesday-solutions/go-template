package zaplog

import (
	"context"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var RequestIdCtxKey = &ContextKey{echo.HeaderXRequestID}

type ContextKey struct {
	Name string
}

var Logger = InitLogger()

func SetLogger(logger *zap.SugaredLogger) *zap.SugaredLogger {
	Logger = logger
	return Logger
}

func getRequestID(c context.Context) string {
	return fmt.Sprintf("\nRequest-ID: %v", c.Value(RequestIdCtxKey))
}
func Info(c context.Context, args ...interface{}) {
	Logger.Info(getRequestID(c), "\n", args, "\n")
}
func Debug(c context.Context, args ...interface{}) {
	Logger.Debug(getRequestID(c), "\n", args, "\n")
}
func InitLogger() *zap.SugaredLogger {
	var zapLogger *zap.Logger
	var err error
	if os.Getenv("ENVIRONMENT_NAME") == "production" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	return zapLogger.Sugar()

}
