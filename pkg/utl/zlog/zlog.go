package zlog

import (
	"context"
	"os"

	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var Logger = initLogger()

type logger struct {
	*zap.SugaredLogger
	// error occurred during creation of logger
	InitError error
}

func Info(c context.Context, msg ...string) {
	Logger.Info("\nRequest-ID: ", c.Value(echo.HeaderXRequestID), "\n", msg)
}
func Debug(c context.Context, msg ...string) {
	Logger.Debug(c.Value(echo.HeaderXRequestID), msg)
}
func initLogger() *logger {
	var zapLogger *zap.Logger
	var err error
	if os.Getenv("ENVIRONMENT_NAME") == "production" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}

	return &logger{
		zapLogger.Sugar(),
		err,
	}
}
