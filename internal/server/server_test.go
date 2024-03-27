package server_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"testing"
	"time"

	"go-template/internal/server"
	"go-template/testutls"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Improve tests
func TestNew(t *testing.T) {
	e := server.New()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
	assert.NotEmpty(t, e)
	response, err := testutls.MakeRequest(
		testutls.RequestParameters{
			E:          e,
			Pathname:   "/",
			HttpMethod: "GET",
		},
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, response["data"], "Go template at your service!🍲")
}

type args struct {
	e                    *echo.Echo
	cfg                  *server.Config
	startServer          func(e *echo.Echo, s *http.Server) (err error)
	startServerCalled    bool
	serverShutDownCalled bool
	shutDownFailed       bool
}

func initValues(startServer func(e *echo.Echo, s *http.Server) error) args {
	config := testutls.MockConfig()
	return args{
		e: server.New(),
		cfg: &server.Config{
			Port:                config.Server.Port,
			ReadTimeoutSeconds:  config.Server.ReadTimeout,
			WriteTimeoutSeconds: config.Server.WriteTimeout,
			Debug:               config.Server.Debug,
		},
		startServer:    startServer,
		shutDownFailed: true,
	}
}
func TestStart(t *testing.T) {
	cases := getTestCases()

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			mockStarServerPatches := mockStartServer(&tt.args)
			mockShutDownPatches, mockSdLoggerPatches := mockShutdownIfNeeded(&tt.args)
			startServerAndInterrupt(tt.args)
			waitForServerShutdownIfNeeded(tt.args)
			assertions(t, tt.args)

			mockStarServerPatches.Reset()
			mockShutDownPatches.Reset()
			mockSdLoggerPatches.Reset()
		})
	}
}

type testCase struct {
	args args
}

func getTestCases() map[string]testCase {
	return map[string]testCase{
		"Success": {
			args: initValues(func(e *echo.Echo, s *http.Server) (err error) {
				return nil
			}),
		},
		"Failure_ServerStartUpFailed": {
			args: initValues(func(e *echo.Echo, s *http.Server) (err error) {
				return fmt.Errorf("error starting up")
			}),
		},
		"Failure_ServerShutDownFailed": {
			args: initValues(func(e *echo.Echo, s *http.Server) (err error) {
				return nil
			}),
		},
	}
}

func mockStartServer(args *args) *Patches {
	patches := ApplyMethod(reflect.TypeOf(args.e), "StartServer", func(e *echo.Echo, s *http.Server) (err error) {
		err = args.startServer(e, s)
		args.startServerCalled = true
		return err
	})
	return patches
}

func mockShutdownIfNeeded(args *args) (mockShutDown *Patches, mockStdLogger *Patches) {
	if args.shutDownFailed {
		mockShutDown = ApplyMethod(reflect.TypeOf(args.e), "Shutdown", func(e *echo.Echo, ctx context.Context) (err error) {
			return fmt.Errorf("error shutting down")
		})
		mockStdLogger = ApplyMethod(
			reflect.TypeOf(args.e.StdLogger), "Fatal",
			func(l *log.Logger, i ...interface{}) {
				args.serverShutDownCalled = true
			})
	}
	return mockShutDown, mockStdLogger
}

func startServerAndInterrupt(args args) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			log.Fatal(err)
		}
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)
		go func() {
			<-sigc
			signal.Stop(sigc)
		}()
		err = proc.Signal(os.Interrupt)
		if err != nil {
			log.Fatal("error")
		}
		time.Sleep(1 * time.Second)
	}()
	server.Start(args.e, args.cfg)
	time.Sleep(400 * time.Millisecond)
}

func waitForServerShutdownIfNeeded(args args) args {
	if args.shutDownFailed {
		time.Sleep(1000 * time.Millisecond) // Adjust time according to your needs
	}
	return args
}

func assertions(t *testing.T, args args) {
	assert.Equal(t, args.startServerCalled, true)
	//	if args.shutDownFailed {
	//		assert.Equal(t, args.serverShutDownCalled, true)
	//	}
}
