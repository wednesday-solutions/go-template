package zaplog

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const (
	ErrorFromProduction = "Error  producton"
	InformationTest     = "test info"
	InfoMessage         = "This is an info log"
	DebugMessage        = "This is a debug log"
	NewCase             = "New Case"
)

type SugaredLogger struct {
	*zap.SugaredLogger
}

func TestInfo(t *testing.T) {
	type args struct {
		c   context.Context
		msg string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: InformationTest,
			args: args{
				c:   context.Background(),
				msg: InfoMessage,
			},
		},
		{
			name: ErrorFromProduction,
			args: args{
				c:   context.Background(),
				msg: InfoMessage,
			},
		},
		{
			name: NewCase,
			args: args{
				c:   context.Background(),
				msg: InfoMessage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			observedLogger := zap.New(observedZapCore).Sugar()

			_ = SetLogger(observedLogger)
			Info(tt.args.c, tt.args.msg)
			assert.Equal(t, 1, observedLogs.Len())
			log := observedLogs.All()[0]
			assert.Equal(t, fmt.Sprintf("\nRequest-ID: <nil>\n[%s]\n", tt.args.msg), log.Message)
			assert.Equal(t, zapcore.Level(0), log.Level)
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		c   context.Context
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: InformationTest,
			args: args{
				c:   context.Background(),
				msg: DebugMessage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			observedZapCore, observedLogs := observer.New(zap.DebugLevel)
			observedLogger := zap.New(observedZapCore).Sugar()
			_ = SetLogger(observedLogger)
			Debug(tt.args.c, tt.args.msg)
			assert.Equal(t, 1, observedLogs.Len())
			log := observedLogs.All()[0]
			assert.Equal(t, fmt.Sprintf("\nRequest-ID: <nil>\n[%s]\n", tt.args.msg), log.Message)
			assert.Equal(t, zapcore.Level(-1), log.Level)
		})
	}
}
