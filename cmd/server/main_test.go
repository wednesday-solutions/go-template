package main_test

import (
	"log"
	"testing"

	main "go-template/cmd/server"
	"go-template/internal/config"
	"go-template/pkg/api"
	"go-template/pkg/utl/convert"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const SuccessCase = "Success"

func TestSetup(t *testing.T) {

	initEnv := func() {

		err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../../"))
		if err != nil {
			log.Fatal(err)
		}
	}
	cases := map[string]struct {
		error   string
		isPanic bool
		init    func()
	}{
		"Failure__envFileNotFound": {
			error:   "open .env.base: no such file or directory",
			isPanic: true,
			init:    initEnv,
		},
		"Failure_NoEnvName": {
			error:   "open .env.base: no such file or directory",
			isPanic: true,
		},
		SuccessCase: {
			isPanic: false,
			init:    initEnv,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			if tt.init != nil {
				tt.init()
			}
			if tt.isPanic {
				assert.PanicsWithValue(t, tt.error, main.Setup, "os.Exit was not called")
			} else {
				apiStarted := false
				loadPatches := ApplyFunc(godotenv.Load, func(...string) error {
					return nil
				})
				apiPatches := ApplyFunc(api.Start, func(cfg *config.Configuration) (*echo.Echo, error) {
					apiStarted = true
					return nil, nil
				})

				defer apiPatches.Reset()
				defer loadPatches.Reset()

				main.Setup()
				assert.Equal(t, apiStarted, true)
			}

		})
	}

}
