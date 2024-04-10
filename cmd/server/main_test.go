package main_test

import (
	"fmt"
	"os"
	"testing"

	main "go-template/cmd/server"
	"go-template/internal/config"
	"go-template/pkg/api"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/joho/godotenv"
	"github.com/keploy/go-sdk/v2/keploy"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const SuccessCase = "Success"

type TestArgs struct {
	setBaseEnv  bool
	patchDotEnv bool
	mockStart   bool
	apiStarted  bool
}

func initEnv(args *TestArgs) *Patches {
	err := keploy.New(keploy.Config{
		Name:           "TestSetup",
		Mode:           keploy.MODE_TEST, // change to MODE_TEST when you run in test mode
		Path:           ".",
		MuteKeployLogs: false,
		Delay:          5,
	})
	if err != nil {
		fmt.Print("error while running keploy", err)
	}
	if args != nil {
		if args.setBaseEnv {
			os.Setenv("ENVIRONMENT_NAME", "")
		}
		if args.patchDotEnv {
			loadPatches := ApplyFunc(godotenv.Load, func(...string) error {
				return nil
			})
			return loadPatches
		}
		if args.mockStart {
			apiPatches := ApplyFunc(api.Start, func(cfg *config.Configuration) (*echo.Echo, error) {
				args.apiStarted = true
				return nil, nil
			})

			return apiPatches
		}
	}

	return nil
}
func TestSetup(t *testing.T) {
	cases := map[string]struct {
		error   string
		isPanic bool
		init    func(*TestArgs) *Patches
		args    *TestArgs
	}{
		"Failure__envFileNotFound": {
			error:   "error loading port from .env",
			isPanic: true,
			init:    initEnv,
			args: &TestArgs{
				patchDotEnv: true,
			},
		},
		SuccessCase: {
			isPanic: false,
			init:    initEnv,
			args: &TestArgs{
				apiStarted: false,
				mockStart:  true,
			},
		},
	}
	defer keploy.KillProcessOnPort()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			if tt.init != nil {
				patches := tt.init(tt.args)
				if patches != nil {
					defer patches.Reset()
				}
			}
			if tt.isPanic {
				assert.PanicsWithValue(t, tt.error, main.Setup, tt.error)
			} else {
				main.Setup()
				assert.Equal(t, tt.args.apiStarted, true)
			}
		})
	}
}
