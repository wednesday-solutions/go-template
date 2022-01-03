package main_test

import (
	"fmt"
	"log"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	main "github.com/wednesday-solutions/go-template/cmd/server"
	"github.com/wednesday-solutions/go-template/internal/config"
	"github.com/wednesday-solutions/go-template/pkg/api"
	"github.com/wednesday-solutions/go-template/testutls"
)

func TestSetup(t *testing.T) {

	initEnv := func() {
		fmt.Print("initing")
		_, _, err := testutls.SetupEnvAndDB(t, testutls.Parameters{EnvFileLocation: "../../.env.local"})
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
			error:   "open .env.local: no such file or directory",
			isPanic: true,
			init:    initEnv,
		},
		"Failure_NoEnvName": {
			error:   "open .env.local: no such file or directory",
			isPanic: true,
		},
		"Success": {
			isPanic: false,
			init:    initEnv,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			if tt.init != nil {
				fmt.Print("initing")
				tt.init()
			}
			if tt.isPanic {
				assert.PanicsWithValue(t, tt.error, main.Setup, "os.Exit was not called")
			} else {
				apiStarted := false
				loadPatches := ApplyFunc(godotenv.Load, func(...string) error {
					return nil
				})
				apiPatches := ApplyFunc(api.Start, func(cfg *config.Configuration) error {
					apiStarted = true
					return nil
				})

				defer apiPatches.Reset()
				defer loadPatches.Reset()

				main.Setup()
				assert.Equal(t, apiStarted, true)
			}

		})
	}

}
